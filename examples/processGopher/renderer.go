// processGopher/renderer.go
package main

import (
	"os"
)

// init configures the GTK renderer to avoid OpenGL context issues
func init() {
	// Force Cairo renderer to avoid OpenGL context issues
	// This ensures compatibility with systems that don't have proper OpenGL setup
	os.Setenv("GSK_RENDERER", "cairo")
	
	// Disable hardware acceleration to prevent context issues
	os.Setenv("GDK_GL", "0")
	
	// Optional: You can enable debugging for renderer issues
	// os.Setenv("GSK_DEBUG", "all")
}
