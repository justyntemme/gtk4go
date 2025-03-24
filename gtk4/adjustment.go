// Package gtk4 provides adjustment functionality for GTK4
// File: gtk4go/gtk4/adjustment.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback function for adjustment value changes
// extern void adjustmentValueChangedCallback(GtkAdjustment *adjustment, gpointer user_data);
//
// // Connect adjustment value-changed signal with callback
// static gulong connectAdjustmentValueChanged(GtkAdjustment *adjustment, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(adjustment), "value-changed", G_CALLBACK(adjustmentValueChangedCallback), user_data);
// }
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// AdjustmentValueChangedCallback represents a callback for adjustment value changed events
type AdjustmentValueChangedCallback func()

var (
	adjustmentCallbacks     = make(map[uintptr]AdjustmentValueChangedCallback)
	adjustmentCallbackMutex sync.RWMutex
)

//export adjustmentValueChangedCallback
func adjustmentValueChangedCallback(adjustment *C.GtkAdjustment, userData C.gpointer) {
	adjustmentCallbackMutex.RLock()
	defer adjustmentCallbackMutex.RUnlock()

	// Convert adjustment pointer to uintptr for lookup
	adjustmentPtr := uintptr(unsafe.Pointer(adjustment))

	// Find and call the callback
	if callback, ok := adjustmentCallbacks[adjustmentPtr]; ok {
		callback()
	}
}

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
func (a *Adjustment) ConnectValueChanged(callback AdjustmentValueChangedCallback) {
	adjustmentCallbackMutex.Lock()
	defer adjustmentCallbackMutex.Unlock()

	// Store callback in map
	adjustmentPtr := uintptr(unsafe.Pointer(a.adjustment))
	adjustmentCallbacks[adjustmentPtr] = callback

	// Connect signal
	C.connectAdjustmentValueChanged(a.adjustment, C.gpointer(unsafe.Pointer(a.adjustment)))
}

// DisconnectValueChanged disconnects the value-changed signal handler
func (a *Adjustment) DisconnectValueChanged() {
	adjustmentCallbackMutex.Lock()
	defer adjustmentCallbackMutex.Unlock()

	// Remove callback from map
	adjustmentPtr := uintptr(unsafe.Pointer(a.adjustment))
	delete(adjustmentCallbacks, adjustmentPtr)
}

// Free frees the adjustment
func (a *Adjustment) Free() {
	if a.adjustment != nil {
		adjustmentCallbackMutex.Lock()
		defer adjustmentCallbackMutex.Unlock()

		// Remove callback from map if exists
		adjustmentPtr := uintptr(unsafe.Pointer(a.adjustment))
		delete(adjustmentCallbacks, adjustmentPtr)

		C.g_object_unref(C.gpointer(unsafe.Pointer(a.adjustment)))
		a.adjustment = nil
	}
}
