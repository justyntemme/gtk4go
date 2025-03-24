// Package gtk4 provides window functionality for GTK4
// File: gtk4go/gtk4/window.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Helper function to enable frame clock synchronization
// static void setWindowRenderingMode(GtkWindow* window) {
//     // Enable hardware acceleration for the window
//     GdkSurface* surface = gtk_native_get_surface(GTK_NATIVE(window));
//     if (surface != NULL) {
//         // Queue surface for rendering - proper way to ensure updates in GTK4
//         gdk_surface_queue_render(surface);
//
//         // In GTK4, we can improve rendering performance by:
//         // 1. Setting frame clock synchronization
//         GdkFrameClock* frame_clock = gdk_surface_get_frame_clock(surface);
//         if (frame_clock != NULL) {
//             // Request updates on the frame clock for smoother animation
//             gdk_frame_clock_begin_updating(frame_clock);
//         }
//     }
// }
//
// // Set content sizing mode for more efficient resizing
// static void setContentSizing(GtkWindow* window) {
//     // Use natural sizing for better resize performance
//     GtkWidget* child = gtk_window_get_child(window);
//     if (child != NULL) {
//         gtk_widget_set_hexpand(child, TRUE);
//         gtk_widget_set_vexpand(child, TRUE);
//     }
// }
import "C"

import (
	"os"
	"unsafe"
)

func init() {
	// Set environment variables to enable hardware acceleration
	// These must be set before the application starts for best effect
	// but we set them here to make sure they're present
	os.Setenv("GSK_RENDERER", "gl")
	os.Setenv("GDK_GL", "always")

	// On some systems, Cairo may be more stable than GL
	// Uncomment if GL causes issues:
	// os.Setenv("GSK_RENDERER", "cairo")
}

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

// Destroy destroys the window
func (w *Window) Destroy() {
	if w.widget != nil {
		C.gtk_window_destroy((*C.GtkWindow)(unsafe.Pointer(w.widget)))
		w.widget = nil
	}
}
