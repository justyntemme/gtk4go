package main

import (
	"fmt"
	"os"
	"strings"
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
	refreshTimer    *time.Timer
	statusLabel     *gtk4.Label
	processList     *gtk4.Box
	processCountLabel *gtk4.Label
	statusMutex     sync.Mutex
	selectedPID     int64 = -1
	processBoxes    = make(map[int64]*gtk4.Box) // Maps PIDs to their box widgets
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
	mainBox := gtk4.NewBox(gtk4.OrientationVertical, 0)
	mainBox.SetMarginTop(10)
	mainBox.SetMarginBottom(10)
	mainBox.SetMarginStart(10)
	mainBox.SetMarginEnd(10)

	// Create header bar
	headerBar := createHeaderBar()
	mainBox.Append(headerBar)

	// Create notebook for different tabs
	notebook := gtk4.NewNotebook()
	mainBox.Append(notebook)
	notebook.SetHExpand(true)
	notebook.SetVExpand(true)

	// Create processes tab
	processesTab := createProcessesTab()
	notebook.AppendPage(processesTab, gtk4.NewLabel("Processes"))

	// Create performance tab
	performanceTab := createPerformanceTab()
	notebook.AppendPage(performanceTab, gtk4.NewLabel("Performance"))

	// Create status bar
	statusBar := createStatusBar()
	mainBox.Append(statusBar)

	return mainBox
}

// createHeaderBar creates the header bar with search and actions
func createHeaderBar() *gtk4.Box {
	headerBox := gtk4.NewBox(gtk4.OrientationHorizontal, 10)
	headerBox.SetMarginTop(5)
	headerBox.SetMarginBottom(5)
	headerBox.SetMarginStart(5)
	headerBox.SetMarginEnd(5)

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
	scrollWin := gtk4.NewScrolledWindow(
		gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyNever),
		gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
		gtk4.WithHExpand(true),
		gtk4.WithVExpand(true)
	)

	// Create a vertical box to hold the process list
	contentBox := gtk4.NewBox(gtk4.OrientationVertical, 10)
	contentBox.SetMarginTop(10)
	contentBox.SetMarginBottom(10)
	contentBox.SetMarginStart(10)
	contentBox.SetMarginEnd(10)

	// Create a header for the process list
	headerBox := gtk4.NewBox(gtk4.OrientationHorizontal, 5)
	headerBox.AddCssClass("process-header")

	// Add column headers
	pidLabel := gtk4.NewLabel("PID")
	pidLabel.SetHExpand(false)
	pidLabel.SetMarginEnd(20)
	pidLabel.AddCssClass("column-header")
	headerBox.Append(pidLabel)

	nameLabel := gtk4.NewLabel("Process Name")
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nameLabel.AddCssClass("column-header")
	headerBox.Append(nameLabel)

	userLabel := gtk4.NewLabel("User")
	userLabel.SetHExpand(false)
	userLabel.SetMarginStart(20)
	userLabel.SetMarginEnd(20)
	userLabel.AddCssClass("column-header")
	headerBox.Append(userLabel)

	cpuLabel := gtk4.NewLabel("CPU %")
	cpuLabel.SetHExpand(false)
	cpuLabel.SetMarginStart(20)
	cpuLabel.SetMarginEnd(20)
	cpuLabel.AddCssClass("column-header")
	headerBox.Append(cpuLabel)

	memLabel := gtk4.NewLabel("Memory")
	memLabel.SetHExpand(false)
	memLabel.SetMarginStart(20)
	memLabel.SetMarginEnd(20)
	memLabel.AddCssClass("column-header")
	headerBox.Append(memLabel)

	contentBox.Append(headerBox)

	// Add a separator
	separator := gtk4.NewSeparator(gtk4.OrientationHorizontal)
	contentBox.Append(separator)

	// Create a box to hold the process items
	processList = gtk4.NewBox(gtk4.OrientationVertical, 0)
	contentBox.Append(processList)

	// Add a label for process count
	countBox := gtk4.NewBox(gtk4.OrientationHorizontal, 5)
	countBox.SetMarginTop(10)
	processCountLabel = gtk4.NewLabel("0 processes")
	countBox.Append(processCountLabel)
	contentBox.Append(countBox)

	// Set the content box as the child of the scrolled window
	scrollWin.SetChild(contentBox)

	return scrollWin
}

// createProcessItem creates a box representing a process item
func createProcessItem(proc ProcessInfo) *gtk4.Box {
	itemBox := gtk4.NewBox(gtk4.OrientationHorizontal, 5)
	itemBox.AddCssClass("process-item")
	
	// We need to store the PID in a way we can retrieve it
	// Since there's no direct SetData method, we'll use a map
	processBoxes[proc.PID] = itemBox

	// We'll use a GestureClick for selection
	clickGesture := gtk4.NewGestureClick()
	itemBox.AddController(clickGesture)
	clickGesture.ConnectReleased(func(n int, x, y float64) {
		// Deselect any previously selected process
		if selectedPID > 0 {
			if prevBox, ok := processBoxes[selectedPID]; ok {
				prevBox.RemoveCssClass("process-selected")
			}
		}
		
		// Update selection and add selected class
		selectedPID = proc.PID
		itemBox.AddCssClass("process-selected")
	})

	// PID column
	pidLabel := gtk4.NewLabel(fmt.Sprintf("%d", proc.PID))
	pidLabel.SetHExpand(false)
	pidLabel.SetMarginEnd(20)
	itemBox.Append(pidLabel)

	// Process name column
	nameLabel := gtk4.NewLabel(proc.Name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	itemBox.Append(nameLabel)

	// Username column
	userLabel := gtk4.NewLabel(proc.Username)
	userLabel.SetHExpand(false)
	userLabel.SetMarginStart(20)
	userLabel.SetMarginEnd(20)
	itemBox.Append(userLabel)

	// CPU column
	cpuLabel := gtk4.NewLabel(fmt.Sprintf("%.1f%%", proc.CPUPercent))
	cpuLabel.SetHExpand(false)
	cpuLabel.SetMarginStart(20)
	cpuLabel.SetMarginEnd(20)
	if proc.CPUPercent > 50 {
		cpuLabel.AddCssClass("cpu-high")
	}
	itemBox.Append(cpuLabel)

	// Memory column
	memLabel := gtk4.NewLabel(formatMemory(proc.MemoryBytes))
	memLabel.SetHExpand(false)
	memLabel.SetMarginStart(20)
	memLabel.SetMarginEnd(20)
	itemBox.Append(memLabel)

	return itemBox
}

// formatMemory formats memory size in a human-readable way
func formatMemory(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// createPerformanceTab creates the performance tab content
func createPerformanceTab() *gtk4.Box {
	box := gtk4.NewBox(gtk4.OrientationVertical, 10)
	box.SetMarginTop(10)
	box.SetMarginBottom(10)
	box.SetMarginStart(10)
	box.SetMarginEnd(10)

	// CPU usage section
	cpuLabel := gtk4.NewLabel("CPU Usage")
	cpuLabel.SetHAlign(gtk4.AlignStart)
	cpuLabel.SetMarginTop(10)
	cpuLabel.AddCssClass("heading")
	box.Append(cpuLabel)

	// CPU usage would go here (placeholder)
	cpuPlaceholder := gtk4.NewLabel("CPU usage graphs will appear here")
	cpuPlaceholder.SetHAlign(gtk4.AlignCenter)
	cpuPlaceholder.SetVAlign(gtk4.AlignCenter)
	cpuPlaceholder.SetMarginTop(20)
	cpuPlaceholder.SetMarginBottom(20)
	cpuPlaceholder.SetMarginStart(20)
	cpuPlaceholder.SetMarginEnd(20)
	box.Append(cpuPlaceholder)

	// Memory usage section
	memLabel := gtk4.NewLabel("Memory Usage")
	memLabel.SetHAlign(gtk4.AlignStart)
	memLabel.SetMarginTop(10)
	memLabel.AddCssClass("heading")
	box.Append(memLabel)

	// Memory usage would go here (placeholder)
	memPlaceholder := gtk4.NewLabel("Memory usage graphs will appear here")
	memPlaceholder.SetHAlign(gtk4.AlignCenter)
	memPlaceholder.SetVAlign(gtk4.AlignCenter)
	memPlaceholder.SetMarginTop(20)
	memPlaceholder.SetMarginBottom(20)
	memPlaceholder.SetMarginStart(20)
	memPlaceholder.SetMarginEnd(20)
	box.Append(memPlaceholder)

	return box
}

// createStatusBar creates the status bar at the bottom
func createStatusBar() *gtk4.Box {
	statusBox := gtk4.NewBox(gtk4.OrientationHorizontal, 5)
	statusBox.SetMarginTop(5)
	statusBox.SetMarginBottom(5)
	statusBox.SetMarginStart(5)
	statusBox.SetMarginEnd(5)

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

	// Clear the process box map
	processBoxes = make(map[int64]*gtk4.Box)

	// Clear the process box map
	processBoxes = make(map[int64]*gtk4.Box)

	// Just create a new box to replace the process list
	newProcessList := gtk4.NewBox(gtk4.OrientationVertical, 0)
	
	// Get the parent - we'll use the name of the box to get it rather than GetParent()
	// Get the content box from the processes tab
	if parent != nil {
		// Replace the old process list with the new one
		parent.Remove(processList)
		parent.Append(newProcessList)
		processList = newProcessList
	}

	// Get the process data
	processes, err := getProcesses()
	if err != nil {
		setStatus(fmt.Sprintf("Error: %v", err))
		return
	}

	// Add each process to the list
	for _, proc := range processes {
		processItem := createProcessItem(proc)
		processList.Append(processItem)
	}

	// Update the process count
	processCountLabel.SetText(fmt.Sprintf("%d processes", len(processes)))

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

	// Get all processes
	processes, err := getProcesses()
	if err != nil {
		setStatus(fmt.Sprintf("Error: %v", err))
		return
	}

	// Filter processes based on search text
	var filteredProcesses []ProcessInfo
	for _, proc := range processes {
		if strings.Contains(strings.ToLower(proc.Name), strings.ToLower(searchText)) {
			filteredProcesses = append(filteredProcesses, proc)
		}
	}

	// Clear the process box map
	processBoxes = make(map[int64]*gtk4.Box)

	// Just create a new box to replace the process list
	newProcessList := gtk4.NewBox(gtk4.OrientationVertical, 0)
	
	// Get the parent of the process list
	parent := processList.GetParent()
	if parent != nil {
		// Replace the old process list with the new one
		parent.Remove(processList)
		parent.Append(newProcessList)
		processList = newProcessList
	}

	// Add filtered processes to the list
	for _, proc := range filteredProcesses {
		processItem := createProcessItem(proc)
		processList.Append(processItem)
	}

	// Update the process count
	processCountLabel.SetText(fmt.Sprintf("%d processes", len(filteredProcesses)))

	setStatus(fmt.Sprintf("Found %d processes matching '%s'", len(filteredProcesses), searchText))
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
		gtk4.DialogFlagModal|gtk4.DialogFlagDestroyWithParent,
		gtk4.MessageTypeWarning,
		gtk4.ButtonsOkCancel,
		fmt.Sprintf("Are you sure you want to terminate process %d?", selectedPID),
	)
	dialog.SetSecondaryText("This may cause data loss if the application is not responding.")

	// Handle dialog response
	dialog.ConnectResponse(func(responseID int) {
		if responseID == int(gtk4.ResponseOk) {
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

// loadAppStyles loads the CSS styles for the application
func loadAppStyles() {
	// Get the CSS provider and add styles
	cssProvider := gtk4.NewCssProvider()
	cssProvider.LoadFromData(GetCSS())

	// Apply to all windows in the application
	gtk4.StyleContextAddProviderForDisplay(
		gtk4.GetDefaultDisplay(),
		cssProvider,
		gtk4.StyleProviderPriorityApplication,
	)
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
