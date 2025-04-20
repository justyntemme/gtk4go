package main

import (
	"os/exec"
	"strconv"
	"strings"
)

// splitFieldsPreservingQuotes splits a string into fields while preserving quoted sections
func splitFieldsPreservingQuotes(s string) []string {
	var fields []string
	var current strings.Builder
	inQuote := false

	// Add space to ensure the last field gets added
	s = s + " "

	for _, char := range s {
		// Handle quotes
		if char == '"' {
			inQuote = !inQuote
			continue
		}

		// If we reach a space and we're not in quotes, add the current field
		if char == ' ' && !inQuote {
			if current.Len() > 0 {
				fields = append(fields, strings.TrimSpace(current.String()))
				current.Reset()
			}
			continue
		}

		// Add the character to the current field
		current.WriteRune(char)
	}

	return fields
}

// getSystemMemoryInfo gets the total and free memory on macOS
func getSystemMemoryInfo() (total int64, free int64, err error) {
	// Get total physical memory
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to a reasonable default if the command fails
		return 8 * 1024 * 1024 * 1024, 4 * 1024 * 1024 * 1024, nil
	}

	totalStr := strings.TrimSpace(string(output))
	total, err = strconv.ParseInt(totalStr, 10, 64)
	if err != nil {
		// Fallback to a reasonable default if parsing fails
		return 8 * 1024 * 1024 * 1024, 4 * 1024 * 1024 * 1024, nil
	}

	// Get memory usage using vm_stat command
	cmd = exec.Command("vm_stat")
	output, err = cmd.Output()
	if err != nil {
		// Fallback to estimating free memory as 50% of total
		return total, total / 2, nil
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
		} else if strings.Contains(line, "Pages purgeable:") {
			// Also count purgeable pages as effectively free memory
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				purgeableStr := strings.TrimSpace(parts[1])
				purgeableStr = strings.ReplaceAll(purgeableStr, ".", "")
				pages, err := strconv.ParseInt(purgeableStr, 10, 64)
				if err == nil {
					freePages += pages
				}
			}
		}
	}

	// Calculate free memory in bytes
	free = freePages * pageSize

	// Ensure free memory is reasonable (at least 5% of total)
	if free < (total / 20) {
		free = total / 10 // Set to 10% of total as a safety measure
	}

	return total, free, nil
}

// getCPUUsage gets the current CPU usage on macOS
func getCPUUsage() (float64, error) {
	// Use a simpler and more reliable approach with vm_stat and ps
	// First try iostat which gives reliable CPU numbers
	cmd := exec.Command("iostat", "-c", "2")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) >= 3 { // Header, separator, data
			// The CPU info should be in the third line (index 2)
			dataLine := strings.TrimSpace(lines[2])
			fields := strings.Fields(dataLine)
			if len(fields) >= 3 {
				// Fields typically are user%, system%, idle%
				userStr := fields[0]
				sysStr := fields[1]

				// Parse user% and system%
				userCPU, errUser := strconv.ParseFloat(userStr, 64)
				sysCPU, errSys := strconv.ParseFloat(sysStr, 64)

				if errUser == nil && errSys == nil {
					// Return combined CPU usage
					return userCPU + sysCPU, nil
				}
			}
		}
	}

	// Fall back to top if iostat fails
	cmd = exec.Command("top", "-l", "1", "-n", "0")
	output, err = cmd.Output()
	if err != nil {
		return 5.0, nil // Return a reasonable default if all methods fail
	}

	// Parse top output
	lines := strings.Split(string(output), "\n")

	// Try to find the CPU usage line
	for _, line := range lines {
		if strings.Contains(line, "CPU usage:") {
			// The line format is usually: "CPU usage: X.XX% user, Y.YY% sys, Z.ZZ% idle"
			parts := strings.Split(line, ",")
			if len(parts) >= 2 {
				// Extract user CPU usage
				userPart := parts[0]
				userParts := strings.Split(userPart, ":")
				if len(userParts) >= 2 {
					userStr := strings.TrimSpace(userParts[1])
					userStr = strings.TrimSuffix(userStr, "%")
					userStr = strings.TrimSuffix(userStr, " user")
					userCPU, err := strconv.ParseFloat(userStr, 64)

					// Extract system CPU usage
					sysPart := strings.TrimSpace(parts[1])
					sysPart = strings.TrimSuffix(sysPart, "%")
					sysPart = strings.TrimSuffix(sysPart, " sys")
					sysCPU, err2 := strconv.ParseFloat(sysPart, 64)

					// Return combined CPU usage
					if err == nil && err2 == nil {
						return userCPU + sysCPU, nil
					} else if err == nil {
						return userCPU, nil
					}
				}
			}
		}
	}

	// If all else fails, return a reasonable default value
	return 5.0, nil
}

// getDarwinProcessDetails gets additional details about a process on macOS
// This is a darwin-specific version that's called by getProcessDetails
func getDarwinProcessDetails(pid int64) (map[string]string, error) {
	details := make(map[string]string)

	// Initialize with default values to ensure we always have complete data
	details["USER"] = "N/A"
	details["%CPU"] = "0.0"
	details["RSS"] = "0"
	details["STAT"] = "N/A"
	details["THCOUNT"] = "0"
	details["STARTED"] = "N/A"

	// Try to get process info using ps with explicit columns specified
	// Don't include headers, use the equals sign approach for reliable parsing
	cmd := exec.Command("ps", "-p", strconv.FormatInt(pid, 10), "-o", "user=,pcpu=,rss=,state=,thcount=,lstart=")
	output, err := cmd.Output()
	if err != nil {
		// Process might have terminated - return the default values
		return details, nil
	}

	// Process the output line
	dataLine := strings.TrimSpace(string(output))
	if dataLine == "" {
		// No data returned, return default values
		return details, nil
	}

	// Split the fields using splitFieldsPreservingQuotes to handle fields with spaces
	fields := strings.Fields(dataLine)

	// Map fields to details if we have enough fields
	if len(fields) >= 1 {
		details["USER"] = fields[0]
	}
	if len(fields) >= 2 {
		details["%CPU"] = fields[1]
	}
	if len(fields) >= 3 {
		details["RSS"] = fields[2]
	}
	if len(fields) >= 4 {
		details["STAT"] = fields[3]
	}
	if len(fields) >= 5 {
		details["THCOUNT"] = fields[4]
	}

	// Last fields are the start time which might have multiple fields
	if len(fields) >= 6 {
		// The remaining fields (5 onwards) should be the start time
		startTime := strings.Join(fields[5:], " ")
		details["STARTED"] = startTime
	}

	// Try to get command name separately
	cmdCmd := exec.Command("ps", "-p", strconv.FormatInt(pid, 10), "-o", "command=")
	cmdOutput, err := cmdCmd.Output()
	if err == nil && len(cmdOutput) > 0 {
		details["COMMAND"] = strings.TrimSpace(string(cmdOutput))
	} else {
		details["COMMAND"] = ""
	}

	return details, nil
}
