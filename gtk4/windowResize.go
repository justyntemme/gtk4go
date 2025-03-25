// Package gtk4 provides window resize detection for GTK4
// File: gtk4go/gtk4/windowResize.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Property notify callbacks for window size changes
// extern void windowPropertyNotifyCallback(GObject *object, GParamSpec *pspec, gpointer user_data);
//
// // Set up window resize detection using property notifications
// static void setupWindowResizeTracking(GtkWindow *window) {
//     // Connect to default-width property changes
//     g_signal_connect(window, "notify::default-width", 
//                     G_CALLBACK(windowPropertyNotifyCallback), window);
//     
//     // Connect to default-height property changes
//     g_signal_connect(window, "notify::default-height", 
//                     G_CALLBACK(windowPropertyNotifyCallback), window);
//
//     // Connect to width-request changes
//     g_signal_connect(window, "notify::width-request", 
//                     G_CALLBACK(windowPropertyNotifyCallback), window);
//     
//     // Connect to height-request changes
//     g_signal_connect(window, "notify::height-request", 
//                     G_CALLBACK(windowPropertyNotifyCallback), window);
//
//     // Surface state changes (maximized, fullscreen, etc.)
//     GdkSurface *surface = gtk_native_get_surface(GTK_NATIVE(window));
//     if (surface) {
//         g_signal_connect(surface, "notify::state", 
//                         G_CALLBACK(windowPropertyNotifyCallback), window);
//     }
// }
//
// // Get current window size via width/height properties
// static void getWindowSize(GtkWindow *window, int *width, int *height) {
//     // Start with default values
//     *width = 0;
//     *height = 0;
//     
//     // Try the surface - most reliable for window dimensions
//     GdkSurface *surface = gtk_native_get_surface(GTK_NATIVE(window));
//     if (surface) {
//         *width = gdk_surface_get_width(surface);
//         *height = gdk_surface_get_height(surface);
//         
//         if (*width > 0 && *height > 0) {
//             return;
//         }
//     }
//     
//     // Try with natural size
//     int natural_width, natural_height;
//     gtk_widget_measure(GTK_WIDGET(window), GTK_ORIENTATION_HORIZONTAL, -1,
//                       NULL, &natural_width, NULL, NULL);
//     gtk_widget_measure(GTK_WIDGET(window), GTK_ORIENTATION_VERTICAL, -1,
//                       NULL, &natural_height, NULL, NULL);
//                       
//     if (natural_width > 0 && natural_height > 0) {
//         *width = natural_width;
//         *height = natural_height;
//         return;
//     }
//     
//     // Get the requested size
//     int request_width, request_height;
//     gtk_widget_get_size_request(GTK_WIDGET(window), &request_width, &request_height);
//     
//     if (request_width > 0 && request_height > 0) {
//         *width = request_width;
//         *height = request_height;
//         return;
//     }
//     
//     // Last resort: use default size
//     gtk_window_get_default_size(window, width, height);
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

//export windowPropertyNotifyCallback
func windowPropertyNotifyCallback(object *C.GObject, pspec *C.GParamSpec, userData C.gpointer) {
	windowPtr := uintptr(unsafe.Pointer(userData))

	windowResizeStateMutex.RLock()
	data, exists := windowResizeState[windowPtr]
	windowResizeStateMutex.RUnlock()

	if !exists {
		return
	}

	now := time.Now()
	lastTime, _ := data.lastResizeTime.Load().(time.Time)
	data.lastResizeTime.Store(now)

	// Get current dimensions using C helper
	var width, height C.int
	C.getWindowSize((*C.GtkWindow)(unsafe.Pointer(userData)), &width, &height)
	
	// Only proceed if we got valid dimensions
	if width <= 0 || height <= 0 {
		return
	}
	
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

	// Get initial window size
	var width, height C.int
	C.getWindowSize((*C.GtkWindow)(unsafe.Pointer(w.widget)), &width, &height)
	data.width = int32(width)
	data.height = int32(height)

	// Store in map
	windowResizeStateMutex.Lock()
	windowResizeState[windowPtr] = data
	windowResizeStateMutex.Unlock()

	// Set up signal connections for resize detection
	C.setupWindowResizeTracking((*C.GtkWindow)(unsafe.Pointer(w.widget)))
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