package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
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
	processListView   *gtk4.ListView
	processListModel  *gtk4.StringList
	selectionModel    *gtk4.SingleSelection
	refreshTimer      *time.Timer
	refreshTimerMutex sync.Mutex
	refreshInProgress atomic.Bool // Flag to prevent concurrent refreshes
	statusLabel       *gtk4.Label
	columnSortOrder   = make(map[int]bool) // true for ascending, false for descending
	statusMutex       sync.Mutex
	selectedPID       atomic.Int64  // Using atomic for thread safety
	processCache      []ProcessInfo // Cache of process data for selection and sorting

	// Performance settings
	autoRefreshEnabled atomic.Bool
	processLimit       = 100 // Limit the number of processes displayed for performance
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

// Add throttling for search operations
var (
	searchThrottleTimer *time.Timer
	searchMutex         sync.Mutex
	lastSearchText      string
)

func main() {
	// Initialize GTK
	if err := gtk4go.Initialize(); err != nil {
		fmt.Printf("Failed to initialize GTK: %v\n", err)
		os.Exit(1)
	}

	// Initialize selectedPID to -1
	selectedPID.Store(-1)

	// Enable auto-refresh by default
	autoRefreshEnabled.Store(true)

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
	win.EnableAcceleratedRendering() // Enable hardware acceleration
	win.OptimizeForResizing()        // Optimize for resizing

	// Create main layout
	mainBox := createMainLayout()

	// Set the window's child to the main box
	win.SetChild(mainBox)

	// Set up window close handler
	win.ConnectCloseRequest(func() bool {
		// Clean up resources
		refreshTimerMutex.Lock()
		if refreshTimer != nil {
			refreshTimer.Stop()
		}
		refreshTimerMutex.Unlock()

		// Shut down the background worker with a timeout
		gtk4go.ShutdownDefaultWorker(2 * time.Second)

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
	stack := gtk4.NewStack(
		gtk4.WithTransitionType(gtk4.StackTransitionTypeSlideLeftRight),
		gtk4.WithTransitionDuration(200),
	)

	// Create processes tab
	processesTab := createProcessesTab()
	stack.AddTitled(processesTab, "processes", "Processes")

	// Create performance tab
	performanceTab := createPerformanceTab()
	stack.AddTitled(performanceTab, "performance", "Performance")

	// Create stack switcher (tabs)
	stackSwitcher := gtk4.NewStackSwitcher(stack)

	// Add a box to hold the stack switcher centered
	switcherBox := gtk4.NewBox(gtk4.OrientationHorizontal, 0)
	switcherBox.Append(stackSwitcher)

	// Add the components to the main layout
	mainBox.Append(switcherBox)
	mainBox.Append(stack)

	// Create status bar
	statusBar := createStatusBar()
	mainBox.Append(statusBar)

	return mainBox
}

// createHeaderBar creates the header bar with search and actions
func createHeaderBar() *gtk4.Box {
	headerBox := gtk4.NewBox(gtk4.OrientationHorizontal, 10)
	headerBox.AddCssClass("header-box")

	// Search entry
	searchEntry := gtk4.NewEntry()
	searchEntry.SetPlaceholderText("Search processes...")
	searchEntry.ConnectChanged(func() {
		// Throttle search to prevent excessive updates
		throttledSearch(searchEntry.GetText())
	})
	headerBox.Append(searchEntry)

	// Refresh button
	refreshButton := gtk4.NewButton("Refresh")
	refreshButton.AddCssClass("refresh-button")
	refreshButton.ConnectClicked(refreshProcessList)
	headerBox.Append(refreshButton)

	// Auto-refresh toggle
	autoRefreshButton := gtk4.NewButton("Auto-refresh: ON")
	autoRefreshButton.ConnectClicked(func() {
		// Toggle auto-refresh
		newValue := !autoRefreshEnabled.Load()
		autoRefreshEnabled.Store(newValue)

		if newValue {
			autoRefreshButton.SetLabel("Auto-refresh: ON")
			// Restart the refresh timer
			startRefreshTimer()
		} else {
			autoRefreshButton.SetLabel("Auto-refresh: OFF")
			// Stop the refresh timer
			refreshTimerMutex.Lock()
			if refreshTimer != nil {
				refreshTimer.Stop()
				refreshTimer = nil
			}
			refreshTimerMutex.Unlock()
		}
	})
	headerBox.Append(autoRefreshButton)

	// End Process button
	endProcessButton := gtk4.NewButton("End Process")
	endProcessButton.AddCssClass("end-process-button")
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
	)

	// Create string list model to display processes
	processListModel = gtk4.NewStringList()

	// Create a selection model for the list
	selectionModel = gtk4.NewSingleSelection(processListModel,
		gtk4.WithAutoselect(false),
	)

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
		label.SetHExpand(true)
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
		boxWidget := listItem.GetChild()
		if listItem.GetSelected() {
			boxWidget.AddCssClass("selected")
		} else {
			boxWidget.RemoveCssClass("selected")
		}
	})

	// Create the list view with selection model and factory
	processListView = gtk4.NewListView(
		selectionModel,
		factory,
		gtk4.WithShowSeparators(true),
		gtk4.WithSingleClickActivate(true),
	)

	// Connect activate signal for item selection
	processListView.ConnectActivate(func(position int) {
		// Use the gtk4go background worker to safely handle selection
		gtk4go.RunInBackground(
			func() (interface{}, error) {
				// This runs in background thread
				processCacheMutex.RLock()
				defer processCacheMutex.RUnlock()

				if position < 0 || position >= len(processCache) {
					return nil, fmt.Errorf("invalid position: %d", position)
				}

				// Return a copy of the process info to avoid race conditions
				return ProcessInfo{
					PID:         processCache[position].PID,
					Name:        processCache[position].Name,
					Username:    processCache[position].Username,
					CPUPercent:  processCache[position].CPUPercent,
					MemoryBytes: processCache[position].MemoryBytes,
					Threads:     processCache[position].Threads,
					State:       processCache[position].State,
					StartTime:   processCache[position].StartTime,
				}, nil
			},
			func(result interface{}, err error) {
				// This runs on UI thread
				if err != nil {
					setStatus(fmt.Sprintf("Process selection error: %v", err))
					selectedPID.Store(-1)
					return
				}

				// Get the process info and update selectedPID
				proc := result.(ProcessInfo)
				selectedPID.Store(proc.PID)
				setStatus(fmt.Sprintf("Selected process: %s (PID: %d)", proc.Name, proc.PID))
			},
		)
	})

	// Add the list view to the scrolled window
	scrollWin.SetChild(processListView)

	return scrollWin
}

// createPerformanceTab creates the performance tab content
func createPerformanceTab() *gtk4.Box {
	box := gtk4.NewBox(gtk4.OrientationVertical, 10)
	box.AddCssClass("performance-tab")

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

	// Use background worker to update performance data periodically
	updatePerformanceData(cpuValueLabel, memValueLabel)

	return box
}

// updatePerformanceData sets up periodic updates of CPU and memory using background worker
func updatePerformanceData(cpuLabel, memLabel *gtk4.Label) {
	// Initial update
	updateCPUAndMemory(cpuLabel, memLabel)

	// Set up a timer to trigger updates
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Schedule update on UI thread
			gtk4go.RunOnUIThread(func() {
				updateCPUAndMemory(cpuLabel, memLabel)
			})
		}
	}()
}

// updateCPUAndMemory updates CPU and memory labels using background worker
func updateCPUAndMemory(cpuLabel, memLabel *gtk4.Label) {
	// Update CPU usage using background worker
	gtk4go.RunInBackground(
		func() (interface{}, error) {
			// This runs in background thread
			return getCPUUsage()
		},
		func(result interface{}, err error) {
			// This runs on UI thread
			if err != nil {
				cpuLabel.SetText("CPU Usage: Error - " + err.Error())
				return
			}

			cpuUsage := result.(float64)
			cpuLabel.SetText(fmt.Sprintf("CPU Usage: %.1f%%", cpuUsage))

			// Update styling based on CPU usage
			cpuLabel.RemoveCssClass("usage-high")
			cpuLabel.AddCssClass("usage-normal")
			if cpuUsage > 80 {
				cpuLabel.RemoveCssClass("usage-normal")
				cpuLabel.AddCssClass("usage-high")
			}
		},
	)

	// Update memory usage using background worker
	gtk4go.RunInBackground(
		func() (interface{}, error) {
			// This runs in background thread
			total, free, err := getSystemMemoryInfo()
			if err != nil {
				return nil, err
			}

			return map[string]int64{
				"total": total,
				"free":  free,
			}, nil
		},
		func(result interface{}, err error) {
			// This runs on UI thread
			if err != nil {
				memLabel.SetText("Memory Usage: Error - " + err.Error())
				return
			}

			memInfo := result.(map[string]int64)
			total := memInfo["total"]
			free := memInfo["free"]

			used := total - free
			usedPercentage := float64(used) / float64(total) * 100

			// Calculate values in MB for display
			totalMB := total / (1024 * 1024)
			usedMB := used / (1024 * 1024)

			memLabel.SetText(fmt.Sprintf("Memory Usage: %d MB / %d MB (%.1f%%)",
				usedMB, totalMB, usedPercentage))

			// Update styling based on memory usage
			memLabel.RemoveCssClass("usage-high")
			memLabel.AddCssClass("usage-normal")
			if usedPercentage > 80 {
				memLabel.RemoveCssClass("usage-normal")
				memLabel.AddCssClass("usage-high")
			}
		},
	)
}

// createStatusBar creates the status bar at the bottom
func createStatusBar() *gtk4.Box {
	statusBox := gtk4.NewBox(gtk4.OrientationHorizontal, 5)
	statusBox.AddCssClass("status-bar")

	// Status label
	statusLabel = gtk4.NewLabel("Ready")
	statusBox.Append(statusLabel)

	return statusBox
}

// startRefreshTimer starts the auto-refresh timer with thread safety
func startRefreshTimer() {
	// If auto-refresh is disabled, don't start the timer
	if !autoRefreshEnabled.Load() {
		return
	}

	refreshTimerMutex.Lock()
	defer refreshTimerMutex.Unlock()

	if refreshTimer != nil {
		refreshTimer.Stop()
	}

	refreshTimer = time.AfterFunc(time.Duration(AUTO_REFRESH_INTERVAL)*time.Second, func() {
		// Schedule refresh on UI thread
		gtk4go.RunOnUIThread(func() {
			// Skip if a refresh is already in progress
			if !refreshInProgress.Load() {
				refreshProcessList()
			}
			// Schedule next refresh
			startRefreshTimer()
		})
	})
}

// refreshProcessList refreshes the process list with current data
func refreshProcessList() {
	// Use atomic flag to prevent concurrent refreshes
	if refreshInProgress.Swap(true) {
		// Already refreshing, skip this request
		return
	}

	setStatus("Refreshing process list...")

	// Use the GTK4Go background worker to get process data without blocking the UI
	gtk4go.QueueBackgroundTask(
		"refresh-processes",
		// Background task function
		func(ctx context.Context, progress func(percent int, message string)) (interface{}, error) {
			// This runs in a background thread
			progress(10, "Loading process list...")

			// Get the process data
			processes, err := getProcesses()
			if err != nil {
				return nil, err
			}

			progress(50, "Sorting processes...")

			// Sort processes by CPU usage (highest first) for quicker access to important processes
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].CPUPercent > processes[j].CPUPercent
			})

			// Limit the number of processes for better performance
			if len(processes) > processLimit {
				processes = processes[:processLimit]
			}

			progress(100, "Completed")
			return processes, nil
		},
		// Completion callback function (runs on UI thread)
		func(result interface{}, err error) {
			// Always clear the flag when done
			defer refreshInProgress.Store(false)

			if err != nil {
				setStatus(fmt.Sprintf("Error: %v", err))
				return
			}

			// Type assertion to get processes
			processes, ok := result.([]ProcessInfo)
			if !ok {
				setStatus("Error: Invalid result type from background task")
				return
			}

			// Clear the current list efficiently
			removeAllListItems()

			// Update the process cache with mutex protection
			processCacheMutex.Lock()
			processCache = processes
			processCacheMutex.Unlock()

			// Add processes to the list efficiently
			addProcessesToList(processes)

			setStatus(fmt.Sprintf("Process list refreshed. %d processes found (displaying %d).",
				len(processes),
				processListModel.GetNItems()))
		},
		// Progress callback function (runs on UI thread)
		func(percent int, message string) {
			setStatus(fmt.Sprintf("Refreshing process list... %d%% %s", percent, message))
		},
	)
}

// Helper functions for efficient list management

// removeAllListItems efficiently removes all items from the list model
func removeAllListItems() {
	// Get the number of items once to avoid potential issues with changing length
	count := processListModel.GetNItems()
	if count > 0 {
		// Just remove all items directly instead of batching to simplify the code
		// and reduce potential for race conditions
		for i := count - 1; i >= 0; i-- {
			if i < processListModel.GetNItems() {
				processListModel.Remove(i)
			}
		}
	}
}

// addProcessesToList efficiently adds processes to the list model
func addProcessesToList(processes []ProcessInfo) {
	// Add processes in a single batch for better performance
	for _, proc := range processes {
		// Safely format the process info
		displayText := formatProcessInfo(proc)

		// Add to the list model
		processListModel.Append(displayText)
	}
}

// formatProcessInfo safely formats process information for display
func formatProcessInfo(proc ProcessInfo) string {
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

	// Calculate memory in MB with a minimum of 1MB to avoid showing 0MB
	memoryMB := proc.MemoryBytes / (1024 * 1024)
	if memoryMB <= 0 && proc.MemoryBytes > 0 {
		memoryMB = 1 // Show at least 1MB if there's any memory usage
	}

	// Format the display string
	return fmt.Sprintf("%d | %s | %s | %.1f%% | %dMB | %d | %s | %s",
		proc.PID,
		proc.Name,
		username,
		proc.CPUPercent,
		memoryMB,
		proc.Threads,
		state,
		startTime,
	)
}

// throttledSearch implements a throttled search to prevent excessive UI updates
func throttledSearch(searchText string) {
	searchMutex.Lock()
	defer searchMutex.Unlock()

	// Save the latest search text
	lastSearchText = searchText

	// Cancel existing timer if any
	if searchThrottleTimer != nil {
		searchThrottleTimer.Stop()
	}

	// Create a new timer that will execute the search after a delay
	searchThrottleTimer = time.AfterFunc(300*time.Millisecond, func() {
		// Get the latest search text safely
		searchMutex.Lock()
		text := lastSearchText
		searchMutex.Unlock()

		// Execute the search on the UI thread
		gtk4go.RunOnUIThread(func() {
			searchProcesses(text)
		})
	})
}

// searchProcesses filters the process list based on search text
func searchProcesses(searchText string) {
	if searchText == "" {
		// If search is empty, just refresh the full list
		refreshProcessList()
		return
	}

	setStatus(fmt.Sprintf("Searching for '%s'...", searchText))

	// Skip if a refresh is already in progress
	if refreshInProgress.Swap(true) {
		return
	}

	// Use the GTK4Go background worker for searching
	gtk4go.QueueBackgroundTask(
		"search-processes",
		// Background task function
		func(ctx context.Context, progress func(percent int, message string)) (interface{}, error) {
			// Get the process data
			progress(10, "Loading process list...")
			processes, err := getProcesses()
			if err != nil {
				return nil, err
			}

			// Filter and collect matching processes
			progress(50, "Filtering processes...")
			var filteredProcesses []ProcessInfo
			count := 0

			// Search with lowercase for case-insensitive matching
			searchLower := strings.ToLower(searchText)

			for _, proc := range processes {
				// Check for cancellation periodically
				select {
				case <-ctx.Done():
					return nil, ctx.Err() // Task was cancelled
				default:
					// Continue processing
				}

				// Limit results for performance
				if count >= processLimit {
					break
				}

				// Simple case-insensitive substring search in process name, username, or PID
				nameLower := strings.ToLower(proc.Name)
				userLower := strings.ToLower(proc.Username)
				pidStr := fmt.Sprintf("%d", proc.PID)

				if strings.Contains(nameLower, searchLower) ||
					strings.Contains(userLower, searchLower) ||
					strings.Contains(pidStr, searchText) {
					filteredProcesses = append(filteredProcesses, proc)
					count++
				}
			}

			progress(100, "Completed")
			return filteredProcesses, nil
		},
		// Completion callback function (runs on UI thread)
		func(result interface{}, err error) {
			// Always clear the flag when done
			defer refreshInProgress.Store(false)

			if err != nil {
				setStatus(fmt.Sprintf("Search error: %v", err))
				return
			}

			// Type assertion to get filtered processes
			filteredProcesses, ok := result.([]ProcessInfo)
			if !ok {
				setStatus("Error: Invalid search result type")
				return
			}

			// Clear the current list efficiently
			removeAllListItems()

			// Update the process cache with the filtered list
			processCacheMutex.Lock()
			processCache = filteredProcesses
			processCacheMutex.Unlock()

			// Add filtered processes to the list
			addProcessesToList(filteredProcesses)

			setStatus(fmt.Sprintf("Found %d processes matching '%s'", len(filteredProcesses), searchText))
		},
		// Progress callback function (runs on UI thread)
		func(percent int, message string) {
			setStatus(fmt.Sprintf("Searching for '%s'... %d%% %s", searchText, percent, message))
		},
	)
}

// endSelectedProcess attempts to terminate the selected process
func endSelectedProcess() {
	// Get selectedPID safely using atomic operation
	pid := selectedPID.Load()
	if pid <= 0 {
		setStatus("No process selected")
		return
	}

	// Show confirmation dialog
	dialog := gtk4.NewMessageDialog(
		nil, // parent window
		gtk4.DialogModal|gtk4.DialogDestroyWithParent,
		gtk4.MessageWarning,
		gtk4.ResponseOk|gtk4.ResponseCancel,
		fmt.Sprintf("Are you sure you want to terminate process %d?", pid),
	)
	// Add secondary text
	secondaryLabel := gtk4.NewLabel("This may cause data loss if the application is not responding.")
	dialog.GetContentArea().Append(secondaryLabel)

	// Handle dialog response
	dialog.ConnectResponse(func(responseId gtk4.ResponseType) {
		if responseId == gtk4.ResponseOk {
			// User confirmed, try to kill the process using the saved pid value
			// Use background worker to avoid blocking UI
			// Capture pid for use in closure
			pidToKill := pid

			gtk4go.QueueBackgroundTask(
				"terminate-process",
				// Background task function
				func(ctx context.Context, progress func(percent int, message string)) (interface{}, error) {
					// This runs in a background thread
					progress(10, "Validating process...")

					// Validate the PID exists before trying to kill it
					processes, err := getProcesses()
					if err != nil {
						return nil, fmt.Errorf("failed to get process list: %v", err)
					}

					// Check if process exists
					found := false
					for _, proc := range processes {
						if proc.PID == pidToKill {
							found = true
							break
						}
					}

					if !found {
						return nil, fmt.Errorf("process %d no longer exists", pidToKill)
					}

					progress(50, "Terminating process...")

					// Terminate the process
					err = terminateProcess(pidToKill)
					if err != nil {
						return nil, err
					}

					progress(100, "Process terminated")
					return pidToKill, nil
				},
				// Completion callback function (runs on UI thread)
				func(result interface{}, err error) {
					if err != nil {
						setStatus(fmt.Sprintf("Error terminating process: %v", err))
					} else {
						killedPid := result.(int64)
						setStatus(fmt.Sprintf("Process %d terminated successfully", killedPid))

						// Reset selected PID if we killed the selected process
						if selectedPID.Load() == killedPid {
							selectedPID.Store(-1)
						}

						// Refresh process list
						refreshProcessList()
					}
				},
				// Progress callback function (runs on UI thread)
				func(percent int, message string) {
					setStatus(fmt.Sprintf("Terminating process %d... %d%% %s",
						pidToKill, percent, message))
				},
			)
		}
		dialog.Destroy()
	})

	dialog.Show()
}

// sortProcessList sorts the process list by the specified column
func sortProcessList(columnID int, ascending bool) {
	// Skip if a refresh is already in progress
	if refreshInProgress.Swap(true) {
		return
	}

	setStatus(fmt.Sprintf("Sorting by column %d...", columnID))

	// Use background worker for sorting
	gtk4go.QueueBackgroundTask(
		"sort-processes",
		// Background task function
		func(ctx context.Context, progress func(percent int, message string)) (interface{}, error) {
			progress(0, "Getting processes...")

			// Get a copy of the process cache to sort
			processCacheMutex.RLock()
			processes := make([]ProcessInfo, len(processCache))
			copy(processes, processCache)
			processCacheMutex.RUnlock()

			progress(25, "Sorting processes...")

			// Sort based on column ID
			sort.SliceStable(processes, func(i, j int) bool {
				var result bool

				switch columnID {
				case COL_PID:
					result = processes[i].PID < processes[j].PID
				case COL_NAME:
					result = processes[i].Name < processes[j].Name
				case COL_USERNAME:
					result = processes[i].Username < processes[j].Username
				case COL_CPU:
					result = processes[i].CPUPercent < processes[j].CPUPercent
				case COL_MEMORY:
					result = processes[i].MemoryBytes < processes[j].MemoryBytes
				case COL_THREADS:
					result = processes[i].Threads < processes[j].Threads
				case COL_STATE:
					result = processes[i].State < processes[j].State
				case COL_STARTED:
					result = processes[i].StartTime < processes[j].StartTime
				default:
					result = processes[i].PID < processes[j].PID
				}

				// Reverse order if descending
				if !ascending {
					result = !result
				}

				return result
			})

			progress(100, "Sorting complete")
			return processes, nil
		},
		// Completion callback function (runs on UI thread)
		func(result interface{}, err error) {
			// Always clear the flag when done
			defer refreshInProgress.Store(false)

			if err != nil {
				setStatus(fmt.Sprintf("Sort error: %v", err))
				return
			}

			// Get sorted processes
			sortedProcesses, ok := result.([]ProcessInfo)
			if !ok {
				setStatus("Error: Invalid sort result type")
				return
			}

			// Clear the current list
			removeAllListItems()

			// Update the process cache
			processCacheMutex.Lock()
			processCache = sortedProcesses
			processCacheMutex.Unlock()

			// Add sorted processes to the list
			addProcessesToList(sortedProcesses)

			// Store sort order
			columnSortOrder[columnID] = ascending

			if ascending {
				setStatus(fmt.Sprintf("Sorted by column %d (%s)", columnID, "ascending"))
			} else {
				setStatus(fmt.Sprintf("Sorted by column %d (%s)", columnID, "descending"))
			}
		},
		// Progress callback function (runs on UI thread)
		func(percent int, message string) {
			setStatus(fmt.Sprintf("Sorting... %d%% %s", percent, message))
		},
	)
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
	// Always run on UI thread
	gtk4go.RunOnUIThread(func() {
		statusMutex.Lock()
		defer statusMutex.Unlock()
		statusLabel.SetText(message)
	})
}
