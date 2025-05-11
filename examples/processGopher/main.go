// processGopher/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/justyntemme/gtk4go"
	"github.com/justyntemme/gtk4go/gtk4"
)

const (
	applicationID    = "com.example.processGopher"
	updateInterval   = 2000 // milliseconds
	applicationTitle = "Process Gopher - System Process Monitor"
)

var (
	procWindow    *ProcessWindow
	workerManager *WorkerManager
)

func main() {
	// Initialize GTK
	if err := gtk4go.Initialize(); err != nil {
		log.Fatalf("Failed to initialize GTK: %v", err)
		os.Exit(1)
	}

	// Create application
	app := gtk4.NewApplication(applicationID)

	// Create and configure the main window
	procWindow = NewProcessWindow()
	
	// Add window to application
	app.AddWindow(procWindow.Window)

	// Create worker manager
	workerManager = NewWorkerManager()

	// Setup periodic updates
	workerManager.StartPeriodicUpdates(procWindow)

	// Run the application
	exitCode := app.Run()
	
	// Cleanup
	workerManager.Shutdown()
	
	os.Exit(exitCode)
}

// updateStatusBar updates the status information
func updateStatusBar(status string, processCount int) {
	if procWindow != nil && procWindow.StatusLabel != nil {
		procWindow.StatusLabel.SetText(fmt.Sprintf("Status: %s | Last Update: %s", 
			status, time.Now().Format("15:04:05")))
	}
	
	if procWindow != nil && procWindow.TotalProcLabel != nil {
		procWindow.TotalProcLabel.SetText(fmt.Sprintf("Total Processes: %d", processCount))
	}
}
