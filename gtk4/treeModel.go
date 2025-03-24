// Package gtk4 provides tree model functionality for GTK4
// File: gtk4go/gtk4/treeModel.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Wrapper for G_VALUE_TYPE macro
// static GType get_value_type(GValue *value) {
//     return G_VALUE_TYPE(value);
// }
//
// // Wrapper for g_value_get_string to handle NULL
// static const char* safe_get_string(GValue *value) {
//     const char* str = g_value_get_string(value);
//     return str ? str : "";
// }
//
// // Wrapper functions for TreeStore set value (to avoid variadic function issues)
// static void tree_store_set_string(GtkTreeStore *store, GtkTreeIter *iter, gint column, const char *value) {
//     gtk_tree_store_set(store, iter, column, value, -1);
// }
//
// static void tree_store_set_int(GtkTreeStore *store, GtkTreeIter *iter, gint column, gint value) {
//     gtk_tree_store_set(store, iter, column, value, -1);
// }
//
// static void tree_store_set_boolean(GtkTreeStore *store, GtkTreeIter *iter, gint column, gboolean value) {
//     gtk_tree_store_set(store, iter, column, value, -1);
// }
//
// static void tree_store_set_float(GtkTreeStore *store, GtkTreeIter *iter, gint column, gfloat value) {
//     gtk_tree_store_set(store, iter, column, value, -1);
// }
//
// static void tree_store_set_double(GtkTreeStore *store, GtkTreeIter *iter, gint column, gdouble value) {
//     gtk_tree_store_set(store, iter, column, value, -1);
// }
//
// // Wrapper functions for ListStore set value (to avoid variadic function issues)
// static void list_store_set_string(GtkListStore *store, GtkTreeIter *iter, gint column, const char *value) {
//     gtk_list_store_set(store, iter, column, value, -1);
// }
//
// static void list_store_set_int(GtkListStore *store, GtkTreeIter *iter, gint column, gint value) {
//     gtk_list_store_set(store, iter, column, value, -1);
// }
//
// static void list_store_set_boolean(GtkListStore *store, GtkTreeIter *iter, gint column, gboolean value) {
//     gtk_list_store_set(store, iter, column, value, -1);
// }
//
// static void list_store_set_float(GtkListStore *store, GtkTreeIter *iter, gint column, gfloat value) {
//     gtk_list_store_set(store, iter, column, value, -1);
// }
//
// static void list_store_set_double(GtkListStore *store, GtkTreeIter *iter, gint column, gdouble value) {
//     gtk_list_store_set(store, iter, column, value, -1);
// }
import "C"

import (
	"runtime"
	"unsafe"
)

// TreeModelFlags defines flags for the TreeModel
type TreeModelFlags int

const (
	// TreeModelListOnly the model is a simple list, not a tree
	TreeModelListOnly TreeModelFlags = C.GTK_TREE_MODEL_LIST_ONLY
	// TreeModelIters the model has persistent iterators
	TreeModelIters TreeModelFlags = C.GTK_TREE_MODEL_ITERS_PERSIST
)

// GType is a numeric type identifier
type GType C.GType

// Define common GTypes
const (
	G_TYPE_STRING  GType = C.G_TYPE_STRING
	G_TYPE_INT     GType = C.G_TYPE_INT
	G_TYPE_BOOLEAN GType = C.G_TYPE_BOOLEAN
	G_TYPE_FLOAT   GType = C.G_TYPE_FLOAT
	G_TYPE_DOUBLE  GType = C.G_TYPE_DOUBLE
)

// TreeIter represents an iterator in a TreeModel
type TreeIter struct {
	iter C.GtkTreeIter
}

// NewTreeIter creates a new TreeIter
func NewTreeIter() *TreeIter {
	return &TreeIter{}
}

// TreePath represents a path to a node in a tree model
type TreePath struct {
	path *C.GtkTreePath
}

// NewTreePath creates a new TreePath
func NewTreePath() *TreePath {
	path := &TreePath{
		path: C.gtk_tree_path_new(),
	}
	runtime.SetFinalizer(path, (*TreePath).Free)
	return path
}

// NewTreePathFromString creates a new TreePath from a string
func NewTreePathFromString(pathStr string) *TreePath {
	cPath := C.CString(pathStr)
	defer C.free(unsafe.Pointer(cPath))

	path := &TreePath{
		path: C.gtk_tree_path_new_from_string(cPath),
	}
	runtime.SetFinalizer(path, (*TreePath).Free)
	return path
}

// ToString converts a TreePath to a string
func (p *TreePath) ToString() string {
	cstr := C.gtk_tree_path_to_string(p.path)
	defer C.g_free(C.gpointer(unsafe.Pointer(cstr)))
	return C.GoString(cstr)
}

// Free frees the TreePath
func (p *TreePath) Free() {
	if p.path != nil {
		C.gtk_tree_path_free(p.path)
		p.path = nil
	}
}

// TreeModel represents a data model for TreeView
type TreeModel interface {
	// GetModelPointer returns the underlying GtkTreeModel pointer
	GetModelPointer() *C.GtkTreeModel

	// GetIter gets an iterator pointing to a path
	GetIter(path *TreePath) (*TreeIter, bool)

	// GetPath gets the path for an iterator
	GetPath(iter *TreeIter) *TreePath

	// GetValue gets the value of a cell
	GetValue(iter *TreeIter, column int) (interface{}, error)

	// GetNColumns gets the number of columns
	GetNColumns() int

	// GetColumnType gets the type of a column
	GetColumnType(column int) GType

	// IterNext moves the iterator to the next row
	IterNext(iter *TreeIter) bool
}

// TreeModelImplementor partially implements TreeModel
type TreeModelImplementor struct {
	model *C.GtkTreeModel
}

// GetModelPointer returns the underlying GtkTreeModel pointer
func (t *TreeModelImplementor) GetModelPointer() *C.GtkTreeModel {
	return t.model
}

// GetIter gets an iterator pointing to a path
func (t *TreeModelImplementor) GetIter(path *TreePath) (*TreeIter, bool) {
	iter := &TreeIter{}
	result := C.gtk_tree_model_get_iter(t.model, &iter.iter, path.path)
	return iter, result == C.TRUE
}

// GetPath gets the path for an iterator
func (t *TreeModelImplementor) GetPath(iter *TreeIter) *TreePath {
	path := &TreePath{
		path: C.gtk_tree_model_get_path(t.model, &iter.iter),
	}
	runtime.SetFinalizer(path, (*TreePath).Free)
	return path
}

// GetValue gets the value of a cell
func (t *TreeModelImplementor) GetValue(iter *TreeIter, column int) (interface{}, error) {
	var value C.GValue
	C.gtk_tree_model_get_value(t.model, &iter.iter, C.gint(column), &value)

	// Handle basic types
	gtype := C.get_value_type(&value)

	var result interface{}

	switch GType(gtype) {
	case G_TYPE_STRING:
		cstr := C.safe_get_string(&value)
		result = C.GoString(cstr)
	case G_TYPE_INT:
		result = int(C.g_value_get_int(&value))
	case G_TYPE_BOOLEAN:
		result = C.g_value_get_boolean(&value) == C.TRUE
	case G_TYPE_FLOAT:
		result = float32(C.g_value_get_float(&value))
	case G_TYPE_DOUBLE:
		result = float64(C.g_value_get_double(&value))
	default:
		result = nil
	}

	C.g_value_unset(&value)
	return result, nil
}

// GetNColumns gets the number of columns
func (t *TreeModelImplementor) GetNColumns() int {
	return int(C.gtk_tree_model_get_n_columns(t.model))
}

// GetColumnType gets the type of a column
func (t *TreeModelImplementor) GetColumnType(column int) GType {
	return GType(C.gtk_tree_model_get_column_type(t.model, C.gint(column)))
}

// IterNext moves the iterator to the next row
func (t *TreeModelImplementor) IterNext(iter *TreeIter) bool {
	return C.gtk_tree_model_iter_next(t.model, &iter.iter) == C.TRUE
}

// TreeStore implements TreeModel for a tree store
type TreeStore struct {
	TreeModelImplementor
	store *C.GtkTreeStore
}

// NewTreeStore creates a new TreeStore
func NewTreeStore(types ...GType) *TreeStore {
	// Convert Go types to C types
	cTypes := make([]C.GType, len(types))
	for i, t := range types {
		cTypes[i] = C.GType(t)
	}

	// Create the tree store
	store := &TreeStore{
		store: C.gtk_tree_store_newv(C.gint(len(types)), (*C.GType)(&cTypes[0])),
	}
	store.model = (*C.GtkTreeModel)(unsafe.Pointer(store.store))

	runtime.SetFinalizer(store, (*TreeStore).Free)
	return store
}

// SetValue sets a value in the TreeStore
func (s *TreeStore) SetValue(iter *TreeIter, column int, value interface{}) {
	// Set value based on type using our C wrappers
	switch v := value.(type) {
	case string:
		cstr := C.CString(v)
		defer C.free(unsafe.Pointer(cstr))
		C.tree_store_set_string(s.store, &iter.iter, C.gint(column), cstr)
	case int:
		C.tree_store_set_int(s.store, &iter.iter, C.gint(column), C.gint(v))
	case bool:
		var cbool C.gboolean
		if v {
			cbool = C.TRUE
		} else {
			cbool = C.FALSE
		}
		C.tree_store_set_boolean(s.store, &iter.iter, C.gint(column), cbool)
	case float32:
		C.tree_store_set_float(s.store, &iter.iter, C.gint(column), C.gfloat(v))
	case float64:
		C.tree_store_set_double(s.store, &iter.iter, C.gint(column), C.gdouble(v))
	}
}

// Append appends a new row to the TreeStore
func (s *TreeStore) Append(parent *TreeIter) *TreeIter {
	iter := &TreeIter{}
	var parentIter *C.GtkTreeIter
	if parent != nil {
		parentIter = &parent.iter
	}
	C.gtk_tree_store_append(s.store, &iter.iter, parentIter)
	return iter
}

// Free frees the TreeStore
func (s *TreeStore) Free() {
	if s.store != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(s.store)))
		s.store = nil
		s.model = nil
	}
}

// ListStore implements TreeModel for a list store
type ListStore struct {
	TreeModelImplementor
	store *C.GtkListStore
}

// NewListStore creates a new ListStore
func NewListStore(types ...GType) *ListStore {
	// Convert Go types to C types
	cTypes := make([]C.GType, len(types))
	for i, t := range types {
		cTypes[i] = C.GType(t)
	}

	// Create the list store
	store := &ListStore{
		store: C.gtk_list_store_newv(C.gint(len(types)), (*C.GType)(&cTypes[0])),
	}
	store.model = (*C.GtkTreeModel)(unsafe.Pointer(store.store))

	runtime.SetFinalizer(store, (*ListStore).Free)
	return store
}

// SetValue sets a value in the ListStore
func (s *ListStore) SetValue(iter *TreeIter, column int, value interface{}) {
	// Set value based on type using our C wrappers
	switch v := value.(type) {
	case string:
		cstr := C.CString(v)
		defer C.free(unsafe.Pointer(cstr))
		C.list_store_set_string(s.store, &iter.iter, C.gint(column), cstr)
	case int:
		C.list_store_set_int(s.store, &iter.iter, C.gint(column), C.gint(v))
	case bool:
		var cbool C.gboolean
		if v {
			cbool = C.TRUE
		} else {
			cbool = C.FALSE
		}
		C.list_store_set_boolean(s.store, &iter.iter, C.gint(column), cbool)
	case float32:
		C.list_store_set_float(s.store, &iter.iter, C.gint(column), C.gfloat(v))
	case float64:
		C.list_store_set_double(s.store, &iter.iter, C.gint(column), C.gdouble(v))
	}
}

// Append appends a new row to the ListStore
func (s *ListStore) Append() *TreeIter {
	iter := &TreeIter{}
	C.gtk_list_store_append(s.store, &iter.iter)
	return iter
}

// Free frees the ListStore
func (s *ListStore) Free() {
	if s.store != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(s.store)))
		s.store = nil
		s.model = nil
	}
}

