# Unified Callback System (UCS) in GTK4Go

## Overview

The Unified Callback System (UCS) is a core component of GTK4Go that provides a centralized, type-safe, and memory-safe way to handle GTK signal callbacks. This system solves several challenges when binding GTK4 (a C library) to Go:

- **Thread safety**: Ensuring GTK UI operations run on the main thread
- **Type safety**: Properly mapping different callback signatures between Go and C
- **Memory management**: Preventing leaks by properly tracking and disconnecting signal handlers
- **Context management**: Maintaining proper object and callback references

The UCS allows Go developers to connect signals to callbacks using an idiomatic Go API while handling the complexities of C-Go interoperability behind the scenes.

## Architecture

The UCS is implemented primarily in `gtk4/callbacks.go` and consists of the following core components:

### Core Components

1. **CallbackManager**: Central manager that maintains maps for tracking callbacks, object handlers, and signal connections
2. **Signal Registration**: Functions for connecting signals to callbacks and generating unique IDs
3. **C-Side Handlers**: C callback functions that bridge GTK signals to Go callbacks
4. **Thread-Safe Execution**: Mechanism to ensure callbacks run on the UI thread
5. **Cleanup Handlers**: Functions to properly disconnect and clean up signal handlers

### Key Data Structures

- **callbacks**: Maps callback IDs to callback data
- **objectHandlers**: Maps object pointers to lists of handler IDs
- **objectCallbacks**: Maps object pointers to signal types to callbacks

## Implementation Details

### Callback Registration Flow

1. A user calls `Connect(object, signal, callback)` on a widget
2. The UCS:
   - Generates a unique ID for the callback
   - Analyzes the callback signature
   - Stores callback data in the `callbacks` map
   - Connects the appropriate C signal handler
   - Tracks the handler for cleanup

### Signal Emission Flow

1. GTK emits a signal
2. The C callback handler receives the signal with the callback ID
3. The handler looks up the Go callback in the `callbacks` map
4. The callback is executed on the UI thread using `execCallback()`

### Cleanup Process

1. When `Disconnect()` or `DisconnectAll()` is called:
   - The signal handler is disconnected using GTK's C API
   - Callback references are removed from all maps
   - Object handler tracking is updated

## How Widgets Use the UCS

### Button

```go
// ConnectClicked connects a callback function to the button's "clicked" signal
func (b *Button) ConnectClicked(callback func()) {
    Connect(b, SignalClicked, callback)
}

// DisconnectClicked disconnects all clicked signal handlers
func (b *Button) DisconnectClicked() {
    DisconnectAll(b)
}
```

The Button widget uses UCS to handle its "clicked" signal. It provides a simple `ConnectClicked` method that wraps the UCS `Connect` function with the appropriate signal type. When the button is clicked, GTK emits the signal, and the UCS ensures the callback executes on the UI thread.

### Entry

```go
// ConnectChanged connects a callback function to the entry's "changed" signal
func (e *Entry) ConnectChanged(callback func()) {
    Connect(e, SignalChanged, callback)
}

// ConnectActivate connects a callback function to the entry's "activate" signal
func (e *Entry) ConnectActivate(callback func()) {
    Connect(e, SignalActivate, callback)
}
```

The Entry widget connects to both "changed" and "activate" signals, allowing callbacks when the text changes or when the user presses Enter. Both signals are managed through the UCS.

### Window

```go
// ConnectCloseRequest connects a callback function to the window's "close-request" signal
func (w *Window) ConnectCloseRequest(callback func() bool) uint64 {
    return Connect(w, SignalCloseRequest, callback)
}
```

The Window widget uses the UCS to handle window close requests. The callback returns a boolean value, which the UCS correctly handles using a specialized C handler that preserves the return value.

### ListView

```go
// ConnectActivate connects a callback for item activation
func (lv *ListView) ConnectActivate(callback ListViewActivateCallback) {
    // Convert to a regular func(int) for the callback handler
    standardCallback := func(position int) {
        callback(position)
    }
    Connect(lv, SignalListActivate, standardCallback)
}
```

ListView uses the UCS for its "activate" signal with a position parameter, demonstrating how UCS handles callbacks with parameters.

### Dialog

```go
// ConnectResponse connects a response callback to the dialog
func (d *Dialog) ConnectResponse(callback DialogResponseCallback) {
    standardCallback := func(responseId ResponseType) {
        callback(responseId)
    }
    // Use direct callback storage for dialog responses
    dialogPtr := uintptr(unsafe.Pointer(d.widget))
    StoreDirectCallback(dialogPtr, SignalDialogResponse, standardCallback)
}
```

Dialogs use a special direct callback storage method within the UCS to handle response signals for buttons clicked in the dialog.

### Adjustment

```go
// ConnectValueChanged connects a callback to the value-changed signal
func (a *Adjustment) ConnectValueChanged(callback func()) {
    Connect(a, SignalValueChanged, callback)
}
```

The Adjustment widget connects to the "value-changed" signal to notify when its value changes.

### ListItemFactory

```go
// ConnectSetup connects a callback for the setup signal
func (f *SignalListItemFactory) ConnectSetup(callback ListItemCallback) {
    Connect(f, SignalSetup, callback)
}

// ConnectBind connects a callback for the bind signal
func (f *SignalListItemFactory) ConnectBind(callback ListItemCallback) {
    Connect(f, SignalBind, callback)
}
```

ListItemFactory uses the UCS to handle the complex lifecycle of list items, with specialized callbacks for setup, bind, unbind, and teardown phases.

### SelectionModel

```go
// ConnectSelectionChanged connects a callback for selection changes
func (m *BaseSelectionModel) ConnectSelectionChanged(callback SelectionChangedCallback) {
    stdCallback := func(position, nItems int) {
        callback(position, nItems)
    }
    modelPtr := uintptr(unsafe.Pointer(m.selectionModel))
    handlerID := C.connectSelectionChanged(m.selectionModel, C.gpointer(unsafe.Pointer(m.selectionModel)))
    StoreCallback(modelPtr, SignalSelectionChanged, stdCallback, handlerID)
}
```

SelectionModel uses the UCS with a specialized storage method for handling selection changes.

### Menu Components

```go
// ConnectDeactivate connects a callback for when the popover is closed
func (pm *PopoverMenu) ConnectDeactivate(callback func()) {
    Connect(pm, SignalDeactivate, callback)
}
```

Menu components use the UCS to handle signals like deactivation when a menu is closed.

### Action Components

Actions use a specialized part of the UCS system with direct callback storage:

```go
// For Action activation
actionPtr := uintptr(unsafe.Pointer(action))
StoreDirectCallback(actionPtr, SignalActionActivate, standardCallback)
```

The UCS handles the unique requirements of action activation with direct pointer matching.

## Advanced Features

### Signal Sources and Disambiguation

The UCS can handle signals with the same name from different sources by tracking the signal source:

```go
// Determine signal source based on object type and signal
source := SourceGeneric
if _, isListView := object.(*ListView); isListView && signal == SignalListActivate {
    source = SourceListView
} else if _, isAction := object.(*Action); isAction && signal == SignalActionActivate {
    source = SourceAction
}
```

This allows proper handling of "activate" signals from both ListView and Action components.

### Memory Management

The UCS automatically tracks handlers to ensure proper cleanup:

```go
// When a widget is destroyed
func (w *Widget) Destroy() {
    DisconnectAll(w) // Automatically disconnects all signals
    // Continue with widget destruction
}
```

This prevents memory leaks and dangling references.

### Thread Safety

The UCS ensures callbacks run on the UI thread using the `execCallback` function:

```go
// execCallback safely executes a callback on the main UI thread
func execCallback(callback interface{}, args ...interface{}) {
    uithread.RunOnUIThread(func() {
        // Execute the callback based on its type
        // ...
    })
}
```

This maintains GTK's thread safety requirements.

## Best Practices

1. **Always call DisconnectAll in Destroy methods**:
   ```go
   func (w *Widget) Destroy() {
       DisconnectAll(w)
       // Continue with destruction
   }
   ```

2. **Return proper handler IDs**:
   ```go
   // Store and return the handler ID for manual disconnection
   func ConnectSignal(callback func()) uint64 {
       return Connect(widget, signal, callback)
   }
   ```

3. **Use type-specific Connect methods**:
   Provide widget-specific methods like `ConnectClicked` rather than asking users to directly use `Connect`.

## Conclusion

The Unified Callback System is a core component of GTK4Go that manages the complex task of bridging GTK signals to Go callbacks. It ensures thread safety, type safety, and memory safety while providing an idiomatic Go API. Understanding how it works helps developers use GTK4Go effectively and extend it with new widget types.