package main

import (
	"../../gtk4"
	"fmt"
	"sync"
)

// Global variables for disk display
var (
	diskCard    *gtk4.Box  // Container for disk information
	currentGrid *gtk4.Grid // Current grid being displayed
	uiMutex     sync.Mutex // Mutex to protect UI operations
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
	autoRefreshButton.AddCssClass("dark-area-btn")

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

	statusBar.Append(spacer)
	statusBar.Append(autoRefreshButton)

	return statusBar
}

// createHardwarePanel builds the hardware information panel
func createHardwarePanel() (*gtk4.Box, *labelMap, *labelMap, *labelMap, *labelMap) {
	// Create main container with scrolling
	containerBox := gtk4.NewBox(gtk4.OrientationVertical, 0)

	scrollWin := gtk4.NewScrolledWindow(
		gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyNever),
		gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
		gtk4.WithPropagateNaturalWidth(true),
		gtk4.WithPropagateNaturalHeight(true),
		gtk4.WithHExpand(true),
		gtk4.WithVExpand(true),
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

	// CPU Grid - horizontal grid with equal column spacing
	cpuGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(4),
		gtk4.WithColumnSpacing(12),
		gtk4.WithRowHomogeneous(false),
		gtk4.WithColumnHomogeneous(true), // Make columns equal width
	)
	cpuGrid.AddCssClass("disk-info-grid")

	cpuLabels := newLabelMap()

	// Create column headers similar to disk/memory sections
	cpuHeaders := []string{"CPU Model", "CPU Cores", "CPU Threads", "CPU Frequency", "CPU Usage"}

	// Add headers to the grid with consistent styling
	for i, header := range cpuHeaders {
		label := gtk4.NewLabel(header)
		label.AddCssClass("disk-header")
		cpuGrid.Attach(label, i, 0, 1, 1)
	}

	// Add a separator row like in other sections
	for i := 0; i < len(cpuHeaders); i++ {
		separator := gtk4.NewLabel("--------")
		separator.AddCssClass("disk-separator")
		cpuGrid.Attach(separator, i, 1, 1, 1)
	}

	// Create value labels for the second row and add to label map
	// CPU Model - Column 0
	cpuModelValue := gtk4.NewLabel("")
	cpuModelValue.AddCssClass("disk-device") // Using disk-device for model name styling
	cpuGrid.Attach(cpuModelValue, 0, 2, 1, 1)
	cpuLabels.add("cpu_model", cpuModelValue)

	// CPU Cores - Column 1
	cpuCoresValue := gtk4.NewLabel("")
	cpuCoresValue.AddCssClass("disk-size")
	cpuGrid.Attach(cpuCoresValue, 1, 2, 1, 1)
	cpuLabels.add("cpu_cores", cpuCoresValue)

	// CPU Threads - Column 2
	cpuThreadsValue := gtk4.NewLabel("")
	cpuThreadsValue.AddCssClass("disk-size")
	cpuGrid.Attach(cpuThreadsValue, 2, 2, 1, 1)
	cpuLabels.add("cpu_threads", cpuThreadsValue)

	// CPU Frequency - Column 3
	cpuFreqValue := gtk4.NewLabel("")
	cpuFreqValue.AddCssClass("disk-avail")
	cpuGrid.Attach(cpuFreqValue, 3, 2, 1, 1)
	cpuLabels.add("cpu_freq", cpuFreqValue)

	// CPU Usage - Column 4
	cpuUsageValue := gtk4.NewLabel("")
	cpuUsageValue.AddCssClass("disk-percent")
	cpuGrid.Attach(cpuUsageValue, 4, 2, 1, 1)
	cpuLabels.add("cpu_usage", cpuUsageValue)

	cpuCard.Append(cpuGrid)
	panel.Append(cpuCard)

	// Create GPU info card
	gpuCard := gtk4.NewBox(gtk4.OrientationVertical, 8)
	gpuCard.AddCssClass("info-card")

	// GPU Section Header
	gpuHeader := gtk4.NewLabel("GPU Information")
	gpuHeader.AddCssClass("card-title")
	gpuCard.Append(gpuHeader)

	// GPU Grid - keep existing vertical layout but ensure consistent styling
	gpuGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(4),
		gtk4.WithColumnSpacing(12),
		gtk4.WithRowHomogeneous(false),
	)
	gpuGrid.AddCssClass("disk-info-grid")

	// Add headers
	gpuHeaders := []string{"Property", "Value"}
	for i, header := range gpuHeaders {
		label := gtk4.NewLabel(header)
		label.AddCssClass("disk-header")
		gpuGrid.Attach(label, i, 0, 1, 1)
	}

	// Add a separator row
	for i := 0; i < len(gpuHeaders); i++ {
		separator := gtk4.NewLabel("--------")
		separator.AddCssClass("disk-separator")
		gpuGrid.Attach(separator, i, 1, 1, 1)
	}

	gpuLabels := newLabelMap()

	// Start adding property rows at row 2 (after headers and separator)
	// GPU Model
	addInfoRow(gpuGrid, 2, "GPU Model:", "", gpuLabels, "gpu_model")
	// GPU Vendor
	addInfoRow(gpuGrid, 3, "GPU Vendor:", "", gpuLabels, "gpu_vendor")
	// GPU Renderer
	addInfoRow(gpuGrid, 4, "GPU Renderer:", "", gpuLabels, "gpu_renderer")
	// GPU Driver
	addInfoRow(gpuGrid, 5, "GPU Driver:", "", gpuLabels, "gpu_driver")
	// GPU OpenGL Version
	addInfoRow(gpuGrid, 6, "OpenGL Version:", "", gpuLabels, "gpu_gl_version")
	// GPU Memory (only for NVIDIA GPUs)
	addInfoRow(gpuGrid, 7, "GPU Memory:", "", gpuLabels, "gpu_memory")
	// GPU Utilization (only for NVIDIA GPUs)
	addInfoRow(gpuGrid, 8, "GPU Utilization:", "", gpuLabels, "gpu_utilization")

	gpuCard.Append(gpuGrid)
	panel.Append(gpuCard)

	// Create Memory info card
	memoryCard := gtk4.NewBox(gtk4.OrientationVertical, 8)
	memoryCard.AddCssClass("info-card")

	// Memory Section Header
	memoryHeader := gtk4.NewLabel("Memory Information")
	memoryHeader.AddCssClass("card-title")
	memoryCard.Append(memoryHeader)

	// Memory Grid - a horizontal grid with all RAM info in one row with equal spacing
	memoryGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(4),
		gtk4.WithColumnSpacing(12),
		gtk4.WithRowHomogeneous(false),
		gtk4.WithColumnHomogeneous(true), // Make columns equal width
	)
	memoryGrid.AddCssClass("disk-info-grid")

	memoryLabels := newLabelMap()

	// Create column headers similar to disk section
	memHeaders := []string{"Total RAM", "Used RAM", "Free RAM", "RAM Usage", "Swap Total", "Swap Used"}

	// Add headers to the grid with consistent styling
	for i, header := range memHeaders {
		label := gtk4.NewLabel(header)
		label.AddCssClass("disk-header")
		memoryGrid.Attach(label, i, 0, 1, 1)
	}

	// Add a separator row like in the disk section
	for i := 0; i < len(memHeaders); i++ {
		separator := gtk4.NewLabel("--------")
		separator.AddCssClass("disk-separator")
		memoryGrid.Attach(separator, i, 1, 1, 1)
	}

	// Create value labels for the second row and add to label map
	// Total RAM - Column 0
	totalRamValue := gtk4.NewLabel("")
	totalRamValue.AddCssClass("disk-size") // Using disk-size for consistent styling
	memoryGrid.Attach(totalRamValue, 0, 2, 1, 1)
	memoryLabels.add("ram_total", totalRamValue)

	// Used RAM - Column 1
	usedRamValue := gtk4.NewLabel("")
	usedRamValue.AddCssClass("disk-used")
	memoryGrid.Attach(usedRamValue, 1, 2, 1, 1)
	memoryLabels.add("ram_used", usedRamValue)

	// Free RAM - Column 2
	freeRamValue := gtk4.NewLabel("")
	freeRamValue.AddCssClass("disk-avail")
	memoryGrid.Attach(freeRamValue, 2, 2, 1, 1)
	memoryLabels.add("ram_free", freeRamValue)

	// RAM Usage - Column 3
	ramUsageValue := gtk4.NewLabel("")
	ramUsageValue.AddCssClass("disk-percent")
	memoryGrid.Attach(ramUsageValue, 3, 2, 1, 1)
	memoryLabels.add("ram_usage", ramUsageValue)

	// Swap Total - Column 4
	swapTotalValue := gtk4.NewLabel("")
	swapTotalValue.AddCssClass("disk-size")
	memoryGrid.Attach(swapTotalValue, 4, 2, 1, 1)
	memoryLabels.add("swap_total", swapTotalValue)

	// Swap Used - Column 5
	swapUsedValue := gtk4.NewLabel("")
	swapUsedValue.AddCssClass("disk-used")
	memoryGrid.Attach(swapUsedValue, 5, 2, 1, 1)
	memoryLabels.add("swap_used", swapUsedValue)

	memoryCard.Append(memoryGrid)
	panel.Append(memoryCard)

	// Create Disk info card - this is protected by a mutex when updated
	uiMutex.Lock()
	diskCard = gtk4.NewBox(gtk4.OrientationVertical, 8)
	diskCard.AddCssClass("info-card")

	// Disk Section Header
	diskHeader := gtk4.NewLabel("Disk Information")
	diskHeader.AddCssClass("card-title")
	diskCard.Append(diskHeader)

	// Create initial grid for disk info
	initialGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(4),
		gtk4.WithColumnSpacing(12),
		gtk4.WithRowHomogeneous(false),
	)
	initialGrid.AddCssClass("disk-info-grid")

	// Add column headers to the grid
	diskHeaders := []string{"Device", "Size", "Used", "Avail", "Use%", "Mount Point"}
	for i, header := range diskHeaders {
		label := gtk4.NewLabel(header)
		label.AddCssClass("disk-header")
		initialGrid.Attach(label, i, 0, 1, 1)
	}

	// Add a separator row
	for i := 0; i < len(diskHeaders); i++ {
		separator := gtk4.NewLabel("--------")
		separator.AddCssClass("disk-separator")
		initialGrid.Attach(separator, i, 1, 1, 1)
	}

	// Add a loading message
	loadingLabel := gtk4.NewLabel("Loading disk information...")
	loadingLabel.AddCssClass("disk-info-message")
	initialGrid.Attach(loadingLabel, 0, 2, 6, 1)

	// Set as current grid and add to card
	currentGrid = initialGrid
	diskCard.Append(currentGrid)
	uiMutex.Unlock()

	// Add card to panel
	panel.Append(diskCard)

	// Create disk labels map (for backward compatibility)
	diskLabels := newLabelMap()

	// Add a placeholder label for text-based info (for backward compatibility)
	placeholderLabel := gtk4.NewLabel("")
	diskLabels.add("disk_info", placeholderLabel)

	// Set the panel as the scrollable content
	scrollWin.SetChild(panel)

	// Add scrolled window to the container box
	containerBox.Append(scrollWin)

	return containerBox, cpuLabels, memoryLabels, diskLabels, gpuLabels
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

	// Spacer to push menu button to the right
	spacer := gtk4.NewBox(gtk4.OrientationHorizontal, 0)
	spacer.SetHExpand(true)

	// Create a gear menu button
	menuButton := gtk4.NewMenuButton()
	menuButton.AddCssClass("dark-area-btn")
	menuButton.SetIconName("emblem-system-symbolic") // Standard GTK gear icon

	// Create menu model for the menu button
	menu := gtk4.NewMenu()

	// Add "Refresh" menu item
	refreshItem := gtk4.NewMenuItem("Refresh", "app.refresh")
	menu.AppendItem(refreshItem)

	// Create a popover menu for the button
	popoverMenu := gtk4.NewPopoverMenu(menu)
	menuButton.SetPopover(popoverMenu)

	// Add elements to header
	headerBar.Append(titleLabel)
	headerBar.Append(spacer)
	headerBar.Append(menuButton)

	return headerBar
}

func createMainLayout(win *gtk4.Window) *gtk4.Box {
	// Create main vertical box
	mainBox := gtk4.NewBox(gtk4.OrientationVertical, 0)

	// Create headerbar for the window
	headerBar := gtk4.NewHeaderBar(
		gtk4.WithShowTitleButtons(true),
		gtk4.WithTitle(TITLE),
	)

	// Create a menu button for the header bar
	menuButton := gtk4.NewMenuButton()
	menuButton.SetIconName("emblem-system-symbolic") // Standard GTK gear icon

	// Create menu model for the menu button
	menu := gtk4.NewMenu()

	// Add "Refresh" menu item
	refreshItem := gtk4.NewMenuItem("Refresh", "app.refresh")
	menu.AppendItem(refreshItem)

	// Create a popover menu for the button
	popoverMenu := gtk4.NewPopoverMenu(menu)
	menuButton.SetPopover(popoverMenu)

	// Add the menu button to the end of the header bar
	headerBar.PackEnd(menuButton)

	// Set the header bar as the window's titlebar
	win.SetTitlebar(headerBar)

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

	hardwarePanel, cpuLabelsMap, memLabelsMap, diskLabelsMap, gpuLabelsMap := createHardwarePanel()
	cpuLabels = cpuLabelsMap
	memoryLabels = memLabelsMap
	diskLabels = diskLabelsMap
	gpuLabels = gpuLabelsMap

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
	err := loadAppStyles()
	if err != nil {
		fmt.Println("Failed to load CSS styles:", err)
	}

	// Initial data load
	refreshAllData()

	return mainBox
}
