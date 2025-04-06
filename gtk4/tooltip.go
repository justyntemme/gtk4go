// Package gtk4 provides tooltip functionality for GTK4
// File: gtk4go/gtk4/tooltip.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Direct tooltip functions - high-level API used by all widgets
// static void widget_set_tooltip_text(GtkWidget *widget, const char *text) {
//     gtk_widget_set_tooltip_text(widget, text);
// }
//
// static void widget_set_tooltip_markup(GtkWidget *widget, const char *markup) {
//     gtk_widget_set_tooltip_markup(widget, markup);
// }
//
// // Set has tooltip
// static void setHasTooltip(GtkWidget *widget, gboolean has_tooltip) {
//     gtk_widget_set_has_tooltip(widget, has_tooltip);
// }
//
// // Get tooltip text
// static char* getTooltipText(GtkWidget *widget) {
//     const char* text = gtk_widget_get_tooltip_text(widget);
//     if (text != NULL) {
//         return g_strdup(text);
//     }
//     return NULL;
// }
//
// // Get tooltip markup
// static char* getTooltipMarkup(GtkWidget *widget) {
//     const char* markup = gtk_widget_get_tooltip_markup(widget);
//     if (markup != NULL) {
//         return g_strdup(markup);
//     }
//     return NULL;
// }
//
// // Get has tooltip
// static gboolean getHasTooltip(GtkWidget *widget) {
//     gboolean has_tooltip;
//     g_object_get(G_OBJECT(widget), "has-tooltip", &has_tooltip, NULL);
//     return has_tooltip;
// }
//
// // Set tooltip delay
// static void setTooltipDelay(GtkWidget *widget, guint delay) {
//     g_object_set(G_OBJECT(widget), "tooltip-delay", delay, NULL);
// }
//
// // Get tooltip delay
// static guint getTooltipDelay(GtkWidget *widget) {
//     guint delay;
//     g_object_get(G_OBJECT(widget), "tooltip-delay", &delay, NULL);
//     return delay;
// }
//
// // Set tooltip icon
// static void setTooltipIcon(GtkTooltip *tooltip, GdkPaintable *paintable) {
//     gtk_tooltip_set_icon(tooltip, (GdkPaintable*)paintable);
// }
import "C"

import (
	"unsafe"
)

// Signal type for tooltip query - comment out to avoid redeclaration
// This is already defined in callbacks.go
// const (
//	SignalQueryTooltip SignalType = "query-tooltip"
// )

// TooltipQueryCallback represents a callback for tooltip query events
type TooltipQueryCallback func(x, y int, keyboardMode bool, tooltip *Tooltip) bool

// Tooltip represents a GTK tooltip
type Tooltip struct {
	tooltip *C.GtkTooltip
}

// SetText sets the tooltip text
func (t *Tooltip) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gtk_tooltip_set_text(t.tooltip, cText)
}

// SetMarkup sets the tooltip markup
func (t *Tooltip) SetMarkup(markup string) {
	cMarkup := C.CString(markup)
	defer C.free(unsafe.Pointer(cMarkup))
	C.gtk_tooltip_set_markup(t.tooltip, cMarkup)
}

// SetIcon sets the tooltip icon
func (t *Tooltip) SetIcon(iconName string) {
	cIconName := C.CString(iconName)
	defer C.free(unsafe.Pointer(cIconName))
	
	// The gtk_tooltip_set_icon_from_icon_name function is safer to use
	C.gtk_tooltip_set_icon_from_icon_name(t.tooltip, cIconName)
}

// SetTipArea sets the area of the widget associated with this tooltip
func (t *Tooltip) SetTipArea(x, y, width, height int) {
	var rect C.GdkRectangle
	rect.x = C.int(x)
	rect.y = C.int(y)
	rect.width = C.int(width)
	rect.height = C.int(height)
	
	C.gtk_tooltip_set_tip_area(t.tooltip, &rect)
}

// SetTooltipText sets a simple text tooltip on a widget
func (w *BaseWidget) SetTooltipText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.widget_set_tooltip_text(w.widget, cText)
}

// SetTooltipMarkup sets a tooltip with markup on a widget
func (w *BaseWidget) SetTooltipMarkup(markup string) {
	cMarkup := C.CString(markup)
	defer C.free(unsafe.Pointer(cMarkup))
	C.widget_set_tooltip_markup(w.widget, cMarkup)
}

// SetHasTooltip sets whether a widget has a tooltip
func (w *BaseWidget) SetHasTooltip(hasTooltip bool) {
	var cHasTooltip C.gboolean
	if hasTooltip {
		cHasTooltip = C.TRUE
	} else {
		cHasTooltip = C.FALSE
	}
	C.setHasTooltip(w.widget, cHasTooltip)
}

// GetTooltipText gets the tooltip text
func (w *BaseWidget) GetTooltipText() string {
	cText := C.getTooltipText(w.widget)
	if cText == nil {
		return ""
	}
	text := C.GoString(cText)
	C.free(unsafe.Pointer(cText))
	return text
}

// GetTooltipMarkup gets the tooltip markup
func (w *BaseWidget) GetTooltipMarkup() string {
	cMarkup := C.getTooltipMarkup(w.widget)
	if cMarkup == nil {
		return ""
	}
	markup := C.GoString(cMarkup)
	C.free(unsafe.Pointer(cMarkup))
	return markup
}

// GetHasTooltip gets whether a widget has a tooltip
func (w *BaseWidget) GetHasTooltip() bool {
	return C.getHasTooltip(w.widget) == C.TRUE
}

// SetTooltipDelay sets the delay before showing the tooltip in milliseconds
func (w *BaseWidget) SetTooltipDelay(delay uint) {
	C.setTooltipDelay(w.widget, C.guint(delay))
}

// GetTooltipDelay gets the delay before showing the tooltip in milliseconds
func (w *BaseWidget) GetTooltipDelay() uint {
	return uint(C.getTooltipDelay(w.widget))
}

// ConnectQueryTooltip connects a callback for the query-tooltip signal
func (w *BaseWidget) ConnectQueryTooltip(callback TooltipQueryCallback) uint64 {
	if callback == nil {
		return 0
	}
	
	// Set has-tooltip property to true
	w.SetHasTooltip(true)
	
	// Wrap the callback to match the expected format for the UCS
	wrappedCallback := func(x int, y int, keyboardMode bool, tooltipPtr uintptr) bool {
		tooltip := &Tooltip{tooltip: (*C.GtkTooltip)(unsafe.Pointer(tooltipPtr))}
		return callback(x, y, keyboardMode, tooltip)
	}
	
	// Register this as a callback with return value
	return Connect(w, SignalQueryTooltip, wrappedCallback)
}

// DisconnectQueryTooltip disconnects the query-tooltip callback
func (w *BaseWidget) DisconnectQueryTooltip() {
	// Get all callbacks for this widget and signal
	widgetPtr := uintptr(unsafe.Pointer(w.widget))
	callbackIDs := getCallbackIDsForSignal(widgetPtr, SignalQueryTooltip)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// TooltipOption configures tooltip behavior for widgets
type TooltipOption func(widget *BaseWidget)

// WithTooltipText sets a simple text tooltip
func WithTooltipText(text string) TooltipOption {
	return func(w *BaseWidget) {
		w.SetTooltipText(text)
	}
}

// WithTooltipMarkup sets a tooltip with markup
func WithTooltipMarkup(markup string) TooltipOption {
	return func(w *BaseWidget) {
		w.SetTooltipMarkup(markup)
	}
}

// WithHasTooltip sets whether a widget has a tooltip
func WithHasTooltip(hasTooltip bool) TooltipOption {
	return func(w *BaseWidget) {
		w.SetHasTooltip(hasTooltip)
	}
}

// WithTooltipDelay sets the delay before showing the tooltip
func WithTooltipDelay(delay uint) TooltipOption {
	return func(w *BaseWidget) {
		w.SetTooltipDelay(delay)
	}
}

// Note: tooltipQueryCallback implementation moved to callbacks.go to avoid duplication

// SetTooltipText is a convenient global function for setting tooltip text on any widget
func SetTooltipText(widget Widget, text string) {
	if w, ok := widget.(interface{ SetTooltipText(string) }); ok {
		w.SetTooltipText(text)
	}
}

// SetTooltipMarkup is a convenient global function for setting tooltip markup on any widget
func SetTooltipMarkup(widget Widget, markup string) {
	if w, ok := widget.(interface{ SetTooltipMarkup(string) }); ok {
		w.SetTooltipMarkup(markup)
	}
}

// Initialize a debug component for tooltips
func init() {
	// Add a new debug component for tooltips
	EnableDebugComponent(DebugComponentTooltip)
}