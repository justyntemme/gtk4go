// Package gtk4 provides list item factory functionality for GTK4
// File: gtk4go/gtk4/listItemFactory.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback functions for list item factory
// extern void listItemSetupCallback(GtkListItemFactory *factory, GtkListItem *list_item, gpointer user_data);
// extern void listItemBindCallback(GtkListItemFactory *factory, GtkListItem *list_item, gpointer user_data);
// extern void listItemUnbindCallback(GtkListItemFactory *factory, GtkListItem *list_item, gpointer user_data);
// extern void listItemTeardownCallback(GtkListItemFactory *factory, GtkListItem *list_item, gpointer user_data);
//
// // Create a signal list item factory
// static GtkListItemFactory* create_signal_list_item_factory() {
//     return gtk_signal_list_item_factory_new();
// }
//
// // Connect factory signals
// static void connect_factory_signals(GtkSignalListItemFactory *factory, gpointer user_data) {
//     if (factory == NULL) return;
//     g_signal_connect(factory, "setup", G_CALLBACK(listItemSetupCallback), user_data);
//     g_signal_connect(factory, "bind", G_CALLBACK(listItemBindCallback), user_data);
//     g_signal_connect(factory, "unbind", G_CALLBACK(listItemUnbindCallback), user_data);
//     g_signal_connect(factory, "teardown", G_CALLBACK(listItemTeardownCallback), user_data);
// }
//
// // GtkListItem helpers
// static void list_item_set_child(GtkListItem *list_item, GtkWidget *child) {
//     if (list_item == NULL || child == NULL) return;
//     gtk_list_item_set_child(list_item, child);
// }
//
// static GtkWidget* list_item_get_child(GtkListItem *list_item) {
//     if (list_item == NULL) return NULL;
//     return gtk_list_item_get_child(list_item);
// }
//
// static void list_item_set_activatable(GtkListItem *list_item, gboolean activatable) {
//     if (list_item == NULL) return;
//     gtk_list_item_set_activatable(list_item, activatable);
// }
//
// static void list_item_set_selectable(GtkListItem *list_item, gboolean selectable) {
//     if (list_item == NULL) return;
//     gtk_list_item_set_selectable(list_item, selectable);
// }
//
// static guint list_item_get_position(GtkListItem *list_item) {
//     if (list_item == NULL) return 0;
//     return gtk_list_item_get_position(list_item);
// }
//
// static gpointer list_item_get_item(GtkListItem *list_item) {
//     if (list_item == NULL) return NULL;
//     return gtk_list_item_get_item(list_item);
// }
//
// // Creates a label widget
// static GtkWidget* create_label(const char* text) {
//     return gtk_label_new(text ? text : "");
// }
//
// // Creates a check button widget
// static GtkWidget* create_check_button() {
//     return gtk_check_button_new();
// }
//
// // Set check button active state
// static void check_button_set_active(GtkCheckButton *button, gboolean active) {
//     if (button == NULL) return;
//     gtk_check_button_set_active(button, active);
// }
//
// // Creates a progress bar widget
// static GtkWidget* create_progress_bar() {
//     return gtk_progress_bar_new();
// }
//
// // Set progress bar fraction
// static void progress_bar_set_fraction(GtkProgressBar *bar, double fraction) {
//     if (bar == NULL) return;
//     gtk_progress_bar_set_fraction(bar, fraction);
// }
//
// // Creates an image widget
// static GtkWidget* create_image() {
//     return gtk_image_new();
// }
import "C"

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

// ListItemFactoryCallbacks holds callbacks for list item factory events
type ListItemFactoryCallbacks struct {
	Setup    func(item *ListItem)
	Bind     func(item *ListItem)
	Unbind   func(item *ListItem)
	Teardown func(item *ListItem)
}

var (
	factoryCallbacks     = make(map[uintptr]*ListItemFactoryCallbacks)
	factoryCallbackMutex sync.RWMutex
)

//export listItemSetupCallback
func listItemSetupCallback(factory *C.GtkListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	factoryCallbackMutex.RLock()
	defer factoryCallbackMutex.RUnlock()

	// Convert factory pointer to uintptr for lookup
	factoryPtr := uintptr(unsafe.Pointer(factory))

	// Find callbacks
	callbacks, ok := factoryCallbacks[factoryPtr]
	if !ok || callbacks.Setup == nil {
		return
	}

	// Create Go wrapper and call the callback
	item := &ListItem{listItem: listItem}
	callbacks.Setup(item)
}

//export listItemBindCallback
func listItemBindCallback(factory *C.GtkListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	factoryCallbackMutex.RLock()
	defer factoryCallbackMutex.RUnlock()

	// Convert factory pointer to uintptr for lookup
	factoryPtr := uintptr(unsafe.Pointer(factory))

	// Find callbacks
	callbacks, ok := factoryCallbacks[factoryPtr]
	if !ok || callbacks.Bind == nil {
		return
	}

	// Create Go wrapper and call the callback
	item := &ListItem{listItem: listItem}
	callbacks.Bind(item)
}

//export listItemUnbindCallback
func listItemUnbindCallback(factory *C.GtkListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	factoryCallbackMutex.RLock()
	defer factoryCallbackMutex.RUnlock()

	// Convert factory pointer to uintptr for lookup
	factoryPtr := uintptr(unsafe.Pointer(factory))

	// Find callbacks
	callbacks, ok := factoryCallbacks[factoryPtr]
	if !ok || callbacks.Unbind == nil {
		return
	}

	// Create Go wrapper and call the callback
	item := &ListItem{listItem: listItem}
	callbacks.Unbind(item)
}

//export listItemTeardownCallback
func listItemTeardownCallback(factory *C.GtkListItemFactory, listItem *C.GtkListItem, userData C.gpointer) {
	factoryCallbackMutex.RLock()
	defer factoryCallbackMutex.RUnlock()

	// Convert factory pointer to uintptr for lookup
	factoryPtr := uintptr(unsafe.Pointer(factory))

	// Find callbacks
	callbacks, ok := factoryCallbacks[factoryPtr]
	if !ok || callbacks.Teardown == nil {
		return
	}

	// Create Go wrapper and call the callback
	item := &ListItem{listItem: listItem}
	callbacks.Teardown(item)
}

// ListItemFactory represents a GTK list item factory
type ListItemFactory struct {
	factory *C.GtkListItemFactory
}

// NewSignalListItemFactory creates a new signal list item factory
func NewSignalListItemFactory(callbacks *ListItemFactoryCallbacks) *ListItemFactory {
	fmt.Println("NewSignalListItemFactory: Creating factory...")
	
	// Create the factory
	factoryPtr := C.create_signal_list_item_factory()
	if factoryPtr == nil {
		fmt.Println("NewSignalListItemFactory: ERROR - Failed to create GTK factory")
		return nil
	}
	
	factory := &ListItemFactory{
		factory: factoryPtr,
	}

	// Store callbacks if provided
	if callbacks != nil {
		factoryCallbackMutex.Lock()
		ptrKey := uintptr(unsafe.Pointer(factory.factory))
		factoryCallbacks[ptrKey] = callbacks
		factoryCallbackMutex.Unlock()

		// Connect signals
		C.connect_factory_signals(
			(*C.GtkSignalListItemFactory)(unsafe.Pointer(factory.factory)),
			C.gpointer(unsafe.Pointer(factory.factory)),
		)
	}

	runtime.SetFinalizer(factory, (*ListItemFactory).Free)
	fmt.Printf("NewSignalListItemFactory: Factory created successfully: %v\n", factory)
	return factory
}

// GetFactory returns the underlying GtkListItemFactory pointer
func (f *ListItemFactory) GetFactory() *C.GtkListItemFactory {
	return f.factory
}

// Free frees the factory
func (f *ListItemFactory) Free() {
	fmt.Printf("ListItemFactory.Free: Cleaning up factory %v\n", f)
	
	if f.factory != nil {
		// Remove callbacks
		factoryCallbackMutex.Lock()
		delete(factoryCallbacks, uintptr(unsafe.Pointer(f.factory)))
		factoryCallbackMutex.Unlock()

		C.g_object_unref(C.gpointer(unsafe.Pointer(f.factory)))
		f.factory = nil
	}
}

// ListItem is a wrapper around a GTK list item
type ListItem struct {
	listItem *C.GtkListItem
}

// SetChild sets the child widget of the list item
func (i *ListItem) SetChild(child Widget) {
	fmt.Printf("ListItem.SetChild: Setting child=%v\n", child)
	
	if i.listItem == nil {
		fmt.Println("ListItem.SetChild: WARNING - i.listItem is nil")
		return
	}
	
	if child == nil {
		fmt.Println("ListItem.SetChild: WARNING - child is nil")
		return
	}
	
	C.list_item_set_child(i.listItem, child.GetWidget())
}

// GetChild gets the child widget of the list item
func (i *ListItem) GetChild() Widget {
	if i.listItem == nil {
		fmt.Println("ListItem.GetChild: WARNING - i.listItem is nil")
		return nil
	}
	
	widget := C.list_item_get_child(i.listItem)
	if widget == nil {
		return nil
	}

	// Basic wrapper - in a real implementation we'd detect the widget type
	return &BaseWidget{widget: widget}
}

// SetActivatable sets whether the item is activatable
func (i *ListItem) SetActivatable(activatable bool) {
	if i.listItem == nil {
		return
	}
	
	var cActivatable C.gboolean
	if activatable {
		cActivatable = C.TRUE
	} else {
		cActivatable = C.FALSE
	}
	C.list_item_set_activatable(i.listItem, cActivatable)
}

// SetSelectable sets whether the item is selectable
func (i *ListItem) SetSelectable(selectable bool) {
	if i.listItem == nil {
		return
	}
	
	var cSelectable C.gboolean
	if selectable {
		cSelectable = C.TRUE
	} else {
		cSelectable = C.FALSE
	}
	C.list_item_set_selectable(i.listItem, cSelectable)
}

// GetPosition gets the position of the item in the model
func (i *ListItem) GetPosition() int {
	if i.listItem == nil {
		return -1
	}
	
	return int(C.list_item_get_position(i.listItem))
}

// GetItem gets the model item represented by this list item
func (i *ListItem) GetItem() interface{} {
	if i.listItem == nil {
		return nil
	}
	
	item := C.list_item_get_item(i.listItem)
	if item == nil {
		return nil
	}

	// Try to get string from the item if it's a GObject
	str := GetStringFromObject((*C.GObject)(item))
	if str != "" {
		return str
	}

	// Otherwise return the pointer
	return uintptr(unsafe.Pointer(item))
}

// ListItemFactoryOption is a function that configures a list item factory
type ListItemFactoryOption func(*ListItemFactory)

// WithSetupCallback sets the setup callback
func WithSetupCallback(callback func(item *ListItem)) ListItemFactoryOption {
	return func(factory *ListItemFactory) {
		factoryCallbackMutex.Lock()
		defer factoryCallbackMutex.Unlock()

		factoryPtr := uintptr(unsafe.Pointer(factory.factory))
		callbacks, ok := factoryCallbacks[factoryPtr]
		if !ok {
			callbacks = &ListItemFactoryCallbacks{}
			factoryCallbacks[factoryPtr] = callbacks
		}
		callbacks.Setup = callback
	}
}

// WithBindCallback sets the bind callback
func WithBindCallback(callback func(item *ListItem)) ListItemFactoryOption {
	return func(factory *ListItemFactory) {
		factoryCallbackMutex.Lock()
		defer factoryCallbackMutex.Unlock()

		factoryPtr := uintptr(unsafe.Pointer(factory.factory))
		callbacks, ok := factoryCallbacks[factoryPtr]
		if !ok {
			callbacks = &ListItemFactoryCallbacks{}
			factoryCallbacks[factoryPtr] = callbacks
		}
		callbacks.Bind = callback
	}
}

// WithUnbindCallback sets the unbind callback
func WithUnbindCallback(callback func(item *ListItem)) ListItemFactoryOption {
	return func(factory *ListItemFactory) {
		factoryCallbackMutex.Lock()
		defer factoryCallbackMutex.Unlock()

		factoryPtr := uintptr(unsafe.Pointer(factory.factory))
		callbacks, ok := factoryCallbacks[factoryPtr]
		if !ok {
			callbacks = &ListItemFactoryCallbacks{}
			factoryCallbacks[factoryPtr] = callbacks
		}
		callbacks.Unbind = callback
	}
}

// WithTeardownCallback sets the teardown callback
func WithTeardownCallback(callback func(item *ListItem)) ListItemFactoryOption {
	return func(factory *ListItemFactory) {
		factoryCallbackMutex.Lock()
		defer factoryCallbackMutex.Unlock()

		factoryPtr := uintptr(unsafe.Pointer(factory.factory))
		callbacks, ok := factoryCallbacks[factoryPtr]
		if !ok {
			callbacks = &ListItemFactoryCallbacks{}
			factoryCallbacks[factoryPtr] = callbacks
		}
		callbacks.Teardown = callback
	}
}

// Helper functions for creating common widget types for list items

// CreateLabel creates a new label widget with the given text
func CreateLabel(text string) *Label {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	widget := C.create_label(cText)
	return &Label{
		BaseWidget: BaseWidget{
			widget: widget,
		},
	}
}

// CreateCheckButton creates a new check button widget
func CreateCheckButton() *CheckButton {
	widget := C.create_check_button()
	return &CheckButton{
		BaseWidget: BaseWidget{
			widget: widget,
		},
	}
}

// SetActive sets whether the check button is active
func (b *CheckButton) SetActive(active bool) {
	var cActive C.gboolean
	if active {
		cActive = C.TRUE
	} else {
		cActive = C.FALSE
	}
	C.check_button_set_active((*C.GtkCheckButton)(unsafe.Pointer(b.widget)), cActive)
}

// CheckButton represents a GTK check button
type CheckButton struct {
	BaseWidget
}

// CreateProgressBar creates a new progress bar widget
func CreateProgressBar() *ProgressBar {
	widget := C.create_progress_bar()
	return &ProgressBar{
		BaseWidget: BaseWidget{
			widget: widget,
		},
	}
}

// ProgressBar represents a GTK progress bar
type ProgressBar struct {
	BaseWidget
}

// SetFraction sets the fraction (0.0 to 1.0) of the progress bar
func (p *ProgressBar) SetFraction(fraction float64) {
	C.progress_bar_set_fraction((*C.GtkProgressBar)(unsafe.Pointer(p.widget)), C.double(fraction))
}

// CreateImage creates a new image widget
func CreateImage() *Image {
	widget := C.create_image()
	return &Image{
		BaseWidget: BaseWidget{
			widget: widget,
		},
	}
}

// Image represents a GTK image widget
type Image struct {
	BaseWidget
}

// Factory helpers for common list item factory patterns

// TextFactory creates a factory that displays text strings
func TextFactory() *ListItemFactory {
	fmt.Println("TextFactory: Creating text factory...")
	
	callbacks := &ListItemFactoryCallbacks{
		Setup: func(item *ListItem) {
			fmt.Println("TextFactory.Setup: Creating label for item")
			// Create a label for this item
			label := CreateLabel("")
			if label == nil {
				fmt.Println("TextFactory.Setup: ERROR - Failed to create label")
				return
			}
			
			fmt.Printf("TextFactory.Setup: Setting label as child for item at position %d\n", 
				item.GetPosition())
			item.SetChild(label)
		},
		Bind: func(item *ListItem) {
			fmt.Printf("TextFactory.Bind: Binding item at position %d\n", item.GetPosition())
			
			// Get the label
			child := item.GetChild()
			if child == nil {
				fmt.Println("TextFactory.Bind: WARNING - item has no child")
				return
			}
			
			label, ok := child.(*Label)
			if !ok {
				fmt.Println("TextFactory.Bind: WARNING - child is not a label")
				return
			}
			
			// Get the text from the item
			modelItem := item.GetItem()
			fmt.Printf("TextFactory.Bind: Model item=%v\n", modelItem)
			
			// Convert to string if possible
			var text string
			switch v := modelItem.(type) {
			case string:
				text = v
				fmt.Printf("TextFactory.Bind: Got string value: %q\n", text)
			case uintptr:
				// Try to get string from GObject
				gobj := (*C.GObject)(unsafe.Pointer(v))
				str := GetStringFromObject(gobj)
				if str != "" {
					text = str
					fmt.Printf("TextFactory.Bind: Got string from GObject: %q\n", text)
				} else {
					text = fmt.Sprintf("Item %d", item.GetPosition())
					fmt.Printf("TextFactory.Bind: Using default text: %q\n", text)
				}
			default:
				// Handle other types or use a default
				text = fmt.Sprintf("Item %d", item.GetPosition())
				fmt.Printf("TextFactory.Bind: Using default text for unknown type: %q\n", text)
			}
			
			// Set the label text
			fmt.Printf("TextFactory.Bind: Setting label text to %q\n", text)
			label.SetText(text)
		},
		Unbind: func(item *ListItem) {
			fmt.Printf("TextFactory.Unbind: Unbinding item at position %d\n", item.GetPosition())
		},
		Teardown: func(item *ListItem) {
			fmt.Printf("TextFactory.Teardown: Tearing down item at position %d\n", item.GetPosition())
		},
	}
	
	factory := NewSignalListItemFactory(callbacks)
	fmt.Printf("TextFactory: Factory created: %v\n", factory)
	if factory == nil {
		fmt.Println("TextFactory: ERROR - Failed to create factory")
	} else if factory.factory == nil {
		fmt.Println("TextFactory: ERROR - Factory has nil factory pointer")
	}
	
	return factory
}

// CheckboxFactory creates a factory that displays checkboxes
func CheckboxFactory() *ListItemFactory {
	fmt.Println("CheckboxFactory: Creating checkbox factory...")
	
	callbacks := &ListItemFactoryCallbacks{
		Setup: func(item *ListItem) {
			fmt.Println("CheckboxFactory.Setup: Creating checkbox for item")
			// Create a check button for this item
			checkButton := CreateCheckButton()
			if checkButton == nil {
				fmt.Println("CheckboxFactory.Setup: ERROR - Failed to create check button")
				return
			}
			
			fmt.Printf("CheckboxFactory.Setup: Setting check button as child for item at position %d\n", 
				item.GetPosition())
			item.SetChild(checkButton)
		},
		Bind: func(item *ListItem) {
			fmt.Printf("CheckboxFactory.Bind: Binding item at position %d\n", item.GetPosition())
			
			// Get the check button
			child := item.GetChild()
			if child == nil {
				fmt.Println("CheckboxFactory.Bind: WARNING - item has no child")
				return
			}
			
			checkButton, ok := child.(*CheckButton)
			if !ok {
				fmt.Println("CheckboxFactory.Bind: WARNING - child is not a check button")
				return
			}
			
			// Get the value from the item
			modelItem := item.GetItem()
			fmt.Printf("CheckboxFactory.Bind: Model item=%v\n", modelItem)
			
			// Convert to bool if possible
			var active bool
			if b, ok := modelItem.(bool); ok {
				active = b
				fmt.Printf("CheckboxFactory.Bind: Got boolean value: %v\n", active)
			} else {
				// Default to false for unknown types
				active = false
				fmt.Printf("CheckboxFactory.Bind: Using default value (false) for unknown type\n")
			}
			
			// Set the check button state
			fmt.Printf("CheckboxFactory.Bind: Setting check button active state to %v\n", active)
			checkButton.SetActive(active)
		},
	}
	
	factory := NewSignalListItemFactory(callbacks)
	fmt.Printf("CheckboxFactory: Factory created: %v\n", factory)
	return factory
}

// ProgressFactory creates a factory that displays progress bars
func ProgressFactory() *ListItemFactory {
	fmt.Println("ProgressFactory: Creating progress factory...")
	
	callbacks := &ListItemFactoryCallbacks{
		Setup: func(item *ListItem) {
			fmt.Println("ProgressFactory.Setup: Creating progress bar for item")
			// Create a progress bar for this item
			progressBar := CreateProgressBar()
			if progressBar == nil {
				fmt.Println("ProgressFactory.Setup: ERROR - Failed to create progress bar")
				return
			}
			
			fmt.Printf("ProgressFactory.Setup: Setting progress bar as child for item at position %d\n", 
				item.GetPosition())
			item.SetChild(progressBar)
		},
		Bind: func(item *ListItem) {
			fmt.Printf("ProgressFactory.Bind: Binding item at position %d\n", item.GetPosition())
			
			// Get the progress bar
			child := item.GetChild()
			if child == nil {
				fmt.Println("ProgressFactory.Bind: WARNING - item has no child")
				return
			}
			
			progressBar, ok := child.(*ProgressBar)
			if !ok {
				fmt.Println("ProgressFactory.Bind: WARNING - child is not a progress bar")
				return
			}
			
			// Get the value from the item
			modelItem := item.GetItem()
			fmt.Printf("ProgressFactory.Bind: Model item=%v\n", modelItem)
			
			// Convert to float if possible
			var fraction float64
			switch v := modelItem.(type) {
			case float64:
				fraction = v
				fmt.Printf("ProgressFactory.Bind: Got float64 value: %v\n", fraction)
			case float32:
				fraction = float64(v)
				fmt.Printf("ProgressFactory.Bind: Got float32 value: %v\n", fraction)
			case int:
				fraction = float64(v) / 100.0 // Assuming percentage
				fmt.Printf("ProgressFactory.Bind: Got int value: %v (converted to %v)\n", v, fraction)
			default:
				// Default to 0 for unknown types
				fraction = 0
				fmt.Printf("ProgressFactory.Bind: Using default value (0) for unknown type\n")
			}
			
			// Ensure value is between 0 and 1
			if fraction < 0 {
				fraction = 0
			} else if fraction > 1 {
				fraction = 1
			}
			
			// Set the progress bar fraction
			fmt.Printf("ProgressFactory.Bind: Setting progress bar fraction to %v\n", fraction)
			progressBar.SetFraction(fraction)
		},
	}
	
	factory := NewSignalListItemFactory(callbacks)
	fmt.Printf("ProgressFactory: Factory created: %v\n", factory)
	return factory
}