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

// EntryOption is a function that configures an entry
type EntryOption func(*Entry)

// Entry represents a GTK entry widget for text input
type Entry struct {
	BaseWidget
}

// NewEntry creates a new GTK entry widget
func NewEntry(options ...EntryOption) *Entry {
	entry := &Entry{
		BaseWidget: BaseWidget{
			widget: C.gtk_entry_new(),
		},
	}

	// Apply options
	for _, option := range options {
		option(entry)
	}

	SetupFinalization(entry, entry.Destroy)
	return entry
}

// WithEntryBuffer creates an entry with a specific buffer
func WithEntryBuffer(buffer *EntryBuffer) EntryOption {
	return func(e *Entry) {
		e.widget = C.gtk_entry_new_with_buffer(buffer.buffer)
	}
}

// WithPlaceholderText sets placeholder text
func WithPlaceholderText(text string) EntryOption {
	return func(e *Entry) {
		e.SetPlaceholderText(text)
	}
}

// WithEditable sets whether the entry is editable
func WithEditable(editable bool) EntryOption {
	return func(e *Entry) {
		e.SetEditable(editable)
	}
}

// SetText sets the text in the entry
func (e *Entry) SetText(text string) {
	WithCString(text, func(cText *C.char) {
		C.gtk_editable_set_text((*C.GtkEditable)(unsafe.Pointer(e.widget)), cText)
	})
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
	WithCString(text, func(cText *C.char) {
		C.gtk_entry_set_placeholder_text((*C.GtkEntry)(unsafe.Pointer(e.widget)), cText)
	})
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
	var ceditable C.gboolean
	if editable {
		ceditable = C.TRUE
	} else {
		ceditable = C.FALSE
	}
	C.gtk_editable_set_editable((*C.GtkEditable)(unsafe.Pointer(e.widget)), ceditable)
}

// SetVisibility sets whether the text is visible or hidden (e.g., for passwords)
func (e *Entry) SetVisibility(visible bool) {
	var cvisible C.gboolean
	if visible {
		cvisible = C.TRUE
	} else {
		cvisible = C.FALSE
	}
	C.gtk_entry_set_visibility((*C.GtkEntry)(unsafe.Pointer(e.widget)), cvisible)
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

// ConnectActivate connects a callback function to the entry's "activate" signal
func (e *Entry) ConnectActivate(callback EntryCallback) {
	entryCallbackMutex.Lock()
	defer entryCallbackMutex.Unlock()

	// Store callback in map
	entryPtr := uintptr(unsafe.Pointer(e.widget))
	entryActivateCallbacks[entryPtr] = callback

	// Connect signal
	C.connectEntryActivate(e.widget, C.gpointer(unsafe.Pointer(e.widget)))
}

// Destroy destroys the entry and cleans up resources
func (e *Entry) Destroy() {
	entryCallbackMutex.Lock()
	defer entryCallbackMutex.Unlock()

	// Remove callbacks from maps if they exist
	entryPtr := uintptr(unsafe.Pointer(e.widget))
	delete(entryChangedCallbacks, entryPtr)
	delete(entryActivateCallbacks, entryPtr)

	// Call base destroy method
	e.BaseWidget.Destroy()
}

// EntryBuffer represents a GTK entry buffer
type EntryBuffer struct {
	buffer *C.GtkEntryBuffer
}

// NewEntryBuffer creates a new entry buffer with initial text
func NewEntryBuffer(initialText string) *EntryBuffer {
	var buffer *C.GtkEntryBuffer

	WithCString(initialText, func(cText *C.char) {
		buffer = C.gtk_entry_buffer_new(cText, C.int(len(initialText)))
	})

	entryBuffer := &EntryBuffer{
		buffer: buffer,
	}

	// Use a simple finalizer
	runtime.SetFinalizer(entryBuffer, func(b *EntryBuffer) {
		b.Free()
	})

	return entryBuffer
}

// SetText sets the text in the buffer
func (b *EntryBuffer) SetText(text string) {
	WithCString(text, func(cText *C.char) {
		C.gtk_entry_buffer_set_text(b.buffer, cText, C.int(len(text)))
	})
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
	if b.buffer != nil {
		C.g_object_unref(C.gpointer(unsafe.Pointer(b.buffer)))
		b.buffer = nil
	}
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
