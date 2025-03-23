// Package gtk4 provides entry widget functionality for GTK4
// File: gtk4go/gtk4/entry.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback functions for entry signals
// extern void entryChangedCallback(GtkEditable *editable, gpointer user_data);
// extern void entryActivateCallback(GtkEntry *entry, gpointer user_data);
//
// // Connect entry signals with callbacks
// static gulong connectEntryChanged(GtkWidget *entry, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(entry), "changed", G_CALLBACK(entryChangedCallback), user_data);
// }
//
// static gulong connectEntryActivate(GtkWidget *entry, gpointer user_data) {
//     return g_signal_connect(G_OBJECT(entry), "activate", G_CALLBACK(entryActivateCallback), user_data);
// }
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// EntryCallback represents a callback for entry events
type EntryCallback func()

var (
	entryChangedCallbacks  = make(map[uintptr]EntryCallback)
	entryActivateCallbacks = make(map[uintptr]EntryCallback)
	entryCallbackMutex     sync.Mutex
)

//export entryChangedCallback
func entryChangedCallback(editable *C.GtkEditable, userData C.gpointer) {
	entryCallbackMutex.Lock()
	defer entryCallbackMutex.Unlock()

	// Convert entry pointer to uintptr for lookup
	entryPtr := uintptr(unsafe.Pointer(editable))

	// Find and call the callback
	if callback, ok := entryChangedCallbacks[entryPtr]; ok {
		callback()
	}
}

//export entryActivateCallback
func entryActivateCallback(entry *C.GtkEntry, userData C.gpointer) {
	entryCallbackMutex.Lock()
	defer entryCallbackMutex.Unlock()

	// Convert entry pointer to uintptr for lookup
	entryPtr := uintptr(unsafe.Pointer(entry))

	// Find and call the callback
	if callback, ok := entryActivateCallbacks[entryPtr]; ok {
		callback()
	}
}

// Entry represents a GTK entry widget for text input
type Entry struct {
	widget *C.GtkWidget
}

// NewEntry creates a new GTK entry widget
func NewEntry() *Entry {
	entry := &Entry{
		widget: C.gtk_entry_new(),
	}
	runtime.SetFinalizer(entry, (*Entry).Destroy)
	return entry
}

// NewEntryWithBuffer creates a new GTK entry widget with a specific buffer
func NewEntryWithBuffer(buffer *EntryBuffer) *Entry {
	entry := &Entry{
		widget: C.gtk_entry_new_with_buffer((*C.GtkEntryBuffer)(unsafe.Pointer(buffer.buffer))),
	}
	runtime.SetFinalizer(entry, (*Entry).Destroy)
	return entry
}

// SetText sets the text in the entry
func (e *Entry) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gtk_editable_set_text((*C.GtkEditable)(unsafe.Pointer(e.widget)), cText)
}

// GetText gets the text from the entry
func (e *Entry) GetText() string {
	cText := C.gtk_editable_get_text((*C.GtkEditable)(unsafe.Pointer(e.widget)))
	if cText == nil {
		return ""
	}
	return C.GoString(cText)
}

// SetPlaceholderText sets the placeholder text shown when the entry is empty
func (e *Entry) SetPlaceholderText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gtk_entry_set_placeholder_text((*C.GtkEntry)(unsafe.Pointer(e.widget)), cText)
}

// GetPlaceholderText gets the placeholder text
func (e *Entry) GetPlaceholderText() string {
	cText := C.gtk_entry_get_placeholder_text((*C.GtkEntry)(unsafe.Pointer(e.widget)))
	if cText == nil {
		return ""
	}
	return C.GoString(cText)
}

// SetEditable sets whether the user can edit the text
func (e *Entry) SetEditable(editable bool) {
	if editable {
		C.gtk_editable_set_editable((*C.GtkEditable)(unsafe.Pointer(e.widget)), C.TRUE)
	} else {
		C.gtk_editable_set_editable((*C.GtkEditable)(unsafe.Pointer(e.widget)), C.FALSE)
	}
}

// SetVisibility sets whether the text is visible or hidden (e.g., for passwords)
func (e *Entry) SetVisibility(visible bool) {
	if visible {
		C.gtk_entry_set_visibility((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.TRUE)
	} else {
		C.gtk_entry_set_visibility((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.FALSE)
	}
}

// SetMaxLength sets the maximum allowed length of the text
func (e *Entry) SetMaxLength(max int) {
	C.gtk_entry_set_max_length((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.int(max))
}

// GetBuffer gets the entry buffer
func (e *Entry) GetBuffer() *EntryBuffer {
	buffer := C.gtk_entry_get_buffer((*C.GtkEntry)(unsafe.Pointer(e.widget)))
	return &EntryBuffer{
		buffer: buffer,
	}
}

// SetBuffer sets the entry buffer
func (e *Entry) SetBuffer(buffer *EntryBuffer) {
	C.gtk_entry_set_buffer((*C.GtkEntry)(unsafe.Pointer(e.widget)), buffer.buffer)
}

// ConnectChanged connects a callback function to the entry's "changed" signal
func (e *Entry) ConnectChanged(callback EntryCallback) {
	entryCallbackMutex.Lock()
	defer entryCallbackMutex.Unlock()

	// Store callback in map
	entryPtr := uintptr(unsafe.Pointer(e.widget))
	entryChangedCallbacks[entryPtr] = callback

	// Connect signal
	C.connectEntryChanged(e.widget, C.gpointer(unsafe.Pointer(e.widget)))
}

// ConnectActivate connects a callback function to the entry's "activate" signal (when Enter is pressed)
func (e *Entry) ConnectActivate(callback EntryCallback) {
	entryCallbackMutex.Lock()
	defer entryCallbackMutex.Unlock()

	// Store callback in map
	entryPtr := uintptr(unsafe.Pointer(e.widget))
	entryActivateCallbacks[entryPtr] = callback

	// Connect signal
	C.connectEntryActivate(e.widget, C.gpointer(unsafe.Pointer(e.widget)))
}

// SetInputPurpose sets the purpose of the entry (e.g., password, URL, etc.)
func (e *Entry) SetInputPurpose(purpose InputPurpose) {
	C.gtk_entry_set_input_purpose((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.GtkInputPurpose(purpose))
}

// SetInputHints sets input hints for the entry
func (e *Entry) SetInputHints(hints InputHints) {
	C.gtk_entry_set_input_hints((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.GtkInputHints(hints))
}

// Destroy destroys the entry
func (e *Entry) Destroy() {
	entryCallbackMutex.Lock()
	defer entryCallbackMutex.Unlock()

	// Remove callbacks from maps if they exist
	entryPtr := uintptr(unsafe.Pointer(e.widget))
	delete(entryChangedCallbacks, entryPtr)
	delete(entryActivateCallbacks, entryPtr)

	// Destroy widget
	C.gtk_widget_unparent(e.widget)
	e.widget = nil
}

// Native returns the underlying GtkWidget pointer
func (e *Entry) Native() uintptr {
	return uintptr(unsafe.Pointer(e.widget))
}

// GetWidget returns the underlying GtkWidget pointer
func (e *Entry) GetWidget() *C.GtkWidget {
	return e.widget
}

// EntryBuffer represents a GTK entry buffer
type EntryBuffer struct {
	buffer *C.GtkEntryBuffer
}

// NewEntryBuffer creates a new entry buffer with initial text
func NewEntryBuffer(initialText string) *EntryBuffer {
	cText := C.CString(initialText)
	defer C.free(unsafe.Pointer(cText))

	buffer := &EntryBuffer{
		buffer: C.gtk_entry_buffer_new(cText, C.int(len(initialText))),
	}
	runtime.SetFinalizer(buffer, (*EntryBuffer).Free)
	return buffer
}

// SetText sets the text in the buffer
func (b *EntryBuffer) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gtk_entry_buffer_set_text(b.buffer, cText, C.int(len(text)))
}

// GetText gets the text from the buffer
func (b *EntryBuffer) GetText() string {
	cText := C.gtk_entry_buffer_get_text(b.buffer)
	if cText == nil {
		return ""
	}
	return C.GoString(cText)
}

// GetLength gets the length of the text in the buffer
func (b *EntryBuffer) GetLength() int {
	return int(C.gtk_entry_buffer_get_length(b.buffer))
}

// Free frees the buffer
func (b *EntryBuffer) Free() {
	C.g_object_unref(C.gpointer(unsafe.Pointer(b.buffer)))
	b.buffer = nil
}

// InputPurpose defines the purpose of an entry
type InputPurpose int

const (
	// InputPurposeFreeForm for normal text entry
	InputPurposeFreeForm InputPurpose = C.GTK_INPUT_PURPOSE_FREE_FORM
	// InputPurposeAlpha for alphabetic entry
	InputPurposeAlpha InputPurpose = C.GTK_INPUT_PURPOSE_ALPHA
	// InputPurposeDigits for digit entry
	InputPurposeDigits InputPurpose = C.GTK_INPUT_PURPOSE_DIGITS
	// InputPurposeNumber for number entry
	InputPurposeNumber InputPurpose = C.GTK_INPUT_PURPOSE_NUMBER
	// InputPurposePhone for phone number entry
	InputPurposePhone InputPurpose = C.GTK_INPUT_PURPOSE_PHONE
	// InputPurposeURL for URL entry
	InputPurposeURL InputPurpose = C.GTK_INPUT_PURPOSE_URL
	// InputPurposeEmail for email entry
	InputPurposeEmail InputPurpose = C.GTK_INPUT_PURPOSE_EMAIL
	// InputPurposeName for name entry
	InputPurposeName InputPurpose = C.GTK_INPUT_PURPOSE_NAME
	// InputPurposePassword for password entry
	InputPurposePassword InputPurpose = C.GTK_INPUT_PURPOSE_PASSWORD
	// InputPurposePin for PIN entry
	InputPurposePin InputPurpose = C.GTK_INPUT_PURPOSE_PIN
)

// InputHints defines input hints for an entry
type InputHints int

const (
	// InputHintsNone for no hints
	InputHintsNone InputHints = C.GTK_INPUT_HINT_NONE
	// InputHintsSpellcheck to enable spellcheck
	InputHintsSpellcheck InputHints = C.GTK_INPUT_HINT_SPELLCHECK
	// InputHintsNoSpellcheck to disable spellcheck
	InputHintsNoSpellcheck InputHints = C.GTK_INPUT_HINT_NO_SPELLCHECK
	// InputHintsWordCompletion to enable word completion
	InputHintsWordCompletion InputHints = C.GTK_INPUT_HINT_WORD_COMPLETION
	// InputHintsLowercase to prefer lowercase
	InputHintsLowercase InputHints = C.GTK_INPUT_HINT_LOWERCASE
)
