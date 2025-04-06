// Package gtk4 provides scrolled window functionality for GTK4
// File: gtk4go/gtk4/scrolledWindow.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// ScrollbarPolicyType defines the visibility policy for scrollbars
type ScrollbarPolicyType int

const (
	// ScrollbarPolicyAlways always shows scrollbars
	ScrollbarPolicyAlways ScrollbarPolicyType = C.GTK_POLICY_ALWAYS
	// ScrollbarPolicyAutomatic shows scrollbars when needed
	ScrollbarPolicyAutomatic ScrollbarPolicyType = C.GTK_POLICY_AUTOMATIC
	// ScrollbarPolicyNever never shows scrollbars
	ScrollbarPolicyNever ScrollbarPolicyType = C.GTK_POLICY_NEVER
)

// ScrolledWindowOption is a function that configures a scrolled window
type ScrolledWindowOption func(*ScrolledWindow)

// ScrolledWindow represents a GTK scrolled window container
type ScrolledWindow struct {
	BaseWidget
}

// NewScrolledWindow creates a new GTK scrolled window container
func NewScrolledWindow(options ...ScrolledWindowOption) *ScrolledWindow {
	scrolledWindow := &ScrolledWindow{
		BaseWidget: BaseWidget{
			widget: C.gtk_scrolled_window_new(),
		},
	}

	// Apply options
	for _, option := range options {
		option(scrolledWindow)
	}

	SetupFinalization(scrolledWindow, scrolledWindow.Destroy)
	return scrolledWindow
}

// WithHScrollbarPolicy sets the horizontal scrollbar policy
func WithHScrollbarPolicy(policy ScrollbarPolicyType) ScrolledWindowOption {
	return func(sw *ScrolledWindow) {
		sw.SetHScrollbarPolicy(policy)
	}
}

// WithVScrollbarPolicy sets the vertical scrollbar policy
func WithVScrollbarPolicy(policy ScrollbarPolicyType) ScrolledWindowOption {
	return func(sw *ScrolledWindow) {
		sw.SetVScrollbarPolicy(policy)
	}
}

// WithPropagateNaturalWidth sets whether to propagate the natural width of the child
func WithPropagateNaturalWidth(propagate bool) ScrolledWindowOption {
	return func(sw *ScrolledWindow) {
		var cpropagate C.gboolean
		if propagate {
			cpropagate = C.TRUE
		} else {
			cpropagate = C.FALSE
		}
		C.gtk_scrolled_window_set_propagate_natural_width(
			(*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)),
			cpropagate,
		)
	}
}

// WithPropagateNaturalHeight sets whether to propagate the natural height of the child
func WithPropagateNaturalHeight(propagate bool) ScrolledWindowOption {
	return func(sw *ScrolledWindow) {
		var cpropagate C.gboolean
		if propagate {
			cpropagate = C.TRUE
		} else {
			cpropagate = C.FALSE
		}
		C.gtk_scrolled_window_set_propagate_natural_height(
			(*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)),
			cpropagate,
		)
	}
}

// WithHExpand sets whether the scrolled window expands horizontally
func WithHExpand(expand bool) ScrolledWindowOption {
	return func(sw *ScrolledWindow) {
		sw.SetHExpand(expand)
	}
}

// WithVExpand sets whether the scrolled window expands vertically
func WithVExpand(expand bool) ScrolledWindowOption {
	return func(sw *ScrolledWindow) {
		sw.SetVExpand(expand)
	}
}

// SetChild sets the child widget of the scrolled window
func (sw *ScrolledWindow) SetChild(child Widget) {
	C.gtk_scrolled_window_set_child(
		(*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)),
		child.GetWidget(),
	)
}

// GetChild gets the child widget of the scrolled window
func (sw *ScrolledWindow) GetChild() Widget {
	widget := C.gtk_scrolled_window_get_child((*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)))

	// Similar to Grid.GetChildAt, we need a proper implementation to convert to Widget
	if widget == nil {
		return nil
	}
	return nil
}

// SetPolicy sets the policy for both horizontal and vertical scrollbars
func (sw *ScrolledWindow) SetPolicy(hPolicy, vPolicy ScrollbarPolicyType) {
	C.gtk_scrolled_window_set_policy(
		(*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)),
		C.GtkPolicyType(hPolicy),
		C.GtkPolicyType(vPolicy),
	)
}

// SetHScrollbarPolicy sets the horizontal scrollbar policy
func (sw *ScrolledWindow) SetHScrollbarPolicy(policy ScrollbarPolicyType) {
	_, vPolicy := sw.GetPolicy()
	sw.SetPolicy(policy, vPolicy)
}

// SetVScrollbarPolicy sets the vertical scrollbar policy
func (sw *ScrolledWindow) SetVScrollbarPolicy(policy ScrollbarPolicyType) {
	hPolicy, _ := sw.GetPolicy()
	sw.SetPolicy(hPolicy, policy)
}

// GetPolicy gets the policy for both horizontal and vertical scrollbars
func (sw *ScrolledWindow) GetPolicy() (hPolicy, vPolicy ScrollbarPolicyType) {
	var chPolicy, cvPolicy C.GtkPolicyType
	C.gtk_scrolled_window_get_policy(
		(*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)),
		&chPolicy,
		&cvPolicy,
	)
	return ScrollbarPolicyType(chPolicy), ScrollbarPolicyType(cvPolicy)
}

// SetPropagateNaturalWidth sets whether to propagate the natural width of the child
func (sw *ScrolledWindow) SetPropagateNaturalWidth(propagate bool) {
	var cpropagate C.gboolean
	if propagate {
		cpropagate = C.TRUE
	} else {
		cpropagate = C.FALSE
	}
	C.gtk_scrolled_window_set_propagate_natural_width(
		(*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)),
		cpropagate,
	)
}

// GetPropagateNaturalWidth gets whether the natural width of the child is propagated
func (sw *ScrolledWindow) GetPropagateNaturalWidth() bool {
	return C.gtk_scrolled_window_get_propagate_natural_width(
		(*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)),
	) == C.TRUE
}

// SetPropagateNaturalHeight sets whether to propagate the natural height of the child
func (sw *ScrolledWindow) SetPropagateNaturalHeight(propagate bool) {
	var cpropagate C.gboolean
	if propagate {
		cpropagate = C.TRUE
	} else {
		cpropagate = C.FALSE
	}
	C.gtk_scrolled_window_set_propagate_natural_height(
		(*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)),
		cpropagate,
	)
}

// GetPropagateNaturalHeight gets whether the natural height of the child is propagated
func (sw *ScrolledWindow) GetPropagateNaturalHeight() bool {
	return C.gtk_scrolled_window_get_propagate_natural_height(
		(*C.GtkScrolledWindow)(unsafe.Pointer(sw.widget)),
	) == C.TRUE
}