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
//
// // New in GTK 4.12: Scroll to API
// static void listViewScrollTo(GtkListView *list_view, guint position, GtkListScrollFlags flags) {
//     #if GTK_CHECK_VERSION(4, 12, 0)
//     #if GTK_CHECK_VERSION(4, 14, 0)
//         // In GTK 4.14+ the function requires a scroll_info parameter
//         GtkScrollInfo *scroll_info = NULL; // Default to NULL for automatic scrolling
//         gtk_list_view_scroll_to(list_view, position, flags, scroll_info);
//     #else
//         // In GTK 4.12-4.13, the function only takes position and flags
//         gtk_list_view_scroll_to(list_view, position, flags);
//     #endif
//     #endif
// }
//
// // New in GTK 4.12: Tab behavior API
// static void listViewSetTabBehavior(GtkListView *list_view, GtkListTabBehavior behavior) {
//     #if GTK_CHECK_VERSION(4, 12, 0)
//     gtk_list_view_set_tab_behavior(list_view, behavior);
//     #endif
// }
//
// static GtkListTabBehavior listViewGetTabBehavior(GtkListView *list_view) {
//     #if GTK_CHECK_VERSION(4, 12, 0)
//     return gtk_list_view_get_tab_behavior(list_view);
//     #else
//     return 0;
//     #endif
// }
//
// // New in GTK 4.12: Header factory API
// static void listViewSetHeaderFactory(GtkListView *list_view, GtkListItemFactory *factory) {
//     #if GTK_CHECK_VERSION(4, 12, 0)
//     gtk_list_view_set_header_factory(list_view, factory);
//     #endif
// }
//
// static GtkListItemFactory* listViewGetHeaderFactory(GtkListView *list_view) {
//     #if GTK_CHECK_VERSION(4, 12, 0)
//     return gtk_list_view_get_header_factory(list_view);
//     #else
//     return NULL;
//     #endif
// }
import "C"

import (
	"unsafe"
)

// ListScrollFlags represents the flags for scrolling to an item in the list
type ListScrollFlags int

const (
	// ListScrollNone indicates no special behavior
	ListScrollNone ListScrollFlags = 0
	// ListScrollFocus means focus the item when scrolling
	ListScrollFocus ListScrollFlags = 1 << 0
	// ListScrollSelect means select the item when scrolling
	ListScrollSelect ListScrollFlags = 1 << 1
)

// ListTabBehavior represents the behavior for tab navigation
type ListTabBehavior int

const (
	// ListTabAll allows tab to focus all items
	ListTabAll ListTabBehavior = 0
	// ListTabItem allows tab to focus only items
	ListTabItem ListTabBehavior = 1
	// ListTabRow allows tab to focus entire rows
	ListTabRow ListTabBehavior = 2
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
	headerFactory  ListItemFactory
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

// WithTabBehavior sets the tab behavior for the list view (GTK 4.12+)
func WithTabBehavior(behavior ListTabBehavior) ListViewOption {
	return func(lv *ListView) {
		C.listViewSetTabBehavior((*C.GtkListView)(unsafe.Pointer(lv.widget)), C.GtkListTabBehavior(behavior))
	}
}

// WithHeaderFactory sets the header factory for the list view (GTK 4.12+)
func WithHeaderFactory(factory ListItemFactory) ListViewOption {
	return func(lv *ListView) {
		if factory != nil {
			C.listViewSetHeaderFactory((*C.GtkListView)(unsafe.Pointer(lv.widget)), factory.GetListItemFactory())
			lv.headerFactory = factory
		}
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

// SetHeaderFactory sets the header factory for the list view (GTK 4.12+)
func (lv *ListView) SetHeaderFactory(factory ListItemFactory) {
	if factory != nil {
		C.listViewSetHeaderFactory((*C.GtkListView)(unsafe.Pointer(lv.widget)), factory.GetListItemFactory())
		lv.headerFactory = factory
	} else {
		C.listViewSetHeaderFactory((*C.GtkListView)(unsafe.Pointer(lv.widget)), nil)
		lv.headerFactory = nil
	}
}

// GetHeaderFactory returns the header factory for the list view (GTK 4.12+)
func (lv *ListView) GetHeaderFactory() ListItemFactory {
	return lv.headerFactory
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

// SetTabBehavior sets the tab behavior for the list view (GTK 4.12+)
func (lv *ListView) SetTabBehavior(behavior ListTabBehavior) {
	C.listViewSetTabBehavior((*C.GtkListView)(unsafe.Pointer(lv.widget)), C.GtkListTabBehavior(behavior))
}

// GetTabBehavior returns the tab behavior for the list view (GTK 4.12+)
func (lv *ListView) GetTabBehavior() ListTabBehavior {
	return ListTabBehavior(C.listViewGetTabBehavior((*C.GtkListView)(unsafe.Pointer(lv.widget))))
}

// ScrollTo scrolls to the item at the given position (GTK 4.12+)
func (lv *ListView) ScrollTo(position int, flags ListScrollFlags) {
	C.listViewScrollTo((*C.GtkListView)(unsafe.Pointer(lv.widget)), C.guint(position), C.GtkListScrollFlags(flags))
}

// ConnectActivate connects a callback for item activation
func (lv *ListView) ConnectActivate(callback ListViewActivateCallback) {
	if callback == nil {
		return
	}

	// To avoid type issues, convert the ListViewActivateCallback to a regular func(int)
	// since that's what the callback handler expects
	rawCallback := func(position int) {
		callback(position)
	}

	// Connect using the raw callback
	Connect(lv, SignalListActivate, rawCallback)

	DebugLog(DebugLevelInfo, DebugComponentListView, 
		"Connected activate callback to ListView %p", unsafe.Pointer(lv.widget))
}

// Destroy overrides BaseWidget's Destroy to clean up list view resources
func (lv *ListView) Destroy() {
	// Clean up callbacks using the unified system
	DisconnectAll(lv)
	
	lv.BaseWidget.Destroy()
}