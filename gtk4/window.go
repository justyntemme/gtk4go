// Package gtk4 provides window functionality for GTK4
// File: gtk4go/gtk4/window.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// WindowOption is a function that configures a window
type WindowOption func(*Window)

// Window represents a GTK window
type Window struct {
	BaseWidget
}

// NewWindow creates a new GTK window with the given title
func NewWindow(title string, options ...WindowOption) *Window {
	window := &Window{
		BaseWidget: BaseWidget{
			widget: C.gtk_window_new(),
		},
	}

	WithCString(title, func(cTitle *C.char) {
		C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(window.widget)), cTitle)
	})

	// Apply default size
	C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(window.widget)), 600, 400)

	// Apply options
	for _, option := range options {
		option(window)
	}

	SetupFinalization(window, window.Destroy)
	return window
}

// WithDefaultSize sets the default window size
func WithDefaultSize(width, height int) WindowOption {
	return func(w *Window) {
		C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(w.widget)), C.int(width), C.int(height))
	}
}

// WithTransientFor sets the parent window
func WithTransientFor(parent *Window) WindowOption {
	return func(w *Window) {
		if parent != nil {
			C.gtk_window_set_transient_for(
				(*C.GtkWindow)(unsafe.Pointer(w.widget)),
				(*C.GtkWindow)(unsafe.Pointer(parent.widget)),
			)
		}
	}
}

// WithModal sets whether the window is modal
func WithModal(modal bool) WindowOption {
	return func(w *Window) {
		var cmodal C.gboolean
		if modal {
			cmodal = C.TRUE
		} else {
			cmodal = C.FALSE
		}
		C.gtk_window_set_modal((*C.GtkWindow)(unsafe.Pointer(w.widget)), cmodal)
	}
}

// SetTitle sets the window title
func (w *Window) SetTitle(title string) {
	WithCString(title, func(cTitle *C.char) {
		C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(w.widget)), cTitle)
	})
}

// SetDefaultSize sets the default window size
func (w *Window) SetDefaultSize(width, height int) {
	C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(w.widget)), C.int(width), C.int(height))
}

// SetChild sets the child widget for the window
func (w *Window) SetChild(child Widget) {
	C.gtk_window_set_child((*C.GtkWindow)(unsafe.Pointer(w.widget)), child.GetWidget())
}

// Show makes the window visible
func (w *Window) Show() {
	C.gtk_widget_set_visible(w.widget, C.TRUE)
}

// Present presents the window to the user (preferred in GTK4)
func (w *Window) Present() {
	C.gtk_window_present((*C.GtkWindow)(unsafe.Pointer(w.widget)))
}

// SetVisible sets the visibility of the window
func (w *Window) SetVisible(visible bool) {
	var cvisible C.gboolean
	if visible {
		cvisible = C.TRUE
	} else {
		cvisible = C.FALSE
	}
	C.gtk_widget_set_visible(w.widget, cvisible)
}

// Destroy destroys the window
func (w *Window) Destroy() {
	if w.widget != nil {
		C.gtk_window_destroy((*C.GtkWindow)(unsafe.Pointer(w.widget)))
		w.widget = nil
	}
}
