// Package gtk4 provides a unified callback management system for GTK4 widgets
// File: gtk4go/gtk4/callbacks.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Exported callback functions (implemented in Go)
// extern void callbackHandler(GObject *object, gpointer data);
// extern void callbackHandlerWithParam(GObject *object, gpointer param, gpointer data);
// extern gboolean callbackHandlerWithReturn(GObject *object, gpointer data);
//
// // Generic function to connect a signal to a handler
// static gulong connectSignal(GObject *object, const char *signal, gboolean hasParam, gboolean hasReturn, guint callbackId) {
//     if (hasReturn) {
//         return g_signal_connect(object, signal, G_CALLBACK(callbackHandlerWithReturn), GUINT_TO_POINTER(callbackId));
//     } else if (hasParam) {
//         return g_signal_connect(object, signal, G_CALLBACK(callbackHandlerWithParam), GUINT_TO_POINTER(callbackId));
//     } else {
//         return g_signal_connect(object, signal, G_CALLBACK(callbackHandler), GUINT_TO_POINTER(callbackId));
//     }
// }
//
// // Function to disconnect a signal
// static void disconnectSignal(GObject *object, gulong handlerId) {
//     if (handlerId > 0) {
//         g_signal_handler_disconnect(object, handlerId);
//     }
// }
import "C"

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Import the main package for UI thread execution
import gtk4go "../../gtk4go"

// SignalType represents the type of GTK signal
type SignalType string

// Common GTK signal types
const (
	// Button signals
	SignalClicked SignalType = "clicked"
	
	// Entry signals
	SignalChanged SignalType = "changed"
	SignalActivate SignalType = "activate"
	
	// Window signals
	SignalCloseRequest SignalType = "close-request"
	
	// Dialog signals
	SignalResponse SignalType = "response"
	
	// ListView signals
	SignalListActivate SignalType = "activate"
	
	// SelectionModel signals
	SignalSelectionChanged SignalType = "selection-changed"
	
	// Adjustment signals
	SignalValueChanged SignalType = "value-changed"
	
	// Action signals
	SignalActionActivate SignalType = "activate"
)

// nextCallbackID is a counter for generating unique callback IDs
var nextCallbackID atomic.Uint64

// CallbackManager handles GTK signal callbacks
type CallbackManager struct {
	// Map from callback ID to callback data
	callbacks     sync.Map
	// Map from object pointer to list of handler IDs
	objectHandlers sync.Map
}

// callbackData stores information about a callback
type callbackData struct {
	callback    interface{}
	objectPtr   uintptr
	signal      SignalType
	hasParam    bool
	hasReturn   bool
	handlerID   C.gulong
}

// Global callback manager
var globalCallbackManager = &CallbackManager{}

// EnableCallbackDebugging enables or disables debug output for callbacks
func EnableCallbackDebugging(enable bool) {
    if enable {
        EnableDebugComponent(DebugComponentCallback)
        SetDebugLevel(DebugLevelVerbose) // Set to verbose level for detailed callback info
    } else {
        DisableDebugComponent(DebugComponentCallback)
    }
}

// Connect connects a signal to a callback function
func Connect(object interface{}, signal SignalType, callback interface{}) (handlerID uint64) {
	// Get the object's pointer
	objectPtr := getObjectPointer(object)
	if objectPtr == 0 {
		DebugLog(DebugLevelError, DebugComponentCallback, "Connect failed: couldn't get object pointer for %T", object)
		return 0 // Invalid object
	}

	// Generate a unique ID for this callback
	id := nextCallbackID.Add(1)

	// Check callback signature to determine parameter and return type
	hasParam, hasReturn := analyzeCallbackSignature(callback)

	// Create callback data
	data := &callbackData{
		callback:    callback,
		objectPtr:   objectPtr,
		signal:      signal,
		hasParam:    hasParam,
		hasReturn:   hasReturn,
		handlerID:   0, // Will be set after connection
	}

	// Connect the signal
	cObject := (*C.GObject)(unsafe.Pointer(objectPtr))
	cSignal := C.CString(string(signal))
	defer C.free(unsafe.Pointer(cSignal))

	// Connect and get handler ID
	handlerId := C.connectSignal(
		cObject,
		cSignal,
		boolToGBoolean(hasParam),
		boolToGBoolean(hasReturn),
		C.guint(id),
	)
	
	// Store the handler ID in the callback data
	data.handlerID = handlerId

	// Store the callback data in the map
	globalCallbackManager.callbacks.Store(id, data)

	// Associate this handler with the object for cleanup
	globalCallbackManager.trackObjectHandler(objectPtr, handlerId)

	DebugLog(DebugLevelInfo, DebugComponentCallback, "Connected signal %s with ID %d to object %p", signal, id, objectPtr)
	
	return id
}

// Disconnect disconnects a signal handler by its ID
func Disconnect(id uint64) {
	// Look up the callback data
	value, ok := globalCallbackManager.callbacks.Load(id)
	if !ok {
		DebugLog(DebugLevelWarning, DebugComponentCallback, "Disconnect failed: callback ID %d not found", id)
		return
	}
	
	data := value.(*callbackData)
	
	// Disconnect the signal
	cObject := (*C.GObject)(unsafe.Pointer(data.objectPtr))
	C.disconnectSignal(cObject, data.handlerID)
	
	// Remove the callback from the map
	globalCallbackManager.callbacks.Delete(id)
	
	// Remove the handler from the object's handler list
	globalCallbackManager.untrackObjectHandler(data.objectPtr, data.handlerID)
	
	DebugLog(DebugLevelInfo, DebugComponentCallback, "Disconnected signal handler ID %d from object %p", id, data.objectPtr)
}

// DisconnectAll disconnects all signal handlers for an object
func DisconnectAll(object interface{}) {
	objectPtr := getObjectPointer(object)
	if objectPtr == 0 {
		DebugLog(DebugLevelWarning, DebugComponentCallback, "DisconnectAll failed: couldn't get object pointer for %T", object)
		return
	}
	
	// Get the object's handlers
	value, ok := globalCallbackManager.objectHandlers.Load(objectPtr)
	if !ok {
		DebugLog(DebugLevelVerbose, DebugComponentCallback, "DisconnectAll: no handlers found for object %p", objectPtr)
		return
	}
	
	handlers := value.([]C.gulong)
	
	// Disconnect each handler
	cObject := (*C.GObject)(unsafe.Pointer(objectPtr))
	for _, handlerId := range handlers {
		C.disconnectSignal(cObject, handlerId)
		DebugLog(DebugLevelVerbose, DebugComponentCallback, "DisconnectAll: disconnected handler ID %d from object %p", handlerId, objectPtr)
	}
	
	// Remove the object from the map
	globalCallbackManager.objectHandlers.Delete(objectPtr)
	
	// Remove all callbacks for this object from the callbacks map
	globalCallbackManager.callbacks.Range(func(key, value interface{}) bool {
		data := value.(*callbackData)
		if data.objectPtr == objectPtr {
			globalCallbackManager.callbacks.Delete(key)
			DebugLog(DebugLevelVerbose, DebugComponentCallback, "DisconnectAll: removed callback ID %d from object %p", key, objectPtr)
		}
		return true
	})
}

// getObjectPointer returns the pointer to the GObject of a GTK widget
func getObjectPointer(object interface{}) uintptr {
	// Handle common GTK widget types
	switch obj := object.(type) {
	case Widget:
		return uintptr(unsafe.Pointer(obj.GetWidget()))
	case *Adjustment:
		return uintptr(unsafe.Pointer(obj.adjustment))
	case *Action:
		return uintptr(unsafe.Pointer(obj.action))
	default:
		// Try to find a GetWidget or Native method using reflection
		val := reflect.ValueOf(object)
		if val.Kind() == reflect.Ptr && !val.IsNil() {
			// Try GetWidget method
			getWidgetMethod := val.MethodByName("GetWidget")
			if getWidgetMethod.IsValid() {
				results := getWidgetMethod.Call(nil)
				if len(results) == 1 {
					return uintptr(unsafe.Pointer(results[0].Pointer()))
				}
			}
			
			// Try Native method
			nativeMethod := val.MethodByName("Native")
			if nativeMethod.IsValid() {
				results := nativeMethod.Call(nil)
				if len(results) == 1 && results[0].Kind() == reflect.Uintptr {
					// Convert uint64 to uintptr safely, even on 32-bit systems
					return uintptr(results[0].Interface().(uintptr))
				}
			}
		}
	}
	
	return 0 // Unable to get pointer
}

// analyzeCallbackSignature determines if a callback takes parameters or returns a value
func analyzeCallbackSignature(callback interface{}) (hasParam bool, hasReturn bool) {
	// Get the type of the callback
	callbackType := reflect.TypeOf(callback)
	
	// Must be a function
	if callbackType.Kind() != reflect.Func {
		return false, false
	}
	
	// Check if it has parameters
	hasParam = callbackType.NumIn() > 0
	
	// Check if it has a return value
	hasReturn = callbackType.NumOut() > 0
	
	return hasParam, hasReturn
}

// trackObjectHandler associates a handler ID with an object
func (m *CallbackManager) trackObjectHandler(objectPtr uintptr, handlerId C.gulong) {
	value, ok := m.objectHandlers.Load(objectPtr)
	var handlers []C.gulong
	if ok {
		handlers = value.([]C.gulong)
	} else {
		handlers = make([]C.gulong, 0, 4) // Pre-allocate space for 4 handlers
	}
	
	handlers = append(handlers, handlerId)
	m.objectHandlers.Store(objectPtr, handlers)
	
	DebugLog(DebugLevelVerbose, DebugComponentCallback, "Tracked handler ID %d for object %p", handlerId, objectPtr)
}

// untrackObjectHandler removes a handler ID from an object
func (m *CallbackManager) untrackObjectHandler(objectPtr uintptr, handlerId C.gulong) {
	value, ok := m.objectHandlers.Load(objectPtr)
	if !ok {
		DebugLog(DebugLevelVerbose, DebugComponentCallback, "untrackObjectHandler: no handlers found for object %p", objectPtr)
		return
	}
	
	handlers := value.([]C.gulong)
	
	// Find and remove the handler
	for i, id := range handlers {
		if id == handlerId {
			// Remove by swapping with the last element and slicing
			handlers[i] = handlers[len(handlers)-1]
			handlers = handlers[:len(handlers)-1]
			break
		}
	}
	
	if len(handlers) == 0 {
		// No more handlers for this object
		m.objectHandlers.Delete(objectPtr)
		DebugLog(DebugLevelVerbose, DebugComponentCallback, "untrackObjectHandler: removed last handler for object %p", objectPtr)
	} else {
		m.objectHandlers.Store(objectPtr, handlers)
		DebugLog(DebugLevelVerbose, DebugComponentCallback, "untrackObjectHandler: %d handlers remaining for object %p", len(handlers), objectPtr)
	}
}

// boolToGBoolean converts a Go bool to a C gboolean
func boolToGBoolean(b bool) C.gboolean {
	if b {
		return C.TRUE
	}
	return C.FALSE
}

// execCallback safely executes a callback on the main UI thread
// to ensure thread safety with GTK
func execCallback(callback interface{}, args ...interface{}) {
	// Execute on UI thread to ensure GTK thread safety
	gtk4go.RunOnUIThread(func() {
		// Execute the callback based on its type
		switch cb := callback.(type) {
		case func():
			cb()
		case func(int):
			if len(args) > 0 {
				if i, ok := args[0].(int); ok {
					cb(i)
				}
			}
		case func(ResponseType):
			if len(args) > 0 {
				if rt, ok := args[0].(ResponseType); ok {
					cb(rt)
				}
			}
		case func() bool:
			cb()
		case func(int, int):
			if len(args) > 1 {
				if i1, ok1 := args[0].(int); ok1 {
					if i2, ok2 := args[1].(int); ok2 {
						cb(i1, i2)
					}
				}
			}
		default:
			DebugLog(DebugLevelError, DebugComponentCallback, "Unsupported callback type: %T", callback)
		}
	})
}

// Exported callback functions for CGo

//export callbackHandler
func callbackHandler(object *C.GObject, data C.gpointer) {
	id := uint64(uintptr(data))
	value, ok := globalCallbackManager.callbacks.Load(id)
	if !ok {
		DebugLog(DebugLevelWarning, DebugComponentCallback, "callbackHandler: callback ID %d not found", id)
		return
	}
	
	callbackData := value.(*callbackData)
	DebugLog(DebugLevelVerbose, DebugComponentCallback, "callbackHandler: executing callback ID %d for signal %s", id, callbackData.signal)
	
	// Call the callback with no parameters on the UI thread
	if callback, ok := callbackData.callback.(func()); ok {
		execCallback(callback)
	} else {
		DebugLog(DebugLevelError, DebugComponentCallback, "callbackHandler: callback has wrong type: %T", callbackData.callback)
	}
}

//export callbackHandlerWithParam
func callbackHandlerWithParam(object *C.GObject, param C.gpointer, data C.gpointer) {
	id := uint64(uintptr(data))
	value, ok := globalCallbackManager.callbacks.Load(id)
	if !ok {
		DebugLog(DebugLevelWarning, DebugComponentCallback, "callbackHandlerWithParam: callback ID %d not found", id)
		return
	}
	
	callbackData := value.(*callbackData)
	paramVal := int(uintptr(param))
	DebugLog(DebugLevelVerbose, DebugComponentCallback, "callbackHandlerWithParam: executing callback ID %d for signal %s with param %v", 
	           id, callbackData.signal, paramVal)
	
	// Handle different callback signatures based on the signal type
	switch callbackData.signal {
	case SignalResponse:
		// For dialog responses, param is the response ID
		if callback, ok := callbackData.callback.(func(ResponseType)); ok {
			execCallback(callback, ResponseType(uintptr(param)))
		}
	case SignalListActivate:
		// For list view activation, param is the position
		if callback, ok := callbackData.callback.(func(int)); ok {
			execCallback(callback, paramVal)
		}
	case SignalSelectionChanged:
		// For selection changed, we have position and count
		// Note: This is a simplification as we're only passing position
		if callback, ok := callbackData.callback.(func(int)); ok {
			execCallback(callback, paramVal)
		} else if callback, ok := callbackData.callback.(func(int, int)); ok {
			// In a real implementation, you'd extract both position and count
			execCallback(callback, paramVal, 0)
		}
	default:
		// For other cases, try to call with an int parameter
		if callback, ok := callbackData.callback.(func(int)); ok {
			execCallback(callback, paramVal)
		} else if callback, ok := callbackData.callback.(func(interface{})); ok {
			// For callbacks that accept any parameter
			execCallback(callback, paramVal)
		} else {
			DebugLog(DebugLevelError, DebugComponentCallback, "callbackHandlerWithParam: callback has wrong type: %T", callbackData.callback)
		}
	}
}

//export callbackHandlerWithReturn
func callbackHandlerWithReturn(object *C.GObject, data C.gpointer) C.gboolean {
	id := uint64(uintptr(data))
	value, ok := globalCallbackManager.callbacks.Load(id)
	if !ok {
		DebugLog(DebugLevelWarning, DebugComponentCallback, "callbackHandlerWithReturn: callback ID %d not found", id)
		return C.FALSE
	}
	
	callbackData := value.(*callbackData)
	DebugLog(DebugLevelVerbose, DebugComponentCallback, "callbackHandlerWithReturn: executing callback ID %d for signal %s", id, callbackData.signal)
	
	// Callbacks with return values need to be executed synchronously
	// to get the return value back to C
	if callback, ok := callbackData.callback.(func() bool); ok {
		// Since we need the return value, we can't use execCallback here
		// Ideally, this should still ensure we're on the UI thread
		result := callback()
		if result {
			return C.TRUE
		}
	} else {
		DebugLog(DebugLevelError, DebugComponentCallback, "callbackHandlerWithReturn: callback has wrong type: %T", callbackData.callback)
	}
	
	return C.FALSE
}

// Initialize the callback system
func init() {
	// Register a finalizer to clean up all callbacks at exit
	runtime.SetFinalizer(globalCallbackManager, func(m *CallbackManager) {
		// Disconnect all signals
		m.callbacks.Range(func(key, value interface{}) bool {
			data := value.(*callbackData)
			cObject := (*C.GObject)(unsafe.Pointer(data.objectPtr))
			C.disconnectSignal(cObject, data.handlerID)
			return true
		})
	})
}

// GetCallbackStats returns statistics about the callback system
func GetCallbackStats() map[string]int {
	stats := make(map[string]int)
	
	// Count callbacks
	callbackCount := 0
	globalCallbackManager.callbacks.Range(func(_, _ interface{}) bool {
		callbackCount++
		return true
	})
	stats["TotalCallbacks"] = callbackCount
	
	// Count objects with handlers
	objectCount := 0
	globalCallbackManager.objectHandlers.Range(func(_, _ interface{}) bool {
		objectCount++
		return true
	})
	stats["ObjectsWithCallbacks"] = objectCount
	
	// Count callback types by signal
	signalCounts := make(map[SignalType]int)
	globalCallbackManager.callbacks.Range(func(_, value interface{}) bool {
		data := value.(*callbackData)
		signalCounts[data.signal]++
		return true
	})
	
	for signal, count := range signalCounts {
		stats[fmt.Sprintf("Signal_%s", signal)] = count
	}
	
	return stats
}