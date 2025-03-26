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
//
// // Helper function to get a string item from a GtkStringObject
// static char* getStringFromObject(gpointer item) {
//     if (item != NULL && GTK_IS_STRING_OBJECT(item)) {
//         return g_strdup(gtk_string_object_get_string(GTK_STRING_OBJECT(item)));
//     }
//     return NULL;
// }
//
// // Helper to check if an object is a GtkStringObject
// static gboolean isStringObject(gpointer item) {
//     return (item != NULL && GTK_IS_STRING_OBJECT(item)) ? TRUE : FALSE;
// }
//
// // Helper for finding a label within a widget hierarchy - GTK4 version
// static GtkLabel* findLabelInChild(GtkWidget* widget) {
//     if (widget == NULL)
//         return NULL;
//         
//     if (GTK_IS_LABEL(widget))
//         return GTK_LABEL(widget);
//         
//     // Check if widget has children (GTK4 approach)
//     GtkWidget* child = gtk_widget_get_first_child(widget);
//     if (child != NULL) {
//         // Try each child
//         while (child != NULL) {
//             if (GTK_IS_LABEL(child))
//                 return GTK_LABEL(child);
//                 
//             // Recursively search this child
//             GtkLabel* label = findLabelInChild(child);
//             if (label != NULL)
//                 return label;
//                 
//             // Try next sibling
//             child = gtk_widget_get_next_sibling(child);
//         }
//     }
//     
//     return NULL;
// }
//
// // Helper to set text on the first label found in a container
// static gboolean setTextOnChildLabel(GtkWidget* container, const char* text) {
//     GtkLabel* label = findLabelInChild(container);
//     if (label != NULL) {
//         gtk_label_set_text(label, text);
//         return TRUE;
//     }
//     return FALSE;
// }
import "C"

import (
	"fmt"
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

	// Return a basic widget wrapper
	return &BaseWidget{widget: widget}
}

// SetChild sets the child widget for the list item
func (li *ListItem) SetChild(child Widget) {
	if child == nil {
		C.listItemSetChild(li.listItem, nil)
	} else {
		C.listItemSetChild(li.listItem, child.GetWidget())
	}
}

// GetItem returns the model item associated with the list item
func (li *ListItem) GetItem() interface{} {
	item := C.listItemGetItem(li.listItem)
	if item == nil {
		return nil
	}

	// Check if it's a string object (common case)
	if C.isStringObject(item) == C.TRUE {
		cstr := C.getStringFromObject(item)
		if cstr != nil {
			// Convert to Go string and free the C string
			str := C.GoString(cstr)
			C.free(unsafe.Pointer(cstr))
			return str
		}
	}

	// For other types, return the raw pointer
	// In a real implementation, we would convert based on the known model type
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

// GetText is a convenience function to get the text from string items
func (li *ListItem) GetText() string {
	// Try to get the item and convert it to a string
	item := li.GetItem()
	
	// If it's already a string, return it directly
	if str, ok := item.(string); ok {
		return str
	}
	
	// If it's a pointer, try to convert via C helper
	if ptr, ok := item.(uintptr); ok {
		if C.isStringObject(C.gpointer(ptr)) == C.TRUE {
			cstr := C.getStringFromObject(C.gpointer(ptr))
			if cstr != nil {
				str := C.GoString(cstr)
				C.free(unsafe.Pointer(cstr))
				return str
			}
		}
	}
	
	// Return empty string if not a string item
	return ""
}

// SetTextOnChildLabel attempts to set text on the first label found in the child widget
func (li *ListItem) SetTextOnChildLabel(text string) bool {
	child := C.listItemGetChild(li.listItem)
	if child == nil {
		return false
	}
	
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	
	result := C.setTextOnChildLabel(child, cText)
	return bool(result != 0)
}

// UpdateChildWithText is a helper function to update a child widget with text
// This function tries several approaches to ensure the text is displayed
func (li *ListItem) UpdateChildWithText() bool {
	// First get the text from the item
	text := li.GetText()
	if text == "" {
		// Try to get a default text based on position
		text = fmt.Sprintf("Item %d", li.GetPosition()+1)
	}
	
	// Try to set the text on a child label
	if li.SetTextOnChildLabel(text) {
		return true
	}
	
	// If that fails and we have a child that's a label, set it directly
	child := li.GetChild()
	if child != nil {
		if label, ok := child.(*Label); ok {
			label.SetText(text)
			return true
		}
	}
	
	// As a last resort, create a new label and set it as the child
	if child == nil {
		label := NewLabel(text)
		li.SetChild(label)
		return true
	}
	
	return false
}