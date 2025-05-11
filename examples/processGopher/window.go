// processGopher/window.go
package main

import (
	"fmt"
	"log"
	"context"

	"github.com/justyntemme/gtk4go/gtk4"
)

// ProcessWindow represents the main application window
type ProcessWindow struct {
	*gtk4.Window
	processList    *gtk4.ListView
	processModel   *gtk4.StringList
	StatusLabel    *gtk4.Label
	TotalProcLabel *gtk4.Label
}

// NewProcessWindow creates and configures the main application window
func NewProcessWindow() *ProcessWindow {
	// Create window
	window := gtk4.NewWindow(applicationTitle)
	window.SetDefaultSize(1000, 700)

	// Try hardware acceleration, but don't fail if it's not available
	// This helps with systems that don't have proper OpenGL context
	// window.EnableAcceleratedRendering()
	// window.OptimizeForResizing()

	// Create main layout and get component references
	layoutComponents := CreateProcessLayout()
	
	// Create process window struct with proper references
	procWindow := &ProcessWindow{
		Window:         window,
		processList:    layoutComponents.ProcessList,
		processModel:   layoutComponents.ProcessModel,
		StatusLabel:    layoutComponents.StatusLabel,
		TotalProcLabel: layoutComponents.TotalProcLabel,
	}
	
	// Set the main layout as window content
	window.SetChild(layoutComponents.MainBox)

	// Connect window events
	procWindow.connectSignals()

	return procWindow
}

// connectSignals connects window event handlers
func (pw *ProcessWindow) connectSignals() {
	// Handle window close request
	pw.ConnectCloseRequest(func() bool {
		// Cleanup before closing
		if workerManager != nil {
			workerManager.StopUpdates()
		}
		return false // Allow the window to close
	})

	// CSS resize optimization can cause OpenGL context issues
	// Commented out to avoid OpenGL-related crashes
	// pw.SetupCSSOptimizedResize()

	// Handle list item selection (optional - for future features)
	if pw.processList != nil {
		pw.processList.ConnectActivate(func(position int) {
			// Handle process selection if needed
			// For example, show detailed process info
		})
	}
}

// UpdateProcessList updates the process list with new data
func (pw *ProcessWindow) UpdateProcessList(processes []ProcessInfo) {
	if pw.processModel == nil {
		fmt.Println("Error: processModel is nil")
		return
	}

	// Clear existing items
	itemCount := int(pw.processModel.GetNItems())
	for i := itemCount - 1; i >= 0; i-- {
		pw.processModel.Remove(i)
	}

	// Add new items
	for i, proc := range processes {
		processText := fmt.Sprintf("%-6d %-30s %-12s %5.1f%% %8.1f",
			proc.PID,
			truncateString(proc.Name, 30),
			truncateString(proc.Username, 12),
			proc.CPUPercent,
			proc.MemoryMB,
		)
		
		pw.processModel.Append(processText)
		
		// Debug: Print the first few items
		if i < 3 {
			fmt.Printf("Added process %d: '%s'\n", i, processText)
		}
	}
	
	// Verify items were added
	newCount := int(pw.processModel.GetNItems())
	fmt.Printf("Updated process list: %d processes\n", newCount)
}

// ShowError displays an error message to the user
func (pw *ProcessWindow) ShowError(title, message string) {
	dialog := gtk4.NewMessageDialog(
		pw.Window,
		gtk4.DialogModal,
		gtk4.MessageError,
		gtk4.ResponseOk,
		message,
	)
	dialog.SetTitle(title)
	
	dialog.ConnectResponse(func(responseId gtk4.ResponseType) {
		dialog.Destroy()
	})
	
	dialog.Show()
}

// ShowAbout displays the about dialog
func (pw *ProcessWindow) ShowAbout() {
	dialog := gtk4.NewDialog("About Process Gopher", pw.Window,
		gtk4.DialogModal|gtk4.DialogDestroyWithParent)
	
	// Get content area
	content := dialog.GetContentArea()
	
	// Create about content
	aboutBox := gtk4.NewBox(gtk4.OrientationVertical, 10)
	aboutBox.AddCssClass("about-dialog")
	
	title := gtk4.NewLabel("Process Gopher v1.0")
	title.AddCssClass("about-title")
	aboutBox.Append(title)
	
	description := gtk4.NewLabel("A simple system process monitor built with Go and GTK4")
	description.AddCssClass("about-description")
	aboutBox.Append(description)
	
	author := gtk4.NewLabel("© 2024 Process Gopher Team")
	author.AddCssClass("about-author")
	aboutBox.Append(author)
	
	content.Append(aboutBox)
	
	// Add OK button
	dialog.AddButton("OK", gtk4.ResponseOk)
	
	// Connect response
	dialog.ConnectResponse(func(responseId gtk4.ResponseType) {
		dialog.Destroy()
	})
	
	dialog.Show()
}

// SetRefreshInterval allows changing the update interval
func (pw *ProcessWindow) SetRefreshInterval(intervalMS int) {
	if workerManager != nil {
		workerManager.SetUpdateInterval(intervalMS)
	}
}

// GetSelectedProcess returns the currently selected process, if any
func (pw *ProcessWindow) GetSelectedProcess() (*ProcessInfo, bool) {
	if pw.processList == nil {
		return nil, false
	}
	
	model := pw.processList.GetModel()
	if model == nil {
		return nil, false
	}
	
	// Get selected item from single selection model
	if singleSelection, ok := model.(*gtk4.SingleSelection); ok {
		selectedPos := singleSelection.GetSelected()
		if selectedPos >= 0 && selectedPos < singleSelection.GetNItems() {
			// Parse the selected text to extract process info
			text := singleSelection.GetItem(selectedPos)
			if str, ok := text.(string); ok {
				// Parse the formatted string back to ProcessInfo
				// This would require implementing a parser or storing ProcessInfo objects directly
				// For now, return nil as this is a future feature
				_ = str // Use the variable to avoid compiler warning
			}
		}
	}
	
	return nil, false
}

// RefreshProcessList manually triggers a process list refresh
func (pw *ProcessWindow) RefreshProcessList() {
	// Queue an immediate background task to refresh the process list
	if workerManager != nil && workerManager.IsRunning() {
		workerManager.worker.QueueTask(
			"manual-refresh",
			func(ctx context.Context, progress func(int, string)) (interface{}, error) {
				processes, err := GetAllProcesses()
				return processes, err
			},
			func(result interface{}, err error) {
				if err != nil {
					log.Printf("Manual refresh error: %v", err)
					updateStatusBar(fmt.Sprintf("Error: %v", err), 0)
				} else if processes, ok := result.([]ProcessInfo); ok {
					pw.UpdateProcessList(processes)
					updateStatusBar("Running", len(processes))
				}
			},
			nil, // No progress callback needed for manual refresh
		)
	}
}

// Destroy cleans up the window and its resources
func (pw *ProcessWindow) Destroy() {
	// Cancel any ongoing updates through worker manager
	if workerManager != nil {
		workerManager.StopUpdates()
	}
	
	// Disconnect all signals
	gtk4.DisconnectAll(pw.Window)
	
	// Clean up model
	if pw.processModel != nil {
		pw.processModel.Destroy()
	}
	
	// Call parent destroy
	pw.Window.Destroy()
}
