// Package main provides system information functionality for both Linux and macOS
// File: info.go - Contains common code and platform-agnostic functions
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/justyntemme/gtk4go"
	"github.com/justyntemme/gtk4go/gtk4"
)

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

// startAutoRefreshTimer starts the auto-refresh timer
func startAutoRefreshTimer() {
	if autoRefreshTimer != nil {
		autoRefreshTimer.Stop()
	}

	autoRefreshTimer = time.AfterFunc(time.Duration(AUTO_REFRESH_INTERVAL)*time.Second, func() {
		// Only refresh if auto-refresh is enabled and not currently refreshing
		if autoRefreshEnabled && refreshAtomicFlag.Load() == 0 {
			// Ensure we call refreshAllData on the UI thread
			gtk4go.RunOnUIThread(func() {
				refreshAllData()
			})
		} else {
			// Restart timer anyway
			startAutoRefreshTimer()
		}
	})
}

// refreshAllData updates all system information
func refreshAllData() {
	// Use atomic operation to check and set isRefreshing
	// This ensures only one refresh can happen at a time
	if !refreshAtomicFlag.CompareAndSwap(0, 1) {
		// Another refresh is already in progress
		return
	}

	// Make sure we update the UI from the UI thread
	gtk4go.RunOnUIThread(func() {
		statusLabel.SetText("Refreshing data...")
	})

	// Use background worker to avoid UI freezing
	gtk4go.RunInBackground(func() (interface{}, error) {
		// Refresh OS Info
		refreshOSInfo(osLabels)

		// Refresh CPU Info
		refreshCPUInfo(cpuLabels)

		// Refresh RAM Info
		refreshRAMInfo(memoryLabels)

		// Refresh GPU Info
		refreshGPUInfo(gpuLabels)

		// Refresh Disk Info
		refreshDiskInfo(diskLabels)

		return "Data refreshed at " + time.Now().Format("15:04:05"), nil
	}, func(result interface{}, err error) {
		// Reset the refreshing flag when done
		refreshAtomicFlag.Store(0)

		// This runs on the UI thread
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

// refreshOSInfo updates the OS information labels
func refreshOSInfo(labels *labelMap) {
	// Update OS Name
	if osName, err := executeCommand("uname", "-s"); err == nil {
		text := strings.TrimSpace(osName)
		labels.update("os_name", text)

		// Add tooltip for potentially long values
		if label, ok := labels.labels["os_name"]; ok && len(text) > 20 {
			label.SetTooltipText(text)
		}
	}

	// Update Kernel Version
	if kernelVersion, err := executeCommand("uname", "-r"); err == nil {
		text := strings.TrimSpace(kernelVersion)
		labels.update("kernel_version", text)

		// Add tooltip for potentially long values
		if label, ok := labels.labels["kernel_version"]; ok && len(text) > 20 {
			label.SetTooltipText(text)
		}
	}

	// Update Distribution or OS Version
	if dist, err := readDistribution(); err == nil {
		text := dist
		labels.update("distribution", text)

		// Add tooltip for potentially long values
		if label, ok := labels.labels["distribution"]; ok && len(text) > 20 {
			label.SetTooltipText(text)
		}
	} else {
		labels.update("distribution", "Unknown (distribution information not available)")
	}

	// Update Architecture
	if arch, err := executeCommand("uname", "-m"); err == nil {
		text := strings.TrimSpace(arch)
		labels.update("architecture", text)

		// Add tooltip for potentially long values
		if label, ok := labels.labels["architecture"]; ok && len(text) > 15 {
			label.SetTooltipText(text)
		}
	}

	// Update Hostname
	if hostname, err := executeCommand("hostname"); err == nil {
		text := strings.TrimSpace(hostname)
		labels.update("hostname", text)

		// Add tooltip for potentially long values
		if label, ok := labels.labels["hostname"]; ok && len(text) > 20 {
			label.SetTooltipText(text)
		}
	}

	// Update Uptime
	if uptime, err := readUptime(); err == nil {
		text := uptime
		labels.update("uptime", text)

		// Add tooltip for potentially long values
		if label, ok := labels.labels["uptime"]; ok && len(text) > 20 {
			label.SetTooltipText(text)
		}
	}

	// Update User
	if user, err := executeCommand("whoami"); err == nil {
		text := strings.TrimSpace(user)
		labels.update("user", text)

		// Add tooltip for potentially long values
		if label, ok := labels.labels["user"]; ok && len(text) > 15 {
			label.SetTooltipText(text)
		}
	}

	// Update Shell
	if shell, ok := os.LookupEnv("SHELL"); ok {
		text := shell
		labels.update("shell", text)

		// Add tooltip for potentially long values
		if label, ok := labels.labels["shell"]; ok && len(text) > 20 {
			label.SetTooltipText(text)
		}
	}
}

// refreshCPUInfo updates the CPU information labels
func refreshCPUInfo(labels *labelMap) {
	// Update CPU Model
	if model, err := readCPUModel(); err == nil {
		// Cross-platform way to get a reasonable display length
		displayText := model
		if len(displayText) > 25 {
			// Try to get just the model name, avoiding excessively long strings
			parts := strings.Fields(model)
			if len(parts) > 3 {
				displayText = strings.Join(parts[:3], " ") + "..."
			} else {
				displayText = model[:22] + "..."
			}
		}

		labels.update("cpu_model", displayText)

		// CPU model strings are often very long, always add tooltip
		if label, ok := labels.labels["cpu_model"]; ok {
			label.SetTooltipText(model)
		}
	}

	// Update CPU Cores and Threads
	cores, threads := getCPUCount()
	labels.update("cpu_cores", fmt.Sprintf("%d", cores))
	labels.update("cpu_threads", fmt.Sprintf("%d", threads))

	// Update CPU Frequency
	if freq, err := readCPUFrequency(); err == nil {
		text := freq
		labels.update("cpu_freq", text)

		// Add tooltip for potentially long values
		if label, ok := labels.labels["cpu_freq"]; ok && len(text) > 15 {
			label.SetTooltipText(text)
		}
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
	// df command is available on both Linux and macOS, but with slightly different flags
	var dfArgs []string

	if runtime.GOOS == "darwin" {
		// macOS df command
		dfArgs = []string{"-h"}
	} else {
		// Linux df command with more detailed output format
		dfArgs = []string{"-h", "--output=source,size,used,avail,pcent,target"}
	}

	// Get disk information using df command
	output, err := executeCommand("df", dfArgs...)

	// Prepare all data off the UI thread (no UI components access)
	var grid *gtk4.Grid
	if err != nil {
		// Create a grid with just an error message, but don't attach to UI yet
		grid = createEmptyDiskGrid("Error getting disk information")
	} else {
		// Parse the output
		lines := strings.Split(output, "\n")
		if len(lines) <= 1 {
			// If there's no data, create a grid with just a message
			grid = createEmptyDiskGrid("No disk information available")
		} else {
			// Create a grid with headers
			grid = createDiskGridWithHeaders()

			// Process each line except header
			rowIndex := 2 // Start at row 2 (after headers and separator)

			// Skip processing the header line
			for i, line := range lines {
				if i == 0 || len(strings.TrimSpace(line)) == 0 {
					// Skip header and empty lines
					continue
				}

				fields := strings.Fields(line)

				// Handle differences between Linux and macOS df output
				if len(fields) >= 6 {
					// Standard output has enough fields
					addDiskRowToGrid(grid, rowIndex, fields)
					rowIndex++
				} else if len(fields) >= 5 && runtime.GOOS == "darwin" {
					// macOS df without the --output flag has 5 fields
					// Convert to the expected format: Filesystem, Size, Used, Avail, Capacity, Mounted on
					adjustedFields := []string{
						fields[0], // Filesystem
						fields[1], // Size
						fields[2], // Used
						fields[3], // Avail
						fields[4], // Capacity
					}

					// Add the mount point if available (sometimes it contains spaces)
					if len(fields) > 5 {
						mountPoint := strings.Join(fields[5:], " ")
						adjustedFields = append(adjustedFields, mountPoint)
					} else {
						adjustedFields = append(adjustedFields, "?") // Unknown mount point
					}

					addDiskRowToGrid(grid, rowIndex, adjustedFields)
					rowIndex++
				}
			}
		}
	}

	// Now that the grid is fully built, schedule its attachment to the UI
	// This will run on the UI thread
	gtk4go.RunOnUIThread(func() {
		updateDiskDisplay(grid)
	})
}

// createEmptyDiskGrid creates a grid with just headers and a message
func createEmptyDiskGrid(message string) *gtk4.Grid {
	grid := gtk4.NewGrid(
		gtk4.WithRowSpacing(4),
		gtk4.WithColumnSpacing(12),
		gtk4.WithRowHomogeneous(false),
	)
	grid.AddCssClass("disk-info-grid")

	// Add headers
	headerLabels := []string{"Device", "Size", "Used", "Avail", "Use%", "Mount Point"}
	for i, header := range headerLabels {
		label := gtk4.NewLabel(header)
		label.AddCssClass("disk-header")
		label.SetHExpand(true)
		grid.Attach(label, i, 0, 1, 1)
	}

	// Add message spanning all columns
	messageLabel := gtk4.NewLabel(message)
	if message == "Error getting disk information" {
		messageLabel.AddCssClass("disk-info-error")
	} else {
		messageLabel.AddCssClass("disk-info-message")
	}
	messageLabel.SetHExpand(true)
	grid.Attach(messageLabel, 0, 1, 6, 1)

	return grid
}

// createDiskGridWithHeaders creates a new grid with headers and separator
func createDiskGridWithHeaders() *gtk4.Grid {
	grid := gtk4.NewGrid(
		gtk4.WithRowSpacing(4),
		gtk4.WithColumnSpacing(12),
		gtk4.WithRowHomogeneous(false),
	)
	grid.AddCssClass("disk-info-grid")

	// Add column headers
	headerLabels := []string{"Device", "Size", "Used", "Avail", "Use%", "Mount Point"}
	for i, header := range headerLabels {
		label := gtk4.NewLabel(header)
		label.AddCssClass("disk-header")
		label.SetHExpand(true)
		grid.Attach(label, i, 0, 1, 1)
	}

	// Add a separator row
	for i := 0; i < len(headerLabels); i++ {
		separator := gtk4.NewLabel("--------")
		separator.AddCssClass("disk-separator")
		separator.SetHExpand(true)
		grid.Attach(separator, i, 1, 1, 1)
	}

	return grid
}

// addDiskRowToGrid adds a row of disk information to the grid with tooltips for truncated values
func addDiskRowToGrid(grid *gtk4.Grid, rowIndex int, fields []string) {
	// Create device label with potential tooltip
	device := fields[0]
	deviceLabel := gtk4.NewLabel(device)
	deviceLabel.AddCssClass("disk-device")
	deviceLabel.SetHExpand(true)

	// Add tooltip if device name is long
	if len(device) > 16 {
		// Store full device name before truncating for display
		fullDevice := device
		// Truncate displayed text
		deviceLabel.SetText(device[:13] + "...")
		// Add tooltip with full device name
		deviceLabel.SetTooltipText(fullDevice)
	}

	grid.Attach(deviceLabel, 0, rowIndex, 1, 1)

	// Size column
	sizeLabel := gtk4.NewLabel(fields[1])
	sizeLabel.AddCssClass("disk-size")
	grid.Attach(sizeLabel, 1, rowIndex, 1, 1)

	// Used column
	usedLabel := gtk4.NewLabel(fields[2])
	usedLabel.AddCssClass("disk-used")
	grid.Attach(usedLabel, 2, rowIndex, 1, 1)

	// Available column
	availLabel := gtk4.NewLabel(fields[3])
	availLabel.AddCssClass("disk-avail")
	grid.Attach(availLabel, 3, rowIndex, 1, 1)

	// Percent column with color coding
	percentLabel := gtk4.NewLabel(fields[4])
	percentLabel.AddCssClass("disk-percent")

	// Add color coding based on usage percentage
	percentValue := strings.TrimSuffix(fields[4], "%")
	if percentVal, err := strconv.Atoi(percentValue); err == nil {
		if percentVal >= 90 {
			percentLabel.AddCssClass("disk-usage-critical")
		} else if percentVal >= 75 {
			percentLabel.AddCssClass("disk-usage-warning")
		} else {
			percentLabel.AddCssClass("disk-usage-normal")
		}
	}

	grid.Attach(percentLabel, 4, rowIndex, 1, 1)

	// Mount point column with potential tooltip
	mountPoint := ""
	if len(fields) > 5 {
		mountPoint = strings.Join(fields[5:], " ")
	}
	mountLabel := gtk4.NewLabel(mountPoint)
	mountLabel.AddCssClass("disk-mount")

	// Add tooltip if mount path is long
	if len(mountPoint) > 20 {
		// Store full path
		fullMount := mountPoint
		// Truncate displayed text
		mountLabel.SetText(mountPoint[:17] + "...")
		// Add tooltip with full path
		mountLabel.SetTooltipText(fullMount)
	}

	grid.Attach(mountLabel, 5, rowIndex, 1, 1)
}

// updateDiskDisplay updates the display with the new grid
func updateDiskDisplay(newGrid *gtk4.Grid) {
	// Lock the UI mutex to prevent concurrent UI modifications
	uiMutex.Lock()
	defer uiMutex.Unlock()

	// Check if diskCard is still valid
	if diskCard == nil {
		return
	}

	// First, remove any existing child from the diskCard
	if currentGrid != nil {
		diskCard.Remove(currentGrid)
	}

	// Add the new grid and set it as the current grid
	diskCard.Append(newGrid)
	currentGrid = newGrid

	// For backward compatibility, update the text label if it exists
	if infoLabel, ok := diskLabels.labels["disk_info"]; ok && infoLabel != nil {
		infoLabel.SetText("Information now displayed in grid format.")
	}
}

// parseMemoryValue parses a memory value string like "1024.00M" and returns the
// value and multiplier (e.g., 1024.00 and 1048576 for MB to bytes)
func parseMemoryValue(valueStr string) (uint64, uint64) {
	var multiplier uint64 = 1

	// Check for unit suffix
	if strings.HasSuffix(valueStr, "K") {
		multiplier = 1024
		valueStr = valueStr[:len(valueStr)-1]
	} else if strings.HasSuffix(valueStr, "M") {
		multiplier = 1024 * 1024
		valueStr = valueStr[:len(valueStr)-1]
	} else if strings.HasSuffix(valueStr, "G") {
		multiplier = 1024 * 1024 * 1024
		valueStr = valueStr[:len(valueStr)-1]
	}

	value, _ := strconv.ParseFloat(valueStr, 64)
	return uint64(value), multiplier
}
