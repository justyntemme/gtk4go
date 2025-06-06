// Package gtk4 provides dialog functionality for GTK4
// File: gtk4go/gtk4/dialog.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// extern void buttonResponseCallback(GtkButton *button, gpointer user_data);
// extern gboolean windowCloseCallback(GtkWindow *window, gpointer user_data);
//
// static void connectButtonResponse(GtkButton *button, int response_id, gpointer dialog_ptr) {
//     // Store response ID as object data on the button
//     g_object_set_data(G_OBJECT(button), "response-id", GINT_TO_POINTER(response_id));
//
//     // Store dialog pointer as object data on the button
//     g_object_set_data(G_OBJECT(button), "dialog-ptr", dialog_ptr);
//
//     // Connect the clicked signal
//     g_signal_connect(button, "clicked", G_CALLBACK(buttonResponseCallback), button);
// }
//
// static void connectWindowClose(GtkWindow *window) {
//     g_signal_connect(window, "close-request", G_CALLBACK(windowCloseCallback), window);
// }
import "C"

import (
	"unsafe"
)

// ResponseType defines standard response IDs
type ResponseType int

const (
	// ResponseNone no response
	ResponseNone ResponseType = C.GTK_RESPONSE_NONE
	// ResponseReject reject the dialog
	ResponseReject ResponseType = C.GTK_RESPONSE_REJECT
	// ResponseAccept accept the dialog
	ResponseAccept ResponseType = C.GTK_RESPONSE_ACCEPT
	// ResponseDeleteEvent dialog was deleted
	ResponseDeleteEvent ResponseType = C.GTK_RESPONSE_DELETE_EVENT
	// ResponseOk affirmative response
	ResponseOk ResponseType = C.GTK_RESPONSE_OK
	// ResponseCancel negative response
	ResponseCancel ResponseType = C.GTK_RESPONSE_CANCEL
	// ResponseClose close response
	ResponseClose ResponseType = C.GTK_RESPONSE_CLOSE
	// ResponseYes yes response
	ResponseYes ResponseType = C.GTK_RESPONSE_YES
	// ResponseNo no response
	ResponseNo ResponseType = C.GTK_RESPONSE_NO
	// ResponseApply apply response
	ResponseApply ResponseType = C.GTK_RESPONSE_APPLY
	// ResponseHelp help response
	ResponseHelp ResponseType = C.GTK_RESPONSE_HELP
)

// DialogFlags defines behavior flags for dialogs
type DialogFlags int

const (
	// DialogModal makes the dialog modal
	DialogModal DialogFlags = 1 << 0
	// DialogDestroyWithParent destroys the dialog when its parent is destroyed
	DialogDestroyWithParent DialogFlags = 1 << 1
	// DialogUseHeaderBar uses a header bar for the dialog
	DialogUseHeaderBar DialogFlags = 1 << 2
)

// DialogResponseCallback represents a callback for dialog response events
type DialogResponseCallback func(responseId ResponseType)

// Custom signal type for dialog responses
const (
	SignalDialogResponse SignalType = "dialog-response"
)

//export buttonResponseCallback
func buttonResponseCallback(button *C.GtkButton, userData C.gpointer) {
	// Get response ID from button data
	responsePtr := C.g_object_get_data((*C.GObject)(unsafe.Pointer(button)), C.CString("response-id"))
	responseId := ResponseType(uintptr(responsePtr))

	// Get dialog pointer from button data
	dialogPtr := uintptr(C.g_object_get_data((*C.GObject)(unsafe.Pointer(button)), C.CString("dialog-ptr")))

	DebugLog(DebugLevelVerbose, DebugComponentDialog, "Button clicked with response %d for dialog %v", responseId, dialogPtr)

	// Lookup callback using the unified callback system
	if callback := GetCallback(dialogPtr, SignalDialogResponse); callback != nil {
		// Use the SafeCallback function to execute the callback
		if typedCallback, ok := callback.(func(ResponseType)); ok {
			SafeCallback(typedCallback, responseId)
		} else {
			DebugLog(DebugLevelError, DebugComponentDialog, 
				"Invalid callback type for dialog response: %T, expected func(ResponseType)", callback)
		}
	}
}

//export windowCloseCallback
func windowCloseCallback(window *C.GtkWindow, userData C.gpointer) C.gboolean {
	windowPtr := uintptr(unsafe.Pointer(window))
	DebugLog(DebugLevelVerbose, DebugComponentDialog, "Window close request for %v", windowPtr)

	// Lookup callback using the unified callback system
	if callback := GetCallback(windowPtr, SignalDialogResponse); callback != nil {
		// Use the SafeCallback function to execute the callback
		if typedCallback, ok := callback.(func(ResponseType)); ok {
			SafeCallback(typedCallback, ResponseDeleteEvent)
		} else {
			DebugLog(DebugLevelError, DebugComponentDialog, 
				"Invalid callback type for dialog response: %T, expected func(ResponseType)", callback)
		}
	}

	// Return FALSE to allow the window to close
	return C.FALSE
}

// Dialog represents a GTK dialog
type Dialog struct {
	Window
	buttonArea  *Box
	contentArea *Box
}

// NewDialog creates a new dialog
func NewDialog(title string, parent *Window, flags DialogFlags) *Dialog {
	// Create base window
	window := NewWindow(title)

	// Set modal and transient parent
	if flags&DialogModal != 0 {
		C.gtk_window_set_modal((*C.GtkWindow)(unsafe.Pointer(window.widget)), C.TRUE)
	}

	if parent != nil {
		C.gtk_window_set_transient_for(
			(*C.GtkWindow)(unsafe.Pointer(window.widget)),
			(*C.GtkWindow)(unsafe.Pointer(parent.widget)),
		)
	}

	if flags&DialogDestroyWithParent != 0 {
		C.gtk_window_set_destroy_with_parent(
			(*C.GtkWindow)(unsafe.Pointer(window.widget)),
			C.TRUE,
		)
	}

	// Create a dialog
	dialog := &Dialog{
		Window: *window,
	}

	// Connect window close handler
	C.connectWindowClose((*C.GtkWindow)(unsafe.Pointer(window.widget)))

	// Create a box for content
	mainBox := NewBox(OrientationVertical, 0)

	// Create content area
	dialog.contentArea = NewBox(OrientationVertical, 10)
	dialog.contentArea.AddCssClass("dialog-content-area")

	// Create button area
	dialog.buttonArea = NewBox(OrientationHorizontal, 6)
	dialog.buttonArea.AddCssClass("dialog-button-area")

	// Set up the button area for dialog buttons
	dialog.buttonArea.SetHomogeneous(false)
	C.gtk_widget_set_halign(dialog.buttonArea.widget, C.GTK_ALIGN_END)

	// Add padding
	C.gtk_widget_set_margin_start(dialog.contentArea.widget, 16)
	C.gtk_widget_set_margin_end(dialog.contentArea.widget, 16)
	C.gtk_widget_set_margin_top(dialog.contentArea.widget, 16)
	C.gtk_widget_set_margin_bottom(dialog.contentArea.widget, 16)

	C.gtk_widget_set_margin_start(dialog.buttonArea.widget, 16)
	C.gtk_widget_set_margin_end(dialog.buttonArea.widget, 16)
	C.gtk_widget_set_margin_top(dialog.buttonArea.widget, 10)
	C.gtk_widget_set_margin_bottom(dialog.buttonArea.widget, 16)

	// Add the areas to the main box
	mainBox.Append(dialog.contentArea)
	mainBox.Append(dialog.buttonArea)

	// Add the main box to the window
	dialog.SetChild(mainBox)

	// Set up default size
	dialog.SetDefaultSize(400, 200)

	return dialog
}

// AddButton adds a button to the dialog
func (d *Dialog) AddButton(text string, responseId ResponseType) *Button {
	// Create a button
	button := NewButton(text)

	// Add it to the button area
	d.buttonArea.Append(button)

	// Connect button to response using C helper
	C.connectButtonResponse(
		(*C.GtkButton)(unsafe.Pointer(button.widget)),
		C.int(responseId),
		C.gpointer(unsafe.Pointer(d.widget)),
	)

	return button
}

// GetContentArea gets the content area of the dialog
func (d *Dialog) GetContentArea() *Box {
	return d.contentArea
}

// ConnectResponse connects a response callback to the dialog
func (d *Dialog) ConnectResponse(callback DialogResponseCallback) {
	// Convert DialogResponseCallback to a standard callback type
	// that the unified callback system can handle
	standardCallback := func(responseId ResponseType) {
		callback(responseId)
	}

	// Store in the unified callback system
	dialogPtr := uintptr(unsafe.Pointer(d.widget))
	
	// Use the store direct callback method to ensure it's stored correctly
	StoreDirectCallback(dialogPtr, SignalDialogResponse, standardCallback)

	DebugLog(DebugLevelInfo, DebugComponentDialog, 
		"Connected response callback to dialog %v", dialogPtr)
}

// DisconnectResponse disconnects the response callback from the dialog
func (d *Dialog) DisconnectResponse() {
	// Get the dialog pointer
	dialogPtr := uintptr(unsafe.Pointer(d.widget))
	
	// Remove callbacks for the dialog response signal
	globalCallbackManager.objectCallbacks.Range(func(key, value interface{}) bool {
		if key.(uintptr) == dialogPtr {
			callbackMap := value.(map[SignalType]interface{})
			delete(callbackMap, SignalDialogResponse)
			
			// If no more callbacks, remove the whole entry
			if len(callbackMap) == 0 {
				globalCallbackManager.objectCallbacks.Delete(key)
			} else {
				// Otherwise update with the modified map
				globalCallbackManager.objectCallbacks.Store(key, callbackMap)
			}
		}
		return true
	})
}

// Destroy overrides Window's Destroy to clean up dialog resources
func (d *Dialog) Destroy() {
	DebugLog(DebugLevelInfo, DebugComponentDialog, "Destroying dialog %v", 
		uintptr(unsafe.Pointer(d.widget)))

	// Remove any callbacks associated with this dialog
	d.DisconnectResponse()
	
	// Call the base Destroy to handle window cleanup
	d.Window.Destroy()
}

// MessageType defines the type of message dialog
type MessageType int

const (
	// MessageInfo for informational messages
	MessageInfo MessageType = iota
	// MessageWarning for warning messages
	MessageWarning
	// MessageQuestion for questions requiring user response
	MessageQuestion
	// MessageError for error messages
	MessageError
	// MessageOther for other messages
	MessageOther
)

// MessageDialog represents a GTK message dialog
type MessageDialog struct {
	Dialog
	messageType MessageType
}

// NewMessageDialog creates a new message dialog
func NewMessageDialog(parent *Window, flags DialogFlags, messageType MessageType, buttons ResponseType, message string) *MessageDialog {
	// Create a dialog
	dialog := NewDialog("", parent, flags)

	// Create a message dialog
	msgDialog := &MessageDialog{
		Dialog:      *dialog,
		messageType: messageType,
	}

	// Set CSS class based on message type
	switch messageType {
	case MessageInfo:
		msgDialog.AddCssClass("info-dialog")
	case MessageWarning:
		msgDialog.AddCssClass("warning-dialog")
	case MessageQuestion:
		msgDialog.AddCssClass("question-dialog")
	case MessageError:
		msgDialog.AddCssClass("error-dialog")
	}

	// Add message label
	msgLabel := NewLabel(message)
	msgLabel.AddCssClass("dialog-message")
	msgDialog.GetContentArea().Append(msgLabel)

	// Add buttons
	if buttons&ResponseOk != 0 {
		msgDialog.AddButton("OK", ResponseOk)
	}
	if buttons&ResponseClose != 0 {
		msgDialog.AddButton("Close", ResponseClose)
	}
	if buttons&ResponseCancel != 0 {
		msgDialog.AddButton("Cancel", ResponseCancel)
	}
	if buttons&ResponseYes != 0 {
		msgDialog.AddButton("Yes", ResponseYes)
	}
	if buttons&ResponseNo != 0 {
		msgDialog.AddButton("No", ResponseNo)
	}

	return msgDialog
}

// FileDialogAction defines the type of file chooser
type FileDialogAction int

const (
	// FileDialogActionOpen for selecting an existing file
	FileDialogActionOpen FileDialogAction = iota
	// FileDialogActionSave for saving a file
	FileDialogActionSave
	// FileDialogActionSelectFolder for selecting a folder
	FileDialogActionSelectFolder
)

// FileDialog represents a GTK file chooser dialog
type FileDialog struct {
	Dialog
	fileEntry  *Entry
	actionType FileDialogAction
}

// NewFileDialog creates a new file chooser dialog
func NewFileDialog(title string, parent *Window, action FileDialogAction) *FileDialog {
	// Create a dialog
	dialog := NewDialog(title, parent, DialogModal)

	// Create a file dialog
	fileDialog := &FileDialog{
		Dialog:     *dialog,
		actionType: action,
	}

	// Add content
	contentArea := fileDialog.GetContentArea()

	// Add label based on action
	var labelText string
	switch action {
	case FileDialogActionOpen:
		labelText = "Select a file to open:"
	case FileDialogActionSave:
		labelText = "Save file as:"
	case FileDialogActionSelectFolder:
		labelText = "Select folder:"
	}

	fileLabel := NewLabel(labelText)
	contentArea.Append(fileLabel)

	// Add entry for file path
	fileDialog.fileEntry = NewEntry()

	// Set placeholder text
	switch action {
	case FileDialogActionOpen:
		fileDialog.fileEntry.SetPlaceholderText("File path")
	case FileDialogActionSave:
		fileDialog.fileEntry.SetPlaceholderText("Enter filename")
	case FileDialogActionSelectFolder:
		fileDialog.fileEntry.SetPlaceholderText("Folder path")
	}

	contentArea.Append(fileDialog.fileEntry)

	// Add appropriate buttons
	switch action {
	case FileDialogActionOpen:
		fileDialog.AddButton("Cancel", ResponseCancel)
		fileDialog.AddButton("Open", ResponseAccept)
	case FileDialogActionSave:
		fileDialog.AddButton("Cancel", ResponseCancel)
		fileDialog.AddButton("Save", ResponseAccept)
	case FileDialogActionSelectFolder:
		fileDialog.AddButton("Cancel", ResponseCancel)
		fileDialog.AddButton("Select", ResponseAccept)
	}

	return fileDialog
}

// GetFilename gets the filename from the file dialog
func (d *FileDialog) GetFilename() string {
	return d.fileEntry.GetText()
}

// SetFilename sets the filename in the file dialog
func (d *FileDialog) SetFilename(filename string) {
	d.fileEntry.SetText(filename)
}