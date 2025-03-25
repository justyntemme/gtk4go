// Package gtk4 provides modern action-based menu functionality for GTK4
// File: gtk4go/gtk4/action.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Action callback function
// extern void actionActivateCallback(GSimpleAction *action, GVariant *parameter, gpointer user_data);
//
// // Connect action activate signal with callback
// static GSimpleAction* createSimpleAction(const char* name, gpointer user_data) {
//     GSimpleAction *action = g_simple_action_new(name, NULL);
//     g_signal_connect(action, "activate", G_CALLBACK(actionActivateCallback), user_data);
//     return action;
// }
//
// // Add action to action map
// static void addActionToMap(GActionMap *map, GAction *action) {
//     g_action_map_add_action(map, action);
// }
//
// // Set application menu bar
// static void setApplicationMenuBar(GtkApplication* app, GMenuModel* menu_model) {
//     gtk_application_set_menubar(app, menu_model);
// }
import "C"

import (
	"sync"
	"unsafe"
)

// ActionCallback represents a callback for action activation
type ActionCallback func()

var (
	actionCallbacks     = make(map[uintptr]ActionCallback)
	actionCallbackMutex sync.RWMutex
)

//export actionActivateCallback
func actionActivateCallback(action *C.GSimpleAction, parameter *C.GVariant, userData C.gpointer) {
	actionCallbackMutex.RLock()
	defer actionCallbackMutex.RUnlock()

	// Convert action pointer to uintptr for lookup
	actionPtr := uintptr(unsafe.Pointer(action))

	// Find and call the callback
	if callback, ok := actionCallbacks[actionPtr]; ok {
		callback()
	}
}

// ActionGroup represents a GTK action group
type ActionGroup interface {
	AddAction(action *Action)
	RemoveAction(name string)
}

// Action represents a GTK action
type Action struct {
	action *C.GSimpleAction
	name   string
}

// NewAction creates a new GTK action
func NewAction(name string, callback ActionCallback) *Action {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	// Create a simple action with no parameter
	action := C.createSimpleAction(cName, nil)

	// Store action in map
	actionPtr := uintptr(unsafe.Pointer(action))
	
	actionCallbackMutex.Lock()
	actionCallbacks[actionPtr] = callback
	actionCallbackMutex.Unlock()

	return &Action{
		action: action,
		name:   name,
	}
}

// GetNative returns the underlying GAction pointer
func (a *Action) GetNative() *C.GAction {
	return (*C.GAction)(unsafe.Pointer(a.action))
}

// GetName returns the action name
func (a *Action) GetName() string {
	return a.name
}

// Free frees resources associated with the action
func (a *Action) Free() {
	if a.action != nil {
		actionCallbackMutex.Lock()
		actionPtr := uintptr(unsafe.Pointer(a.action))
		delete(actionCallbacks, actionPtr)
		actionCallbackMutex.Unlock()
		
		C.g_object_unref(C.gpointer(unsafe.Pointer(a.action)))
		a.action = nil
	}
}

// ApplicationActionGroup implements the ActionGroup interface for Application
type ApplicationActionGroup struct {
	app *Application
}

// AddAction adds an action to the application
func (a *ApplicationActionGroup) AddAction(action *Action) {
	C.addActionToMap((*C.GActionMap)(unsafe.Pointer(a.app.app)), action.GetNative())
}

// RemoveAction removes an action from the application
func (a *ApplicationActionGroup) RemoveAction(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.g_action_map_remove_action((*C.GActionMap)(unsafe.Pointer(a.app.app)), cName)
}

// GetActionGroup returns the application's action group
func (a *Application) GetActionGroup() ActionGroup {
	return &ApplicationActionGroup{app: a}
}

// SetMenuBar sets the application menu bar
func (a *Application) SetMenuBar(menu *Menu) {
    // In the menu implementation for GTK4, we should set the menu model
    // directly on the application. For a valid implementation, we'd need to
    // ensure the application is registered first, but for this PoC we'll
    // use the direct approach.
    C.setApplicationMenuBar(a.app, menu.GetMenuModel())
}

// Popover represents a GTK popover
type Popover struct {
	BaseWidget
}

// NewPopover creates a new GTK popover
func NewPopover() *Popover {
	popover := &Popover{
		BaseWidget: BaseWidget{
			widget: C.gtk_popover_new(),
		},
	}

	SetupFinalization(popover, popover.Destroy)
	return popover
}

// SetChild sets the child widget for the popover
func (p *Popover) SetChild(child Widget) {
	C.gtk_popover_set_child(
		(*C.GtkPopover)(unsafe.Pointer(p.widget)),
		child.GetWidget(),
	)
}

// SetPointingTo sets the rectangle the popover points to
func (p *Popover) SetPointingTo(rect *C.GdkRectangle) {
	C.gtk_popover_set_pointing_to(
		(*C.GtkPopover)(unsafe.Pointer(p.widget)),
		rect,
	)
}

// SetPosition sets the position of the popover relative to the parent widget
func (p *Popover) SetPosition(position C.GtkPositionType) {
	C.gtk_popover_set_position(
		(*C.GtkPopover)(unsafe.Pointer(p.widget)),
		position,
	)
}

// SetAutohide sets whether to hide the popover when clicked outside
func (p *Popover) SetAutohide(autohide bool) {
	var cautohide C.gboolean
	if autohide {
		cautohide = C.TRUE
	} else {
		cautohide = C.FALSE
	}
	C.gtk_popover_set_autohide(
		(*C.GtkPopover)(unsafe.Pointer(p.widget)),
		cautohide,
	)
}

// SetDefaultWidget sets the default widget for the popover
func (p *Popover) SetDefaultWidget(widget Widget) {
	C.gtk_popover_set_default_widget(
		(*C.GtkPopover)(unsafe.Pointer(p.widget)),
		widget.GetWidget(),
	)
}

// Popup shows the popover
func (p *Popover) Popup() {
	C.gtk_popover_popup((*C.GtkPopover)(unsafe.Pointer(p.widget)))
}

// Popdown hides the popover
func (p *Popover) Popdown() {
	C.gtk_popover_popdown((*C.GtkPopover)(unsafe.Pointer(p.widget)))
}