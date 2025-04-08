//go:build linux
// +build linux

// Package main provides Linux-specific system information functionality
// File: info_linux.go - Contains Linux-specific implementations
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

// readUptime reads the system uptime on Linux
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

// readCPUModel reads the CPU model on Linux
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

// getCPUCount returns the number of physical cores and threads on Linux
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

// readCPUFrequency reads the CPU frequency on Linux
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

// getCPUUsage gets the CPU usage percentage on Linux
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

// getMemoryInfo gets memory information (total, used, free) on Linux
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

// getSwapInfo gets swap information (total, used, free) on Linux
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

// refreshGPUInfo updates the GPU information labels for Linux
func refreshGPUInfo(labels *labelMap) {
	// Helper function to truncate long text
	truncateText := func(text string, maxLength int) string {
		if len(text) > maxLength {
			return text[:maxLength-3] + "..."
		}
		return text
	}

	// First try to use lshw to get GPU information
	lshwSuccess := false
	if _, err := executeCommand("which", "lshw"); err == nil {
		// Run lshw to get display information
		lshwOutput, err := executeCommand("lshw", "-C", "Display")
		if err == nil && len(lshwOutput) > 0 {
			// Parse the lshw output
			lines := strings.Split(lshwOutput, "\n")

			// Variables to store extracted information
			var product, vendor, driver, resolution, memory string

			// Parse each line for relevant information
			for _, line := range lines {
				line = strings.TrimSpace(line)

				if strings.Contains(line, "product:") && len(strings.Split(line, "product:")) > 1 {
					product = strings.TrimSpace(strings.Split(line, "product:")[1])
				} else if strings.Contains(line, "vendor:") && len(strings.Split(line, "vendor:")) > 1 {
					vendor = strings.TrimSpace(strings.Split(line, "vendor:")[1])
				} else if strings.Contains(line, "configuration:") && len(strings.Split(line, "configuration:")) > 1 {
					config := strings.TrimSpace(strings.Split(line, "configuration:")[1])

					// Extract driver and resolution from configuration
					configParts := strings.Split(config, " ")
					for _, part := range configParts {
						if strings.HasPrefix(part, "driver=") {
							driver = strings.TrimPrefix(part, "driver=")
						} else if strings.HasPrefix(part, "resolution=") {
							resolution = strings.TrimPrefix(part, "resolution=")
						}
					}
				} else if strings.Contains(line, "memory:") && len(strings.Split(line, "memory:")) > 1 {
					// Just take the first memory entry as an indication
					if memory == "" {
						parts := strings.Split(line, "memory:")
						memory = strings.TrimSpace(parts[1])
						// Cut off at first space if there are multiple entries
						if spaceIdx := strings.Index(memory, " "); spaceIdx != -1 {
							memory = memory[:spaceIdx]
						}
					}
				}
			}

			// Update labels with the information we found
			if product != "" {
				displayText := truncateText(product, 35)
				labels.update("gpu_model", displayText)

				// Add tooltip for full text if truncated
				if label, ok := labels.labels["gpu_model"]; ok && len(product) > 35 {
					label.SetTooltipText(product)
				}

				lshwSuccess = true
			}

			if vendor != "" {
				displayText := truncateText(vendor, 30)
				labels.update("gpu_vendor", displayText)

				// Add tooltip for full text if truncated
				if label, ok := labels.labels["gpu_vendor"]; ok && len(vendor) > 30 {
					label.SetTooltipText(vendor)
				}
			}

			if driver != "" {
				// displayText := truncateText(driver, 30)
				displayText := strings.Fields(driver)[0]
				labels.update("gpu_driver", displayText)

				// Add tooltip for full text if truncated
				if label, ok := labels.labels["gpu_driver"]; ok && len(driver) > 30 {
					label.SetTooltipText(driver)
				}
			}

			// If we have resolution information, use it for renderer
			if resolution != "" {
				displayText := truncateText("Resolution: "+resolution, 30)
				labels.update("gpu_renderer", displayText)

				// Add tooltip for full text if truncated
				if label, ok := labels.labels["gpu_renderer"]; ok && len(resolution) > 25 {
					label.SetTooltipText("Resolution: " + resolution)
				}
			}

			// Try to use memory information if available
			if memory != "" {
				displayText := truncateText("Memory: "+memory, 30)
				labels.update("gpu_memory", displayText)

				// Add tooltip for full text if truncated
				if label, ok := labels.labels["gpu_memory"]; ok && len(memory) > 25 {
					label.SetTooltipText("Memory: " + memory)
				}
			}
		}
	}

	// Check which labels we need to populate with fallback methods
	// Only use fallback methods for GPU model if lshw wasn't successful
	modelLabel, hasModel := labels.labels["gpu_model"]
	modelEmpty := (!hasModel || modelLabel == nil || modelLabel.GetText() == "") && !lshwSuccess

	vendorLabel, hasVendor := labels.labels["gpu_vendor"]
	vendorEmpty := !hasVendor || vendorLabel == nil || vendorLabel.GetText() == ""

	rendererLabel, hasRenderer := labels.labels["gpu_renderer"]
	rendererEmpty := !hasRenderer || rendererLabel == nil || rendererLabel.GetText() == ""

	driverLabel, hasDriver := labels.labels["gpu_driver"]
	driverEmpty := !hasDriver || driverLabel == nil || driverLabel.GetText() == ""

	memoryLabel, hasMemory := labels.labels["gpu_memory"]
	memoryEmpty := !hasMemory || memoryLabel == nil || memoryLabel.GetText() == ""

	// Fall back to lspci for GPU model if needed
	if modelEmpty {
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
	}

	// Try to get OpenGL information using glxinfo
	if _, err := executeCommand("which", "glxinfo"); err == nil {
		// Only get vendor if needed
		if vendorEmpty {
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
		}

		// Only get renderer if needed
		if rendererEmpty {
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
		}

		// Always try to get OpenGL version
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
				// Only update driver if needed
				if driverEmpty {
					driverText := "NVIDIA " + strings.TrimSpace(parts[1])
					displayDriver := truncateText(driverText, 30)
					labels.update("gpu_driver", displayDriver)

					// Add tooltip for full text if truncated
					if label, ok := labels.labels["gpu_driver"]; ok && len(driverText) > 30 {
						label.SetTooltipText(driverText)
					}
				}

				// Only update memory if needed
				if memoryEmpty {
					memoryText := strings.TrimSpace(parts[2])
					displayMemory := truncateText(memoryText, 30)
					labels.update("gpu_memory", displayMemory)

					// Add tooltip for full text if truncated
					if label, ok := labels.labels["gpu_memory"]; ok && len(memoryText) > 30 {
						label.SetTooltipText(memoryText)
					}
				}

				// Always update utilization as it's real-time data
				utilizationText := strings.TrimSpace(parts[3])
				displayUtil := truncateText(utilizationText, 30)
				labels.update("gpu_utilization", displayUtil)

				// Add tooltip for full text if truncated
				if label, ok := labels.labels["gpu_utilization"]; ok && len(utilizationText) > 30 {
					label.SetTooltipText(utilizationText)
				}
			}
		}
	} else if driverEmpty {
		// Try to get driver info from lspci if needed and nvidia-smi is not available
		if _, err := executeCommand("which", "lspci"); err == nil {
			driverCmd := "lspci -v | grep -A10 -i 'vga\\|3d' | grep 'Kernel driver in use'"
			if driver, err := executeCommand("bash", "-c", driverCmd); err == nil {
				parts := strings.SplitN(driver, ":", 2)
				if len(parts) >= 2 {
					driverText := strings.TrimSpace(parts[1])
					// displayText := truncateText(driverText, 30)
					displayText := strings.Fields(driverText)[0]
					labels.update("gpu_driver", displayText)

					// Add tooltip for full text if truncated
					if label, ok := labels.labels["gpu_driver"]; ok && len(driverText) > 30 {
						label.SetTooltipText(driverText)
					}
				}
			} else {
				labels.update("gpu_driver", "Unknown")
			}
		}
	}
}