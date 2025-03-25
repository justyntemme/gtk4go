// Package gtk4 provides list view functionality for GTK4
// File: gtk4go/gtk4/listview.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // ListView operations
// static GtkWidget* createListView(GtkSelectionModel *model, GtkListItemFactory *factory) {
//     return gtk_list_view_new(model, factory);
// }
//
// static void listViewSetModel(GtkListView *list_view, GtkSelectionModel *model) {
//     gtk_list_view_set_model(list_view, model);
// }
//
// static GtkSelectionModel* listViewGetModel(GtkListView *list_view) {
//     return gtk_list_view_get_model(list_view);
// }
//
// static void listViewSetFactory(GtkListView *list_view, GtkListItemFactory *factory) {
//     gtk_list_view_set_factory(list_view, factory);
// }
//
// static GtkListItemFactory* listViewGetFactory(GtkListView *list_view) {
//     return gtk_list_view_get_factory(list_view);
// }
//
// static void listViewSetShowSeparators(GtkListView *list_view, gboolean show_separators) {
//     gtk_list_view_set_show_separators(list_view, show_separators);
// }
//
// static gboolean listViewGetShowSeparators(GtkListView *list_view) {
//     return gtk_list_view_get_show_separators(list_view);
// }
//
// static void listViewSetSingleClickActivate(GtkListView *list_view, gboolean single_click_activate) {
//     gtk_list_view_set_single_click_activate(list_view, single_click_activate);
// }
//
// static gboolean listViewGetSingleClickActivate(GtkListView *list_view) {
//     return gtk_list_view_get_single_click_activate(list_view);
// }
//
// static void listViewSetEnableRubberband(GtkListView *list_view, gboolean enable_rubberband) {
//     gtk_list_view_set_enable_rubberband(list_view, enable_rubberband);
// }
//
// static gboolean listViewGetEnableRubberband(GtkListView *list_view) {
//     return gtk_list_view_get_enable_rubberband(list_view);
// }
import "C"

import (
	"unsafe"
)

// ListViewActivateCallback represents a callback for list view item activation
type ListViewActivateCallback func(position int)

// ListViewOption is a function that configures a list view
type ListViewOption func(*ListView)

// ListView represents a GTK list view widget
type ListView struct {
	BaseWidget
	selectionModel SelectionModel
	factory        ListItemFactory
}

// NewListView creates a new GTK list view
func NewListView(selectionModel SelectionModel, factory ListItemFactory, options ...ListViewOption) *ListView {
	var widget *C.GtkWidget

	// Create list view with selection model and factory if provided
	if selectionModel != nil && factory != nil {
		widget = C.createListView(selectionModel.GetSelectionModel(), factory.GetListItemFactory())
	} else {
		widget = C.createListView(nil, nil)
	}

	listView := &ListView{
		BaseWidget: BaseWidget{
			widget: widget,
		},
		selectionModel: selectionModel,
		factory:        factory,
	}

	// Apply options
	for _, option := range options {
		option(listView)
	}

	SetupFinalization(listView, listView.Destroy)
	return listView
}

// WithShowSeparators sets whether to show separators between items
func WithShowSeparators(showSeparators bool) ListViewOption {
	return func(lv *ListView) {
		var cshowSeparators C.gboolean
		if showSeparators {
			cshowSeparators = C.TRUE
		} else {
			cshowSeparators = C.FALSE
		}
		C.listViewSetShowSeparators((*C.GtkListView)(unsafe.Pointer(lv.widget)), cshowSeparators)
	}
}

// WithSingleClickActivate sets whether items are activated on single click
func WithSingleClickActivate(singleClickActivate bool) ListViewOption {
	return func(lv *ListView) {
		var csingleClickActivate C.gboolean
		if singleClickActivate {
			csingleClickActivate = C.TRUE
		} else {
			csingleClickActivate = C.FALSE
		}
		C.listViewSetSingleClickActivate((*C.GtkListView)(unsafe.Pointer(lv.widget)), csingleClickActivate)
	}
}

// WithEnableRubberband sets whether to enable rubberband selection
func WithEnableRubberband(enableRubberband bool) ListViewOption {
	return func(lv *ListView) {
		var cenableRubberband C.gboolean
		if enableRubberband {
			cenableRubberband = C.TRUE
		} else {
			cenableRubberband = C.FALSE
		}
		C.listViewSetEnableRubberband((*C.GtkListView)(unsafe.Pointer(lv.widget)), cenableRubberband)
	}
}

// SetModel sets the selection model for the list view
func (lv *ListView) SetModel(model SelectionModel) {
	if model != nil {
		C.listViewSetModel((*C.GtkListView)(unsafe.Pointer(lv.widget)), model.GetSelectionModel())
		lv.selectionModel = model
	} else {
		C.listViewSetModel((*C.GtkListView)(unsafe.Pointer(lv.widget)), nil)
		lv.selectionModel = nil
	}
}

// GetModel returns the selection model for the list view
func (lv *ListView) GetModel() SelectionModel {
	return lv.selectionModel
}

// SetFactory sets the list item factory for the list view
func (lv *ListView) SetFactory(factory ListItemFactory) {
	if factory != nil {
		C.listViewSetFactory((*C.GtkListView)(unsafe.Pointer(lv.widget)), factory.GetListItemFactory())
		lv.factory = factory
	} else {
		C.listViewSetFactory((*C.GtkListView)(unsafe.Pointer(lv.widget)), nil)
		lv.factory = nil
	}
}

// GetFactory returns the list item factory for the list view
func (lv *ListView) GetFactory() ListItemFactory {
	return lv.factory
}

// SetShowSeparators sets whether to show separators between items
func (lv *ListView) SetShowSeparators(showSeparators bool) {
	var cshowSeparators C.gboolean
	if showSeparators {
		cshowSeparators = C.TRUE
	} else {
		cshowSeparators = C.FALSE
	}
	C.listViewSetShowSeparators((*C.GtkListView)(unsafe.Pointer(lv.widget)), cshowSeparators)
}

// GetShowSeparators returns whether separators are shown between items
func (lv *ListView) GetShowSeparators() bool {
	return C.listViewGetShowSeparators((*C.GtkListView)(unsafe.Pointer(lv.widget))) != 0
}

// SetSingleClickActivate sets whether items are activated on single click
func (lv *ListView) SetSingleClickActivate(singleClickActivate bool) {
	var csingleClickActivate C.gboolean
	if singleClickActivate {
		csingleClickActivate = C.TRUE
	} else {
		csingleClickActivate = C.FALSE
	}
	C.listViewSetSingleClickActivate((*C.GtkListView)(unsafe.Pointer(lv.widget)), csingleClickActivate)
}

// GetSingleClickActivate returns whether items are activated on single click
func (lv *ListView) GetSingleClickActivate() bool {
	return C.listViewGetSingleClickActivate((*C.GtkListView)(unsafe.Pointer(lv.widget))) != 0
}

// SetEnableRubberband sets whether to enable rubberband selection
func (lv *ListView) SetEnableRubberband(enableRubberband bool) {
	var cenableRubberband C.gboolean
	if enableRubberband {
		cenableRubberband = C.TRUE
	} else {
		cenableRubberband = C.FALSE
	}
	C.listViewSetEnableRubberband((*C.GtkListView)(unsafe.Pointer(lv.widget)), cenableRubberband)
}

// GetEnableRubberband returns whether rubberband selection is enabled
func (lv *ListView) GetEnableRubberband() bool {
	return C.listViewGetEnableRubberband((*C.GtkListView)(unsafe.Pointer(lv.widget))) != 0
}

// ConnectActivate connects a callback for item activation
// This needs to use a wrapper to correctly handle the integer parameter
func (lv *ListView) ConnectActivate(callback ListViewActivateCallback) {
	if callback == nil {
		return
	}

	// We need to use a wrapper function to properly adapt the parameter type
	Connect(lv, SignalListActivate, func(param interface{}) {
		// The callback system in callbacks.go passes parameters as interface{}
		// We need to convert it to the correct type (int) for our callback
		if positionVal, ok := param.(int); ok {
			// Log successful connection
			DebugLog(DebugLevelVerbose, DebugComponentListView, 
				"ListView activate callback triggered for position: %d", positionVal)
			callback(positionVal)
		} else {
			// Handle potential type mismatches and debug info
			DebugLog(DebugLevelError, DebugComponentListView, 
				"Type mismatch on position parameter: expected int, got %T", param)
			// Default to 0 if the parameter type is not what we expect
			callback(0)
		}
	})
}

// Destroy overrides BaseWidget's Destroy to clean up list view resources
func (lv *ListView) Destroy() {
	// Clean up callbacks using the unified system
	DisconnectAll(lv)
	
	lv.BaseWidget.Destroy()
}