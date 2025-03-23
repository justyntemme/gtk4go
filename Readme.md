# GTK4Go

GTK4Go is a Go binding library for GTK4 (GIMP Toolkit version 4). It provides a lightweight, idiomatic Go API for creating GUI applications using GTK4.

## Architecture Overview

GTK4Go is designed as a thin wrapper around the GTK4 C API using cgo. The library follows a component-based architecture, where each GTK widget is represented by a corresponding Go struct that encapsulates the underlying C implementation.

### Core Components

1. **Initialization System**
   - The main package (`gtk4go`) handles GTK4 initialization.
   - Automatic initialization occurs on import.
   - Manual initialization is also available via `Initialize()`.

2. **Widget System**
   - Each GTK widget is represented by a Go struct in the `gtk4` package.
   - All widgets implement a common interface with methods like `GetWidget()` to access the underlying C pointers.
   - Memory management is handled through Go's finalizers to ensure proper resource cleanup.

3. **Event Handling**
   - Signal connections are implemented using callback functions.
   - Callback management is done through handler maps with mutex protection.
   - The event loop is managed by GTK's application system rather than the deprecated `gtk_main()`.

4. **Component Hierarchy**
   ```
   GtkApplication
   └── GtkWindow
       └── Container Widgets (e.g., GtkBox)
           └── Child Widgets (e.g., GtkButton, GtkLabel)
   ```

### Current Implementation

The library currently implements the following components:

- **Application**: Core application functionality
- **Window**: Basic window management
- **Box**: Container for horizontal and vertical layouts
- **Button**: Clickable button with event handling
- **Label**: Text display widget

### Design Principles

1. **Idiomatic Go**: The API is designed to feel natural to Go developers while maintaining access to GTK's power.
2. **Safety**: Memory management and resource cleanup are automated where possible.
3. **Performance**: The binding is designed to be lightweight with minimal overhead.
4. **Simplicity**: The API aims to simplify GTK development without hiding its capabilities.

## Usage Example

```go
package main

import (
	"fmt"
	"log"
	"os"
	"github.com/justyntemme/gtk4go"
	"github.com/justyntemme/gtk4go/gtk4"
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
		fmt.Println("Button clicked!")
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
```

## Implementation Plan

### Step 1: Entry Widget Implementation

The Entry widget should be implemented next because:
- It provides essential text input capabilities required by most GUI applications
- It's relatively simple to implement but adds significant functionality
- It complements the existing components (Window, Box, Button, Label)

**Implementation Details:**
1. Create `gtk4/entry.go` file with the following structure:
   - C bindings for GtkEntry functions
   - Entry struct with widget pointer
   - Constructor functions (NewEntry, NewEntryWithBuffer)
   - Text manipulation methods (SetText, GetText, etc.)
   - Signal handling for "changed" and "activate" events
   - Property methods (SetPlaceholder, SetEditable, etc.)

**Implementation Considerations:**
- Implement proper memory management for text buffers
- Ensure signal callbacks are properly handled
- Consider implementing EntryBuffer as a separate type

**Example API:**
```go
// NewEntry creates a new text entry widget
func NewEntry() *Entry

// SetText sets the entry text
func (e *Entry) SetText(text string)

// GetText gets the entry text
func (e *Entry) GetText() string

// ConnectChanged connects a callback to text-changed event
func (e *Entry) ConnectChanged(callback func())

// SetPlaceholderText sets placeholder text
func (e *Entry) SetPlaceholderText(text string)
```

### Step 2: StyleContext and CSS Support

CSS styling should be implemented next because:
- It provides a way to customize the appearance of all widgets
- It's a core feature of modern GTK4 applications
- It will be required by more complex widgets later

**Implementation Details:**
1. Create `gtk4/css.go` file with:
   - CSS Provider implementation
   - StyleContext methods
   - Functions to load and apply CSS
   - Widget style property methods

2. Add style-related methods to existing widgets

**Implementation Considerations:**
- Ensure CSS files can be loaded from strings and files
- Implement proper error handling for CSS parsing
- Consider the scoping of CSS (application-wide vs. widget-specific)

**Example API:**
```go
// LoadCSSFromFile loads CSS from a file
func LoadCSSFromFile(path string) (*CSSProvider, error)

// LoadCSSFromString loads CSS from a string
func LoadCSSFromString(css string) (*CSSProvider, error)

// AddStyleClass adds a CSS class to a widget
func (w *Widget) AddStyleClass(className string)

// RemoveStyleClass removes a CSS class from a widget
func (w *Widget) RemoveStyleClass(className string)
```

### Step 3: Dialog Implementation

Dialogs should be implemented next because:
- They provide critical functionality for user interaction
- They're required for many common operations (file selection, alerts, etc.)
- They build upon the basic window system already implemented

**Implementation Details:**
1. Create `gtk4/dialog.go` with:
   - Base Dialog implementation
   - Common dialog types (MessageDialog, FileChooserDialog)
   - Response handling system

2. Implement specific dialog variants:
   - AlertDialog for simple notifications
   - ConfirmDialog for yes/no questions
   - FileChooserDialog for file operations

**Implementation Considerations:**
- Handle modal and non-modal behavior
- Implement proper signal handling for responses
- Consider higher-level convenience functions

**Example API:**
```go
// NewDialog creates a new dialog
func NewDialog(title string, parent *Window, flags DialogFlags) *Dialog

// AddButton adds a button to the dialog
func (d *Dialog) AddButton(text string, responseId ResponseType)

// Run runs the dialog modally
func (d *Dialog) Run() ResponseType

// ShowMessageDialog shows a simple message dialog
func ShowMessageDialog(parent *Window, messageType MessageType, 
                      title, message string) ResponseType

// ShowFileChooserDialog shows a file chooser dialog
func ShowFileChooserDialog(parent *Window, title string, 
                         action FileChooserAction) (string, bool)
```

### Step 4: Grid and Layout Containers

Grid and advanced layout containers should be implemented next because:
- They provide essential layout capabilities beyond simple boxes
- They're required for more complex UI designs
- They build upon the container system already in place

**Implementation Details:**
1. Create `gtk4/grid.go` with:
   - Grid container implementation
   - Cell placement and spanning methods
   - Row/column properties and alignment

2. Implement additional containers:
   - Paned (split view) container
   - Stack container (for tab-like interfaces)
   - Viewport for scrollable content

**Implementation Considerations:**
- Implement proper child positioning and sizing
- Consider alignment and expansion properties
- Add comprehensive layout methods

**Example API:**
```go
// NewGrid creates a new grid container
func NewGrid() *Grid

// Attach attaches a widget to a grid cell
func (g *Grid) Attach(child Widget, left, top, width, height int)

// SetRowSpacing sets spacing between rows
func (g *Grid) SetRowSpacing(spacing int)

// SetColumnHomogeneous sets whether columns are homogeneous
func (g *Grid) SetColumnHomogeneous(homogeneous bool)
```

### Step 5: ListView and TreeView Implementation

ListView and TreeView should be implemented next because:
- They're essential for displaying collections of data
- They're complex widgets that build upon previous implementations
- They require model/view architecture that can be reused

**Implementation Details:**
1. Create base model interfaces:
   - ListModel for simple lists
   - TreeModel for hierarchical data

2. Implement view widgets:
   - ListView for flat lists of items
   - TreeView for hierarchical data
   - Cell renderer system for customizing appearance

**Implementation Considerations:**
- Implement proper model/view separation
- Consider memory management for large datasets
- Implement selection handling and signals

**Example API:**
```go
// NewListView creates a new list view
func NewListView(model ListModel) *ListView

// NewTreeView creates a new tree view
func NewTreeView(model TreeModel) *TreeView

// AddColumn adds a column to a tree view
func (t *TreeView) AddColumn(title string, renderer CellRenderer, column int)

// GetSelection gets the selection model
func (t *TreeView) GetSelection() *TreeSelection
```

### Step 6: Menu and Action System

Menu and action system should be implemented next because:
- They provide standardized command handling for applications
- They integrate with modern UI paradigms (app menus, popover menus)
- They build upon signal handling already implemented

**Implementation Details:**
1. Create action system:
   - ActionGroup for organizing actions
   - Action for individual commands
   - ActionMap for hierarchical organization

2. Implement menu components:
   - MenuModel for describing menu structure
   - MenuBar, PopoverMenu, and MenuButton widgets

**Implementation Considerations:**
- Implement GActions properly for modern GTK4 design
- Consider application-wide vs. window-specific actions
- Implement proper keyboard accelerators

**Example API:**
```go
// NewAction creates a new application action
func NewAction(name string, callback func()) *Action

// AddAction adds an action to an action map
func (a *Application) AddAction(action *Action)

// NewMenuModel creates a menu model from a menu description
func NewMenuModel(description string) (*MenuModel, error)

// SetMenuModel sets the application menu model
func (a *Application) SetMenuModel(model *MenuModel)
```

### Step 7: GtkBuilder and UI File Support

GtkBuilder should be implemented next because:
- It allows for UI definitions to be loaded from XML files
- It enables visual design tools to be used with the library
- It provides a more declarative approach to UI development

**Implementation Details:**
1. Create `gtk4/builder.go` with:
   - Builder implementation for loading UI files
   - Object mapping system to connect Go objects to UI elements
   - Signal connection system for UI-defined signals

**Implementation Considerations:**
- Implement proper type conversion between GTK and Go types
- Consider how to expose object properties
- Implement error handling for malformed UI files

**Example API:**
```go
// NewBuilder creates a new builder
func NewBuilder() *Builder

// AddFromFile adds objects from a UI file
func (b *Builder) AddFromFile(filename string) error

// AddFromString adds objects from a UI string
func (b *Builder) AddFromString(uiString string) error

// GetObject gets an object by ID
func (b *Builder) GetObject(id string) (interface{}, error)

// ConnectSignals connects signals to handler functions
func (b *Builder) ConnectSignals(handlers map[string]interface{}) error
```

### Step 8: Advanced Features and Refinement

After implementing the core components, focus on advanced features and refinements:

1. **Clipboard Support**
   - Implement clipboard operations (copy, paste, drag-and-drop)

2. **File System Integration**
   - Implement file monitoring and operations

3. **Application State and Settings**
   - Implement GSettings binding for persistent storage

4. **Internationalization**
   - Add support for translation and locale-specific formatting

5. **Accessibility Features**
   - Implement accessibility interfaces for screen readers

6. **Testing and Documentation**
   - Create comprehensive test suite
   - Generate API documentation
   - Develop example applications

This implementation plan provides a structured approach to completing the GTK4Go library. By following these steps in order, you'll build a comprehensive, usable library with a natural progression from basic to advanced features.
