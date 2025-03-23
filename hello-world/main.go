package main

import (
	"../../gtk4go"
	"../gtk4"
	"log"
	"os"
)

func main() {
	// Initialize GTK (this is also done automatically on import)
	if err := gtk4go.Initialize(); err != nil {
		log.Fatalf("Failed to initialize GTK: %v", err)
	}

	// Create a new application
	app := gtk4.NewApplication("com.example.HelloWorld")

	// Create a window
	win := gtk4.NewWindow("Hello GTK4 from Go!")
	win.SetDefaultSize(400, 300)

	// Create a label with text
	lbl := gtk4.NewLabel("Hello, World!")

	// Add the label to the window
	win.SetChild(lbl)

	// Add the window to the application (this connects activate signal)
	app.AddWindow(win)

	// Note: Don't call win.Show() or win.Present() here!
	// The window will be shown automatically when the application activates

	// Run the application
	os.Exit(app.Run())
}
