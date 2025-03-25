// Package gtk4 provides button functionality for GTK4
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

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
	if cLabel == nil {
		return ""
	}
	return C.GoString(cLabel)
}

// ConnectClicked connects a callback function to the button's "clicked" signal
func (b *Button) ConnectClicked(callback func()) {
	// Use the new callback system from callbacks.go
	Connect(b, SignalClicked, callback)
}

// DisconnectClicked disconnects all clicked signal handlers
func (b *Button) DisconnectClicked() {
	// Use the DisconnectAll function from the unified callback system
	// for the specific signal
	// TODO: If needed, we could add a method to the callback system to disconnect by signal
	DisconnectAll(b)
}

// Destroy destroys the button and cleans up resources
func (b *Button) Destroy() {
	// Disconnect all signals for this widget
	DisconnectAll(b)

	// Call base destroy method
	b.BaseWidget.Destroy()
}

// SetIconName sets the icon for the button
func (b *Button) SetIconName(iconName string) {
	WithCString(iconName, func(cIconName *C.char) {
		// Create a new image with the icon name
		image := C.gtk_image_new_from_icon_name(cIconName)

		// Set the image on the button
		C.gtk_button_set_child((*C.GtkButton)(unsafe.Pointer(b.widget)), image)
	})
}

// SetChild sets the child widget for the button (instead of a label)
func (b *Button) SetChild(child Widget) {
	C.gtk_button_set_child((*C.GtkButton)(unsafe.Pointer(b.widget)), child.GetWidget())
}

// GetChild gets the child widget of the button
func (b *Button) GetChild() Widget {
	widget := C.gtk_button_get_child((*C.GtkButton)(unsafe.Pointer(b.widget)))
	if widget == nil {
		return nil
	}

	// Create a BaseWidget wrapper for the child
	// Note: In a real implementation, we would determine the widget type
	return &BaseWidget{widget: widget}
}

// SetHasFrame sets whether the button has a visible frame
func (b *Button) SetHasFrame(hasFrame bool) {
	var cHasFrame C.gboolean
	if hasFrame {
		cHasFrame = C.TRUE
	} else {
		cHasFrame = C.FALSE
	}
	C.gtk_button_set_has_frame((*C.GtkButton)(unsafe.Pointer(b.widget)), cHasFrame)
}

// GetHasFrame gets whether the button has a visible frame
func (b *Button) GetHasFrame() bool {
	return C.gtk_button_get_has_frame((*C.GtkButton)(unsafe.Pointer(b.widget))) == C.TRUE
}

