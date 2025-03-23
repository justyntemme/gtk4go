// Package gtk4 provides label functionality for GTK4
// File: gtk4go/gtk4/label.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// Label represents a GTK label
type Label struct {
	widget *C.GtkWidget
}

// NewLabel creates a new GTK label with the given text
func NewLabel(text string) *Label {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	label := &Label{
		widget: C.gtk_label_new(cText),
	}
	runtime.SetFinalizer(label, (*Label).Destroy)
	return label
}

// SetText sets the label text
func (l *Label) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gtk_label_set_text((*C.GtkLabel)(unsafe.Pointer(l.widget)), cText)
}

// SetMarkup sets the label markup
func (l *Label) SetMarkup(markup string) {
	cMarkup := C.CString(markup)
	defer C.free(unsafe.Pointer(cMarkup))
	C.gtk_label_set_markup((*C.GtkLabel)(unsafe.Pointer(l.widget)), cMarkup)
}

// Destroy destroys the label
func (l *Label) Destroy() {
	// Instead of using gtk_widget_destroy directly, use a safer approach for GTK4
	if l.widget != nil {
		// Using gtk_widget_unparent as a safer alternative in GTK4
		C.gtk_widget_unparent(l.widget)
		l.widget = nil
	}
}

// Native returns the underlying GtkWidget pointer
func (l *Label) Native() uintptr {
	return uintptr(unsafe.Pointer(l.widget))
}

// GetWidget returns the underlying GtkWidget pointer
func (l *Label) GetWidget() *C.GtkWidget {
	return l.widget
}
