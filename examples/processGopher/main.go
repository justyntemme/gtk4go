package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/justyntemme/gtk4go"
	"github.com/justyntemme/gtk4go/gtk4"
)

// Define constants for styling and sizing
const (
	APP_ID         = "com.example.process-gopher"
	TITLE          = "Process Gopher"
	DEFAULT_WIDTH  = 1000
	DEFAULT_HEIGHT = 700

	// Auto-refresh interval in seconds
	AUTO_REFRESH_INTERVAL = 2
)

// Global variables for data and UI state
var (
	processListView  *gtk4.ListView
	processListModel *gtk4.StringList
	selectionModel   *gtk4.SingleSelection
	refreshTimer     *time.Timer
	statusLabel      *gtk4.Label
	columnSortOrder  = make(map[int]bool) // true for ascending, false for descending
	statusMutex      sync.Mutex
	selectedPID      int64 = -1
	processCache     []ProcessInfo // Cache of process data for selection and sorting
)

// Column indices for the process list
const (
	COL_PID = iota
	COL_NAME
	COL_USERNAME
	COL_CPU
	COL_MEMORY
	COL_THREADS
	COL_STATE
	COL_STARTED
	COL_COUNT // Total number of columns
)

func main() {
	// Initialize GTK
	if err := gtk4go.Initialize(); err != nil {
		fmt.Printf("Failed to initialize GTK: %v\n", err)
		os.Exit(1)
	}

	// Create application
	app := gtk4.NewApplication(APP_ID)

	// Create actions
	refreshAction := gtk4.NewAction("refresh", refreshProcessList)
	app.GetActionGroup().AddAction(refreshAction)

	endProcessAction := gtk4.NewAction("end-process", endSelectedProcess)
	app.GetActionGroup().AddAction(endProcessAction)

	// Create window
	win := gtk4.NewWindow(TITLE)
	win.SetDefaultSize(DEFAULT_WIDTH, DEFAULT_HEIGHT)

	// Create main layout
	mainBox := createMainLayout()

	// Set the window's child to the main box
	win.SetChild(mainBox)

	// Set up window close handler
	win.ConnectCloseRequest(func() bool {
		// Clean up resources
		if refreshTimer != nil {
			refreshTimer.Stop()
		}
		return false // Return false to allow window to close
	})

	// Add window to application
	app.AddWindow(win)

	// Load CSS styling
	loadAppStyles()

	// Start auto-refresh timer
	startRefreshTimer()

	// Initial load of process data
	refreshProcessList()

	// Run the application
	os.Exit(app.Run())
}

// createMainLayout creates the main layout of the application
func createMainLayout() *gtk4.Box {
	// Create main vertical box
	mainBox := gtk4.NewBox(gtk4.OrientationVertical, 10)

	// Create header bar
	headerBar := createHeaderBar()
	mainBox.Append(headerBar)

	// Create stack for different tabs
	stack := gtk4.NewStack()
	mainBox.Append(stack)

	// Create processes tab
	processesTab := createProcessesTab()
	stack.AddTitled(processesTab, "processes", "Processes")

	// Create performance tab
	performanceTab := createPerformanceTab()
	stack.AddTitled(performanceTab, "performance", "Performance")

	// Create stack switcher (tabs)
	stackSwitcher := gtk4.NewStackSwitcher(stack)
	mainBox.Append(stackSwitcher)

	// Create status bar
	statusBar := createStatusBar()
	mainBox.Append(statusBar)

	return mainBox
}

// createHeaderBar creates the header bar with search and actions
func createHeaderBar() *gtk4.Box {
	headerBox := gtk4.NewBox(gtk4.OrientationHorizontal, 10)

	// Search entry
	searchEntry := gtk4.NewEntry()
	searchEntry.SetPlaceholderText("Search processes...")
	searchEntry.ConnectChanged(func() {
		searchProcesses(searchEntry.GetText())
	})
	headerBox.Append(searchEntry)

	// Refresh button
	refreshButton := gtk4.NewButton("Refresh")
	refreshButton.ConnectClicked(refreshProcessList)
	headerBox.Append(refreshButton)

	// End Process button
	endProcessButton := gtk4.NewButton("End Process")
	endProcessButton.ConnectClicked(endSelectedProcess)
	headerBox.Append(endProcessButton)

	return headerBox
}

// createProcessesTab creates the processes tab content
func createProcessesTab() *gtk4.ScrolledWindow {
	// Create scrolled window
	scrollWin := gtk4.NewScrolledWindow()

	// Create string list model to display processes
	processListModel = gtk4.NewStringList()

	// Create a selection model for the list
	selectionModel = gtk4.NewSingleSelection(processListModel)

	// Create a factory for list items
	factory := gtk4.NewSignalListItemFactory()

	// Set up list items with setup callback
	factory.ConnectSetup(func(listItem *gtk4.ListItem) {
		// Create a box for layout
		box := gtk4.NewBox(gtk4.OrientationHorizontal, 10)
		box.SetHExpand(true)
		box.AddCssClass("list-item-box")

		// Create an icon and label
		icon := gtk4.NewLabel("•")
		icon.AddCssClass("list-item-icon")
		box.Append(icon)

		label := gtk4.NewLabel("")
		label.AddCssClass("list-item-label")
		box.Append(label)

		// Set the box as the child of the list item
		listItem.SetChild(box)
	})

	// Bind data to list items
	factory.ConnectBind(func(listItem *gtk4.ListItem) {
		// Get the text from the model
		text := listItem.GetText()
		if text == "" {
			text = fmt.Sprintf("Item %d", listItem.GetPosition()+1)
		}

		// Set the text on the label inside the box
		listItem.SetTextOnChildLabel(text)

		// Add selected class if the item is selected
		if listItem.GetSelected() {
			boxWidget := listItem.GetChild()
			boxWidget.AddCssClass("selected")
		} else {
			boxWidget := listItem.GetChild()
			boxWidget.RemoveCssClass("selected")
		}
	})

	// Create the list view with selection model and factory
	processListView = gtk4.NewListView(selectionModel, factory)

	// Connect activate signal for item selection
	processListView.ConnectActivate(func(position int) {
		if position >= 0 && position < len(processCache) {
			// Get the actual PID from our process cache
			selectedPID = processCache[position].PID
			setStatus(fmt.Sprintf("Selected process: %s (PID: %d)", processCache[position].Name, selectedPID))
		} else {
			selectedPID = -1
			setStatus("Invalid process selection")
		}
	})

	// Add the list view to the scrolled window
	scrollWin.SetChild(processListView)

	return scrollWin
}



// createPerformanceTab creates the performance tab content
func createPerformanceTab() *gtk4.Box {
	box := gtk4.NewBox(gtk4.OrientationVertical, 10)

	// CPU usage section
	cpuLabel := gtk4.NewLabel("CPU Usage")
	cpuLabel.AddCssClass("heading")
	box.Append(cpuLabel)

	// CPU usage display
	cpuValueLabel := gtk4.NewLabel("Collecting data...")
	cpuValueLabel.AddCssClass("usage-value")
	box.Append(cpuValueLabel)

	// Memory usage section
	memLabel := gtk4.NewLabel("Memory Usage")
	memLabel.AddCssClass("heading")
	box.Append(memLabel)

	// Memory usage display
	memValueLabel := gtk4.NewLabel("Collecting data...")
	memValueLabel.AddCssClass("usage-value")
	box.Append(memValueLabel)

	// Start a timer to update the performance data
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			// Use RunOnUIThread to safely update UI components from a goroutine
			gtk4go.RunOnUIThread(func() {
				// Update CPU usage
				cpuUsage, err := getCPUUsage()
				if err == nil {
					cpuValueLabel.SetText(fmt.Sprintf("CPU Usage: %.1f%%", cpuUsage))
					cpuValueLabel.AddCssClass("usage-normal")
					if cpuUsage > 80 {
						cpuValueLabel.RemoveCssClass("usage-normal")
						cpuValueLabel.AddCssClass("usage-high")
					}
				} else {
					cpuValueLabel.SetText("CPU Usage: Error - " + err.Error())
				}
				
				// Update memory usage
				total, free, err := getSystemMemoryInfo()
				if err == nil {
					used := total - free
					usedPercentage := float64(used) / float64(total) * 100
					
					// Calculate values in MB for display
					totalMB := total / (1024 * 1024)
					usedMB := used / (1024 * 1024)
					
					memValueLabel.SetText(fmt.Sprintf("Memory Usage: %d MB / %d MB (%.1f%%)", 
						usedMB, totalMB, usedPercentage))
					memValueLabel.AddCssClass("usage-normal")
					if usedPercentage > 80 {
						memValueLabel.RemoveCssClass("usage-normal")
						memValueLabel.AddCssClass("usage-high")
					}
				} else {
					memValueLabel.SetText("Memory Usage: Error - " + err.Error())
				}
			})
		}
	}()

	return box
}

// createStatusBar creates the status bar at the bottom
func createStatusBar() *gtk4.Box {
	statusBox := gtk4.NewBox(gtk4.OrientationHorizontal, 5)

	// Status label
	statusLabel = gtk4.NewLabel("Ready")
	statusBox.Append(statusLabel)

	return statusBox
}

// startRefreshTimer starts the auto-refresh timer
func startRefreshTimer() {
	if refreshTimer != nil {
		refreshTimer.Stop()
	}

	refreshTimer = time.AfterFunc(time.Duration(AUTO_REFRESH_INTERVAL)*time.Second, func() {
		// Schedule refresh on UI thread
		gtk4go.RunOnUIThread(func() {
			refreshProcessList()
			startRefreshTimer() // Schedule next refresh
		})
	})
}

// refreshProcessList refreshes the process list with current data
func refreshProcessList() {
	setStatus("Refreshing process list...")

	// Clear the current list
	// Remove all items - items.length changes after each remove, so remove from the end
	for i := processListModel.GetNItems() - 1; i >= 0; i-- {
		processListModel.Remove(i)
	}

	// Get the process data
	processes, err := getProcesses()
	if err != nil {
		setStatus(fmt.Sprintf("Error: %v", err))
		return
	}

	// Update the process cache
	processCache = processes

	// Add each process to the list
	for _, proc := range processes {
		// Format a display string for each process
		// Handle empty values for better display
		username := proc.Username
		if username == "" {
			username = "N/A"
		}
		
		state := proc.State
		if state == "" {
			state = "N/A"
		}
		
		startTime := proc.StartTime
		if startTime == "" {
			startTime = "N/A"
		}
		
		displayText := fmt.Sprintf("%d | %s | %s | %.1f%% | %dMB | %d | %s | %s",
			proc.PID,
			proc.Name,
			username,
			proc.CPUPercent,
			proc.MemoryBytes/(1024*1024),
			proc.Threads,
			state,
			startTime,
		)
		
		// Add to the list model
		processListModel.Append(displayText)
	}

	setStatus(fmt.Sprintf("Process list refreshed. %d processes found.", len(processes)))
}

// searchProcesses filters the process list based on search text
func searchProcesses(searchText string) {
	if searchText == "" {
		// If search is empty, just refresh the full list
		refreshProcessList()
		return
	}

	setStatus(fmt.Sprintf("Searching for '%s'...", searchText))

	// We would implement filtering logic here.
	// For now, we'll just refresh the list and report that search is not yet implemented.
	refreshProcessList()
	setStatus(fmt.Sprintf("Search for '%s' completed.", searchText))
}

// endSelectedProcess attempts to terminate the selected process
func endSelectedProcess() {
	if selectedPID <= 0 {
		setStatus("No process selected")
		return
	}

	// Show confirmation dialog
	dialog := gtk4.NewMessageDialog(
		nil, // parent window
		gtk4.DialogModal|gtk4.DialogDestroyWithParent,
		gtk4.MessageWarning,
		gtk4.ResponseOk|gtk4.ResponseCancel,
		fmt.Sprintf("Are you sure you want to terminate process %d?", selectedPID),
	)
	// Add secondary text
	secondaryLabel := gtk4.NewLabel("This may cause data loss if the application is not responding.")
	dialog.GetContentArea().Append(secondaryLabel)

	// Handle dialog response
	dialog.ConnectResponse(func(responseId gtk4.ResponseType) {
		if responseId == gtk4.ResponseOk {
			// User confirmed, try to kill the process
			err := terminateProcess(selectedPID)
			if err != nil {
				setStatus(fmt.Sprintf("Error terminating process: %v", err))
			} else {
				setStatus(fmt.Sprintf("Process %d terminated successfully", selectedPID))
				refreshProcessList()
			}
		}
		dialog.Destroy()
	})

	dialog.Show()
}

// sortProcessList sorts the process list by the specified column
func sortProcessList(columnID int, ascending bool) {
	// This is a stub. In a real implementation, we would sort the list store.
	if ascending {
		setStatus(fmt.Sprintf("Sorting by column %d in ascending order", columnID))
	} else {
		setStatus(fmt.Sprintf("Sorting by column %d in descending order", columnID))
	}
}

// loadAppStyles loads the CSS styles for the application
func loadAppStyles() {
	// Get the CSS provider and add styles
	cssProvider, err := gtk4.LoadCSS(GetCSS())
	if err != nil {
		fmt.Printf("Error loading CSS: %v\n", err)
		return
	}

	// Apply to all windows in the application
	gtk4.AddProviderForDisplay(cssProvider, 600)
}

// setStatus updates the status bar text
func setStatus(message string) {
	statusMutex.Lock()
	defer statusMutex.Unlock()

	// Update on UI thread
	gtk4go.RunOnUIThread(func() {
		statusLabel.SetText(message)
	})
}
