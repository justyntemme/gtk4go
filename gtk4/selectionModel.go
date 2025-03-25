// Package gtk4 provides selection model functionality for GTK4
// File: gtk4go/gtk4/selectionmodel.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Selection model callbacks
// extern void selectionChangedCallback(GtkSelectionModel *model, guint position, guint n_items, gpointer user_data);
//
// // Connect selection changed signal
// static gulong connectSelectionChanged(GtkSelectionModel *model, gpointer user_data) {
//     return g_signal_connect(model, "selection-changed", G_CALLBACK(selectionChangedCallback), user_data);
// }
//
// // SingleSelection operations
// static GtkSingleSelection* createSingleSelection(GListModel *model) {
//     return gtk_single_selection_new(model);
// }
//
// static void setSingleSelectionModel(GtkSingleSelection *selection, GListModel *model) {
//     gtk_single_selection_set_model(selection, model);
// }
//
// static void setSingleSelectionSelected(GtkSingleSelection *selection, guint position) {
//     gtk_single_selection_set_selected(selection, position);
// }
//
// static guint getSingleSelectionSelected(GtkSingleSelection *selection) {
//     return gtk_single_selection_get_selected(selection);
// }
//
// static void setSingleSelectionAutoselect(GtkSingleSelection *selection, gboolean autoselect) {
//     gtk_single_selection_set_autoselect(selection, autoselect);
// }
//
// static gboolean getSingleSelectionAutoselect(GtkSingleSelection *selection) {
//     return gtk_single_selection_get_autoselect(selection);
// }
//
// // MultiSelection operations
// static GtkMultiSelection* createMultiSelection(GListModel *model) {
//     return gtk_multi_selection_new(model);
// }
//
// static void setMultiSelectionModel(GtkMultiSelection *selection, GListModel *model) {
//     gtk_multi_selection_set_model(selection, model);
// }
//
// // NoSelection operations
// static GtkNoSelection* createNoSelection(GListModel *model) {
//     return gtk_no_selection_new(model);
// }
//
// static void setNoSelectionModel(GtkNoSelection *selection, GListModel *model) {
//     gtk_no_selection_set_model(selection, model);
// }
//
// // Common selection model operations
// static gboolean selectionModelIsSelected(GtkSelectionModel *model, guint position) {
//     return gtk_selection_model_is_selected(model, position);
// }
//
// static GtkBitset* selectionModelGetSelection(GtkSelectionModel *model) {
//     return gtk_selection_model_get_selection(model);
// }
//
// static void selectionModelSelectItem(GtkSelectionModel *model, guint position, gboolean unselect_rest) {
//     gtk_selection_model_select_item(model, position, unselect_rest);
// }
//
// static void selectionModelUnselectItem(GtkSelectionModel *model, guint position) {
//     gtk_selection_model_unselect_item(model, position);
// }
//
// static void selectionModelSelectRange(GtkSelectionModel *model, guint position, guint n_items, gboolean unselect_rest) {
//     gtk_selection_model_select_range(model, position, n_items, unselect_rest);
// }
//
// static void selectionModelUnselectRange(GtkSelectionModel *model, guint position, guint n_items) {
//     gtk_selection_model_unselect_range(model, position, n_items);
// }
//
// static void selectionModelSelectAll(GtkSelectionModel *model) {
//     gtk_selection_model_select_all(model);
// }
//
// static void selectionModelUnselectAll(GtkSelectionModel *model) {
//     gtk_selection_model_unselect_all(model);
// }
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// SelectionChangedCallback represents a callback for selection changes
type SelectionChangedCallback func(position, nItems int)

var (
	selectionCallbacks     = make(map[uintptr]SelectionChangedCallback)
	selectionCallbackMutex sync.RWMutex
)

//export selectionChangedCallback
func selectionChangedCallback(model *C.GtkSelectionModel, position, nItems C.guint, userData C.gpointer) {
	selectionCallbackMutex.RLock()
	defer selectionCallbackMutex.RUnlock()

	// Convert model pointer to uintptr for lookup
	modelPtr := uintptr(unsafe.Pointer(model))

	// Find and call the callback
	if callback, ok := selectionCallbacks[modelPtr]; ok {
		callback(int(position), int(nItems))
	}
}

// SelectionModel is an interface for GTK selection models
type SelectionModel interface {
	ListModel

	// GetSelectionModel returns the underlying GtkSelectionModel pointer
	GetSelectionModel() *C.GtkSelectionModel

	// IsSelected returns whether the item at the given position is selected
	IsSelected(position int) bool

	// SelectItem selects the item at the given position
	SelectItem(position int, unselectRest bool)

	// UnselectItem unselects the item at the given position
	UnselectItem(position int)

	// SelectRange selects a range of items
	SelectRange(position, nItems int, unselectRest bool)

	// UnselectRange unselects a range of items
	UnselectRange(position, nItems int)

	// SelectAll selects all items
	SelectAll()

	// UnselectAll unselects all items
	UnselectAll()

	// ConnectSelectionChanged connects a callback for selection changes
	ConnectSelectionChanged(callback SelectionChangedCallback)
}

// BaseSelectionModel provides common functionality for selection models
type BaseSelectionModel struct {
	BaseListModel
	selectionModel *C.GtkSelectionModel
	sourceModel    ListModel // The source model for this selection model
}

// GetSelectionModel returns the underlying GtkSelectionModel pointer
func (m *BaseSelectionModel) GetSelectionModel() *C.GtkSelectionModel {
	return m.selectionModel
}

// IsSelected returns whether the item at the given position is selected
func (m *BaseSelectionModel) IsSelected(position int) bool {
	return C.selectionModelIsSelected(m.selectionModel, C.guint(position)) != 0
}

// SelectItem selects the item at the given position
func (m *BaseSelectionModel) SelectItem(position int, unselectRest bool) {
	var cunselectRest C.gboolean
	if unselectRest {
		cunselectRest = C.TRUE
	} else {
		cunselectRest = C.FALSE
	}
	C.selectionModelSelectItem(m.selectionModel, C.guint(position), cunselectRest)
}

// UnselectItem unselects the item at the given position
func (m *BaseSelectionModel) UnselectItem(position int) {
	C.selectionModelUnselectItem(m.selectionModel, C.guint(position))
}

// SelectRange selects a range of items
func (m *BaseSelectionModel) SelectRange(position, nItems int, unselectRest bool) {
	var cunselectRest C.gboolean
	if unselectRest {
		cunselectRest = C.TRUE
	} else {
		cunselectRest = C.FALSE
	}
	C.selectionModelSelectRange(m.selectionModel, C.guint(position), C.guint(nItems), cunselectRest)
}

// UnselectRange unselects a range of items
func (m *BaseSelectionModel) UnselectRange(position, nItems int) {
	C.selectionModelUnselectRange(m.selectionModel, C.guint(position), C.guint(nItems))
}

// SelectAll selects all items
func (m *BaseSelectionModel) SelectAll() {
	C.selectionModelSelectAll(m.selectionModel)
}

// UnselectAll unselects all items
func (m *BaseSelectionModel) UnselectAll() {
	C.selectionModelUnselectAll(m.selectionModel)
}

// GetItem returns the item at the given position by delegating to the source model
func (m *BaseSelectionModel) GetItem(position int) interface{} {
	if m.sourceModel != nil {
		return m.sourceModel.GetItem(position)
	}
	// Fallback to the BaseListModel implementation if no source model
	return m.BaseListModel.GetItem(position)
}

// ConnectSelectionChanged connects a callback for selection changes
func (m *BaseSelectionModel) ConnectSelectionChanged(callback SelectionChangedCallback) {
	if callback == nil {
		return
	}

	selectionCallbackMutex.Lock()
	defer selectionCallbackMutex.Unlock()

	// Store the callback in the map
	modelPtr := uintptr(unsafe.Pointer(m.selectionModel))
	selectionCallbacks[modelPtr] = callback

	// Connect the signal
	C.connectSelectionChanged(m.selectionModel, C.gpointer(unsafe.Pointer(m.selectionModel)))
}

// Destroy frees resources associated with the selection model
func (m *BaseSelectionModel) Destroy() {
	if m.selectionModel != nil {
		selectionCallbackMutex.Lock()
		delete(selectionCallbacks, uintptr(unsafe.Pointer(m.selectionModel)))
		selectionCallbackMutex.Unlock()
	}

	m.BaseListModel.Destroy()
	m.selectionModel = nil
	m.sourceModel = nil
}

// SingleSelection is a selection model that allows selecting a single item
type SingleSelection struct {
	BaseSelectionModel
	singleSelection *C.GtkSingleSelection
}

// SingleSelectionOption is a function that configures a single selection
type SingleSelectionOption func(*SingleSelection)

// NewSingleSelection creates a new single selection model
func NewSingleSelection(model ListModel, options ...SingleSelectionOption) *SingleSelection {
	var singleSelection *C.GtkSingleSelection
	if model != nil {
		singleSelection = C.createSingleSelection(model.GetListModel())
	} else {
		singleSelection = C.createSingleSelection(nil)
	}

	selection := &SingleSelection{
		BaseSelectionModel: BaseSelectionModel{
			BaseListModel: BaseListModel{
				model: (*C.GListModel)(unsafe.Pointer(singleSelection)),
			},
			selectionModel: (*C.GtkSelectionModel)(unsafe.Pointer(singleSelection)),
			sourceModel:    model,
		},
		singleSelection: singleSelection,
	}

	// Apply options
	for _, option := range options {
		option(selection)
	}

	runtime.SetFinalizer(selection, (*SingleSelection).Destroy)
	return selection
}

// WithAutoselect sets whether the selection should automatically select an item
func WithAutoselect(autoselect bool) SingleSelectionOption {
	return func(s *SingleSelection) {
		var cautoselect C.gboolean
		if autoselect {
			cautoselect = C.TRUE
		} else {
			cautoselect = C.FALSE
		}
		C.setSingleSelectionAutoselect(s.singleSelection, cautoselect)
	}
}

// WithInitialSelection sets the initially selected item
func WithInitialSelection(position int) SingleSelectionOption {
	return func(s *SingleSelection) {
		C.setSingleSelectionSelected(s.singleSelection, C.guint(position))
	}
}

// SetModel sets the model for the selection
func (s *SingleSelection) SetModel(model ListModel) {
	if model != nil {
		C.setSingleSelectionModel(s.singleSelection, model.GetListModel())
		s.sourceModel = model
	} else {
		C.setSingleSelectionModel(s.singleSelection, nil)
		s.sourceModel = nil
	}
}

// GetSelected returns the position of the selected item
func (s *SingleSelection) GetSelected() int {
	return int(C.getSingleSelectionSelected(s.singleSelection))
}

// SetSelected sets the selected item
func (s *SingleSelection) SetSelected(position int) {
	C.setSingleSelectionSelected(s.singleSelection, C.guint(position))
}

// SetAutoselect sets whether the selection should automatically select an item
func (s *SingleSelection) SetAutoselect(autoselect bool) {
	var cautoselect C.gboolean
	if autoselect {
		cautoselect = C.TRUE
	} else {
		cautoselect = C.FALSE
	}
	C.setSingleSelectionAutoselect(s.singleSelection, cautoselect)
}

// GetAutoselect returns whether the selection automatically selects an item
func (s *SingleSelection) GetAutoselect() bool {
	return C.getSingleSelectionAutoselect(s.singleSelection) != 0
}

// GetItem delegates to the source model to get an item at a specific position
func (s *SingleSelection) GetItem(position int) interface{} {
	return s.BaseSelectionModel.GetItem(position)
}

// Destroy frees resources associated with the single selection
func (s *SingleSelection) Destroy() {
	s.BaseSelectionModel.Destroy()
	s.singleSelection = nil
}

// MultiSelection is a selection model that allows selecting multiple items
type MultiSelection struct {
	BaseSelectionModel
	multiSelection *C.GtkMultiSelection
}

// NewMultiSelection creates a new multi-selection model
func NewMultiSelection(model ListModel) *MultiSelection {
	var multiSelection *C.GtkMultiSelection
	if model != nil {
		multiSelection = C.createMultiSelection(model.GetListModel())
	} else {
		multiSelection = C.createMultiSelection(nil)
	}

	selection := &MultiSelection{
		BaseSelectionModel: BaseSelectionModel{
			BaseListModel: BaseListModel{
				model: (*C.GListModel)(unsafe.Pointer(multiSelection)),
			},
			selectionModel: (*C.GtkSelectionModel)(unsafe.Pointer(multiSelection)),
			sourceModel:    model,
		},
		multiSelection: multiSelection,
	}

	runtime.SetFinalizer(selection, (*MultiSelection).Destroy)
	return selection
}

// SetModel sets the model for the selection
func (s *MultiSelection) SetModel(model ListModel) {
	if model != nil {
		C.setMultiSelectionModel(s.multiSelection, model.GetListModel())
		s.sourceModel = model
	} else {
		C.setMultiSelectionModel(s.multiSelection, nil)
		s.sourceModel = nil
	}
}

// GetItem delegates to the source model to get an item at a specific position
func (s *MultiSelection) GetItem(position int) interface{} {
	return s.BaseSelectionModel.GetItem(position)
}

// Destroy frees resources associated with the multi selection
func (s *MultiSelection) Destroy() {
	s.BaseSelectionModel.Destroy()
	s.multiSelection = nil
}

// NoSelection is a selection model that doesn't allow selecting items
type NoSelection struct {
	BaseSelectionModel
	noSelection *C.GtkNoSelection
}

// NewNoSelection creates a new no-selection model
func NewNoSelection(model ListModel) *NoSelection {
	var noSelection *C.GtkNoSelection
	if model != nil {
		noSelection = C.createNoSelection(model.GetListModel())
	} else {
		noSelection = C.createNoSelection(nil)
	}

	selection := &NoSelection{
		BaseSelectionModel: BaseSelectionModel{
			BaseListModel: BaseListModel{
				model: (*C.GListModel)(unsafe.Pointer(noSelection)),
			},
			selectionModel: (*C.GtkSelectionModel)(unsafe.Pointer(noSelection)),
			sourceModel:    model,
		},
		noSelection: noSelection,
	}

	runtime.SetFinalizer(selection, (*NoSelection).Destroy)
	return selection
}

// SetModel sets the model for the selection
func (s *NoSelection) SetModel(model ListModel) {
	if model != nil {
		C.setNoSelectionModel(s.noSelection, model.GetListModel())
		s.sourceModel = model
	} else {
		C.setNoSelectionModel(s.noSelection, nil)
		s.sourceModel = nil
	}
}

// GetItem delegates to the source model to get an item at a specific position
func (s *NoSelection) GetItem(position int) interface{} {
	return s.BaseSelectionModel.GetItem(position)
}

// Destroy frees resources associated with the no selection
func (s *NoSelection) Destroy() {
	s.BaseSelectionModel.Destroy()
	s.noSelection = nil
}