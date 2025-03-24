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
//     return g_signal_connect(G_OBJECT(list_view), "activate", G_CALLBACK(listViewItemActivatedCallback), user_data);
// }
//
// // Selection handling
// extern void listViewSelectionModelChangedCallback(GtkSelectionModel *model, guint position, guint n_items, gpointer user_data);
//
// // Connect selection model signals
// static gulong list_view_connect_selection_changed(GtkSelectionModel *model, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(model), "selection-changed", G_CALLBACK(listViewSelectionModelChangedCallback), user_data);
// }
//
// // ListView creation and configuration
// static GtkWidget* list_view_create_widget(GtkSelectionModel *model, GtkListItemFactory *factory) {
//     return gtk_list_view_new(model, factory);
// }
//
// static void list_view_set_model(GtkListView *list_view, GtkSelectionModel *model) {
//     gtk_list_view_set_model(list_view, model);
// }
//
// static void list_view_set_factory(GtkListView *list_view, GtkListItemFactory *factory) {
//     gtk_list_view_set_factory(list_view, factory);
// }
//
// static void list_view_set_show_separators(GtkListView *list_view, gboolean show_separators) {
//     gtk_list_view_set_show_separators(list_view, show_separators);
// }
//
// static void list_view_set_single_click_activate(GtkListView *list_view, gboolean single_click_activate) {
//     gtk_list_view_set_single_click_activate(list_view, single_click_activate);
// }
//
// static void list_view_set_enable_rubberband(GtkListView *list_view, gboolean enable_rubberband) {
//     gtk_list_view_set_enable_rubberband(list_view, enable_rubberband);
// }
//
// // Selection model creation
// static GtkSelectionModel* list_view_create_single_selection(GListModel *model) {
//     return GTK_SELECTION_MODEL(gtk_single_selection_new(model));
// }
//
// static GtkSelectionModel* list_view_create_multi_selection(GListModel *model) {
//     return GTK_SELECTION_MODEL(gtk_multi_selection_new(model));
// }
//
// static GtkSelectionModel* list_view_create_no_selection(GListModel *model) {
//     return GTK_SELECTION_MODEL(gtk_no_selection_new(model));
// }
//
// // Single selection helpers
// static void list_view_single_selection_set_selected(GtkSingleSelection *selection, guint position) {
//     gtk_single_selection_set_selected(selection, position);
// }
//
// static guint list_view_single_selection_get_selected(GtkSingleSelection *selection) {
//     return gtk_single_selection_get_selected(selection);
// }
//
// static gpointer list_view_single_selection_get_selected_item(GtkSingleSelection *selection) {
//     return gtk_single_selection_get_selected_item(selection);
// }
//
// // Multi selection helpers
// static void list_view_multi_selection_select_item(GtkSelectionModel *model, guint position, gboolean select) {
//     gtk_selection_model_select_item(model, position, select);
// }
//
// static void list_view_multi_selection_select_range(GtkSelectionModel *model, guint position, guint n_items, gboolean select) {
//     gtk_selection_model_select_range(model, position, n_items, select);
// }
//
// static void list_view_multi_selection_select_all(GtkSelectionModel *model) {
//     gtk_selection_model_select_all(model);
// }
//
// static void list_view_multi_selection_unselect_all(GtkSelectionModel *model) {
//     gtk_selection_model_unselect_all(model);
// }
//
// static gboolean list_view_multi_selection_is_selected(GtkSelectionModel *model, guint position) {
//     return gtk_selection_model_is_selected(model, position);
// }
import "C"

import (
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
}

// NewListView creates a new GTK list view widget with the given model
func NewListView(model ListModel, options ...ListViewOption) *ListView {
	// Create a list view with a default configuration
	listView := &ListView{
		model: model,
	}

	// Create a default factory if none provided
	listView.factory = TextFactory()

	// Set up the selection model based on the model
	var glistModel *C.GListModel
	if glm, ok := model.(*GListModel); ok {
		// Use the wrapped GListModel
		glistModel = glm.model
	} else if ls, ok := model.(*ListStore); ok {
		// Use the GListModel from ListStore
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
	} else {
		// For a custom ListModel implementation, we'd need to create a GListModel wrapper
		// This is complex and would need careful implementation
		// For now, we create an empty ListStore as a placeholder
		ls := NewListStore(G_TYPE_OBJECT)
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
	}

	// Create a single selection model by default
	listView.selectionModel = C.list_view_create_single_selection(glistModel)

	// Create the widget
	listView.widget = C.list_view_create_widget(
		listView.selectionModel,
		listView.factory.factory,
	)

	// Apply options
	for _, option := range options {
		option(listView)
	}

	SetupFinalization(listView, listView.Destroy)
	return listView
}

// WithSelectionMode sets the selection mode
func WithSelectionMode(mode SelectionMode) ListViewOption {
	return func(lv *ListView) {
		// Get the current GListModel from the selection model
		var glistModel *C.GListModel
		
		// This is a simplified implementation
		// In a real implementation, you would extract the model properly
		if lv.model != nil {
			if glm, ok := lv.model.(*GListModel); ok {
				glistModel = glm.model
			} else if ls, ok := lv.model.(*ListStore); ok {
				glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
			}
		}
		
		// If no model available, do nothing
		if glistModel == nil {
			return
		}
		
		// Clean up the old selection model
		if lv.selectionModel != nil {
			C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
		}
		
		// Create a new selection model based on the mode
		switch mode {
		case SelectionModeSingle:
			lv.selectionModel = C.list_view_create_single_selection(glistModel)
		case SelectionModeMultiple:
			lv.selectionModel = C.list_view_create_multi_selection(glistModel)
		case SelectionModeNone:
			lv.selectionModel = C.list_view_create_no_selection(glistModel)
		}
		
		// Update the list view
		C.list_view_set_model((*C.GtkListView)(unsafe.Pointer(lv.widget)), lv.selectionModel)
	}
}

// WithFactory sets the factory for creating list items
func WithFactory(factory *ListItemFactory) ListViewOption {
	return func(lv *ListView) {
		// Clean up the old factory if it's our default
		if lv.factory != nil && lv.factory != factory {
			lv.factory.Free()
		}
		
		lv.factory = factory
		C.list_view_set_factory((*C.GtkListView)(unsafe.Pointer(lv.widget)), factory.factory)
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
		C.list_view_set_show_separators((*C.GtkListView)(unsafe.Pointer(lv.widget)), cShow)
	}
}

// WithSingleClickActivate sets whether items can be activated with a single click
func WithSingleClickActivate(enable bool) ListViewOption {
	return func(lv *ListView) {
		var cEnable C.gboolean
		if enable {
			cEnable = C.TRUE
		} else {
			cEnable = C.FALSE
		}
		C.list_view_set_single_click_activate((*C.GtkListView)(unsafe.Pointer(lv.widget)), cEnable)
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
		C.list_view_set_enable_rubberband((*C.GtkListView)(unsafe.Pointer(lv.widget)), cEnable)
	}
}

// ConnectItemActivated connects a callback to the item-activated signal
func (lv *ListView) ConnectItemActivated(callback ListViewItemActivatedCallback) {
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
	lv.model = model
	
	// Create a GListModel wrapper if needed
	var glistModel *C.GListModel
	if glm, ok := model.(*GListModel); ok {
		glistModel = glm.model
	} else if ls, ok := model.(*ListStore); ok {
		glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
	} else {
		// For a custom ListModel implementation, we'd need to create a GListModel wrapper
		// This is complex and would need careful implementation
		// For now, return without updating
		return
	}
	
	// Clean up the old selection model
	if lv.selectionModel != nil {
		// First disconnect any callbacks
		listViewCallbackMutex.Lock()
		delete(selectionChangedCallbacks, uintptr(unsafe.Pointer(lv.selectionModel)))
		listViewCallbackMutex.Unlock()
		
		C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
	}
	
	// Create a new selection model (default to single selection)
	lv.selectionModel = C.list_view_create_single_selection(glistModel)
	
	// Update the list view
	C.list_view_set_model((*C.GtkListView)(unsafe.Pointer(lv.widget)), lv.selectionModel)
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
	return lv.model.GetItem(int(position))
}

// GetSelectedPosition returns the position of the selected item or -1 if no item is selected
func (lv *ListView) GetSelectedPosition() int {
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
	// Select each position
	for _, pos := range positions {
		if pos >= 0 && pos < lv.model.GetNItems() {
			C.list_view_multi_selection_select_item(lv.selectionModel, C.guint(pos), C.TRUE)
		}
	}
}

// UnselectItems unselects the items at the specified positions
func (lv *ListView) UnselectItems(positions []int) {
	// Unselect each position
	for _, pos := range positions {
		if pos >= 0 && pos < lv.model.GetNItems() {
			C.list_view_multi_selection_select_item(lv.selectionModel, C.guint(pos), C.FALSE)
		}
	}
}

// SelectRange selects a range of items
func (lv *ListView) SelectRange(start, count int) {
	// Check range is valid
	if start < 0 || count <= 0 || start+count > lv.model.GetNItems() {
		return
	}
	
	// Select the range
	C.list_view_multi_selection_select_range(lv.selectionModel, C.guint(start), C.guint(count), C.TRUE)
}

// UnselectRange unselects a range of items
func (lv *ListView) UnselectRange(start, count int) {
	// Check range is valid
	if start < 0 || count <= 0 || start+count > lv.model.GetNItems() {
		return
	}
	
	// Unselect the range
	C.list_view_multi_selection_select_range(lv.selectionModel, C.guint(start), C.guint(count), C.FALSE)
}

// SelectAll selects all items
func (lv *ListView) SelectAll() {
	C.list_view_multi_selection_select_all(lv.selectionModel)
}

// UnselectAll unselects all items
func (lv *ListView) UnselectAll() {
	C.list_view_multi_selection_unselect_all(lv.selectionModel)
}

// Destroy cleans up resources
func (lv *ListView) Destroy() {
	listViewCallbackMutex.Lock()
	defer listViewCallbackMutex.Unlock()

	// Remove callbacks
	delete(listViewCallbacks, uintptr(unsafe.Pointer(lv.widget)))
	
	if lv.selectionModel != nil {
		delete(selectionChangedCallbacks, uintptr(unsafe.Pointer(lv.selectionModel)))
		C.g_object_unref(C.gpointer(unsafe.Pointer(lv.selectionModel)))
		lv.selectionModel = nil
	}
	
	// Clean up factory if we created it
	if lv.factory != nil {
		lv.factory.Free()
		lv.factory = nil
	}

	// Call base destroy
	lv.BaseWidget.Destroy()
}