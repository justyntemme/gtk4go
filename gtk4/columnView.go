// Package gtk4 provides ColumnView functionality for GTK4
// File: gtk4go/gtk4/columnView.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback functions for column view
// extern void columnViewActivateCallback(GtkColumnView *view, guint position, gpointer user_data);
//
// // Connect column view signals
// static gulong connectColumnViewActivate(GtkColumnView *view, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(view), "activate", G_CALLBACK(columnViewActivateCallback), user_data);
// }
//
// // ColumnView creation and configuration
// static GtkWidget* create_column_view(GtkSelectionModel *model) {
//     return gtk_column_view_new(model);
// }
//
// // ColumnView operations
// static void column_view_append_column(GtkColumnView *view, GtkColumnViewColumn *column) {
//     gtk_column_view_append_column(view, column);
// }
//
// static GListModel* column_view_get_columns(GtkColumnView *view) {
//     return gtk_column_view_get_columns(view);
// }
//
// static GtkColumnViewColumn* column_view_get_column_at_position(GtkColumnView *view, guint position) {
//     GListModel* columns = gtk_column_view_get_columns(view);
//     gpointer item = g_list_model_get_item(columns, position);
//     if (item == NULL) {
//         return NULL;
//     }
//     return GTK_COLUMN_VIEW_COLUMN(item);
// }
//
// static void column_view_remove_column(GtkColumnView *view, GtkColumnViewColumn *column) {
//     gtk_column_view_remove_column(view, column);
// }
//
// static void column_view_set_model(GtkColumnView *view, GtkSelectionModel *model) {
//     gtk_column_view_set_model(view, model);
// }
//
// static void column_view_set_show_row_separators(GtkColumnView *view, gboolean show_separators) {
//     gtk_column_view_set_show_row_separators(view, show_separators);
// }
//
// static void column_view_set_show_column_separators(GtkColumnView *view, gboolean show_separators) {
//     gtk_column_view_set_show_column_separators(view, show_separators);
// }
//
// static void column_view_set_reorderable(GtkColumnView *view, gboolean reorderable) {
//     gtk_column_view_set_reorderable(view, reorderable);
// }
//
// static void column_view_set_enable_rubberband(GtkColumnView *view, gboolean enable) {
//     gtk_column_view_set_enable_rubberband(view, enable);
// }
//
// // ColumnViewColumn creation and configuration
// static GtkColumnViewColumn* create_column_view_column(const char* title, GtkListItemFactory *factory) {
//     return gtk_column_view_column_new(title, factory);
// }
//
// static void column_view_column_set_title(GtkColumnViewColumn *column, const char* title) {
//     gtk_column_view_column_set_title(column, title);
// }
//
// static void column_view_column_set_factory(GtkColumnViewColumn *column, GtkListItemFactory *factory) {
//     gtk_column_view_column_set_factory(column, factory);
// }
//
// static void column_view_column_set_resizable(GtkColumnViewColumn *column, gboolean resizable) {
//     gtk_column_view_column_set_resizable(column, resizable);
// }
//
// static void column_view_column_set_expand(GtkColumnViewColumn *column, gboolean expand) {
//     gtk_column_view_column_set_expand(column, expand);
// }
//
// static void column_view_column_set_fixed_width(GtkColumnViewColumn *column, int width) {
//     gtk_column_view_column_set_fixed_width(column, width);
// }
//
// static void column_view_column_set_visible(GtkColumnViewColumn *column, gboolean visible) {
//     gtk_column_view_column_set_visible(column, visible);
// }
//
// static void column_view_column_set_header_menu(GtkColumnViewColumn *column, GMenuModel *menu) {
//     gtk_column_view_column_set_header_menu(column, menu);
// }
//
// // Selection model creation
// static GtkSelectionModel* create_single_selection(GListModel *model) {
//     return GTK_SELECTION_MODEL(gtk_single_selection_new(model));
// }
//
// static GtkSelectionModel* create_multi_selection(GListModel *model) {
//     return GTK_SELECTION_MODEL(gtk_multi_selection_new(model));
// }
//
// static GtkSelectionModel* create_no_selection(GListModel *model) {
//     return GTK_SELECTION_MODEL(gtk_no_selection_new(model));
// }
//
// // Sort handling
// typedef struct {
//     int column_id;
//     int direction; // 0 = ascending, 1 = descending
// } SortInfo;
//
// extern int column_view_sort_func(gpointer a, gpointer b, gpointer user_data);
//
// static void column_view_column_set_sorter(GtkColumnViewColumn *column, GtkSorter *sorter) {
//     gtk_column_view_column_set_sorter(column, sorter);
// }
//
// // Create a custom sorter function that uses the column ID
// static GtkSorter* create_custom_sorter(int direction) {
//     // Create sort info
//     SortInfo* sort_info = (SortInfo*)malloc(sizeof(SortInfo));
//     sort_info->direction = direction;
//
//     // Create custom sorter with our callback
//     GtkCustomSorter* sorter = gtk_custom_sorter_new((GCompareDataFunc)column_view_sort_func, sort_info, free);
//     return GTK_SORTER(sorter);
// }
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// ColumnViewActivatedCallback represents a callback for column view activated events
type ColumnViewActivatedCallback func(position int)

var (
	columnViewCallbacks     = make(map[uintptr]ColumnViewActivatedCallback)
	columnViewCallbackMutex sync.RWMutex
)

//export columnViewActivateCallback
func columnViewActivateCallback(view *C.GtkColumnView, position C.guint, userData C.gpointer) {
	columnViewCallbackMutex.RLock()
	defer columnViewCallbackMutex.RUnlock()

	// Convert view pointer to uintptr for lookup
	viewPtr := uintptr(unsafe.Pointer(view))

	// Find and call the callback
	if callback, ok := columnViewCallbacks[viewPtr]; ok {
		callback(int(position))
	}
}

//export column_view_sort_func
func column_view_sort_func(a, b, userData C.gpointer) C.int {
	// Extract sort info
	sortInfo := (*C.SortInfo)(userData)
	direction := int(sortInfo.direction)

	// In a real implementation, we would extract values from a and b
	// and compare them. This is a simplified version.
	if direction == 0 {
		// Ascending
		return C.int(0)
	} else {
		// Descending
		return C.int(0)
	}
}

// ColumnViewOption is a function that configures a column view
type ColumnViewOption func(*ColumnView)

// ColumnView represents a GTK column view widget
type ColumnView struct {
	BaseWidget
	selectionModel *C.GtkSelectionModel
	model          ListModel
	columns        []*ColumnViewColumn
}

// NewColumnView creates a new GTK column view widget with the given model
func NewColumnView(model ListModel, options ...ColumnViewOption) *ColumnView {
	// Create a column view with a default configuration
	columnView := &ColumnView{
		model:   model,
		columns: make([]*ColumnViewColumn, 0),
	}

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
	columnView.selectionModel = C.create_single_selection(glistModel)

	// Create the widget
	columnView.widget = C.create_column_view(columnView.selectionModel)

	// Apply options
	for _, option := range options {
		option(columnView)
	}

	SetupFinalization(columnView, columnView.Destroy)
	return columnView
}

// WithColumnSelectionMode sets the selection mode for the column view

func WithColumnSelectionMode(mode SelectionMode) ColumnViewOption {
	return func(cv *ColumnView) {
		// Get the current GListModel from the selection model
		var glistModel *C.GListModel

		// This is a simplified implementation
		// In a real implementation, you would extract the model properly
		if cv.model != nil {
			if glm, ok := cv.model.(*GListModel); ok {
				glistModel = glm.model
			} else if ls, ok := cv.model.(*ListStore); ok {
				glistModel = (*C.GListModel)(unsafe.Pointer(ls.store))
			}
		}

		// If no model available, do nothing
		if glistModel == nil {
			return
		}

		// Create a new selection model before cleaning up the old one
		var newSelectionModel *C.GtkSelectionModel

		// Create a new selection model based on the mode
		switch mode {
		case SelectionModeSingle:
			newSelectionModel = C.create_single_selection(glistModel)
		case SelectionModeMultiple:
			newSelectionModel = C.create_multi_selection(glistModel)
		case SelectionModeNone:
			newSelectionModel = C.create_no_selection(glistModel)
		default:
			// Default to single selection if mode is invalid
			newSelectionModel = C.create_single_selection(glistModel)
		}

		// Check if new model was created successfully
		if newSelectionModel == nil {
			return
		}

		// Clean up the old selection model AFTER creating the new one
		if cv.selectionModel != nil {
			C.g_object_unref(C.gpointer(unsafe.Pointer(cv.selectionModel)))
		}

		// Store the new selection model
		cv.selectionModel = newSelectionModel

		// Only set the model on the column view if we have a valid widget
		if cv.widget != nil {
			// This is the line that's crashing - ensure we have valid pointers
			C.column_view_set_model(
				(*C.GtkColumnView)(unsafe.Pointer(cv.widget)),
				cv.selectionModel)
		}
	}
}

// WithShowRowSeparators sets whether to show row separators
func WithShowRowSeparators(show bool) ColumnViewOption {
	return func(cv *ColumnView) {
		var cShow C.gboolean
		if show {
			cShow = C.TRUE
		} else {
			cShow = C.FALSE
		}
		C.column_view_set_show_row_separators((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), cShow)
	}
}

// WithShowColumnSeparators sets whether to show column separators
func WithShowColumnSeparators(show bool) ColumnViewOption {
	return func(cv *ColumnView) {
		var cShow C.gboolean
		if show {
			cShow = C.TRUE
		} else {
			cShow = C.FALSE
		}
		C.column_view_set_show_column_separators((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), cShow)
	}
}

// WithReorderable sets whether columns can be reordered
func WithReorderable(reorderable bool) ColumnViewOption {
	return func(cv *ColumnView) {
		var cReorderable C.gboolean
		if reorderable {
			cReorderable = C.TRUE
		} else {
			cReorderable = C.FALSE
		}
		C.column_view_set_reorderable((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), cReorderable)
	}
}

// WithColumnRubberband sets whether rubberband selection is enabled
func WithColumnRubberband(enable bool) ColumnViewOption {
	return func(cv *ColumnView) {
		var cEnable C.gboolean
		if enable {
			cEnable = C.TRUE
		} else {
			cEnable = C.FALSE
		}
		C.column_view_set_enable_rubberband((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), cEnable)
	}
}

// ConnectActivate connects a callback to the activate signal
func (cv *ColumnView) ConnectActivate(callback ColumnViewActivatedCallback) {
	columnViewCallbackMutex.Lock()
	defer columnViewCallbackMutex.Unlock()

	// Store callback in map
	viewPtr := uintptr(unsafe.Pointer(cv.widget))
	columnViewCallbacks[viewPtr] = callback

	// Connect signal
	C.connectColumnViewActivate((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), C.gpointer(unsafe.Pointer(cv.widget)))
}

// AppendColumn appends a column to the column view
func (cv *ColumnView) AppendColumn(column *ColumnViewColumn) {
	C.column_view_append_column(
		(*C.GtkColumnView)(unsafe.Pointer(cv.widget)),
		column.column,
	)

	// Store the column
	cv.columns = append(cv.columns, column)
}

// RemoveColumn removes a column from the column view
func (cv *ColumnView) RemoveColumn(column *ColumnViewColumn) {
	C.column_view_remove_column(
		(*C.GtkColumnView)(unsafe.Pointer(cv.widget)),
		column.column,
	)

	// Remove from our slice
	for i, col := range cv.columns {
		if col == column {
			cv.columns = append(cv.columns[:i], cv.columns[i+1:]...)
			break
		}
	}
}

// GetColumn gets the column at the specified position
func (cv *ColumnView) GetColumn(position int) *ColumnViewColumn {
	if position < 0 || position >= len(cv.columns) {
		return nil
	}

	return cv.columns[position]
}

// GetColumnCount gets the number of columns
func (cv *ColumnView) GetColumnCount() int {
	return len(cv.columns)
}

// GetColumnByPosition gets a column at the specified position directly from GTK
func (cv *ColumnView) GetColumnByPosition(position int) *ColumnViewColumn {
	column := C.column_view_get_column_at_position(
		(*C.GtkColumnView)(unsafe.Pointer(cv.widget)),
		C.guint(position),
	)

	if column == nil {
		return nil
	}

	// Create a wrapper
	return &ColumnViewColumn{
		column: column,
	}
}

// GetModel returns the model used by the column view
func (cv *ColumnView) GetModel() ListModel {
	return cv.model
}

// SetModel sets the model for the column view
func (cv *ColumnView) SetModel(model ListModel) {
	cv.model = model

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
	if cv.selectionModel != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(cv.selectionModel)))
	}

	// Create a new selection model (default to single selection)
	cv.selectionModel = C.create_single_selection(glistModel)

	// Update the column view
	C.column_view_set_model((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), cv.selectionModel)
}

// SetShowRowSeparators sets whether to show row separators
func (cv *ColumnView) SetShowRowSeparators(show bool) {
	var cShow C.gboolean
	if show {
		cShow = C.TRUE
	} else {
		cShow = C.FALSE
	}
	C.column_view_set_show_row_separators((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), cShow)
}

// SetShowColumnSeparators sets whether to show column separators
func (cv *ColumnView) SetShowColumnSeparators(show bool) {
	var cShow C.gboolean
	if show {
		cShow = C.TRUE
	} else {
		cShow = C.FALSE
	}
	C.column_view_set_show_column_separators((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), cShow)
}

// SetReorderable sets whether columns can be reordered
func (cv *ColumnView) SetReorderable(reorderable bool) {
	var cReorderable C.gboolean
	if reorderable {
		cReorderable = C.TRUE
	} else {
		cReorderable = C.FALSE
	}
	C.column_view_set_reorderable((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), cReorderable)
}

// SetRubberbandSelection sets whether rubberband selection is enabled
func (cv *ColumnView) SetRubberbandSelection(enable bool) {
	var cEnable C.gboolean
	if enable {
		cEnable = C.TRUE
	} else {
		cEnable = C.FALSE
	}
	C.column_view_set_enable_rubberband((*C.GtkColumnView)(unsafe.Pointer(cv.widget)), cEnable)
}

// Destroy cleans up resources
func (cv *ColumnView) Destroy() {
	columnViewCallbackMutex.Lock()
	defer columnViewCallbackMutex.Unlock()

	// Remove callbacks
	delete(columnViewCallbacks, uintptr(unsafe.Pointer(cv.widget)))

	if cv.selectionModel != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(cv.selectionModel)))
		cv.selectionModel = nil
	}

	// Release columns
	for _, column := range cv.columns {
		column.Free()
	}
	cv.columns = nil

	// Call base destroy
	cv.BaseWidget.Destroy()
}

// ColumnViewColumnOption is a function that configures a column view column
type ColumnViewColumnOption func(*ColumnViewColumn)

// ColumnViewColumn represents a GTK column view column
type ColumnViewColumn struct {
	column  *C.GtkColumnViewColumn
	factory *ListItemFactory
}

// NewColumnViewColumn creates a new GTK column view column
func NewColumnViewColumn(title string, factory *ListItemFactory, options ...ColumnViewColumnOption) *ColumnViewColumn {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	column := &ColumnViewColumn{
		column:  C.create_column_view_column(cTitle, factory.factory),
		factory: factory,
	}

	// Apply options
	for _, option := range options {
		option(column)
	}

	runtime.SetFinalizer(column, (*ColumnViewColumn).Free)
	return column
}

// WithResizable sets whether the column is resizable
func WithResizable(resizable bool) ColumnViewColumnOption {
	return func(col *ColumnViewColumn) {
		var cResizable C.gboolean
		if resizable {
			cResizable = C.TRUE
		} else {
			cResizable = C.FALSE
		}
		C.column_view_column_set_resizable(col.column, cResizable)
	}
}

// WithExpand sets whether the column should expand
func WithExpand(expand bool) ColumnViewColumnOption {
	return func(col *ColumnViewColumn) {
		var cExpand C.gboolean
		if expand {
			cExpand = C.TRUE
		} else {
			cExpand = C.FALSE
		}
		C.column_view_column_set_expand(col.column, cExpand)
	}
}

// WithFixedWidth sets a fixed width for the column
func WithFixedWidth(width int) ColumnViewColumnOption {
	return func(col *ColumnViewColumn) {
		C.column_view_column_set_fixed_width(col.column, C.int(width))
	}
}

// WithVisible sets whether the column is visible
func WithVisible(visible bool) ColumnViewColumnOption {
	return func(col *ColumnViewColumn) {
		var cVisible C.gboolean
		if visible {
			cVisible = C.TRUE
		} else {
			cVisible = C.FALSE
		}
		C.column_view_column_set_visible(col.column, cVisible)
	}
}

// WithSorter sets a sorter for the column
func WithSorter(ascending bool) ColumnViewColumnOption {
	return func(col *ColumnViewColumn) {
		direction := 0
		if !ascending {
			direction = 1
		}

		sorter := C.create_custom_sorter(C.int(direction))
		C.column_view_column_set_sorter(col.column, sorter)
	}
}

// SetTitle sets the title of the column
func (col *ColumnViewColumn) SetTitle(title string) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.column_view_column_set_title(col.column, cTitle)
}

// SetFactory sets the factory for the column
func (col *ColumnViewColumn) SetFactory(factory *ListItemFactory) {
	// Store the old factory for cleanup
	oldFactory := col.factory

	// Set the new factory
	col.factory = factory
	C.column_view_column_set_factory(col.column, factory.factory)

	// Clean up the old factory if different
	if oldFactory != nil && oldFactory != factory {
		oldFactory.Free()
	}
}

// SetResizable sets whether the column is resizable
func (col *ColumnViewColumn) SetResizable(resizable bool) {
	var cResizable C.gboolean
	if resizable {
		cResizable = C.TRUE
	} else {
		cResizable = C.FALSE
	}
	C.column_view_column_set_resizable(col.column, cResizable)
}

// SetExpand sets whether the column should expand
func (col *ColumnViewColumn) SetExpand(expand bool) {
	var cExpand C.gboolean
	if expand {
		cExpand = C.TRUE
	} else {
		cExpand = C.FALSE
	}
	C.column_view_column_set_expand(col.column, cExpand)
}

// SetFixedWidth sets a fixed width for the column
func (col *ColumnViewColumn) SetFixedWidth(width int) {
	C.column_view_column_set_fixed_width(col.column, C.int(width))
}

// SetVisible sets whether the column is visible
func (col *ColumnViewColumn) SetVisible(visible bool) {
	var cVisible C.gboolean
	if visible {
		cVisible = C.TRUE
	} else {
		cVisible = C.FALSE
	}
	C.column_view_column_set_visible(col.column, cVisible)
}

// SetSorter sets a sorter for the column
func (col *ColumnViewColumn) SetSorter(ascending bool) {
	direction := 0
	if !ascending {
		direction = 1
	}

	sorter := C.create_custom_sorter(C.int(direction))
	C.column_view_column_set_sorter(col.column, sorter)
}

// Free frees the column
func (col *ColumnViewColumn) Free() {
	if col.column != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(col.column)))
		col.column = nil
	}

	// Clean up factory
	if col.factory != nil {
		col.factory.Free()
		col.factory = nil
	}
}

// Helper functions for creating common column types

// TextColumn creates a column that displays text values
func TextColumn(title string, columnID int, options ...ColumnViewColumnOption) *ColumnViewColumn {
	// Create a factory for rendering text
	factory := TextFactory()

	// Create the column
	column := NewColumnViewColumn(title, factory, options...)

	// Add a sorter
	column.SetSorter(true)

	return column
}

// CheckboxColumn creates a column that displays checkbox values
func CheckboxColumn(title string, columnID int, options ...ColumnViewColumnOption) *ColumnViewColumn {
	// Create a factory for rendering checkboxes
	factory := CheckboxFactory()

	// Create the column
	column := NewColumnViewColumn(title, factory, options...)

	// Add a sorter
	column.SetSorter(true)

	return column
}

// ProgressColumn creates a column that displays progress bars
func ProgressColumn(title string, columnID int, options ...ColumnViewColumnOption) *ColumnViewColumn {
	// Create a factory for rendering progress bars
	factory := ProgressFactory()

	// Create the column
	column := NewColumnViewColumn(title, factory, options...)

	// Add a sorter
	column.SetSorter(true)

	return column
}

// CustomColumn creates a column with a custom factory
func CustomColumn(title string, factory *ListItemFactory, columnID int, options ...ColumnViewColumnOption) *ColumnViewColumn {
	// Create the column
	column := NewColumnViewColumn(title, factory, options...)

	// Add a sorter
	column.SetSorter(true)

	return column
}
