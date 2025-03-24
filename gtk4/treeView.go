// Package gtk4 provides tree view functionality for GTK4
// File: gtk4go/gtk4/treeView.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback function for tree view selection changed
// extern void treeViewSelectionChangedCallback(GtkTreeSelection *selection, gpointer user_data);
//
// // Connect tree view selection changed signal with callback
// static gulong connectTreeViewSelectionChanged(GtkTreeSelection *selection, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(selection), "changed", G_CALLBACK(treeViewSelectionChangedCallback), user_data);
// }
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// TreeViewSelectionChangedCallback represents a callback for tree view selection changed events
type TreeViewSelectionChangedCallback func()

var (
	treeViewSelectionCallbacks = make(map[uintptr]TreeViewSelectionChangedCallback)
	treeViewCallbackMutex      sync.RWMutex
)

//export treeViewSelectionChangedCallback
func treeViewSelectionChangedCallback(selection *C.GtkTreeSelection, userData C.gpointer) {
	treeViewCallbackMutex.RLock()
	defer treeViewCallbackMutex.RUnlock()

	// Convert selection pointer to uintptr for lookup
	selectionPtr := uintptr(unsafe.Pointer(selection))

	// Find and call the callback
	if callback, ok := treeViewSelectionCallbacks[selectionPtr]; ok {
		callback()
	}
}

// TreeViewOption is a function that configures a tree view
type TreeViewOption func(*TreeView)

// TreeView represents a GTK tree view widget
type TreeView struct {
	BaseWidget
	selection *TreeSelection
}

// NewTreeView creates a new GTK tree view widget
func NewTreeView(model TreeModel, options ...TreeViewOption) *TreeView {
	var widget *C.GtkWidget

	if model != nil {
		widget = C.gtk_tree_view_new_with_model(model.GetModelPointer())
	} else {
		widget = C.gtk_tree_view_new()
	}

	tv := &TreeView{
		BaseWidget: BaseWidget{
			widget: widget,
		},
	}

	// Create selection
	sel := C.gtk_tree_view_get_selection((*C.GtkTreeView)(unsafe.Pointer(tv.widget)))
	tv.selection = &TreeSelection{
		selection: sel,
	}

	// Apply options
	for _, option := range options {
		option(tv)
	}

	SetupFinalization(tv, tv.Destroy)
	return tv
}

// WithHeaders configures whether headers are visible
func WithHeaders(visible bool) TreeViewOption {
	return func(tv *TreeView) {
		var cVisible C.gboolean
		if visible {
			cVisible = C.TRUE
		} else {
			cVisible = C.FALSE
		}
		C.gtk_tree_view_set_headers_visible((*C.GtkTreeView)(unsafe.Pointer(tv.widget)), cVisible)
	}
}

// WithMultipleSelection enables multiple selection
func WithMultipleSelection() TreeViewOption {
	return func(tv *TreeView) {
		C.gtk_tree_selection_set_mode(tv.selection.selection, C.GTK_SELECTION_MULTIPLE)
	}
}

// GetSelection gets the tree selection
func (tv *TreeView) GetSelection() *TreeSelection {
	return tv.selection
}

// AppendColumn appends a column to the tree view
func (tv *TreeView) AppendColumn(column *TreeViewColumn) int {
	return int(C.gtk_tree_view_append_column((*C.GtkTreeView)(unsafe.Pointer(tv.widget)), column.column))
}

// SetModel sets the model for the tree view
func (tv *TreeView) SetModel(model TreeModel) {
	if model != nil {
		C.gtk_tree_view_set_model((*C.GtkTreeView)(unsafe.Pointer(tv.widget)), model.GetModelPointer())
	} else {
		C.gtk_tree_view_set_model((*C.GtkTreeView)(unsafe.Pointer(tv.widget)), nil)
	}
}

// GetModel gets the model for the tree view
func (tv *TreeView) GetModel() TreeModel {
	modelPtr := C.gtk_tree_view_get_model((*C.GtkTreeView)(unsafe.Pointer(tv.widget)))
	if modelPtr == nil {
		return nil
	}

	// Create a TreeModelImplementor wrapper
	return &TreeModelImplementor{
		model: modelPtr,
	}
}

// SetHeadersVisible sets whether headers are visible
func (tv *TreeView) SetHeadersVisible(visible bool) {
	var cVisible C.gboolean
	if visible {
		cVisible = C.TRUE
	} else {
		cVisible = C.FALSE
	}
	C.gtk_tree_view_set_headers_visible((*C.GtkTreeView)(unsafe.Pointer(tv.widget)), cVisible)
}

// ExpandAll expands all rows in the tree view
func (tv *TreeView) ExpandAll() {
	C.gtk_tree_view_expand_all((*C.GtkTreeView)(unsafe.Pointer(tv.widget)))
}

// CollapseAll collapses all rows in the tree view
func (tv *TreeView) CollapseAll() {
	C.gtk_tree_view_collapse_all((*C.GtkTreeView)(unsafe.Pointer(tv.widget)))
}

// Destroy overrides the BaseWidget Destroy to clean up resources
func (tv *TreeView) Destroy() {
	// Clean up selection callback if registered
	if tv.selection != nil {
		tv.selection.DisconnectChanged()
		tv.selection = nil
	}

	// Call base Destroy
	tv.BaseWidget.Destroy()
}

// TreeSelection represents a GTK tree selection
type TreeSelection struct {
	selection *C.GtkTreeSelection
}

// SelectionMode defines the type of selection
type SelectionMode int

const (
	// SelectionNone no selection is possible
	SelectionNone SelectionMode = C.GTK_SELECTION_NONE
	// SelectionSingle only one item can be selected
	SelectionSingle SelectionMode = C.GTK_SELECTION_SINGLE
	// SelectionBrowse browse mode, only one item can be selected, changes as the pointer moves
	SelectionBrowse SelectionMode = C.GTK_SELECTION_BROWSE
	// SelectionMultiple multiple items can be selected
	SelectionMultiple SelectionMode = C.GTK_SELECTION_MULTIPLE
)

// SetMode sets the selection mode
func (ts *TreeSelection) SetMode(mode SelectionMode) {
	C.gtk_tree_selection_set_mode(ts.selection, C.GtkSelectionMode(mode))
}

// GetMode gets the selection mode
func (ts *TreeSelection) GetMode() SelectionMode {
	return SelectionMode(C.gtk_tree_selection_get_mode(ts.selection))
}

// GetSelected gets the selected row
func (ts *TreeSelection) GetSelected() (TreeModel, *TreeIter, bool) {
	var model *C.GtkTreeModel
	iter := &TreeIter{}

	result := C.gtk_tree_selection_get_selected(ts.selection, &model, &iter.iter)
	if result == C.FALSE {
		return nil, nil, false
	}

	// Create a TreeModel wrapper for the model
	treeModel := &TreeModelImplementor{
		model: model,
	}

	return treeModel, iter, true
}

// SelectIter selects a row by iterator
func (ts *TreeSelection) SelectIter(iter *TreeIter) {
	C.gtk_tree_selection_select_iter(ts.selection, &iter.iter)
}

// UnselectIter unselects a row by iterator
func (ts *TreeSelection) UnselectIter(iter *TreeIter) {
	C.gtk_tree_selection_unselect_iter(ts.selection, &iter.iter)
}

// ConnectChanged connects a callback function to the selection's "changed" signal
func (ts *TreeSelection) ConnectChanged(callback TreeViewSelectionChangedCallback) {
	treeViewCallbackMutex.Lock()
	defer treeViewCallbackMutex.Unlock()

	// Store callback in map
	selectionPtr := uintptr(unsafe.Pointer(ts.selection))
	treeViewSelectionCallbacks[selectionPtr] = callback

	// Connect signal
	C.connectTreeViewSelectionChanged(ts.selection, C.gpointer(unsafe.Pointer(ts.selection)))
}

// DisconnectChanged disconnects the changed signal handler
func (ts *TreeSelection) DisconnectChanged() {
	treeViewCallbackMutex.Lock()
	defer treeViewCallbackMutex.Unlock()

	// Remove callback from map
	selectionPtr := uintptr(unsafe.Pointer(ts.selection))
	delete(treeViewSelectionCallbacks, selectionPtr)
}

// TreeViewColumn represents a column in a GTK tree view
type TreeViewColumn struct {
	column *C.GtkTreeViewColumn
}

// NewTreeViewColumn creates a new tree view column
func NewTreeViewColumn(title string, renderer CellRenderer, attributes ...ColumnAttribute) *TreeViewColumn {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	column := &TreeViewColumn{
		column: C.gtk_tree_view_column_new(),
	}

	C.gtk_tree_view_column_set_title(column.column, cTitle)

	if renderer != nil {
		C.gtk_tree_view_column_pack_start(column.column, renderer.GetCellRenderer(), C.TRUE)

		// Add attributes
		for _, attr := range attributes {
			cProperty := C.CString(attr.Property)
			defer C.free(unsafe.Pointer(cProperty))
			C.gtk_tree_view_column_add_attribute(column.column, renderer.GetCellRenderer(), cProperty, C.gint(attr.Column))
		}
	}

	runtime.SetFinalizer(column, (*TreeViewColumn).Free)
	return column
}

// SetTitle sets the column title
func (c *TreeViewColumn) SetTitle(title string) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.gtk_tree_view_column_set_title(c.column, cTitle)
}

// SetResizable sets whether the column is resizable
func (c *TreeViewColumn) SetResizable(resizable bool) {
	var cResizable C.gboolean
	if resizable {
		cResizable = C.TRUE
	} else {
		cResizable = C.FALSE
	}
	C.gtk_tree_view_column_set_resizable(c.column, cResizable)
}

// SetReorderable sets whether the column is reorderable
func (c *TreeViewColumn) SetReorderable(reorderable bool) {
	var cReorderable C.gboolean
	if reorderable {
		cReorderable = C.TRUE
	} else {
		cReorderable = C.FALSE
	}
	C.gtk_tree_view_column_set_reorderable(c.column, cReorderable)
}

// SetSortColumnID sets the column to sort by
func (c *TreeViewColumn) SetSortColumnID(sortColumnID int) {
	C.gtk_tree_view_column_set_sort_column_id(c.column, C.gint(sortColumnID))
}

// SetClickable sets whether the column is clickable
func (c *TreeViewColumn) SetClickable(clickable bool) {
	var cClickable C.gboolean
	if clickable {
		cClickable = C.TRUE
	} else {
		cClickable = C.FALSE
	}
	C.gtk_tree_view_column_set_clickable(c.column, cClickable)
}

// Free frees the column
func (c *TreeViewColumn) Free() {
	if c.column != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(c.column)))
		c.column = nil
	}
}

// ColumnAttribute associates a cell renderer property with a model column
type ColumnAttribute struct {
	Property string
	Column   int
}

// Attr creates a column attribute
func Attr(property string, column int) ColumnAttribute {
	return ColumnAttribute{
		Property: property,
		Column:   column,
	}
}

