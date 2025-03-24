// Package gtk4 provides grid layout functionality for GTK4
// File: gtk4go/gtk4/grid.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// GridOption is a function that configures a grid
type GridOption func(*Grid)

// Grid represents a GTK grid container
type Grid struct {
	BaseWidget
}

// NewGrid creates a new GTK grid container
func NewGrid(options ...GridOption) *Grid {
	grid := &Grid{
		BaseWidget: BaseWidget{
			widget: C.gtk_grid_new(),
		},
	}

	// Apply options
	for _, option := range options {
		option(grid)
	}

	SetupFinalization(grid, grid.Destroy)
	return grid
}

// WithRowSpacing sets spacing between rows
func WithRowSpacing(spacing int) GridOption {
	return func(g *Grid) {
		C.gtk_grid_set_row_spacing((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.guint(spacing))
	}
}

// WithColumnSpacing sets spacing between columns
func WithColumnSpacing(spacing int) GridOption {
	return func(g *Grid) {
		C.gtk_grid_set_column_spacing((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.guint(spacing))
	}
}

// WithRowHomogeneous sets whether all rows have the same height
func WithRowHomogeneous(homogeneous bool) GridOption {
	return func(g *Grid) {
		var chomogeneous C.gboolean
		if homogeneous {
			chomogeneous = C.TRUE
		} else {
			chomogeneous = C.FALSE
		}
		C.gtk_grid_set_row_homogeneous((*C.GtkGrid)(unsafe.Pointer(g.widget)), chomogeneous)
	}
}

// WithColumnHomogeneous sets whether all columns have the same width
func WithColumnHomogeneous(homogeneous bool) GridOption {
	return func(g *Grid) {
		var chomogeneous C.gboolean
		if homogeneous {
			chomogeneous = C.TRUE
		} else {
			chomogeneous = C.FALSE
		}
		C.gtk_grid_set_column_homogeneous((*C.GtkGrid)(unsafe.Pointer(g.widget)), chomogeneous)
	}
}

// Attach attaches a widget to the grid at the given position
func (g *Grid) Attach(child Widget, column, row, width, height int) {
	C.gtk_grid_attach(
		(*C.GtkGrid)(unsafe.Pointer(g.widget)),
		child.GetWidget(),
		C.int(column),
		C.int(row),
		C.int(width),
		C.int(height),
	)
}

// AttachNextTo attaches a widget to the grid, next to another widget
func (g *Grid) AttachNextTo(child, sibling Widget, side GridPosition, width, height int) {
	C.gtk_grid_attach_next_to(
		(*C.GtkGrid)(unsafe.Pointer(g.widget)),
		child.GetWidget(),
		sibling.GetWidget(),
		C.GtkPositionType(side),
		C.int(width),
		C.int(height),
	)
}

// InsertRow inserts a row at the specified position
func (g *Grid) InsertRow(position int) {
	C.gtk_grid_insert_row((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.int(position))
}

// InsertColumn inserts a column at the specified position
func (g *Grid) InsertColumn(position int) {
	C.gtk_grid_insert_column((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.int(position))
}

// RemoveRow removes a row from the grid
func (g *Grid) RemoveRow(position int) {
	C.gtk_grid_remove_row((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.int(position))
}

// RemoveColumn removes a column from the grid
func (g *Grid) RemoveColumn(position int) {
	C.gtk_grid_remove_column((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.int(position))
}

// GetChildAt gets the child at the specified position
func (g *Grid) GetChildAt(column, row int) Widget {
	widget := C.gtk_grid_get_child_at(
		(*C.GtkGrid)(unsafe.Pointer(g.widget)),
		C.int(column),
		C.int(row),
	)

	// We can't directly return a Widget from the C pointer
	// Instead, we would need to wrap it in an appropriate Go struct
	// This is a simplified implementation that returns nil
	if widget == nil {
		return nil
	}

	// In a real implementation, we would determine the type of widget
	// and return an appropriate Go wrapper
	return nil
}

// SetRowSpacing sets the amount of space between rows
func (g *Grid) SetRowSpacing(spacing int) {
	C.gtk_grid_set_row_spacing((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.guint(spacing))
}

// GetRowSpacing gets the amount of space between rows
func (g *Grid) GetRowSpacing() int {
	return int(C.gtk_grid_get_row_spacing((*C.GtkGrid)(unsafe.Pointer(g.widget))))
}

// SetColumnSpacing sets the amount of space between columns
func (g *Grid) SetColumnSpacing(spacing int) {
	C.gtk_grid_set_column_spacing((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.guint(spacing))
}

// GetColumnSpacing gets the amount of space between columns
func (g *Grid) GetColumnSpacing() int {
	return int(C.gtk_grid_get_column_spacing((*C.GtkGrid)(unsafe.Pointer(g.widget))))
}

// SetRowHomogeneous sets whether all rows should be the same height
func (g *Grid) SetRowHomogeneous(homogeneous bool) {
	var chomogeneous C.gboolean
	if homogeneous {
		chomogeneous = C.TRUE
	} else {
		chomogeneous = C.FALSE
	}
	C.gtk_grid_set_row_homogeneous((*C.GtkGrid)(unsafe.Pointer(g.widget)), chomogeneous)
}

// GetRowHomogeneous gets whether all rows are the same height
func (g *Grid) GetRowHomogeneous() bool {
	return C.gtk_grid_get_row_homogeneous((*C.GtkGrid)(unsafe.Pointer(g.widget))) == C.TRUE
}

// SetColumnHomogeneous sets whether all columns should be the same width
func (g *Grid) SetColumnHomogeneous(homogeneous bool) {
	var chomogeneous C.gboolean
	if homogeneous {
		chomogeneous = C.TRUE
	} else {
		chomogeneous = C.FALSE
	}
	C.gtk_grid_set_column_homogeneous((*C.GtkGrid)(unsafe.Pointer(g.widget)), chomogeneous)
}

// GetColumnHomogeneous gets whether all columns are the same width
func (g *Grid) GetColumnHomogeneous() bool {
	return C.gtk_grid_get_column_homogeneous((*C.GtkGrid)(unsafe.Pointer(g.widget))) == C.TRUE
}

// GridPosition defines the position relative to another widget
type GridPosition int

const (
	// PositionLeft positions a widget to the left of another widget
	PositionLeft GridPosition = C.GTK_POS_LEFT
	// PositionRight positions a widget to the right of another widget
	PositionRight GridPosition = C.GTK_POS_RIGHT
	// PositionTop positions a widget above another widget
	PositionTop GridPosition = C.GTK_POS_TOP
	// PositionBottom positions a widget below another widget
	PositionBottom GridPosition = C.GTK_POS_BOTTOM
)
