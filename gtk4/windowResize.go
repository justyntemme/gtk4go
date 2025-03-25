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
	// Atomic fields for thread-safe access without mutex
	isResizing      atomic.Bool     // Whether window is currently being resized
	width           atomic.Int32    // Current window width
	height          atomic.Int32    // Current window height
	watcherID       atomic.Int64    // ID of the current resize watcher goroutine
	
	// Time values stored in atomic.Value
	lastResizeTime  atomic.Value    // time.Time - Last time a resize event was detected
	resizeStartTime atomic.Value    // time.Time - When resize operation began
	
	// Configuration fields (protected by mutex)
	mu                 sync.RWMutex  // Mutex for non-atomic fields
	onResizeStart      func()        // Function to call when resize starts
	onResizeEnd        func()        // Function to call when resize ends
	resizeEndThreshold time.Duration // Threshold to detect end of resize operation
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
	data.lastResizeTime.Store(now)

	// Get current dimensions using C helper
	var width, height C.int
	C.getWindowSize((*C.GtkWindow)(unsafe.Pointer(userData)), &width, &height)
	
	// Only proceed if we got valid dimensions
	if width <= 0 || height <= 0 {
		return
	}
	
	// Load old dimensions atomically and store new dimensions atomically
	oldWidth := data.width.Load()
	oldHeight := data.height.Load()
	newWidth := int32(width)
	newHeight := int32(height)
	
	// Store the new dimensions using atomic operations
	data.width.Store(newWidth)
	data.height.Store(newHeight)

	// Check if dimensions actually changed
	if newWidth == oldWidth && newHeight == oldHeight {
		return
	}

	// Is this a new resize operation?
	wasResizing := data.isResizing.Load()
	if !wasResizing {
		// Mark as resizing using atomic operation
		data.isResizing.Store(true)
		data.resizeStartTime.Store(now)

		// Call resize start handler if set (thread-safe access)
		data.mu.RLock()
		onResizeStart := data.onResizeStart
		data.mu.RUnlock()
		
		if onResizeStart != nil {
			onResizeStart()
		}
	}

	// Ensure a watcher goroutine is running
	startResizeWatcher(data, windowPtr)
}

// startResizeWatcher ensures a single watcher goroutine is running to detect resize completion
func startResizeWatcher(data *windowResizeData, windowPtr uintptr) {
	// Generate a unique watcher ID
	newWatcherID := time.Now().UnixNano()
	
	// Try to set the new watcher ID using atomic operation, and get the old one
	oldWatcherID := data.watcherID.Swap(newWatcherID)
	
	// If oldWatcherID was 0, no watcher was running
	if oldWatcherID == 0 {
		go resizeWatcherGoroutine(data, windowPtr, newWatcherID)
	}
	// If oldWatcherID was non-zero, a watcher is already running
	// It will detect it's been replaced when it checks its ID
}

// resizeWatcherGoroutine monitors for resize completion
func resizeWatcherGoroutine(data *windowResizeData, windowPtr uintptr, myID int64) {
	// Get the threshold (thread-safe)
	data.mu.RLock()
	threshold := data.resizeEndThreshold
	data.mu.RUnlock()
	
	// Wait for a short period to see if resize continues
	time.Sleep(threshold)
	
	// Check if we're still the active watcher using atomic load
	if data.watcherID.Load() != myID {
		return // Another watcher has taken over
	}
	
	// Mark that no watcher is running using atomic operation
	data.watcherID.Store(int64(0))
	
	// Check if no new resize events have occurred
	lastResizeTime, ok := data.lastResizeTime.Load().(time.Time)
	if !ok {
		return // Type assertion failed
	}
	
	if time.Since(lastResizeTime) >= threshold {
		// Resize has ended, but only if we're still in resize state
		if data.isResizing.Swap(false) { // returns old value and sets to false atomically
			// We were resizing and now we're not
			
			// Call resize end handler if set (thread-safe)
			data.mu.RLock()
			onResizeEnd := data.onResizeEnd
			data.mu.RUnlock()
			
			if onResizeEnd != nil {
				onResizeEnd()
			}
		}
	} else {
		// Still getting resize events, start another watcher
		startResizeWatcher(data, windowPtr)
	}
}

// SetupResizeDetection sets up resize detection for a window
func (w *Window) SetupResizeDetection(onResizeStart, onResizeEnd ResizeCallback) {
	windowPtr := uintptr(unsafe.Pointer(w.widget))

	// Create resize data with thread-safe initialization
	data := &windowResizeData{
		resizeEndThreshold: 200 * time.Millisecond, // Default threshold
	}

	// Store callbacks with mutex protection
	data.mu.Lock()
	data.onResizeStart = onResizeStart
	data.onResizeEnd = onResizeEnd
	data.mu.Unlock()

	// Initialize atomic values
	data.lastResizeTime.Store(time.Time{})
	data.resizeStartTime.Store(time.Time{})
	data.watcherID.Store(int64(0))

	// Get initial window size
	var width, height C.int
	C.getWindowSize((*C.GtkWindow)(unsafe.Pointer(w.widget)), &width, &height)
	
	// Store the initial dimensions using atomic operations
	data.width.Store(int32(width))
	data.height.Store(int32(height))

	// Store in map with mutex protection
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

	// Use atomic operation to get the resizing state
	return data.isResizing.Load()
}

// GetSize returns the current window size
func (w *Window) GetSize() (width, height int32) {
	windowPtr := uintptr(unsafe.Pointer(w.widget))

	windowResizeStateMutex.RLock()
	data, exists := windowResizeState[windowPtr]
	windowResizeStateMutex.RUnlock()

	if !exists {
		return 0, 0
	}

	// Use atomic operations to get the dimensions
	return data.width.Load(), data.height.Load()
}

// SetResizeEndThreshold sets the time threshold for detecting the end of a resize operation
func (w *Window) SetResizeEndThreshold(threshold time.Duration) {
	windowPtr := uintptr(unsafe.Pointer(w.widget))

	windowResizeStateMutex.RLock()
	data, exists := windowResizeState[windowPtr]
	windowResizeStateMutex.RUnlock()

	if !exists {
		return
	}

	// Use mutex for non-atomic field
	data.mu.Lock()
	data.resizeEndThreshold = threshold
	data.mu.Unlock()
}

// CleanupResizeDetection cleans up resize detection for a window
func (w *Window) CleanupResizeDetection() {
	windowPtr := uintptr(unsafe.Pointer(w.widget))

	windowResizeStateMutex.Lock()
	delete(windowResizeState, windowPtr)
	windowResizeStateMutex.Unlock()
}