// Package gtk4 provides ListView widget functionality for GTK4
// File: gtk4go/gtk4/listView.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback functions for list view signals
// extern void listViewItemActivatedCallback(GtkListView *listView, guint position, gpointer user_data);
// extern void listViewSelectionChangedCallback(GtkSelectionModel *selectionModel, guint position, guint nItems, gpointer user_data);
//
// // Connect list view signals
// static gulong connectListViewItemActivated(GtkWidget *view, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(view), "activate", G_CALLBACK(listViewItemActivatedCallback), user_data);
// }
//
// static gulong connectListViewSelectionChanged(GtkSelectionModel *model, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(model), "selection-changed", G_CALLBACK(listViewSelectionChangedCallback), user_data);
// }
//
// // Setup list item - make static to avoid linker conflicts
// static void setupListItem(GtkListItem *item, gpointer data) {
//     // Add a label to the list item
//     GtkWidget *label = gtk_label_new("");
//     gtk_list_item_set_child(item, label);
// }
//
// // Bind list item - make static to avoid linker conflicts
// static void bindListItem(GtkListItem *item, gpointer data) {
//     // Get the data from the list item
//     GtkWidget *label = gtk_list_item_get_child(item);
//     
//     // Get the position of the item
//     unsigned int position = gtk_list_item_get_position(item);
//     
//     // Set the text of the label (this would come from the model in real use)
//     char buf[256];
//     snprintf(buf, sizeof(buf), "Item %u", position);
//     gtk_label_set_text(GTK_LABEL(label), buf);
// }
//
// // Create list item factory
// static GtkListItemFactory* createListItemFactory() {
//     GtkSignalListItemFactory *factory = (GtkSignalListItemFactory*)gtk_signal_list_item_factory_new();
//     g_signal_connect(factory, "setup", G_CALLBACK(setupListItem), NULL);
//     g_signal_connect(factory, "bind", G_CALLBACK(bindListItem), NULL);
//     return GTK_LIST_ITEM_FACTORY(factory);
// }
//
// // Helper functions to check selection model types
// static gboolean isSingleSelection(GtkSelectionModel *model) {
//     return GTK_IS_SINGLE_SELECTION(model) ? TRUE : FALSE;
// }
//
// static gboolean isMultiSelection(GtkSelectionModel *model) {
//     return GTK_IS_MULTI_SELECTION(model) ? TRUE : FALSE;
// }
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// ListItemActivatedCallback represents a callback for list item activated events
type ListItemActivatedCallback func(position int)

// ListViewSelectionChangedCallback represents a callback for selection changed events
type ListViewSelectionChangedCallback func(position int, nItems int)

var (
	listViewItemActivatedCallbacks = make(map[uintptr]ListItemActivatedCallback)
	listViewSelectionChangedCallbacks = make(map[uintptr]ListViewSelectionChangedCallback)
	listViewCallbackMutex sync.RWMutex
)

//export listViewItemActivatedCallback
func listViewItemActivatedCallback(listView *C.GtkListView, position C.guint, userData C.gpointer) {
	listViewCallbackMutex.RLock()
	defer listViewCallbackMutex.RUnlock()

	// Convert list view pointer to uintptr for lookup
	listViewPtr := uintptr(unsafe.Pointer(listView))

	// Find and call the callback
	if callback, ok := listViewItemActivatedCallbacks[listViewPtr]; ok {
		callback(int(position))
	}
}

//export listViewSelectionChangedCallback
func listViewSelectionChangedCallback(selectionModel *C.GtkSelectionModel, position C.guint, nItems C.guint, userData C.gpointer) {
	listViewCallbackMutex.RLock()
	defer listViewCallbackMutex.RUnlock()

	// Convert selection model pointer to uintptr for lookup
	modelPtr := uintptr(unsafe.Pointer(selectionModel))

	// Find and call the callback
	if callback, ok := listViewSelectionChangedCallbacks[modelPtr]; ok {
		callback(int(position), int(nItems))
	}
}

// ListViewOption is a function that configures a list view
type ListViewOption func(*ListView)

// ListView represents a GTK list view widget
type ListView struct {
	BaseWidget
	model ListModel
	selectionModel *C.GtkSelectionModel
	goListModel *goListModel // Wrapper around our Go ListModel
}

// NewListView creates a new GTK list view widget with the given model
func NewListView(model ListModel, options ...ListViewOption) *ListView {
	// Create a Go wrapper around our model
	goModel := &goListModel{
		model: model,
	}

	// Create a GTK list view
	listView := &ListView{
		model: model,
		goListModel: goModel,
	}

	// Create GtkNoSelection initially to simplify setup
	var selectionModel *C.GtkSelectionModel
	selectionModel = (*C.GtkSelectionModel)(unsafe.Pointer(
		C.gtk_no_selection_new((*C.GListModel)(unsafe.Pointer(goModel.createListModel()))),
	))
	listView.selectionModel = selectionModel

	// Create the actual list view widget
	listView.widget = C.gtk_list_view_new(
		selectionModel,
		C.createListItemFactory(),
	)

	// Apply options
	for _, option := range options {
		option(listView)
	}

	// Set up proper finalization
	SetupFinalization(listView, listView.Destroy)
	return listView
}

// WithListSingleSelection configures the list view for single selection
func WithListSingleSelection() ListViewOption {
	return func(lv *ListView) {
		// Remove current selection model
		if lv.selectionModel != nil {
			C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
		}

		// Create a new single selection model
		lv.selectionModel = (*C.GtkSelectionModel)(unsafe.Pointer(
			C.gtk_single_selection_new((*C.GListModel)(unsafe.Pointer(lv.goListModel.createListModel()))),
		))

		// Update the list view
		C.gtk_list_view_set_model(
			(*C.GtkListView)(unsafe.Pointer(lv.widget)),
			lv.selectionModel,
		)
	}
}

// WithListMultiSelection configures the list view for multiple selection
func WithListMultiSelection() ListViewOption {
	return func(lv *ListView) {
		// Remove current selection model
		if lv.selectionModel != nil {
			C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
		}

		// Create a new multi selection model
		lv.selectionModel = (*C.GtkSelectionModel)(unsafe.Pointer(
			C.gtk_multi_selection_new((*C.GListModel)(unsafe.Pointer(lv.goListModel.createListModel()))),
		))

		// Update the list view
		C.gtk_list_view_set_model(
			(*C.GtkListView)(unsafe.Pointer(lv.widget)),
			lv.selectionModel,
		)
	}
}

// ConnectItemActivated connects a callback to the item-activated signal
func (lv *ListView) ConnectItemActivated(callback ListItemActivatedCallback) {
	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Store callback in map
	listViewPtr := uintptr(unsafe.Pointer(lv.widget))
	listViewItemActivatedCallbacks[listViewPtr] = callback

	// Connect signal
	C.connectListViewItemActivated(lv.widget, C.gpointer(unsafe.Pointer(lv.widget)))
}

// ConnectSelectionChanged connects a callback to the selection-changed signal
func (lv *ListView) ConnectSelectionChanged(callback ListViewSelectionChangedCallback) {
	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Store callback in map using selection model as key
	modelPtr := uintptr(unsafe.Pointer(lv.selectionModel))
	listViewSelectionChangedCallbacks[modelPtr] = callback

	// Connect signal
	C.connectListViewSelectionChanged(lv.selectionModel, C.gpointer(unsafe.Pointer(lv.selectionModel)))
}

// GetModel returns the model used by the list view
func (lv *ListView) GetModel() ListModel {
	return lv.model
}

// GetSelectedItem returns the selected item or nil if no item is selected
func (lv *ListView) GetSelectedItem() interface{} {
	// Check if we're using a single selection
	singleSelection := C.isSingleSelection(lv.selectionModel)
	if singleSelection == C.FALSE {
		// Not a single selection model
		return nil
	}

	// Get the selected item from the single selection model
	selected := C.gtk_single_selection_get_selected(
		(*C.GtkSingleSelection)(unsafe.Pointer(lv.selectionModel)),
	)
	
	if selected == C.GTK_INVALID_LIST_POSITION {
		return nil
	}

	// Return the item from our Go model
	return lv.model.GetItem(int(selected))
}

// GetSelectedItems returns the selected items (for multi-selection)
func (lv *ListView) GetSelectedItems() []interface{} {
	// Check if we're using a multi selection
	multiSelection := C.isMultiSelection(lv.selectionModel)
	if multiSelection == C.FALSE {
		// Not a multi selection model
		if item := lv.GetSelectedItem(); item != nil {
			return []interface{}{item}
		}
		return []interface{}{}
	}

	// Get the number of items in the model
	var items []interface{}
	n := lv.model.GetNItems()
	
	// For each item, check if it's selected
	for i := 0; i < n; i++ {
		isSelected := C.gtk_selection_model_is_selected(
			lv.selectionModel,
			C.guint(i),
		)
		
		if isSelected == C.TRUE {
			items = append(items, lv.model.GetItem(i))
		}
	}
	
	return items
}

// Destroy cleans up resources
func (lv *ListView) Destroy() {
	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Remove callbacks
	listViewPtr := uintptr(unsafe.Pointer(lv.widget))
	delete(listViewItemActivatedCallbacks, listViewPtr)
	modelPtr := uintptr(unsafe.Pointer(lv.selectionModel))
	delete(listViewSelectionChangedCallbacks, modelPtr)

	// Clean up selection model
	if lv.selectionModel != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
		lv.selectionModel = nil
	}

	// Call base destroy
	lv.BaseWidget.Destroy()
}

// goListModel is a wrapper around a Go ListModel that implements GListModel
type goListModel struct {
	model ListModel
	listModel *C.GObject // The actual GListModel instance
}

// createListModel creates a new GListModel from our Go model
func (m *goListModel) createListModel() *C.GObject {
	// In a real implementation, we would create a custom GListModel using GObject
	// For this example, we'll use a placeholder string list model
	items := C.g_list_store_new(C.g_type_from_name(C.CString("gchararray")))
	
	// Add placeholders for each item in our model
	n := m.model.GetNItems()
	for i := 0; i < n; i++ {
		cstr := C.CString(fmt.Sprintf("Item %d", i))
		gstr := C.g_strdup(cstr)
		C.free(unsafe.Pointer(cstr))
		
		C.g_list_store_append(items, C.gpointer(unsafe.Pointer(gstr)))
	}
	
	m.listModel = (*C.GObject)(unsafe.Pointer(items))
	return m.listModel
}