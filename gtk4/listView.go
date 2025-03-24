// Package gtk4 provides ListView functionality for GTK4
// File: gtk4go/gtk4/listView.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback functions for list view
// extern void listViewItemActivatedCallback(GtkListView *list_view, guint position, gpointer user_data);
// extern void listViewSelectionChangedCallback(GtkSelectionModel *model, guint position, guint n_items, gpointer user_data);
//
// // Connect list view signals
// static gulong listViewConnectItemActivated(GtkListView *list_view, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(list_view), "activate", G_CALLBACK(listViewItemActivatedCallback), user_data);
// }
//
// static gulong listViewConnectSelectionChanged(GtkSelectionModel *model, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(model), "selection-changed", G_CALLBACK(listViewSelectionChangedCallback), user_data);
// }
//
// // Create a default GListStore with GObject type
// static GListStore* createDefaultStore() {
//     return g_list_store_new(G_TYPE_OBJECT);
// }
//
// // Append to a GListStore
// static void listStoreAppend(GListStore* store, gpointer item) {
//     if (store != NULL && item != NULL) {
//         g_list_store_append(store, item);
//     }
// }
//
// // ListView creation and configuration
// static GtkWidget* createListView(GtkSelectionModel *model, GtkListItemFactory *factory) {
//     return gtk_list_view_new(model, factory);
// }
//
// // Get item from selection model
// static gpointer getSelectedItem(GtkSelectionModel *model) {
//     if (model == NULL) return NULL;
//     guint selected = gtk_single_selection_get_selected(GTK_SINGLE_SELECTION(model));
//     if (selected == GTK_INVALID_LIST_POSITION) return NULL;
//     return g_list_model_get_item(gtk_single_selection_get_model(GTK_SINGLE_SELECTION(model)), selected);
// }
//
// // ListView configuration
// static void listViewSetShowSeparators(GtkListView *list_view, gboolean show_separators) {
//     gtk_list_view_set_show_separators(list_view, show_separators);
// }
//
// static void listViewSetSingleClickActivate(GtkListView *list_view, gboolean single_click_activate) {
//     gtk_list_view_set_single_click_activate(list_view, single_click_activate);
// }
//
// static void listViewSetEnableRubberband(GtkListView *list_view, gboolean enable_rubberband) {
//     gtk_list_view_set_enable_rubberband(list_view, enable_rubberband);
// }
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// ListView activation callback
type ListViewItemActivatedCallback func(position int)

// ListView selection changed callback
type ListViewSelectionChangedCallback func(position int, nItems int)

// Callback storage
var (
	listViewActivatedCallbacks  = make(map[uintptr]ListViewItemActivatedCallback)
	listViewSelectionCallbacks  = make(map[uintptr]ListViewSelectionChangedCallback)
	listViewCallbackMutex       sync.RWMutex
)

//export listViewItemActivatedCallback
func listViewItemActivatedCallback(listView *C.GtkListView, position C.guint, userData C.gpointer) {
	listViewCallbackMutex.RLock()
	defer listViewCallbackMutex.RUnlock()

	viewPtr := uintptr(unsafe.Pointer(listView))
	if callback, ok := listViewActivatedCallbacks[viewPtr]; ok {
		callback(int(position))
	}
}

//export listViewSelectionChangedCallback
func listViewSelectionChangedCallback(model *C.GtkSelectionModel, position C.guint, nItems C.guint, userData C.gpointer) {
	listViewCallbackMutex.RLock()
	defer listViewCallbackMutex.RUnlock()

	modelPtr := uintptr(unsafe.Pointer(model))
	if callback, ok := listViewSelectionCallbacks[modelPtr]; ok {
		callback(int(position), int(nItems))
	}
}

// ListViewOption is a function that configures a ListView
type ListViewOption func(*ListView)

// WithFactory sets the factory for the ListView
func WithFactory(factory *ListItemFactory) ListViewOption {
	return func(lv *ListView) {
		lv.factory = factory
        // No need to update C factory here as it's set during listview creation
	}
}

// WithSelectionMode sets the selection mode for the ListView
func WithSelectionMode(mode SelectionMode) ListViewOption {
	return func(lv *ListView) {
		lv.selectionMode = mode
        // Selection mode is applied during model initialization
	}
}

// WithShowSeparators sets whether to show separators between items
func WithShowSeparators(show bool) ListViewOption {
	return func(lv *ListView) {
		var cShow C.gboolean
		if show {
			cShow = C.TRUE
		} else {
			cShow = C.FALSE
		}
		
		if lv.widget != nil {
		    C.listViewSetShowSeparators((*C.GtkListView)(unsafe.Pointer(lv.widget)), cShow)
		}
	}
}

// WithSingleClickActivate sets whether items are activated with a single click
func WithSingleClickActivate(singleClick bool) ListViewOption {
	return func(lv *ListView) {
		var cSingleClick C.gboolean
		if singleClick {
			cSingleClick = C.TRUE
		} else {
			cSingleClick = C.FALSE
		}
		
		if lv.widget != nil {
		    C.listViewSetSingleClickActivate((*C.GtkListView)(unsafe.Pointer(lv.widget)), cSingleClick)
		}
	}
}

// WithRubberband sets whether rubberband selection is enabled
func WithRubberband(enable bool) ListViewOption {
	return func(lv *ListView) {
		var cEnable C.gboolean
		if enable {
			cEnable = C.TRUE
		} else {
			cEnable = C.FALSE
		}
		
		if lv.widget != nil {
		    C.listViewSetEnableRubberband((*C.GtkListView)(unsafe.Pointer(lv.widget)), cEnable)
		}
	}
}

// ListView represents a GTK ListView widget
type ListView struct {
	BaseWidget
	model         ListModel
	factory       *ListItemFactory
	selectionMode SelectionMode
	selectionModel *C.GtkSelectionModel
}

// NewListView creates a new GTK ListView with a model and options
func NewListView(model ListModel, options ...ListViewOption) *ListView {
	fmt.Println("Creating new ListView...")
	
	// Create the ListView instance
	listView := &ListView{
		model:         model,
		selectionMode: SelectionModeSingle, // Default selection mode
	}
	
	// Apply factory option first if provided
	for _, option := range options {
		option(listView)
	}
	
	// Ensure we have a factory
	if listView.factory == nil {
		fmt.Println("No factory provided, creating default text factory")
		listView.factory = TextFactory()
	}

	// Create selection model based on the model
	var glistModel *C.GListModel
	if glm, ok := model.(*GListModel); ok {
		glistModel = glm.model
	} else if ls, ok := model.(*ListStore); ok {
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
	} else {
		// For a custom ListModel implementation, use a default store
		fmt.Println("Custom model implementation not supported fully, using default store")
		store := C.createDefaultStore()
		glistModel = (*C.GListModel)(unsafe.Pointer(store))
	}
	
	// Create appropriate selection model based on mode
	if glistModel != nil {
		switch listView.selectionMode {
		case SelectionModeSingle:
			listView.selectionModel = (*C.GtkSelectionModel)(unsafe.Pointer(C.gtk_single_selection_new(glistModel)))
		case SelectionModeMultiple:
			listView.selectionModel = (*C.GtkSelectionModel)(unsafe.Pointer(C.gtk_multi_selection_new(glistModel)))
		case SelectionModeNone:
			listView.selectionModel = (*C.GtkSelectionModel)(unsafe.Pointer(C.gtk_no_selection_new(glistModel)))
		default:
			listView.selectionModel = (*C.GtkSelectionModel)(unsafe.Pointer(C.gtk_single_selection_new(glistModel)))
		}
	} else {
		fmt.Println("WARNING: Could not create GListModel, ListView will not function properly")
		return nil
	}
	
	// Create the ListView widget
	if listView.factory != nil && listView.factory.factory != nil {
		listView.widget = C.createListView(
			listView.selectionModel,
			listView.factory.factory,
		)
	} else {
		fmt.Println("ERROR: Factory is nil or has nil factory pointer")
		return nil
	}
	
	if listView.widget == nil {
		fmt.Println("ERROR: Failed to create ListView widget")
		return nil
	}

	// Apply remaining options
	for _, option := range options {
		option(listView)
	}
	
	SetupFinalization(listView, listView.Destroy)
	fmt.Println("ListView created successfully")
	return listView
}

// ConnectItemActivated connects a callback for item activation
func (lv *ListView) ConnectItemActivated(callback ListViewItemActivatedCallback) {
	if lv.widget == nil {
		return
	}

	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Store callback
	viewPtr := uintptr(unsafe.Pointer(lv.widget))
	listViewActivatedCallbacks[viewPtr] = callback

	// Connect signal
	C.listViewConnectItemActivated(
		(*C.GtkListView)(unsafe.Pointer(lv.widget)),
		C.gpointer(unsafe.Pointer(lv.widget)),
	)
}

// ConnectSelectionChanged connects a callback for selection changes
func (lv *ListView) ConnectSelectionChanged(callback ListViewSelectionChangedCallback) {
	if lv.selectionModel == nil {
		return
	}

	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Store callback
	modelPtr := uintptr(unsafe.Pointer(lv.selectionModel))
	listViewSelectionCallbacks[modelPtr] = callback

	// Connect signal
	C.listViewConnectSelectionChanged(
		lv.selectionModel,
		C.gpointer(unsafe.Pointer(lv.selectionModel)),
	)
}

// GetSelectedItem gets the currently selected item from the model
func (lv *ListView) GetSelectedItem() interface{} {
	if lv.selectionModel == nil {
		return nil
	}

	// Get selected item from C
	itemPtr := C.getSelectedItem(lv.selectionModel)
	if itemPtr == nil {
		return nil
	}
	
	// We need to unref the item when we're done with it
	defer C.g_object_unref(itemPtr)

	// Try to get string from GObject
	str := GetStringFromObject((*C.GObject)(itemPtr))
	if str != "" {
		return str
	}

	// If not a string, return as is
	return itemPtr
}

// SetModel sets a new model for the list view
func (lv *ListView) SetModel(model ListModel) {
	// This is a simplified implementation - a full implementation would need
	// to update the C widget properly
	lv.model = model
}

// Destroy destroys the list view and frees resources
func (lv *ListView) Destroy() {
	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Remove callbacks
	if lv.widget != nil {
		delete(listViewActivatedCallbacks, uintptr(unsafe.Pointer(lv.widget)))
	}
	
	if lv.selectionModel != nil {
		delete(listViewSelectionCallbacks, uintptr(unsafe.Pointer(lv.selectionModel)))
	}

	// Free selection model
	if lv.selectionModel != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
		lv.selectionModel = nil
	}

	// Clean up base widget
	lv.BaseWidget.Destroy()
}