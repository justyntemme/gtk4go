// Package gtk4 provides window functionality for GTK4
// File: gtk4go/gtk4/window.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// Window represents a GTK window
type Window struct {
	widget *C.GtkWidget
}

// NewWindow creates a new GTK window with the given title
func NewWindow(title string) *Window {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	window := &Window{
		widget: C.gtk_window_new(),
	}
	C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(window.widget)), cTitle)
	C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(window.widget)), 600, 400)

	runtime.SetFinalizer(window, (*Window).Destroy)
	return window
}

// SetTitle sets the window title
func (w *Window) SetTitle(title string) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(w.widget)), cTitle)
}

// SetDefaultSize sets the default window size
func (w *Window) SetDefaultSize(width, height int) {
	C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(w.widget)), C.int(width), C.int(height))
}

// SetChild sets the child widget for the window
func (w *Window) SetChild(child interface{}) {
	// This is a simple implementation; in a real library you would need to handle
	// different widget types properly
	if c, ok := child.(interface{ GetWidget() *C.GtkWidget }); ok {
		C.gtk_window_set_child((*C.GtkWindow)(unsafe.Pointer(w.widget)), c.GetWidget())
	}
}

// Show makes the window visible (use Present instead in GTK4)
func (w *Window) Show() {
	// Deprecated in GTK4, use SetVisible or Present instead
	C.gtk_widget_set_visible(w.widget, C.TRUE)
}

// Present presents the window to the user (preferred in GTK4)
func (w *Window) Present() {
	C.gtk_window_present((*C.GtkWindow)(unsafe.Pointer(w.widget)))
}

// SetVisible sets the visibility of the window
func (w *Window) SetVisible(visible bool) {
	if visible {
		C.gtk_widget_set_visible(w.widget, C.TRUE)
	} else {
		C.gtk_widget_set_visible(w.widget, C.FALSE)
	}
}

// Destroy destroys the window
func (w *Window) Destroy() {
	C.gtk_window_destroy((*C.GtkWindow)(unsafe.Pointer(w.widget)))
}

// Connect connects a signal to the window
func (w *Window) Connect(signal string, callback interface{}) {
	// A real implementation would handle signals properly
	// This is a placeholder
}

// Native returns the underlying GtkWidget pointer
func (w *Window) Native() uintptr {
	return uintptr(unsafe.Pointer(w.widget))
}

// GetWidget returns the underlying GtkWidget pointer
func (w *Window) GetWidget() *C.GtkWidget {
	return w.widget
}
