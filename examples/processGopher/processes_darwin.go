package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// getSystemMemoryInfo gets the total and free memory on macOS
func getSystemMemoryInfo() (total int64, free int64, err error) {
	// Get total physical memory
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get total memory: %v", err)
	}
	
	totalStr := strings.TrimSpace(string(output))
	total, err = strconv.ParseInt(totalStr, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse total memory: %v", err)
	}

	// Get memory usage using vm_stat command
	cmd = exec.Command("vm_stat")
	output, err = cmd.Output()
	if err != nil {
		return total, 0, fmt.Errorf("failed to get memory stats: %v", err)
	}

	// Parse vm_stat output to calculate free memory
	lines := strings.Split(string(output), "\n")
	pageSize := int64(4096) // Default page size on macOS is 4KB
	freePages := int64(0)
	
	for _, line := range lines {
		if strings.Contains(line, "Pages free:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				freeStr := strings.TrimSpace(parts[1])
				freeStr = strings.ReplaceAll(freeStr, ".", "")
				pages, err := strconv.ParseInt(freeStr, 10, 64)
				if err == nil {
					freePages += pages
				}
			}
		} else if strings.Contains(line, "Pages inactive:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				inactiveStr := strings.TrimSpace(parts[1])
				inactiveStr = strings.ReplaceAll(inactiveStr, ".", "")
				pages, err := strconv.ParseInt(inactiveStr, 10, 64)
				if err == nil {
					freePages += pages
				}
			}
		}
	}

	// Calculate free memory in bytes
	free = freePages * pageSize
	return total, free, nil
}

// getCPUUsage gets the current CPU usage on macOS
func getCPUUsage() (float64, error) {
	// Run top command to get CPU usage
	cmd := exec.Command("top", "-l", "1", "-n", "0")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get CPU usage: %v", err)
	}

	// Parse top output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "CPU usage:") {
			// Extract the user CPU percentage
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "user," && i > 0 {
					// User CPU percentage should be right before "user,"
					percentStr := strings.TrimSuffix(fields[i-1], "%")
					percent, err := strconv.ParseFloat(percentStr, 64)
					if err == nil {
						return percent, nil
					}
				}
			}
		}
	}

	return 0, fmt.Errorf("could not parse CPU usage from top output")
}

// getProcessDetails gets additional details about a process on macOS
func getProcessDetails(pid int64) (map[string]string, error) {
	details := make(map[string]string)

	// Get process info using ps
	cmd := exec.Command("ps", "-p", strconv.FormatInt(pid, 10), "-o", "pid,ppid,pri,nice,args")
	output, err := cmd.Output()
	if err != nil {
		return details, fmt.Errorf("failed to get process details: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return details, fmt.Errorf("no process info found")
	}

	// Parse the header and data
	headerLine := lines[0]
	dataLine := lines[1]

	headers := strings.Fields(headerLine)
	data := strings.Fields(dataLine)

	// Match headers to data
	for i, header := range headers {
		if i < len(data) {
			details[header] = data[i]
		}
	}

	// If the last field is the command, handle it specially to get the full command
	if len(headers) > 0 && headers[len(headers)-1] == "ARGS" {
		cmdIndex := len(headers) - 1
		if cmdIndex < len(data) {
			// Reconstruct the full command
			cmdStart := strings.Index(dataLine, data[cmdIndex])
			if cmdStart >= 0 {
				details["ARGS"] = dataLine[cmdStart:]
			}
		}
	}

	return details, nil
}
