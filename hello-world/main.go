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

	// Create buttons with labels
	helloBtn := gtk4.NewButton("Say Hello")
	aboutBtn := gtk4.NewButton("About")
	fileBtn := gtk4.NewButton("Open File")

	// Add CSS classes to the buttons
	helloBtn.AddCssClass("square-button")
	aboutBtn.AddCssClass("square-button")
	fileBtn.AddCssClass("square-button")

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
		.dialog-content-area {
			padding: 16px;
		}
		.dialog-button-area {
			padding: 8px;
			background-color: #f0f0f0;
		}
		.dialog-message {
			font-size: 14px;
			padding: 8px;
		}
		.info-dialog .dialog-message {
			color: #0066cc;
		}
		.warning-dialog .dialog-message {
			color: #ff6600;
		}
		.error-dialog .dialog-message {
			color: #cc0000;
		}
		.question-dialog .dialog-message {
			color: #006633;
		}
		.title {
			font-size: 18px;
			font-weight: bold;
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

	// Connect hello button click event
	helloBtn.ConnectClicked(func() {
		name := entry.GetText()
		if name == "" {
			name = "World"
		}

		// Create a simple info dialog
		dialog := gtk4.NewMessageDialog(
			win,
			gtk4.DialogModal,
			gtk4.MessageInfo,
			gtk4.ResponseOk,
			fmt.Sprintf("Hello, %s! Nice to meet you.", name),
		)
		dialog.SetTitle("Greeting")

		// Connect response handler before showing the dialog
		responseChan := make(chan bool)

		dialog.ConnectResponse(func(responseId gtk4.ResponseType) {
			fmt.Printf("Dialog response: %d\n", responseId)
			dialog.Destroy()
			responseChan <- true
		})

		// Show the dialog
		dialog.Show()

		// Wait for response in a goroutine to not block the UI
		go func() {
			<-responseChan
			fmt.Println("Dialog closed")
		}()

		resultLbl.SetText(fmt.Sprintf("Hello, %s!", name))
	})

	// Connect about button click event
	aboutBtn.ConnectClicked(func() {
		// Create a custom about dialog
		dialog := gtk4.NewDialog("About This Application", win, gtk4.DialogModal|gtk4.DialogDestroyWithParent)

		// Get the content area of the dialog
		content := dialog.GetContentArea()

		// Add some content to the dialog
		titleLabel := gtk4.NewLabel("GTK4Go Demo Application")
		titleLabel.AddCssClass("title")
		descLabel := gtk4.NewLabel("This is a simple demonstration of GTK4 bindings for Go.")
		versionLabel := gtk4.NewLabel("Version: 1.0.0")

		// Add widgets to the content area
		content.Append(titleLabel)
		content.Append(descLabel)
		content.Append(versionLabel)

		// Add padding to the content area
		content.SetSpacing(10)

		// Add OK button to the dialog
		dialog.AddButton("OK", gtk4.ResponseOk)

		// Connect response handler before showing the dialog
		responseChan := make(chan bool)

		dialog.ConnectResponse(func(responseId gtk4.ResponseType) {
			fmt.Printf("About dialog response: %d\n", responseId)
			dialog.Destroy()
			responseChan <- true
		})

		// Show the dialog
		dialog.Show()

		// Wait for response in a goroutine to not block the UI
		go func() {
			<-responseChan
			fmt.Println("About dialog closed")
		}()
	})

	// Connect file button click event
	fileBtn.ConnectClicked(func() {
		// Show a confirmation dialog
		confirmDialog := gtk4.NewMessageDialog(
			win,
			gtk4.DialogModal,
			gtk4.MessageQuestion,
			gtk4.ResponseYes|gtk4.ResponseNo,
			"Do you want to open a file?",
		)
		confirmDialog.SetTitle("Confirm Action")

		// Connect response handler for the confirmation dialog
		confirmChan := make(chan bool)
		var confirmed bool

		confirmDialog.ConnectResponse(func(responseId gtk4.ResponseType) {
			fmt.Printf("Confirm dialog response: %d\n", responseId)
			confirmed = (responseId == gtk4.ResponseYes)
			confirmDialog.Destroy()
			confirmChan <- true
		})

		// Show the confirmation dialog
		confirmDialog.Show()

		// Wait for confirmation response and then potentially show file dialog
		go func() {
			<-confirmChan
			fmt.Println("Confirm dialog closed, confirmed:", confirmed)

			if confirmed {
				// Create file open dialog
				fileDialog := gtk4.NewFileDialog("Select a File", win, gtk4.FileDialogOpen)

				// Connect response handler for the file dialog
				fileChan := make(chan bool)
				var selectedFilename string

				fileDialog.ConnectResponse(func(responseId gtk4.ResponseType) {
					fmt.Printf("File dialog response: %d\n", responseId)
					if responseId == gtk4.ResponseAccept {
						selectedFilename = fileDialog.GetFilename()
						if selectedFilename != "" {
							// This updates the UI with the selected filename
							resultLbl.SetText(fmt.Sprintf("Selected file: %s", selectedFilename))

							// Print to the console/terminal
							fmt.Printf("User selected file: %s\n", selectedFilename)

							// Here you could also process the file if needed
							// For example, read its contents:
							// fileContent, err := os.ReadFile(selectedFilename)
							// if err == nil {
							//     fmt.Printf("File content: %s\n", string(fileContent))
							// }
						}
					}
					fileDialog.Destroy()
					fileChan <- true
				})

				// Show the file dialog
				fileDialog.Show()

				// Wait for file dialog response
				go func() {
					<-fileChan
					fmt.Println("File dialog closed")
					if selectedFilename != "" {
						fmt.Printf("Selected file path: %s\n", selectedFilename)

						// You could trigger additional processing here
						// processFile(selectedFilename)
					}
				}()
			}
		}()
	})

	// Create a horizontal button box for the buttons
	buttonBox := gtk4.NewBox(gtk4.OrientationHorizontal, 5)
	buttonBox.Append(helloBtn)
	buttonBox.Append(aboutBtn)
	buttonBox.Append(fileBtn)

	// Add widgets to the main box with proper spacing
	box.Append(lbl)
	box.Append(entry)
	box.Append(buttonBox)
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
