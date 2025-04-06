package main

import (
	"../../gtk4"
	"fmt"
)

// addInfoRow adds a row to an info grid with key/value pair
func addInfoRow(grid *gtk4.Grid, row int, key string, value string, labels *labelMap, labelKey string) {
	keyLabel := gtk4.NewLabel(key)
	keyLabel.AddCssClass("info-key")
	grid.Attach(keyLabel, 0, row, 1, 1)

	valueLabel := gtk4.NewLabel(value)
	valueLabel.AddCssClass("info-value")
	grid.Attach(valueLabel, 1, row, 1, 1)

	labels.add(labelKey, valueLabel)
}

// createStatusBar builds the status bar at the bottom of the window
func createStatusBar() *gtk4.Box {
	statusBar := gtk4.NewBox(gtk4.OrientationHorizontal, 8)
	statusBar.AddCssClass("status-bar")

	// Status indicator
	statusLabel = gtk4.NewLabel("Ready")
	statusLabel.AddCssClass("status-label")

	// Add auto-refresh toggle button
	autoRefreshButton := gtk4.NewButton("Auto-refresh: On")
	autoRefreshButton.AddCssClass("toggle-button")

	// Connect button to toggle auto-refresh state
	autoRefreshButton.ConnectClicked(func() {
		autoRefreshEnabled = !autoRefreshEnabled
		if autoRefreshEnabled {
			autoRefreshButton.SetLabel("Auto-refresh: On")
			// Start auto-refresh if interval is set
			if AUTO_REFRESH_INTERVAL > 0 {
				startAutoRefreshTimer()
			}
		} else {
			autoRefreshButton.SetLabel("Auto-refresh: Off")
			// Stop auto-refresh timer
			if autoRefreshTimer != nil {
				autoRefreshTimer.Stop()
			}
		}
	})

	// Add elements to status bar
	statusBar.Append(statusLabel)

	// Create right-aligned area for refresh info
	spacer := gtk4.NewBox(gtk4.OrientationHorizontal, 0)
	spacer.SetHExpand(true)

	// Last updated info
	lastUpdatedLabel := gtk4.NewLabel("Last updated: Never")
	lastUpdatedLabel.AddCssClass("update-time")

	statusBar.Append(spacer)
	statusBar.Append(autoRefreshButton)
	statusBar.Append(lastUpdatedLabel)

	return statusBar
}

// createHardwarePanel builds the hardware information panel
func createHardwarePanel() (*gtk4.Box, *labelMap, *labelMap, *labelMap) {
	// Create main container with scrolling
	containerBox := gtk4.NewBox(gtk4.OrientationVertical, 0)

	scrollWin := gtk4.NewScrolledWindow(
		gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyNever),
		gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
		gtk4.WithPropagateNaturalWidth(true), gtk4.WithPropagateNaturalHeight(true),
	)

	panel := gtk4.NewBox(gtk4.OrientationVertical, 16)
	panel.AddCssClass("content-panel")

	// Section title
	titleLabel := gtk4.NewLabel("Hardware Information")
	titleLabel.AddCssClass("panel-title")
	panel.Append(titleLabel)

	// Create CPU info card
	cpuCard := gtk4.NewBox(gtk4.OrientationVertical, 8)
	cpuCard.AddCssClass("info-card")

	// CPU Section Header
	cpuHeader := gtk4.NewLabel("CPU Information")
	cpuHeader.AddCssClass("card-title")
	cpuCard.Append(cpuHeader)

	// CPU Grid
	cpuGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(8),
		gtk4.WithColumnSpacing(24),
		gtk4.WithRowHomogeneous(false),
	)
	cpuGrid.AddCssClass("info-grid")

	cpuLabels := newLabelMap()

	// CPU Model
	addInfoRow(cpuGrid, 0, "CPU Model:", "", cpuLabels, "cpu_model")

	// CPU Cores
	addInfoRow(cpuGrid, 1, "CPU Cores:", "", cpuLabels, "cpu_cores")

	// CPU Threads
	addInfoRow(cpuGrid, 2, "CPU Threads:", "", cpuLabels, "cpu_threads")

	// CPU Frequency
	addInfoRow(cpuGrid, 3, "CPU Frequency:", "", cpuLabels, "cpu_freq")

	// CPU Usage
	addInfoRow(cpuGrid, 4, "CPU Usage:", "", cpuLabels, "cpu_usage")

	cpuCard.Append(cpuGrid)
	panel.Append(cpuCard)

	// Create Memory info card
	memoryCard := gtk4.NewBox(gtk4.OrientationVertical, 8)
	memoryCard.AddCssClass("info-card")

	// Memory Section Header
	memoryHeader := gtk4.NewLabel("Memory Information")
	memoryHeader.AddCssClass("card-title")
	memoryCard.Append(memoryHeader)

	// Memory Grid
	memoryGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(8),
		gtk4.WithColumnSpacing(24),
		gtk4.WithRowHomogeneous(false),
	)
	memoryGrid.AddCssClass("info-grid")

	memoryLabels := newLabelMap()

	// Total RAM
	addInfoRow(memoryGrid, 0, "Total RAM:", "", memoryLabels, "ram_total")

	// Used RAM
	addInfoRow(memoryGrid, 1, "Used RAM:", "", memoryLabels, "ram_used")

	// Free RAM
	addInfoRow(memoryGrid, 2, "Free RAM:", "", memoryLabels, "ram_free")

	// RAM Usage
	addInfoRow(memoryGrid, 3, "RAM Usage:", "", memoryLabels, "ram_usage")

	// Swap Total
	addInfoRow(memoryGrid, 4, "Swap Total:", "", memoryLabels, "swap_total")

	// Swap Used
	addInfoRow(memoryGrid, 5, "Swap Used:", "", memoryLabels, "swap_used")

	memoryCard.Append(memoryGrid)
	panel.Append(memoryCard)

	// Create Disk info card
	diskCard := gtk4.NewBox(gtk4.OrientationVertical, 8)
	diskCard.AddCssClass("info-card")

	// Disk Section Header
	diskHeader := gtk4.NewLabel("Disk Information")
	diskHeader.AddCssClass("card-title")
	diskCard.Append(diskHeader)

	// Create box for disk info
	diskBox := gtk4.NewBox(gtk4.OrientationVertical, 8)
	diskBox.AddCssClass("info-grid")

	// Disk Storage Label (will contain formatted disk info)
	diskInfoLabel := gtk4.NewLabel("")
	diskInfoLabel.AddCssClass("disk-info")
	diskBox.Append(diskInfoLabel)

	diskLabels := newLabelMap()
	diskLabels.add("disk_info", diskInfoLabel)

	diskCard.Append(diskBox)
	panel.Append(diskCard)

	// Set the panel as the scrollable content
	scrollWin.SetChild(panel)

	// Add scrolled window to the container box
	containerBox.Append(scrollWin)

	return containerBox, cpuLabels, memoryLabels, diskLabels
}

// createSystemPanel builds the system information panel
func createSystemPanel() (*gtk4.Box, *labelMap) {
	// Create main container
	panel := gtk4.NewBox(gtk4.OrientationVertical, 16)
	panel.AddCssClass("content-panel")

	// Section title
	titleLabel := gtk4.NewLabel("System Information")
	titleLabel.AddCssClass("panel-title")
	panel.Append(titleLabel)

	// Create info card
	card := gtk4.NewBox(gtk4.OrientationVertical, 8)
	card.AddCssClass("info-card")

	// Create grid for info items
	grid := gtk4.NewGrid(
		gtk4.WithRowSpacing(8),
		gtk4.WithColumnSpacing(24),
		gtk4.WithRowHomogeneous(false),
	)
	grid.AddCssClass("info-grid")

	// Add labels map to store references
	labels := newLabelMap()

	// OS Name
	addInfoRow(grid, 0, "OS Name:", "", labels, "os_name")

	// Kernel Version
	addInfoRow(grid, 1, "Kernel Version:", "", labels, "kernel_version")

	// Distribution
	addInfoRow(grid, 2, "Distribution:", "", labels, "distribution")

	// Architecture
	addInfoRow(grid, 3, "Architecture:", "", labels, "architecture")

	// Hostname
	addInfoRow(grid, 4, "Hostname:", "", labels, "hostname")

	// Uptime
	addInfoRow(grid, 5, "Uptime:", "", labels, "uptime")

	// User
	addInfoRow(grid, 6, "Current User:", "", labels, "user")

	// Shell
	addInfoRow(grid, 7, "Default Shell:", "", labels, "shell")

	// Add grid to card
	card.Append(grid)

	// Add card to panel
	panel.Append(card)

	return panel, labels
}

// createSidebar builds the navigation sidebar
func createSidebar(stack *gtk4.Stack) *gtk4.Box {
	sidebar := gtk4.NewBox(gtk4.OrientationVertical, 0)
	sidebar.AddCssClass("sidebar")
	sidebar.SetHExpand(false)

	// System Info button
	systemBtn := gtk4.NewButton("System Info")
	systemBtn.AddCssClass("sidebar-button")
	systemBtn.AddCssClass("sidebar-button-selected")

	// Hardware button
	hardwareBtn := gtk4.NewButton("Hardware")
	hardwareBtn.AddCssClass("sidebar-button")

	// Connect System button click handler to show the system panel
	systemBtn.ConnectClicked(func() {
		// Update visual selection state
		systemBtn.AddCssClass("sidebar-button-selected")
		hardwareBtn.RemoveCssClass("sidebar-button-selected")

		// Switch stack to system panel
		stack.SetVisibleChildName("system")
	})

	// Connect Hardware button click handler to show the hardware panel
	hardwareBtn.ConnectClicked(func() {
		// Update visual selection state
		hardwareBtn.AddCssClass("sidebar-button-selected")
		systemBtn.RemoveCssClass("sidebar-button-selected")

		// Switch stack to hardware panel
		stack.SetVisibleChildName("hardware")
	})

	// Add buttons to sidebar
	sidebar.Append(systemBtn)
	sidebar.Append(hardwareBtn)

	// Add spacing at the bottom of the sidebar
	spacer := gtk4.NewBox(gtk4.OrientationVertical, 0)
	spacer.SetVExpand(true)
	sidebar.Append(spacer)

	return sidebar
}

// createHeaderBar builds the application header bar
func createHeaderBar() *gtk4.Box {
	headerBar := gtk4.NewBox(gtk4.OrientationHorizontal, 0)
	headerBar.AddCssClass("header-bar")

	// App Title
	titleLabel := gtk4.NewLabel(TITLE)
	titleLabel.AddCssClass("header-title")

	// Spacer to push refresh button to the right
	spacer := gtk4.NewBox(gtk4.OrientationHorizontal, 0)
	spacer.SetHExpand(true)

	// Refresh button with icon
	refreshButton := gtk4.NewButton("Refresh")
	refreshButton.AddCssClass("refresh-button")

	// Connect refresh button click
	refreshButton.ConnectClicked(func() {
		refreshAllData()
	})

	// Add elements to header
	headerBar.Append(titleLabel)
	headerBar.Append(spacer)
	headerBar.Append(refreshButton)

	return headerBar
}

func createMainLayout(win *gtk4.Window) *gtk4.Box {
	// Create main vertical box
	mainBox := gtk4.NewBox(gtk4.OrientationVertical, 0)

	// Create header with application title and controls
	header := createHeaderBar()
	mainBox.Append(header)

	// Create a horizontal box for sidebar and content
	contentBox := gtk4.NewBox(gtk4.OrientationHorizontal, 0)

	// Create stack for different views
	stack := gtk4.NewStack(
		gtk4.WithTransitionType(gtk4.StackTransitionTypeSlideLeftRight),
		gtk4.WithTransitionDuration(200),
	)

	// Create each info panel
	systemPanel, osLabelsMap := createSystemPanel()
	osLabels = osLabelsMap

	hardwarePanel, cpuLabelsMap, memLabelsMap, diskLabelsMap := createHardwarePanel()
	cpuLabels = cpuLabelsMap
	memoryLabels = memLabelsMap
	diskLabels = diskLabelsMap

	// Add panels to stack
	stack.AddTitled(systemPanel, "system", "System")
	stack.AddTitled(hardwarePanel, "hardware", "Hardware")

	// Create sidebar for navigation (pass stack for navigation)
	sidebar := createSidebar(stack)

	// Create bottom status bar
	statusBox := createStatusBar()

	// Add sidebar and stack to content box
	contentBox.Append(sidebar)
	contentBox.Append(stack)

	// Add content and status bar to main box
	mainBox.Append(contentBox)
	mainBox.Append(statusBox)

	// Load CSS styling
	testing()
	err := loadAppStyles()
	if err != nil {
		fmt.Println("We should probably implement logging TODO")
	}

	// Initial data load
	refreshAllData()

	return mainBox
}
