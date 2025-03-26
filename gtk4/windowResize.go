// Package gtk4 provides window resize detection for GTK4 using the unified callback system
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
	"sync/atomic"
	"time"
	"unsafe"
	
	// Import the core uithread package for thread-safe operations
	"../core/uithread"
)

// windowResizeState stores state information for resize detection
type windowResizeState struct {
	// Atomic fields for thread-safe access
	isResizing      atomic.Bool  // Whether window is currently being resized
	width           atomic.Int32 // Current window width
	height          atomic.Int32 // Current window height
	lastResizeTime  atomic.Int64 // Last time a resize event was detected (Unix nano)
	resizeStartTime atomic.Int64 // When resize operation began (Unix nano)
}

// Global map of window pointers to resize state
var windowResizeStates = make(map[uintptr]*windowResizeState)

//export windowPropertyNotifyCallback
func windowPropertyNotifyCallback(object *C.GObject, pspec *C.GParamSpec, userData C.gpointer) {
	windowPtr := uintptr(unsafe.Pointer(userData))
	
	// Get or create state for this window
	state, ok := windowResizeStates[windowPtr]
	if !ok {
		// Skip if window is not being tracked
		return
	}
	
	now := time.Now()
	state.lastResizeTime.Store(now.UnixNano())

	// Get current dimensions using C helper
	var width, height C.int
	C.getWindowSize((*C.GtkWindow)(unsafe.Pointer(userData)), &width, &height)
	
	// Only proceed if we got valid dimensions
	if width <= 0 || height <= 0 {
		return
	}
	
	// Store the old and new dimensions
	oldWidth := state.width.Load()
	oldHeight := state.height.Load()
	newWidth := int32(width)
	newHeight := int32(height)
	
	// Store the new dimensions
	state.width.Store(newWidth)
	state.height.Store(newHeight)

	// Check if dimensions actually changed
	if newWidth == oldWidth && newHeight == oldHeight {
		return
	}

	// Is this a new resize operation?
	wasResizing := state.isResizing.Load()
	if !wasResizing {
		// Mark as resizing
		state.isResizing.Store(true)
		state.resizeStartTime.Store(now.UnixNano())

		// Trigger resize start callback via the unified callback system
		if callback := GetCallback(windowPtr, SignalResizeStart); callback != nil {
			// Execute the callback
			SafeCallback(callback)
		}
	} else {
		// Trigger resize update callback via the unified callback system
		if callback := GetCallback(windowPtr, SignalResizeUpdate); callback != nil {
			// Execute the callback
			SafeCallback(callback)
		}
	}
	
	// Start or restart resize end detection
	go detectResizeEnd(windowPtr)
}

// detectResizeEnd waits for a period without resize events and then triggers the resize end callback
func detectResizeEnd(windowPtr uintptr) {
	// Default threshold for resize end detection (200ms)
	threshold := 200 * time.Millisecond
	
	// Sleep for threshold duration
	time.Sleep(threshold)
	
	// Get state
	state, ok := windowResizeStates[windowPtr]
	if !ok {
		return // Window no longer being tracked
	}
	
	// Check if resize has ended (no new events during threshold period)
	lastResizeTime := time.Unix(0, state.lastResizeTime.Load())
	if time.Since(lastResizeTime) >= threshold && state.isResizing.Load() {
		// Mark resize as ended
		state.isResizing.Store(false)
		
		// Trigger resize end callback via the unified callback system
		if callback := GetCallback(windowPtr, SignalResizeEnd); callback != nil {
			// Run on UI thread using our safe callback mechanism
			uithread.RunOnUIThread(func() {
				// Execute the callback safely
				SafeCallback(callback)
			})
		}
	}
}

// SetupResizeDetection sets up resize detection for a window
func (w *Window) SetupResizeDetection() {
	windowPtr := uintptr(unsafe.Pointer(w.widget))
	
	// Create resize state if not already exists
	if _, ok := windowResizeStates[windowPtr]; !ok {
		// Create resize state
		state := &windowResizeState{}
		
		// Store initial window size
		var width, height C.int
		C.getWindowSize((*C.GtkWindow)(unsafe.Pointer(w.widget)), &width, &height)
		state.width.Store(int32(width))
		state.height.Store(int32(height))
		
		// Store state in global map
		windowResizeStates[windowPtr] = state
		
		// Set up property notification in C
		C.setupWindowResizeTracking((*C.GtkWindow)(unsafe.Pointer(w.widget)))
	}
}

// ConnectResizeStart connects a callback for when resize starts
func (w *Window) ConnectResizeStart(callback func()) uint64 {
	// Ensure resize detection is set up
	w.SetupResizeDetection()
	
	// Use the unified callback system
	return Connect(w, SignalResizeStart, callback)
}

// ConnectResizeEnd connects a callback for when resize ends
func (w *Window) ConnectResizeEnd(callback func()) uint64 {
	// Ensure resize detection is set up
	w.SetupResizeDetection()
	
	// Use the unified callback system
	return Connect(w, SignalResizeEnd, callback)
}

// ConnectResizeUpdate connects a callback for resize updates
func (w *Window) ConnectResizeUpdate(callback func()) uint64 {
	// Ensure resize detection is set up
	w.SetupResizeDetection()
	
	// Use the unified callback system
	return Connect(w, SignalResizeUpdate, callback)
}

// DisconnectResizeStart disconnects the resize start callback
func (w *Window) DisconnectResizeStart() {
	// Get all callbacks for this window
	windowPtr := uintptr(unsafe.Pointer(w.widget))
	callbackIDs := getCallbackIDsForSignal(windowPtr, SignalResizeStart)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// DisconnectResizeEnd disconnects the resize end callback
func (w *Window) DisconnectResizeEnd() {
	// Get all callbacks for this window
	windowPtr := uintptr(unsafe.Pointer(w.widget))
	callbackIDs := getCallbackIDsForSignal(windowPtr, SignalResizeEnd)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// DisconnectResizeUpdate disconnects the resize update callback
func (w *Window) DisconnectResizeUpdate() {
	// Get all callbacks for this window
	windowPtr := uintptr(unsafe.Pointer(w.widget))
	callbackIDs := getCallbackIDsForSignal(windowPtr, SignalResizeUpdate)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// IsResizing returns true if the window is currently being resized
func (w *Window) IsResizing() bool {
	windowPtr := uintptr(unsafe.Pointer(w.widget))
	state, ok := windowResizeStates[windowPtr]
	if !ok {
		return false
	}
	return state.isResizing.Load()
}

// GetSize returns the current window size
func (w *Window) GetSize() (width, height int32) {
	windowPtr := uintptr(unsafe.Pointer(w.widget))
	state, ok := windowResizeStates[windowPtr]
	if !ok {
		return 0, 0
	}
	return state.width.Load(), state.height.Load()
}

// CleanupResizeDetection cleans up resize detection for a window
func (w *Window) CleanupResizeDetection() {
	// Remove resize state
	delete(windowResizeStates, uintptr(unsafe.Pointer(w.widget)))
}

// SetupCSSOptimizedResize sets up CSS optimization during resize
func (w *Window) SetupCSSOptimizedResize() {
	// Set up resize callbacks for CSS optimization
	w.ConnectResizeStart(func() {
		// Optimize all global CSS providers
		optimizeAllProviders()

		// Switch to lightweight CSS
		display := C.gdk_display_get_default()
		useResizeCSSProvider(display)
	})
	
	w.ConnectResizeEnd(func() {
		// Reset all global CSS providers
		resetAllProviders()

		// Restore normal CSS
		display := C.gdk_display_get_default()
		restoreOriginalCSSProvider(display, nil)
	})
}