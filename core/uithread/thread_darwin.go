//go:build darwin
// +build darwin

package uithread

import (
	"runtime"
)

func init() {
	// On macOS, we need to lock the main thread for Cocoa/AppKit integration
	// This ensures proper handling of UI events and prevents crashes
	runtime.LockOSThread()
}