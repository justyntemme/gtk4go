//go:build darwin
// +build darwin

package gtk4go

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func init() {
	// macOS-specific environment settings
	
	// Lock the main thread for proper Cocoa/AppKit integration
	// This is critical for macOS UI stability
	runtime.LockOSThread()
	
	// Check if Metal is supported on this macOS system
	if hasMetalSupport() {
		// Use GL renderer which can leverage Metal on macOS
		os.Setenv("GSK_RENDERER", "gl")
		os.Setenv("GDK_GL", "always")
	} else {
		// Fallback to Cairo renderer for older macOS systems
		os.Setenv("GSK_RENDERER", "cairo")
		os.Setenv("GDK_GL", "0")
	}
	
	// Reduce animation complexity for better performance
	os.Setenv("GTK_ANIMATION_TIMEOUT_FACTOR", "2")
	
	// Set DPI handling for Retina displays only if needed
	// This is detected automatically in GTK4, so we don't force it
	
	// Use native decorations for modern macOS appearance
	os.Setenv("GTK_THEME", "Adwaita")
}

// hasMetalSupport checks if the system supports Metal graphics API
func hasMetalSupport() bool {
	// Check macOS version - Metal requires 10.14+ (Mojave or later)
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	version := strings.TrimSpace(string(output))
	
	// Simple version check (a more robust implementation would parse the version)
	// Metal is well supported on macOS 10.14 and later
	if strings.HasPrefix(version, "10.") {
		// Extract minor version number
		parts := strings.Split(version, ".")
		if len(parts) > 1 {
			if minor := parts[1]; minor != "" {
				if minor >= "14" {
					return true
				}
			}
		}
		return false
	}
	
	// macOS 11+ (Big Sur and later) fully supports Metal
	return true
}