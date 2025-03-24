// Package gtk4 provides viewport widget functionality for GTK4
// File: gtk4go/gtk4/viewport.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// ViewportOption is a function that configures a viewport
type ViewportOption func(*Viewport)

// Viewport represents a GTK viewport container for scrolling
type Viewport struct {
	BaseWidget
}

// NewViewport creates a new GTK viewport container
func NewViewport(options ...ViewportOption) *Viewport {
	// Pass nil for both adjustments - GTK will create default ones
	var hadjustment, vadjustment *C.GtkAdjustment
	hadjustment = nil
	vadjustment = nil

	viewport := &Viewport{
		BaseWidget: BaseWidget{
			widget: C.gtk_viewport_new(hadjustment, vadjustment),
		},
	}

	// Apply options
	for _, option := range options {
		option(viewport)
	}

	SetupFinalization(viewport, viewport.Destroy)
	return viewport
}

// WithScrollToFocus sets whether the viewport should bring a widget into view when it receives focus
func WithScrollToFocus(scroll bool) ViewportOption {
	return func(v *Viewport) {
		var cscroll C.gboolean
		if scroll {
			cscroll = C.TRUE
		} else {
			cscroll = C.FALSE
		}
		C.gtk_viewport_set_scroll_to_focus((*C.GtkViewport)(unsafe.Pointer(v.widget)), cscroll)
	}
}

// SetChild sets the child widget of the viewport
func (v *Viewport) SetChild(child Widget) {
	C.gtk_viewport_set_child((*C.GtkViewport)(unsafe.Pointer(v.widget)), child.GetWidget())
}

// GetChild gets the child widget of the viewport
func (v *Viewport) GetChild() Widget {
	widget := C.gtk_viewport_get_child((*C.GtkViewport)(unsafe.Pointer(v.widget)))

	// Similar to Grid.GetChildAt, we need a proper implementation to convert to Widget
	if widget == nil {
		return nil
	}
	return nil
}

// SetScrollToFocus sets whether the viewport should bring a widget into view when it receives focus
func (v *Viewport) SetScrollToFocus(scroll bool) {
	var cscroll C.gboolean
	if scroll {
		cscroll = C.TRUE
	} else {
		cscroll = C.FALSE
	}
	C.gtk_viewport_set_scroll_to_focus((*C.GtkViewport)(unsafe.Pointer(v.widget)), cscroll)
}

// GetScrollToFocus gets whether the viewport brings a widget into view when it receives focus
func (v *Viewport) GetScrollToFocus() bool {
	return C.gtk_viewport_get_scroll_to_focus((*C.GtkViewport)(unsafe.Pointer(v.widget))) == C.TRUE
}

// ScrollTo scrolls the viewport to make the child widget visible
// Note: In GTK4, scrolling to a specific child is typically handled through
// the parent ScrolledWindow. This is a helper method that finds the
// appropriate scroll adjustments.
func (v *Viewport) ScrollTo(widget Widget) {
	// In GTK4, we need to use the scroll adjustments
	// This is a simplified implementation
	// The exact implementation would depend on the specific use case

	// This method works indirectly by requesting the widget to grab focus,
	// which will trigger scroll_to_focus if it's enabled
	C.gtk_widget_grab_focus(widget.GetWidget())
}
