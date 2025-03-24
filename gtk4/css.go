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
//
// // Helper to mark CSS provider for optimization
// static void _go_css_provider_set_optimization(GtkCssProvider *provider, gboolean optimize) {
//     g_object_set_data(G_OBJECT(provider), "optimize-rendering", GINT_TO_POINTER(optimize ? 1 : 0));
// }
import "C"

import (
	"os"
	"runtime"
	"sync"
	"unsafe"
)

// Global CSS provider cache to avoid recreating providers
var (
	// Cache CSS providers by content
	cssProviderCache = make(map[string]*CSSProvider)
	cssProviderMutex sync.RWMutex

	// Track global providers for optimization
	globalProviders     = make([]*CSSProvider, 0, 5)
	globalProviderMutex sync.RWMutex

	// Lightweight CSS for resize operations
	resizeCSSProvider *CSSProvider
)

func init() {
	// Create a lightweight CSS provider for resize operations
	initResizeCSS := `
		/* Minimal CSS during resize - only essential rules */
		window, dialog { background-color: #f5f5f5; }
		button { padding: 2px; }
		entry { padding: 2px; }
		label { padding: 0; }
	`
	var err error
	resizeCSSProvider, err = loadCSS(initResizeCSS)
	if err != nil {
		// Fall back to empty provider if there's an error
		resizeCSSProvider = newCSSProvider()
	}
}

// CSSProvider represents a GTK CSS provider
type CSSProvider struct {
	provider *C.GtkCssProvider
	cssData  string // Store original CSS for cache lookups
}

// newCSSProvider creates a new GTK CSS provider
func newCSSProvider() *CSSProvider {
	provider := &CSSProvider{
		provider: C.gtk_css_provider_new(),
	}
	runtime.SetFinalizer(provider, (*CSSProvider).free)
	return provider
}

// loadFromData loads CSS data from a string
func (p *CSSProvider) loadFromData(cssData string) error {
	// Store the CSS data for cache lookups
	p.cssData = cssData

	cCssData := C.CString(cssData)
	defer C.free(unsafe.Pointer(cCssData))

	// Use the helper function that calls the GTK4 API
	C._go_css_provider_load_from_string(p.provider, cCssData)
	return nil
}

// loadFromFile loads CSS data from a file
func (p *CSSProvider) loadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return p.loadFromData(string(data))
}

// free frees the CSS provider
func (p *CSSProvider) free() {
	if p.provider != nil {
		// Remove from global providers list if present
		globalProviderMutex.Lock()
		for i, provider := range globalProviders {
			if provider == p {
				// Remove without preserving order
				globalProviders[i] = globalProviders[len(globalProviders)-1]
				globalProviders = globalProviders[:len(globalProviders)-1]
				break
			}
		}
		globalProviderMutex.Unlock()

		// Remove from cache if present
		if p.cssData != "" {
			cssProviderMutex.Lock()
			delete(cssProviderCache, p.cssData)
			cssProviderMutex.Unlock()
		}

		C.g_object_unref(C.gpointer(unsafe.Pointer(p.provider)))
		p.provider = nil
	}
}

// setOptimization enables or disables rendering optimization for this provider
func (p *CSSProvider) setOptimization(optimize bool) {
	var cOptimize C.gboolean
	if optimize {
		cOptimize = C.TRUE
	} else {
		cOptimize = C.FALSE
	}
	C._go_css_provider_set_optimization(p.provider, cOptimize)
}

// optimizeAllProviders enables optimization for all global CSS providers
func optimizeAllProviders() {
	globalProviderMutex.RLock()
	providers := make([]*CSSProvider, len(globalProviders))
	copy(providers, globalProviders)
	globalProviderMutex.RUnlock()

	for _, provider := range providers {
		provider.setOptimization(true)
	}
}

// resetAllProviders disables optimization for all global CSS providers
func resetAllProviders() {
	globalProviderMutex.RLock()
	providers := make([]*CSSProvider, len(globalProviders))
	copy(providers, globalProviders)
	globalProviderMutex.RUnlock()

	for _, provider := range providers {
		provider.setOptimization(false)
	}
}

// useResizeCSSProvider temporarily switches to a lightweight CSS provider during resize
func useResizeCSSProvider(display *C.GdkDisplay) *C.GtkCssProvider {
	if display == nil {
		display = C.gdk_display_get_default()
	}

	// Instead of trying to get the existing provider, we'll just
	// add our lightweight provider with a higher priority
	if resizeCSSProvider != nil {
		C.gtk_style_context_add_provider_for_display(display,
			(*C.GtkStyleProvider)(unsafe.Pointer(resizeCSSProvider.provider)),
			C.guint(priorityResize)) // Higher priority
	}

	// Return nil since we can't get the original provider directly
	return nil
}

// restoreOriginalCSSProvider restores the CSS after resize
func restoreOriginalCSSProvider(display *C.GdkDisplay, original *C.GtkCssProvider) {
	if display == nil {
		display = C.gdk_display_get_default()
	}

	// Just remove the resize provider
	if resizeCSSProvider != nil {
		C.gtk_style_context_remove_provider_for_display(display,
			(*C.GtkStyleProvider)(unsafe.Pointer(resizeCSSProvider.provider)))
	}

	// We don't need to restore the original provider since we never removed it,
	// we just temporarily added a higher-priority provider
}

// AddProviderForDisplay adds a CSS provider to the default display
func AddProviderForDisplay(provider *CSSProvider, priority uint) {
	display := C.gdk_display_get_default()
	C.gtk_style_context_add_provider_for_display(display,
		(*C.GtkStyleProvider)(unsafe.Pointer(provider.provider)),
		C.guint(priority))

	// Add to global providers list for optimization
	globalProviderMutex.Lock()
	globalProviders = append(globalProviders, provider)
	globalProviderMutex.Unlock()
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
	// priorityApplication is the priority for application-specific styles
	priorityApplication StylePriority = 600
	// priorityUser is the priority for user-specific styles
	priorityUser StylePriority = 800
	// priorityTheme is the priority for theme styles
	priorityTheme StylePriority = 400
	// prioritySetting is the priority for settings styles
	prioritySetting StylePriority = 500
	// priorityFallback is the priority for fallback styles
	priorityFallback StylePriority = 1
	// priorityResize is a higher priority used during resize operations
	priorityResize StylePriority = 900
)

// loadCSS is a convenience function to create a provider and load CSS from a string with caching
func loadCSS(cssData string) (*CSSProvider, error) {
	// Check cache first
	cssProviderMutex.RLock()
	provider, exists := cssProviderCache[cssData]
	cssProviderMutex.RUnlock()

	if exists {
		return provider, nil
	}

	// Create new provider if not in cache
	provider = newCSSProvider()
	err := provider.loadFromData(cssData)
	if err != nil {
		return nil, err
	}

	// Store in cache
	cssProviderMutex.Lock()
	cssProviderCache[cssData] = provider
	cssProviderMutex.Unlock()

	return provider, nil
}

// LoadCSS is a public convenience function to create a provider and load CSS from a string
func LoadCSS(cssData string) (*CSSProvider, error) {
	return loadCSS(cssData)
}

// LoadCSSFromFile is a convenience function to create a provider and load CSS from a file
func LoadCSSFromFile(filepath string) (*CSSProvider, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return loadCSS(string(data))
}
