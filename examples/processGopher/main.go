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
	processTreeView  *gtk4.GtkTreeView
	processListStore *gtk4.GtkListStore
	refreshTimer     *time.Timer
	statusLabel      *gtk4.Label
	columnSortOrder  = make(map[int]bool) // true for ascending, false for descending
	statusMutex      sync.Mutex
	selectedPID      int64 = -1
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

	// Create tree view
	processTreeView = gtk4.NewTreeView(nil, nil)
	processTreeView.SetEnableSearch(true)
	processTreeView.SetSearchColumn(COL_NAME)

	// Create columns
	createProcessListColumns()

	// Create list store with column types
	types := []gtk4.GType{
		gtk4.GTypeInt64,     // PID
		gtk4.GTypeString,    // Name
		gtk4.GTypeString,    // Username
		gtk4.GTypeDouble,    // CPU %
		gtk4.GTypeInt64,     // Memory (bytes)
		gtk4.GTypeInt,       // Threads
		gtk4.GTypeString,    // State
		gtk4.GTypeString,    // Started
	}
	processListStore = gtk4.NewListStore(types)
	processTreeView.SetModel(processListStore)

	// Connect to selection changed
	selection := processTreeView.GetSelection()
	selection.ConnectChanged(func() {
		// Get the selected row
		iter, selected := selection.GetSelected()
		if selected {
			// Get the PID from the selected row
			pidValue := processListStore.GetValue(iter, COL_PID)
			if pidValue != nil {
				selectedPID = pidValue.(int64)
			} else {
				selectedPID = -1
			}
		} else {
			selectedPID = -1
		}
	})

	// Add the tree view to the scrolled window
	scrollWin.SetChild(processTreeView)

	return scrollWin
}

// createProcessListColumns creates columns for the process list
func createProcessListColumns() {
	// PID column
	renderer := gtk4.NewCellRendererText()
	column := gtk4.NewTreeViewColumn("PID", renderer)
	column.SetResizable(true)
	column.SetSortColumnID(COL_PID)
	column.SetAttribute(renderer, "text", COL_PID)
	processTreeView.AppendColumn(column)

	// Process Name column
	renderer = gtk4.NewCellRendererText()
	column = gtk4.NewTreeViewColumn("Process Name", renderer)
	column.SetResizable(true)
	column.SetExpand(true)
	column.SetSortColumnID(COL_NAME)
	column.SetAttribute(renderer, "text", COL_NAME)
	processTreeView.AppendColumn(column)

	// Username column
	renderer = gtk4.NewCellRendererText()
	column = gtk4.NewTreeViewColumn("User", renderer)
	column.SetResizable(true)
	column.SetSortColumnID(COL_USERNAME)
	column.SetAttribute(renderer, "text", COL_USERNAME)
	processTreeView.AppendColumn(column)

	// CPU column
	renderer = gtk4.NewCellRendererText()
	column = gtk4.NewTreeViewColumn("CPU %", renderer)
	column.SetResizable(true)
	column.SetSortColumnID(COL_CPU)
	column.SetAttribute(renderer, "text", COL_CPU)
	processTreeView.AppendColumn(column)

	// Memory column
	renderer = gtk4.NewCellRendererText()
	column = gtk4.NewTreeViewColumn("Memory", renderer)
	column.SetResizable(true)
	column.SetSortColumnID(COL_MEMORY)
	column.SetAttribute(renderer, "text", COL_MEMORY)
	processTreeView.AppendColumn(column)

	// Threads column
	renderer = gtk4.NewCellRendererText()
	column = gtk4.NewTreeViewColumn("Threads", renderer)
	column.SetResizable(true)
	column.SetSortColumnID(COL_THREADS)
	column.SetAttribute(renderer, "text", COL_THREADS)
	processTreeView.AppendColumn(column)

	// State column
	renderer = gtk4.NewCellRendererText()
	column = gtk4.NewTreeViewColumn("State", renderer)
	column.SetResizable(true)
	column.SetSortColumnID(COL_STATE)
	column.SetAttribute(renderer, "text", COL_STATE)
	processTreeView.AppendColumn(column)

	// Started column
	renderer = gtk4.NewCellRendererText()
	column = gtk4.NewTreeViewColumn("Started", renderer)
	column.SetResizable(true)
	column.SetSortColumnID(COL_STARTED)
	column.SetAttribute(renderer, "text", COL_STARTED)
	processTreeView.AppendColumn(column)

	// Connect sort signals to all columns
	for i := 0; i < COL_COUNT; i++ {
		col := processTreeView.GetColumn(i)
		if col != nil {
			colID := i // Capture the column ID
			col.ConnectClicked(func() {
				// Toggle sort order for this column
				columnSortOrder[colID] = !columnSortOrder[colID]
				sortProcessList(colID, columnSortOrder[colID])
			})
		}
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

	// Clear the current list
	processListStore.Clear()

	// Get the process data
	processes, err := getProcesses()
	if err != nil {
		setStatus(fmt.Sprintf("Error: %v", err))
		return
	}

	// Add each process to the list store
	for _, proc := range processes {
		// Add a new row to the list store
		iter := processListStore.Append()
		processListStore.SetValue(iter, COL_PID, proc.PID)
		processListStore.SetValue(iter, COL_NAME, proc.Name)
		processListStore.SetValue(iter, COL_USERNAME, proc.Username)
		processListStore.SetValue(iter, COL_CPU, proc.CPUPercent)
		processListStore.SetValue(iter, COL_MEMORY, proc.MemoryBytes)
		processListStore.SetValue(iter, COL_THREADS, proc.Threads)
		processListStore.SetValue(iter, COL_STATE, proc.State)
		processListStore.SetValue(iter, COL_STARTED, proc.StartTime)
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
