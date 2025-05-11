// processGopher/info.go
package main

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessInfo holds essential information about a process
type ProcessInfo struct {
	PID        int32
	Name       string
	Username   string
	CPUPercent float64
	MemoryMB   float32
	CreateTime string
}

// GetAllProcesses retrieves information about all running processes
func GetAllProcesses() ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var procInfos []ProcessInfo
	
	// Use a channel to collect results from concurrent goroutines
	processChan := make(chan ProcessInfo, len(processes))
	var wg sync.WaitGroup
	
	// Limit concurrent process info fetching to avoid overwhelming the system
	semaphore := make(chan struct{}, 10)
	
	for _, proc := range processes {
		wg.Add(1)
		go func(p *process.Process) {
			defer wg.Done()
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			info, err := getProcessInfo(p)
			if err != nil {
				// Skip processes we can't access
				return
			}
			processChan <- *info
		}(proc)
	}
	
	// Start a goroutine to close the channel when all are done
	go func() {
		wg.Wait()
		close(processChan)
	}()
	
	// Collect results
	for info := range processChan {
		procInfos = append(procInfos, info)
	}

	// Sort by memory usage (descending)
	sort.Slice(procInfos, func(i, j int) bool {
		return procInfos[i].MemoryMB > procInfos[j].MemoryMB
	})

	return procInfos, nil
}

// GetAllProcessesWithProgress retrieves information about all running processes
// with progress reporting for background tasks
func GetAllProcessesWithProgress(progressCallback func(current, total int, message string)) ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	total := len(processes)
	var procInfos []ProcessInfo
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Limit concurrent process info fetching
	semaphore := make(chan struct{}, 10)
	completed := int32(0)
	
	for _, proc := range processes {
		wg.Add(1)
		go func(p *process.Process) {
			defer wg.Done()
			
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			info, err := getProcessInfo(p)
			if err != nil {
				// Still count as completed even if error
				current := atomic.AddInt32(&completed, 1)
				if progressCallback != nil && current%10 == 0 {
					progressCallback(int(current), total, fmt.Sprintf("Fetched %d/%d processes", current, total))
				}
				return
			}
			
			mu.Lock()
			procInfos = append(procInfos, *info)
			mu.Unlock()
			
			current := atomic.AddInt32(&completed, 1)
			if progressCallback != nil && current%10 == 0 {
				progressCallback(int(current), total, fmt.Sprintf("Fetched %d/%d processes", current, total))
			}
		}(proc)
	}
	
	wg.Wait()
	
	// Sort by memory usage (descending)
	sort.Slice(procInfos, func(i, j int) bool {
		return procInfos[i].MemoryMB > procInfos[j].MemoryMB
	})

	return procInfos, nil
}

// getProcessInfo extracts detailed information from a process
func getProcessInfo(proc *process.Process) (*ProcessInfo, error) {
	pid := proc.Pid

	name, err := proc.Name()
	if err != nil {
		name = "Unknown"
	}

	username, err := proc.Username()
	if err != nil {
		username = "Unknown"
	}

	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		cpuPercent = 0.0
	}

	memInfo, err := proc.MemoryInfo()
	var memoryMB float32
	if err == nil {
		memoryMB = float32(memInfo.RSS) / (1024 * 1024) // Convert to MB
	}

	createTime, err := proc.CreateTime()
	var createTimeStr string
	if err == nil {
		createTimeStr = time.Unix(createTime/1000, 0).Format("2006-01-02 15:04:05")
	} else {
		createTimeStr = "Unknown"
	}

	return &ProcessInfo{
		PID:        pid,
		Name:       name,
		Username:   username,
		CPUPercent: cpuPercent,
		MemoryMB:   memoryMB,
		CreateTime: createTimeStr,
	}, nil
}

// FormatProcessInfo creates a string representation of process info for display
func FormatProcessInfo(info ProcessInfo) string {
	return fmt.Sprintf("PID: %-6d | %-20s | %-10s | CPU: %5.1f%% | Mem: %7.1f MB",
		info.PID, truncateString(info.Name, 20), truncateString(info.Username, 10),
		info.CPUPercent, info.MemoryMB)
}

// truncateString limits string length and adds ellipsis if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// GetProcessCount returns the total number of running processes
func GetProcessCount() (int, error) {
	processes, err := process.Processes()
	if err != nil {
		return 0, err
	}
	return len(processes), nil
}
