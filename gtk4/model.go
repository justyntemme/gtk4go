// Package gtk4 provides model functionality for GTK4
// File: gtk4go/gtk4/model.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Wrapper functions for GValue type checking (C macros)
// static gboolean value_holds_string(GValue* value) {
//     return G_VALUE_HOLDS_STRING(value);
// }
//
// static gboolean value_holds_int(GValue* value) {
//     return G_VALUE_HOLDS_INT(value);
// }
//
// static gboolean value_holds_boolean(GValue* value) {
//     return G_VALUE_HOLDS_BOOLEAN(value);
// }
//
// static gboolean value_holds_double(GValue* value) {
//     return G_VALUE_HOLDS_DOUBLE(value);
// }
//
// static gboolean value_holds_float(GValue* value) {
//     return G_VALUE_HOLDS_FLOAT(value);
// }
//
// static GType value_get_type(GValue* value) {
//     return G_VALUE_TYPE(value);
// }
//
// // Wrapper for value extraction helpers to avoid variadic function issues
// static const char* value_get_string_safe(GValue* value) {
//     if (!G_VALUE_HOLDS_STRING(value)) return NULL;
//     const char* str = g_value_get_string(value);
//     return str ? str : "";
// }
//
// static int value_get_int_safe(GValue* value) {
//     if (!G_VALUE_HOLDS_INT(value)) return 0;
//     return g_value_get_int(value);
// }
//
// static gboolean value_get_boolean_safe(GValue* value) {
//     if (!G_VALUE_HOLDS_BOOLEAN(value)) return FALSE;
//     return g_value_get_boolean(value);
// }
//
// static double value_get_double_safe(GValue* value) {
//     if (!G_VALUE_HOLDS_DOUBLE(value)) return 0.0;
//     return g_value_get_double(value);
// }
//
// static float value_get_float_safe(GValue* value) {
//     if (!G_VALUE_HOLDS_FLOAT(value)) return 0.0f;
//     return g_value_get_float(value);
// }
//
// // Wrapper functions for value setting
// static void value_set_string(GValue* value, const char* str) {
//     g_value_set_string(value, str);
// }
//
// static void value_set_int(GValue* value, int i) {
//     g_value_set_int(value, i);
// }
//
// static void value_set_boolean(GValue* value, gboolean b) {
//     g_value_set_boolean(value, b);
// }
//
// static void value_set_double(GValue* value, double d) {
//     g_value_set_double(value, d);
// }
//
// static void value_set_float(GValue* value, float f) {
//     g_value_set_float(value, f);
// }
//
// // GListModel helpers
// static guint g_list_model_get_n_items_safe(GListModel* model) {
//     return model ? g_list_model_get_n_items(model) : 0;
// }
//
// static GType g_list_model_get_item_type_safe(GListModel* model) {
//     return model ? g_list_model_get_item_type(model) : G_TYPE_NONE;
// }
//
// static gpointer g_list_model_get_item_safe(GListModel* model, guint position) {
//     return model ? g_list_model_get_item(model, position) : NULL;
// }
//
// // Create a GListStore
// static GListStore* create_list_store(GType item_type) {
//     return g_list_store_new(item_type);
// }
//
// // Append to a GListStore
// static void list_store_append(GListStore* store, gpointer item) {
//     g_list_store_append(store, item);
// }
//
// // Remove from a GListStore
// static void list_store_remove(GListStore* store, guint position) {
//     g_list_store_remove(store, position);
// }
//
// // Insert into a GListStore
// static void list_store_insert(GListStore* store, guint position, gpointer item) {
//     g_list_store_insert(store, position, item);
// }
//
// // Create a GtkTreeListModel
// static GtkTreeListModel* create_tree_list_model(GListModel* root, gboolean passthrough,
//                               gboolean autoexpand, GtkTreeListModelCreateModelFunc create_func,
//                               gpointer user_data, GDestroyNotify user_destroy) {
//     return gtk_tree_list_model_new(root, passthrough, autoexpand, create_func, user_data, user_destroy);
// }
//
// // GtkTreeListRow helpers
// static GtkTreeListRow* get_tree_list_row(GtkTreeListModel* model, guint position) {
//     GObject* obj = G_OBJECT(g_list_model_get_item(G_LIST_MODEL(model), position));
//     return GTK_TREE_LIST_ROW(obj);
// }
//
// static GListModel* tree_list_row_get_children(GtkTreeListRow* row) {
//     return row ? gtk_tree_list_row_get_children(row) : NULL;
// }
//
// static gpointer tree_list_row_get_item(GtkTreeListRow* row) {
//     return row ? gtk_tree_list_row_get_item(row) : NULL;
// }
//
// static gboolean tree_list_row_is_expandable(GtkTreeListRow* row) {
//     return row ? gtk_tree_list_row_is_expandable(row) : FALSE;
// }
//
// static gboolean tree_list_row_get_expanded(GtkTreeListRow* row) {
//     return row ? gtk_tree_list_row_get_expanded(row) : FALSE;
// }
//
// static void tree_list_row_set_expanded(GtkTreeListRow* row, gboolean expanded) {
//     if (row) gtk_tree_list_row_set_expanded(row, expanded);
// }
import "C"

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"
)

// GType is a numeric type identifier
type GType C.GType

// Define common GTypes
const (
	G_TYPE_NONE    GType = C.G_TYPE_NONE
	G_TYPE_STRING  GType = C.G_TYPE_STRING
	G_TYPE_INT     GType = C.G_TYPE_INT
	G_TYPE_BOOLEAN GType = C.G_TYPE_BOOLEAN
	G_TYPE_FLOAT   GType = C.G_TYPE_FLOAT
	G_TYPE_DOUBLE  GType = C.G_TYPE_DOUBLE
	G_TYPE_OBJECT  GType = C.G_TYPE_OBJECT
)

// ListModel is an interface for list models
type ListModel interface {
	// GetNItems returns the number of items in the model
	GetNItems() int

	// GetItem returns the item at the given position
	GetItem(position int) interface{}

	// GetItemType returns the type of items in the model
	GetItemType() GType
}

// GListModel is a wrapper around a GListModel
type GListModel struct {
	model *C.GListModel
}

// NewGListModel creates a wrapper around an existing GListModel
func NewGListModel(model *C.GListModel) *GListModel {
	if model == nil {
		return nil
	}

	result := &GListModel{
		model: model,
	}

	runtime.SetFinalizer(result, (*GListModel).free)
	return result
}

// free frees the GListModel
func (m *GListModel) free() {
	if m.model != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(m.model)))
		m.model = nil
	}
}

// GetNItems returns the number of items in the model
func (m *GListModel) GetNItems() int {
	return int(C.g_list_model_get_n_items_safe(m.model))
}

// GetItemType returns the type of items in the model
func (m *GListModel) GetItemType() GType {
	return GType(C.g_list_model_get_item_type_safe(m.model))
}

// GetItem returns the item at the given position
func (m *GListModel) GetItem(position int) interface{} {
	item := C.g_list_model_get_item_safe(m.model, C.guint(position))
	if item == nil {
		return nil
	}

	// We need to unref the item when we're done with it
	defer C.g_object_unref(item)

	// Convert to appropriate Go type based on item type
	// This depends on the specific use case and may need customization
	return castGObjectToGoValue((*C.GObject)(item))
}

// castGObjectToGoValue converts a GObject to a Go value
// This is a placeholder and would need to be implemented based on your needs
func castGObjectToGoValue(obj *C.GObject) interface{} {
	// In a real implementation, you would check the object type
	// and convert appropriately. For now, we just return the pointer.
	return uintptr(unsafe.Pointer(obj))
}

// ListStore is a wrapper around a GListStore
type ListStore struct {
	store *C.GListStore
	model *GListModel // Wrapper for model interface
}

// NewListStore creates a new ListStore with the given item type
func NewListStore(itemType GType) *ListStore {
	store := &ListStore{
		store: C.create_list_store(C.GType(itemType)),
	}

	// Create model wrapper
	store.model = NewGListModel((*C.GListModel)(unsafe.Pointer(store.store)))

	runtime.SetFinalizer(store, (*ListStore).Free)
	return store
}

// Append adds an item to the end of the list
func (s *ListStore) Append(item interface{}) {
	// Convert item to GObject based on its type
	// This is a simplified version and would need to be expanded
	itemPtr := convertGoValueToGObject(item)
	if itemPtr != nil {
		C.list_store_append(s.store, C.gpointer(itemPtr))
	}
}

// Insert inserts an item at the specified position
func (s *ListStore) Insert(position int, item interface{}) {
	itemPtr := convertGoValueToGObject(item)
	if itemPtr != nil {
		C.list_store_insert(s.store, C.guint(position), C.gpointer(itemPtr))
	}
}

// Remove removes the item at the specified position
func (s *ListStore) Remove(position int) {
	C.list_store_remove(s.store, C.guint(position))
}

// GetModel returns the underlying GListModel wrapper
func (s *ListStore) GetModel() *GListModel {
	return s.model
}

// Free frees the ListStore
func (s *ListStore) Free() {
	if s.store != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(s.store)))
		s.store = nil
	}

	s.model = nil
}

// GetNItems returns the number of items in the model
func (s *ListStore) GetNItems() int {
	return s.model.GetNItems()
}

// GetItem returns the item at the given position
func (s *ListStore) GetItem(position int) interface{} {
	return s.model.GetItem(position)
}

// GetItemType returns the type of items in the model
func (s *ListStore) GetItemType() GType {
	return s.model.GetItemType()
}

// convertGoValueToGObject converts a Go value to a GObject
// This is a placeholder and would need to be implemented based on your needs
func convertGoValueToGObject(value interface{}) unsafe.Pointer {
	// In a real implementation, you would create appropriate GObjects
	// based on the Go type. For now, we just return nil.
	return nil
}

// TreeListModel is a wrapper around a GTK TreeListModel
type TreeListModel struct {
	model     *C.GtkTreeListModel
	listModel *GListModel // Wrapper for model interface
}

// TreeListCreateModelFunc is the type for a function that creates child models for tree items
type TreeListCreateModelFunc func(item interface{}) ListModel

// treeListCreateModelCallback is the C callback for creating child models
//
//export treeListCreateModelCallback
func treeListCreateModelCallback(item *C.gpointer, userData C.gpointer) *C.GListModel {
	// Extract Go function pointer from user data
	// This is complex and would need careful implementation
	// For now, return nil
	return nil
}

// NewTreeListModel creates a new TreeListModel
func NewTreeListModel(root ListModel, passthrough bool, autoexpand bool, createFunc TreeListCreateModelFunc) *TreeListModel {
	// Get the GListModel from the root if it's our implementation
	var rootModel *C.GListModel
	if gm, ok := root.(*GListModel); ok {
		rootModel = gm.model
	} else {
		// Create a wrapper around a custom ListModel implementation
		// This is complex and would need careful implementation
		// For now, return nil
		return nil
	}

	// Convert booleans to C booleans
	var cPassthrough, cAutoexpand C.gboolean
	if passthrough {
		cPassthrough = C.TRUE
	}
	if autoexpand {
		cAutoexpand = C.TRUE
	}

	// Create the TreeListModel
	// Note: The create_func callback is complex and would need careful implementation
	// For now, we pass nil
	model := C.create_tree_list_model(rootModel, cPassthrough, cAutoexpand, nil, nil, nil)

	result := &TreeListModel{
		model: model,
	}

	// Create model wrapper
	result.listModel = NewGListModel((*C.GListModel)(unsafe.Pointer(model)))

	runtime.SetFinalizer(result, (*TreeListModel).Free)
	return result
}

// GetModel returns the underlying GListModel wrapper
func (t *TreeListModel) GetModel() *GListModel {
	return t.listModel
}

// Free frees the TreeListModel
func (t *TreeListModel) Free() {
	if t.model != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(t.model)))
		t.model = nil
	}

	t.listModel = nil
}

// GetNItems returns the number of items in the model
func (t *TreeListModel) GetNItems() int {
	return t.listModel.GetNItems()
}

// GetItem returns the item at the given position
func (t *TreeListModel) GetItem(position int) interface{} {
	return t.listModel.GetItem(position)
}

// GetItemType returns the type of items in the model
func (t *TreeListModel) GetItemType() GType {
	return t.listModel.GetItemType()
}

// TreeListRow is a wrapper around a GTK TreeListRow
type TreeListRow struct {
	row *C.GtkTreeListRow
}

// NewTreeListRow creates a wrapper around a GTK TreeListRow
func NewTreeListRow(row *C.GtkTreeListRow) *TreeListRow {
	if row == nil {
		return nil
	}

	result := &TreeListRow{
		row: row,
	}

	runtime.SetFinalizer(result, (*TreeListRow).Free)
	return result
}

// GetRow returns the TreeListRow at the specified position in a TreeListModel
func (t *TreeListModel) GetRow(position int) *TreeListRow {
	row := C.get_tree_list_row(t.model, C.guint(position))
	return NewTreeListRow(row)
}

// GetChildren returns the children of the row as a GListModel
func (r *TreeListRow) GetChildren() *GListModel {
	children := C.tree_list_row_get_children(r.row)
	if children == nil {
		return nil
	}

	return NewGListModel(children)
}

// GetItem returns the item represented by the row
func (r *TreeListRow) GetItem() interface{} {
	item := C.tree_list_row_get_item(r.row)
	if item == nil {
		return nil
	}

	// We need to unref the item when we're done with it
	defer C.g_object_unref(item)

	// Convert to appropriate Go type
	return castGObjectToGoValue((*C.GObject)(item))
}

// IsExpandable returns whether the row is expandable
func (r *TreeListRow) IsExpandable() bool {
	return C.tree_list_row_is_expandable(r.row) == C.TRUE
}

// GetExpanded returns whether the row is expanded
func (r *TreeListRow) GetExpanded() bool {
	return C.tree_list_row_get_expanded(r.row) == C.TRUE
}

// SetExpanded sets whether the row is expanded
func (r *TreeListRow) SetExpanded(expanded bool) {
	var cExpanded C.gboolean
	if expanded {
		cExpanded = C.TRUE
	}
	C.tree_list_row_set_expanded(r.row, cExpanded)
}

// Free frees the TreeListRow
func (r *TreeListRow) Free() {
	if r.row != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(r.row)))
		r.row = nil
	}
}

// SimpleListModel is a basic Go implementation of the ListModel interface
type SimpleListModel struct {
	items    []interface{}
	itemType GType
}

// NewSimpleListModel creates a new SimpleListModel with the given item type
func NewSimpleListModel(itemType GType) *SimpleListModel {
	return &SimpleListModel{
		items:    make([]interface{}, 0),
		itemType: itemType,
	}
}

// GetNItems returns the number of items in the model
func (m *SimpleListModel) GetNItems() int {
	return len(m.items)
}

// GetItem returns the item at the given position
func (m *SimpleListModel) GetItem(position int) interface{} {
	if position >= 0 && position < len(m.items) {
		return m.items[position]
	}
	return nil
}

// GetItemType returns the type of items in the model
func (m *SimpleListModel) GetItemType() GType {
	return m.itemType
}

// AddItem adds an item to the model
func (m *SimpleListModel) AddItem(item interface{}) {
	m.items = append(m.items, item)
}

// ClearItems removes all items from the model
func (m *SimpleListModel) ClearItems() {
	m.items = make([]interface{}, 0)
}

// GetModelFromWidget gets a model from a widget, if applicable
func GetModelFromWidget(widget interface{}) interface{} {
	// Try different model getter methods based on widget type
	if w, ok := widget.(interface{ GetModel() interface{} }); ok {
		return w.GetModel()
	}

	// Legacy method for TreeModel compatibility
	if w, ok := widget.(interface{ GetModel() *GListModel }); ok {
		return w.GetModel()
	}

	// Legacy method for TreeModel compatibility
	if w, ok := widget.(interface{ GetModel() ListModel }); ok {
		return w.GetModel()
	}

	return nil
}

// SetModelForWidget sets a model for a widget, if applicable
func SetModelForWidget(widget interface{}, model interface{}) bool {
	// Try different model setter methods based on widget and model types
	if w, ok := widget.(interface{ SetModel(interface{}) bool }); ok {
		return w.SetModel(model)
	}

	// Legacy method for TreeModel compatibility
	if w, ok := widget.(interface{ SetModel(*GListModel) }); ok {
		if m, ok := model.(*GListModel); ok {
			w.SetModel(m)
			return true
		}
	}

	// Legacy method for TreeModel compatibility
	if w, ok := widget.(interface{ SetModel(ListModel) }); ok {
		if m, ok := model.(ListModel); ok {
			w.SetModel(m)
			return true
		}
	}

	return false
}

// Helper functions for working with models

// ValueFromGValue converts a GValue to a Go value based on its type
func ValueFromGValue(value *C.GValue) interface{} {
	if value == nil {
		return nil
	}

	// Use safe wrapper functions to extract values based on type
	switch {
	case C.value_holds_string(value) != 0:
		cstr := C.value_get_string_safe(value)
		if cstr == nil {
			return ""
		}
		return C.GoString(cstr)
	case C.value_holds_int(value) != 0:
		return int(C.value_get_int_safe(value))
	case C.value_holds_boolean(value) != 0:
		return C.value_get_boolean_safe(value) == C.TRUE
	case C.value_holds_double(value) != 0:
		return float64(C.value_get_double_safe(value))
	case C.value_holds_float(value) != 0:
		return float32(C.value_get_float_safe(value))
	default:
		return nil
	}
}

// SetGValueFromValue sets a GValue from a Go value
func SetGValueFromValue(gvalue *C.GValue, value interface{}) bool {
	if gvalue == nil {
		return false
	}

	// Initialize if needed
	if C.value_get_type(gvalue) == 0 {
		// Guess the type based on the value's type
		switch value.(type) {
		case string:
			C.g_value_init(gvalue, C.G_TYPE_STRING)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			C.g_value_init(gvalue, C.G_TYPE_INT)
		case bool:
			C.g_value_init(gvalue, C.G_TYPE_BOOLEAN)
		case float64:
			C.g_value_init(gvalue, C.G_TYPE_DOUBLE)
		case float32:
			C.g_value_init(gvalue, C.G_TYPE_FLOAT)
		default:
			return false
		}
	}

	// Set the value based on its type
	switch v := value.(type) {
	case string:
		cstr := C.CString(v)
		defer C.free(unsafe.Pointer(cstr))
		C.value_set_string(gvalue, cstr)
		return true
	case int:
		C.value_set_int(gvalue, C.int(v))
		return true
	case bool:
		var cbool C.gboolean
		if v {
			cbool = C.TRUE
		} else {
			cbool = C.FALSE
		}
		C.value_set_boolean(gvalue, cbool)
		return true
	case float64:
		C.value_set_double(gvalue, C.double(v))
		return true
	case float32:
		C.value_set_float(gvalue, C.float(v))
		return true
	default:
		return false
	}
}

// Model-related variables
var (
	modelSystemInitialized = false
)

// CleanupModelResources does model resource cleanup
func CleanupModelResources() {
	// Nothing special to clean up in this implementation
	modelSystemInitialized = false
}

// InitializeModelSystem initializes model systems
func InitializeModelSystem() {
	if modelSystemInitialized {
		return
	}

	// No special initialization needed in this implementation
	modelSystemInitialized = true
}

// GenericModelData represents a generic data structure for models
type GenericModelData struct {
	Type      GType
	Value     interface{}
	UserData  interface{}
	Reference bool
}

// NewGenericModelData creates a new GenericModelData
func NewGenericModelData(value interface{}) *GenericModelData {
	var dataType GType

	// Determine type
	switch value.(type) {
	case string:
		dataType = G_TYPE_STRING
	case int:
		dataType = G_TYPE_INT
	case bool:
		dataType = G_TYPE_BOOLEAN
	case float32:
		dataType = G_TYPE_FLOAT
	case float64:
		dataType = G_TYPE_DOUBLE
	default:
		dataType = 0 // Unknown type
	}

	return &GenericModelData{
		Type:      dataType,
		Value:     value,
		Reference: false,
	}
}

// These functions should be added to the model.go file
// or extracted to a separate module if appropriate

// ListStore is a wrapper around a GListStore
type ListStore struct {
	store *C.GListStore
	model *GListModel // Wrapper for model interface
}

// NewListStore creates a new ListStore with the given item type
func NewListStore(itemType GType) *ListStore {
	fmt.Printf("NewListStore: Creating store with item type %v\n", itemType)

	// Validate item type - this is critical to prevent segfaults
	// GTK requires item_type to be G_TYPE_OBJECT or a subclass
	if itemType != G_TYPE_OBJECT && itemType != G_TYPE_STRING {
		fmt.Printf("NewListStore: WARNING - item type %v is not G_TYPE_OBJECT, forcing to G_TYPE_OBJECT\n", itemType)
		itemType = G_TYPE_OBJECT
	}

	store := &ListStore{
		store: C.create_list_store(C.GType(itemType)),
	}

	if store.store == nil {
		fmt.Println("NewListStore: ERROR - Failed to create GListStore")
		return nil
	}

	// Create model wrapper
	store.model = NewGListModel((*C.GListModel)(unsafe.Pointer(store.store)))

	runtime.SetFinalizer(store, (*ListStore).Free)
	fmt.Println("NewListStore: ListStore created successfully")
	return store
}

// Append adds an item to the end of the list
func (s *ListStore) Append(item interface{}) {
	fmt.Printf("ListStore.Append: Appending item %v\n", item)

	if s.store == nil {
		fmt.Println("ListStore.Append: WARNING - s.store is nil")
		return
	}

	// Convert item to GObject based on its type
	var itemPtr unsafe.Pointer

	switch v := item.(type) {
	case string:
		// Create a GObject from the string
		obj := CreateStringObject(v)
		if obj == nil {
			fmt.Println("ListStore.Append: ERROR - Failed to create GObject from string")
			return
		}
		itemPtr = unsafe.Pointer(obj)
	case int, float64, float32, bool:
		// For basic types, create a string representation and wrap as GObject
		obj := CreateStringObject(fmt.Sprintf("%v", v))
		if obj == nil {
			fmt.Println("ListStore.Append: ERROR - Failed to create GObject from value")
			return
		}
		itemPtr = unsafe.Pointer(obj)
	case unsafe.Pointer:
		// Use pointer directly
		itemPtr = v
	case uintptr:
		// Convert uintptr to unsafe.Pointer
		itemPtr = unsafe.Pointer(v)
	default:
		fmt.Printf("ListStore.Append: WARNING - Unsupported item type %T\n", item)
		return
	}

	if itemPtr != nil {
		C.list_store_append(s.store, C.gpointer(itemPtr))
		fmt.Println("ListStore.Append: Item appended successfully")
	}
}

// Insert inserts an item at the specified position
func (s *ListStore) Insert(position int, item interface{}) {
	fmt.Printf("ListStore.Insert: Inserting item %v at position %d\n", item, position)

	if s.store == nil {
		fmt.Println("ListStore.Insert: WARNING - s.store is nil")
		return
	}

	// Convert item to GObject based on its type (similar to Append)
	var itemPtr unsafe.Pointer

	switch v := item.(type) {
	case string:
		// Create a GObject from the string
		obj := CreateStringObject(v)
		if obj == nil {
			fmt.Println("ListStore.Insert: ERROR - Failed to create GObject from string")
			return
		}
		itemPtr = unsafe.Pointer(obj)
	case unsafe.Pointer:
		// Use pointer directly
		itemPtr = v
	case uintptr:
		// Convert uintptr to unsafe.Pointer
		itemPtr = unsafe.Pointer(v)
	default:
		fmt.Printf("ListStore.Insert: WARNING - Unsupported item type %T\n", item)
		return
	}

	if itemPtr != nil {
		C.list_store_insert(s.store, C.guint(position), C.gpointer(itemPtr))
		fmt.Printf("ListStore.Insert: Item inserted successfully at position %d\n", position)
	}
}

// Remove removes the item at the specified position
func (s *ListStore) Remove(position int) {
	fmt.Printf("ListStore.Remove: Removing item at position %d\n", position)

	if s.store == nil {
		fmt.Println("ListStore.Remove: WARNING - s.store is nil")
		return
	}

	C.list_store_remove(s.store, C.guint(position))
	fmt.Printf("ListStore.Remove: Item removed from position %d\n", position)
}

// GetModel returns the underlying GListModel wrapper
func (s *ListStore) GetModel() *GListModel {
	return s.model
}

// Free frees the ListStore
func (s *ListStore) Free() {
	fmt.Println("ListStore.Free: Cleaning up resources")

	if s.store != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(s.store)))
		s.store = nil
	}

	s.model = nil
}

// GetNItems returns the number of items in the model
func (s *ListStore) GetNItems() int {
	if s.model == nil {
		return 0
	}
	return s.model.GetNItems()
}

// GetItem returns the item at the given position
func (s *ListStore) GetItem(position int) interface{} {
	if s.model == nil {
		return nil
	}

	item := s.model.GetItem(position)

	// If item is a pointer, try to convert it to a string
	if ptr, ok := item.(uintptr); ok {
		gobj := (*C.GObject)(unsafe.Pointer(ptr))
		str := GetStringFromObject(gobj)
		if str != "" {
			return str
		}
	}

	return item
}

// GetItemType returns the type of items in the model
func (s *ListStore) GetItemType() GType {
	if s.model == nil {
		return G_TYPE_NONE
	}
	return s.model.GetItemType()
}

// AppendString adds a string to a ListStore
func (s *ListStore) AppendString(text string) {
	fmt.Printf("AppendString: Adding string %q to ListStore\n", text)

	if s.store == nil {
		fmt.Println("AppendString: WARNING - s.store is nil")
		return
	}

	obj := CreateStringObject(text)
	if obj == nil {
		fmt.Println("AppendString: ERROR - Failed to create string object")
		return
	}

	C.list_store_append(s.store, C.gpointer(unsafe.Pointer(obj)))
	fmt.Println("AppendString: String added successfully")
}

// CreateStringObject creates a GObject from a string for use in ListStore
func CreateStringObject(text string) *C.GObject {
	fmt.Printf("CreateStringObject: Creating GObject for string %q\n", text)

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

