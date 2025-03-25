// Package gtk4 provides box layout functionality for GTK4
// File: gtk4go/gtk4/box.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
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

// BoxOption is a function that configures a box
type BoxOption func(*Box)

// Box represents a GTK box container
type Box struct {
	BaseWidget
}

// NewBox creates a new GTK box with the given orientation
func NewBox(orientation Orientation, spacing int, options ...BoxOption) *Box {
	box := &Box{
		BaseWidget: BaseWidget{
			widget: C.gtk_box_new(C.GtkOrientation(orientation), C.int(spacing)),
		},
	}

	// Apply options
	for _, option := range options {
		option(box)
	}

	SetupFinalization(box, box.Destroy)
	return box
}

// WithSpacing sets the spacing between children
func WithSpacing(spacing int) BoxOption {
	return func(b *Box) {
		C.gtk_box_set_spacing((*C.GtkBox)(unsafe.Pointer(b.widget)), C.int(spacing))
	}
}

// WithHomogeneous sets whether all children get the same space
func WithHomogeneous(homogeneous bool) BoxOption {
	return func(b *Box) {
		var chomogeneous C.gboolean
		if homogeneous {
			chomogeneous = C.TRUE
		} else {
			chomogeneous = C.FALSE
		}
		C.gtk_box_set_homogeneous((*C.GtkBox)(unsafe.Pointer(b.widget)), chomogeneous)
	}
}

// Append adds a widget to the end of the box
func (b *Box) Append(child Widget) {
	C.gtk_box_append((*C.GtkBox)(unsafe.Pointer(b.widget)), child.GetWidget())
}

// Prepend adds a widget to the start of the box
func (b *Box) Prepend(child Widget) {
	C.gtk_box_prepend((*C.GtkBox)(unsafe.Pointer(b.widget)), child.GetWidget())
}

// Remove removes a widget from the box
func (b *Box) Remove(child Widget) {
	C.gtk_box_remove((*C.GtkBox)(unsafe.Pointer(b.widget)), child.GetWidget())
}

// SetSpacing sets the spacing between children
func (b *Box) SetSpacing(spacing int) {
	C.gtk_box_set_spacing((*C.GtkBox)(unsafe.Pointer(b.widget)), C.int(spacing))
}

// SetHomogeneous sets whether all children get the same space
func (b *Box) SetHomogeneous(homogeneous bool) {
	var chomogeneous C.gboolean
	if homogeneous {
		chomogeneous = C.TRUE
	} else {
		chomogeneous = C.FALSE
	}
	C.gtk_box_set_homogeneous((*C.GtkBox)(unsafe.Pointer(b.widget)), chomogeneous)
}

// SetHExpand sets whether the box expands horizontally
func (b *Box) SetHExpand(expand bool) {
	var cexpand C.gboolean
	if expand {
		cexpand = C.TRUE
	} else {
		cexpand = C.FALSE
	}
	C.gtk_widget_set_hexpand(b.widget, cexpand)
}

// SetVExpand sets whether the box expands vertically
func (b *Box) SetVExpand(expand bool) {
	var cexpand C.gboolean
	if expand {
		cexpand = C.TRUE
	} else {
		cexpand = C.FALSE
	}
	C.gtk_widget_set_vexpand(b.widget, cexpand)
}