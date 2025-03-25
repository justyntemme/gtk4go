// Package gtk4 provides button functionality for GTK4
// File: gtk4go/gtk4/button.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback function for button clicks
// extern void buttonClickedCallback(GtkButton *button, gpointer user_data);
//
// // Connect button click signal with callback
// static gulong connectButtonClicked(GtkWidget *button, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(button), "clicked", G_CALLBACK(buttonClickedCallback), user_data);
// }
import "C"

import (
	"sync"
	"unsafe"
)

// ButtonClickedCallback represents a callback for button clicked events
type ButtonClickedCallback func()

// Thread-safe callback registry using RWMutex for efficient locking
type buttonCallbackRegistry struct {
	sync.RWMutex
	callbacks map[uintptr]ButtonClickedCallback
}

// Global button callback registry
var buttonRegistry = &buttonCallbackRegistry{
	callbacks: make(map[uintptr]ButtonClickedCallback),
}

// set adds or updates a callback in the registry
func (r *buttonCallbackRegistry) set(key uintptr, callback ButtonClickedCallback) {
	r.Lock()
	defer r.Unlock()
	r.callbacks[key] = callback
}

// get retrieves a callback by key
func (r *buttonCallbackRegistry) get(key uintptr) (ButtonClickedCallback, bool) {
	r.RLock()
	defer r.RUnlock()
	callback, ok := r.callbacks[key]
	return callback, ok
}

// delete removes a callback from the registry
func (r *buttonCallbackRegistry) delete(key uintptr) {
	r.Lock()
	defer r.Unlock()
	delete(r.callbacks, key)
}

//export buttonClickedCallback
func buttonClickedCallback(button *C.GtkButton, userData C.gpointer) {
	// Convert button pointer to uintptr for lookup
	buttonPtr := uintptr(unsafe.Pointer(button))

	// Get callback - use read-only lock for better performance
	if callback, ok := buttonRegistry.get(buttonPtr); ok {
		callback()
	}
}

// ButtonOption is a function that configures a button
type ButtonOption func(*Button)

// Button represents a GTK button
type Button struct {
	BaseWidget
}

// NewButton creates a new GTK button with the given label
func NewButton(label string, options ...ButtonOption) *Button {
	var widget *C.GtkWidget

	WithCString(label, func(cLabel *C.char) {
		widget = C.gtk_button_new_with_label(cLabel)
	})

	button := &Button{
		BaseWidget: BaseWidget{
			widget: widget,
		},
	}

	// Apply options
	for _, option := range options {
		option(button)
	}

	SetupFinalization(button, button.Destroy)
	return button
}

// WithMnemonic creates a button with mnemonic support
func WithMnemonic(label string) ButtonOption {
	return func(b *Button) {
		WithCString(label, func(cLabel *C.char) {
			b.widget = C.gtk_button_new_with_mnemonic(cLabel)
		})
	}
}

// SetLabel sets the button's label
func (b *Button) SetLabel(label string) {
	WithCString(label, func(cLabel *C.char) {
		C.gtk_button_set_label((*C.GtkButton)(unsafe.Pointer(b.widget)), cLabel)
	})
}

// GetLabel gets the button's label
func (b *Button) GetLabel() string {
	cLabel := C.gtk_button_get_label((*C.GtkButton)(unsafe.Pointer(b.widget)))
	return C.GoString(cLabel)
}

// ConnectClicked connects a callback function to the button's "clicked" signal
func (b *Button) ConnectClicked(callback ButtonClickedCallback) {
	// Store callback in registry
	buttonPtr := uintptr(unsafe.Pointer(b.widget))
	buttonRegistry.set(buttonPtr, callback)

	// Connect signal
	C.connectButtonClicked(b.widget, C.gpointer(unsafe.Pointer(b.widget)))
}

// DisconnectClicked disconnects the clicked signal handler
func (b *Button) DisconnectClicked() {
	// Remove callback from registry
	buttonPtr := uintptr(unsafe.Pointer(b.widget))
	buttonRegistry.delete(buttonPtr)
}

// Destroy destroys the button and cleans up resources
func (b *Button) Destroy() {
	// Remove callback from registry
	buttonPtr := uintptr(unsafe.Pointer(b.widget))
	buttonRegistry.delete(buttonPtr)

	// Call base destroy method
	b.BaseWidget.Destroy()
}