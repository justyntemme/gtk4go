package main

import (
	"../../gtk4go"
	"../gtk4"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
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

func main() {
	// Force software rendering to avoid OpenGL crashes
	// These environment variables need to be set BEFORE GTK is initialized
	os.Setenv("GSK_RENDERER", "cairo") // Use Cairo renderer instead of GL
	os.Setenv("GDK_GL", "0")           // Disable OpenGL
	os.Setenv("GDK_BACKEND", "x11")    // Force X11 backend which is more stable

	// Initialize GTK
	if err := gtk4go.Initialize(); err != nil {
		fmt.Printf("Failed to initialize GTK: %v\n", err)
		os.Exit(1)
	}

	// Create application
	app := gtk4.NewApplication("com.example.linuxsysteminfo")

	// Create window
	win := gtk4.NewWindow("Linux System Info")
	win.SetDefaultSize(800, 600)
	// Disable hardware acceleration to prevent crashes
	win.DisableAcceleratedRendering()

	// Create main vertical box
	mainBox := gtk4.NewBox(gtk4.OrientationVertical, 10)

	// Create header with title
	headerLabel := gtk4.NewLabel("Linux System Information")
	headerLabel.AddCssClass("header-label")
	mainBox.Append(headerLabel)

	// Create stack for tabs
	stack := gtk4.NewStack(
		gtk4.WithTransitionType(gtk4.StackTransitionTypeSlideLeftRight),
		gtk4.WithTransitionDuration(200),
	)

	// Create OS Info tab
	osInfoBox, osLabels := createOSInfoTab()
	stack.AddTitled(osInfoBox, "os-info", "OS Info")

	// Create Hardware tab
	hardwareBox, cpuLabels, ramLabels, diskLabel := createHardwareTab()
	stack.AddTitled(hardwareBox, "hardware", "Hardware")

	// Create stack switcher (tab bar)
	stackSwitcher := gtk4.NewStackSwitcher(stack)
	stackSwitcher.AddCssClass("stack-switcher")

	// Add the stack switcher and stack to the main box
	mainBox.Append(stackSwitcher)
	mainBox.Append(stack)

	// Load CSS for styling
	cssProvider, err := gtk4.LoadCSS(`
		.header-label {
			font-size: 24px;
			font-weight: bold;
			padding: 15px;
			color: #2a76c6;
		}
		.stack-switcher {
			padding: 8px;
		}
		.info-label {
			font-size: 14px;
			padding: 4px;
		}
		.info-category {
			font-size: 18px;
			font-weight: bold;
			padding: 10px;
			color: #2a76c6;
			border-bottom: 1px solid #cccccc;
			margin-bottom: 10px;
		}
		.info-grid {
			margin: 15px;
		}
		.info-key {
			font-weight: bold;
			color: #333333;
			padding-right: 10px;
		}
		.info-value {
			color: #0066cc;
		}
		.info-section {
			background-color: #f5f5f5;
			border-radius: 4px;
			padding: 12px;
			margin: 8px;
		}
		.usage-bar {
			min-height: 20px;
			border-radius: 3px;
			background-color: #e0e0e0;
		}
		.usage-bar-fill {
			background-color: #4caf50;
			border-radius: 3px;
		}
		.refresh-button {
			padding: 8px 16px;
			background-color: #3584e4;
			color: white;
			font-weight: bold;
			border-radius: 4px;
		}
		.refresh-button:hover {
			background-color: #1c71d8;
		}
		.status-label {
			font-style: italic;
			color: #666666;
			padding: 5px;
		}
	`)
	if err != nil {
		fmt.Printf("Failed to load CSS: %v\n", err)
	} else {
		// Apply CSS provider to the entire application
		gtk4.AddProviderForDisplay(cssProvider, 600)
	}

	// Create refresh button
	refreshButton := gtk4.NewButton("Refresh Data")
	refreshButton.AddCssClass("refresh-button")

	// Status label
	statusLabel := gtk4.NewLabel("Ready")
	statusLabel.AddCssClass("status-label")

	// Create bottom bar with refresh button and status
	bottomBar := gtk4.NewBox(gtk4.OrientationHorizontal, 10)
	bottomBar.Append(refreshButton)
	bottomBar.Append(statusLabel)
	mainBox.Append(bottomBar)

	// Set up refresh functionality
	refreshButton.ConnectClicked(func() {
		statusLabel.SetText("Refreshing data...")

		// Use background worker to avoid UI freezing
		gtk4go.RunInBackground(func() (interface{}, error) {
			// Refresh OS Info
			refreshOSInfo(osLabels)

			// Refresh CPU Info
			refreshCPUInfo(cpuLabels)

			// Refresh RAM Info
			refreshRAMInfo(ramLabels)

			// Refresh Disk Info
			refreshDiskInfo(diskLabel)

			return "Data refreshed at " + time.Now().Format("15:04:05"), nil
		}, func(result interface{}, err error) {
			if err != nil {
				statusLabel.SetText("Error refreshing data: " + err.Error())
			} else {
				statusLabel.SetText(result.(string))
			}
		})
	})

	// Auto-refresh when application starts
	win.ConnectCloseRequest(func() bool {
		return false // Return false to allow window to close
	})

	// Initial data load
	refreshOSInfo(osLabels)
	refreshCPUInfo(cpuLabels)
	refreshRAMInfo(ramLabels)
	refreshDiskInfo(diskLabel)

	// Set the window's child to the main box
	win.SetChild(mainBox)

	// Add window to application
	app.AddWindow(win)

	// Run the application
	os.Exit(app.Run())
}

// createOSInfoTab creates the OS Information tab
func createOSInfoTab() (*gtk4.Box, *labelMap) {
	// Create main box
	osInfoBox := gtk4.NewBox(gtk4.OrientationVertical, 10)
	osInfoBox.AddCssClass("info-section")

	// Create section header
	osHeader := gtk4.NewLabel("Operating System Information")
	osHeader.AddCssClass("info-category")
	osInfoBox.Append(osHeader)

	// Create grid for labels
	grid := gtk4.NewGrid(
		gtk4.WithRowSpacing(8),
		gtk4.WithColumnSpacing(15),
		gtk4.WithRowHomogeneous(false),
	)
	grid.AddCssClass("info-grid")

	// Add row headers
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

	osInfoBox.Append(grid)

	return osInfoBox, labels
}

// createHardwareTab creates the Hardware Information tab
func createHardwareTab() (*gtk4.Box, *labelMap, *labelMap, *labelMap) {
	// Create main box
	hardwareBox := gtk4.NewBox(gtk4.OrientationVertical, 10)

	// CPU Info Section
	cpuSection := gtk4.NewBox(gtk4.OrientationVertical, 5)
	cpuSection.AddCssClass("info-section")

	cpuHeader := gtk4.NewLabel("CPU Information")
	cpuHeader.AddCssClass("info-category")
	cpuSection.Append(cpuHeader)

	// Create grid for CPU labels
	cpuGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(8),
		gtk4.WithColumnSpacing(15),
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

	cpuSection.Append(cpuGrid)
	hardwareBox.Append(cpuSection)

	// RAM Info Section
	ramSection := gtk4.NewBox(gtk4.OrientationVertical, 5)
	ramSection.AddCssClass("info-section")

	ramHeader := gtk4.NewLabel("Memory Information")
	ramHeader.AddCssClass("info-category")
	ramSection.Append(ramHeader)

	// Create grid for RAM labels
	ramGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(8),
		gtk4.WithColumnSpacing(15),
		gtk4.WithRowHomogeneous(false),
	)
	ramGrid.AddCssClass("info-grid")

	ramLabels := newLabelMap()

	// Total RAM
	addInfoRow(ramGrid, 0, "Total RAM:", "", ramLabels, "ram_total")

	// Used RAM
	addInfoRow(ramGrid, 1, "Used RAM:", "", ramLabels, "ram_used")

	// Free RAM
	addInfoRow(ramGrid, 2, "Free RAM:", "", ramLabels, "ram_free")

	// RAM Usage
	addInfoRow(ramGrid, 3, "RAM Usage:", "", ramLabels, "ram_usage")

	// Swap Total
	addInfoRow(ramGrid, 4, "Swap Total:", "", ramLabels, "swap_total")

	// Swap Used
	addInfoRow(ramGrid, 5, "Swap Used:", "", ramLabels, "swap_used")

	ramSection.Append(ramGrid)
	hardwareBox.Append(ramSection)

	// Disk Info Section
	diskSection := gtk4.NewBox(gtk4.OrientationVertical, 5)
	diskSection.AddCssClass("info-section")

	diskHeader := gtk4.NewLabel("Disk Information")
	diskHeader.AddCssClass("info-category")
	diskSection.Append(diskHeader)

	// Create container for disk info
	diskLabels := newLabelMap()

	// We'll create a scrolled window for disk information since there could be many disks
	scrollWin := gtk4.NewScrolledWindow(
		gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyNever),
		gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
	)

	diskBox := gtk4.NewBox(gtk4.OrientationVertical, 5)
	diskBox.AddCssClass("info-grid")
	scrollWin.SetChild(diskBox)

	diskSection.Append(scrollWin)
	hardwareBox.Append(diskSection)

	return hardwareBox, cpuLabels, ramLabels, diskLabels
}

// addInfoRow adds a row to the info grid with a key/value pair
func addInfoRow(grid *gtk4.Grid, row int, key string, value string, labels *labelMap, labelKey string) {
	keyLabel := gtk4.NewLabel(key)
	keyLabel.AddCssClass("info-key")
	grid.Attach(keyLabel, 0, row, 1, 1)

	valueLabel := gtk4.NewLabel(value)
	valueLabel.AddCssClass("info-value")
	grid.Attach(valueLabel, 1, row, 1, 1)

	labels.add(labelKey, valueLabel)
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

	// Create formatted output
	var formattedOutput strings.Builder
	for i, line := range lines {
		if i == 0 || len(strings.TrimSpace(line)) == 0 {
			// Skip header or empty lines
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 5 {
			// Format: Device, Size, Used, Available, Usage %, Mount Point
			deviceInfo := fmt.Sprintf("Device: %s\nSize: %s\nUsed: %s\nAvailable: %s\nUsage: %s\nMount: %s\n\n",
				fields[0], fields[1], fields[2], fields[3], fields[4],
				strings.Join(fields[5:], " "))
			formattedOutput.WriteString(deviceInfo)
		}
	}

	labels.update("disk_info", formattedOutput.String())
}

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

