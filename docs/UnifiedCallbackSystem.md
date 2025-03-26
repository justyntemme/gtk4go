# Implementation of the Unified Callback System in GTK4Go

## Overview

GTK4Go utilizes a unified callback system (UCS) to manage signals and events emitted by GTK widgets. This system is designed to provide a type-safe, thread-safe, and memory-safe way to connect Go functions to GTK signals. The core logic for the UCS is located in `gtk4/callbacks.go`.

## Implementation Details

The UCS implementation revolves around a `CallbackManager` which contains the following:

- `callbacks`: A `sync.Map` that stores callback data, mapping a unique callback ID to a `callbackData` struct. The `callbackData` contains the original Go callback function, the object's pointer, the signal type, and handler IDs for disconnection. This allows looking up the original Go callback function when a signal is emitted from the C side.

- `objectHandlers`: A `sync.Map` that stores a list of handler IDs associated with each object pointer. This allows efficient disconnection of all signal handlers for a specific widget when the widget is destroyed or no longer needed.

- `objectCallbacks`: A `sync.Map` that stores callbacks associated with an object and signal type. This facilitates direct lookup of callbacks based on an object pointer and a signal type. It's used by some parts of the code to directly access a callback, for example, in the action implementation.

The key steps in using the UCS are:

1.  **Connecting a Signal:**

    - The `Connect` function generates a unique ID for each callback using `nextCallbackID`.
    - It retrieves the pointer to the GTK object using `getObjectPointer()`. This function handles various GTK object types, including widgets and adjustments.
    - It creates a `callbackData` struct containing the callback, object pointer, and signal type.
    - It stores the `callbackData` in the `callbacks` map, using the unique ID as the key.
    - It uses `C.connectSignal()` to connect the GTK signal to a generic C callback handler (`callbackHandler`, `callbackHandlerWithParam`, or `callbackHandlerWithReturn`), passing the unique ID as user data. The C callback functions are simple bridges to the Go side.
    - The ID is used in the C callback handler to look up the original Go callback.
    - The callback is associated with the object using the `trackObjectHandler` function.

2.  **Signal Emission:**

    - When a GTK signal is emitted, the appropriate C callback handler is invoked.
    - The C callback handler receives the unique callback ID as user data.
    - The C callback handler looks up the `callbackData` in the `callbacks` map using the ID.
    - The C callback handler retrieves the Go callback function from the `callbackData`.
    - The C callback handler calls the Go callback function, wrapping in `execCallback()` function which safely execute a callback on the main UI thread using the `uithread` package's `RunOnUIThread` function.

3.  **Disconnecting a Signal:**

    - The `Disconnect` function retrieves the `callbackData` using the callback ID.
    - It uses `C.disconnectSignal()` to disconnect the GTK signal handler using the `handlerID` which is stored in the callback data.
    - It removes the `callbackData` from the `callbacks` map.
    - It untracks the object handler with `untrackObjectHandler` function.

4.  **Direct Callback Storing:**
    - The `StoreDirectCallback` is used to bypass connecting a signal through the normal `Connect` function, but stores the callback directly for a given pointer. This is necessary for scenarios, such as actions, when the callback needs to be associated with a specific action pointer and the callback is not directly connected.

## Why Application Doesn't Use the Unified Callback System Directly

The `Application` struct in `gtk4/application.go` does not directly use the unified callback system for the `activate` signal. Instead, it uses a direct C callback (`activateCallback`).

```go
// Activate callback
static void activateCallback(GtkApplication* app, gpointer user_data) {
    ActivateData* data = (ActivateData*)user_data;
    GtkWidget* window = (GtkWidget*)data->window;

    // Set application
    gtk_window_set_application(GTK_WINDOW(window), app);

    // Show the window
    gtk_widget_set_visible(window, TRUE);
}

// Connect activate signal
static void connect_activate(GtkApplication* app, GtkWidget* window) {
    ActivateData* data = malloc(sizeof(ActivateData));
    data->window = window;
    data->app = app;

    g_signal_connect(app, "activate", G_CALLBACK(activateCallback), data);
}
```
