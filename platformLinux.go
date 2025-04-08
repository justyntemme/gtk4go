//go:build linux
// +build linux

package gtk4go

import (
	"os"
	"os/exec"
	"strings"
)

func init() {
	// Detect Wayland vs X11 and set appropriate renderer
	if isWaylandSession() {
		// Wayland works best with GL renderer
		os.Setenv("GSK_RENDERER", "gl")
		os.Setenv("GDK_GL", "always")
		
		// Force Wayland backend
		os.Setenv("GDK_BACKEND", "wayland")
	} else {
		// X11 has better compatibility with Cairo renderer for some systems
		// but can use GL if available
		if hasGLSupport() {
			os.Setenv("GSK_RENDERER", "gl")
			os.Setenv("GDK_GL", "always")
		} else {
			os.Setenv("GSK_RENDERER", "cairo")
		}
		
		// Force X11 backend for traditional X11 sessions
		os.Setenv("GDK_BACKEND", "x11")
	}
}

// isWaylandSession returns true if running under Wayland
func isWaylandSession() bool {
	// Check environment variables first
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return true
	}
	
	if os.Getenv("XDG_SESSION_TYPE") == "wayland" {
		return true
	}
	
	// If environment variables aren't set, try to detect session type
	cmd := exec.Command("loginctl", "show-session", "$XDG_SESSION_ID", "-p", "Type")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "wayland") {
		return true
	}
	
	return false
}

// hasGLSupport checks if the system has OpenGL support
func hasGLSupport() bool {
	// Check for presence of GL libraries
	glPaths := []string{
		"/usr/lib/libGL.so",
		"/usr/lib/x86_64-linux-gnu/libGL.so",
		"/usr/lib64/libGL.so",
	}
	
	for _, path := range glPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	
	// Try to detect Mesa or other GL implementation
	cmd := exec.Command("glxinfo")
	if err := cmd.Run(); err == nil {
		return true
	}
	
	// Check if X server has GLX extension
	cmd = exec.Command("xdpyinfo")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "GLX") {
		return true
	}
	
	return false
}