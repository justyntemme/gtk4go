// Package gtk4 provides ListView functionality for GTK4
// File: gtk4go/gtk4/listView.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback functions for list view
// extern void listViewItemActivatedCallback(GtkListView *list_view, guint position, gpointer user_data);
//
// // Connect list view signals
// static gulong list_view_connect_item_activated(GtkListView *list_view, gpointer user_data) {
//     if (list_view == NULL) return 0;
//     return g_signal_connect(G_OBJECT(list_view), "activate", G_CALLBACK(listViewItemActivatedCallback), user_data);
// }
//
// // Selection handling
// extern void listViewSelectionModelChangedCallback(GtkSelectionModel *model, guint position, guint n_items, gpointer user_data);
//
// // Connect selection model signals
// static gulong list_view_connect_selection_changed(GtkSelectionModel *model, gpointer user_data) {
//     if (model == NULL) return 0;
//     return g_signal_connect(G_OBJECT(model), "selection-changed", G_CALLBACK(listViewSelectionModelChangedCallback), user_data);
// }
//
// // Create a default GListStore with GObject type
// static GListStore* create_default_store() {
//     return g_list_store_new(G_TYPE_OBJECT);
// }
//
// // ListView creation and configuration
// static GtkWidget* list_view_create_widget(GtkSelectionModel *model, GtkListItemFactory *factory) {
//     if (model == NULL || factory == NULL) {
//         g_warning("ListView creation failed - null model or factory");
//         return NULL;
//     }
//     return gtk_list_view_new(model, factory);
// }
//
// static void list_view_set_model(GtkListView *list_view, GtkSelectionModel *model) {
//     if (list_view == NULL || model == NULL) return;
//     gtk_list_view_set_model(list_view, model);
// }
//
// static void list_view_set_factory(GtkListView *list_view, GtkListItemFactory *factory) {
//     if (list_view == NULL || factory == NULL) return;
//     gtk_list_view_set_factory(list_view, factory);
// }
//
// static void list_view_set_show_separators(GtkListView *list_view, gboolean show_separators) {
//     if (list_view == NULL) return;
//     gtk_list_view_set_show_separators(list_view, show_separators);
// }
//
// static void list_view_set_single_click_activate(GtkListView *list_view, gboolean single_click_activate) {
//     if (list_view == NULL) return;
//     gtk_list_view_set_single_click_activate(list_view, single_click_activate);
// }
//
// static void list_view_set_enable_rubberband(GtkListView *list_view, gboolean enable_rubberband) {
//     if (list_view == NULL) return;
//     gtk_list_view_set_enable_rubberband(list_view, enable_rubberband);
// }
//
// // Selection model creation - ensure item_type is G_TYPE_OBJECT or subclass
// static GtkSelectionModel* list_view_create_single_selection(GListModel *model) {
//     if (model == NULL) return NULL;
//     return GTK_SELECTION_MODEL(gtk_single_selection_new(model));
// }
//
// static GtkSelectionModel* list_view_create_multi_selection(GListModel *model) {
//     if (model == NULL) return NULL;
//     return GTK_SELECTION_MODEL(gtk_multi_selection_new(model));
// }
//
// static GtkSelectionModel* list_view_create_no_selection(GListModel *model) {
//     if (model == NULL) return NULL;
//     return GTK_SELECTION_MODEL(gtk_no_selection_new(model));
// }
//
// // Single selection helpers
// static void list_view_single_selection_set_selected(GtkSingleSelection *selection, guint position) {
//     if (selection == NULL) return;
//     gtk_single_selection_set_selected(selection, position);
// }
//
// static guint list_view_single_selection_get_selected(GtkSingleSelection *selection) {
//     if (selection == NULL) return GTK_INVALID_LIST_POSITION;
//     return gtk_single_selection_get_selected(selection);
// }
//
// static gpointer list_view_single_selection_get_selected_item(GtkSingleSelection *selection) {
//     if (selection == NULL) return NULL;
//     return gtk_single_selection_get_selected_item(selection);
// }
//
// // Multi selection helpers
// static void list_view_multi_selection_select_item(GtkSelectionModel *model, guint position, gboolean select) {
//     if (model == NULL) return;
//     gtk_selection_model_select_item(model, position, select);
// }
//
// static void list_view_multi_selection_select_range(GtkSelectionModel *model, guint position, guint n_items, gboolean select) {
//     if (model == NULL) return;
//     gtk_selection_model_select_range(model, position, n_items, select);
// }
//
// static void list_view_multi_selection_select_all(GtkSelectionModel *model) {
//     if (model == NULL) return;
//     gtk_selection_model_select_all(model);
// }
//
// static void list_view_multi_selection_unselect_all(GtkSelectionModel *model) {
//     if (model == NULL) return;
//     gtk_selection_model_unselect_all(model);
// }
//
// static gboolean list_view_multi_selection_is_selected(GtkSelectionModel *model, guint position) {
//     if (model == NULL) return FALSE;
//     return gtk_selection_model_is_selected(model, position);
// }
//
// // Create a GObject from a string - helper for ListView
// static GObject* create_string_object(const char* text) {
//     if (text == NULL) return NULL;
//     
//     // Create a label widget (which is a GObject) to hold the string
//     GtkWidget* label = gtk_label_new(text);
//     if (label == NULL) return NULL;
//     
//     // Store the original string as a property
//     g_object_set_data_full(G_OBJECT(label), "text", g_strdup(text), g_free);
//     
//     return G_OBJECT(label);
// }
//
// // Extract string from a GObject - helper for ListView
// static const char* get_string_from_object(GObject* obj) {
//     if (obj == NULL) return "";
//     
//     // Try to get the string data we stored
//     const char* text = g_object_get_data(obj, "text");
//     if (text != NULL) return text;
//     
//     // If that fails, try to get text from a label
//     if (GTK_IS_LABEL(obj)) {
//         return gtk_label_get_text(GTK_LABEL(obj));
//     }
//     
//     return "";
// }
import "C"

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"unsafe"
)

// ListViewItemActivatedCallback represents a callback for list item activated events
type ListViewItemActivatedCallback func(position int)

// SelectionChangedCallback represents a callback for selection changed events
type SelectionChangedCallback func(position int, nItems int)

var (
	listViewCallbacks          = make(map[uintptr]ListViewItemActivatedCallback)
	selectionChangedCallbacks  = make(map[uintptr]SelectionChangedCallback)
	listViewCallbackMutex      sync.RWMutex
)

//export listViewItemActivatedCallback
func listViewItemActivatedCallback(listView *C.GtkListView, position C.guint, userData C.gpointer) {
	listViewCallbackMutex.RLock()
	defer listViewCallbackMutex.RUnlock()

	// Convert list view pointer to uintptr for lookup
	listViewPtr := uintptr(unsafe.Pointer(listView))

	// Find and call the callback
	if callback, ok := listViewCallbacks[listViewPtr]; ok {
		callback(int(position))
	}
}

//export listViewSelectionModelChangedCallback
func listViewSelectionModelChangedCallback(model *C.GtkSelectionModel, position C.guint, nItems C.guint, userData C.gpointer) {
	listViewCallbackMutex.RLock()
	defer listViewCallbackMutex.RUnlock()

	// Convert model pointer to uintptr for lookup
	modelPtr := uintptr(unsafe.Pointer(model))

	// Find and call the callback
	if callback, ok := selectionChangedCallbacks[modelPtr]; ok {
		callback(int(position), int(nItems))
	}
}

// SelectionMode defines the type of selection supported
type SelectionMode int

const (
	// SelectionModeNone no selection is possible
	SelectionModeNone SelectionMode = iota
	// SelectionModeSingle only one item can be selected
	SelectionModeSingle
	// SelectionModeMultiple multiple items can be selected
	SelectionModeMultiple
)

// ListViewOption is a function that configures a list view
type ListViewOption func(*ListView)

// ListView represents a GTK list view widget
type ListView struct {
	BaseWidget
	selectionModel *C.GtkSelectionModel
	factory        *ListItemFactory
	model          ListModel
	ownedObjects   []interface{} // Keep references to prevent GC
}

// NewListView creates a new GTK list view widget with the given model
func NewListView(model ListModel, options ...ListViewOption) *ListView {
	fmt.Println("NewListView: Creating ListView...")
	
	// Create a list view with a default configuration
	listView := &ListView{
		model:        model,
		ownedObjects: make([]interface{}, 0),
	}

	// Create a default factory 
	fmt.Println("NewListView: Creating default factory...")
	listView.factory = TextFactory()
	if listView.factory == nil {
		fmt.Println("NewListView: ERROR - Failed to create default factory")
		return nil
	}
	fmt.Printf("NewListView: Default factory created: %v\n", listView.factory)
	listView.ownedObjects = append(listView.ownedObjects, listView.factory)

	// Set up the selection model based on the model
	var glistModel *C.GListModel
	if model == nil {
		fmt.Println("NewListView: WARNING - model is nil, creating default ListStore")
		// Create a default model since model is nil
		ls := NewListStore(G_TYPE_OBJECT)
		if ls == nil {
			fmt.Println("NewListView: ERROR - Failed to create default ListStore")
			return nil
		}
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
		listView.model = ls
		listView.ownedObjects = append(listView.ownedObjects, ls)
	} else if glm, ok := model.(*GListModel); ok {
		fmt.Println("NewListView: Using provided GListModel")
		// Use the wrapped GListModel
		glistModel = glm.model
	} else if ls, ok := model.(*ListStore); ok {
		fmt.Println("NewListView: Using provided ListStore")
		// Use the GListModel from ListStore
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
	} else {
		fmt.Println("NewListView: WARNING - Unknown model type, creating default ListStore")
		// For a custom ListModel implementation, create an empty ListStore as a placeholder
		ls := NewListStore(G_TYPE_OBJECT)
		if ls == nil {
			fmt.Println("NewListView: ERROR - Failed to create wrapper ListStore")
			return nil
		}
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
		listView.model = ls
		listView.ownedObjects = append(listView.ownedObjects, ls)
	}

	// Create a single selection model by default
	fmt.Println("NewListView: Creating selection model...")
	listView.selectionModel = C.list_view_create_single_selection(glistModel)
	if listView.selectionModel == nil {
		fmt.Println("NewListView: ERROR - Failed to create selection model")
		return nil
	}
	fmt.Println("NewListView: Selection model created successfully")

	// Create the widget
	fmt.Println("NewListView: Creating ListView widget...")
	listView.widget = C.list_view_create_widget(
		listView.selectionModel,
		listView.factory.factory,
	)
	
	if listView.widget == nil {
		fmt.Println("NewListView: ERROR - Failed to create ListView widget")
		return nil
	}
	fmt.Println("NewListView: ListView widget created successfully")

	// Apply options
	for i, option := range options {
		fmt.Printf("NewListView: Applying option %d...\n", i)
		option(listView)
	}
	fmt.Println("NewListView: All options applied")

	SetupFinalization(listView, listView.Destroy)
	fmt.Println("NewListView: ListView created successfully")
	return listView
}

// WithSelectionMode sets the selection mode
func WithSelectionMode(mode SelectionMode) ListViewOption {
	return func(lv *ListView) {
		fmt.Printf("WithSelectionMode: Setting mode=%v\n", mode)
		
		// Skip if widget is not created yet
		if lv.widget == nil {
			fmt.Println("WithSelectionMode: WARNING - lv.widget is nil, cannot set mode")
			return
		}

		// Get the GListModel from our current model
		var glistModel *C.GListModel
		if lv.model == nil {
			fmt.Println("WithSelectionMode: WARNING - lv.model is nil, creating default ListStore")
			// Create a default model
			ls := NewListStore(G_TYPE_OBJECT)
			if ls == nil {
				fmt.Println("WithSelectionMode: ERROR - Failed to create ListStore")
				return
			}
			glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
			lv.model = ls
			lv.ownedObjects = append(lv.ownedObjects, ls)
		} else if glm, ok := lv.model.(*GListModel); ok {
			glistModel = glm.model
		} else if ls, ok := lv.model.(*ListStore); ok {
			glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
		} else {
			fmt.Println("WithSelectionMode: WARNING - Unknown model type, creating default ListStore")
			// Cannot get GListModel, so create a new one
			ls := NewListStore(G_TYPE_OBJECT)
			if ls == nil {
				fmt.Println("WithSelectionMode: ERROR - Failed to create wrapper ListStore")
				return
			}
			glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
			lv.model = ls
			lv.ownedObjects = append(lv.ownedObjects, ls)
		}
		
		// Clean up the old selection model
		if lv.selectionModel != nil {
			fmt.Println("WithSelectionMode: Cleaning up old selection model")
			C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
		}
		
		// Create a new selection model based on the mode
		fmt.Printf("WithSelectionMode: Creating new selection model for mode %v\n", mode)
		var newModel *C.GtkSelectionModel
		switch mode {
		case SelectionModeSingle:
			newModel = C.list_view_create_single_selection(glistModel)
		case SelectionModeMultiple:
			newModel = C.list_view_create_multi_selection(glistModel)
		case SelectionModeNone:
			newModel = C.list_view_create_no_selection(glistModel)
		}

		// Only update if we got a valid model
		if newModel != nil {
			lv.selectionModel = newModel
			fmt.Println("WithSelectionMode: Setting new selection model on ListView")
			C.list_view_set_model((*C.GtkListView)(unsafe.Pointer(lv.widget)), lv.selectionModel)
			fmt.Println("WithSelectionMode: Selection mode set successfully")
		} else {
			fmt.Println("WithSelectionMode: ERROR - Failed to create new selection model")
		}
	}
}

// WithFactory sets the factory for creating list items
func WithFactory(factory *ListItemFactory) ListViewOption {
	return func(lv *ListView) {
		fmt.Printf("WithFactory: Setting factory=%v\n", factory)
		
		// Validate factory
		if factory == nil {
			fmt.Println("WithFactory: WARNING - factory is nil, ignoring")
			return
		}
		
		if factory.factory == nil {
			fmt.Println("WithFactory: WARNING - factory.factory is nil, ignoring")
			return
		}
		
		// Validate widget
		if lv.widget == nil {
			fmt.Println("WithFactory: WARNING - lv.widget is nil, cannot set factory")
			return
		}
		
		// Clean up the old factory if we own it
		if lv.factory != nil && lv.factory != factory {
			fmt.Println("WithFactory: Removing reference to old factory")
			// Find and remove from owned objects
			for i, obj := range lv.ownedObjects {
				if fac, ok := obj.(*ListItemFactory); ok && fac == lv.factory {
					lv.ownedObjects = append(lv.ownedObjects[:i], lv.ownedObjects[i+1:]...)
					break
				}
			}
		}
		
		fmt.Println("WithFactory: Setting new factory")
		lv.factory = factory
		
		fmt.Println("WithFactory: Setting factory on ListView")
		C.list_view_set_factory((*C.GtkListView)(unsafe.Pointer(lv.widget)), factory.factory)
		fmt.Println("WithFactory: Factory set successfully")
	}
}

// WithShowSeparators sets whether to show separators between items
func WithShowSeparators(show bool) ListViewOption {
	return func(lv *ListView) {
		fmt.Printf("WithShowSeparators: Setting show=%v\n", show)
		
		if lv.widget == nil {
			fmt.Println("WithShowSeparators: WARNING - lv.widget is nil, ignoring")
			return
		}
		
		var cShow C.gboolean
		if show {
			cShow = C.TRUE
		} else {
			cShow = C.FALSE
		}
		
		C.list_view_set_show_separators((*C.GtkListView)(unsafe.Pointer(lv.widget)), cShow)
		fmt.Println("WithShowSeparators: Show separators set successfully")
	}
}

// WithSingleClickActivate sets whether items can be activated with a single click
func WithSingleClickActivate(enable bool) ListViewOption {
	return func(lv *ListView) {
		fmt.Printf("WithSingleClickActivate: Setting enable=%v\n", enable)
		
		if lv.widget == nil {
			fmt.Println("WithSingleClickActivate: WARNING - lv.widget is nil, ignoring")
			return
		}
		
		var cEnable C.gboolean
		if enable {
			cEnable = C.TRUE
		} else {
			cEnable = C.FALSE
		}
		
		C.list_view_set_single_click_activate((*C.GtkListView)(unsafe.Pointer(lv.widget)), cEnable)
		fmt.Println("WithSingleClickActivate: Single click activate set successfully")
	}
}

// WithRubberband sets whether rubberband selection is enabled
func WithRubberband(enable bool) ListViewOption {
	return func(lv *ListView) {
		fmt.Printf("WithRubberband: Setting enable=%v\n", enable)
		
		if lv.widget == nil {
			fmt.Println("WithRubberband: WARNING - lv.widget is nil, ignoring")
			return
		}
		
		var cEnable C.gboolean
		if enable {
			cEnable = C.TRUE
		} else {
			cEnable = C.FALSE
		}
		
		C.list_view_set_enable_rubberband((*C.GtkListView)(unsafe.Pointer(lv.widget)), cEnable)
		fmt.Println("WithRubberband: Rubberband selection set successfully")
	}
}

// ConnectItemActivated connects a callback to the item-activated signal
func (lv *ListView) ConnectItemActivated(callback ListViewItemActivatedCallback) {
	if lv.widget == nil {
		fmt.Println("ConnectItemActivated: WARNING - lv.widget is nil, ignoring")
		return
	}
	
	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Store callback in map
	listViewPtr := uintptr(unsafe.Pointer(lv.widget))
	listViewCallbacks[listViewPtr] = callback

	// Connect signal
	C.list_view_connect_item_activated((*C.GtkListView)(unsafe.Pointer(lv.widget)), C.gpointer(unsafe.Pointer(lv.widget)))
}

// ConnectSelectionChanged connects a callback to the selection-changed signal
func (lv *ListView) ConnectSelectionChanged(callback SelectionChangedCallback) {
	if lv.selectionModel == nil {
		fmt.Println("ConnectSelectionChanged: WARNING - lv.selectionModel is nil, ignoring")
		return
	}
	
	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Store callback in map
	modelPtr := uintptr(unsafe.Pointer(lv.selectionModel))
	selectionChangedCallbacks[modelPtr] = callback

	// Connect signal
	C.list_view_connect_selection_changed(lv.selectionModel, C.gpointer(unsafe.Pointer(lv.selectionModel)))
}

// GetModel returns the model used by the list view
func (lv *ListView) GetModel() ListModel {
	return lv.model
}

// SetModel sets the model for the list view
func (lv *ListView) SetModel(model ListModel) {
	fmt.Printf("SetModel: Setting model=%v\n", model)
	
	if lv.widget == nil {
		fmt.Println("SetModel: WARNING - lv.widget is nil, ignoring")
		return
	}
	
	lv.model = model
	
	// Create a GListModel wrapper if needed
	var glistModel *C.GListModel
	if model == nil {
		fmt.Println("SetModel: WARNING - model is nil, creating default ListStore")
		// Create a default model
		ls := NewListStore(G_TYPE_OBJECT)
		if ls == nil {
			fmt.Println("SetModel: ERROR - Failed to create ListStore")
			return
		}
		lv.model = ls
		lv.ownedObjects = append(lv.ownedObjects, ls)
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
	} else if glm, ok := model.(*GListModel); ok {
		glistModel = glm.model
	} else if ls, ok := model.(*ListStore); ok {
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
	} else {
		fmt.Println("SetModel: WARNING - Unknown model type, creating wrapper ListStore")
		// For a custom ListModel implementation, we'd need to create a GListModel wrapper
		ls := NewListStore(G_TYPE_OBJECT)
		if ls == nil {
			fmt.Println("SetModel: ERROR - Failed to create wrapper ListStore")
			return
		}
		// Copy items from custom model to ListStore
		// (not implemented here)
		lv.ownedObjects = append(lv.ownedObjects, ls)
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
	}
	
	// Clean up the old selection model
	if lv.selectionModel != nil {
		fmt.Println("SetModel: Cleaning up old selection model")
		// First disconnect any callbacks
		listViewCallbackMutex.Lock()
		delete(selectionChangedCallbacks, uintptr(unsafe.Pointer(lv.selectionModel)))
		listViewCallbackMutex.Unlock()
		
		C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
	}
	
	// Create a new selection model (default to single selection)
	fmt.Println("SetModel: Creating new selection model")
	lv.selectionModel = C.list_view_create_single_selection(glistModel)
	if lv.selectionModel == nil {
		fmt.Println("SetModel: ERROR - Failed to create selection model")
		return
	}
	
	// Update the list view
	fmt.Println("SetModel: Setting new model on ListView")
	C.list_view_set_model((*C.GtkListView)(unsafe.Pointer(lv.widget)), lv.selectionModel)
	fmt.Println("SetModel: Model set successfully")
}

// SetSelectionMode sets the selection mode
func (lv *ListView) SetSelectionMode(mode SelectionMode) {
	WithSelectionMode(mode)(lv)
}

// SetFactory sets the factory for creating list items
func (lv *ListView) SetFactory(factory *ListItemFactory) {
	WithFactory(factory)(lv)
}

// SetShowSeparators sets whether to show separators between items
func (lv *ListView) SetShowSeparators(show bool) {
	WithShowSeparators(show)(lv)
}

// SetSingleClickActivate sets whether items can be activated with a single click
func (lv *ListView) SetSingleClickActivate(enable bool) {
	WithSingleClickActivate(enable)(lv)
}

// SetRubberbandSelection sets whether rubberband selection is enabled
func (lv *ListView) SetRubberbandSelection(enable bool) {
	WithRubberband(enable)(lv)
}

// Selection management for single selection mode

// GetSelectedItem returns the selected item or nil if no item is selected
func (lv *ListView) GetSelectedItem() interface{} {
	if lv.widget == nil || lv.selectionModel == nil || lv.model == nil {
		return nil
	}
	
	// Check if using single selection
	singleSelection := (*C.GtkSingleSelection)(unsafe.Pointer(lv.selectionModel))
	if singleSelection == nil {
		return nil
	}
	
	// Get the selected position
	position := C.list_view_single_selection_get_selected(singleSelection)
	if position == C.GTK_INVALID_LIST_POSITION {
		return nil
	}
	
	// Return the item from the model
	if int(position) >= lv.model.GetNItems() {
		return nil
	}
	return lv.model.GetItem(int(position))
}

// GetSelectedPosition returns the position of the selected item or -1 if no item is selected
func (lv *ListView) GetSelectedPosition() int {
	if lv.widget == nil || lv.selectionModel == nil {
		return -1
	}
	
	// Check if using single selection
	singleSelection := (*C.GtkSingleSelection)(unsafe.Pointer(lv.selectionModel))
	if singleSelection == nil {
		return -1
	}
	
	// Get the selected position
	position := C.list_view_single_selection_get_selected(singleSelection)
	if position == C.GTK_INVALID_LIST_POSITION {
		return -1
	}
	
	return int(position)
}

// SelectItem selects the item at the specified position
func (lv *ListView) SelectItem(position int) {
	if lv.widget == nil || lv.selectionModel == nil || lv.model == nil {
		return
	}
	
	// Check if using single selection
	singleSelection := (*C.GtkSingleSelection)(unsafe.Pointer(lv.selectionModel))
	if singleSelection == nil {
		return
	}
	
	// Check position is in range
	if position < 0 || position >= lv.model.GetNItems() {
		return
	}
	
	// Select the item
	C.list_view_single_selection_set_selected(singleSelection, C.guint(position))
}

// Selection management for multi selection mode

// GetSelectedItems returns the selected items (for multi-selection)
func (lv *ListView) GetSelectedItems() []interface{} {
	if lv.widget == nil || lv.selectionModel == nil || lv.model == nil {
		return nil
	}
	
	// Get the number of items in the model
	var items []interface{}
	n := lv.model.GetNItems()
	
	// For each item, check if it's selected
	for i := 0; i < n; i++ {
		if C.list_view_multi_selection_is_selected(lv.selectionModel, C.guint(i)) == C.TRUE {
			items = append(items, lv.model.GetItem(i))
		}
	}
	
	return items
}

// GetSelectedPositions returns the positions of selected items
func (lv *ListView) GetSelectedPositions() []int {
	if lv.widget == nil || lv.selectionModel == nil || lv.model == nil {
		return nil
	}
	
	// Get the number of items in the model
	var positions []int
	n := lv.model.GetNItems()
	
	// For each item, check if it's selected
	for i := 0; i < n; i++ {
		if C.list_view_multi_selection_is_selected(lv.selectionModel, C.guint(i)) == C.TRUE {
			positions = append(positions, i)
		}
	}
	
	return positions
}

// SelectItems selects the items at the specified positions
func (lv *ListView) SelectItems(positions []int) {
	if lv.widget == nil || lv.selectionModel == nil || lv.model == nil {
		return
	}
	
	// Select each position
	for _, pos := range positions {
		if pos >= 0 && pos < lv.model.GetNItems() {
			C.list_view_multi_selection_select_item(lv.selectionModel, C.guint(pos), C.TRUE)
		}
	}
}

// UnselectItems unselects the items at the specified positions
func (lv *ListView) UnselectItems(positions []int) {
	if lv.widget == nil || lv.selectionModel == nil || lv.model == nil {
		return
	}
	
	// Unselect each position
	for _, pos := range positions {
		if pos >= 0 && pos < lv.model.GetNItems() {
			C.list_view_multi_selection_select_item(lv.selectionModel, C.guint(pos), C.FALSE)
		}
	}
}

// SelectRange selects a range of items
func (lv *ListView) SelectRange(start, count int) {
	if lv.widget == nil || lv.selectionModel == nil || lv.model == nil {
		return
	}
	
	// Check range is valid
	if start < 0 || count <= 0 || start+count > lv.model.GetNItems() {
		return
	}
	
	// Select the range
	C.list_view_multi_selection_select_range(lv.selectionModel, C.guint(start), C.guint(count), C.TRUE)
}

// UnselectRange unselects a range of items
func (lv *ListView) UnselectRange(start, count int) {
	if lv.widget == nil || lv.selectionModel == nil || lv.model == nil {
		return
	}
	
	// Check range is valid
	if start < 0 || count <= 0 || start+count > lv.model.GetNItems() {
		return
	}
	
	// Unselect the range
	C.list_view_multi_selection_select_range(lv.selectionModel, C.guint(start), C.guint(count), C.FALSE)
}

// SelectAll selects all items
func (lv *ListView) SelectAll() {
	if lv.widget == nil || lv.selectionModel == nil {
		return
	}
	
	C.list_view_multi_selection_select_all(lv.selectionModel)
}

// UnselectAll unselects all items
func (lv *ListView) UnselectAll() {
	if lv.widget == nil || lv.selectionModel == nil {
		return
	}
	
	C.list_view_multi_selection_unselect_all(lv.selectionModel)
}

// Destroy cleans up resources
func (lv *ListView) Destroy() {
	fmt.Println("ListView.Destroy: Cleaning up resources")
	
	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Remove callbacks
	if lv.widget != nil {
		delete(listViewCallbacks, uintptr(unsafe.Pointer(lv.widget)))
	}
	
	if lv.selectionModel != nil {
		delete(selectionChangedCallbacks, uintptr(unsafe.Pointer(lv.selectionModel)))
		C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
		lv.selectionModel = nil
	}
	
	// Clean up factory if we own it
	if lv.factory != nil {
		for _, obj := range lv.ownedObjects {
			if fac, ok := obj.(*ListItemFactory); ok && fac == lv.factory {
				fac.Free()
				break
			}
		}
		lv.factory = nil
	}
	
	// Clear owned objects
	lv.ownedObjects = nil

	// Call base destroy
	lv.BaseWidget.Destroy()
}

// ListStore specific helper functions

// CreateStringObject creates a GObject from a string for use in ListStore
func CreateStringObject(text string) *C.GObject {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	
	obj := C.create_string_object(cText)
	if obj == nil {
		fmt.Fprintf(os.Stderr, "WARNING: Failed to create GObject for string\n")
		return nil
	}
	
	return obj
}

// GetStringFromObject extracts string from a GObject
func GetStringFromObject(obj *C.GObject) string {
	if obj == nil {
		return ""
	}
	
	cStr := C.get_string_from_object(obj)
	if cStr == nil {
		return ""
	}
	
	return C.GoString(cStr)
}

// AppendString adds a string to a ListStore
func (s *ListStore) AppendString(text string) {
	fmt.Printf("AppendString: Adding string %q to ListStore\n", text)
	
	obj := CreateStringObject(text)
	if obj == nil {
		fmt.Println("AppendString: ERROR - Failed to create string object")
		return
	}
	
	C.list_store_append(s.store, C.gpointer(unsafe.Pointer(obj)))
	fmt.Println("AppendString: String added successfully")
}