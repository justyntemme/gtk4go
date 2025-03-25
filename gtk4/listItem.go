// Package gtk4 provides list item functionality for GTK4
// File: gtk4go/gtk4/listitem.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// static GtkWidget* listItemGetChild(GtkListItem *list_item) {
//     return gtk_list_item_get_child(list_item);
// }
//
// static void listItemSetChild(GtkListItem *list_item, GtkWidget *child) {
//     gtk_list_item_set_child(list_item, child);
// }
//
// static gpointer listItemGetItem(GtkListItem *list_item) {
//     return gtk_list_item_get_item(list_item);
// }
//
// static guint listItemGetPosition(GtkListItem *list_item) {
//     return gtk_list_item_get_position(list_item);
// }
//
// static gboolean listItemGetSelected(GtkListItem *list_item) {
//     return gtk_list_item_get_selected(list_item);
// }
import "C"

import (
	"unsafe"
)

// ListItem represents a GTK list item
type ListItem struct {
	listItem *C.GtkListItem
}

// GetChild returns the child widget of the list item
func (li *ListItem) GetChild() Widget {
	widget := C.listItemGetChild(li.listItem)
	if widget == nil {
		return nil
	}

	// Note: This would need to return a Go wrapper for the widget
	// For now, return a generic BaseWidget as a placeholder
	return &BaseWidget{widget: widget}
}

// SetChild sets the child widget for the list item
func (li *ListItem) SetChild(child Widget) {
	C.listItemSetChild(li.listItem, child.GetWidget())
}

// GetItem returns the model item associated with the list item
func (li *ListItem) GetItem() interface{} {
	item := C.listItemGetItem(li.listItem)
	if item == nil {
		return nil
	}

	// The actual implementation would need to convert the C item
	// to the appropriate Go type based on the model being used
	// For now, just return the raw pointer
	return uintptr(unsafe.Pointer(item))
}

// GetPosition returns the position of the list item in the model
func (li *ListItem) GetPosition() int {
	return int(C.listItemGetPosition(li.listItem))
}

// GetSelected returns whether the list item is selected
func (li *ListItem) GetSelected() bool {
	return bool(C.listItemGetSelected(li.listItem) != 0)
}