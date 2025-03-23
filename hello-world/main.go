package main

import (
	"../../gtk4go"
	"../gtk4"
	"fmt"
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

	// Create a vertical box container with 10px spacing
	box := gtk4.NewBox(gtk4.OrientationVertical, 10)

	// Create a label with text
	lbl := gtk4.NewLabel("Hello, World!")

	// Create a button with label
	btn := gtk4.NewButton("Click Me")

	// Connect button click event
	btn.ConnectClicked(func() {
		fmt.Println("clicked")
	})

	// Add widgets to the box
	box.Append(lbl)
	box.Append(btn)

	// Add the box to the window
	win.SetChild(box)

	// Add the window to the application
	app.AddWindow(win)

	// Run the application
	os.Exit(app.Run())
}
