//go:build darwin
// +build darwin

// Package main provides macOS-specific system information functionality
// File: info_darwin.go - Contains macOS-specific implementations
package main

import (
	"bufio"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// readDistribution reads the macOS version information
func readDistribution() (string, error) {
	// Get macOS product name
	productName, err := executeCommand("sw_vers", "-productName")
	if err != nil {
		return "", err
	}

	// Get macOS version
	productVersion, err := executeCommand("sw_vers", "-productVersion")
	if err != nil {
		return "", err
	}

	// Get macOS build
	buildVersion, err := executeCommand("sw_vers", "-buildVersion")
	if err != nil {
		return "", err
	}

	// Clean up and combine the results
	productName = strings.TrimSpace(productName)
	productVersion = strings.TrimSpace(productVersion)
	buildVersion = strings.TrimSpace(buildVersion)

	return fmt.Sprintf("%s %s (%s)", productName, productVersion, buildVersion), nil
}

// readUptime reads the system uptime on macOS
func readUptime() (string, error) {
	uptimeOutput, err := executeCommand("uptime")
	if err != nil {
		return "", err
	}

	// Parse the uptime output, which looks like:
	// "9:45  up 10 days,  2:14, 5 users, load averages: 1.67 2.01 2.16"
	parts := strings.Split(uptimeOutput, "up ")
	if len(parts) < 2 {
		return "Unknown", fmt.Errorf("unexpected uptime format")
	}

	uptimePart := strings.Split(parts[1], ",")[0]
	return strings.TrimSpace(uptimePart), nil
}

// readCPUModel reads the CPU model on macOS
func readCPUModel() (string, error) {
	model, err := executeCommand("sysctl", "-n", "machdep.cpu.brand_string")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(model), nil
}

// getCPUCount returns the number of physical cores and threads on macOS
func getCPUCount() (int, int) {
	// Get physical CPU count
	physicalOutput, err := executeCommand("sysctl", "-n", "hw.physicalcpu")
	physicals := 1 // Default to 1 if we can't determine
	if err == nil {
		if count, err := strconv.Atoi(strings.TrimSpace(physicalOutput)); err == nil {
			physicals = count
		}
	}

	// Get logical CPU count
	logicalOutput, err := executeCommand("sysctl", "-n", "hw.logicalcpu")
	threads := runtime.NumCPU() // Use Go's runtime as fallback
	if err == nil {
		if count, err := strconv.Atoi(strings.TrimSpace(logicalOutput)); err == nil {
			threads = count
		}
	}

	return physicals, threads
}

// readCPUFrequency reads the CPU frequency on macOS
func readCPUFrequency() (string, error) {
	// Get CPU frequency in Hz
	freqOutput, err := executeCommand("sysctl", "-n", "hw.cpufrequency")
	if err != nil {
		return "Unknown", err
	}

	// Convert to GHz
	freq, err := strconv.ParseInt(strings.TrimSpace(freqOutput), 10, 64)
	if err != nil {
		return "Unknown", err
	}

	freqGHz := float64(freq) / 1_000_000_000
	return fmt.Sprintf("%.2f GHz", freqGHz), nil
}

// getCPUUsage gets the CPU usage percentage on macOS
func getCPUUsage() (float64, error) {
	// On macOS, we can use top in a different way
	output, err := executeCommand("top", "-l", "1", "-n", "0")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "CPU usage") {
			// Parse line like: "CPU usage: 7.35% user, 14.39% sys, 78.25% idle"
			idleIndex := strings.Index(line, "% idle")
			if idleIndex > 0 {
				// Extract the idle percentage
				idleStart := strings.LastIndex(line[:idleIndex], " ") + 1
				idleStr := line[idleStart:idleIndex]
				idle, err := strconv.ParseFloat(idleStr, 64)
				if err == nil {
					return 100.0 - idle, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("CPU usage not found")
}

// getMemoryInfo gets memory information (total, used, free) on macOS
func getMemoryInfo() (uint64, uint64, uint64, error) {
	// Get total memory
	totalOutput, err := executeCommand("sysctl", "-n", "hw.memsize")
	if err != nil {
		return 0, 0, 0, err
	}

	total, err := strconv.ParseUint(strings.TrimSpace(totalOutput), 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	// Get memory usage using vm_stat
	vmStatOutput, err := executeCommand("vm_stat")
	if err != nil {
		return 0, 0, 0, err
	}

	// Parse vm_stat output
	pageSize := uint64(4096) // Default page size
	free := uint64(0)

	scanner := bufio.NewScanner(strings.NewReader(vmStatOutput))
	for scanner.Scan() {
		line := scanner.Text()

		// Get "Pages free"
		if strings.HasPrefix(line, "Pages free:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				valueStr := strings.TrimSpace(strings.TrimSuffix(parts[1], "."))
				value, err := strconv.ParseUint(valueStr, 10, 64)
				if err == nil {
					free += value * pageSize
				}
			}
		}

		// Also include "Pages inactive" for a better estimate of free memory
		if strings.HasPrefix(line, "Pages inactive:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				valueStr := strings.TrimSpace(strings.TrimSuffix(parts[1], "."))
				value, err := strconv.ParseUint(valueStr, 10, 64)
				if err == nil {
					free += value * pageSize
				}
			}
		}
	}

	// Calculate used memory
	used := total - free

	return total, used, free, nil
}

// getSwapInfo gets swap information (total, used, free) on macOS
func getSwapInfo() (uint64, uint64, uint64, error) {
	// Get swap info using sysctl
	output, err := executeCommand("sysctl", "-n", "vm.swapusage")
	if err != nil {
		return 0, 0, 0, err
	}

	// Parse output like: "total = 2048.00M used = 1017.75M free = 1030.25M"
	var total, used, free uint64

	parts := strings.Fields(output)
	for i, part := range parts {
		if part == "total" && i+2 < len(parts) {
			valueStr := parts[i+2]
			value, multiplier := parseMemoryValue(valueStr)
			total = value * multiplier
		} else if part == "used" && i+2 < len(parts) {
			valueStr := parts[i+2]
			value, multiplier := parseMemoryValue(valueStr)
			used = value * multiplier
		} else if part == "free" && i+2 < len(parts) {
			valueStr := parts[i+2]
			value, multiplier := parseMemoryValue(valueStr)
			free = value * multiplier
		}
	}

	return total, used, free, nil
}

// refreshGPUInfo updates the GPU information labels for macOS
func refreshGPUInfo(labels *labelMap) {
	// Use system_profiler to get GPU information on macOS
	output, err := executeCommand("system_profiler", "SPDisplaysDataType")
	if err != nil {
		labels.update("gpu_model", "Error getting GPU information")
		return
	}

	// Parse system_profiler output
	var gpuModel, gpuVendor, gpuRenderer, gpuMemory string

	lines := strings.Split(output, "\n")
	inGraphicsCard := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if we're in a Graphics section
		if strings.Contains(line, "Graphics") && strings.HasSuffix(line, ":") {
			inGraphicsCard = true
			continue
		}

		// Skip if not in a Graphics section
		if !inGraphicsCard {
			continue
		}

		// Check for various GPU properties
		if strings.Contains(line, "Chipset Model:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 {
				gpuModel = strings.TrimSpace(parts[1])
				// If model contains vendor, extract it
				if strings.Contains(gpuModel, "AMD") {
					gpuVendor = "AMD"
				} else if strings.Contains(gpuModel, "NVIDIA") {
					gpuVendor = "NVIDIA"
				} else if strings.Contains(gpuModel, "Intel") {
					gpuVendor = "Intel"
				} else if strings.Contains(gpuModel, "Apple") {
					gpuVendor = "Apple"
				}
			}
		} else if strings.Contains(line, "Vendor:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 {
				gpuVendor = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(line, "Device ID:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 {
				gpuRenderer = "Device ID: " + strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(line, "VRAM") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 {
				gpuMemory = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(line, "Metal:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 {
				if gpuRenderer == "" {
					gpuRenderer = "Metal: " + strings.TrimSpace(parts[1])
				}
			}
		} else if strings.Contains(line, "Resolution:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 && gpuRenderer == "" {
				gpuRenderer = "Resolution: " + strings.TrimSpace(parts[1])
			}
		}
	}

	// Update labels with the information we found
	if gpuModel != "" {
		labels.update("gpu_model", gpuModel)
	} else {
		labels.update("gpu_model", "Unknown GPU")
	}

	if gpuVendor != "" {
		labels.update("gpu_vendor", gpuVendor)
	}

	if gpuRenderer != "" {
		labels.update("gpu_renderer", gpuRenderer)
	}

	if gpuMemory != "" {
		labels.update("gpu_memory", gpuMemory)
	}

	// Get Metal or OpenGL driver info
	driverInfo, err := executeCommand("system_profiler", "SPDisplaysDataType", "-detailLevel", "mini")
	if err == nil {
		if strings.Contains(driverInfo, "Metal:") {
			labels.update("gpu_driver", "Metal")
		} else {
			labels.update("gpu_driver", "OpenGL")
		}
	}

	// Try to get OpenGL version
	output, err = executeCommand("defaults", "read", "/Library/Preferences/com.apple.opengl", "GLVersion")
	if err == nil {
		labels.update("gpu_gl_version", strings.TrimSpace(output))
	}
}

