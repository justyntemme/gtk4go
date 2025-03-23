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
	lbl := gtk4.NewLabel("Enter your name:")

	// Create a text entry widget
	entry := gtk4.NewEntry()
	entry.SetPlaceholderText("Type your name here")

	// Create a second label for displaying the entered text
	resultLbl := gtk4.NewLabel("Hello, World!")

	// Create a button with label
	btn := gtk4.NewButton("Say Hello")

	// Add a CSS class to the button for styling
	btn.AddCssClass("square-button")

	// Load CSS for styling
	cssProvider, err := gtk4.LoadCSS(`
		.square-button {
			border-radius: 0;
			padding: 8px 16px;
			background-color: #3584e4;
			color: white;
			font-weight: bold;
		}
		.square-button:hover {
			background-color: #1c71d8;
		}
		window {
			background-color: #f6f5f4;
		}
		entry {
			padding: 8px;
			margin: 4px 0;
		}
		label {
			margin: 4px 0;
		}
	`)
	if err != nil {
		log.Printf("Failed to load CSS: %v", err)
	} else {
		// Apply CSS provider to the entire application
		gtk4.AddProviderForDisplay(cssProvider, uint(gtk4.PriorityApplication))
	}

	// Connect entry activate event (when Enter is pressed)
	entry.ConnectActivate(func() {
		name := entry.GetText()
		if name == "" {
			name = "World"
		}
		resultLbl.SetText(fmt.Sprintf("Hello, %s!", name))
	})

	// Connect button click event
	btn.ConnectClicked(func() {
		name := entry.GetText()
		if name == "" {
			name = "World"
		}
		resultLbl.SetText(fmt.Sprintf("Hello, %s!", name))
		fmt.Println("Said hello to:", name)
	})

	// Connect entry changed event
	entry.ConnectChanged(func() {
		fmt.Println("Text changed to:", entry.GetText())
	})

	// Add widgets to the box with proper spacing
	box.Append(lbl)
	box.Append(entry)
	box.Append(btn)
	box.Append(resultLbl)

	// Add some spacing to make the layout more attractive
	box.SetSpacing(15)

	// Add the box to the window
	win.SetChild(box)

	// Add the window to the application
	app.AddWindow(win)

	// Run the application
	os.Exit(app.Run())
}
