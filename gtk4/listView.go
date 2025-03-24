// In listView.go, add the missing list_store_append function declaration 
// to the C preamble section. Place this near other list_store functions 
// around line 49 after "create_default_store" and before "list_view_create_widget"

// Package gtk4 provides ListView functionality for GTK4
// File: gtk4go/gtk4/listView.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Signal callback functions for list view
// extern void listViewItemActivatedCallback(GtkListView *list_view, guint position, gpointer user_data);
//
// // Connect list view signals
// static gulong list_view_connect_item_activated(GtkListView *list_view, gpointer user_data) {
//     if (list_view == NULL) return 0;
//     return g_signal_connect(G_OBJECT(list_view), "activate", G_CALLBACK(listViewItemActivatedCallback), user_data);
// }
//
// // Selection handling
// extern void listViewSelectionModelChangedCallback(GtkSelectionModel *model, guint position, guint n_items, gpointer user_data);
//
// // Connect selection model signals
// static gulong list_view_connect_selection_changed(GtkSelectionModel *model, gpointer user_data) {
//     if (model == NULL) return 0;
//     return g_signal_connect(G_OBJECT(model), "selection-changed", G_CALLBACK(listViewSelectionModelChangedCallback), user_data);
// }
//
// // Create a default GListStore with GObject type
// static GListStore* create_default_store() {
//     return g_list_store_new(G_TYPE_OBJECT);
// }
//
// // Append to a GListStore - ADD THIS FUNCTION
// static void list_store_append(GListStore* store, gpointer item) {
//     if (store == NULL) return;
//     g_list_store_append(store, item);
// }
//
// // ListView creation and configuration
// static GtkWidget* list_view_create_widget(GtkSelectionModel *model, GtkListItemFactory *factory) {
//     if (model == NULL || factory == NULL) {
//         g_warning("ListView creation failed - null model or factory");
//         return NULL;
//     }
//     return gtk_list_view_new(model, factory);
// }
// 
// // Rest of the file remains the same...