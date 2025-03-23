// Package gtk4 provides widget extension functionality for GTK4
// File: gtk4go/gtk4/widgetCCCExtensions.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Helper function to get gboolean as int
// static int _go_widget_has_css_class(GtkWidget *widget, const char *class_name) {
//     return (int)gtk_widget_has_css_class(widget, class_name);
// }
import "C"

import (
	"unsafe"
)

// AddCssClass adds a CSS class to a widget
func (w *Window) AddCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_add_css_class(w.widget, cClassName)
}

// RemoveCssClass removes a CSS class from a widget
func (w *Window) RemoveCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_remove_css_class(w.widget, cClassName)
}

// HasCssClass checks if a widget has a CSS class
func (w *Window) HasCssClass(className string) bool {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	return C._go_widget_has_css_class(w.widget, cClassName) != 0
}

// AddCssClass adds a CSS class to a widget
func (b *Box) AddCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_add_css_class(b.widget, cClassName)
}

// RemoveCssClass removes a CSS class from a widget
func (b *Box) RemoveCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_remove_css_class(b.widget, cClassName)
}

// HasCssClass checks if a widget has a CSS class
func (b *Box) HasCssClass(className string) bool {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	return C._go_widget_has_css_class(b.widget, cClassName) != 0
}

// AddCssClass adds a CSS class to a widget
func (b *Button) AddCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_add_css_class(b.widget, cClassName)
}

// RemoveCssClass removes a CSS class from a widget
func (b *Button) RemoveCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_remove_css_class(b.widget, cClassName)
}

// HasCssClass checks if a widget has a CSS class
func (b *Button) HasCssClass(className string) bool {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	return C._go_widget_has_css_class(b.widget, cClassName) != 0
}

// AddCssClass adds a CSS class to a widget
func (l *Label) AddCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_add_css_class(l.widget, cClassName)
}

// RemoveCssClass removes a CSS class from a widget
func (l *Label) RemoveCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_remove_css_class(l.widget, cClassName)
}

// HasCssClass checks if a widget has a CSS class
func (l *Label) HasCssClass(className string) bool {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	return C._go_widget_has_css_class(l.widget, cClassName) != 0
}

// AddCssClass adds a CSS class to a widget
func (e *Entry) AddCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_add_css_class(e.widget, cClassName)
}

// RemoveCssClass removes a CSS class from a widget
func (e *Entry) RemoveCssClass(className string) {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	C.gtk_widget_remove_css_class(e.widget, cClassName)
}

// HasCssClass checks if a widget has a CSS class
func (e *Entry) HasCssClass(className string) bool {
	cClassName := C.CString(className)
	defer C.free(unsafe.Pointer(cClassName))
	return C._go_widget_has_css_class(e.widget, cClassName) != 0
}
