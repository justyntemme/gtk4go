// Package gtk4 provides modern action-based menu functionality for GTK4
// File: gtk4go/gtk4/action.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Export our action activation callback
// extern void actionActivateCallback(GSimpleAction *action, GVariant *parameter, gpointer user_data);
//
// // Create simple action without connecting a callback
// static GSimpleAction* createSimpleAction(const char* name) {
//     return g_simple_action_new(name, NULL);
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
//
// // Helper function to explicitly connect the activate signal
// static gulong connectActionActivate(GSimpleAction *action, gpointer callback_data) {
//     return g_signal_connect(action, "activate", G_CALLBACK(actionActivateCallback), callback_data);
// }
import "C"

import (
	"unsafe"
)

// ActionCallback represents a callback for action activation
type ActionCallback func()

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
	action := C.createSimpleAction(cName)

	// Create the Action instance
	a := &Action{
		action: action,
		name:   name,
	}

	// Connect the callback using the DIRECT C pointer of the action
	// This is critical - we need to register the callback with the exact pointer
	// that will be used in the actionActivateCallback function
	actionPtr := uintptr(unsafe.Pointer(action))

	// Debug the action pointer for reference
	DebugLog(DebugLevelInfo, DebugComponentAction, "Registering action %s with pointer %p", name, unsafe.Pointer(action))

	// Convert the ActionCallback to a standard func() before storing
	// This is crucial because the callback system doesn't recognize ActionCallback
	// but it does recognize plain func()
	standardCallback := func() {
		callback()
	}

	// Store the callback directly with the action pointer as the key
	StoreDirectCallback(actionPtr, SignalActionActivate, standardCallback)

	// Also connect directly to the action's activate signal
	// This ensures the signal is triggered when a menu item is clicked
	C.connectActionActivate(action, C.gpointer(unsafe.Pointer(action)))

	return a
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
		// Disconnect all callbacks using the unified callback system
		// Use the direct action pointer for disconnection
		actionPtr := uintptr(unsafe.Pointer(a.action))

		// Remove from objectCallbacks map
		globalCallbackManager.objectCallbacks.Delete(actionPtr)

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

//export actionActivateCallback
func actionActivateCallback(action *C.GSimpleAction, parameter *C.GVariant, userData C.gpointer) {
	DebugLog(DebugLevelVerbose, DebugComponentAction, "Action activated: %p", unsafe.Pointer(action))

	// Convert action pointer to uintptr for lookup
	actionPtr := uintptr(unsafe.Pointer(action))

	// Get the callback from the objectCallbacks map directly
	var callback interface{}

	// Try to find the callback for this action pointer
	callbackMapObj, ok := globalCallbackManager.objectCallbacks.Load(actionPtr)
	if ok {
		callbackMap := callbackMapObj.(map[SignalType]interface{})
		callback = callbackMap[SignalActionActivate]
	}

	if callback != nil {
		// Execute the callback
		DebugLog(DebugLevelInfo, DebugComponentAction, "Found callback for action: %p", unsafe.Pointer(action))
		SafeCallback(callback)
	} else {
		DebugLog(DebugLevelWarning, DebugComponentAction,
			"No callback found for action: %p (action may not be registered correctly)", unsafe.Pointer(action))

		// Dump the current action registrations for debugging
		dumpActionCallbacks()
	}
}

// Debug helper to dump all registered action callbacks
func dumpActionCallbacks() {
	callbackCount := 0
	globalCallbackManager.objectCallbacks.Range(func(key, value interface{}) bool {
		ptr := key.(uintptr)
		callbackMap := value.(map[SignalType]interface{})

		// See if it has an action activate signal
		if cb, ok := callbackMap[SignalActionActivate]; ok {
			callbackCount++
			DebugLog(DebugLevelInfo, DebugComponentAction,
				"Registered action callback: ptr=%v, callback=%T", ptr, cb)
		}
		return true
	})

	DebugLog(DebugLevelInfo, DebugComponentAction,
		"Found %d registered action callbacks", callbackCount)
}
