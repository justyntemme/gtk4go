// Package gtk4 provides paned container functionality for GTK4
// File: gtk4go/gtk4/paned.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// PanedOption is a function that configures a paned container
type PanedOption func(*Paned)

// Paned represents a GTK paned container (a split view)
type Paned struct {
	BaseWidget
}

// NewPaned creates a new GTK paned container with the given orientation
func NewPaned(orientation Orientation, options ...PanedOption) *Paned {
	paned := &Paned{
		BaseWidget: BaseWidget{
			widget: C.gtk_paned_new(C.GtkOrientation(orientation)),
		},
	}

	// Apply options
	for _, option := range options {
		option(paned)
	}

	SetupFinalization(paned, paned.Destroy)
	return paned
}

// WithPosition sets the initial position of the divider
func WithPosition(position int) PanedOption {
	return func(p *Paned) {
		C.gtk_paned_set_position((*C.GtkPaned)(unsafe.Pointer(p.widget)), C.int(position))
	}
}

// WithWideHandle sets whether to use a wide handle
func WithWideHandle(wide bool) PanedOption {
	return func(p *Paned) {
		var cwide C.gboolean
		if wide {
			cwide = C.TRUE
		} else {
			cwide = C.FALSE
		}
		C.gtk_paned_set_wide_handle((*C.GtkPaned)(unsafe.Pointer(p.widget)), cwide)
	}
}

// SetStartChild sets the start (top/left) child widget
func (p *Paned) SetStartChild(child Widget) {
	C.gtk_paned_set_start_child((*C.GtkPaned)(unsafe.Pointer(p.widget)), child.GetWidget())
}

// SetEndChild sets the end (bottom/right) child widget
func (p *Paned) SetEndChild(child Widget) {
	C.gtk_paned_set_end_child((*C.GtkPaned)(unsafe.Pointer(p.widget)), child.GetWidget())
}

// GetStartChild gets the start child widget
func (p *Paned) GetStartChild() Widget {
	widget := C.gtk_paned_get_start_child((*C.GtkPaned)(unsafe.Pointer(p.widget)))
	// Similar to Grid.GetChildAt, we need a proper implementation to convert to Widget
	if widget == nil {
		return nil
	}
	return nil
}

// GetEndChild gets the end child widget
func (p *Paned) GetEndChild() Widget {
	widget := C.gtk_paned_get_end_child((*C.GtkPaned)(unsafe.Pointer(p.widget)))
	// Similar to Grid.GetChildAt, we need a proper implementation to convert to Widget
	if widget == nil {
		return nil
	}
	return nil
}

// SetPosition sets the position of the divider
func (p *Paned) SetPosition(position int) {
	C.gtk_paned_set_position((*C.GtkPaned)(unsafe.Pointer(p.widget)), C.int(position))
}

// GetPosition gets the position of the divider
func (p *Paned) GetPosition() int {
	return int(C.gtk_paned_get_position((*C.GtkPaned)(unsafe.Pointer(p.widget))))
}

// SetWideHandle sets whether to use a wide handle
func (p *Paned) SetWideHandle(wide bool) {
	var cwide C.gboolean
	if wide {
		cwide = C.TRUE
	} else {
		cwide = C.FALSE
	}
	C.gtk_paned_set_wide_handle((*C.GtkPaned)(unsafe.Pointer(p.widget)), cwide)
}

// GetWideHandle gets whether a wide handle is used
func (p *Paned) GetWideHandle() bool {
	return C.gtk_paned_get_wide_handle((*C.GtkPaned)(unsafe.Pointer(p.widget))) == C.TRUE
}

// SetStartChildResizable sets whether the start child is resizable
func (p *Paned) SetStartChildResizable(resizable bool) {
	var cresizable C.gboolean
	if resizable {
		cresizable = C.TRUE
	} else {
		cresizable = C.FALSE
	}
	C.gtk_paned_set_resize_start_child((*C.GtkPaned)(unsafe.Pointer(p.widget)), cresizable)
}

// GetStartChildResizable gets whether the start child is resizable
func (p *Paned) GetStartChildResizable() bool {
	return C.gtk_paned_get_resize_start_child((*C.GtkPaned)(unsafe.Pointer(p.widget))) == C.TRUE
}

// SetEndChildResizable sets whether the end child is resizable
func (p *Paned) SetEndChildResizable(resizable bool) {
	var cresizable C.gboolean
	if resizable {
		cresizable = C.TRUE
	} else {
		cresizable = C.FALSE
	}
	C.gtk_paned_set_resize_end_child((*C.GtkPaned)(unsafe.Pointer(p.widget)), cresizable)
}

// GetEndChildResizable gets whether the end child is resizable
func (p *Paned) GetEndChildResizable() bool {
	return C.gtk_paned_get_resize_end_child((*C.GtkPaned)(unsafe.Pointer(p.widget))) == C.TRUE
}

// SetShrinkStartChild sets whether the start child can be made smaller than its requisition
func (p *Paned) SetShrinkStartChild(shrink bool) {
	var cshrink C.gboolean
	if shrink {
		cshrink = C.TRUE
	} else {
		cshrink = C.FALSE
	}
	C.gtk_paned_set_shrink_start_child((*C.GtkPaned)(unsafe.Pointer(p.widget)), cshrink)
}

// GetShrinkStartChild gets whether the start child can be made smaller than its requisition
func (p *Paned) GetShrinkStartChild() bool {
	return C.gtk_paned_get_shrink_start_child((*C.GtkPaned)(unsafe.Pointer(p.widget))) == C.TRUE
}

// SetShrinkEndChild sets whether the end child can be made smaller than its requisition
func (p *Paned) SetShrinkEndChild(shrink bool) {
	var cshrink C.gboolean
	if shrink {
		cshrink = C.TRUE
	} else {
		cshrink = C.FALSE
	}
	C.gtk_paned_set_shrink_end_child((*C.GtkPaned)(unsafe.Pointer(p.widget)), cshrink)
}

// GetShrinkEndChild gets whether the end child can be made smaller than its requisition
func (p *Paned) GetShrinkEndChild() bool {
	return C.gtk_paned_get_shrink_end_child((*C.GtkPaned)(unsafe.Pointer(p.widget))) == C.TRUE
}
