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
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// uiThreadID tracks the ID of the UI thread
var uiThreadID int64

// dispatchQueue is a channel for functions to be executed on the UI thread
var dispatchQueue = make(chan func(), 100)

// initialized tracks whether GTK has been initialized
var (
	initialized bool
	initMutex   sync.Mutex
	idleHandles sync.Map // Maps uint64 keys to idle handles
	nextIdleKey uint64
)

// Initialize ensures GTK is initialized and starts the dispatch queue.
// This is automatically called when importing the package.
func Initialize() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return nil
	}

	// Store the UI thread ID
	uiThreadID = threadID()

	// Check if GTK is already initialized
	if C.gtk_is_initialized() == C.FALSE {
		// Initialize GTK
		if C.gtk_init_check() == C.FALSE {
			return fmt.Errorf("failed to initialize GTK")
		}
	}

	// Start the dispatch queue processor
	go processDispatchQueue()

	initialized = true
	return nil
}

// IsUIThread returns true if the current goroutine is running on the UI thread
func IsUIThread() bool {
	return threadID() == atomic.LoadInt64(&uiThreadID)
}

// RunOnUIThread schedules a function to be executed on the UI thread.
// If called from the UI thread, the function is executed immediately.
func RunOnUIThread(fn func()) {
	if IsUIThread() {
		fn()
		return
	}
	dispatchQueue <- fn
}

// MustRunOnUIThread panics if not called from the UI thread
func MustRunOnUIThread() {
	if !IsUIThread() {
		panic("This function must be called from the UI thread")
	}
}

// threadID returns a unique identifier for the current OS thread
func threadID() int64 {
	var id int64
	// This func will be executed on the current OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	// Use the memory address of a local variable as a proxy for thread ID
	id = int64(uintptr(unsafe.Pointer(&id)))
	return id
}

// processDispatchQueue processes functions in the dispatch queue
func processDispatchQueue() {
	for fn := range dispatchQueue {
		invokeOnUIThread(fn)
	}
}

// invokeOnUIThread schedules a Go function to be executed on the UI thread
func invokeOnUIThread(fn func()) {
	// Get a unique key for this function
	key := atomic.AddUint64(&nextIdleKey, 1)

	// Store the function in the idle handles map
	idleHandles.Store(key, fn)

	// Schedule the function to be executed on the UI thread
	C.addIdleFunction(C.gpointer(uintptr(key)))
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
	mu     sync.Mutex
)

// SafeUIOperation executes a function safely on the UI thread
// and returns when the operation is complete
func SafeUIOperation(operation func()) {
	if IsUIThread() {
		operation()
		return
	}

	// Use a channel to synchronize
	done := make(chan struct{})

	RunOnUIThread(func() {
		operation()
		close(done)
	})

	// Wait for the operation to complete
	<-done
}

// init initializes the GTK4 library.
func init() {
	runtime.LockOSThread()
	// Initialize GTK
	Initialize()
}
