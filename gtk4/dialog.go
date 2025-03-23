// Package gtk4 provides dialog functionality for GTK4
// File: gtk4go/gtk4/dialog.go
package gtk4

/*
#cgo pkg-config: gtk4
#include <gtk/gtk.h>
#include <stdlib.h>

// Response callback for close requests
extern gboolean go_window_close_request_callback(GtkWindow *window, gpointer user_data);

// Connect window close-request signal
static gulong connect_window_close_request(GtkWindow *window, gpointer user_data) {
    return g_signal_connect(G_OBJECT(window), "close-request", G_CALLBACK(go_window_close_request_callback), user_data);
}

// Button click callback
extern void go_dialog_button_clicked(GtkButton *button, gpointer user_data);

// Connect button clicked signal for dialog buttons
static void connect_dialog_button_clicked(GtkButton *button, int response_id, gpointer dialog_ptr) {
    // Store response ID as user data (we'll cast this back to int in Go)
    g_signal_connect(G_OBJECT(button), "clicked", G_CALLBACK(go_dialog_button_clicked),
                     GINT_TO_POINTER(response_id));

    // Set dialog pointer as object data on the button
    g_object_set_data(G_OBJECT(button), "dialog-ptr", dialog_ptr);
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"sync"
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

var (
	dialogCallbacks     = make(map[uintptr]DialogResponseCallback)
	dialogCallbackMutex sync.Mutex
	debug               = false // Set to true to enable debug logging
)

// Debug logging helper
func debugLog(format string, args ...interface{}) {
	if debug {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// Go callback for window close requests
//
//export go_window_close_request_callback
func go_window_close_request_callback(window *C.GtkWindow, userData C.gpointer) C.gboolean {
	windowPtr := uintptr(unsafe.Pointer(window))
	debugLog("Window close request for %v", windowPtr)

	dialogCallbackMutex.Lock()
	callback, exists := dialogCallbacks[windowPtr]
	dialogCallbackMutex.Unlock()

	if exists {
		debugLog("Found callback for window %v, sending DeleteEvent", windowPtr)
		// Execute callback outside the mutex lock to prevent deadlocks
		callback(ResponseDeleteEvent)
	} else {
		debugLog("No callback found for window %v", windowPtr)
	}

	// Return FALSE to allow the window to close
	return C.FALSE
}

// Go callback for dialog button clicks
//
//export go_dialog_button_clicked
func go_dialog_button_clicked(button *C.GtkButton, userData C.gpointer) {
	// Extract response ID from user data (which is an integer cast to pointer)
	responseId := ResponseType(uintptr(userData))

	// Get the dialog pointer from the button's object data
	dialogPtr := uintptr(unsafe.Pointer(C.g_object_get_data(
		(*C.GObject)(unsafe.Pointer(button)),
		(*C.gchar)(unsafe.Pointer(C.CString("dialog-ptr"))),
	)))

	debugLog("Button clicked with response %d for dialog %v", responseId, dialogPtr)

	dialogCallbackMutex.Lock()
	callback, exists := dialogCallbacks[dialogPtr]
	dialogCallbackMutex.Unlock()

	if exists {
		debugLog("Found callback for dialog %v, sending response %d", dialogPtr, responseId)
		// Execute callback outside the mutex lock to prevent deadlocks
		callback(responseId)
	} else {
		debugLog("No callback found for dialog %v", dialogPtr)
	}
}

// Dialog represents a GTK4 dialog window
type Dialog struct {
	Window
	buttonArea  *Box
	contentArea *Box
}

// NewDialog creates a new dialog window
func NewDialog(title string, parent *Window, flags DialogFlags) *Dialog {
	// Create a new window
	var parentPtr *C.GtkWindow
	if parent != nil {
		parentPtr = (*C.GtkWindow)(unsafe.Pointer(parent.widget))
	}

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	// Create a window with the appropriate properties
	windowWidget := C.gtk_window_new()
	C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(windowWidget)), cTitle)

	// Set the window to be modal if requested
	if flags&DialogModal != 0 {
		C.gtk_window_set_modal((*C.GtkWindow)(unsafe.Pointer(windowWidget)), C.TRUE)
	}

	// Set the window to be destroyed with parent if requested
	if flags&DialogDestroyWithParent != 0 {
		C.gtk_window_set_destroy_with_parent((*C.GtkWindow)(unsafe.Pointer(windowWidget)), C.TRUE)
	}

	// Set the transient parent
	if parent != nil {
		C.gtk_window_set_transient_for((*C.GtkWindow)(unsafe.Pointer(windowWidget)), parentPtr)
	}

	// Create a dialog object
	dialog := &Dialog{
		Window: Window{
			widget: windowWidget,
		},
	}

	// Connect the close request callback
	windowPtr := uintptr(unsafe.Pointer(windowWidget))
	C.connect_window_close_request((*C.GtkWindow)(unsafe.Pointer(windowWidget)), C.gpointer(windowPtr))

	// Create a vertical box container for the dialog
	mainBox := NewBox(OrientationVertical, 0)

	// Create a content area
	dialog.contentArea = NewBox(OrientationVertical, 10)
	dialog.contentArea.AddCssClass("dialog-content-area")

	// Add padding to the content area
	C.gtk_widget_set_margin_start(dialog.contentArea.widget, 16)
	C.gtk_widget_set_margin_end(dialog.contentArea.widget, 16)
	C.gtk_widget_set_margin_top(dialog.contentArea.widget, 16)
	C.gtk_widget_set_margin_bottom(dialog.contentArea.widget, 16)

	// Create a button area (horizontal box at the bottom)
	dialog.buttonArea = NewBox(OrientationHorizontal, 6)
	dialog.buttonArea.AddCssClass("dialog-button-area")

	// Add CSS for button area to have buttons aligned to the right
	dialog.buttonArea.SetHomogeneous(false)
	C.gtk_widget_set_halign(dialog.buttonArea.widget, C.GTK_ALIGN_END)
	C.gtk_widget_set_margin_start(dialog.buttonArea.widget, 16)
	C.gtk_widget_set_margin_end(dialog.buttonArea.widget, 16)
	C.gtk_widget_set_margin_top(dialog.buttonArea.widget, 16)
	C.gtk_widget_set_margin_bottom(dialog.buttonArea.widget, 16)

	// Add the content and button areas to the main box
	mainBox.Append(dialog.contentArea)
	mainBox.Append(dialog.buttonArea)

	// Set the main box as the child of the window
	dialog.SetChild(mainBox)

	// Set a reasonable default size
	dialog.SetDefaultSize(350, 150)

	// Set up finalizer for cleanup
	runtime.SetFinalizer(dialog, (*Dialog).Destroy)

	debugLog("Created new dialog %v", dialog.Native())

	return dialog
}

// AddButton adds a button to the dialog
func (d *Dialog) AddButton(text string, responseId ResponseType) *Button {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	// Create a button with the given text
	button := NewButton(text)

	// Add the button to the button area
	d.buttonArea.Append(button)

	// Connect the button to trigger the response
	C.connect_dialog_button_clicked(
		(*C.GtkButton)(unsafe.Pointer(button.widget)),
		C.int(responseId),
		C.gpointer(unsafe.Pointer(d.widget)),
	)

	debugLog("Added button with response %d to dialog %v", responseId, d.Native())

	return button
}

// GetContentArea gets the content area of the dialog
func (d *Dialog) GetContentArea() *Box {
	return d.contentArea
}

// GetButtonArea gets the button area of the dialog
func (d *Dialog) GetButtonArea() *Box {
	return d.buttonArea
}

// ConnectResponse connects a response callback to the dialog
func (d *Dialog) ConnectResponse(callback DialogResponseCallback) {
	dialogPtr := d.Native()

	dialogCallbackMutex.Lock()
	defer dialogCallbackMutex.Unlock()

	// Store the callback
	dialogCallbacks[dialogPtr] = callback

	debugLog("Connected response callback to dialog %v", dialogPtr)
}

// Destroy overrides Window's Destroy to clean up dialog resources
func (d *Dialog) Destroy() {
	dialogPtr := d.Native()
	debugLog("Destroying dialog %v", dialogPtr)

	dialogCallbackMutex.Lock()
	// Remove callback
	delete(dialogCallbacks, dialogPtr)
	dialogCallbackMutex.Unlock()

	// Call the parent's Destroy method
	d.Window.Destroy()
}

// MessageDialog represents a GTK message dialog
type MessageDialog struct {
	Dialog
	messageType MessageType
}

// MessageType defines the type of message
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

// ButtonsType defines standard button combinations
type ButtonsType int

const (
	// ButtonsNone no buttons
	ButtonsNone ButtonsType = iota
	// ButtonsOk OK button
	ButtonsOk
	// ButtonsClose Close button
	ButtonsClose
	// ButtonsCancel Cancel button
	ButtonsCancel
	// ButtonsYesNo Yes and No buttons
	ButtonsYesNo
	// ButtonsOkCancel OK and Cancel buttons
	ButtonsOkCancel
)

// NewMessageDialog creates a new message dialog
func NewMessageDialog(parent *Window, flags DialogFlags, messageType MessageType, buttons ResponseType, message string) *MessageDialog {
	// Create a dialog
	dialog := &MessageDialog{
		Dialog:      *NewDialog("", parent, flags),
		messageType: messageType,
	}

	// Set the appropriate CSS class for the message type
	switch messageType {
	case MessageInfo:
		dialog.AddCssClass("info-dialog")
	case MessageWarning:
		dialog.AddCssClass("warning-dialog")
	case MessageQuestion:
		dialog.AddCssClass("question-dialog")
	case MessageError:
		dialog.AddCssClass("error-dialog")
	}

	// Add the message to the content area
	messageLabel := NewLabel(message)
	messageLabel.AddCssClass("dialog-message")
	dialog.GetContentArea().Append(messageLabel)

	// Add buttons based on the button type
	if buttons&ResponseOk != 0 {
		dialog.AddButton("OK", ResponseOk)
	}
	if buttons&ResponseClose != 0 {
		dialog.AddButton("Close", ResponseClose)
	}
	if buttons&ResponseCancel != 0 {
		dialog.AddButton("Cancel", ResponseCancel)
	}
	if buttons&ResponseYes != 0 {
		dialog.AddButton("Yes", ResponseYes)
	}
	if buttons&ResponseNo != 0 {
		dialog.AddButton("No", ResponseNo)
	}

	debugLog("Created new message dialog %v with type %d", dialog.Native(), messageType)

	return dialog
}

// SetMarkup sets the message using markup
func (d *MessageDialog) SetMarkup(markup string) {
	// Create a new label with markup
	messageLabel := NewLabel("")
	messageLabel.SetMarkup(markup)
	messageLabel.AddCssClass("dialog-message")

	// Get the content area
	contentArea := d.GetContentArea()

	// Since we don't have direct access to the children, we'll just add the new label
	contentArea.Append(messageLabel)
}

// FileDialog represents a file selection dialog
type FileDialog struct {
	Dialog
	fileEntry  *Entry
	actionType FileDialogAction
}

// FileDialogAction defines the type of file dialog
type FileDialogAction int

const (
	// FileDialogOpen for opening files
	FileDialogOpen FileDialogAction = iota
	// FileDialogSave for saving files
	FileDialogSave
	// FileDialogSelectFolder for selecting folders
	FileDialogSelectFolder
)

// NewFileDialog creates a new file dialog
func NewFileDialog(title string, parent *Window, action FileDialogAction) *FileDialog {
	// Create a dialog
	dialog := &FileDialog{
		Dialog:     *NewDialog(title, parent, DialogModal),
		actionType: action,
	}

	// Add content for file selection
	contentArea := dialog.GetContentArea()

	// Create a label
	var labelText string
	switch action {
	case FileDialogOpen:
		labelText = "Select a file to open:"
	case FileDialogSave:
		labelText = "Save file as:"
	case FileDialogSelectFolder:
		labelText = "Select folder:"
	}

	fileLabel := NewLabel(labelText)
	contentArea.Append(fileLabel)

	// Create an entry for the file path
	dialog.fileEntry = NewEntry()

	// Set placeholder text based on action
	switch action {
	case FileDialogOpen:
		dialog.fileEntry.SetPlaceholderText("File path")
	case FileDialogSave:
		dialog.fileEntry.SetPlaceholderText("Enter filename")
	case FileDialogSelectFolder:
		dialog.fileEntry.SetPlaceholderText("Folder path")
	}

	contentArea.Append(dialog.fileEntry)

	// Add appropriate buttons
	switch action {
	case FileDialogOpen:
		dialog.AddButton("Cancel", ResponseCancel)
		dialog.AddButton("Open", ResponseAccept)
	case FileDialogSave:
		dialog.AddButton("Cancel", ResponseCancel)
		dialog.AddButton("Save", ResponseAccept)
	case FileDialogSelectFolder:
		dialog.AddButton("Cancel", ResponseCancel)
		dialog.AddButton("Select", ResponseAccept)
	}

	debugLog("Created new file dialog %v with action %d", dialog.Native(), action)

	return dialog
}

// GetFilename gets the filename from the entry
func (d *FileDialog) GetFilename() string {
	return d.fileEntry.GetText()
}

// SetFilename sets the filename in the entry
func (d *FileDialog) SetFilename(filename string) {
	d.fileEntry.SetText(filename)
}

// Convenience Functions

// ShowMessageDialog shows a message dialog and returns the response
func ShowMessageDialog(parent *Window, messageType MessageType, title, message string) ResponseType {
	dialog := NewMessageDialog(parent, DialogModal, messageType, ResponseOk, message)
	dialog.SetTitle(title)
	dialog.Show()

	var response ResponseType
	done := make(chan bool)

	dialog.ConnectResponse(func(responseId ResponseType) {
		debugLog("Message dialog response: %d", responseId)
		response = responseId
		dialog.Destroy()
		done <- true
	})

	<-done
	return response
}

// ShowConfirmDialog shows a confirmation dialog and returns true if confirmed
func ShowConfirmDialog(parent *Window, title, message string) bool {
	dialog := NewMessageDialog(parent, DialogModal, MessageQuestion, ResponseYes|ResponseNo, message)
	dialog.SetTitle(title)
	dialog.Show()

	var confirmed bool
	done := make(chan bool)

	dialog.ConnectResponse(func(responseId ResponseType) {
		debugLog("Confirm dialog response: %d", responseId)
		confirmed = (responseId == ResponseYes)
		dialog.Destroy()
		done <- true
	})

	<-done
	return confirmed
}

// ShowFileOpenDialog shows a file open dialog and returns the selected filename
func ShowFileOpenDialog(parent *Window, title string) (string, bool) {
	dialog := NewFileDialog(title, parent, FileDialogOpen)
	dialog.Show()

	var filename string
	var selected bool
	done := make(chan bool)

	dialog.ConnectResponse(func(responseId ResponseType) {
		debugLog("File open dialog response: %d", responseId)
		if responseId == ResponseAccept {
			filename = dialog.GetFilename()
			selected = true
		} else {
			selected = false
		}
		dialog.Destroy()
		done <- true
	})

	<-done
	return filename, selected
}

// ShowFileSaveDialog shows a file save dialog and returns the selected filename
func ShowFileSaveDialog(parent *Window, title string) (string, bool) {
	dialog := NewFileDialog(title, parent, FileDialogSave)
	dialog.Show()

	var filename string
	var selected bool
	done := make(chan bool)

	dialog.ConnectResponse(func(responseId ResponseType) {
		debugLog("File save dialog response: %d", responseId)
		if responseId == ResponseAccept {
			filename = dialog.GetFilename()
			selected = true
		} else {
			selected = false
		}
		dialog.Destroy()
		done <- true
	})

	<-done
	return filename, selected
}

// ShowFolderSelectDialog shows a folder selection dialog and returns the selected folder
func ShowFolderSelectDialog(parent *Window, title string) (string, bool) {
	dialog := NewFileDialog(title, parent, FileDialogSelectFolder)
	dialog.Show()

	var folder string
	var selected bool
	done := make(chan bool)

	dialog.ConnectResponse(func(responseId ResponseType) {
		debugLog("Folder select dialog response: %d", responseId)
		if responseId == ResponseAccept {
			folder = dialog.GetFilename()
			selected = true
		} else {
			selected = false
		}
		dialog.Destroy()
		done <- true
	})

	<-done
	return folder, selected
}
