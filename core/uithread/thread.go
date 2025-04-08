// Package uithread provides utilities for managing UI thread operations for GTK.
// File: gtk4go/core/uithread/thread.go
package uithread

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// uiThreadID tracks the ID of the UI thread
var uiThreadID int64

// dispatchQueue is a channel for functions to be executed on the UI thread
var dispatchQueue = make(chan func(), 100)

// initialized tracks whether the system has been initialized
var (
	initialized bool
	initMutex   sync.Mutex
)

// RegisterIdleHandler allows the GTK package to register its implementation
// of the idle function that executes callbacks on the UI thread
var RegisterIdleHandler func(fn func())

// Global variables to manage idle functions
var (
	idleFunctions = sync.Map{}
	nextIdleKey   = atomic.Uint64{}
)

// Initialize initializes the UI thread handling system
func Initialize() {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return
	}

	// Platform-specific initialization is performed in the init function
	// of the platform-specific files (thread_darwin.go, thread_linux.go)
	
	// Store the UI thread ID - this must be done on the main thread
	// Note that on macOS, the OS thread is already locked by the platform-specific init
	uiThreadID = threadID()

	// Initialize platform-specific idle handler
	initPlatformIdleHandler()

	// Start the dispatch queue processor
	go processDispatchQueue()

	initialized = true
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
		// Use registered idle handler if available, otherwise direct call
		if RegisterIdleHandler != nil {
			RegisterIdleHandler(fn)
		} else {
			// Direct call is less ideal but works as fallback
			fn()
		}
	}
}

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

func init() {
	// Initialize the UI thread handling
	Initialize()
}