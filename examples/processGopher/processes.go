package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	
	"github.com/mitchellh/go-ps"
)

// ProcessInfo contains information about a process
type ProcessInfo struct {
	PID         int64
	Name        string
	Username    string
	CPUPercent  float64
	MemoryBytes int64
	Threads     int
	State       string
	StartTime   string
}

// getProcesses returns a list of running processes using the go-ps package
func getProcesses() ([]ProcessInfo, error) {
	// Get the list of processes using go-ps
	processes, err := ps.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get process list: %v", err)
	}

	// Convert to our ProcessInfo format
	result := make([]ProcessInfo, 0, len(processes))
	for _, proc := range processes {
		// Get additional process details based on platform
		cpuPercent := 0.0
		memBytes := int64(0)
		threads := 0
		state := ""
		username := ""
		startTime := ""

		// Get platform-specific process details
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			// On Unix systems, we can use ps to get additional details
			details, err := getProcessDetails(int64(proc.Pid()))
			if err == nil {
				username = details["USER"]
				state = details["STAT"]
				
				// Try to parse CPU percentage
				if cpuStr, ok := details["%CPU"]; ok {
					cpuPercent, _ = strconv.ParseFloat(cpuStr, 64)
				}
				
				// Try to parse memory usage
				if memStr, ok := details["RSS"]; ok {
					memKB, _ := strconv.ParseInt(memStr, 10, 64)
					memBytes = memKB * 1024 // Convert KB to bytes
				}
				
				// Try to parse thread count
				if threadStr, ok := details["THCOUNT"]; ok {
					threads, _ = strconv.Atoi(threadStr)
				} else if threadStr, ok := details["NLWP"]; ok {
					// On Linux, thread count is in NLWP
					threads, _ = strconv.Atoi(threadStr)
				}
				
				// Get start time
				if timeStr, ok := details["STARTED"]; ok {
					startTime = timeStr
				}
			}
		}

		// Create the process info
		info := ProcessInfo{
			PID:         int64(proc.Pid()),
			Name:        proc.Executable(),
			Username:    username,
			CPUPercent:  cpuPercent,
			MemoryBytes: memBytes,
			Threads:     threads,
			State:       state,
			StartTime:   startTime,
		}

		result = append(result, info)
	}

	return result, nil
}

// getProcessDetails uses ps to get additional details about a process
func getProcessDetails(pid int64) (map[string]string, error) {
	details := make(map[string]string)
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// For macOS, use the Darwin-specific function
		return getDarwinProcessDetails(pid)
	case "linux":
		// For Linux
		cmd = exec.Command("ps", "-p", strconv.FormatInt(pid, 10), "-o", "pid,user,%cpu,rss,stat,nlwp,start")
	default:
		// Generic fallback
		cmd = exec.Command("ps", "-p", strconv.FormatInt(pid, 10), "-o", "pid,user,%cpu,rss,state")
	}

	output, err := cmd.Output()
	if err != nil {
		return details, fmt.Errorf("failed to get process details: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return details, fmt.Errorf("no process info found")
	}

	// Parse the header and data
	headerLine := strings.TrimSpace(lines[0])
	dataLine := strings.TrimSpace(lines[1])

	// Split header and data into fields
	headerFields := strings.Fields(headerLine)
	dataFields := strings.Fields(dataLine)

	// Map headers to data
	for i, header := range headerFields {
		if i < len(dataFields) {
			details[header] = dataFields[i]
		}
	}

	return details, nil
}

// formatStartTime converts the ps start time format to a more readable format
func formatStartTime(timeStr string) string {
	// Parse time string like "Mon Jan 2 15:04:05 2006"
	t, err := time.Parse("Mon Jan 2 15:04:05 2006", timeStr)
	if err != nil {
		return timeStr // Return original if parsing fails
	}

	// Format to a more readable format
	return t.Format("2006-01-02 15:04:05")
}

// terminateProcess attempts to terminate a process by sending a SIGTERM signal
func terminateProcess(pid int64) error {
	cmd := exec.Command("kill", strconv.FormatInt(pid, 10))
	return cmd.Run()
}

// ps -eo "pid,user,%cpu,rss,state,thcount,lstart,command"
