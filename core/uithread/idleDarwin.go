//go:build darwin
// +build darwin

package uithread

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
/*
#import <Cocoa/Cocoa.h>

// Forward declaration of Go callback
extern void goCallbackHandler(void* data);

// Function to dispatch to the main queue
static void dispatchToMainQueue(void* data) {
    dispatch_async(dispatch_get_main_queue(), ^{
        goCallbackHandler(data);
    });
}
*/
import "C"
import "unsafe"

// initPlatformIdleHandler initializes the platform-specific idle handler
func initPlatformIdleHandler() {
    // Use the Cocoa dispatch queue for macOS
    RegisterIdleHandler = func(fn func()) {
        // Store function in Go-side map to prevent GC
        idleKey := nextIdleKey.Add(1)
        idleFunctions.Store(idleKey, fn)
        
        // Use Cocoa's dispatch_async with main queue
        C.dispatchToMainQueue(unsafe.Pointer(uintptr(idleKey)))
    }
}

//export goCallbackHandler
func goCallbackHandler(data unsafe.Pointer) {
    idleKey := uint64(uintptr(data))
    
    // Get function from map
    if fnVal, ok := idleFunctions.Load(idleKey); ok {
        // Remove it after retrieval
        idleFunctions.Delete(idleKey)
        
        // Call the function
        if fn, ok := fnVal.(func()); ok {
            fn()
        }
    }
}