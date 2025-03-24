// Package gtk4 provides window resize detection for GTK4
// File: gtk4go/gtk4/windowResize.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Callback for window size-allocate events
// extern void windowSizeAllocateCallback(GtkWidget *widget, int width, int height, gpointer user_data);
//
// // Connect size-allocate signal for resize detection
// static void connectSizeAllocate(GtkWidget *window) {
//     g_signal_connect(window, "size-allocate", G_CALLBACK(windowSizeAllocateCallback), window);
// }
import "C"

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	// Track resize state for each window
	windowResizeState      = make(map[uintptr]*windowResizeData)
	windowResizeStateMutex sync.RWMutex
)

// windowResizeData tracks resize state for a window
type windowResizeData struct {
	isResizing      atomic.Bool
	lastResizeTime  atomic.Value // time.Time
	resizeStartTime atomic.Value // time.Time
	width           int32
	height          int32

	// Function to call when resize starts
	onResizeStart func()

	// Function to call when resize ends
	onResizeEnd func()

	// Resize threshold to detect end of resize operation
	resizeEndThreshold time.Duration
}

// ResizeCallback is a function called when resize state changes
type ResizeCallback func()

//export windowSizeAllocateCallback
func windowSizeAllocateCallback(widget *C.GtkWidget, width, height C.int, userData C.gpointer) {
	windowPtr := uintptr(unsafe.Pointer(widget))

	windowResizeStateMutex.RLock()
	data, exists := windowResizeState[windowPtr]
	windowResizeStateMutex.RUnlock()

	if !exists {
		return
	}

	now := time.Now()
	lastTime, _ := data.lastResizeTime.Load().(time.Time)
	data.lastResizeTime.Store(now)

	// Get current dimensions
	newWidth := int32(width)
	newHeight := int32(height)
	oldWidth := atomic.LoadInt32(&data.width)
	oldHeight := atomic.LoadInt32(&data.height)

	// Store new dimensions
	atomic.StoreInt32(&data.width, newWidth)
	atomic.StoreInt32(&data.height, newHeight)

	// Check if dimensions actually changed
	if newWidth == oldWidth && newHeight == oldHeight {
		return
	}

	// Is this a new resize operation?
	wasResizing := data.isResizing.Load()
	if !wasResizing {
		// Mark as resizing
		data.isResizing.Store(true)
		data.resizeStartTime.Store(now)

		// Call resize start handler if set
		if data.onResizeStart != nil {
			data.onResizeStart()
		}
	}

	// Start a goroutine to detect when resize is done
	// Only start a new one if this is a different resize operation
	if !lastTime.IsZero() && now.Sub(lastTime) > 100*time.Millisecond {
		go func() {
			// Wait for a short period to see if resize continues
			time.Sleep(data.resizeEndThreshold)

			// Check if no new resize events have occurred
			lastResizeTime, _ := data.lastResizeTime.Load().(time.Time)
			if time.Since(lastResizeTime) >= data.resizeEndThreshold {
				// Resize has ended
				if data.isResizing.Load() {
					data.isResizing.Store(false)

					// Call resize end handler if set
					if data.onResizeEnd != nil {
						data.onResizeEnd()
					}
				}
			}
		}()
	}
}

// SetupResizeDetection sets up resize detection for a window
func (w *Window) SetupResizeDetection(onResizeStart, onResizeEnd ResizeCallback) {
	windowPtr := uintptr(unsafe.Pointer(w.widget))

	// Create resize data
	data := &windowResizeData{
		onResizeStart:      onResizeStart,
		onResizeEnd:        onResizeEnd,
		resizeEndThreshold: 200 * time.Millisecond, // Default threshold
	}

	// Initialize atomic values
	data.lastResizeTime.Store(time.Time{})
	data.resizeStartTime.Store(time.Time{})

	// Store in map
	windowResizeStateMutex.Lock()
	windowResizeState[windowPtr] = data
	windowResizeStateMutex.Unlock()

	// Connect size-allocate signal
	C.connectSizeAllocate(w.widget)
}

// SetupCSSOptimizedResize sets up CSS optimization during resize
func (w *Window) SetupCSSOptimizedResize() {
	// Set up resize detection with CSS optimization
	w.SetupResizeDetection(
		// On resize start
		func() {
			// Optimize all global CSS providers
			optimizeAllProviders()

			// Switch to lightweight CSS
			display := C.gdk_display_get_default()
			useResizeCSSProvider(display)
		},
		// On resize end
		func() {
			// Reset all global CSS providers
			resetAllProviders()

			// Restore normal CSS
			display := C.gdk_display_get_default()
			restoreOriginalCSSProvider(display, nil)
		},
	)
}

// IsResizing returns true if the window is currently being resized
func (w *Window) IsResizing() bool {
	windowPtr := uintptr(unsafe.Pointer(w.widget))

	windowResizeStateMutex.RLock()
	data, exists := windowResizeState[windowPtr]
	windowResizeStateMutex.RUnlock()

	if !exists {
		return false
	}

	return data.isResizing.Load()
}

// CleanupResizeDetection cleans up resize detection for a window
func (w *Window) CleanupResizeDetection() {
	windowPtr := uintptr(unsafe.Pointer(w.widget))

	windowResizeStateMutex.Lock()
	delete(windowResizeState, windowPtr)
	windowResizeStateMutex.Unlock()
}
