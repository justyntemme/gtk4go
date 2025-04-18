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
// extern gboolean tooltipQueryCallback(GtkWidget *widget, gint x, gint y, gboolean keyboard_mode, GtkTooltip *tooltip, gpointer user_data);
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
// // Connect tooltip query signal specifically
// static gulong connectTooltipQuery(GtkWidget *widget, guint callbackId) {
//     return g_signal_connect(widget, "query-tooltip", G_CALLBACK(tooltipQueryCallback), GUINT_TO_POINTER(callbackId));
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

	// Import core uithread package
	"github.com/justyntemme/gtk4go/core/uithread"
)

// SignalType represents the type of GTK signal
type SignalType string

// SignalSource represents the source of a signal (to distinguish between signals with the same name)
type SignalSource int

const (
	// Signal sources
	SourceGeneric SignalSource = iota
	SourceListView
	SourceAction
)

// Common GTK signal types
const (
	// Button signals
	SignalClicked SignalType = "clicked"

	// Entry signals
	SignalChanged  SignalType = "changed"
	SignalActivate SignalType = "activate"

	// Window signals
	SignalCloseRequest SignalType = "close-request"

	// Window resize signals
	SignalResizeStart  SignalType = "resize-start"
	SignalResizeEnd    SignalType = "resize-end"
	SignalResizeUpdate SignalType = "resize-update"

	// Dialog signals
	SignalResponse SignalType = "response"

	// ListView signals - same name as action signal but different context
	SignalListActivate SignalType = "activate"

	// SelectionModel signals
	SignalSelectionChanged SignalType = "selection-changed"

	// Adjustment signals
	SignalValueChanged SignalType = "value-changed"

	// Action signals - same name as list signal but different context
	SignalActionActivate SignalType = "activate"

	// Tooltip signals
	SignalQueryTooltip SignalType = "query-tooltip"
)

// Import debug components defined in debug.go
// Only defining SignalQueryTooltip here since it's related to callbacks

// nextCallbackID is a counter for generating unique callback IDs
var nextCallbackID atomic.Uint64

// CallbackManager handles GTK signal callbacks
type CallbackManager struct {
	// Map from callback ID to callback data
	callbacks sync.Map
	// Map from object pointer to list of handler IDs
	objectHandlers sync.Map
	// Map from object pointer to map of signal type to callback data
	objectCallbacks sync.Map
}

// callbackData stores information about a callback
type callbackData struct {
	callback  interface{}
	objectPtr uintptr
	signal    SignalType
	source    SignalSource // Added field to track signal source
	hasParam  bool
	hasReturn bool
	handlerID C.gulong
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

	// Determine signal source based on object type and signal
	source := SourceGeneric
	if _, isListView := object.(*ListView); isListView && signal == SignalListActivate {
		source = SourceListView
	} else if _, isAction := object.(*Action); isAction && signal == SignalActionActivate {
		source = SourceAction
	}

	// Create callback data
	data := &callbackData{
		callback:  callback,
		objectPtr: objectPtr,
		signal:    signal,
		source:    source,
		hasParam:  hasParam,
		hasReturn: hasReturn,
		handlerID: 0, // Will be set after connection
	}

	// Connect the signal
	cObject := (*C.GObject)(unsafe.Pointer(objectPtr))
	cSignal := C.CString(string(signal))
	defer C.free(unsafe.Pointer(cSignal))

	// Connect and get handler ID - special case for tooltip query signal
	var handlerId C.gulong
	if signal == SignalQueryTooltip {
		handlerId = C.connectTooltipQuery((*C.GtkWidget)(unsafe.Pointer(objectPtr)), C.guint(id))
	} else {
		// Connect regular signal
		handlerId = C.connectSignal(
			cObject,
			cSignal,
			boolToGBoolean(hasParam),
			boolToGBoolean(hasReturn),
			C.guint(id),
		)
	}

	// Store the handler ID in the callback data
	data.handlerID = handlerId

	// Store the callback data in the map
	globalCallbackManager.callbacks.Store(id, data)

	// Associate this handler with the object for cleanup
	globalCallbackManager.trackObjectHandler(objectPtr, handlerId)

	// Store callback by object and signal for direct lookups
	globalCallbackManager.storeObjectCallback(objectPtr, signal, callback)

	DebugLog(DebugLevelInfo, DebugComponentCallback, "Connected signal %s with ID %d to object %p (source: %d)",
		signal, id, objectPtr, source)

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

	// Remove the callback from the maps
	globalCallbackManager.callbacks.Delete(id)
	globalCallbackManager.removeObjectCallback(data.objectPtr, data.signal)

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

	// Remove the object from the maps
	globalCallbackManager.objectHandlers.Delete(objectPtr)
	globalCallbackManager.objectCallbacks.Delete(objectPtr)

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

// GetCallback retrieves a callback for a specific object and signal
func GetCallback(objectPtr uintptr, signal SignalType) interface{} {
	objectCallbacksValue, ok := globalCallbackManager.objectCallbacks.Load(objectPtr)
	if !ok {
		return nil
	}

	objectCallbacks := objectCallbacksValue.(map[SignalType]interface{})
	callback, ok := objectCallbacks[signal]
	if !ok {
		return nil
	}

	return callback
}

// getCallbackIDsForSignal returns all callback IDs for a specific object and signal
func getCallbackIDsForSignal(objectPtr uintptr, signal SignalType) []uint64 {
	var ids []uint64

	// Scan all callbacks for matches
	globalCallbackManager.callbacks.Range(func(id, value interface{}) bool {
		data := value.(*callbackData)
		if data.objectPtr == objectPtr && data.signal == signal {
			ids = append(ids, id.(uint64))
		}
		return true
	})

	return ids
}

// storeObjectCallback stores a callback by object pointer and signal type
func (m *CallbackManager) storeObjectCallback(objectPtr uintptr, signal SignalType, callback interface{}) {
	objectCallbacksValue, ok := m.objectCallbacks.Load(objectPtr)
	var objectCallbacks map[SignalType]interface{}

	if !ok {
		objectCallbacks = make(map[SignalType]interface{})
	} else {
		objectCallbacks = objectCallbacksValue.(map[SignalType]interface{})
	}

	objectCallbacks[signal] = callback
	m.objectCallbacks.Store(objectPtr, objectCallbacks)
}

// removeObjectCallback removes a callback by object pointer and signal type
func (m *CallbackManager) removeObjectCallback(objectPtr uintptr, signal SignalType) {
	objectCallbacksValue, ok := m.objectCallbacks.Load(objectPtr)
	if !ok {
		return
	}

	objectCallbacks := objectCallbacksValue.(map[SignalType]interface{})
	delete(objectCallbacks, signal)

	if len(objectCallbacks) == 0 {
		m.objectCallbacks.Delete(objectPtr)
	} else {
		m.objectCallbacks.Store(objectPtr, objectCallbacks)
	}
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
	case *SignalListItemFactory:
		return uintptr(unsafe.Pointer(obj.factory))
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

			// Try GetListItemFactory method for ListItemFactory
			getFactoryMethod := val.MethodByName("GetListItemFactory")
			if getFactoryMethod.IsValid() {
				results := getFactoryMethod.Call(nil)
				if len(results) == 1 {
					return uintptr(unsafe.Pointer(results[0].Pointer()))
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

	// Special case for tooltips which return a boolean
	if callbackType.NumOut() > 0 && callbackType.Out(0).Kind() == reflect.Bool {
		hasReturn = true
	}

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
	uithread.RunOnUIThread(func() {
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
		// Support for ListItemCallback and its equivalent function type
		case ListItemCallback:
			if len(args) > 0 {
				if li, ok := args[0].(*ListItem); ok {
					cb(li)
				} else {
					DebugLog(DebugLevelError, DebugComponentCallback,
						"ListItemCallback called with invalid argument type: %T, expected *ListItem", args[0])
				}
			}
		case func(*ListItem):
			if len(args) > 0 {
				if li, ok := args[0].(*ListItem); ok {
					cb(li)
				} else {
					DebugLog(DebugLevelError, DebugComponentCallback,
						"func(*ListItem) called with invalid argument type: %T, expected *ListItem", args[0])
				}
			}
		// Support for tooltip query callbacks
		case func(int, int, bool, uintptr) bool:
			if len(args) >= 4 {
				if x, ok1 := args[0].(int); ok1 {
					if y, ok2 := args[1].(int); ok2 {
						if keyboard, ok3 := args[2].(bool); ok3 {
							if tooltipPtr, ok4 := args[3].(uintptr); ok4 {
								cb(x, y, keyboard, tooltipPtr)
							}
						}
					}
				}
			}
		case func(int, int, bool, *Tooltip) bool:
			if len(args) >= 4 {
				if x, ok1 := args[0].(int); ok1 {
					if y, ok2 := args[1].(int); ok2 {
						if keyboard, ok3 := args[2].(bool); ok3 {
							if tooltip, ok4 := args[3].(*Tooltip); ok4 {
								cb(x, y, keyboard, tooltip)
							}
						}
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
	DebugLog(DebugLevelVerbose, DebugComponentCallback, "callbackHandlerWithParam: executing callback ID %d for signal %s with param %v (source: %d)",
		id, callbackData.signal, paramVal, callbackData.source)

	// Check for ListItemCallback specifically
	if _, isListItemCallback := callbackData.callback.(ListItemCallback); isListItemCallback {
		DebugLog(DebugLevelInfo, DebugComponentCallback,
			"Found ListItemCallback, wrapping list item pointer %v", uintptr(param))

		// Create a ListItem wrapper for the GtkListItem pointer
		listItem := &ListItem{listItem: (*C.GtkListItem)(unsafe.Pointer(param))}

		// Execute the callback with the ListItem
		execCallback(callbackData.callback, listItem)
		return
	}

	// Similarly handle func(*ListItem) type
	if _, isFuncListItem := callbackData.callback.(func(*ListItem)); isFuncListItem {
		DebugLog(DebugLevelInfo, DebugComponentCallback,
			"Found func(*ListItem), wrapping list item pointer %v", uintptr(param))

		// Create a ListItem wrapper for the GtkListItem pointer
		listItem := &ListItem{listItem: (*C.GtkListItem)(unsafe.Pointer(param))}

		// Execute the callback with the ListItem
		execCallback(callbackData.callback, listItem)
		return
	}

	// Handle different callback signatures based on the signal type and source
	switch {
	case callbackData.signal == SignalResponse:
		// For dialog responses, param is the response ID
		if callback, ok := callbackData.callback.(func(ResponseType)); ok {
			execCallback(callback, ResponseType(uintptr(param)))
		}

	case callbackData.signal == SignalSelectionChanged:
		// For selection changed, we have position and count
		if callback, ok := callbackData.callback.(func(int)); ok {
			execCallback(callback, paramVal)
		} else if callback, ok := callbackData.callback.(func(int, int)); ok {
			// In a real implementation, you'd extract both position and count
			execCallback(callback, paramVal, 0)
		}

	case callbackData.signal == SignalListActivate && callbackData.source == SourceListView:
		// For ListView activation - check for multiple possible types
		// First try direct function type
		if callback, ok := callbackData.callback.(func(int)); ok {
			DebugLog(DebugLevelInfo, DebugComponentListView,
				"Executing list activate callback for position %d", paramVal)
			execCallback(callback, paramVal)
		} else if callback, ok := callbackData.callback.(ListViewActivateCallback); ok {
			// Then try the specific callback type
			DebugLog(DebugLevelInfo, DebugComponentListView,
				"Executing ListViewActivateCallback for position %d", paramVal)
			execCallback(func(pos int) {
				callback(pos)
			}, paramVal)
		} else {
			// Log error if neither type matches
			DebugLog(DebugLevelError, DebugComponentListView,
				"ListActivate callback has wrong type: %T, expected func(int) or ListViewActivateCallback",
				callbackData.callback)
		}

	case callbackData.signal == SignalActionActivate && callbackData.source == SourceAction:
		// For Action activation
		if callback, ok := callbackData.callback.(func()); ok {
			DebugLog(DebugLevelInfo, DebugComponentAction,
				"Executing action activate callback")
			execCallback(callback)
		} else {
			DebugLog(DebugLevelError, DebugComponentAction,
				"ActionActivate callback has wrong type: %T, expected func()", callbackData.callback)
		}

	case callbackData.signal == SignalActivate:
		// For general activation (Entry, etc)
		if callback, ok := callbackData.callback.(func()); ok {
			execCallback(callback)
		} else {
			DebugLog(DebugLevelError, DebugComponentCallback,
				"Activate callback has wrong type: %T, expected func()", callbackData.callback)
		}

	default:
		// For other cases, try to call with an int parameter
		if callback, ok := callbackData.callback.(func(int)); ok {
			execCallback(callback, paramVal)
		} else if callback, ok := callbackData.callback.(func(interface{})); ok {
			// For callbacks that accept any parameter
			execCallback(callback, paramVal)
		} else if callback, ok := callbackData.callback.(func()); ok {
			// Try no parameter callback as last resort
			execCallback(callback)
		} else {
			DebugLog(DebugLevelError, DebugComponentCallback,
				"callbackHandlerWithParam: callback has wrong type: %T", callbackData.callback)
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

//export tooltipQueryCallback
func tooltipQueryCallback(widget *C.GtkWidget, x C.gint, y C.gint, keyboardMode C.gboolean, tooltip *C.GtkTooltip, userData C.gpointer) C.gboolean {
	// Convert userData to callback ID
	id := uint64(uintptr(userData))

	// Find the callback data
	value, ok := globalCallbackManager.callbacks.Load(id)
	if !ok {
		DebugLog(DebugLevelWarning, DebugComponentTooltip,
			"tooltipQueryCallback: callback ID %d not found", id)
		// Default to showing the tooltip
		return C.TRUE
	}

	callbackData := value.(*callbackData)
	DebugLog(DebugLevelVerbose, DebugComponentTooltip,
		"tooltipQueryCallback: executing callback ID %d for signal %s", id, callbackData.signal)

	// Check for the appropriate callback type (for the UCS)
	if callback, ok := callbackData.callback.(func(int, int, bool, uintptr) bool); ok {
		// Execute the callback directly as we need the return value
		result := callback(
			int(x),
			int(y),
			keyboardMode == C.TRUE,
			uintptr(unsafe.Pointer(tooltip)),
		)

		if result {
			return C.TRUE
		}
		return C.FALSE
	} else if callback, ok := callbackData.callback.(func(int, int, bool, *Tooltip) bool); ok {
		// Create a Tooltip wrapper and call the callback
		tooltipObj := &Tooltip{tooltip: tooltip}
		result := callback(
			int(x),
			int(y),
			keyboardMode == C.TRUE,
			tooltipObj,
		)

		if result {
			return C.TRUE
		}
		return C.FALSE
	} else {
		DebugLog(DebugLevelError, DebugComponentTooltip,
			"tooltipQueryCallback: callback has wrong type: %T", callbackData.callback)
	}

	// Default to showing the tooltip
	return C.TRUE
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

// SafeCallback safely executes a callback, ensuring it runs on the UI thread
func SafeCallback(callback interface{}, args ...interface{}) {
	execCallback(callback, args...)
}

// StoreCallback is a helper function to store a callback in the UCS
func StoreCallback(ptr uintptr, signal SignalType, callback interface{}, handlerID C.gulong) {
	// Store the callback in the UCS
	callbackMap := make(map[SignalType]interface{})
	callbackMap[signal] = callback
	globalCallbackManager.objectCallbacks.Store(ptr, callbackMap)

	// Track handler ID for cleanup
	handlers := []C.gulong{handlerID}
	globalCallbackManager.objectHandlers.Store(ptr, handlers)
}

// StoreDirectCallback is a helper function to directly store a callback for a pointer
// This bypasses the normal Connect mechanism to ensure direct pointer matching
func StoreDirectCallback(ptr uintptr, signal SignalType, callback interface{}) {
	// Access the global callback manager's maps directly
	globalCallbackManager.objectCallbacks.Store(ptr, map[SignalType]interface{}{
		signal: callback,
	})

	DebugLog(DebugLevelInfo, DebugComponentCallback,
		"Directly stored callback for pointer %v and signal %s", ptr, signal)
}

// RunOnUIThread runs a function on the UI thread
func RunOnUIThread(fn func()) {
	uithread.RunOnUIThread(fn)
}

