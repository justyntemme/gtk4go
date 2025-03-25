package main

import (
	"../../gtk4go"
	"../gtk4"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Initialize GTK (this is also done automatically on import)
	if err := gtk4go.Initialize(); err != nil {
		log.Fatalf("Failed to initialize GTK: %v", err)
	}

	// Create a new application
	app := gtk4.NewApplication("com.example.HelloWorld")

	// Create a window with optimized rendering
	win := gtk4.NewWindow("Hello GTK4 from Go!")
	win.SetDefaultSize(900, 650)

	// Enable hardware-accelerated rendering
	win.EnableAcceleratedRendering()

	// Set up CSS optimization during window resize
	win.SetupCSSOptimizedResize()

	// Optimize for resizing specifically
	win.OptimizeForResizing()

	// Create a vertical box container as the main layout
	mainBox := gtk4.NewBox(gtk4.OrientationVertical, 10)

	// Create a menu bar for the application
	menuBar := gtk4.NewMenuBar()
	mainBox.Append(menuBar)

	// Create a title label
	titleLabel := gtk4.NewLabel("GTK4Go Demo Application")
	titleLabel.AddCssClass("title")
	mainBox.Append(titleLabel)

	// Create a horizontal paned container to split the UI
	paned := gtk4.NewPaned(gtk4.OrientationHorizontal,
		gtk4.WithPosition(350),
		gtk4.WithWideHandle(true),
	)

	// ---- LEFT SIDE OF PANED ----
	leftBox := gtk4.NewBox(gtk4.OrientationVertical, 10)

	// Add a menu button with popup menu
	menuButton := gtk4.NewMenuButton()
	menuButton.SetLabel("Quick Actions")
	leftBox.Append(menuButton)

	// Input section
	inputLabel := gtk4.NewLabel("Enter your name:")
	entry := gtk4.NewEntry()
	entry.SetPlaceholderText("Type your name here")
	resultLbl := gtk4.NewLabel("Hello, World!")

	// Now use a Grid for button layout
	buttonsGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(10),
		gtk4.WithColumnSpacing(10),
		gtk4.WithColumnHomogeneous(true),
	)

	// Create buttons
	helloBtn := gtk4.NewButton("Say Hello")
	aboutBtn := gtk4.NewButton("About")
	fileBtn := gtk4.NewButton("Open File")
	longTaskBtn := gtk4.NewButton("Run Long Task")

	// Add buttons to grid (col, row, width, height)
	buttonsGrid.Attach(helloBtn, 0, 0, 1, 1)
	buttonsGrid.Attach(aboutBtn, 1, 0, 1, 1)
	buttonsGrid.Attach(fileBtn, 0, 1, 1, 1)
	buttonsGrid.Attach(longTaskBtn, 1, 1, 1, 1)

	// Progress label
	progressLbl := gtk4.NewLabel("Ready")
	progressLbl.AddCssClass("progress-label")

	// Add widgets to left box
	leftBox.Append(inputLabel)
	leftBox.Append(entry)
	leftBox.Append(buttonsGrid)
	leftBox.Append(resultLbl)
	leftBox.Append(progressLbl)

	// ---- RIGHT SIDE OF PANED ----

	// Create a Stack for different content pages
	rightStack := gtk4.NewStack(
		gtk4.WithTransitionType(gtk4.StackTransitionTypeSlideLeftRight),
		gtk4.WithTransitionDuration(200),
	)

	// Stack Page 1: Info Page
	infoBox := gtk4.NewBox(gtk4.OrientationVertical, 10)
	infoBox.Append(gtk4.NewLabel("GTK4Go Information"))
	infoBox.Append(gtk4.NewLabel("This demo showcases the new layout containers:"))

	// Use a grid to display information about widgets
	infoGrid := gtk4.NewGrid(
		gtk4.WithRowSpacing(5),
		gtk4.WithColumnSpacing(10),
		gtk4.WithRowHomogeneous(false),
	)

	// Add headers
	infoGrid.Attach(gtk4.NewLabel("Widget"), 0, 0, 1, 1)
	infoGrid.Attach(gtk4.NewLabel("Description"), 1, 0, 1, 1)

	// Add widget information rows
	widgets := []string{"Grid", "Paned", "Stack", "StackSwitcher", "ScrolledWindow", "ListView"}
	descriptions := []string{
		"Arranges widgets in rows and columns",
		"Divides space between two widgets with adjustable separator",
		"Shows one widget at a time with transitions",
		"Provides buttons to switch between stack pages",
		"Provides scrolling for large content",
		"Displays items from a list model with customizable presentation",
	}

	for i, widget := range widgets {
		widgetLabel := gtk4.NewLabel(widget)
		widgetLabel.AddCssClass("info-widget")
		descLabel := gtk4.NewLabel(descriptions[i])
		descLabel.AddCssClass("info-desc")

		infoGrid.Attach(widgetLabel, 0, i+1, 1, 1)
		infoGrid.Attach(descLabel, 1, i+1, 1, 1)
	}

	infoBox.Append(infoGrid)
	rightStack.AddTitled(infoBox, "info", "Information")

	// Stack Page 2: Log Page with ScrolledWindow
	scrollWin := gtk4.NewScrolledWindow(
		gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
		gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAlways),
		gtk4.WithPropagateNaturalHeight(false), // Don't propagate natural height to allow scrolling
	)

	// Create a vertical box for log entries
	logBox := gtk4.NewBox(gtk4.OrientationVertical, 5)

	// Add some sample log entries
	for i := 1; i <= 10; i++ {
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%d] Log entry #%d", i, i))
		logEntry.AddCssClass("log-entry")
		logBox.Append(logEntry)
	}

	scrollWin.SetChild(logBox)
	rightStack.AddTitled(scrollWin, "logs", "Logs")

	// Stack Page 3: Help Page
	helpBox := gtk4.NewBox(gtk4.OrientationVertical, 10)
	helpBox.Append(gtk4.NewLabel("Help Information"))

	helpText := gtk4.NewLabel(`
Using this application:

1. Enter your name in the text field
2. Click "Say Hello" to see a greeting
3. Click "About" to learn about the app
4. Click "Open File" to select a file
5. Click "Run Long Task" to see a background task
6. Go to ListView tab to see the new ListView widget in action

This demo showcases GTK4Go's layout containers and widgets.
	`)

	helpBox.Append(helpText)
	rightStack.AddTitled(helpBox, "help", "Help")

	// Stack Page 4: ListView Example (NEW)
	listViewBox := gtk4.NewBox(gtk4.OrientationVertical, 10)
	listViewTitle := gtk4.NewLabel("ListView Example")
	listViewTitle.AddCssClass("subtitle")
	listViewBox.Append(listViewTitle)
	
	// Add description
	listViewDesc := gtk4.NewLabel("This demonstrates the new ListView widget with data binding and selection")
	listViewBox.Append(listViewDesc)

	// Create controls for ListView
	listViewControls := gtk4.NewBox(gtk4.OrientationHorizontal, 6)
	listViewControls.AddCssClass("controls-box")

	// Add button
	addItemBtn := gtk4.NewButton("Add Item")
	listViewControls.Append(addItemBtn)

	// Remove button
	removeItemBtn := gtk4.NewButton("Remove Selected")
	listViewControls.Append(removeItemBtn)

	// Clear button
	clearItemsBtn := gtk4.NewButton("Clear All")
	listViewControls.Append(clearItemsBtn)
	
	listViewBox.Append(listViewControls)

	// Create a string list model with sample data
	listModel := gtk4.NewStringList()
	for i := 1; i <= 15; i++ {
		listModel.Append(fmt.Sprintf("List Item %d", i))
	}

	// Create a selection model (SingleSelection)
	selectionModel := gtk4.NewSingleSelection(listModel, 
		gtk4.WithAutoselect(true),
		gtk4.WithInitialSelection(0),
	)

	// Create a list item factory
	factory := gtk4.NewSignalListItemFactory()

	// Setup list items with setup callback
	factory.ConnectSetup(func(listItem *gtk4.ListItem) {
		// Create a box for layout
		box := gtk4.NewBox(gtk4.OrientationHorizontal, 10)
		box.SetHExpand(true)
		box.AddCssClass("list-item-box")

		// Create an icon for visual interest
		icon := gtk4.NewLabel("•")
		icon.AddCssClass("list-item-icon")
		box.Append(icon)
		
		// Create a label for the item text with initial text
		// We'll use CSS to control the appearance based on position
		label := gtk4.NewLabel("List Item")
		label.AddCssClass("list-item-label")
		box.Append(label)
		
		// Set the box as the child of the list item
		listItem.SetChild(box)
	})

	// Bind data to list items with bind callback
	factory.ConnectBind(func(listItem *gtk4.ListItem) {
		// Get the box from the list item
		boxWidget := listItem.GetChild()
		
		// Get the position for reference
		position := listItem.GetPosition()
		
		// In a real implementation, we'd find the label inside the box and update its text
		// For the demo, we'll modify both style classes to reflect selection state
		
		// Add selected class if the item is selected
		if listItem.GetSelected() {
			boxWidget.AddCssClass("selected")
		} else {
			boxWidget.RemoveCssClass("selected")
		}
		
		// Since we can't easily update children, we'll set a CSS class with the position
		// and use that in the CSS to show different styles for different items
		boxWidget.AddCssClass(fmt.Sprintf("item-%d", position))
	})

	// Create list view with the selection model and factory
	listView := gtk4.NewListView(selectionModel, factory,
		gtk4.WithShowSeparators(true),
		gtk4.WithSingleClickActivate(true),
	)

	// Connect the list view activate signal
	listView.ConnectActivate(func(position int) {
		// Log the activation
		activateLog := fmt.Sprintf("[%s] ListView: Item activated at position %d", 
			time.Now().Format("15:04:05"), position)
		
		// Create a log entry
		logEntry := gtk4.NewLabel(activateLog)
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)
		
		// Show a dialog with the selected item
		messageDialog := gtk4.NewMessageDialog(
			win,
			gtk4.DialogModal,
			gtk4.MessageInfo,
			gtk4.ResponseOk,
			fmt.Sprintf("You selected item at position %d", position),
		)
		messageDialog.SetTitle("ListView Item Selected")
		
		// Connect response handler
		messageDialog.ConnectResponse(func(responseId gtk4.ResponseType) {
			messageDialog.Destroy()
		})
		
		// Show the dialog
		messageDialog.Show()
	})

	// Create a scrolled window to contain the list view
	listViewScroll := gtk4.NewScrolledWindow(
		gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyNever),
		gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
	)
	listViewScroll.SetChild(listView)
	listViewScroll.AddCssClass("list-view-container")
	
	// Add the list view scrolled window to the box
	listViewBox.Append(listViewScroll)
	
	// Connect the add item button
	addItemBtn.ConnectClicked(func() {
		// Add a new item to the list model
		newItem := fmt.Sprintf("New List Item %d", listModel.GetNItems()+1)
		listModel.Append(newItem)
		
		// Log the action
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Added new list item: %s", 
			time.Now().Format("15:04:05"), newItem))
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)
	})
	
	// Connect the remove item button
	removeItemBtn.ConnectClicked(func() {
		// Get the selected position
		selectedPos := selectionModel.GetSelected()
		
		// Check if there's a valid selection
		if selectedPos >= 0 && selectedPos < listModel.GetNItems() {
			// Get the item text before removing it
			itemText := fmt.Sprintf("Item %d", selectedPos+1)
			
			// Remove the item
			listModel.Remove(selectedPos)
			
			// Log the action
			logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Removed list item: %s", 
				time.Now().Format("15:04:05"), itemText))
			logEntry.AddCssClass("log-entry")
			logBox.Prepend(logEntry)
		}
	})
	
	// Connect the clear items button
	clearItemsBtn.ConnectClicked(func() {
		// Clear all items by removing them one by one from the end
		for i := listModel.GetNItems() - 1; i >= 0; i-- {
			listModel.Remove(i)
		}
		
		// Log the action
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Cleared all list items", 
			time.Now().Format("15:04:05")))
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)
		
		// Add a default item back
		listModel.Append("List Empty")
	})

	// Add the ListView page to the stack
	rightStack.AddTitled(listViewBox, "listview", "ListView")

	// Create a stack switcher for the right stack
	stackSwitcher := gtk4.NewStackSwitcher(rightStack)

	// Create a box to hold the stack switcher and stack
	rightBox := gtk4.NewBox(gtk4.OrientationVertical, 5)
	rightBox.Append(stackSwitcher)
	rightBox.Append(rightStack)

	// Add left and right sides to the paned container
	paned.SetStartChild(leftBox)
	paned.SetEndChild(rightBox)

	// Add paned container to main box
	mainBox.Append(paned)

	// Add CSS classes to the buttons
	helloBtn.AddCssClass("square-button")
	aboutBtn.AddCssClass("square-button")
	fileBtn.AddCssClass("square-button")
	longTaskBtn.AddCssClass("square-button")

	// Load CSS for styling - using caching for better performance
	cssProvider, err := gtk4.LoadCSS(`
		.title {
			font-size: 18px;
			font-weight: bold;
			padding: 10px;
			color: #2a76c6;
		}
		.subtitle {
			font-size: 16px;
			font-weight: bold;
			padding: 8px;
			color: #2a76c6;
		}
		.square-button {
			border-radius: 4px;
			padding: 8px 16px;
			background-color: #3584e4;
			color: white;
			font-weight: bold;
		}
		.square-button:hover {
			background-color: #1c71d8;
		}
		.square-button.disabled {
			opacity: 0.6;
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
		.progress-label {
			font-style: italic;
			color: #666666;
		}
		.info-widget {
			font-weight: bold;
			color: #2a76c6;
		}
		.info-desc {
			color: #333333;
		}
		.log-entry {
			font-family: monospace;
			padding: 2px 5px;
			text-align: left;
			border-bottom: 1px solid #e0e0e0;
		}
		.log-entry:nth-child(odd) {
			background-color: #f5f5f5;
		}
		.list-view-container {
			border: 1px solid #cccccc;
			border-radius: 4px;
			min-height: 250px;
		}
		.list-item-box {
			padding: 8px;
			border-radius: 3px;
			transition: background-color 0.2s ease;
		}
		.list-item-box.selected {
			background-color: #3584e4;
			color: white;
		}
		.list-item-icon {
			font-size: 16px;
			color: #3584e4;
			font-weight: bold;
		}
		.list-item-box.selected .list-item-icon {
			color: white;
		}
		.list-item-label {
			font-size: 14px;
		}
		/* Add position-based styling using the item-X classes */
		.list-item-box.item-0 .list-item-label:after { content: " 1"; }
		.list-item-box.item-1 .list-item-label:after { content: " 2"; }
		.list-item-box.item-2 .list-item-label:after { content: " 3"; }
		.list-item-box.item-3 .list-item-label:after { content: " 4"; }
		.list-item-box.item-4 .list-item-label:after { content: " 5"; }
		.list-item-box.item-5 .list-item-label:after { content: " 6"; }
		.list-item-box.item-6 .list-item-label:after { content: " 7"; }
		.list-item-box.item-7 .list-item-label:after { content: " 8"; }
		.list-item-box.item-8 .list-item-label:after { content: " 9"; }
		.list-item-box.item-9 .list-item-label:after { content: " 10"; }
		.list-item-box.item-10 .list-item-label:after { content: " 11"; }
		.list-item-box.item-11 .list-item-label:after { content: " 12"; }
		.list-item-box.item-12 .list-item-label:after { content: " 13"; }
		.list-item-box.item-13 .list-item-label:after { content: " 14"; }
		.list-item-box.item-14 .list-item-label:after { content: " 15"; }
		.list-item-box.item-15 .list-item-label:after { content: " 16"; }
		.list-item-box.item-16 .list-item-label:after { content: " 17"; }
		.list-item-box.item-17 .list-item-label:after { content: " 18"; }
		.list-item-box.item-18 .list-item-label:after { content: " 19"; }
		.list-item-box.item-19 .list-item-label:after { content: " 20"; }
		/* Add more for additional items as needed */
		.controls-box {
			padding: 8px;
			margin-bottom: 8px;
			background-color: #f0f0f0;
			border-radius: 4px;
		}
	`)
	if err != nil {
		log.Printf("Failed to load CSS: %v", err)
	} else {
		// Apply CSS provider to the entire application
		gtk4.AddProviderForDisplay(cssProvider, 600)
	}

	// Define functions for common operations to be shared between buttons and menu items
	sayHello := func() {
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
		dialog.ConnectResponse(func(responseId gtk4.ResponseType) {
			fmt.Printf("Dialog response: %d\n", responseId)
			dialog.Destroy() // Destroy the dialog when done

			// Add log entry for the action
			logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Greeted %s", time.Now().Format("15:04:05"), name))
			logEntry.AddCssClass("log-entry")
			logBox.Prepend(logEntry)
		})

		// Show the dialog
		dialog.Show()

		resultLbl.SetText(fmt.Sprintf("Hello, %s!", name))
	}

	showAboutDialog := func() {
		// Create a custom about dialog
		dialog := gtk4.NewDialog("About This Application", win, gtk4.DialogModal|gtk4.DialogDestroyWithParent)

		// Get the content area of the dialog
		content := dialog.GetContentArea()

		// Add some content to the dialog
		titleLabel := gtk4.NewLabel("GTK4Go Demo Application")
		titleLabel.AddCssClass("title")
		descLabel := gtk4.NewLabel("This is a simple demonstration of GTK4 bindings for Go.")
		versionLabel := gtk4.NewLabel("Version: 1.0.0")
		featuresLabel := gtk4.NewLabel("New Features: ListView with data binding and selection")

		// Add widgets to the content area
		content.Append(titleLabel)
		content.Append(descLabel)
		content.Append(versionLabel)
		content.Append(featuresLabel)

		// Add padding to the content area
		content.SetSpacing(10)

		// Add OK button to the dialog
		dialog.AddButton("OK", gtk4.ResponseOk)

		// Connect response handler
		dialog.ConnectResponse(func(responseId gtk4.ResponseType) {
			fmt.Printf("About dialog response: %d\n", responseId)
			dialog.Destroy()

			// Add log entry for the action
			logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Opened About dialog", time.Now().Format("15:04:05")))
			logEntry.AddCssClass("log-entry")
			logBox.Prepend(logEntry)

			// Switch to logs tab
			rightStack.SetVisibleChildName("logs")
		})

		// Show the dialog
		dialog.Show()
	}

	showOpenFileDialog := func() {
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
		confirmDialog.ConnectResponse(func(responseId gtk4.ResponseType) {
			fmt.Printf("Confirm dialog response: %d\n", responseId)
			confirmed := (responseId == gtk4.ResponseYes)
			confirmDialog.Destroy()

			if confirmed {
				// Create file open dialog
				fileDialog := gtk4.NewFileDialog("Select a File", win, gtk4.FileDialogActionOpen)

				// Connect response handler for the file dialog
				fileDialog.ConnectResponse(func(responseId gtk4.ResponseType) {
					fmt.Printf("File dialog response: %d\n", responseId)
					if responseId == gtk4.ResponseAccept {
						filename := fileDialog.GetFilename()
						if filename != "" {
							// Update UI
							resultLbl.SetText(fmt.Sprintf("Selected file: %s", filename))

							// Log the selection
							fmt.Printf("Selected file: %s\n", filename)

							// Add log entry for the action
							logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Selected file: %s",
								time.Now().Format("15:04:05"), filename))
							logEntry.AddCssClass("log-entry")
							logBox.Prepend(logEntry)

							// Switch to logs tab
							rightStack.SetVisibleChildName("logs")
						}
					}
					fileDialog.Destroy()
				})

				// Show the file dialog
				fileDialog.Show()
			}
		})

		// Show the confirmation dialog
		confirmDialog.Show()
	}

	clearInput := func() {
		// Clear the entry field
		entry.SetText("")
		resultLbl.SetText("Hello, World!")
		
		// Add log entry for the action
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Cleared input", time.Now().Format("15:04:05")))
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)
	}

	showInfoTab := func() {
		rightStack.SetVisibleChildName("info")
	}

	showLogsTab := func() {
		rightStack.SetVisibleChildName("logs")
	}

	showHelpTab := func() {
		rightStack.SetVisibleChildName("help")
	}
	
	showListViewTab := func() {
		rightStack.SetVisibleChildName("listview")
	}

	runLongTask := func() {
		// Only implement if a task is not already running
		if cancelFunc != nil {
			// Cancel the current task
			cancelFunc()
			cancelFunc = nil
			longTaskBtn.SetLabel("Run Long Task")
			progressLbl.SetText("Task cancelled")
			return
		}

		// Update UI to show task is starting
		longTaskBtn.SetLabel("Cancel Task")
		longTaskBtn.AddCssClass("disabled")
		progressLbl.SetText("Starting task...")
		progressLbl.AddCssClass("progress-label")

		// Add log entry for starting the task
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Started long task", time.Now().Format("15:04:05")))
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)

		// Switch to logs tab
		rightStack.SetVisibleChildName("logs")

		// Start a background task
		cancelFunc = gtk4go.QueueBackgroundTask(
			"long-task",
			func(ctx context.Context, progress func(percent int, message string)) (interface{}, error) {
				// Simulate a long task with 10 steps
				for i := 0; i <= 100; i += 10 {
					// Check for cancellation
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					default:
						// Continue processing
					}

					// Update progress
					progress(i, fmt.Sprintf("Processing step %d of 10", i/10))

					// Add log entry for each step
					progressMsg := fmt.Sprintf("Task step %d of 10 completed", i/10)
					gtk4go.RunOnUIThread(func() {
						logStep := gtk4.NewLabel(fmt.Sprintf("[%s] %s",
							time.Now().Format("15:04:05"), progressMsg))
						logStep.AddCssClass("log-entry")
						logBox.Prepend(logStep)
					})

					// Simulate work
					time.Sleep(500 * time.Millisecond)
				}

				// Return some result data
				return "Task completed successfully!", nil
			},
			func(result interface{}, err error) {
				// Reset button
				longTaskBtn.SetLabel("Run Long Task")
				longTaskBtn.RemoveCssClass("disabled")

				// Update result based on success or failure
				if err != nil {
					if err == context.Canceled {
						progressLbl.SetText("Task was cancelled")

						// Add log entry for cancellation
						logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Task cancelled",
							time.Now().Format("15:04:05")))
						logEntry.AddCssClass("log-entry")
						logBox.Prepend(logEntry)
					} else {
						progressLbl.SetText(fmt.Sprintf("Error: %v", err))

						// Add log entry for error
						logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Task error: %v",
							time.Now().Format("15:04:05"), err))
						logEntry.AddCssClass("log-entry")
						logBox.Prepend(logEntry)
					}
				} else {
					progressLbl.SetText(result.(string))

					// Add log entry for completion
					logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] %s",
						time.Now().Format("15:04:05"), result.(string)))
					logEntry.AddCssClass("log-entry")
					logBox.Prepend(logEntry)
				}

				// Clear the cancel function
				cancelFunc = nil
			},
			func(percent int, message string) {
				// This runs on the UI thread to show progress
				progressLbl.SetText(fmt.Sprintf("%d%% - %s", percent, message))
			},
		)
	}

	exitApp := func() {
		// Exit the application
		os.Exit(0)
	}

	// Connect button click events to the functions
	helloBtn.ConnectClicked(sayHello)
	aboutBtn.ConnectClicked(showAboutDialog)
	fileBtn.ConnectClicked(showOpenFileDialog)
	longTaskBtn.ConnectClicked(runLongTask)

	// Get the application's action group
	actionGroup := app.GetActionGroup()

	// Create actions for the menu items
	sayHelloAction := gtk4.NewAction("say_hello", sayHello)
	actionGroup.AddAction(sayHelloAction)

	openAction := gtk4.NewAction("open", showOpenFileDialog)
	actionGroup.AddAction(openAction)

	saveAction := gtk4.NewAction("save", func() {
		// Implement save functionality (placeholder)
		fmt.Println("Save action triggered")
		resultLbl.SetText("Save action triggered (not implemented)")
		
		// Add log entry for the action
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Save action triggered (not implemented)", 
			time.Now().Format("15:04:05")))
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)
	})
	actionGroup.AddAction(saveAction)

	clearAction := gtk4.NewAction("clear", clearInput)
	actionGroup.AddAction(clearAction)

	aboutAction := gtk4.NewAction("about", showAboutDialog)
	actionGroup.AddAction(aboutAction)

	logsAction := gtk4.NewAction("logs", showLogsTab)
	actionGroup.AddAction(logsAction)

	infoAction := gtk4.NewAction("info", showInfoTab)
	actionGroup.AddAction(infoAction)

	helpAction := gtk4.NewAction("help", showHelpTab)
	actionGroup.AddAction(helpAction)
	
	listViewAction := gtk4.NewAction("listview", showListViewTab)
	actionGroup.AddAction(listViewAction)
	
	taskAction := gtk4.NewAction("task", runLongTask)
	actionGroup.AddAction(taskAction)

	exitAction := gtk4.NewAction("exit", exitApp)
	actionGroup.AddAction(exitAction)

	// Create application menu
	menu := gtk4.NewMenu()

	// Create File menu
	fileMenu := gtk4.NewMenu()
	fileOpenItem := gtk4.NewMenuItem("Open", "app.open")
	fileSaveItem := gtk4.NewMenuItem("Save", "app.save")
	fileExitItem := gtk4.NewMenuItem("Exit", "app.exit")
	fileMenu.AppendItem(fileOpenItem)
	fileMenu.AppendItem(fileSaveItem)
	fileMenu.AppendItem(fileExitItem)
	menu.AppendSubmenu("File", fileMenu)

	// Create Edit menu
	editMenu := gtk4.NewMenu()
	editClearItem := gtk4.NewMenuItem("Clear", "app.clear")
	editMenu.AppendItem(editClearItem)
	menu.AppendSubmenu("Edit", editMenu)

	// Create View menu
	viewMenu := gtk4.NewMenu()
	viewLogsItem := gtk4.NewMenuItem("Show Logs", "app.logs")
	viewInfoItem := gtk4.NewMenuItem("Show Info", "app.info")
	viewHelpItem := gtk4.NewMenuItem("Show Help", "app.help")
	viewListViewItem := gtk4.NewMenuItem("Show ListView", "app.listview")
	viewMenu.AppendItem(viewLogsItem)
	viewMenu.AppendItem(viewInfoItem)
	viewMenu.AppendItem(viewHelpItem)
	viewMenu.AppendItem(viewListViewItem)
	menu.AppendSubmenu("View", viewMenu)
	
	// Create Tools menu
	toolsMenu := gtk4.NewMenu()
	toolsTaskItem := gtk4.NewMenuItem("Run Task", "app.task")
	toolsMenu.AppendItem(toolsTaskItem)
	menu.AppendSubmenu("Tools", toolsMenu)

	// Create Help menu
	helpMenu := gtk4.NewMenu()
	helpAboutItem := gtk4.NewMenuItem("About", "app.about")
	helpMenu.AppendItem(helpAboutItem)
	menu.AppendSubmenu("Help", helpMenu)

	// Set the menubar's menu model
	menuBar.SetMenuModel(menu)

	// Create a menu model for the menu button
	quickMenu := gtk4.NewMenu()
	quickHelloItem := gtk4.NewMenuItem("Say Hello", "app.say_hello")
	quickOpenItem := gtk4.NewMenuItem("Open File", "app.open")
	quickListViewItem := gtk4.NewMenuItem("Show ListView", "app.listview")
	quickAboutItem := gtk4.NewMenuItem("About", "app.about")
	quickExitItem := gtk4.NewMenuItem("Exit", "app.exit")
	
	quickMenu.AppendItem(quickHelloItem)
	quickMenu.AppendItem(quickOpenItem)
	quickMenu.AppendItem(quickListViewItem)
	quickMenu.AppendItem(quickAboutItem)
	quickMenu.AppendItem(quickExitItem)
	
	// Create a popover menu for the menu button and connect it
	popoverMenu := gtk4.NewPopoverMenu(quickMenu)
	menuButton.SetPopover(popoverMenu)

	// Add the main box to the window
	win.SetChild(mainBox)

	// Add the window to the application
	app.AddWindow(win)

	// Show instructions on how to test the window performance
	log.Println("Running Hello World with menus and optimized window performance.")
	log.Println("Try using the menu bar and menu button to access application features.")
	log.Println("Check out the new ListView tab to see ListView widget in action.")

	// Run the application
	os.Exit(app.Run())

	// Clean up background workers at exit
	gtk4go.DefaultWorker.Stop()
}

// Variable for long task cancellation
var cancelFunc context.CancelFunc