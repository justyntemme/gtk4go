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
	"sync/atomic"
	"unsafe"
)

// ButtonClickedCallback represents a callback for button clicked events
type ButtonClickedCallback func()

// Thread-safe callback registry using atomic.Value for lock-free reads
type buttonCallbackRegistry struct {
	// Use atomic.Value to store map - allows lock-free reads
	callbacks atomic.Value
}

// newButtonCallbackRegistry creates a new registry
func newButtonCallbackRegistry() *buttonCallbackRegistry {
	r := &buttonCallbackRegistry{}
	r.callbacks.Store(make(map[uintptr]ButtonClickedCallback))
	return r
}

// set adds or updates a callback in the registry (thread-safe)
func (r *buttonCallbackRegistry) set(key uintptr, callback ButtonClickedCallback) {
	// Create a new map with the updated value - copy-on-write pattern
	current := r.callbacks.Load().(map[uintptr]ButtonClickedCallback)
	updated := make(map[uintptr]ButtonClickedCallback, len(current)+1)

	// Copy existing entries
	for k, v := range current {
		updated[k] = v
	}

	// Add or update the entry
	updated[key] = callback

	// Atomically replace the map
	r.callbacks.Store(updated)
}

// get retrieves a callback by key (lock-free, thread-safe read)
func (r *buttonCallbackRegistry) get(key uintptr) (ButtonClickedCallback, bool) {
	current := r.callbacks.Load().(map[uintptr]ButtonClickedCallback)
	callback, ok := current[key]
	return callback, ok
}

// delete removes a callback from the registry (thread-safe)
func (r *buttonCallbackRegistry) delete(key uintptr) {
	current := r.callbacks.Load().(map[uintptr]ButtonClickedCallback)

	// Check if key exists
	if _, exists := current[key]; !exists {
		return
	}

	// Copy-on-write for deletion
	updated := make(map[uintptr]ButtonClickedCallback, len(current)-1)
	for k, v := range current {
		if k != key {
			updated[k] = v
		}
	}

	r.callbacks.Store(updated)
}

// Global button callback registry
var buttonRegistry = newButtonCallbackRegistry()

//export buttonClickedCallback
func buttonClickedCallback(button *C.GtkButton, userData C.gpointer) {
	// Convert button pointer to uintptr for lookup
	buttonPtr := uintptr(unsafe.Pointer(button))

	// Get callback without locking - atomic read
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
	// Store callback in registry using atomic.Value (no locks needed for reads)
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
