// Package gtk4 provides selection mode definitions for GTK4
// File: gtk4go/gtk4/selection.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
import "C"

// SelectionMode defines how items can be selected in a view
type SelectionMode int

const (
	// SelectionModeNone indicates no item can be selected
	SelectionModeNone SelectionMode = 0
	// SelectionModeSingle indicates only one item can be selected
	SelectionModeSingle SelectionMode = 1
	// SelectionModeMultiple indicates multiple items can be selected
	SelectionModeMultiple SelectionMode = 2
)
