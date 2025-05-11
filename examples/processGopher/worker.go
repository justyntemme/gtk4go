// processGopher/worker.go
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/justyntemme/gtk4go"
)

// WorkerManager manages background tasks for the application
type WorkerManager struct {
	worker         *gtk4go.BackgroundWorker
	updateCancel   context.CancelFunc
	updateInterval time.Duration
	mu             sync.RWMutex
}

// NewWorkerManager creates a new worker manager
func NewWorkerManager() *WorkerManager {
	// Create a dedicated background worker for process monitoring
	worker := gtk4go.NewBackgroundWorker(2) // Use 2 workers for parallel tasks
	
	return &WorkerManager{
		worker:         worker,
		updateInterval: time.Duration(updateInterval) * time.Millisecond,
	}
}

// StartPeriodicUpdates begins the automatic process list updates
func (wm *WorkerManager) StartPeriodicUpdates(window *ProcessWindow) {
	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	
	wm.mu.Lock()
	wm.updateCancel = cancel
	wm.mu.Unlock()

	// Initial load with progress reporting
	wm.worker.QueueTask(
		"initial-process-load",
		func(taskCtx context.Context, progress func(int, string)) (interface{}, error) {
			processes, err := GetAllProcessesWithProgress(func(current, total int, message string) {
				percent := int(float64(current) / float64(total) * 100)
				progress(percent, message)
			})
			return processes, err
		},
		func(result interface{}, err error) {
			if err != nil {
				log.Printf("Initial process load error: %v", err)
				updateStatusBar(fmt.Sprintf("Error: %v", err), 0)
			} else if processes, ok := result.([]ProcessInfo); ok {
				window.UpdateProcessList(processes)
				updateStatusBar("Running", len(processes))
			}
		},
		func(percent int, message string) {
			updateStatusBar(fmt.Sprintf("Loading: %s", message), 0)
		},
	)

	// Start periodic updates
	wm.worker.QueueTask(
		"process-updater",
		func(taskCtx context.Context, progress func(int, string)) (interface{}, error) {
			ticker := time.NewTicker(wm.updateInterval)
			defer ticker.Stop()

			updateCount := 0
			for {
				select {
				case <-taskCtx.Done():
					return fmt.Sprintf("Updates stopped after %d iterations", updateCount), nil
				case <-ctx.Done():
					return fmt.Sprintf("Application shutting down after %d iterations", updateCount), nil
				case <-ticker.C:
					updateCount++
					
					// Report progress
					if updateCount%10 == 0 {
						progress(updateCount, fmt.Sprintf("Completed %d updates", updateCount))
					}
					
					// Fetch processes in background thread
					processes, err := GetAllProcesses()
					if err != nil {
						log.Printf("Error getting processes (update %d): %v", updateCount, err)
						// Update UI with error status on UI thread
						gtk4go.RunOnUIThread(func() {
							updateStatusBar(fmt.Sprintf("Error: %v", err), 0)
						})
						continue
					}

					// Update UI on UI thread
					gtk4go.RunOnUIThread(func() {
						window.UpdateProcessList(processes)
						updateStatusBar("Running", len(processes))
					})
				}
			}
		},
		func(result interface{}, err error) {
			if err != nil {
				log.Printf("Update task error: %v", err)
			} else if resultStr, ok := result.(string); ok {
				log.Printf("Update task completed: %s", resultStr)
			}
		},
		func(percent int, message string) {
			// Could be used for update status indication
			log.Printf("Update progress: %s", message)
		},
	)
}

// StopUpdates stops the periodic updates
func (wm *WorkerManager) StopUpdates() {
	wm.mu.Lock()
	if wm.updateCancel != nil {
		wm.updateCancel()
		wm.updateCancel = nil
	}
	wm.mu.Unlock()
}

// SetUpdateInterval changes the update interval
func (wm *WorkerManager) SetUpdateInterval(intervalMS int) {
	wm.mu.Lock()
	wm.updateInterval = time.Duration(intervalMS) * time.Millisecond
	wm.mu.Unlock()
	
	// Restart updates with new interval
	window := procWindow
	wm.StopUpdates()
	wm.StartPeriodicUpdates(window)
}

// Shutdown gracefully shuts down the worker
func (wm *WorkerManager) Shutdown() {
	wm.StopUpdates()
	
	// Allow up to 5 seconds for tasks to complete
	success := wm.worker.Shutdown(5 * time.Second)
	if !success {
		log.Printf("Warning: Not all background tasks completed within timeout")
	}
}

// GetActiveWorkerCount returns the number of active workers
func (wm *WorkerManager) GetActiveWorkerCount() int {
	return wm.worker.GetActiveWorkerCount()
}

// IsRunning returns whether the worker is currently running
func (wm *WorkerManager) IsRunning() bool {
	return wm.worker.IsRunning()
}
