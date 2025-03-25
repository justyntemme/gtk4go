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
   - The event loop is managed by GTK's application system.

4. **Background Processing**
   - The `BackgroundWorker` provides asynchronous operation capabilities.
   - Tasks can be queued for background processing with progress updates.
   - Results are safely delivered back to the UI thread.

### Component Hierarchy

```
GtkApplication
└── GtkWindow
    └── Container Widgets (e.g., GtkBox, GtkGrid, GtkPaned)
        └── Child Widgets (e.g., GtkButton, GtkLabel, GtkEntry)
```

### Current Implementation

The library currently implements the following components:

- **Application**: Core application functionality
- **Window**: Basic window management
- **Box**: Container for horizontal and vertical layouts
- **Button**: Clickable button with event handling
- **Label**: Text display widget
- **Entry**: Text input widget
- **Grid**: Grid layout container
- **Paned**: Split view container
- **Stack**: Stack container for showing one widget at a time
- **StackSwitcher**: UI control for switching stack pages
- **ScrolledWindow**: Container for scrollable content
- **Viewport**: Container for viewing a portion of a larger area
- **Dialog**: Base dialog and common dialog types
- **Adjustment**: Value adjustment for ranged widgets
- **CSS**: Styling support for widgets

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

### Next Steps

#### 1. ListView and Implementation

ListView and should be implemented next because:

- They're essential for displaying collections of data
- They're complex widgets that build upon previous implementations
- They require model/view architecture that can be reused

**Implementation Details:**

1. Create base model interfaces:

   - ListModel for simple lists

2. Implement view widgets:
   - ListView for flat lists of items
   - Cell renderer system for customizing appearance

**Implementation Considerations:**

- Implement proper model/view separation
- Consider memory management for large datasets
- Implement selection handling and signals

**Example API:**

```go
// NewListView creates a new list view
func NewListView(model ListModel) *ListView



// GetSelection gets the selection model
```

#### 2. Menu and Action System

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

#### 3. GtkBuilder and UI File Support

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

#### 4. Advanced Features and Refinement

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

## Development Status

GTK4Go is currently in a Proof of Concept (PoC) phase. The core architecture and several important widgets have been implemented, but the library is not yet complete or production-ready. Contributions and feedback are welcome.

## License

[MIT License](LICENSE)
