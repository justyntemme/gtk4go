package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ProcessInfo contains information about a process
type ProcessInfo struct {
	PID        int64
	Name       string
	Username   string
	CPUPercent float64
	MemoryBytes int64
	Threads    int
	State      string
	StartTime  string
}

// getProcesses returns a list of running processes
func getProcesses() ([]ProcessInfo, error) {
	// Get the process data using ps command
	cmd := exec.Command("ps", "-eo", "pid,user,%cpu,rss,state,thcount,lstart,command")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute ps command: %v", err)
	}

	// Parse the output
	processes := []ProcessInfo{}
	lines := strings.Split(string(output), "\n")
	
	// Skip the header line and parse each process line
	for i, line := range lines {
		if i == 0 || len(line) == 0 {
			continue // Skip header and empty lines
		}

		// Process the line - this is a simplified parsing logic
		// In a real app, we would handle potential formatting issues
		fields := splitPSOutput(line)
		if len(fields) < 8 {
			continue // Not enough fields
		}

		pid, err := strconv.ParseInt(strings.TrimSpace(fields[0]), 10, 64)
		if err != nil {
			continue // Invalid PID
		}

		cpuPercent, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
		if err != nil {
			cpuPercent = 0.0 // Default to 0 if parsing fails
		}

		memKB, err := strconv.ParseInt(strings.TrimSpace(fields[3]), 10, 64)
		if err != nil {
			memKB = 0 // Default to 0 if parsing fails
		}
		// Convert KB to bytes
		memBytes := memKB * 1024

		threads, err := strconv.Atoi(strings.TrimSpace(fields[5]))
		if err != nil {
			threads = 0 // Default to 0 if parsing fails
		}

		// Parse the start time - combining fields 6-10 which contain the date
		startTimeStr := strings.Join(fields[6:11], " ")
		startTime := formatStartTime(startTimeStr)

		// Get the command name - should be the last field
		commandStr := fields[11]
		// Extract just the executable name from the command
		cmdParts := strings.Split(commandStr, "/")
		procName := cmdParts[len(cmdParts)-1]
		// Remove any arguments
		procName = strings.Split(procName, " ")[0]

		// Create the process info
		proc := ProcessInfo{
			PID:        pid,
			Name:       procName,
			Username:   strings.TrimSpace(fields[1]),
			CPUPercent: cpuPercent,
			MemoryBytes: memBytes,
			Threads:    threads,
			State:      strings.TrimSpace(fields[4]),
			StartTime:  startTime,
		}
		
		processes = append(processes, proc)
	}

	return processes, nil
}

// splitPSOutput splits the ps command output line into fields
// This is more complex than a simple split due to the format of the ps output
func splitPSOutput(line string) []string {
	fields := []string{}
	
	// PID (field 0)
	pidEnd := strings.IndexFunc(line, func(r rune) bool {
		return r != ' ' && !isDigit(r)
	})
	if pidEnd == -1 {
		pidEnd = len(line)
	}
	fields = append(fields, line[:pidEnd])
	line = strings.TrimSpace(line[pidEnd:])
	
	// Username (field 1)
	usernameEnd := findFirstSpace(line)
	fields = append(fields, line[:usernameEnd])
	line = strings.TrimSpace(line[usernameEnd:])
	
	// CPU% (field 2)
	cpuEnd := findFirstSpace(line)
	fields = append(fields, line[:cpuEnd])
	line = strings.TrimSpace(line[cpuEnd:])
	
	// RSS (field 3)
	rssEnd := findFirstSpace(line)
	fields = append(fields, line[:rssEnd])
	line = strings.TrimSpace(line[rssEnd:])
	
	// State (field 4)
	stateEnd := findFirstSpace(line)
	fields = append(fields, line[:stateEnd])
	line = strings.TrimSpace(line[stateEnd:])
	
	// Thread count (field 5)
	threadEnd := findFirstSpace(line)
	fields = append(fields, line[:threadEnd])
	line = strings.TrimSpace(line[threadEnd:])
	
	// The rest is the start time (6 fields) and command
	// Start time format: "Ddd Mmm DD HH:MM:SS YYYY"
	parts := strings.SplitN(line, " ", 7)
	if len(parts) < 7 {
		// Not enough parts, add what we have
		fields = append(fields, parts...)
		// Pad with empty strings to ensure we have at least 12 fields
		for len(fields) < 12 {
			fields = append(fields, "")
		}
	} else {
		// Add the 6 date/time parts
		fields = append(fields, parts[:6]...)
		// The last part is the command
		fields = append(fields, parts[6])
	}
	
	return fields
}

// findFirstSpace finds the index of the first space in a string
func findFirstSpace(s string) int {
	for i, r := range s {
		if r == ' ' {
			return i
		}
	}
	return len(s)
}

// isDigit returns true if the rune is a digit
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
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
