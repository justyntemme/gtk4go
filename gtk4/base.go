// Package gtk4 provides base widget functionality for GTK4
// File: gtk4go/gtk4/base.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// Widget defines the common interface for GTK widgets
type Widget interface {
	// GetWidget returns the underlying GtkWidget pointer
	GetWidget() *C.GtkWidget

	// Native returns the underlying pointer as uintptr
	Native() uintptr

	// Destroy releases the widget resources
	Destroy()

	// AddCssClass adds a CSS class to the widget
	AddCssClass(className string)

	// RemoveCssClass removes a CSS class from the widget
	RemoveCssClass(className string)

	// HasCssClass checks if the widget has a CSS class
	HasCssClass(className string) bool
}

// BaseWidget provides common functionality for GTK widgets
type BaseWidget struct {
	widget *C.GtkWidget
}

// GetWidget returns the underlying GtkWidget pointer
func (w *BaseWidget) GetWidget() *C.GtkWidget {
	return w.widget
}

// Native returns the underlying GtkWidget pointer as uintptr
func (w *BaseWidget) Native() uintptr {
	return uintptr(unsafe.Pointer(w.widget))
}

// Destroy destroys the widget
func (w *BaseWidget) Destroy() {
	if w.widget != nil {
		C.gtk_widget_unparent(w.widget)
		w.widget = nil
	}
}

// AddCssClass adds a CSS class to the widget
func (w *BaseWidget) AddCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_add_css_class(w.widget, cClassName)
}

// RemoveCssClass removes a CSS class from the widget
func (w *BaseWidget) RemoveCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_remove_css_class(w.widget, cClassName)
}

// HasCssClass checks if the widget has a CSS class
func (w *BaseWidget) HasCssClass(className string) bool {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	return C.gtk_widget_has_css_class(w.widget, cClassName) == 1
}

// WithCString executes a function with a C string that is automatically freed
func WithCString(s string, fn func(*C.char)) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	fn(cs)
}

// CastToGObject casts a widget pointer to a GObject pointer
func CastToGObject(widget *C.GtkWidget) *C.GObject {
	return (*C.GObject)(unsafe.Pointer(widget))
}

// SetupFinalization sets up proper finalization for a widget
func SetupFinalization(widget Widget, destroyFunc func()) {
	runtime.SetFinalizer(widget, func(w Widget) {
		destroyFunc()
	})
}

// SignalHandler manages signals and callbacks for GTK widgets
type SignalHandler struct {
	callbacks     map[uintptr]map[string]interface{}
	callbackMutex sync.Mutex
}

// NewSignalHandler creates a new signal handler
func NewSignalHandler() *SignalHandler {
	return &SignalHandler{
		callbacks: make(map[uintptr]map[string]interface{}),
	}
}

// Connect connects a callback to a signal
func (s *SignalHandler) Connect(widget uintptr, signal string, callback interface{}) {
	s.callbackMutex.Lock()
	defer s.callbackMutex.Unlock()

	if _, ok := s.callbacks[widget]; !ok {
		s.callbacks[widget] = make(map[string]interface{})
	}

	s.callbacks[widget][signal] = callback
}

// Disconnect disconnects all callbacks for a widget
func (s *SignalHandler) Disconnect(widget uintptr) {
	s.callbackMutex.Lock()
	defer s.callbackMutex.Unlock()

	delete(s.callbacks, widget)
}

// Get retrieves a callback for a widget and signal
func (s *SignalHandler) Get(widget uintptr, signal string) interface{} {
	s.callbackMutex.Lock()
	defer s.callbackMutex.Unlock()

	if callbackMap, ok := s.callbacks[widget]; ok {
		if callback, ok := callbackMap[signal]; ok {
			return callback
		}
	}

	return nil
}

// GTKError represents an error in GTK operations
type GTKError struct {
	Op  string
	Err error
}

// Error implements the error interface
func (e *GTKError) Error() string {
	if e.Err != nil {
		return "gtk4go: " + e.Op + ": " + e.Err.Error()
	}
	return "gtk4go: " + e.Op
}
