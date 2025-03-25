// File: gtk4go/gtk4/windowActions.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Add action to window
// static void add_action_to_window(GtkWindow* window, GAction* action) {
//     g_action_map_add_action(G_ACTION_MAP(window), action);
// }
//
// // Remove action from window
// static void remove_action_from_window(GtkWindow* window, const char* name) {
//     g_action_map_remove_action(G_ACTION_MAP(window), name);
// }
import "C"

import (
	"unsafe"
)

// WindowActionGroup implements the ActionGroup interface for Window
type WindowActionGroup struct {
	win *Window
}

// AddAction adds an action to the window
func (w *WindowActionGroup) AddAction(action *Action) {
	C.add_action_to_window((*C.GtkWindow)(unsafe.Pointer(w.win.widget)), action.GetNative())
}

// RemoveAction removes an action from the window
func (w *WindowActionGroup) RemoveAction(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.remove_action_from_window((*C.GtkWindow)(unsafe.Pointer(w.win.widget)), cName)
}

// GetActionGroup returns the window's action group
func (w *Window) GetActionGroup() ActionGroup {
	return &WindowActionGroup{win: w}
}

