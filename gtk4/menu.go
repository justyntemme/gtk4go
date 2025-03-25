// Package gtk4 provides modern menu functionality for GTK4
// File: gtk4go/gtk4/menu.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Helper functions for GMenu operations
// static GMenuItem* create_menu_item(const char* label, const char* action) {
//     GMenuItem* item = g_menu_item_new(label, action);
//     return item;
// }
//
// static GMenu* create_menu() {
//     return g_menu_new();
// }
//
// static void append_item_to_menu(GMenu* menu, GMenuItem* item) {
//     g_menu_append_item(menu, item);
// }
//
// static void append_submenu_to_menu(GMenu* menu, const char* label, GMenu* submenu) {
//     g_menu_append_submenu(menu, label, G_MENU_MODEL(submenu));
// }
//
// // PopoverMenu helper functions
// static GtkWidget* create_popover_menu_from_model(GMenuModel* model) {
//     return gtk_popover_menu_new_from_model(model);
// }
//
// static void set_popover_menu_model(GtkPopoverMenu* popover, GMenuModel* model) {
//     gtk_popover_menu_set_menu_model(popover, model);
// }
//
// // Additional helper for set parent
// static void set_popover_parent(GtkPopover* popover, GtkWidget* parent) {
//     gtk_widget_set_parent(GTK_WIDGET(popover), parent);
// }
//
// // Signal constants for menus
// static const char* SIGNAL_ACTIVATE = "activate";
// static const char* SIGNAL_SELECTION_CHANGED = "selection-changed";
// static const char* SIGNAL_DEACTIVATE = "deactivate";
import "C"

import (
	"unsafe"
)

// MenuItemActivateCallback represents a callback for menu item activation
type MenuItemActivateCallback func()

// MenuItem represents a menu item
type MenuItem struct {
	item *C.GMenuItem
	name string
}

// NewMenuItem creates a new menu item
func NewMenuItem(label, action string) *MenuItem {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	
	cAction := C.CString(action)
	defer C.free(unsafe.Pointer(cAction))
	
	item := C.create_menu_item(cLabel, cAction)
	
	return &MenuItem{
		item: item,
		name: action,
	}
}

// GetNative returns the underlying GMenuItem pointer
func (mi *MenuItem) GetNative() *C.GMenuItem {
	return mi.item
}

// GetName returns the action name of the menu item
func (mi *MenuItem) GetName() string {
	return mi.name
}

// ConnectActivate connects a callback for when the menu item is activated
func (mi *MenuItem) ConnectActivate(callback MenuItemActivateCallback) uint64 {
	// Use the unified callback system from callback.go
	return Connect(mi, SignalType(C.GoString(C.SIGNAL_ACTIVATE)), callback)
}

// Menu represents a GTK menu
type Menu struct {
	menu *C.GMenu
}

// NewMenu creates a new GTK menu
func NewMenu() *Menu {
	menu := &Menu{
		menu: C.create_menu(),
	}
	return menu
}

// AppendItem adds a menu item to the menu
func (m *Menu) AppendItem(item *MenuItem) {
	C.append_item_to_menu(m.menu, item.item)
}

// AppendSubmenu adds a submenu to the menu
func (m *Menu) AppendSubmenu(label string, submenu *Menu) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	
	C.append_submenu_to_menu(m.menu, cLabel, submenu.menu)
}

// GetMenuModel returns the underlying GMenuModel
func (m *Menu) GetMenuModel() *C.GMenuModel {
	return (*C.GMenuModel)(unsafe.Pointer(m.menu))
}

// GetNative returns the underlying GMenu pointer for callback registration
func (m *Menu) GetNative() uintptr {
	return uintptr(unsafe.Pointer(m.menu))
}

// MenuBar represents a GTK menu bar
type MenuBar struct {
	BaseWidget
}

// NewMenuBar creates a new GTK menu bar
func NewMenuBar() *MenuBar {
	menuBar := &MenuBar{
		BaseWidget: BaseWidget{
			widget: C.gtk_popover_menu_bar_new_from_model(nil),
		},
	}

	SetupFinalization(menuBar, menuBar.Destroy)
	return menuBar
}

// SetMenuModel sets the menu model for the menu bar
func (mb *MenuBar) SetMenuModel(menu *Menu) {
	C.gtk_popover_menu_bar_set_menu_model(
		(*C.GtkPopoverMenuBar)(unsafe.Pointer(mb.widget)),
		menu.GetMenuModel(),
	)
}

// ConnectSelectionChanged connects a callback for selection changes in the menu bar
func (mb *MenuBar) ConnectSelectionChanged(callback func()) uint64 {
	return Connect(mb, SignalType(C.GoString(C.SIGNAL_SELECTION_CHANGED)), callback)
}

// PopoverMenu represents a GTK popover menu
type PopoverMenu struct {
	BaseWidget
}

// NewPopoverMenu creates a new GTK popover menu from a menu model
func NewPopoverMenu(menu *Menu) *PopoverMenu {
	popoverMenu := &PopoverMenu{
		BaseWidget: BaseWidget{
			widget: C.create_popover_menu_from_model(menu.GetMenuModel()),
		},
	}

	SetupFinalization(popoverMenu, popoverMenu.Destroy)
	return popoverMenu
}

// SetMenuModel sets the menu model for the popover menu
func (pm *PopoverMenu) SetMenuModel(menu *Menu) {
	C.set_popover_menu_model(
		(*C.GtkPopoverMenu)(unsafe.Pointer(pm.widget)),
		menu.GetMenuModel(),
	)
}

// SetParent sets the parent widget for the popover
func (pm *PopoverMenu) SetParent(parent Widget) {
	C.set_popover_parent(
		(*C.GtkPopover)(unsafe.Pointer(pm.widget)),
		parent.GetWidget(),
	)
}

// Popup shows the popover
func (pm *PopoverMenu) Popup() {
	C.gtk_popover_popup((*C.GtkPopover)(unsafe.Pointer(pm.widget)))
}

// Popdown hides the popover
func (pm *PopoverMenu) Popdown() {
	C.gtk_popover_popdown((*C.GtkPopover)(unsafe.Pointer(pm.widget)))
}

// ConnectDeactivate connects a callback for when the popover is closed
func (pm *PopoverMenu) ConnectDeactivate(callback func()) uint64 {
	return Connect(pm, SignalType(C.GoString(C.SIGNAL_DEACTIVATE)), callback)
}

// Destroy overrides BaseWidget's Destroy to clean up resources
func (pm *PopoverMenu) Destroy() {
	// Clean up all callbacks using the unified system
	DisconnectAll(pm)
	
	// Call the base method
	pm.BaseWidget.Destroy()
}

// MenuButton represents a GTK menu button
type MenuButton struct {
	BaseWidget
}

// NewMenuButton creates a new GTK menu button
func NewMenuButton() *MenuButton {
	menuButton := &MenuButton{
		BaseWidget: BaseWidget{
			widget: C.gtk_menu_button_new(),
		},
	}

	SetupFinalization(menuButton, menuButton.Destroy)
	return menuButton
}

// SetMenuModel sets the menu model for the menu button
func (mb *MenuButton) SetMenuModel(menu *Menu) {
	C.gtk_menu_button_set_menu_model(
		(*C.GtkMenuButton)(unsafe.Pointer(mb.widget)),
		menu.GetMenuModel(),
	)
}

// SetLabel sets the label for the menu button
func (mb *MenuButton) SetLabel(label string) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	
	C.gtk_menu_button_set_label(
		(*C.GtkMenuButton)(unsafe.Pointer(mb.widget)),
		cLabel,
	)
}

// SetPopover sets a popover for the menu button
func (mb *MenuButton) SetPopover(popover *PopoverMenu) {
	C.gtk_menu_button_set_popover(
		(*C.GtkMenuButton)(unsafe.Pointer(mb.widget)),
		popover.widget,
	)
}

// Destroy overrides BaseWidget's Destroy to clean up resources
func (mb *MenuButton) Destroy() {
	// Clean up all callbacks using the unified system
	DisconnectAll(mb)
	
	// Call the base method
	mb.BaseWidget.Destroy()
}