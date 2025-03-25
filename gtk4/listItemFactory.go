// Package gtk4 provides list item factory functionality for GTK4
// File: gtk4go/gtk4/listitemfactory.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal list item factory callbacks
// extern void setupListItemCallback(GtkSignalListItemFactory *factory, GtkListItem *list_item, gpointer user_data);
// extern void bindListItemCallback(GtkSignalListItemFactory *factory, GtkListItem *list_item, gpointer user_data);
// extern void unbindListItemCallback(GtkSignalListItemFactory *factory, GtkListItem *list_item, gpointer user_data);
// extern void teardownListItemCallback(GtkSignalListItemFactory *factory, GtkListItem *list_item, gpointer user_data);
//
// // Connect signals for list item factory
// static GtkSignalListItemFactory* createSignalListItemFactory() {
//     return (GtkSignalListItemFactory*)gtk_signal_list_item_factory_new();
// }
//
// static gulong connectSetupListItem(GtkSignalListItemFactory *factory, gpointer user_data) {
//     return g_signal_connect(factory, "setup", G_CALLBACK(setupListItemCallback), user_data);
// }
//
// static gulong connectBindListItem(GtkSignalListItemFactory *factory, gpointer user_data) {
//     return g_signal_connect(factory, "bind", G_CALLBACK(bindListItemCallback), user_data);
// }
//
// static gulong connectUnbindListItem(GtkSignalListItemFactory *factory, gpointer user_data) {
//     return g_signal_connect(factory, "unbind", G_CALLBACK(unbindListItemCallback), user_data);
// }
//
// static gulong connectTeardownListItem(GtkSignalListItemFactory *factory, gpointer user_data) {
//     return g_signal_connect(factory, "teardown", G_CALLBACK(teardownListItemCallback), user_data);
// }
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// ListItemCallback represents a callback for list item operations
type ListItemCallback func(listItem *ListItem)

// ListItemCallbackType defines the types of callbacks for list item factory
type ListItemCallbackType int

const (
	// ListItemCallbackSetup is called when a new list item is created
	ListItemCallbackSetup ListItemCallbackType = iota
	// ListItemCallbackBind is called when a list item is bound to a model item
	ListItemCallbackBind
	// ListItemCallbackUnbind is called when a list item is unbound from a model item
	ListItemCallbackUnbind
	// ListItemCallbackTeardown is called when a list item is destroyed
	ListItemCallbackTeardown
)

var (
	// Map of factory pointers to maps of callback types to callbacks
	factoryCallbacks     = make(map[uintptr]map[ListItemCallbackType]ListItemCallback)
	factoryCallbackMutex sync.RWMutex
)

//export setupListItemCallback
func setupListItemCallback(factory *C.GtkSignalListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	factoryCallbackMutex.RLock()
	defer factoryCallbackMutex.RUnlock()

	// Convert factory pointer to uintptr for lookup
	factoryPtr := uintptr(unsafe.Pointer(factory))

	// Find and call the callback
	if callbacks, ok := factoryCallbacks[factoryPtr]; ok {
		if callback, ok := callbacks[ListItemCallbackSetup]; ok {
			callback(&ListItem{listItem: listItem})
		}
	}
}

//export bindListItemCallback
func bindListItemCallback(factory *C.GtkSignalListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	factoryCallbackMutex.RLock()
	defer factoryCallbackMutex.RUnlock()

	// Convert factory pointer to uintptr for lookup
	factoryPtr := uintptr(unsafe.Pointer(factory))

	// Find and call the callback
	if callbacks, ok := factoryCallbacks[factoryPtr]; ok {
		if callback, ok := callbacks[ListItemCallbackBind]; ok {
			callback(&ListItem{listItem: listItem})
		}
	}
}

//export unbindListItemCallback
func unbindListItemCallback(factory *C.GtkSignalListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	factoryCallbackMutex.RLock()
	defer factoryCallbackMutex.RUnlock()

	// Convert factory pointer to uintptr for lookup
	factoryPtr := uintptr(unsafe.Pointer(factory))

	// Find and call the callback
	if callbacks, ok := factoryCallbacks[factoryPtr]; ok {
		if callback, ok := callbacks[ListItemCallbackUnbind]; ok {
			callback(&ListItem{listItem: listItem})
		}
	}
}

//export teardownListItemCallback
func teardownListItemCallback(factory *C.GtkSignalListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	factoryCallbackMutex.RLock()
	defer factoryCallbackMutex.RUnlock()

	// Convert factory pointer to uintptr for lookup
	factoryPtr := uintptr(unsafe.Pointer(factory))

	// Find and call the callback
	if callbacks, ok := factoryCallbacks[factoryPtr]; ok {
		if callback, ok := callbacks[ListItemCallbackTeardown]; ok {
			callback(&ListItem{listItem: listItem})
		}
	}
}

// ListItemFactory is an interface for factories that create list items
type ListItemFactory interface {
	// GetListItemFactory returns the underlying GtkListItemFactory pointer
	GetListItemFactory() *C.GtkListItemFactory

	// Destroy frees resources associated with the factory
	Destroy()
}

// SignalListItemFactory is a factory that creates list items using signals
type SignalListItemFactory struct {
	factory *C.GtkSignalListItemFactory
}

// NewSignalListItemFactory creates a new signal list item factory
func NewSignalListItemFactory() *SignalListItemFactory {
	factory := &SignalListItemFactory{
		factory: C.createSignalListItemFactory(),
	}

	// Initialize the callback map for this factory
	factoryCallbackMutex.Lock()
	factoryCallbacks[uintptr(unsafe.Pointer(factory.factory))] = make(map[ListItemCallbackType]ListItemCallback)
	factoryCallbackMutex.Unlock()

	runtime.SetFinalizer(factory, (*SignalListItemFactory).Destroy)
	return factory
}

// GetListItemFactory returns the underlying GtkListItemFactory pointer
func (f *SignalListItemFactory) GetListItemFactory() *C.GtkListItemFactory {
	return (*C.GtkListItemFactory)(unsafe.Pointer(f.factory))
}

// ConnectSetup connects a callback for the setup signal
func (f *SignalListItemFactory) ConnectSetup(callback ListItemCallback) {
	if callback == nil {
		return
	}

	factoryCallbackMutex.Lock()
	defer factoryCallbackMutex.Unlock()

	// Store the callback in the map
	factoryPtr := uintptr(unsafe.Pointer(f.factory))
	factoryCallbacks[factoryPtr][ListItemCallbackSetup] = callback

	// Connect the signal
	C.connectSetupListItem(f.factory, C.gpointer(unsafe.Pointer(f.factory)))
}

// ConnectBind connects a callback for the bind signal
func (f *SignalListItemFactory) ConnectBind(callback ListItemCallback) {
	if callback == nil {
		return
	}

	factoryCallbackMutex.Lock()
	defer factoryCallbackMutex.Unlock()

	// Store the callback in the map
	factoryPtr := uintptr(unsafe.Pointer(f.factory))
	factoryCallbacks[factoryPtr][ListItemCallbackBind] = callback

	// Connect the signal
	C.connectBindListItem(f.factory, C.gpointer(unsafe.Pointer(f.factory)))
}

// ConnectUnbind connects a callback for the unbind signal
func (f *SignalListItemFactory) ConnectUnbind(callback ListItemCallback) {
	if callback == nil {
		return
	}

	factoryCallbackMutex.Lock()
	defer factoryCallbackMutex.Unlock()

	// Store the callback in the map
	factoryPtr := uintptr(unsafe.Pointer(f.factory))
	factoryCallbacks[factoryPtr][ListItemCallbackUnbind] = callback

	// Connect the signal
	C.connectUnbindListItem(f.factory, C.gpointer(unsafe.Pointer(f.factory)))
}

// ConnectTeardown connects a callback for the teardown signal
func (f *SignalListItemFactory) ConnectTeardown(callback ListItemCallback) {
	if callback == nil {
		return
	}

	factoryCallbackMutex.Lock()
	defer factoryCallbackMutex.Unlock()

	// Store the callback in the map
	factoryPtr := uintptr(unsafe.Pointer(f.factory))
	factoryCallbacks[factoryPtr][ListItemCallbackTeardown] = callback

	// Connect the signal
	C.connectTeardownListItem(f.factory, C.gpointer(unsafe.Pointer(f.factory)))
}

// Destroy frees resources associated with the factory
func (f *SignalListItemFactory) Destroy() {
	if f.factory != nil {
		factoryCallbackMutex.Lock()
		delete(factoryCallbacks, uintptr(unsafe.Pointer(f.factory)))
		factoryCallbackMutex.Unlock()

		C.g_object_unref(C.gpointer(unsafe.Pointer(f.factory)))
		f.factory = nil
	}
}