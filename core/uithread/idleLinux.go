//go:build linux
// +build linux

package uithread

// initPlatformIdleHandler initializes the platform-specific idle handler for Linux
// On Linux, we'll let the GTK main.go code set up the idle handler
func initPlatformIdleHandler() {
    // For Linux, we'll use GTK's g_idle_add system
    // This will be set during gtk4go.Initialize()
    // So we intentionally leave RegisterIdleHandler as nil here
}