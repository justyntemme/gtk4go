package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	
	"github.com/mitchellh/go-ps"
)

// Process caching to reduce system calls
var (
	processInfoCache     []ProcessInfo
	lastProcessQueryTime time.Time
	processCacheMutex    sync.RWMutex
	processCacheDuration = 1 * time.Second // Cache process info for 1 second
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
	// Use a simple process cache to avoid excessive process querying with proper locking
	processCacheMutex.RLock()
	cacheAge := time.Since(lastProcessQueryTime)
	cacheValid := cacheAge < processCacheDuration && len(processInfoCache) > 0
	processCacheMutex.RUnlock()
	
	if cacheValid {
		// Use cached data if it's recent enough and not empty
		processCacheMutex.RLock()
		cache := make([]ProcessInfo, len(processInfoCache))
		copy(cache, processInfoCache) // Make a copy to avoid race conditions
		processCacheMutex.RUnlock()
		return cache, nil
	}

	// Get the list of processes using go-ps
	processes, err := ps.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get process list: %v", err)
	}

	// Create a map of already processed PIDs to avoid duplicates
	processedPIDs := make(map[int]bool)

	// Convert to our ProcessInfo format
	result := make([]ProcessInfo, 0, len(processes))
	for _, proc := range processes {
		// Skip if we've already processed this PID
		pid := proc.Pid()
		if processedPIDs[pid] {
			continue
		}
		processedPIDs[pid] = true

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
			details, err := getProcessDetails(int64(pid))
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

		// Get process name - fallback to executable if name is empty
		name := proc.Executable()
		if name == "" {
			name = fmt.Sprintf("Process %d", pid)
		}

		// Create the process info
		info := ProcessInfo{
			PID:         int64(pid),
			Name:        name,
			Username:    username,
			CPUPercent:  cpuPercent,
			MemoryBytes: memBytes,
			Threads:     threads,
			State:       state,
			StartTime:   startTime,
		}

		result = append(result, info)
	}

	// Update cache with proper locking
	processCacheMutex.Lock()
	processInfoCache = make([]ProcessInfo, len(result))
	copy(processInfoCache, result)
	lastProcessQueryTime = time.Now()
	processCacheMutex.Unlock()

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
		// For Linux - use column names without header for more reliable parsing
		cmd = exec.Command("ps", "-p", strconv.FormatInt(pid, 10), "-o", "user=,pcpu=,rss=,stat=,nlwp=,lstart=")
	default:
		// Generic fallback
		cmd = exec.Command("ps", "-p", strconv.FormatInt(pid, 10), "-o", "user=,pcpu=,rss=,state=")
	}

	output, err := cmd.Output()
	if err != nil {
		// Initialize with default values
		details["USER"] = "N/A"
		details["%CPU"] = "0.0"
		details["RSS"] = "0"
		details["STAT"] = "N/A"
		details["NLWP"] = "0"
		details["STARTED"] = "N/A"
		return details, nil
	}

	// Process output line - since we specified '=' format, there are no headers
	dataLine := strings.TrimSpace(string(output))
	if dataLine == "" {
		// Initialize with default values
		details["USER"] = "N/A"
		details["%CPU"] = "0.0"
		details["RSS"] = "0"
		details["STAT"] = "N/A"
		details["NLWP"] = "0"
		details["STARTED"] = "N/A"
		return details, nil
	}

	// Split fields
	fields := strings.Fields(dataLine)

	// Map fields to appropriate keys based on the output format we requested
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
		details["NLWP"] = fields[4]
	}
	if len(fields) >= 6 {
		// The remaining fields (5 onwards) might be part of the start time
		startTime := strings.Join(fields[5:], " ")
		details["STARTED"] = startTime
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
	// For safety, validate the PID before attempting to kill
	if pid <= 0 {
		return fmt.Errorf("invalid PID: %d", pid)
	}

	// Use go-ps to check if the process exists first
	processes, err := ps.Processes()
	if err != nil {
		return fmt.Errorf("failed to get process list: %v", err)
	}

	// Check if the process exists
	exists := false
	for _, proc := range processes {
		if int64(proc.Pid()) == pid {
			exists = true
			break
		}
	}

	if !exists {
		return fmt.Errorf("process with PID %d does not exist", pid)
	}

	// Kill the process using the kill command
	cmd := exec.Command("kill", strconv.FormatInt(pid, 10))
	return cmd.Run()
}
