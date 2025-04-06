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

	// Update Distribution
	if fileExists("/etc/os-release") {
		if dist, err := readDistribution(); err == nil {
			text := dist
			labels.update("distribution", text)
			
			// Add tooltip for potentially long values
			if label, ok := labels.labels["distribution"]; ok && len(text) > 20 {
				label.SetTooltipText(text)
			}
		}
	} else {
		labels.update("distribution", "Unknown (os-release not found)")
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
		text := model
		labels.update("cpu_model", text)
		
		// CPU model strings are often very long, always add tooltip
		if label, ok := labels.labels["cpu_model"]; ok {
			label.SetTooltipText(text)
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

// refreshGPUInfo updates the GPU information labels
func refreshGPUInfo(labels *labelMap) {
	// Helper function to truncate long text
	truncateText := func(text string, maxLength int) string {
		if len(text) > maxLength {
			return text[:maxLength-3] + "..."
		}
		return text
	}

	// Try to get GPU information using lspci
	if _, err := executeCommand("which", "lspci"); err == nil {
		// Extract GPU info using grep
		gpuLines, err := executeCommand("bash", "-c", "lspci | grep -i 'vga\\|3d\\|2d'")
		if err == nil && len(gpuLines) > 0 {
			// Set primary GPU model
			lines := strings.Split(gpuLines, "\n")
			if len(lines) > 0 {
				// Extract GPU name from the first line
				parts := strings.SplitN(lines[0], ":", 2)
				if len(parts) >= 2 {
					model := strings.TrimSpace(parts[1])
					displayText := truncateText(model, 35)
					labels.update("gpu_model", displayText)
					
					// Always add tooltip for GPU model as they're typically long
					if label, ok := labels.labels["gpu_model"]; ok && len(model) > 35 {
						label.SetTooltipText(model) // Show full text in tooltip
					}
				}
			}
		} else {
			labels.update("gpu_model", "No dedicated GPU detected")
		}
	} else {
		labels.update("gpu_model", "GPU detection not available (lspci not found)")
	}

	// Try to get OpenGL information using glxinfo
	if _, err := executeCommand("which", "glxinfo"); err == nil {
		// Extract OpenGL vendor
		vendorCmd := "glxinfo | grep 'OpenGL vendor'"
		if vendor, err := executeCommand("bash", "-c", vendorCmd); err == nil {
			parts := strings.SplitN(vendor, ":", 2)
			if len(parts) >= 2 {
				vendorText := strings.TrimSpace(parts[1])
				displayText := truncateText(vendorText, 30)
				labels.update("gpu_vendor", displayText)
				
				// Add tooltip for full text if truncated
				if label, ok := labels.labels["gpu_vendor"]; ok && len(vendorText) > 30 {
					label.SetTooltipText(vendorText)
				}
			}
		}

		// Extract OpenGL renderer
		rendererCmd := "glxinfo | grep 'OpenGL renderer'"
		if renderer, err := executeCommand("bash", "-c", rendererCmd); err == nil {
			parts := strings.SplitN(renderer, ":", 2)
			if len(parts) >= 2 {
				rendererText := strings.TrimSpace(parts[1])
				displayText := truncateText(rendererText, 30)
				labels.update("gpu_renderer", displayText)
				
				// Add tooltip for full text if truncated
				if label, ok := labels.labels["gpu_renderer"]; ok && len(rendererText) > 30 {
					label.SetTooltipText(rendererText)
				}
			}
		}

		// Extract OpenGL version
		versionCmd := "glxinfo | grep 'OpenGL version'"
		if version, err := executeCommand("bash", "-c", versionCmd); err == nil {
			parts := strings.SplitN(version, ":", 2)
			if len(parts) >= 2 {
				versionText := strings.TrimSpace(parts[1])
				displayText := truncateText(versionText, 30)
				labels.update("gpu_gl_version", displayText)
				
				// Add tooltip for full text if truncated
				if label, ok := labels.labels["gpu_gl_version"]; ok && len(versionText) > 30 {
					label.SetTooltipText(versionText)
				}
			}
		}
	} else {
		labels.update("gpu_gl_version", "OpenGL info not available (glxinfo not found)")
	}

	// Try to get NVIDIA-specific information if available
	if _, err := executeCommand("which", "nvidia-smi"); err == nil {
		// NVIDIA GPU detected, get additional info
		if nvInfo, err := executeCommand("nvidia-smi", "--query-gpu=name,driver_version,memory.total,utilization.gpu", "--format=csv,noheader"); err == nil {
			parts := strings.Split(nvInfo, ",")
			if len(parts) >= 4 {
				driverText := "NVIDIA " + strings.TrimSpace(parts[1])
				memoryText := strings.TrimSpace(parts[2])
				utilizationText := strings.TrimSpace(parts[3])
				
				displayDriver := truncateText(driverText, 30)
				displayMemory := truncateText(memoryText, 30)
				displayUtil := truncateText(utilizationText, 30)
				
				labels.update("gpu_driver", displayDriver)
				labels.update("gpu_memory", displayMemory)
				labels.update("gpu_utilization", displayUtil)
				
				// Add tooltips for truncated values
				if label, ok := labels.labels["gpu_driver"]; ok && len(driverText) > 30 {
					label.SetTooltipText(driverText)
				}
				if label, ok := labels.labels["gpu_memory"]; ok && len(memoryText) > 30 {
					label.SetTooltipText(memoryText)
				}
				if label, ok := labels.labels["gpu_utilization"]; ok && len(utilizationText) > 30 {
					label.SetTooltipText(utilizationText)
				}
			}
		}
	} else {
		// Try to get driver info from lspci
		if _, err := executeCommand("which", "lspci"); err == nil {
			driverCmd := "lspci -v | grep -A10 -i 'vga\\|3d' | grep 'Kernel driver in use'"
			if driver, err := executeCommand("bash", "-c", driverCmd); err == nil {
				parts := strings.SplitN(driver, ":", 2)
				if len(parts) >= 2 {
					driverText := strings.TrimSpace(parts[1])
					displayText := truncateText(driverText, 30)
					labels.update("gpu_driver", displayText)
					
					// Add tooltip for full text if truncated
					if label, ok := labels.labels["gpu_driver"]; ok && len(driverText) > 30 {
						label.SetTooltipText(driverText)
					}
				}
			}
		} else {
			labels.update("gpu_driver", "Unknown")
		}
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
			for i, line := range lines {
				if i == 0 || len(strings.TrimSpace(line)) == 0 {
					// Skip header and empty lines
					continue
				}

				fields := strings.Fields(line)
				if len(fields) >= 6 {
					// Add this disk entry to the grid with tooltips
					addDiskRowToGrid(grid, rowIndex, fields)
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
		grid.Attach(label, i, 0, 1, 1)
	}

	// Add message spanning all columns
	messageLabel := gtk4.NewLabel(message)
	if message == "Error getting disk information" {
		messageLabel.AddCssClass("disk-info-error")
	} else {
		messageLabel.AddCssClass("disk-info-message")
	}
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
		grid.Attach(label, i, 0, 1, 1)
	}

	// Add a separator row
	for i := 0; i < len(headerLabels); i++ {
		separator := gtk4.NewLabel("--------")
		separator.AddCssClass("disk-separator")
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
	mountPoint := strings.Join(fields[5:], " ")
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