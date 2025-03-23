// Package gtk4 provides box layout functionality for GTK4
// File: gtk4go/gtk4/box.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// Orientation defines the orientation of a Box
type Orientation int

const (
	// OrientationHorizontal for horizontal layout
	OrientationHorizontal Orientation = C.GTK_ORIENTATION_HORIZONTAL
	// OrientationVertical for vertical layout
	OrientationVertical Orientation = C.GTK_ORIENTATION_VERTICAL
)

// Box represents a GTK box container
type Box struct {
	widget *C.GtkWidget
}

// NewBox creates a new GTK box with the given orientation
func NewBox(orientation Orientation, spacing int) *Box {
	box := &Box{
		widget: C.gtk_box_new(C.GtkOrientation(orientation), C.int(spacing)),
	}
	runtime.SetFinalizer(box, (*Box).Destroy)
	return box
}

// Append adds a widget to the end of the box
func (b *Box) Append(child interface{}) {
	if c, ok := child.(interface{ GetWidget() *C.GtkWidget }); ok {
		C.gtk_box_append((*C.GtkBox)(unsafe.Pointer(b.widget)), c.GetWidget())
	}
}

// Prepend adds a widget to the start of the box
func (b *Box) Prepend(child interface{}) {
	if c, ok := child.(interface{ GetWidget() *C.GtkWidget }); ok {
		C.gtk_box_prepend((*C.GtkBox)(unsafe.Pointer(b.widget)), c.GetWidget())
	}
}

// Remove removes a widget from the box
func (b *Box) Remove(child interface{}) {
	if c, ok := child.(interface{ GetWidget() *C.GtkWidget }); ok {
		C.gtk_box_remove((*C.GtkBox)(unsafe.Pointer(b.widget)), c.GetWidget())
	}
}

// SetSpacing sets the spacing between children
func (b *Box) SetSpacing(spacing int) {
	C.gtk_box_set_spacing((*C.GtkBox)(unsafe.Pointer(b.widget)), C.int(spacing))
}

// SetHomogeneous sets whether all children get the same space
func (b *Box) SetHomogeneous(homogeneous bool) {
	boolVal := C.gboolean(0)
	if homogeneous {
		boolVal = C.gboolean(1)
	}
	C.gtk_box_set_homogeneous((*C.GtkBox)(unsafe.Pointer(b.widget)), boolVal)
}

// Destroy destroys the box
func (b *Box) Destroy() {
	C.gtk_widget_unparent(b.widget)
	b.widget = nil
}

// Native returns the underlying GtkWidget pointer
func (b *Box) Native() uintptr {
	return uintptr(unsafe.Pointer(b.widget))
}

// GetWidget returns the underlying GtkWidget pointer
func (b *Box) GetWidget() *C.GtkWidget {
	return b.widget
}
