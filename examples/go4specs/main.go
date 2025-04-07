package main

import (
	"../../../gtk4go"
	"../../gtk4"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Define constants for styling and sizing
const (
	APP_ID         = "com.example.system-info"
	TITLE          = "System Info"
	DEFAULT_WIDTH  = 900
	DEFAULT_HEIGHT = 600

	// Auto-refresh interval in seconds (0 = disabled)
	AUTO_REFRESH_INTERVAL = 30
)

// labelMap stores references to labels for updating with thread safety
type labelMap struct {
	mu     sync.RWMutex
	labels map[string]*gtk4.Label
}

func newLabelMap() *labelMap {
	return &labelMap{
		labels: make(map[string]*gtk4.Label),
	}
}

func (lm *labelMap) add(key string, label *gtk4.Label) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.labels[key] = label
}

func (lm *labelMap) update(key string, value string) {
	lm.mu.RLock()
	label, ok := lm.labels[key]
	lm.mu.RUnlock()

	if ok && label != nil {
		// Update UI on the UI thread
		gtk4go.RunOnUIThread(func() {
			label.SetText(value)
		})
	}
}

// Global variables for data and UI state
var (
	osLabels           *labelMap
	cpuLabels          *labelMap
	memoryLabels       *labelMap
	diskLabels         *labelMap
	gpuLabels          *labelMap
	statusLabel        *gtk4.Label
	autoRefreshEnabled bool = true
	lastRefreshTime    time.Time
	autoRefreshTimer   *time.Timer
	appInstance        *gtk4.Application // Store application instance globally

	// Thread-safety improvements
	refreshAtomicFlag atomic.Int32 // 0 = not refreshing, 1 = refreshing
)

func main() {
	os.Setenv("GSK_RENDERER", "cairo")
	os.Setenv("GDK_GL", "0")
	gtk4.EnableCallbackDebugging(true)

	// Initialize GTK
	if err := gtk4go.Initialize(); err != nil {
		fmt.Printf("Failed to initialize GTK: %v\n", err)
		os.Exit(1)
	}

	// Create application
	app := gtk4.NewApplication(APP_ID)
	appInstance = app // Store application instance globally

	// Create the refresh action
	refreshAction := gtk4.NewAction("refresh", func() {
		refreshAllData()
	})
	app.GetActionGroup().AddAction(refreshAction)

	// Create window
	win := gtk4.NewWindow(TITLE)
	win.SetDefaultSize(DEFAULT_WIDTH, DEFAULT_HEIGHT)

	// Create main layout
	mainBox := createMainLayout(win)

	// Set the window's child to the main box
	win.SetChild(mainBox)

	// Set up window close handler
	win.ConnectCloseRequest(func() bool {
		// Clean up resources
		if autoRefreshTimer != nil {
			autoRefreshTimer.Stop()
		}
		return false // Return false to allow window to close
	})

	// Add window to application
	app.AddWindow(win)

	// Start auto-refresh timer if enabled
	if AUTO_REFRESH_INTERVAL > 0 {
		startAutoRefreshTimer()
	}

	// Run the application
	os.Exit(app.Run())
}

// Functions for refreshing data and other operations will be here
// I'm omitting them as they are part of the original code and not changed

