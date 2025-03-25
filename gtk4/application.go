// Package gtk4 provides application functionality for GTK4
// File: gtk4go/gtk4/application.go
package gtk4

// #cgo pkg-config: gtk4
// #include <gtk/gtk.h>
// #include <stdlib.h>
//
// // Callback struct to pass Go function through C
// typedef struct {
//     void* window;
//     void* app;
// } ActivateData;
//
// // Activate callback
// static void activateCallback(GtkApplication* app, gpointer user_data) {
//     ActivateData* data = (ActivateData*)user_data;
//     GtkWidget* window = (GtkWidget*)data->window;
//
//     // Set application
//     gtk_window_set_application(GTK_WINDOW(window), app);
//
//     // Show the window
//     gtk_widget_set_visible(window, TRUE);
// }
//
// // Connect activate signal
// static void connect_activate(GtkApplication* app, GtkWidget* window) {
//     ActivateData* data = malloc(sizeof(ActivateData));
//     data->window = window;
//     data->app = app;
//
//     g_signal_connect(app, "activate", G_CALLBACK(activateCallback), data);
// }
import "C"

import (
	"runtime"
	"unsafe"
)

// ApplicationOption is a function that configures an application
type ApplicationOption func(*Application)

// Application represents a GTK application
type Application struct {
	app *C.GtkApplication
}

// NewApplication creates a new GTK application with the given ID
func NewApplication(id string, options ...ApplicationOption) *Application {
	cID := C.CString(id)
	defer C.free(unsafe.Pointer(cID))

	app := &Application{
		app: C.gtk_application_new(cID, C.G_APPLICATION_DEFAULT_FLAGS),
	}
	
	// Apply options
	for _, option := range options {
		option(app)
	}
	
	runtime.SetFinalizer(app, (*Application).Destroy)
	return app
}

// WithAction adds an action to the application at creation time
func WithAction(action *Action) ApplicationOption {
	return func(a *Application) {
		a.GetActionGroup().AddAction(action)
	}
}

// WithMenuBar sets the application menu bar at creation time
func WithMenuBar(menu *Menu) ApplicationOption {
	return func(a *Application) {
		a.SetMenuBar(menu)
	}
}

// AddWindow adds a window to the application and connects the activate signal
func (a *Application) AddWindow(window any) {
	if w, ok := window.(interface{ GetWidget() *C.GtkWidget }); ok {
		// Connect the activate signal to handle window display
		C.connect_activate(a.app, w.GetWidget())
	}
}

// Run runs the application
func (a *Application) Run() int {
	status := C.g_application_run((*C.GApplication)(unsafe.Pointer(a.app)), 0, nil)
	return int(status)
}

// Destroy destroys the application
func (a *Application) Destroy() {
	C.g_object_unref(C.gpointer(unsafe.Pointer(a.app)))
}