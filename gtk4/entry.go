// Package gtk4 provides entry widget functionality for GTK4
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Helper function to properly handle max length with correct types
// static void set_max_length(GtkEntryBuffer *buffer, unsigned int max_length) {
//     gtk_entry_buffer_set_max_length(buffer, (gsize)max_length);
// }
//
// // Helper function to get max length with correct type conversion
// static unsigned int get_max_length(GtkEntryBuffer *buffer) {
//     return (unsigned int)gtk_entry_buffer_get_max_length(buffer);
// }
import "C"

import (
	"runtime"
	"unsafe"
)

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

// GetEditable gets whether the user can edit the text
func (e *Entry) GetEditable() bool {
	return C.gtk_editable_get_editable((*C.GtkEditable)(unsafe.Pointer(e.widget))) == C.TRUE
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

// GetVisibility gets whether the text is visible
func (e *Entry) GetVisibility() bool {
	return C.gtk_entry_get_visibility((*C.GtkEntry)(unsafe.Pointer(e.widget))) == C.TRUE
}

// ConnectChanged connects a callback function to the entry's "changed" signal
func (e *Entry) ConnectChanged(callback func()) {
	Connect(e, SignalChanged, callback)
}

// ConnectActivate connects a callback function to the entry's "activate" signal
func (e *Entry) ConnectActivate(callback func()) {
	Connect(e, SignalActivate, callback)
}

// DisconnectChanged disconnects the changed signal handler
func (e *Entry) DisconnectChanged() {
	// Since we don't have a specific disconnect function for a single signal type,
	// we'll have to disconnect all signal handlers
	DisconnectAll(e)
}

// DisconnectActivate disconnects the activate signal handler
func (e *Entry) DisconnectActivate() {
	// Since we don't have a specific disconnect function for a single signal type,
	// we'll have to disconnect all signal handlers
	DisconnectAll(e)
}

// Destroy destroys the entry and cleans up resources
func (e *Entry) Destroy() {
	// Disconnect all signals
	DisconnectAll(e)
	
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

// SetMaxLength sets the maximum length of the text in the buffer
func (b *EntryBuffer) SetMaxLength(length uint) {
	// Use our helper function to handle type conversion correctly
	C.set_max_length(b.buffer, C.uint(length))
}

// GetMaxLength gets the maximum length of the text in the buffer
func (b *EntryBuffer) GetMaxLength() uint {
	// Use our helper function to handle type conversion correctly
	return uint(C.get_max_length(b.buffer))
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

// SetInputPurpose sets the purpose of the entry
func (e *Entry) SetInputPurpose(purpose InputPurpose) {
	C.gtk_entry_set_input_purpose((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.GtkInputPurpose(purpose))
}

// GetInputPurpose gets the purpose of the entry
func (e *Entry) GetInputPurpose() InputPurpose {
	return InputPurpose(C.gtk_entry_get_input_purpose((*C.GtkEntry)(unsafe.Pointer(e.widget))))
}

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

// SetInputHints sets the input hints for the entry
func (e *Entry) SetInputHints(hints InputHints) {
	C.gtk_entry_set_input_hints((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.GtkInputHints(hints))
}

// GetInputHints gets the input hints for the entry
func (e *Entry) GetInputHints() InputHints {
	return InputHints(C.gtk_entry_get_input_hints((*C.GtkEntry)(unsafe.Pointer(e.widget))))
}

// SetAlignment sets the alignment for the entry text
func (e *Entry) SetAlignment(xalign float32) {
	C.gtk_entry_set_alignment((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.gfloat(xalign))
}

// GetAlignment gets the alignment for the entry text
func (e *Entry) GetAlignment() float32 {
	return float32(C.gtk_entry_get_alignment((*C.GtkEntry)(unsafe.Pointer(e.widget))))
}

// SetProgressFraction sets the current fraction of the task that's been completed
func (e *Entry) SetProgressFraction(fraction float64) {
	C.gtk_entry_set_progress_fraction((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.gdouble(fraction))
}

// GetProgressFraction gets the current fraction of the task that's been completed
func (e *Entry) GetProgressFraction() float64 {
	return float64(C.gtk_entry_get_progress_fraction((*C.GtkEntry)(unsafe.Pointer(e.widget))))
}

// SetProgressPulseStep sets the fraction of total entry width to move the progress bouncing block
func (e *Entry) SetProgressPulseStep(fraction float64) {
	C.gtk_entry_set_progress_pulse_step((*C.GtkEntry)(unsafe.Pointer(e.widget)), C.gdouble(fraction))
}

// GetProgressPulseStep gets the fraction of total entry width to move the progress bouncing block
func (e *Entry) GetProgressPulseStep() float64 {
	return float64(C.gtk_entry_get_progress_pulse_step((*C.GtkEntry)(unsafe.Pointer(e.widget))))
}

// ProgressPulse causes the entry's progress indicator to enter "activity mode"
func (e *Entry) ProgressPulse() {
	C.gtk_entry_progress_pulse((*C.GtkEntry)(unsafe.Pointer(e.widget)))
}

// SetEnableUndo sets whether the user can undo/redo entry edits
func (e *Entry) SetEnableUndo(enabled bool) {
	var cenabled C.gboolean
	if enabled {
		cenabled = C.TRUE
	} else {
		cenabled = C.FALSE
	}
	C.gtk_editable_set_enable_undo((*C.GtkEditable)(unsafe.Pointer(e.widget)), cenabled)
}

// GetEnableUndo gets whether the user can undo/redo entry edits
func (e *Entry) GetEnableUndo() bool {
	return C.gtk_editable_get_enable_undo((*C.GtkEditable)(unsafe.Pointer(e.widget))) == C.TRUE
}