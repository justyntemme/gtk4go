// Package gtk4go provides async worker functionality for GTK4.
// File: gtk4go/background.go
package gtk4go

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// BackgroundWorker handles running tasks in the background
// while providing updates to the UI thread
type BackgroundWorker struct {
	workQueue chan *WorkItem
	stopChan  chan struct{}
	wg        sync.WaitGroup
	isRunning bool
	workerID  int
	mu        sync.Mutex
}

// WorkStatus represents the status of a background task
type WorkStatus int

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
	status      WorkStatus
	result      interface{}
	err         error
	ctx         context.Context
	cancelFunc  context.CancelFunc
	progressMu  sync.Mutex
	lastUpdate  time.Time
	updateDelay time.Duration // Minimum time between progress updates to prevent UI spam
}

// NewBackgroundWorker creates a new background worker
func NewBackgroundWorker(numWorkers int) *BackgroundWorker {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	worker := &BackgroundWorker{
		workQueue: make(chan *WorkItem, 100),
		stopChan:  make(chan struct{}),
		isRunning: true,
	}

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

	for {
		select {
		case <-w.stopChan:
			return
		case item := <-w.workQueue:
			if item == nil {
				continue
			}

			// Mark as running
			item.status = StatusRunning

			// Create a progress function that runs on UI thread
			progressFunc := func(percent int, message string) {
				item.progressMu.Lock()
				// Only update at most once per minimum delay to avoid flooding the UI
				now := time.Now()
				if now.Sub(item.lastUpdate) < item.updateDelay {
					item.progressMu.Unlock()
					return
				}
				item.lastUpdate = now
				item.progressMu.Unlock()

				if item.OnProgress != nil {
					RunOnUIThread(func() {
						item.OnProgress(percent, message)
					})
				}
			}

			// Execute the task
			result, err := item.Task(item.ctx, progressFunc)

			// Check if cancelled
			if item.ctx.Err() == context.Canceled {
				item.status = StatusCancelled
				item.err = context.Canceled
			} else if err != nil {
				item.status = StatusFailed
				item.err = err
			} else {
				item.status = StatusCompleted
				item.result = result
			}

			// Execute completion callback on UI thread
			if item.OnComplete != nil {
				RunOnUIThread(func() {
					item.OnComplete(item.result, item.err)
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
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isRunning {
		if onComplete != nil {
			RunOnUIThread(func() {
				onComplete(nil, fmt.Errorf("worker is not running"))
			})
		}
		return func() {}
	}

	// Create cancellable context
	ctx, cancelFunc := context.WithCancel(context.Background())

	item := &WorkItem{
		ID:          id,
		Task:        task,
		OnProgress:  onProgress,
		OnComplete:  onComplete,
		status:      StatusPending,
		ctx:         ctx,
		cancelFunc:  cancelFunc,
		lastUpdate:  time.Now(),
		updateDelay: 100 * time.Millisecond, // Update UI at most every 100ms
	}

	// Queue the work
	select {
	case w.workQueue <- item:
		// Successfully queued
	default:
		// Queue is full, execute completion with error
		if onComplete != nil {
			RunOnUIThread(func() {
				onComplete(nil, fmt.Errorf("work queue is full"))
			})
		}
	}

	return cancelFunc
}

// Stop stops the worker and waits for all tasks to complete
func (w *BackgroundWorker) Stop() {
	w.mu.Lock()
	if !w.isRunning {
		w.mu.Unlock()
		return
	}
	w.isRunning = false
	close(w.stopChan)
	w.mu.Unlock()

	// Wait for all workers to finish
	w.wg.Wait()
}

// DefaultWorker is the default background worker
var DefaultWorker = NewBackgroundWorker(4) // Create with 4 worker goroutines

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

// init ensures we clean up the default worker when the program exits
func init() {
	// This doesn't actually get called since Go doesn't have a clean shutdown hook,
	// but keeping it here as a reminder that in a real app, you'd want to call
	// DefaultWorker.Stop() at program exit
}
