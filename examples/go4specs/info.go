package main

import (
	"../../../gtk4go"
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
