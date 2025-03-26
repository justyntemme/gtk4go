// Package gtk4 provides adjustment functionality for GTK4
// File: gtk4go/gtk4/adjustment.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// AdjustmentOption is a function that configures an adjustment
type AdjustmentOption func(*Adjustment)

// Adjustment represents a GTK adjustment
type Adjustment struct {
	adjustment *C.GtkAdjustment
}

// NewAdjustment creates a new GTK adjustment
func NewAdjustment(value, lower, upper, stepIncrement, pageIncrement, pageSize float64, options ...AdjustmentOption) *Adjustment {
	adjustment := &Adjustment{
		adjustment: C.gtk_adjustment_new(
			C.gdouble(value),
			C.gdouble(lower),
			C.gdouble(upper),
			C.gdouble(stepIncrement),
			C.gdouble(pageIncrement),
			C.gdouble(pageSize),
		),
	}

	// Apply options
	for _, option := range options {
		option(adjustment)
	}

	runtime.SetFinalizer(adjustment, (*Adjustment).Free)
	return adjustment
}

// WithValue sets the initial value of the adjustment
func WithValue(value float64) AdjustmentOption {
	return func(a *Adjustment) {
		a.SetValue(value)
	}
}

// WithRange sets the range of the adjustment
func WithRange(lower, upper float64) AdjustmentOption {
	return func(a *Adjustment) {
		a.SetLower(lower)
		a.SetUpper(upper)
	}
}

// WithStepIncrement sets the step increment of the adjustment
func WithStepIncrement(step float64) AdjustmentOption {
	return func(a *Adjustment) {
		a.SetStepIncrement(step)
	}
}

// WithPageIncrement sets the page increment of the adjustment
func WithPageIncrement(page float64) AdjustmentOption {
	return func(a *Adjustment) {
		a.SetPageIncrement(page)
	}
}

// WithPageSize sets the page size of the adjustment
func WithPageSize(size float64) AdjustmentOption {
	return func(a *Adjustment) {
		a.SetPageSize(size)
	}
}

// GetValue gets the value of the adjustment
func (a *Adjustment) GetValue() float64 {
	return float64(C.gtk_adjustment_get_value(a.adjustment))
}

// SetValue sets the value of the adjustment
func (a *Adjustment) SetValue(value float64) {
	C.gtk_adjustment_set_value(a.adjustment, C.gdouble(value))
}

// GetLower gets the lower bound of the adjustment
func (a *Adjustment) GetLower() float64 {
	return float64(C.gtk_adjustment_get_lower(a.adjustment))
}

// SetLower sets the lower bound of the adjustment
func (a *Adjustment) SetLower(lower float64) {
	C.gtk_adjustment_set_lower(a.adjustment, C.gdouble(lower))
}

// GetUpper gets the upper bound of the adjustment
func (a *Adjustment) GetUpper() float64 {
	return float64(C.gtk_adjustment_get_upper(a.adjustment))
}

// SetUpper sets the upper bound of the adjustment
func (a *Adjustment) SetUpper(upper float64) {
	C.gtk_adjustment_set_upper(a.adjustment, C.gdouble(upper))
}

// GetStepIncrement gets the step increment of the adjustment
func (a *Adjustment) GetStepIncrement() float64 {
	return float64(C.gtk_adjustment_get_step_increment(a.adjustment))
}

// SetStepIncrement sets the step increment of the adjustment
func (a *Adjustment) SetStepIncrement(step float64) {
	C.gtk_adjustment_set_step_increment(a.adjustment, C.gdouble(step))
}

// GetPageIncrement gets the page increment of the adjustment
func (a *Adjustment) GetPageIncrement() float64 {
	return float64(C.gtk_adjustment_get_page_increment(a.adjustment))
}

// SetPageIncrement sets the page increment of the adjustment
func (a *Adjustment) SetPageIncrement(page float64) {
	C.gtk_adjustment_set_page_increment(a.adjustment, C.gdouble(page))
}

// GetPageSize gets the page size of the adjustment
func (a *Adjustment) GetPageSize() float64 {
	return float64(C.gtk_adjustment_get_page_size(a.adjustment))
}

// SetPageSize sets the page size of the adjustment
func (a *Adjustment) SetPageSize(size float64) {
	C.gtk_adjustment_set_page_size(a.adjustment, C.gdouble(size))
}

// ConnectValueChanged connects a callback to the value-changed signal
func (a *Adjustment) ConnectValueChanged(callback func()) {
	// Use the unified callback system
	Connect(a, SignalValueChanged, callback)
}

// DisconnectValueChanged disconnects the value-changed signal handler
func (a *Adjustment) DisconnectValueChanged() {
	// Use the unified callback system to find and disconnect all callbacks
	// for this object and signal
	adjustmentPtr := uintptr(unsafe.Pointer(a.adjustment))
	callbackIDs := getCallbackIDsForSignal(adjustmentPtr, SignalValueChanged)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// Free frees the adjustment
func (a *Adjustment) Free() {
	if a.adjustment != nil {
		// Disconnect all signal handlers
		DisconnectAll(a)
		
		C.g_object_unref(C.gpointer(unsafe.Pointer(a.adjustment)))
		a.adjustment = nil
	}
}