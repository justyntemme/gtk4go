# GTK4Go Widgets Guide

This document provides a practical guide to using the widgets available in GTK4Go, with examples drawn from the included hello-world application.

## Table of Contents

- [Application](#application)
- [Window](#window)
- [Box](#box)
- [Label](#label)
- [Button](#button)
- [Entry](#entry)
- [Grid](#grid)
- [Paned](#paned)
- [Stack and StackSwitcher](#stack-and-stackswitcher)
- [ScrolledWindow](#scrolledwindow)
- [ListView and Models](#listview-and-models)
- [Dialog](#dialog)
- [Menu Components](#menu-components)
- [CSS Styling](#css-styling)
- [Background Tasks](#background-tasks)

## Application

The `Application` widget is the entry point for GTK applications. It handles application-wide functionality such as command-line arguments, application ID, and window management.

```go
// Create a new application
app := gtk4.NewApplication("com.example.HelloWorld")

// Add a window to the application
app.AddWindow(win)

// Run the application and get exit code
os.Exit(app.Run())
```

The application ID should be a unique, reverse-domain name for your application.

## Window

The `Window` widget is the main container for your application's user interface.

```go
// Create a window with a title
win := gtk4.NewWindow("Hello GTK4 from Go!")

// Set default size
win.SetDefaultSize(1300, 950)

// Enable hardware-accelerated rendering
win.EnableAcceleratedRendering()

// Optimize for resizing
win.OptimizeForResizing()

// Set up CSS optimization during window resize
win.SetupCSSOptimizedResize()

// Set the main child widget
win.SetChild(mainBox)

// Connect close request signal
win.ConnectCloseRequest(func() bool {
    // Return true to prevent closing, false to allow
    return false
})

// Show the window
win.Show()
```

The `Window` widget includes performance optimizations for smooth rendering and resizing.

## Box

The `Box` widget arranges child widgets in a horizontal or vertical line.

```go
// Create a vertical box with 10px spacing
mainBox := gtk4.NewBox(gtk4.OrientationVertical, 10)

// Create a horizontal box with 6px spacing
listViewControls := gtk4.NewBox(gtk4.OrientationHorizontal, 6)

// Add a CSS class
listViewControls.AddCssClass("controls-box")

// Add a child widget
mainBox.Append(menuBar)
mainBox.Append(titleLabel)
```

Box widgets can be nested to create complex layouts, and they automatically manage the size of their children.

## Label

The `Label` widget displays text.

```go
// Create a simple label
titleLabel := gtk4.NewLabel("GTK4Go Demo Application")

// Add a CSS class for styling
titleLabel.AddCssClass("title")

// Create a label with markup
helpText := gtk4.NewLabel(`
Using this application:

1. Enter your name in the text field
2. Click "Say Hello" to see a greeting
...
`)
```

Labels support basic formatting and can display multi-line text.

## Button

The `Button` widget is a clickable control that triggers an action.

```go
// Create a button with a label
helloBtn := gtk4.NewButton("Say Hello")

// Add CSS class for styling
helloBtn.AddCssClass("square-button")

// Connect clicked signal
helloBtn.ConnectClicked(func() {
    // Code to run when button is clicked
    name := entry.GetText()
    resultLbl.SetText(fmt.Sprintf("Hello, %s!", name))
})
```

Buttons can have text labels or icons, and they emit the "clicked" signal when activated.

## Entry

The `Entry` widget is a single-line text input field.

```go
// Create a new entry
entry := gtk4.NewEntry()

// Set placeholder text
entry.SetPlaceholderText("Type your name here")

// Get entered text
name := entry.GetText()

// Set text programmatically
entry.SetText("")

// Connect activate signal (triggered when Enter is pressed)
entry.ConnectActivate(func() {
    sayHello()
})

// Connect changed signal (triggered when text changes)
entry.ConnectChanged(func() {
    // React to text changes
})
```

The Entry widget is used for user text input and can be configured with various input modes and validation.

## Grid

The `Grid` widget arranges child widgets in a table-like layout.

```go
// Create a grid with options
buttonsGrid := gtk4.NewGrid(
    gtk4.WithRowSpacing(10),
    gtk4.WithColumnSpacing(10),
    gtk4.WithColumnHomogeneous(true),
)

// Add widgets to specific positions (column, row, width, height)
buttonsGrid.Attach(helloBtn, 0, 0, 1, 1)
buttonsGrid.Attach(aboutBtn, 1, 0, 1, 1)
buttonsGrid.Attach(fileBtn, 0, 1, 1, 1)
buttonsGrid.Attach(longTaskBtn, 1, 1, 1, 1)
```

Grid is useful for creating form layouts and other structured arrangements of widgets.

## Paned

The `Paned` widget contains two child widgets with an adjustable divider between them.

```go
// Create a horizontal paned container
paned := gtk4.NewPaned(gtk4.OrientationHorizontal,
    gtk4.WithPosition(350),  // Initial position of divider
    gtk4.WithWideHandle(true),  // Use a wider handle for easier grabbing
)

// Set the start (left) and end (right) children
paned.SetStartChild(leftBox)
paned.SetEndChild(rightBox)
```

Paned containers are perfect for creating resizable split views.

## Stack and StackSwitcher

The `Stack` widget shows one child at a time, with animated transitions between them. The `StackSwitcher` provides buttons to switch between stack pages.

```go
// Create a stack with transition animation
rightStack := gtk4.NewStack(
    gtk4.WithTransitionType(gtk4.StackTransitionTypeSlideLeftRight),
    gtk4.WithTransitionDuration(200),
)

// Add titled pages to the stack
rightStack.AddTitled(infoBox, "info", "Information")
rightStack.AddTitled(scrollWin, "logs", "Logs")
rightStack.AddTitled(helpBox, "help", "Help")
rightStack.AddTitled(listViewBox, "listview", "ListView")

// Switch to a specific page
rightStack.SetVisibleChildName("logs")

// Create a stack switcher for the stack
stackSwitcher := gtk4.NewStackSwitcher(rightStack)

// Add the stack switcher and stack to a layout
rightBox.Append(stackSwitcher)
rightBox.Append(rightStack)
```

Stack and StackSwitcher work together to provide a tabbed interface.

## ScrolledWindow

The `ScrolledWindow` widget adds scrollbars around another widget when its content exceeds the visible area.

```go
// Create a scrolled window
scrollWin := gtk4.NewScrolledWindow(
    gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
    gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAlways),
    gtk4.WithPropagateNaturalHeight(false), // Don't propagate natural height to allow scrolling
)

// Set the child widget that will be scrollable
scrollWin.SetChild(logBox)
```

ScrolledWindow is essential for handling content that may not fit in the available space.

## ListView and Models

The `ListView` widget displays a list of items using a data model and a factory to create item widgets.

```go
// Create a string list model with data
listModel := gtk4.NewStringList()
for i := 1; i <= 15; i++ {
    listModel.Append(fmt.Sprintf("List Item %d", i))
}

// Create a selection model to handle item selection
selectionModel := gtk4.NewSingleSelection(listModel,
    gtk4.WithAutoselect(false),
    gtk4.WithInitialSelection(0),
)

// Create a factory for list items
factory := gtk4.NewSignalListItemFactory()

// Set up list items with setup callback
factory.ConnectSetup(func(listItem *gtk4.ListItem) {
    // Create a box for layout
    box := gtk4.NewBox(gtk4.OrientationHorizontal, 10)
    box.SetHExpand(true)
    box.AddCssClass("list-item-box")

    // Create an icon and label
    icon := gtk4.NewLabel("•")
    icon.AddCssClass("list-item-icon")
    box.Append(icon)

    label := gtk4.NewLabel("")
    label.AddCssClass("list-item-label")
    box.Append(label)

    // Set the box as the child of the list item
    listItem.SetChild(box)
})

// Bind data to list items
factory.ConnectBind(func(listItem *gtk4.ListItem) {
    // Get the text from the model
    text := listItem.GetText()
    if text == "" {
        text = fmt.Sprintf("Item %d", listItem.GetPosition()+1)
    }

    // Set the text on the label inside the box
    listItem.SetTextOnChildLabel(text)

    // Add selected class if the item is selected
    if listItem.GetSelected() {
        boxWidget := listItem.GetChild()
        boxWidget.AddCssClass("selected")
    } else {
        boxWidget := listItem.GetChild()
        boxWidget.RemoveCssClass("selected")
    }
})

// Create the list view with selection model and factory
listView := gtk4.NewListView(selectionModel, factory,
    gtk4.WithShowSeparators(true),
    gtk4.WithSingleClickActivate(true),
)

// Connect activate signal
listView.ConnectActivate(func(position int) {
    // Do something when an item is activated
})

// Add items to the model
listModel.Append("New Item")

// Remove items from the model
listModel.Remove(position)

// Get the number of items
count := listModel.GetNItems()

// Get a specific item
item := listModel.GetString(position)

// Set a selected item
selectionModel.SetSelected(position)

// Get selected item
selectedPos := selectionModel.GetSelected()
```

ListView is a modern, flexible list widget that separates data from presentation.

## Dialog

GTK4Go provides several dialog types for common interactions.

### MessageDialog

```go
// Create a message dialog
messageDialog := gtk4.NewMessageDialog(
    win,                   // Parent window
    gtk4.DialogModal,      // Flags
    gtk4.MessageInfo,      // Message type
    gtk4.ResponseOk,       // Buttons
    "You selected item at position 10",  // Message
)
messageDialog.SetTitle("ListView Item Selected")

// Connect response handler
messageDialog.ConnectResponse(func(responseId gtk4.ResponseType) {
    messageDialog.Destroy()
})

// Show the dialog
messageDialog.Show()
```

### Custom Dialog

```go
// Create a custom dialog
dialog := gtk4.NewDialog("About This Application", win, 
    gtk4.DialogModal|gtk4.DialogDestroyWithParent)

// Get the content area
content := dialog.GetContentArea()

// Add widgets to the content area
titleLabel := gtk4.NewLabel("GTK4Go Demo Application")
content.Append(titleLabel)

// Add buttons
dialog.AddButton("OK", gtk4.ResponseOk)

// Connect response signal
dialog.ConnectResponse(func(responseId gtk4.ResponseType) {
    dialog.Destroy()
})

// Show the dialog
dialog.Show()
```

### FileDialog

```go
// Create a file dialog
fileDialog := gtk4.NewFileDialog("Select a File", win, gtk4.FileDialogActionOpen)

// Connect response handler
fileDialog.ConnectResponse(func(responseId gtk4.ResponseType) {
    if responseId == gtk4.ResponseAccept {
        filename := fileDialog.GetFilename()
        // Use the selected filename
    }
    fileDialog.Destroy()
})

// Show the dialog
fileDialog.Show()
```

Dialogs are modal windows that request input from the user or display information.

## Menu Components

GTK4Go provides several components for creating menus.

### Menu and MenuItems

```go
// Create a menu
menu := gtk4.NewMenu()

// Create menu items with actions
fileOpenItem := gtk4.NewMenuItem("Open", "app.open")
fileSaveItem := gtk4.NewMenuItem("Save", "app.save")

// Add items to menu
menu.AppendItem(fileOpenItem)
menu.AppendItem(fileSaveItem)

// Create submenu
fileMenu := gtk4.NewMenu()
menu.AppendSubmenu("File", fileMenu)
```

### MenuBar

```go
// Create a menu bar
menuBar := gtk4.NewMenuBar()

// Set the menu model
menuBar.SetMenuModel(menu)

// Add to layout
mainBox.Append(menuBar)
```

### MenuButton and PopoverMenu

```go
// Create a menu button
menuButton := gtk4.NewMenuButton()
menuButton.SetLabel("Quick Actions")

// Create a menu for the button
quickMenu := gtk4.NewMenu()
quickHelloItem := gtk4.NewMenuItem("Say Hello", "app.say_hello")
quickMenu.AppendItem(quickHelloItem)

// Create a popover menu and connect to button
popoverMenu := gtk4.NewPopoverMenu(quickMenu)
menuButton.SetPopover(popoverMenu)
```

### Actions

```go
// Get action group
actionGroup := app.GetActionGroup()

// Create an action
sayHelloAction := gtk4.NewAction("say_hello", sayHello)
actionGroup.AddAction(sayHelloAction)
```

Menu components allow you to create application menus, context menus, and more.

## CSS Styling

GTK4Go supports CSS styling for widgets.

```go
// Load CSS from string
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
`)
if err != nil {
    log.Printf("Failed to load CSS: %v", err)
} else {
    // Apply CSS provider globally
    gtk4.AddProviderForDisplay(cssProvider, 600)
}

// Add CSS classes to widgets
titleLabel.AddCssClass("title")
button.AddCssClass("square-button")

// Remove CSS classes
button.RemoveCssClass("disabled")

// Check if a widget has a CSS class
if button.HasCssClass("disabled") {
    // Widget has the class
}
```

CSS styling allows you to customize the appearance of your application's widgets.

## Background Tasks

GTK4Go provides a background task system for running operations without blocking the UI.

```go
// Queue a background task
cancelFunc = gtk4go.QueueBackgroundTask(
    "task-id",  // Task ID
    func(ctx context.Context, progress func(percent int, message string)) (interface{}, error) {
        // Task code runs in background
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
            
            // Simulate work
            time.Sleep(500 * time.Millisecond)
        }
        
        // Return result
        return "Task completed successfully!", nil
    },
    func(result interface{}, err error) {
        // Completion callback runs on UI thread
        if err != nil {
            // Handle error
        } else {
            // Use result
            resultStr := result.(string)
        }
    },
    func(percent int, message string) {
        // Progress callback runs on UI thread
        progressLbl.SetText(fmt.Sprintf("%d%% - %s", percent, message))
    },
)

// Cancel a running task
if cancelFunc != nil {
    cancelFunc()
    cancelFunc = nil
}
```

Background tasks allow you to run long operations without freezing the UI, with progress updates and proper cancellation support.

## Best Practices

1. **Use builder pattern with options**: Most widgets support a functional options pattern for configuration.

   ```go
   grid := gtk4.NewGrid(
       gtk4.WithRowSpacing(10),
       gtk4.WithColumnSpacing(10),
       gtk4.WithColumnHomogeneous(true),
   )
   ```

2. **Organize complex UIs with containers**: Use nested containers like Box, Grid, and Paned to create well-structured layouts.

3. **Apply CSS classes for styling**: Use CSS classes instead of inline styles for better maintainability.

4. **Properly disconnect signals**: When destroying widgets, ensure signals are disconnected to prevent memory leaks.

5. **Run long operations in background**: Use the background task system for operations that might block the UI.

6. **Optimize windows for performance**: Enable hardware acceleration and resize optimization for smooth UIs.

   ```go
   win.EnableAcceleratedRendering()
   win.OptimizeForResizing()
   win.SetupCSSOptimizedResize()
   ```

7. **Use ListModels for data-driven UIs**: Separate your data from presentation using the modern ListView widget.