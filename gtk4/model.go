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
import "C"
import "unsafe"

// ListModel is an interface for list models
// Kept here for compatibility with listView.go
// Using int instead of uint for compatibility
type ListModel interface {
	// GetNItems returns the number of items in the model
	GetNItems() int
	
	// GetItem returns the item at the given position
	GetItem(position int) interface{}
	
	// GetItemType returns the type of items in the model
	GetItemType() GType
}

// SimpleListModel is a basic implementation of the ListModel interface
type SimpleListModel struct {
	items     []interface{}
	itemType  GType
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
	if w, ok := widget.(interface{ GetModel() TreeModel }); ok {
		return w.GetModel()
	}
	if w, ok := widget.(interface{ GetModel() ListModel }); ok {
		return w.GetModel()
	}
	if w, ok := widget.(interface{ GetModel() interface{} }); ok {
		return w.GetModel()
	}
	return nil
}

// SetModelForWidget sets a model for a widget, if applicable
func SetModelForWidget(widget interface{}, model interface{}) bool {
	// Try different model setter methods based on widget and model types
	if w, ok := widget.(interface{ SetModel(TreeModel) }); ok {
		if m, ok := model.(TreeModel); ok {
			w.SetModel(m)
			return true
		}
	}
	if w, ok := widget.(interface{ SetModel(ListModel) }); ok {
		if m, ok := model.(ListModel); ok {
			w.SetModel(m)
			return true
		}
	}
	if w, ok := widget.(interface{ SetModel(interface{}) bool }); ok {
		return w.SetModel(model)
	}
	return false
}

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