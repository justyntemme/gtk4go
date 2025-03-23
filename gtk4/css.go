// Package gtk4 provides CSS styling functionality for GTK4
// File: gtk4go/gtk4/css.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Helper to convert Go string to CSS string
// static void _go_css_provider_load_from_string(GtkCssProvider *provider, const char *css_string) {
//     gtk_css_provider_load_from_string(provider, css_string);
// }
import "C"

import (
	"os"
	"runtime"
	"unsafe"
)

// CSSProvider represents a GTK CSS provider
type CSSProvider struct {
	provider *C.GtkCssProvider
}

// NewCSSProvider creates a new GTK CSS provider
func NewCSSProvider() *CSSProvider {
	provider := &CSSProvider{
		provider: C.gtk_css_provider_new(),
	}
	runtime.SetFinalizer(provider, (*CSSProvider).Free)
	return provider
}

// LoadFromData loads CSS data from a string
func (p *CSSProvider) LoadFromData(cssData string) error {
	cCssData := C.CString(cssData)
	defer C.free(unsafe.Pointer(cCssData))

	// Use the helper function that calls the GTK4 API
	C._go_css_provider_load_from_string(p.provider, cCssData)
	return nil
}

// LoadFromFile loads CSS data from a file
func (p *CSSProvider) LoadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return p.LoadFromData(string(data))
}

// Free frees the CSS provider
func (p *CSSProvider) Free() {
	C.g_object_unref(C.gpointer(unsafe.Pointer(p.provider)))
	p.provider = nil
}

// AddProviderForDisplay adds a CSS provider to the default display
func AddProviderForDisplay(provider *CSSProvider, priority uint) {
	display := C.gdk_display_get_default()
	C.gtk_style_context_add_provider_for_display(display,
		(*C.GtkStyleProvider)(unsafe.Pointer(provider.provider)),
		C.guint(priority))
}

// Widget CSS class methods - using modern GTK4 API

// AddStyleClass adds a CSS class to a widget
func AddStyleClass(widget interface{}, className string) {
	if w, ok := widget.(interface{ GetWidget() *C.GtkWidget }); ok {
		cClassName := C.CString(className)
		defer C.free(unsafe.Pointer(cClassName))
		C.gtk_widget_add_css_class(w.GetWidget(), cClassName)
	}
}

// RemoveStyleClass removes a CSS class from a widget
func RemoveStyleClass(widget interface{}, className string) {
	if w, ok := widget.(interface{ GetWidget() *C.GtkWidget }); ok {
		cClassName := C.CString(className)
		defer C.free(unsafe.Pointer(cClassName))
		C.gtk_widget_remove_css_class(w.GetWidget(), cClassName)
	}
}

// HasStyleClass checks if a widget has a CSS class
func HasStyleClass(widget interface{}, className string) bool {
	if w, ok := widget.(interface{ GetWidget() *C.GtkWidget }); ok {
		cClassName := C.CString(className)
		defer C.free(unsafe.Pointer(cClassName))
		return C.gtk_widget_has_css_class(w.GetWidget(), cClassName) == 1
	}
	return false
}

// StylePriority defines the priority levels for CSS providers
type StylePriority uint

const (
	// PriorityApplication is the priority for application-specific styles
	PriorityApplication StylePriority = 600
	// PriorityUser is the priority for user-specific styles
	PriorityUser StylePriority = 800
	// PriorityTheme is the priority for theme styles
	PriorityTheme StylePriority = 400
	// PrioritySetting is the priority for settings styles
	PrioritySetting StylePriority = 500
	// PriorityFallback is the priority for fallback styles
	PriorityFallback StylePriority = 1
)

// LoadCSS is a convenience function to create a provider and load CSS from a string
func LoadCSS(cssData string) (*CSSProvider, error) {
	provider := NewCSSProvider()
	err := provider.LoadFromData(cssData)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// LoadCSSFromFile is a convenience function to create a provider and load CSS from a file
func LoadCSSFromFile(filepath string) (*CSSProvider, error) {
	provider := NewCSSProvider()
	err := provider.LoadFromFile(filepath)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
