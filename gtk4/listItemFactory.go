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
	"unsafe"
)

// ListItemCallback represents a callback for list item operations
type ListItemCallback func(listItem *ListItem)

// Define signal types for list item factory - using GTK's actual signal names
const (
	SignalSetup    SignalType = "setup"
	SignalBind     SignalType = "bind"
	SignalUnbind   SignalType = "unbind"
	SignalTeardown SignalType = "teardown"
)

//export setupListItemCallback
func setupListItemCallback(factory *C.GtkSignalListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	// Get factory pointer for lookup in the unified callback system
	factoryPtr := uintptr(unsafe.Pointer(factory))
	
	// Create a Go wrapper for the list item
	goListItem := &ListItem{listItem: listItem}
	
	// Find callback using the unified callback system
	if callback := GetCallback(factoryPtr, SignalSetup); callback != nil {
		// The modified SafeCallback function in callbacks.go now handles ListItemCallback
		SafeCallback(callback, goListItem)
	}
}

//export bindListItemCallback
func bindListItemCallback(factory *C.GtkSignalListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	// Get factory pointer for lookup in the unified callback system
	factoryPtr := uintptr(unsafe.Pointer(factory))
	
	// Create a Go wrapper for the list item
	goListItem := &ListItem{listItem: listItem}
	
	// Find callback using the unified callback system
	if callback := GetCallback(factoryPtr, SignalBind); callback != nil {
		// The modified SafeCallback function in callbacks.go now handles ListItemCallback
		SafeCallback(callback, goListItem)
	} else {
		// If no callback is registered, try a default implementation
		// that just sets the text on the child label
		goListItem.UpdateChildWithText()
	}
}

//export unbindListItemCallback
func unbindListItemCallback(factory *C.GtkSignalListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	// Get factory pointer for lookup in the unified callback system
	factoryPtr := uintptr(unsafe.Pointer(factory))
	
	// Create a Go wrapper for the list item
	goListItem := &ListItem{listItem: listItem}
	
	// Find callback using the unified callback system
	if callback := GetCallback(factoryPtr, SignalUnbind); callback != nil {
		// The modified SafeCallback function in callbacks.go now handles ListItemCallback
		SafeCallback(callback, goListItem)
	}
}

//export teardownListItemCallback
func teardownListItemCallback(factory *C.GtkSignalListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	// Get factory pointer for lookup in the unified callback system
	factoryPtr := uintptr(unsafe.Pointer(factory))
	
	// Create a Go wrapper for the list item
	goListItem := &ListItem{listItem: listItem}
	
	// Find callback using the unified callback system
	if callback := GetCallback(factoryPtr, SignalTeardown); callback != nil {
		// The modified SafeCallback function in callbacks.go now handles ListItemCallback
		SafeCallback(callback, goListItem)
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

	// Use the standard Connect function with the correct GTK signal name
	Connect(f, SignalSetup, callback)
}

// ConnectBind connects a callback for the bind signal
func (f *SignalListItemFactory) ConnectBind(callback ListItemCallback) {
	if callback == nil {
		return
	}

	// Use the standard Connect function with the correct GTK signal name
	Connect(f, SignalBind, callback)
}

// ConnectUnbind connects a callback for the unbind signal
func (f *SignalListItemFactory) ConnectUnbind(callback ListItemCallback) {
	if callback == nil {
		return
	}

	// Use the standard Connect function with the correct GTK signal name
	Connect(f, SignalUnbind, callback)
}

// ConnectTeardown connects a callback for the teardown signal
func (f *SignalListItemFactory) ConnectTeardown(callback ListItemCallback) {
	if callback == nil {
		return
	}

	// Use the standard Connect function with the correct GTK signal name
	Connect(f, SignalTeardown, callback)
}

// DisconnectSetup disconnects the setup signal callback
func (f *SignalListItemFactory) DisconnectSetup() {
	factoryPtr := uintptr(unsafe.Pointer(f.factory))
	callbackIDs := getCallbackIDsForSignal(factoryPtr, SignalSetup)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// DisconnectBind disconnects the bind signal callback
func (f *SignalListItemFactory) DisconnectBind() {
	factoryPtr := uintptr(unsafe.Pointer(f.factory))
	callbackIDs := getCallbackIDsForSignal(factoryPtr, SignalBind)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// DisconnectUnbind disconnects the unbind signal callback
func (f *SignalListItemFactory) DisconnectUnbind() {
	factoryPtr := uintptr(unsafe.Pointer(f.factory))
	callbackIDs := getCallbackIDsForSignal(factoryPtr, SignalUnbind)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// DisconnectTeardown disconnects the teardown signal callback
func (f *SignalListItemFactory) DisconnectTeardown() {
	factoryPtr := uintptr(unsafe.Pointer(f.factory))
	callbackIDs := getCallbackIDsForSignal(factoryPtr, SignalTeardown)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// Destroy frees resources associated with the factory
func (f *SignalListItemFactory) Destroy() {
	if f.factory != nil {
		// Disconnect all signal handlers using the unified callback system
		DisconnectAll(f)
		
		C.g_object_unref(C.gpointer(unsafe.Pointer(f.factory)))
		f.factory = nil
	}
}