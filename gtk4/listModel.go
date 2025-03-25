// Package gtk4 provides list model functionality for GTK4
// File: gtk4go/gtk4/listmodel.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // ListModel callbacks
// extern void listModelItemsChangedCallback(GListModel *model, guint position, guint removed, guint added, gpointer user_data);
//
// // Connect signal for items-changed
// static gulong connectListModelItemsChanged(GListModel *model, gpointer user_data) {
//     return g_signal_connect(model, "items-changed", G_CALLBACK(listModelItemsChangedCallback), user_data);
// }
//
// // StringList operations
// static GtkStringList* createStringList() {
//     return gtk_string_list_new(NULL);
// }
//
// static void stringListAppend(GtkStringList *list, const char *string) {
//     gtk_string_list_append(list, string);
// }
//
// static void stringListRemove(GtkStringList *list, guint position) {
//     gtk_string_list_remove(list, position);
// }
//
// static guint stringListGetNItems(GtkStringList *list) {
//     return g_list_model_get_n_items(G_LIST_MODEL(list));
// }
//
// static const char* stringListGetString(GtkStringList *list, guint position) {
//     const char* result = NULL;
//     GtkStringObject *obj = GTK_STRING_OBJECT(g_list_model_get_item(G_LIST_MODEL(list), position));
//     if (obj != NULL) {
//         result = gtk_string_object_get_string(obj);
//         g_object_unref(obj);
//     }
//     return result;
// }
//
// // ListStore operations
// static GListStore* createListStore(GType item_type) {
//     return g_list_store_new(item_type);
// }
//
// static void listStoreAppend(GListStore *store, gpointer item) {
//     g_list_store_append(store, item);
// }
//
// static void listStoreRemove(GListStore *store, guint position) {
//     g_list_store_remove(store, position);
// }
//
// static gpointer listModelGetItem(GListModel *model, guint position) {
//     return g_list_model_get_item(model, position);
// }
//
// static guint listModelGetNItems(GListModel *model) {
//     return g_list_model_get_n_items(model);
// }
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// ListModelItemsChangedCallback represents a callback for list model changes
type ListModelItemsChangedCallback func(position, removed, added int)

var (
	listModelCallbacks     = make(map[uintptr]ListModelItemsChangedCallback)
	listModelCallbackMutex sync.RWMutex
)

//export listModelItemsChangedCallback
func listModelItemsChangedCallback(model *C.GListModel, position, removed, added C.guint, userData C.gpointer) {
	listModelCallbackMutex.RLock()
	defer listModelCallbackMutex.RUnlock()

	// Convert model pointer to uintptr for lookup
	modelPtr := uintptr(unsafe.Pointer(model))

	// Find and call the callback
	if callback, ok := listModelCallbacks[modelPtr]; ok {
		callback(int(position), int(removed), int(added))
	}
}

// ListModel is an interface for GTK list models
type ListModel interface {
	// GetListModel returns the underlying GListModel pointer
	GetListModel() *C.GListModel

	// GetNItems returns the number of items in the model
	GetNItems() int

	// GetItem returns the item at the given position
	GetItem(position int) interface{}

	// ConnectItemsChanged connects a callback for list model changes
	ConnectItemsChanged(callback ListModelItemsChangedCallback)

	// Destroy frees resources associated with the list model
	Destroy()
}

// BaseListModel provides common functionality for list models
type BaseListModel struct {
	model *C.GListModel
}

// GetListModel returns the underlying GListModel pointer
func (m *BaseListModel) GetListModel() *C.GListModel {
	return m.model
}

// GetNItems returns the number of items in the model
func (m *BaseListModel) GetNItems() int {
	return int(C.listModelGetNItems(m.model))
}

// ConnectItemsChanged connects a callback for list model changes
func (m *BaseListModel) ConnectItemsChanged(callback ListModelItemsChangedCallback) {
	if callback == nil {
		return
	}

	listModelCallbackMutex.Lock()
	defer listModelCallbackMutex.Unlock()

	// Store the callback in the map
	modelPtr := uintptr(unsafe.Pointer(m.model))
	listModelCallbacks[modelPtr] = callback

	// Connect the signal
	C.connectListModelItemsChanged(m.model, C.gpointer(unsafe.Pointer(m.model)))
}

// Destroy frees resources associated with the list model
func (m *BaseListModel) Destroy() {
	if m.model != nil {
		listModelCallbackMutex.Lock()
		delete(listModelCallbacks, uintptr(unsafe.Pointer(m.model)))
		listModelCallbackMutex.Unlock()

		C.g_object_unref(C.gpointer(unsafe.Pointer(m.model)))
		m.model = nil
	}
}

// StringList is a list model for strings
type StringList struct {
	BaseListModel
	stringList *C.GtkStringList
}

// NewStringList creates a new string list
func NewStringList() *StringList {
	stringList := C.createStringList()
	list := &StringList{
		BaseListModel: BaseListModel{
			model: (*C.GListModel)(unsafe.Pointer(stringList)),
		},
		stringList: stringList,
	}

	runtime.SetFinalizer(list, (*StringList).Destroy)
	return list
}

// Append adds a string to the list
func (l *StringList) Append(str string) {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	C.stringListAppend(l.stringList, cStr)
}

// Remove removes a string from the list at the given position
func (l *StringList) Remove(position int) {
	if position >= 0 && position < l.GetNItems() {
		C.stringListRemove(l.stringList, C.guint(position))
	}
}

// GetString returns the string at the given position
func (l *StringList) GetString(position int) string {
	if position < 0 || position >= l.GetNItems() {
		return ""
	}

	cStr := C.stringListGetString(l.stringList, C.guint(position))
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// GetItem returns the item at the given position as an interface{}
func (l *StringList) GetItem(position int) interface{} {
	return l.GetString(position)
}

// Destroy frees resources associated with the string list
func (l *StringList) Destroy() {
	l.BaseListModel.Destroy()
	l.stringList = nil
}

// ListStore is a generic list store
type ListStore struct {
	BaseListModel
	listStore *C.GListStore
	itemType  C.GType
	items     []interface{} // Keep Go references to items
}

// NewListStore creates a new list store with the given item type
func NewListStore(itemType C.GType) *ListStore {
	listStore := C.createListStore(itemType)
	store := &ListStore{
		BaseListModel: BaseListModel{
			model: (*C.GListModel)(unsafe.Pointer(listStore)),
		},
		listStore: listStore,
		itemType:  itemType,
		items:     make([]interface{}, 0),
	}

	runtime.SetFinalizer(store, (*ListStore).Destroy)
	return store
}

// Append adds an item to the list store
// Note: The implementation of this method depends on the type of items stored
// and would need customization for practical use
func (s *ListStore) Append(item interface{}) {
	// This is a simplified implementation that would need to be adapted
	// based on the actual item types being stored
	var cItem C.gpointer
	// Convert item to appropriate C pointer based on type
	// (Implementation would depend on the actual item types)

	C.listStoreAppend(s.listStore, cItem)
	s.items = append(s.items, item) // Store Go reference
}

// Remove removes an item from the list store at the given position
func (s *ListStore) Remove(position int) {
	if position >= 0 && position < len(s.items) {
		C.listStoreRemove(s.listStore, C.guint(position))
		// Remove the Go reference
		s.items = append(s.items[:position], s.items[position+1:]...)
	}
}

// GetItem returns the item at the given position
func (s *ListStore) GetItem(position int) interface{} {
	if position < 0 || position >= len(s.items) {
		return nil
	}
	return s.items[position]
}

// Destroy frees resources associated with the list store
func (s *ListStore) Destroy() {
	s.BaseListModel.Destroy()
	s.listStore = nil
	s.items = nil
}