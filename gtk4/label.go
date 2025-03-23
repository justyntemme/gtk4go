// Package gtk4 provides label functionality for GTK4
// File: gtk4go/gtk4/label.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// LabelOption is a function that configures a label
type LabelOption func(*Label)

// Label represents a GTK label
type Label struct {
	BaseWidget
}

// NewLabel creates a new GTK label with the given text
func NewLabel(text string, options ...LabelOption) *Label {
	var widget *C.GtkWidget

	WithCString(text, func(cText *C.char) {
		widget = C.gtk_label_new(cText)
	})

	label := &Label{
		BaseWidget: BaseWidget{
			widget: widget,
		},
	}

	// Apply options
	for _, option := range options {
		option(label)
	}

	SetupFinalization(label, label.Destroy)
	return label
}

// WithMarkup configures a label to use markup
func WithMarkup(markup string) LabelOption {
	return func(l *Label) {
		l.SetMarkup(markup)
	}
}

// WithSelectable makes the label selectable
func WithSelectable(selectable bool) LabelOption {
	return func(l *Label) {
		var cselectable C.gboolean
		if selectable {
			cselectable = C.TRUE
		} else {
			cselectable = C.FALSE
		}
		C.gtk_label_set_selectable((*C.GtkLabel)(unsafe.Pointer(l.widget)), cselectable)
	}
}

// SetText sets the label text
func (l *Label) SetText(text string) {
	WithCString(text, func(cText *C.char) {
		C.gtk_label_set_text((*C.GtkLabel)(unsafe.Pointer(l.widget)), cText)
	})
}

// SetMarkup sets the label markup
func (l *Label) SetMarkup(markup string) {
	WithCString(markup, func(cMarkup *C.char) {
		C.gtk_label_set_markup((*C.GtkLabel)(unsafe.Pointer(l.widget)), cMarkup)
	})
}

// GetText gets the label text
func (l *Label) GetText() string {
	cText := C.gtk_label_get_text((*C.GtkLabel)(unsafe.Pointer(l.widget)))
	return C.GoString(cText)
}
