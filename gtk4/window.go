// Package gtk4 provides window functionality for GTK4
// File: gtk4go/gtk4/window.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Helper function to enable frame clock synchronization with platform-specific optimizations
// static void setWindowRenderingMode(GtkWindow* window) {
//     // Get surface for the window
//     GdkSurface* surface = gtk_native_get_surface(GTK_NATIVE(window));
//     if (surface == NULL) {
//         return;
//     }
//
//     // Queue surface for rendering - proper way to ensure updates in GTK4
//     gdk_surface_queue_render(surface);
//
// #ifdef __APPLE__
//     // macOS-specific rendering optimizations
//     // Use Cairo renderer which is more stable on macOS
//     // Avoid excessive animations and transitions
//     g_object_set(gtk_settings_get_default(), 
//                 "gtk-enable-animations", FALSE, 
//                 NULL);
// #else
//     // Linux-specific rendering optimizations
//     // Set frame clock synchronization for smoother animations on Linux
//     GdkFrameClock* frame_clock = gdk_surface_get_frame_clock(surface);
//     if (frame_clock != NULL) {
//         gdk_frame_clock_begin_updating(frame_clock);
//     }
// #endif
// }
//
// // Set content sizing mode for more efficient resizing with platform specifics
// static void setContentSizing(GtkWindow* window) {
//     // Use natural sizing for better resize performance
//     GtkWidget* child = gtk_window_get_child(window);
//     if (child != NULL) {
//         gtk_widget_set_hexpand(child, TRUE);
//         gtk_widget_set_vexpand(child, TRUE);
//
// #ifdef __APPLE__
//         // macOS needs specific margin handling for window decorations
//         gtk_widget_set_margin_start(child, 0);
//         gtk_widget_set_margin_end(child, 0);
//         gtk_widget_set_margin_top(child, 0);
//         gtk_widget_set_margin_bottom(child, 0);
// #endif
//     }
// }
//
// // Set a widget as the window titlebar
// static void setWindowTitlebar(GtkWindow *window, GtkWidget *titlebar) {
//     gtk_window_set_titlebar(window, titlebar);
// }
import "C"

import (
	"unsafe"
)

// WindowOption is a function that configures a window
type WindowOption func(*Window)

// Window represents a GTK window
type Window struct {
	BaseWidget
	isAcceleratedRendering bool
}

// NewWindow creates a new GTK window with the given title
func NewWindow(title string, options ...WindowOption) *Window {
	window := &Window{
		BaseWidget: BaseWidget{
			widget: C.gtk_window_new(),
		},
		isAcceleratedRendering: true,
	}

	WithCString(title, func(cTitle *C.char) {
		C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(window.widget)), cTitle)
	})

	// Apply default size
	C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(window.widget)), 600, 400)

	// Enable hardware acceleration and optimized rendering mode
	window.EnableAcceleratedRendering()

	// Apply options
	for _, option := range options {
		option(window)
	}

	SetupFinalization(window, window.Destroy)
	return window
}

// EnableAcceleratedRendering enables hardware-accelerated rendering and optimizes
// the window for better resize performance
func (w *Window) EnableAcceleratedRendering() {
	if !w.isAcceleratedRendering {
		return
	}

	// Use the C helper function to set up optimal rendering
	C.setWindowRenderingMode((*C.GtkWindow)(unsafe.Pointer(w.widget)))

	// Ensure widget is realized first
	C.gtk_widget_realize(w.widget)
}

// DisableAcceleratedRendering disables hardware acceleration
func (w *Window) DisableAcceleratedRendering() {
	w.isAcceleratedRendering = false
	// Additional code to disable acceleration if needed
}

// OptimizeForResizing applies optimizations specifically for window resizing
func (w *Window) OptimizeForResizing() {
	// Tell GTK to optimize layout calculations during resize
	C.gtk_window_set_resizable((*C.GtkWindow)(unsafe.Pointer(w.widget)), C.TRUE)

	// Configure child widget sizing for better performance
	C.setContentSizing((*C.GtkWindow)(unsafe.Pointer(w.widget)))
}

// SetTitlebar sets a widget as the window's titlebar
func (w *Window) SetTitlebar(titlebar Widget) {
    C.setWindowTitlebar((*C.GtkWindow)(unsafe.Pointer(w.widget)), titlebar.GetWidget())
}

// WithTitlebar sets a widget as the window's titlebar at creation time
func WithTitlebar(titlebar Widget) WindowOption {
    return func(w *Window) {
        w.SetTitlebar(titlebar)
    }
}

// WithDefaultSize sets the default window size
func WithDefaultSize(width, height int) WindowOption {
	return func(w *Window) {
		C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(w.widget)), C.int(width), C.int(height))
	}
}

// WithTransientFor sets the parent window
func WithTransientFor(parent *Window) WindowOption {
	return func(w *Window) {
		if parent != nil {
			C.gtk_window_set_transient_for(
				(*C.GtkWindow)(unsafe.Pointer(w.widget)),
				(*C.GtkWindow)(unsafe.Pointer(parent.widget)),
			)
		}
	}
}

// WithModal sets whether the window is modal
func WithModal(modal bool) WindowOption {
	return func(w *Window) {
		var cmodal C.gboolean
		if modal {
			cmodal = C.TRUE
		} else {
			cmodal = C.FALSE
		}
		C.gtk_window_set_modal((*C.GtkWindow)(unsafe.Pointer(w.widget)), cmodal)
	}
}

// WithAcceleratedRendering sets whether to use hardware acceleration
func WithAcceleratedRendering(enabled bool) WindowOption {
	return func(w *Window) {
		w.isAcceleratedRendering = enabled
		if enabled {
			w.EnableAcceleratedRendering()
		} else {
			w.DisableAcceleratedRendering()
		}
	}
}

// SetTitle sets the window title
func (w *Window) SetTitle(title string) {
	WithCString(title, func(cTitle *C.char) {
		C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(w.widget)), cTitle)
	})
}

// SetDefaultSize sets the default window size
func (w *Window) SetDefaultSize(width, height int) {
	C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(w.widget)), C.int(width), C.int(height))
}

// SetChild sets the child widget for the window
func (w *Window) SetChild(child Widget) {
	C.gtk_window_set_child((*C.GtkWindow)(unsafe.Pointer(w.widget)), child.GetWidget())

	// Re-apply optimization for resizing when child changes
	C.setContentSizing((*C.GtkWindow)(unsafe.Pointer(w.widget)))
}

// Show makes the window visible
func (w *Window) Show() {
	// Ensure hardware acceleration is set up before showing
	w.EnableAcceleratedRendering()
	C.gtk_widget_set_visible(w.widget, C.TRUE)
}

// Present presents the window to the user (preferred in GTK4)
func (w *Window) Present() {
	// Ensure hardware acceleration is set up before presenting
	w.EnableAcceleratedRendering()
	C.gtk_window_present((*C.GtkWindow)(unsafe.Pointer(w.widget)))
}

// SetVisible sets the visibility of the window
func (w *Window) SetVisible(visible bool) {
	// Ensure hardware acceleration is set up if becoming visible
	if visible {
		w.EnableAcceleratedRendering()
	}

	var cvisible C.gboolean
	if visible {
		cvisible = C.TRUE
	} else {
		cvisible = C.FALSE
	}
	C.gtk_widget_set_visible(w.widget, cvisible)
}

// ConnectCloseRequest connects a callback function to the window's "close-request" signal
// The callback should return true to stop the default handling of the signal (prevent closing),
// or false to allow the default handling (allow closing)
func (w *Window) ConnectCloseRequest(callback func() bool) uint64 {
	return Connect(w, SignalCloseRequest, callback)
}

// DisconnectCloseRequest disconnects the close-request signal handler
func (w *Window) DisconnectCloseRequest() {
	// Get all callbacks for this window
	windowPtr := uintptr(unsafe.Pointer(w.widget))
	callbackIDs := getCallbackIDsForSignal(windowPtr, SignalCloseRequest)
	
	// Disconnect each callback
	for _, id := range callbackIDs {
		Disconnect(id)
	}
}

// Destroy destroys the window and cleans up resources
func (w *Window) Destroy() {
	if w.widget != nil {
		// Disconnect all signals for this window
		DisconnectAll(w)
		
		// Clean up window resize detection if it was set up
		delete(windowResizeStates, uintptr(unsafe.Pointer(w.widget)))

		// Destroy the window
		C.gtk_window_destroy((*C.GtkWindow)(unsafe.Pointer(w.widget)))
		w.widget = nil
	}
}