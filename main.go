// Package gtk4go provides Go bindings to GTK4.
// File: gtk4go/main.go
package gtk4go

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"runtime"
	"sync"
)

// Initialize ensures GTK is initialized.
// This is automatically called when importing the package.
func Initialize() error {
	// Check if GTK is already initialized
	if C.gtk_is_initialized() == C.FALSE {
		// Initialize GTK
		if C.gtk_init_check() == C.FALSE {
			return fmt.Errorf("failed to initialize GTK")
		}
	}
	return nil
}

// Note: gtk_main and gtk_main_quit are removed in GTK4, use GtkApplication instead
// These functions are kept for compatibility but will log a warning if used

// Main is deprecated in GTK4. Use Application.Run() instead.
func Main() {
	// In GTK4, this is not available. Applications should use GtkApplication instead.
	fmt.Println("Warning: gtk_main() is not available in GTK4. Use GtkApplication instead.")
}

// MainQuit is deprecated in GTK4. Use g_application_quit() instead.
func MainQuit() {
	// In GTK4, this is not available. Applications should use GtkApplication instead.
	fmt.Println("Warning: gtk_main_quit() is not available in GTK4. Use GtkApplication instead.")
}

// Events returns the global GTK events channel.
func Events() chan any {
	return events
}

var (
	events = make(chan any, 10)
	mu     sync.Mutex
)

// init initializes the GTK4 library.
func init() {
	runtime.LockOSThread()
	// Initialize GTK
	Initialize()
}
