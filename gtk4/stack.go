// Package gtk4 provides stack container functionality for GTK4
// File: gtk4go/gtk4/stack.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// StackTransitionType defines the type of animation used when transitioning between pages in a Stack
type StackTransitionType int

const (
	// StackTransitionTypeNone no transition
	StackTransitionTypeNone StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_NONE
	// StackTransitionTypeCrossfade crossfade transition
	StackTransitionTypeCrossfade StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_CROSSFADE
	// StackTransitionTypeSlideRight slide from left to right
	StackTransitionTypeSlideRight StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_SLIDE_RIGHT
	// StackTransitionTypeSlideLeft slide from right to left
	StackTransitionTypeSlideLeft StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_SLIDE_LEFT
	// StackTransitionTypeSlideUp slide from bottom to top
	StackTransitionTypeSlideUp StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_SLIDE_UP
	// StackTransitionTypeSlideDown slide from top to bottom
	StackTransitionTypeSlideDown StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_SLIDE_DOWN
	// StackTransitionTypeSlideLeftRight slide from left or right according to the children order
	StackTransitionTypeSlideLeftRight StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_SLIDE_LEFT_RIGHT
	// StackTransitionTypeSlideUpDown slide from top or bottom according to the children order
	StackTransitionTypeSlideUpDown StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_SLIDE_UP_DOWN
	// StackTransitionTypeOverUp slide the new widget over the old one from bottom to top
	StackTransitionTypeOverUp StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_OVER_UP
	// StackTransitionTypeOverDown slide the new widget over the old one from top to bottom
	StackTransitionTypeOverDown StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_OVER_DOWN
	// StackTransitionTypeOverLeft slide the new widget over the old one from right to left
	StackTransitionTypeOverLeft StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_OVER_LEFT
	// StackTransitionTypeOverRight slide the new widget over the old one from left to right
	StackTransitionTypeOverRight StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_OVER_RIGHT
	// StackTransitionTypeUnderUp slide the old widget under the new one from bottom to top
	StackTransitionTypeUnderUp StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_UNDER_UP
	// StackTransitionTypeUnderDown slide the old widget under the new one from top to bottom
	StackTransitionTypeUnderDown StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_UNDER_DOWN
	// StackTransitionTypeUnderLeft slide the old widget under the new one from right to left
	StackTransitionTypeUnderLeft StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_UNDER_LEFT
	// StackTransitionTypeUnderRight slide the old widget under the new one from left to right
	StackTransitionTypeUnderRight StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_UNDER_RIGHT
	// StackTransitionTypeOverUpDown slide the new widget over the old one from bottom or top according to the order
	StackTransitionTypeOverUpDown StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_OVER_UP_DOWN
	// StackTransitionTypeOverDownUp slide the new widget over the old one from top or bottom according to the order
	StackTransitionTypeOverDownUp StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_OVER_DOWN_UP
	// StackTransitionTypeOverLeftRight slide the new widget over the old one from right or left according to the order
	StackTransitionTypeOverLeftRight StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_OVER_LEFT_RIGHT
	// StackTransitionTypeOverRightLeft slide the new widget over the old one from left or right according to the order
	StackTransitionTypeOverRightLeft StackTransitionType = C.GTK_STACK_TRANSITION_TYPE_OVER_RIGHT_LEFT
)

// StackOption is a function that configures a stack
type StackOption func(*Stack)

// Stack represents a GTK stack container
type Stack struct {
	BaseWidget
}

// NewStack creates a new GTK stack container
func NewStack(options ...StackOption) *Stack {
	stack := &Stack{
		BaseWidget: BaseWidget{
			widget: C.gtk_stack_new(),
		},
	}

	// Apply options
	for _, option := range options {
		option(stack)
	}

	SetupFinalization(stack, stack.Destroy)
	return stack
}

// WithTransitionType sets the transition type
func WithTransitionType(transitionType StackTransitionType) StackOption {
	return func(s *Stack) {
		C.gtk_stack_set_transition_type(
			(*C.GtkStack)(unsafe.Pointer(s.widget)),
			C.GtkStackTransitionType(transitionType),
		)
	}
}

// WithTransitionDuration sets the transition duration
func WithTransitionDuration(duration uint) StackOption {
	return func(s *Stack) {
		C.gtk_stack_set_transition_duration(
			(*C.GtkStack)(unsafe.Pointer(s.widget)),
			C.guint(duration),
		)
	}
}

// WithHHomogeneous sets whether all children have the same width
func WithHHomogeneous(homogeneous bool) StackOption {
	return func(s *Stack) {
		var chomogeneous C.gboolean
		if homogeneous {
			chomogeneous = C.TRUE
		} else {
			chomogeneous = C.FALSE
		}
		C.gtk_stack_set_hhomogeneous((*C.GtkStack)(unsafe.Pointer(s.widget)), chomogeneous)
	}
}

// WithVHomogeneous sets whether all children have the same height
func WithVHomogeneous(homogeneous bool) StackOption {
	return func(s *Stack) {
		var chomogeneous C.gboolean
		if homogeneous {
			chomogeneous = C.TRUE
		} else {
			chomogeneous = C.FALSE
		}
		C.gtk_stack_set_vhomogeneous((*C.GtkStack)(unsafe.Pointer(s.widget)), chomogeneous)
	}
}

// AddNamed adds a child to the stack with the given name
func (s *Stack) AddNamed(child Widget, name string) {
	WithCString(name, func(cName *C.char) {
		C.gtk_stack_add_named(
			(*C.GtkStack)(unsafe.Pointer(s.widget)),
			child.GetWidget(),
			cName,
		)
	})
}

// AddTitled adds a child to the stack with the given name and title
func (s *Stack) AddTitled(child Widget, name, title string) {
	WithCString(name, func(cName *C.char) {
		WithCString(title, func(cTitle *C.char) {
			C.gtk_stack_add_titled(
				(*C.GtkStack)(unsafe.Pointer(s.widget)),
				child.GetWidget(),
				cName,
				cTitle,
			)
		})
	})
}

// Remove removes a child from the stack
func (s *Stack) Remove(child Widget) {
	C.gtk_stack_remove((*C.GtkStack)(unsafe.Pointer(s.widget)), child.GetWidget())
}

// GetChildByName gets the child with the given name
func (s *Stack) GetChildByName(name string) Widget {
	var widget *C.GtkWidget
	WithCString(name, func(cName *C.char) {
		widget = C.gtk_stack_get_child_by_name(
			(*C.GtkStack)(unsafe.Pointer(s.widget)),
			cName,
		)
	})

	if widget == nil {
		return nil
	}

	// Similar to Grid.GetChildAt, we need a proper implementation to convert to Widget
	return nil
}

// SetVisibleChild sets the visible child by widget reference
func (s *Stack) SetVisibleChild(child Widget) {
	C.gtk_stack_set_visible_child(
		(*C.GtkStack)(unsafe.Pointer(s.widget)),
		child.GetWidget(),
	)
}

// SetVisibleChildName sets the visible child by name
func (s *Stack) SetVisibleChildName(name string) {
	WithCString(name, func(cName *C.char) {
		C.gtk_stack_set_visible_child_name(
			(*C.GtkStack)(unsafe.Pointer(s.widget)),
			cName,
		)
	})
}

// GetVisibleChildName gets the name of the visible child
func (s *Stack) GetVisibleChildName() string {
	cName := C.gtk_stack_get_visible_child_name((*C.GtkStack)(unsafe.Pointer(s.widget)))
	if cName == nil {
		return ""
	}
	return C.GoString(cName)
}

// SetTransitionType sets the type of animation used for transitions
func (s *Stack) SetTransitionType(transitionType StackTransitionType) {
	C.gtk_stack_set_transition_type(
		(*C.GtkStack)(unsafe.Pointer(s.widget)),
		C.GtkStackTransitionType(transitionType),
	)
}

// GetTransitionType gets the type of animation used for transitions
func (s *Stack) GetTransitionType() StackTransitionType {
	return StackTransitionType(C.gtk_stack_get_transition_type((*C.GtkStack)(unsafe.Pointer(s.widget))))
}

// SetTransitionDuration sets the duration of the transition in milliseconds
func (s *Stack) SetTransitionDuration(duration uint) {
	C.gtk_stack_set_transition_duration(
		(*C.GtkStack)(unsafe.Pointer(s.widget)),
		C.guint(duration),
	)
}

// GetTransitionDuration gets the duration of the transition in milliseconds
func (s *Stack) GetTransitionDuration() uint {
	return uint(C.gtk_stack_get_transition_duration((*C.GtkStack)(unsafe.Pointer(s.widget))))
}

// SetInterpolateSize sets whether the stack should interpolate its size during transitions
func (s *Stack) SetInterpolateSize(interpolate bool) {
	var cinterpolate C.gboolean
	if interpolate {
		cinterpolate = C.TRUE
	} else {
		cinterpolate = C.FALSE
	}
	C.gtk_stack_set_interpolate_size((*C.GtkStack)(unsafe.Pointer(s.widget)), cinterpolate)
}

// GetInterpolateSize gets whether the stack is set to interpolate its size during transitions
func (s *Stack) GetInterpolateSize() bool {
	return C.gtk_stack_get_interpolate_size((*C.GtkStack)(unsafe.Pointer(s.widget))) == C.TRUE
}

// SetHHomogeneous sets whether the stack allocates the same width for all children
func (s *Stack) SetHHomogeneous(homogeneous bool) {
	var chomogeneous C.gboolean
	if homogeneous {
		chomogeneous = C.TRUE
	} else {
		chomogeneous = C.FALSE
	}
	C.gtk_stack_set_hhomogeneous((*C.GtkStack)(unsafe.Pointer(s.widget)), chomogeneous)
}

// GetHHomogeneous gets whether the stack allocates the same width for all children
func (s *Stack) GetHHomogeneous() bool {
	return C.gtk_stack_get_hhomogeneous((*C.GtkStack)(unsafe.Pointer(s.widget))) == C.TRUE
}

// SetVHomogeneous sets whether the stack allocates the same height for all children
func (s *Stack) SetVHomogeneous(homogeneous bool) {
	var chomogeneous C.gboolean
	if homogeneous {
		chomogeneous = C.TRUE
	} else {
		chomogeneous = C.FALSE
	}
	C.gtk_stack_set_vhomogeneous((*C.GtkStack)(unsafe.Pointer(s.widget)), chomogeneous)
}

// GetVHomogeneous gets whether the stack allocates the same height for all children
func (s *Stack) GetVHomogeneous() bool {
	return C.gtk_stack_get_vhomogeneous((*C.GtkStack)(unsafe.Pointer(s.widget))) == C.TRUE
}
