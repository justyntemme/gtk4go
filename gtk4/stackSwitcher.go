// Package gtk4 provides stack switcher functionality for GTK4
// File: gtk4go/gtk4/stackSwitcher.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// StackSwitcherOption is a function that configures a stack switcher
type StackSwitcherOption func(*StackSwitcher)

// StackSwitcher represents a GTK stack switcher
type StackSwitcher struct {
	BaseWidget
}

// NewStackSwitcher creates a new GTK stack switcher
func NewStackSwitcher(stack *Stack, options ...StackSwitcherOption) *StackSwitcher {
	switcher := &StackSwitcher{
		BaseWidget: BaseWidget{
			widget: C.gtk_stack_switcher_new(),
		},
	}

	// Set the stack if provided
	if stack != nil {
		C.gtk_stack_switcher_set_stack(
			(*C.GtkStackSwitcher)(unsafe.Pointer(switcher.widget)),
			(*C.GtkStack)(unsafe.Pointer(stack.widget)),
		)
	}

	// Apply options
	for _, option := range options {
		option(switcher)
	}

	SetupFinalization(switcher, switcher.Destroy)
	return switcher
}

// WithStack sets the stack to control
func WithStack(stack *Stack) StackSwitcherOption {
	return func(ss *StackSwitcher) {
		if stack != nil {
			C.gtk_stack_switcher_set_stack(
				(*C.GtkStackSwitcher)(unsafe.Pointer(ss.widget)),
				(*C.GtkStack)(unsafe.Pointer(stack.widget)),
			)
		}
	}
}

// SetStack sets the stack to control
func (ss *StackSwitcher) SetStack(stack *Stack) {
	if stack != nil {
		C.gtk_stack_switcher_set_stack(
			(*C.GtkStackSwitcher)(unsafe.Pointer(ss.widget)),
			(*C.GtkStack)(unsafe.Pointer(stack.widget)),
		)
	}
}

// GetStack gets the stack being controlled
func (ss *StackSwitcher) GetStack() *Stack {
	stackPtr := C.gtk_stack_switcher_get_stack((*C.GtkStackSwitcher)(unsafe.Pointer(ss.widget)))
	if stackPtr == nil {
		return nil
	}

	// Create a new Stack instance wrapping the C pointer
	return &Stack{
		BaseWidget: BaseWidget{
			widget: (*C.GtkWidget)(unsafe.Pointer(stackPtr)),
		},
	}
}
