// Package gtk4 provides cell renderer functionality for GTK4
// File: gtk4go/gtk4/cellRenderer.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Wrapper functions for cell renderer property setting
// static void set_cell_text(GtkCellRenderer *renderer, const char *text) {
//     g_object_set(renderer, "text", text, NULL);
// }
//
// static void set_cell_editable(GtkCellRenderer *renderer, gboolean editable) {
//     g_object_set(renderer, "editable", editable, NULL);
// }
//
// static void set_cell_active(GtkCellRenderer *renderer, gboolean active) {
//     g_object_set(renderer, "active", active, NULL);
// }
//
// static void set_cell_radio(GtkCellRenderer *renderer, gboolean radio) {
//     g_object_set(renderer, "radio", radio, NULL);
// }
//
// static void set_cell_value(GtkCellRenderer *renderer, int value) {
//     g_object_set(renderer, "value", value, NULL);
// }
//
// // Signal callback function for cell renderer toggled
// extern void cellRendererToggledCallback(GtkCellRendererToggle *renderer, gchar *path, gpointer user_data);
//
// // Connect cell renderer toggled signal with callback
// static gulong connectCellRendererToggled(GtkCellRendererToggle *renderer, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(renderer), "toggled", G_CALLBACK(cellRendererToggledCallback), user_data);
// }
//
// // Signal callback function for cell renderer edited
// extern void cellRendererEditedCallback(GtkCellRendererText *renderer, gchar *path, gchar *new_text, gpointer user_data);
//
// // Connect cell renderer edited signal with callback
// static gulong connectCellRendererEdited(GtkCellRendererText *renderer, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(renderer), "edited", G_CALLBACK(cellRendererEditedCallback), user_data);
// }
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// CellRendererToggledCallback represents a callback for cell renderer toggled events
type CellRendererToggledCallback func(path string)

// CellRendererEditedCallback represents a callback for cell renderer edited events
type CellRendererEditedCallback func(path string, newText string)

var (
	cellRendererCallbacks     = make(map[uintptr]interface{})
	cellRendererCallbackMutex sync.RWMutex
)

//export cellRendererToggledCallback
func cellRendererToggledCallback(renderer *C.GtkCellRendererToggle, pathPtr *C.gchar, userData C.gpointer) {
	cellRendererCallbackMutex.RLock()
	defer cellRendererCallbackMutex.RUnlock()

	// Convert renderer pointer to uintptr for lookup
	rendererPtr := uintptr(unsafe.Pointer(renderer))

	// Find and call the callback
	if callback, ok := cellRendererCallbacks[rendererPtr].(CellRendererToggledCallback); ok {
		// Convert C path to Go string
		path := C.GoString((*C.char)(unsafe.Pointer(pathPtr)))
		callback(path)
	}
}

//export cellRendererEditedCallback
func cellRendererEditedCallback(renderer *C.GtkCellRendererText, pathPtr *C.gchar, newTextPtr *C.gchar, userData C.gpointer) {
	cellRendererCallbackMutex.RLock()
	defer cellRendererCallbackMutex.RUnlock()

	// Convert renderer pointer to uintptr for lookup
	rendererPtr := uintptr(unsafe.Pointer(renderer))

	// Find and call the callback
	if callback, ok := cellRendererCallbacks[rendererPtr].(CellRendererEditedCallback); ok {
		// Convert C strings to Go strings
		path := C.GoString((*C.char)(unsafe.Pointer(pathPtr)))
		var newText string
		if newTextPtr != nil {
			newText = C.GoString((*C.char)(unsafe.Pointer(newTextPtr)))
		}
		callback(path, newText)
	}
}

// CellRenderer is an interface for all cell renderers
type CellRenderer interface {
	// GetCellRenderer returns the underlying GtkCellRenderer pointer
	GetCellRenderer() *C.GtkCellRenderer
}

// BaseCellRenderer provides common functionality for cell renderers
type BaseCellRenderer struct {
	renderer *C.GtkCellRenderer
}

// GetCellRenderer returns the underlying GtkCellRenderer pointer
func (r *BaseCellRenderer) GetCellRenderer() *C.GtkCellRenderer {
	return r.renderer
}

// CellRendererText represents a text cell renderer
type CellRendererText struct {
	BaseCellRenderer
}

// NewCellRendererText creates a new text cell renderer
func NewCellRendererText() *CellRendererText {
	text := &CellRendererText{
		BaseCellRenderer: BaseCellRenderer{
			renderer: C.gtk_cell_renderer_text_new(),
		},
	}
	runtime.SetFinalizer(text, (*CellRendererText).Free)
	return text
}

// SetText sets the text property
func (r *CellRendererText) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	// Use the C wrapper function instead of direct g_object_set call
	C.set_cell_text(r.renderer, cText)
}

// SetEditable sets whether the cell is editable
func (r *CellRendererText) SetEditable(editable bool) {
	var cEditable C.gboolean
	if editable {
		cEditable = C.TRUE
	} else {
		cEditable = C.FALSE
	}

	// Use the C wrapper function instead of direct g_object_set call
	C.set_cell_editable(r.renderer, cEditable)
}

// ConnectEdited connects a callback function to the "edited" signal
func (r *CellRendererText) ConnectEdited(callback CellRendererEditedCallback) {
	cellRendererCallbackMutex.Lock()
	defer cellRendererCallbackMutex.Unlock()

	// Store callback in map
	rendererPtr := uintptr(unsafe.Pointer(r.renderer))
	cellRendererCallbacks[rendererPtr] = callback

	// Connect signal
	C.connectCellRendererEdited((*C.GtkCellRendererText)(unsafe.Pointer(r.renderer)),
		C.gpointer(unsafe.Pointer(r.renderer)))
}

// DisconnectEdited disconnects the edited signal handler
func (r *CellRendererText) DisconnectEdited() {
	cellRendererCallbackMutex.Lock()
	defer cellRendererCallbackMutex.Unlock()

	// Remove callback from map
	rendererPtr := uintptr(unsafe.Pointer(r.renderer))
	delete(cellRendererCallbacks, rendererPtr)
}

// Free frees the renderer
func (r *CellRendererText) Free() {
	if r.renderer != nil {
		// Clean up callback
		rendererPtr := uintptr(unsafe.Pointer(r.renderer))
		cellRendererCallbackMutex.Lock()
		delete(cellRendererCallbacks, rendererPtr)
		cellRendererCallbackMutex.Unlock()

		C.g_object_unref(C.gpointer(unsafe.Pointer(r.renderer)))
		r.renderer = nil
	}
}

// CellRendererToggle represents a toggle cell renderer
type CellRendererToggle struct {
	BaseCellRenderer
}

// NewCellRendererToggle creates a new toggle cell renderer
func NewCellRendererToggle() *CellRendererToggle {
	toggle := &CellRendererToggle{
		BaseCellRenderer: BaseCellRenderer{
			renderer: C.gtk_cell_renderer_toggle_new(),
		},
	}
	runtime.SetFinalizer(toggle, (*CellRendererToggle).Free)
	return toggle
}

// SetActive sets the active property
func (r *CellRendererToggle) SetActive(active bool) {
	var cActive C.gboolean
	if active {
		cActive = C.TRUE
	} else {
		cActive = C.FALSE
	}

	// Use the C wrapper function instead of direct g_object_set call
	C.set_cell_active(r.renderer, cActive)
}

// SetRadio sets whether the toggle is displayed as a radio button
func (r *CellRendererToggle) SetRadio(radio bool) {
	var cRadio C.gboolean
	if radio {
		cRadio = C.TRUE
	} else {
		cRadio = C.FALSE
	}

	// Use the C wrapper function instead of direct g_object_set call
	C.set_cell_radio(r.renderer, cRadio)
}

// ConnectToggled connects a callback to a toggle cell renderer's "toggled" signal
func (r *CellRendererToggle) ConnectToggled(callback CellRendererToggledCallback) {
	cellRendererCallbackMutex.Lock()
	defer cellRendererCallbackMutex.Unlock()

	// Store callback in map
	rendererPtr := uintptr(unsafe.Pointer(r.renderer))
	cellRendererCallbacks[rendererPtr] = callback

	// Connect signal
	C.connectCellRendererToggled((*C.GtkCellRendererToggle)(unsafe.Pointer(r.renderer)),
		C.gpointer(unsafe.Pointer(r.renderer)))
}

// DisconnectToggled disconnects the toggled signal handler
func (r *CellRendererToggle) DisconnectToggled() {
	cellRendererCallbackMutex.Lock()
	defer cellRendererCallbackMutex.Unlock()

	// Remove callback from map
	rendererPtr := uintptr(unsafe.Pointer(r.renderer))
	delete(cellRendererCallbacks, rendererPtr)
}

// Free frees the renderer
func (r *CellRendererToggle) Free() {
	if r.renderer != nil {
		// Clean up callback
		rendererPtr := uintptr(unsafe.Pointer(r.renderer))
		cellRendererCallbackMutex.Lock()
		delete(cellRendererCallbacks, rendererPtr)
		cellRendererCallbackMutex.Unlock()

		C.g_object_unref(C.gpointer(unsafe.Pointer(r.renderer)))
		r.renderer = nil
	}
}

// CellRendererPixbuf represents a pixbuf cell renderer
type CellRendererPixbuf struct {
	BaseCellRenderer
}

// NewCellRendererPixbuf creates a new pixbuf cell renderer
func NewCellRendererPixbuf() *CellRendererPixbuf {
	pixbuf := &CellRendererPixbuf{
		BaseCellRenderer: BaseCellRenderer{
			renderer: C.gtk_cell_renderer_pixbuf_new(),
		},
	}
	runtime.SetFinalizer(pixbuf, (*CellRendererPixbuf).Free)
	return pixbuf
}

// Free frees the renderer
func (r *CellRendererPixbuf) Free() {
	if r.renderer != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(r.renderer)))
		r.renderer = nil
	}
}

// CellRendererProgress represents a progress cell renderer
type CellRendererProgress struct {
	BaseCellRenderer
}

// NewCellRendererProgress creates a new progress cell renderer
func NewCellRendererProgress() *CellRendererProgress {
	progress := &CellRendererProgress{
		BaseCellRenderer: BaseCellRenderer{
			renderer: C.gtk_cell_renderer_progress_new(),
		},
	}
	runtime.SetFinalizer(progress, (*CellRendererProgress).Free)
	return progress
}

// SetValue sets the value property
func (r *CellRendererProgress) SetValue(value int) {
	// Use the C wrapper function instead of direct g_object_set call
	C.set_cell_value(r.renderer, C.int(value))
}

// Free frees the renderer
func (r *CellRendererProgress) Free() {
	if r.renderer != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(r.renderer)))
		r.renderer = nil
	}
}

