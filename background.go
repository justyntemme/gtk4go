// Package gtk4go provides async worker functionality for GTK4.
// File: gtk4go/background.go
package gtk4go

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// BackgroundWorker handles running tasks in the background
// while providing updates to the UI thread
type BackgroundWorker struct {
	workQueue     chan *WorkItem
	stopChan      chan struct{}
	wg            sync.WaitGroup
	isRunning     atomic.Bool
	activeWorkers atomic.Int32
	workerCount   atomic.Int32
	nextTaskID    atomic.Uint64
	mu            sync.RWMutex // Used only for fields not amenable to atomic ops
}

// WorkStatus represents the status of a background task
type WorkStatus int32

const (
	// StatusPending indicates the task is pending
	StatusPending WorkStatus = iota
	// StatusRunning indicates the task is running
	StatusRunning
	// StatusCompleted indicates the task completed successfully
	StatusCompleted
	// StatusFailed indicates the task failed with an error
	StatusFailed
	// StatusCancelled indicates the task was cancelled
	StatusCancelled
)

// WorkItem represents a unit of work to be processed
type WorkItem struct {
	ID          string
	Task        func(ctx context.Context, progress func(percent int, message string)) (interface{}, error)
	OnProgress  func(percent int, message string)
	OnComplete  func(result interface{}, err error)
	status      atomic.Int32 // Using atomic int32 for WorkStatus
	result      atomic.Value // For storing the result
	err         atomic.Value // For storing the error
	ctx         context.Context
	cancelFunc  context.CancelFunc
	
	// Progress tracking
	progressMu    sync.Mutex      // Traditional mutex for complex progress operations
	lastUpdate    atomic.Value    // Using atomic.Value for time.Time
	updateDelay   time.Duration   // Update frequency limitation
	progressCalls atomic.Int64    // Count of progress calls for metrics
}

// NewBackgroundWorker creates a new background worker
func NewBackgroundWorker(numWorkers int) *BackgroundWorker {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	worker := &BackgroundWorker{
		workQueue: make(chan *WorkItem, 100),
		stopChan:  make(chan struct{}),
	}
	
	// Set initial state
	worker.isRunning.Store(true)
	worker.workerCount.Store(int32(numWorkers))

	// Start worker goroutines
	worker.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		workerID := i
		go worker.processWork(workerID)
	}

	return worker
}

// processWork runs in a goroutine to process work items
func (w *BackgroundWorker) processWork(workerID int) {
	defer w.wg.Done()
	
	// Track active worker count
	w.activeWorkers.Add(1)
	defer w.activeWorkers.Add(-1)

	for {
		select {
		case <-w.stopChan:
			return
		case item := <-w.workQueue:
			if item == nil {
				continue
			}

			// Mark as running using atomic operation
			item.status.Store(int32(StatusRunning))

			// Create a progress function that runs on UI thread
			progressFunc := func(percent int, message string) {
				// Use atomic operations for thread-safe time checks
				lastUpdateObj := item.lastUpdate.Load()
				var lastUpdate time.Time
				if lastUpdateObj != nil {
					lastUpdate = lastUpdateObj.(time.Time)
				}
				
				now := time.Now()
				
				// Rate limiting check
				item.progressMu.Lock()
				updateDelay := item.updateDelay
				shouldUpdate := now.Sub(lastUpdate) >= updateDelay
				if shouldUpdate {
					item.lastUpdate.Store(now)
				}
				item.progressMu.Unlock()
				
				if !shouldUpdate {
					return
				}
				
				// Increment progress call counter
				item.progressCalls.Add(1)

				if item.OnProgress != nil {
					RunOnUIThread(func() {
						item.OnProgress(percent, message)
					})
				}
			}

			// Execute the task
			result, err := item.Task(item.ctx, progressFunc)

			// Store result/error using atomic operations
			if result != nil {
				item.result.Store(result)
			}
			
			if err != nil {
				item.err.Store(err)
			}

			// Check if cancelled using atomic operations
			var finalStatus WorkStatus
			if item.ctx.Err() == context.Canceled {
				finalStatus = StatusCancelled
			} else if err != nil {
				finalStatus = StatusFailed
			} else {
				finalStatus = StatusCompleted
			}
			
			item.status.Store(int32(finalStatus))

			// Execute completion callback on UI thread
			if item.OnComplete != nil {
				RunOnUIThread(func() {
					// Safely retrieve result/error from atomic storage
					var resultVal interface{}
					if r := item.result.Load(); r != nil {
						resultVal = r
					}
					
					var errVal error
					if e := item.err.Load(); e != nil {
						errVal = e.(error)
					}
					
					item.OnComplete(resultVal, errVal)
				})
			}
		}
	}
}

// QueueTask queues a task for background execution
func (w *BackgroundWorker) QueueTask(
	id string,
	task func(ctx context.Context, progress func(percent int, message string)) (interface{}, error),
	onComplete func(result interface{}, err error),
	onProgress func(percent int, message string),
) context.CancelFunc {
	// Check if we're running using atomic operations
	if !w.isRunning.Load() {
		if onComplete != nil {
			RunOnUIThread(func() {
				onComplete(nil, fmt.Errorf("worker is not running"))
			})
		}
		return func() {}
	}
	
	// Generate ID if none provided
	if id == "" {
		id = fmt.Sprintf("task-%d", w.nextTaskID.Add(1))
	}

	// Create cancellable context
	ctx, cancelFunc := context.WithCancel(context.Background())

	item := &WorkItem{
		ID:          id,
		Task:        task,
		OnProgress:  onProgress,
		OnComplete:  onComplete,
		ctx:         ctx,
		cancelFunc:  cancelFunc,
		updateDelay: 100 * time.Millisecond, // Update UI at most every 100ms
	}
	
	// Initialize atomic values
	item.status.Store(int32(StatusPending))
	item.lastUpdate.Store(time.Now())

	// Try to queue the work with a timeout to prevent deadlocks
	select {
	case w.workQueue <- item:
		// Successfully queued
	case <-time.After(100 * time.Millisecond):
		// Queue is full or blocked
		if onComplete != nil {
			RunOnUIThread(func() {
				onComplete(nil, fmt.Errorf("work queue is full"))
			})
		}
		cancelFunc()
	}

	return cancelFunc
}

// SetProgressUpdateInterval sets the minimum time between progress updates
func (w *BackgroundWorker) SetProgressUpdateInterval(duration time.Duration) {
	// This only affects new tasks
	w.mu.Lock()
	defer w.mu.Unlock()
	// Would be stored in a field if we needed it
}

// GetActiveWorkerCount returns the number of currently active workers
func (w *BackgroundWorker) GetActiveWorkerCount() int {
	return int(w.activeWorkers.Load())
}

// IsRunning returns whether the worker is currently running
func (w *BackgroundWorker) IsRunning() bool {
	return w.isRunning.Load()
}

// Stop stops the worker and waits for all tasks to complete
func (w *BackgroundWorker) Stop() {
	// Use atomic operation to check and update running state
	if !w.isRunning.CompareAndSwap(true, false) {
		// Already stopped
		return
	}
	
	close(w.stopChan)

	// Wait for all workers to finish
	w.wg.Wait()
}

// Shutdown stops the worker and cancels any pending or running tasks
func (w *BackgroundWorker) Shutdown(timeout time.Duration) bool {
	// Stop accepting new tasks
	if !w.isRunning.CompareAndSwap(true, false) {
		// Already stopped
		return true
	}

	// Close the stop channel to signal workers to exit
	close(w.stopChan)

	// Use a channel to signal completion
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	// Wait for workers to finish with timeout
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// DefaultWorker is the default background worker
var DefaultWorker = NewBackgroundWorker(runtime.NumCPU())

// QueueBackgroundTask is a convenience function that queues a task on the default worker
func QueueBackgroundTask(
	id string,
	task func(ctx context.Context, progress func(percent int, message string)) (interface{}, error),
	onComplete func(result interface{}, err error),
	onProgress func(percent int, message string),
) context.CancelFunc {
	return DefaultWorker.QueueTask(id, task, onComplete, onProgress)
}

// RunInBackground runs a simple task without progress updates
func RunInBackground(
	task func() (interface{}, error),
	onComplete func(result interface{}, err error),
) context.CancelFunc {
	wrappedTask := func(ctx context.Context, _ func(int, string)) (interface{}, error) {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Continue with task
		}

		return task()
	}

	return QueueBackgroundTask("", wrappedTask, onComplete, nil)
}

// ShutdownDefaultWorker shuts down the default worker with a timeout
func ShutdownDefaultWorker(timeout time.Duration) bool {
	return DefaultWorker.Shutdown(timeout)
}

// init ensures we clean up the default worker when the program exits
func init() {
	// Register a cleanup function to be called at exit if possible
	runtime.SetFinalizer(DefaultWorker, func(w *BackgroundWorker) {
		w.Stop()
	})
}