package main

import (
	"../../../gtk4go"
	"../../gtk4"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Define constants for styling and sizing
const (
	APP_ID         = "com.example.system-info"
	TITLE          = "System Info"
	DEFAULT_WIDTH  = 900
	DEFAULT_HEIGHT = 600

	// Auto-refresh interval in seconds (0 = disabled)
	AUTO_REFRESH_INTERVAL = 30
)

// labelMap stores references to labels for updating
type labelMap struct {
	labels map[string]*gtk4.Label
}

func newLabelMap() *labelMap {
	return &labelMap{
		labels: make(map[string]*gtk4.Label),
	}
}

func (lm *labelMap) add(key string, label *gtk4.Label) {
	lm.labels[key] = label
}

func (lm *labelMap) update(key string, value string) {
	if label, ok := lm.labels[key]; ok {
		label.SetText(value)
	}
}

// Global variables for data and UI state
var (
	osLabels           *labelMap
	cpuLabels          *labelMap
	memoryLabels       *labelMap
	diskLabels         *labelMap
	statusLabel        *gtk4.Label
	autoRefreshEnabled bool = true
	isRefreshing       bool = false
	lastRefreshTime    time.Time
	autoRefreshTimer   *time.Timer
)

func main() {
	os.Setenv("GSK_RENDERER", "cairo")
	os.Setenv("GDK_GL", "0")

	// Initialize GTK
	if err := gtk4go.Initialize(); err != nil {
		fmt.Printf("Failed to initialize GTK: %v\n", err)
		os.Exit(1)
	}

	// Create application
	app := gtk4.NewApplication(APP_ID)

	// Create window
	win := gtk4.NewWindow(TITLE)
	win.SetDefaultSize(DEFAULT_WIDTH, DEFAULT_HEIGHT)

	// Create main layout
	mainBox := createMainLayout(win)

	// Set the window's child to the main box
	win.SetChild(mainBox)

	// Set up window close handler
	win.ConnectCloseRequest(func() bool {
		// Clean up resources
		if autoRefreshTimer != nil {
			autoRefreshTimer.Stop()
		}
		return false // Return false to allow window to close
	})

	// Add window to application
	app.AddWindow(win)

	// Start auto-refresh timer if enabled
	if AUTO_REFRESH_INTERVAL > 0 {
		startAutoRefreshTimer()
	}

	// Run the application
	os.Exit(app.Run())
}

// startAutoRefreshTimer starts the auto-refresh timer
func startAutoRefreshTimer() {
	if autoRefreshTimer != nil {
		autoRefreshTimer.Stop()
	}

	autoRefreshTimer = time.AfterFunc(time.Duration(AUTO_REFRESH_INTERVAL)*time.Second, func() {
		// Only refresh if auto-refresh is enabled and not currently refreshing
		if autoRefreshEnabled && !isRefreshing {
			gtk4go.RunOnUIThread(func() {
				refreshAllData()
				// Restart timer for next refresh
				startAutoRefreshTimer()
			})
		} else {
			// Restart timer anyway
			startAutoRefreshTimer()
		}
	})
}

// createMainLayout builds the main app UI structure
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
	loadAppStyles()

	// Initial data load
	refreshAllData()

	return mainBox
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

// refreshAllData updates all system information
func refreshAllData() {
	if isRefreshing {
		return
	}

	isRefreshing = true
	statusLabel.SetText("Refreshing data...")

	// Use background worker to avoid UI freezing
	gtk4go.RunInBackground(func() (interface{}, error) {
		// Refresh OS Info
		refreshOSInfo(osLabels)

		// Refresh CPU Info
		refreshCPUInfo(cpuLabels)

		// Refresh RAM Info
		refreshRAMInfo(memoryLabels)

		// Refresh Disk Info
		refreshDiskInfo(diskLabels)

		return "Data refreshed at " + time.Now().Format("15:04:05"), nil
	}, func(result interface{}, err error) {
		isRefreshing = false
		lastRefreshTime = time.Now()

		if err != nil {
			statusLabel.SetText("Error refreshing data: " + err.Error())
		} else {
			statusLabel.SetText("Ready")

			// Find the lastUpdatedLabel in the status bar and update it
			gtk4go.RunOnUIThread(func() {
				// Find all labels in the window that have the "update-time" class
				// For a real app, we'd store a reference to this label
				// This is a simplified approach
				updateTimeStr := "Last updated: " + time.Now().Format("15:04:05")
				statusLabel.SetText("Ready - " + updateTimeStr)
			})
		}
	})
}

// loadAppStyles loads CSS styles for the application
func loadAppStyles() {
	cssProvider, err := gtk4.LoadCSS(`
		window {
			background-color: #f5f5f5;
		}
		
		.header-bar {
			background-color: #3584e4;
			color: white;
			padding: 8px 16px;
			min-height: 48px;
		}
		
		.header-title {
			font-size: 18px;
			font-weight: bold;
			color: white;
		}
		
		.refresh-button {
			padding: 8px 16px;
			background-color: rgba(255, 255, 255, 0.1);
			color: white;
			border-radius: 4px;
		}
		
		.sidebar {
			background-color: #323232;
			min-width: 200px;
			padding: 0;
		}
		
		.sidebar-button {
			background-color: transparent;
			color: #eeeeee;
			border-radius: 0;
			border-left: 4px solid transparent;
			padding: 16px;
			margin: 0;
		}
		
		.sidebar-button:hover {
			background-color: rgba(255, 255, 255, 0.1);
		}
		
		.sidebar-button-selected {
			background-color: rgba(255, 255, 255, 0.15);
			border-left: 4px solid #3584e4;
			font-weight: bold;
		}
		
		.content-panel {
			padding: 24px;
			background-color: #fafafa;
		}
		
		.panel-title {
			font-size: 22px;
			font-weight: bold;
			margin-bottom: 16px;
			color: #303030;
		}
		
		.info-card {
			background-color: white;
			border-radius: 8px;
			padding: 16px;
			margin-bottom: 16px;
			box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		}
		
		.card-title {
			font-size: 16px;
			font-weight: bold;
			margin-bottom: 8px;
			color: #303030;
		}
		
		.info-grid {
			margin: 8px 0;
		}
		
		.info-key {
			font-weight: normal;
			color: #707070;
			padding-right: 16px;
		}
		
		.info-value {
			font-weight: bold;
			color: #303030;
		}
		
		.disk-info {
			font-family: monospace;
			padding: 12px;
			border-radius: 4px;
			background-color: #f5f5f5;
		}
		
		.status-bar {
			background-color: #323232;
			color: #eeeeee;
			padding: 8px 16px;
			border-top: 1px solid #444444;
		}
		
		.status-label {
			color: #eeeeee;
		}
		
		.update-time {
			color: #bbbbbb;
			font-size: 12px;
		}
		
		.toggle-button {
			background-color: rgba(255, 255, 255, 0.1);
			border-radius: 4px;
			padding: 4px 8px;
			color: #eeeeee;
			font-size: 12px;
		}
		
		.toggle-button:hover {
			background-color: rgba(255, 255, 255, 0.2);
		}
	`)

	if err != nil {
		fmt.Printf("Failed to load CSS: %v\n", err)
	} else {
		// Apply CSS provider to the entire application
		gtk4.AddProviderForDisplay(cssProvider, 600)
	}
}

// refreshOSInfo updates the OS information labels
func refreshOSInfo(labels *labelMap) {
	// Update OS Name
	if osName, err := executeCommand("uname", "-s"); err == nil {
		labels.update("os_name", strings.TrimSpace(osName))
	}

	// Update Kernel Version
	if kernelVersion, err := executeCommand("uname", "-r"); err == nil {
		labels.update("kernel_version", strings.TrimSpace(kernelVersion))
	}

	// Update Distribution
	if fileExists("/etc/os-release") {
		if dist, err := readDistribution(); err == nil {
			labels.update("distribution", dist)
		}
	} else {
		labels.update("distribution", "Unknown (os-release not found)")
	}

	// Update Architecture
	if arch, err := executeCommand("uname", "-m"); err == nil {
		labels.update("architecture", strings.TrimSpace(arch))
	}

	// Update Hostname
	if hostname, err := executeCommand("hostname"); err == nil {
		labels.update("hostname", strings.TrimSpace(hostname))
	}

	// Update Uptime
	if uptime, err := readUptime(); err == nil {
		labels.update("uptime", uptime)
	}

	// Update User
	if user, err := executeCommand("whoami"); err == nil {
		labels.update("user", strings.TrimSpace(user))
	}

	// Update Shell
	if shell, ok := os.LookupEnv("SHELL"); ok {
		labels.update("shell", shell)
	}
}

// refreshCPUInfo updates the CPU information labels
func refreshCPUInfo(labels *labelMap) {
	// Update CPU Model
	if model, err := readCPUModel(); err == nil {
		labels.update("cpu_model", model)
	}

	// Update CPU Cores and Threads
	cores, threads := getCPUCount()
	labels.update("cpu_cores", fmt.Sprintf("%d", cores))
	labels.update("cpu_threads", fmt.Sprintf("%d", threads))

	// Update CPU Frequency
	if freq, err := readCPUFrequency(); err == nil {
		labels.update("cpu_freq", freq)
	}

	// Update CPU Usage
	if usage, err := getCPUUsage(); err == nil {
		labels.update("cpu_usage", fmt.Sprintf("%.1f%%", usage))
	}
}

// refreshRAMInfo updates the RAM information labels
func refreshRAMInfo(labels *labelMap) {
	// Get memory info
	if total, used, free, err := getMemoryInfo(); err == nil {
		usagePercent := float64(used) / float64(total) * 100

		labels.update("ram_total", formatBytes(total))
		labels.update("ram_used", formatBytes(used))
		labels.update("ram_free", formatBytes(free))
		labels.update("ram_usage", fmt.Sprintf("%.1f%%", usagePercent))
	}

	// Get swap info
	if total, used, _, err := getSwapInfo(); err == nil {
		labels.update("swap_total", formatBytes(total))
		labels.update("swap_used", formatBytes(used))
	}
}

// refreshDiskInfo updates the disk information
func refreshDiskInfo(labels *labelMap) {
	// Get disk information using df command
	output, err := executeCommand("df", "-h", "--output=source,size,used,avail,pcent,target")
	if err != nil {
		labels.update("disk_info", "Error getting disk information")
		return
	}

	// Parse the output
	lines := strings.Split(output, "\n")
	if len(lines) <= 1 {
		labels.update("disk_info", "No disk information available")
		return
	}

	// Create formatted output with pretty table
	var formattedOutput strings.Builder

	// Create header
	formattedOutput.WriteString(fmt.Sprintf("%-16s %-8s %-8s %-8s %-6s %-s\n",
		"Device", "Size", "Used", "Avail", "Use%", "Mount Point"))
	formattedOutput.WriteString(strings.Repeat("-", 80) + "\n")

	// Process each line except header
	for i, line := range lines {
		if i == 0 || len(strings.TrimSpace(line)) == 0 {
			// Skip header and empty lines
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 6 {
			// Format each line for better readability
			device := fields[0]
			if len(device) > 16 {
				device = device[:13] + "..."
			}

			mountPoint := strings.Join(fields[5:], " ")
			if len(mountPoint) > 20 {
				mountPoint = mountPoint[:17] + "..."
			}

			formattedOutput.WriteString(fmt.Sprintf("%-16s %-8s %-8s %-8s %-6s %-s\n",
				device, fields[1], fields[2], fields[3], fields[4], mountPoint))
		}
	}

	labels.update("disk_info", formattedOutput.String())
}

// Helper functions for system information

// executeCommand executes a command and returns its output
func executeCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing command %s: %w", command, err)
	}
	return string(output), nil
}

// fileExists checks if a file exists
func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// readDistribution reads the Linux distribution information
func readDistribution() (string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			// Extract value between quotes
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			value := parts[1]
			// Remove quotes if present
			value = strings.Trim(value, "\"")
			return value, nil
		}
	}

	return "Unknown", fmt.Errorf("distribution not found in os-release")
}

// readUptime reads the system uptime
func readUptime() (string, error) {
	file, err := os.Open("/proc/uptime")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) >= 1 {
			uptime, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return "", err
			}

			// Convert to human-readable format
			days := int(uptime / 86400)
			hours := int(uptime/3600) % 24
			minutes := int(uptime/60) % 60

			if days > 0 {
				return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes), nil
			}
			return fmt.Sprintf("%d hours, %d minutes", hours, minutes), nil
		}
	}

	return "", fmt.Errorf("failed to parse uptime")
}

// readCPUModel reads the CPU model
func readCPUModel() (string, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "Unknown", fmt.Errorf("CPU model not found in cpuinfo")
}

// getCPUCount returns the number of physical cores and threads
func getCPUCount() (int, int) {
	threads := runtime.NumCPU()

	// Try to get physical core count from /proc/cpuinfo
	physicalsMap := make(map[string]bool)

	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return threads, threads // Fallback to logical cores
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "physical id") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				physicalsMap[strings.TrimSpace(parts[1])] = true
			}
		}
	}

	physicals := len(physicalsMap)
	if physicals == 0 {
		physicals = 1 // Fallback to at least 1 physical CPU
	}

	return physicals, threads
}

// readCPUFrequency reads the CPU frequency
func readCPUFrequency() (string, error) {
	// First try to get from /proc/cpuinfo
	file, err := os.Open("/proc/cpuinfo")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "cpu MHz") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					freq, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
					if err == nil {
						return fmt.Sprintf("%.2f GHz", freq/1000), nil
					}
				}
			}
		}
	}

	// Fallback to using lscpu
	output, err := executeCommand("lscpu")
	if err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "CPU MHz") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					freq, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
					if err == nil {
						return fmt.Sprintf("%.2f GHz", freq/1000), nil
					}
				}
			}
		}
	}

	return "Unknown", fmt.Errorf("CPU frequency not found")
}

// getCPUUsage gets the CPU usage percentage
func getCPUUsage() (float64, error) {
	// Using a simple approach with top command
	output, err := executeCommand("top", "-bn1")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Cpu(s)") {
			parts := strings.Split(line, ",")
			for _, part := range parts {
				if strings.Contains(part, "id") {
					// Extract idle percentage
					idlePart := strings.TrimSpace(part)
					idle, err := strconv.ParseFloat(strings.Split(idlePart, " ")[0], 64)
					if err == nil {
						return 100.0 - idle, nil
					}
				}
			}
		}
	}

	return 0, fmt.Errorf("CPU usage not found")
}

// getMemoryInfo gets memory information (total, used, free)
func getMemoryInfo() (uint64, uint64, uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()

	var total, free, available uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				total, _ = strconv.ParseUint(parts[1], 10, 64)
				total *= 1024 // Convert from KB to bytes
			}
		} else if strings.HasPrefix(line, "MemFree:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				free, _ = strconv.ParseUint(parts[1], 10, 64)
				free *= 1024 // Convert from KB to bytes
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				available, _ = strconv.ParseUint(parts[1], 10, 64)
				available *= 1024 // Convert from KB to bytes
			}
		}
	}

	used := total - available

	return total, used, free, nil
}

// getSwapInfo gets swap information (total, used, free)
func getSwapInfo() (uint64, uint64, uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()

	var total, free uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "SwapTotal:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				total, _ = strconv.ParseUint(parts[1], 10, 64)
				total *= 1024 // Convert from KB to bytes
			}
		} else if strings.HasPrefix(line, "SwapFree:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				free, _ = strconv.ParseUint(parts[1], 10, 64)
				free *= 1024 // Convert from KB to bytes
			}
		}
	}

	used := total - free

	return total, used, free, nil
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
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
