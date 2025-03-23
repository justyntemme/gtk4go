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
	"runtime"
	"sync"
	"unsafe"
)

// ButtonClickedCallback represents a callback for button clicked events
type ButtonClickedCallback func()

var (
	buttonCallbacks     = make(map[uintptr]ButtonClickedCallback)
	buttonCallbackMutex sync.Mutex
)

//export buttonClickedCallback
func buttonClickedCallback(button *C.GtkButton, userData C.gpointer) {
	buttonCallbackMutex.Lock()
	defer buttonCallbackMutex.Unlock()

	// Convert button pointer to uintptr for lookup
	buttonPtr := uintptr(unsafe.Pointer(button))

	// Find and call the callback
	if callback, ok := buttonCallbacks[buttonPtr]; ok {
		callback()
	}
}

// Button represents a GTK button
type Button struct {
	widget *C.GtkWidget
}

// NewButton creates a new GTK button with the given label
func NewButton(label string) *Button {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))

	button := &Button{
		widget: C.gtk_button_new_with_label(cLabel),
	}
	runtime.SetFinalizer(button, (*Button).Destroy)
	return button
}

// NewButtonWithMnemonic creates a new GTK button with a mnemonic label
func NewButtonWithMnemonic(label string) *Button {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))

	button := &Button{
		widget: C.gtk_button_new_with_mnemonic(cLabel),
	}
	runtime.SetFinalizer(button, (*Button).Destroy)
	return button
}

// SetLabel sets the button's label
func (b *Button) SetLabel(label string) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	C.gtk_button_set_label((*C.GtkButton)(unsafe.Pointer(b.widget)), cLabel)
}

// GetLabel gets the button's label
func (b *Button) GetLabel() string {
	cLabel := C.gtk_button_get_label((*C.GtkButton)(unsafe.Pointer(b.widget)))
	return C.GoString(cLabel)
}

// ConnectClicked connects a callback function to the button's "clicked" signal
func (b *Button) ConnectClicked(callback ButtonClickedCallback) {
	buttonCallbackMutex.Lock()
	defer buttonCallbackMutex.Unlock()

	// Store callback in map
	buttonPtr := uintptr(unsafe.Pointer(b.widget))
	buttonCallbacks[buttonPtr] = callback

	// Connect signal
	C.connectButtonClicked(b.widget, C.gpointer(unsafe.Pointer(b.widget)))
}

// Destroy destroys the button
func (b *Button) Destroy() {
	buttonCallbackMutex.Lock()
	defer buttonCallbackMutex.Unlock()

	// Remove callback from map if exists
	buttonPtr := uintptr(unsafe.Pointer(b.widget))
	delete(buttonCallbacks, buttonPtr)

	// Destroy widget
	C.gtk_widget_unparent(b.widget)
	b.widget = nil
}

// Native returns the underlying GtkWidget pointer
func (b *Button) Native() uintptr {
	return uintptr(unsafe.Pointer(b.widget))
}

// GetWidget returns the underlying GtkWidget pointer
func (b *Button) GetWidget() *C.GtkWidget {
	return b.widget
}
