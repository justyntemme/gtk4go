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
	win.SetDefaultSize(800, 600)

	// Enable hardware-accelerated rendering
	win.EnableAcceleratedRendering()

	// Set up CSS optimization during window resize
	win.SetupCSSOptimizedResize()

	// Optimize for resizing specifically
	win.OptimizeForResizing()

	// Create a vertical box container as the main layout
	mainBox := gtk4.NewBox(gtk4.OrientationVertical, 10)

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
	widgets := []string{"Grid", "Paned", "Stack", "StackSwitcher", "ScrolledWindow", "TreeView"}
	descriptions := []string{
		"Arranges widgets in rows and columns",
		"Divides space between two widgets with adjustable separator",
		"Shows one widget at a time with transitions",
		"Provides buttons to switch between stack pages",
		"Provides scrolling for large content",
		"Displays data in list or tree format",
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
6. Try the TreeView tab to see list and tree views

This demo showcases GTK4Go's layout containers and widgets.
	`)

	helpBox.Append(helpText)
	rightStack.AddTitled(helpBox, "help", "Help")

	// Stack Page 4: TreeView Demo
	treeViewBox := gtk4.NewBox(gtk4.OrientationVertical, 10)
	treeViewBox.Append(gtk4.NewLabel("TreeView Demo"))
	
	// Create a paned container to show two types of TreeViews
	treeViewPaned := gtk4.NewPaned(gtk4.OrientationHorizontal, 
		gtk4.WithPosition(380),
		gtk4.WithWideHandle(true),
	)
	
	// LEFT SIDE: List TreeView
	listBox := gtk4.NewBox(gtk4.OrientationVertical, 5)
	listLabel := gtk4.NewLabel("Simple List View:")
	listLabel.AddCssClass("tree-header")
	listBox.Append(listLabel)
	
	// Create a scrolled window for the list view
	listScrollWin := gtk4.NewScrolledWindow(
		gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
		gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
	)
	
	// Create a list store model for the list view
	listStore := gtk4.NewListStore(
		gtk4.G_TYPE_STRING, // Name
		gtk4.G_TYPE_INT,    // Value
		gtk4.G_TYPE_BOOLEAN, // Active
	)
	
	// Add data to list store
	for i := 0; i < 10; i++ {
		iter := listStore.Append()
		listStore.SetValue(iter, 0, fmt.Sprintf("Item %d", i+1))
		listStore.SetValue(iter, 1, (i+1)*10)
		listStore.SetValue(iter, 2, i%2 == 0)
	}
	
	// Create the list view
	listView := gtk4.NewTreeView(listStore, gtk4.WithHeaders(true))
	
	// Create renderers
	textRenderer := gtk4.NewCellRendererText()
	valueRenderer := gtk4.NewCellRendererText()
	toggleRenderer := gtk4.NewCellRendererToggle()
	
	// Make toggle renderer interactive
	toggleRenderer.SetActive(true) // Default state
	toggleRenderer.ConnectToggled(func(path string) {
		treePath := gtk4.NewTreePathFromString(path)
		iter, ok := listStore.GetIter(treePath)
		if !ok {
			return
		}
		
		// Get current value
		value, _ := listStore.GetValue(iter, 2)
		active := value.(bool)
		
		// Toggle value
		listStore.SetValue(iter, 2, !active)
		
		// Add log entry for the action
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Toggled item %s to %v", 
			time.Now().Format("15:04:05"), path, !active))
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)
	})
	
	// Create columns
	nameColumn := gtk4.NewTreeViewColumn("Name", textRenderer, gtk4.Attr("text", 0))
	valueColumn := gtk4.NewTreeViewColumn("Value", valueRenderer, gtk4.Attr("text", 1))
	activeColumn := gtk4.NewTreeViewColumn("Active", toggleRenderer, gtk4.Attr("active", 2))
	
	// Make columns resizable
	nameColumn.SetResizable(true)
	valueColumn.SetResizable(true)
	activeColumn.SetResizable(true)
	
	// Add columns to the tree view
	listView.AppendColumn(nameColumn)
	listView.AppendColumn(valueColumn)
	listView.AppendColumn(activeColumn)
	
	// Create a status label for selection feedback
	listStatusLabel := gtk4.NewLabel("No selection")
	listStatusLabel.AddCssClass("status-label")
	
	// Set up selection handler
	listSelection := listView.GetSelection()
	listSelection.ConnectChanged(func() {
		model, iter, ok := listSelection.GetSelected()
		if ok {
			nameVal, _ := model.GetValue(iter, 0)
			valueVal, _ := model.GetValue(iter, 1)
			activeVal, _ := model.GetValue(iter, 2)
			
			statusText := fmt.Sprintf("Selected: %s, Value: %d, Active: %v",
				nameVal.(string), valueVal.(int), activeVal.(bool))
			listStatusLabel.SetText(statusText)
			
			// Add to log
			logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] %s", 
				time.Now().Format("15:04:05"), statusText))
			logEntry.AddCssClass("log-entry")
			logBox.Prepend(logEntry)
		} else {
			listStatusLabel.SetText("No selection")
		}
	})
	
	// Add list view to scrolled window
	listScrollWin.SetChild(listView)
	
	// Create button toolbar for list actions
	listButtons := gtk4.NewBox(gtk4.OrientationHorizontal, 5)
	
	addButton := gtk4.NewButton("Add Item")
	addButton.ConnectClicked(func() {
		iter := listStore.Append()
		
		// Count existing items
		count := 0
		currentIter, hasIter := listStore.GetIter(gtk4.NewTreePath())
		for hasIter {
			count++
			hasIter = listStore.IterNext(currentIter)
		}
		
		// Add new item
		newName := fmt.Sprintf("New Item %d", count)
		listStore.SetValue(iter, 0, newName)
		listStore.SetValue(iter, 1, count*5)
		listStore.SetValue(iter, 2, true)
		
		// Add log entry
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Added new item: %s", 
			time.Now().Format("15:04:05"), newName))
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)
	})
	
	listButtons.Append(addButton)
	
	// Add list view components to the list box
	listBox.Append(listScrollWin)
	listBox.Append(listStatusLabel)
	listBox.Append(listButtons)
	
	// RIGHT SIDE: Tree TreeView
	treeBox := gtk4.NewBox(gtk4.OrientationVertical, 5)
	treeLabel := gtk4.NewLabel("Hierarchical Tree View:")
	treeLabel.AddCssClass("tree-header")
	treeBox.Append(treeLabel)
	
	// Create a scrolled window for the tree view
	treeScrollWin := gtk4.NewScrolledWindow(
		gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
		gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
	)
	
	// Create a tree store model
	treeStore := gtk4.NewTreeStore(
		gtk4.G_TYPE_STRING, // Category/Item
		gtk4.G_TYPE_INT,    // Count
	)
	
	// Add data to tree store
	categories := []string{"Electronics", "Books", "Clothing", "Food"}
	items := [][]string{
		{"Laptop", "Phone", "Tablet", "Camera"},
		{"Fiction", "Non-fiction", "Reference", "Magazines"},
		{"Shirts", "Pants", "Shoes", "Accessories"},
		{"Fruits", "Vegetables", "Meat", "Dairy"},
	}
	counts := [][]int{
		{5, 10, 7, 3},
		{12, 8, 5, 3},
		{20, 15, 10, 5},
		{8, 12, 6, 9},
	}
	
	for i, category := range categories {
		// Add category (parent node)
		parentIter := treeStore.Append(nil)
		treeStore.SetValue(parentIter, 0, category)
		
		// Calculate total count
		totalCount := 0
		for _, count := range counts[i] {
			totalCount += count
		}
		treeStore.SetValue(parentIter, 1, totalCount)
		
		// Add items (child nodes)
		for j, item := range items[i] {
			childIter := treeStore.Append(parentIter)
			treeStore.SetValue(childIter, 0, item)
			treeStore.SetValue(childIter, 1, counts[i][j])
		}
	}
	
	// Create the tree view
	treeView := gtk4.NewTreeView(treeStore, gtk4.WithHeaders(true))
	
	// Create renderers
	treeCatRenderer := gtk4.NewCellRendererText()
	treeCountRenderer := gtk4.NewCellRendererText()
	
	// Create columns
	catColumn := gtk4.NewTreeViewColumn("Category/Item", treeCatRenderer, gtk4.Attr("text", 0))
	countColumn := gtk4.NewTreeViewColumn("Count", treeCountRenderer, gtk4.Attr("text", 1))
	
	// Configure columns
	catColumn.SetResizable(true)
	countColumn.SetResizable(true)
	
	// Add columns to the tree view
	treeView.AppendColumn(catColumn)
	treeView.AppendColumn(countColumn)
	
	// Create a status label for tree selection
	treeStatusLabel := gtk4.NewLabel("No selection")
	treeStatusLabel.AddCssClass("status-label")
	
	// Set up selection handler
	treeSelection := treeView.GetSelection()
	treeSelection.ConnectChanged(func() {
		model, iter, ok := treeSelection.GetSelected()
		if ok {
			nameVal, _ := model.GetValue(iter, 0)
			countVal, _ := model.GetValue(iter, 1)
			
			statusText := fmt.Sprintf("Selected: %s, Count: %d",
				nameVal.(string), countVal.(int))
			treeStatusLabel.SetText(statusText)
			
			// Add to log
			logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] %s", 
				time.Now().Format("15:04:05"), statusText))
			logEntry.AddCssClass("log-entry")
			logBox.Prepend(logEntry)
		} else {
			treeStatusLabel.SetText("No selection")
		}
	})
	
	// Add tree view to scrolled window
	treeScrollWin.SetChild(treeView)
	
	// Create button toolbar for tree actions
	treeButtons := gtk4.NewBox(gtk4.OrientationHorizontal, 5)
	
	expandButton := gtk4.NewButton("Expand All")
	expandButton.ConnectClicked(func() {
		treeView.ExpandAll()
		
		// Add log entry
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Expanded all tree nodes", 
			time.Now().Format("15:04:05")))
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)
	})
	
	collapseButton := gtk4.NewButton("Collapse All")
	collapseButton.ConnectClicked(func() {
		treeView.CollapseAll()
		
		// Add log entry
		logEntry := gtk4.NewLabel(fmt.Sprintf("[%s] Collapsed all tree nodes", 
			time.Now().Format("15:04:05")))
		logEntry.AddCssClass("log-entry")
		logBox.Prepend(logEntry)
	})
	
	treeButtons.Append(expandButton)
	treeButtons.Append(collapseButton)
	
	// Add tree view components to the tree box
	treeBox.Append(treeScrollWin)
	treeBox.Append(treeStatusLabel)
	treeBox.Append(treeButtons)
	
	// Add list and tree views to paned container
	treeViewPaned.SetStartChild(listBox)
	treeViewPaned.SetEndChild(treeBox)
	
	// Add paned container to the tree view box
	treeViewBox.Append(treeViewPaned)
	
	// Add TreeView tab to right stack
	rightStack.AddTitled(treeViewBox, "treeview", "TreeView")

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
		.tree-header {
			font-weight: bold;
			color: #333333;
			margin-top: 5px;
		}
		.status-label {
			font-style: italic;
			color: #666666;
			padding: 5px;
			border-top: 1px solid #e0e0e0;
		}
	`)
	if err != nil {
		log.Printf("Failed to load CSS: %v", err)
	} else {
		// Apply CSS provider to the entire application
		gtk4.AddProviderForDisplay(cssProvider, 600)
	}

	// Set up event handlers

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
	})

	// Connect long task button click event
	var cancelFunc context.CancelFunc

	longTaskBtn.ConnectClicked(func() {
		// Check if a task is already running
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
				// This runs in a background goroutine

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
				// This runs on the UI thread when task is completed or fails

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
	})

	// Add the main box to the window
	win.SetChild(mainBox)

	// Add the window to the application
	app.AddWindow(win)

	// Show instructions on how to test the window performance
	log.Println("Running Hello World with optimized window performance.")
	log.Println("Try resizing and moving the window to test the improvements.")

	// Run the application
	os.Exit(app.Run())

	// Clean up background workers at exit
	gtk4go.DefaultWorker.Stop()
}