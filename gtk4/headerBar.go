// Package gtk4 provides header bar functionality for GTK4
// File: gtk4go/gtk4/headerBar.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Helper functions for GtkHeaderBar
// static GtkWidget* createHeaderBar() {
//     return gtk_header_bar_new();
// }
//
// static void setHeaderBarShowTitleButtons(GtkHeaderBar *header_bar, gboolean setting) {
//     gtk_header_bar_set_show_title_buttons(header_bar, setting);
// }
//
// // In GTK4, we need to use label widgets for title and subtitle
// static void setHeaderBarTitle(GtkHeaderBar *header_bar, const char *title) {
//     // Remove any existing title widget
//     GtkWidget *current = gtk_header_bar_get_title_widget(header_bar);
//     if (current != NULL) {
//         gtk_header_bar_set_title_widget(header_bar, NULL);
//     }
//     
//     // Create a new label with the title
//     GtkWidget *label = gtk_label_new(title);
//     // Make it bold and larger
//     PangoAttrList *attrs = pango_attr_list_new();
//     pango_attr_list_insert(attrs, pango_attr_weight_new(PANGO_WEIGHT_BOLD));
//     pango_attr_list_insert(attrs, pango_attr_scale_new(1.2));
//     gtk_label_set_attributes(GTK_LABEL(label), attrs);
//     pango_attr_list_unref(attrs);
//     
//     // Set it as the title widget
//     gtk_header_bar_set_title_widget(header_bar, label);
// }
//
// static void setHeaderBarTitleWidget(GtkHeaderBar *header_bar, GtkWidget *title_widget) {
//     gtk_header_bar_set_title_widget(header_bar, title_widget);
// }
//
// static void setHeaderBarDecorationLayout(GtkHeaderBar *header_bar, const char *layout) {
//     gtk_header_bar_set_decoration_layout(header_bar, layout);
// }
//
// static void packStart(GtkHeaderBar *header_bar, GtkWidget *child) {
//     gtk_header_bar_pack_start(header_bar, child);
// }
//
// static void packEnd(GtkHeaderBar *header_bar, GtkWidget *child) {
//     gtk_header_bar_pack_end(header_bar, child);
// }
import "C"

import (
    "unsafe"
)

// HeaderBarOption is a function that configures a header bar
type HeaderBarOption func(*HeaderBar)

// HeaderBar represents a GTK header bar
type HeaderBar struct {
    BaseWidget
}

// NewHeaderBar creates a new GTK header bar
func NewHeaderBar(options ...HeaderBarOption) *HeaderBar {
    headerBar := &HeaderBar{
        BaseWidget: BaseWidget{
            widget: C.createHeaderBar(),
        },
    }

    // Apply options
    for _, option := range options {
        option(headerBar)
    }

    SetupFinalization(headerBar, headerBar.Destroy)
    return headerBar
}

// WithShowTitleButtons sets whether to show title buttons (minimize, maximize, close)
func WithShowTitleButtons(show bool) HeaderBarOption {
    return func(hb *HeaderBar) {
        var cshow C.gboolean
        if show {
            cshow = C.TRUE
        } else {
            cshow = C.FALSE
        }
        C.setHeaderBarShowTitleButtons((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), cshow)
    }
}

// WithTitle sets the header bar title
func WithTitle(title string) HeaderBarOption {
    return func(hb *HeaderBar) {
        cTitle := C.CString(title)
        defer C.free(unsafe.Pointer(cTitle))
        C.setHeaderBarTitle((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), cTitle)
    }
}

// WithTitleWidget sets a custom widget as the header bar title
func WithTitleWidget(titleWidget Widget) HeaderBarOption {
    return func(hb *HeaderBar) {
        C.setHeaderBarTitleWidget((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), titleWidget.GetWidget())
    }
}

// WithDecorationLayout sets the header bar decoration layout
func WithDecorationLayout(layout string) HeaderBarOption {
    return func(hb *HeaderBar) {
        cLayout := C.CString(layout)
        defer C.free(unsafe.Pointer(cLayout))
        C.setHeaderBarDecorationLayout((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), cLayout)
    }
}

// SetShowTitleButtons sets whether to show title buttons (minimize, maximize, close)
func (hb *HeaderBar) SetShowTitleButtons(show bool) {
    var cshow C.gboolean
    if show {
        cshow = C.TRUE
    } else {
        cshow = C.FALSE
    }
    C.setHeaderBarShowTitleButtons((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), cshow)
}

// SetTitle sets the header bar title
func (hb *HeaderBar) SetTitle(title string) {
    cTitle := C.CString(title)
    defer C.free(unsafe.Pointer(cTitle))
    C.setHeaderBarTitle((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), cTitle)
}

// SetTitleWidget sets a custom widget as the header bar title
func (hb *HeaderBar) SetTitleWidget(titleWidget Widget) {
    C.setHeaderBarTitleWidget((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), titleWidget.GetWidget())
}

// SetDecorationLayout sets the header bar decoration layout
func (hb *HeaderBar) SetDecorationLayout(layout string) {
    cLayout := C.CString(layout)
    defer C.free(unsafe.Pointer(cLayout))
    C.setHeaderBarDecorationLayout((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), cLayout)
}

// PackStart adds a widget to the start of the header bar
func (hb *HeaderBar) PackStart(child Widget) {
    C.packStart((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), child.GetWidget())
}

// PackEnd adds a widget to the end of the header bar
func (hb *HeaderBar) PackEnd(child Widget) {
    C.packEnd((*C.GtkHeaderBar)(unsafe.Pointer(hb.widget)), child.GetWidget())
}