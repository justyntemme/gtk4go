// Package gtk4go provides Go bindings to GTK4.
// File: gtk4go/main.go
package gtk4go

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // C callback for idle functions
// extern gboolean idleCallback(gpointer user_data);
//
// // Add an idle function to be called on the main loop
// static guint addIdleFunction(gpointer user_data) {
//     // Use GSourceFunc signature (gboolean (*)(gpointer)) explicitly
//     return g_idle_add((GSourceFunc)idleCallback, user_data);
// }
//
// // Remove a source from the main loop
// static void removeSource(guint source_id) {
//     g_source_remove(source_id);
// }
import "C"

import (
	"fmt"
	"sync"
	"sync/atomic"

	// Import the core uithread package
	"./core/uithread"
)

// initialized tracks whether GTK has been initialized
var (
	initialized bool
	initMutex   sync.Mutex
	idleHandles sync.Map // Maps uint64 keys to idle handles
	nextIdleKey atomic.Uint64
)

// Initialize ensures GTK is initialized and starts the dispatch queue.
// This is automatically called when importing the package.
func Initialize() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return nil
	}

	// Check if GTK is already initialized
	if C.gtk_is_initialized() == C.FALSE {
		// Initialize GTK
		if C.gtk_init_check() == C.FALSE {
			return fmt.Errorf("failed to initialize GTK")
		}
	}

	// Register the GTK idle handler with the uithread package
	uithread.RegisterIdleHandler = func(fn func()) {
		// Get a unique key for this function
		key := nextIdleKey.Add(1)

		// Store the function in the idle handles map
		idleHandles.Store(key, fn)

		// Schedule the function to be executed on the UI thread
		C.addIdleFunction(C.gpointer(uintptr(key)))
	}

	initialized = true
	return nil
}

// RunOnUIThread schedules a function to be executed on the UI thread.
func RunOnUIThread(fn func()) {
	uithread.RunOnUIThread(fn)
}

// IsUIThread returns true if the current goroutine is running on the UI thread
func IsUIThread() bool {
	return uithread.IsUIThread()
}

//export idleCallback
func idleCallback(userData C.gpointer) C.gboolean {
	// Get the key from the user data
	key := uint64(uintptr(userData))

	// Get the function from the idle handles map
	fnVal, ok := idleHandles.Load(key)
	if !ok {
		return C.FALSE
	}

	// Remove the function from the map
	idleHandles.Delete(key)

	// Call the function
	fn := fnVal.(func())
	fn()

	// Return FALSE to remove the idle function
	return C.FALSE
}

// Events returns the global GTK events channel.
func Events() chan any {
	return events
}

var (
	events = make(chan any, 10)
)

// init initializes the GTK4 library.
func init() {
	// Initialize GTK
	Initialize()
}