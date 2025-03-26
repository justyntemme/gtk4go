# GTK4Go: Modern GTK4 Bindings for Go

GTK4Go provides comprehensive Go bindings for GTK4, enabling you to build native, cross-platform GUI applications using Go. This library focuses on providing type-safe, memory-safe, and Go-idiomatic access to GTK4's rich widget set and functionality.

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)

## Features

- **Complete GTK4 API coverage**: Access the full power of GTK4 from Go
- **Type-safe bindings**: Properly typed interfaces for GTK4 components
- **Automatic memory management**: Automatic clean-up of GTK resources
- **Background worker support**: Run long tasks without blocking the UI
- **Thread-safe**: Proper handling of GTK's UI thread requirements
- **Modern widgets**: Support for GTK4's new widgets like ListView
- **CSS styling**: Comprehensive styling support with performance optimizations

## Architecture

GTK4Go is structured around several key components that work together to provide a complete GTK4 experience from Go:

### Core Components

1. **CGo Interface Layer**
   - Uses CGo to interface with GTK4's C libraries
   - Exposes GTK functionality through idiomatic Go interfaces
   - Handles type conversion between Go and C types

2. **UI Thread Management**
   - GTK requires all UI operations to happen on the main thread
   - Implemented in `core/uithread` package
   - Provides a message queue system for thread-safe UI operations
   - Automatically routes callbacks and UI updates to the main thread

3. **Unified Callback System**
   - Centralized callback management for all GTK4 signals
   - Maps C signal emissions to Go callback functions
   - Handles memory management for callbacks to prevent leaks
   - Supports various callback signatures for different signal types

4. **Background Worker**
   - Non-blocking background task execution
   - Progress reporting back to the UI thread
   - Cancellation support
   - Error handling

5. **Widget Hierarchy**
   - Object-oriented widget hierarchy mirroring GTK4's structure
   - Common base types with specific implementations
   - Builder pattern with option functions for widget creation

### Memory Management

GTK4Go employs several strategies to ensure proper memory management:

- **Automatic Finalization**: Uses Go's finalizers to clean up widget resources
- **Explicit Destroy Methods**: All widgets have Destroy() methods to explicitly free resources
- **Reference Tracking**: Maintains maps of object references to prevent premature garbage collection
- **Signal Handler Cleanup**: Automatically disconnects signal handlers when widgets are destroyed

### Process Architecture

```
+----------------------------------+
| Go Application                   |
|                                  |
|  +------------------------------+|
|  | GTK4Go API                   ||
|  |                              ||
|  |  +-----------+ +------------+||
|  |  | Widget    | | Background |||
|  |  | Hierarchy | | Worker     |||
|  |  +-----------+ +------------+||
|  |                              ||
|  |  +-----------+ +------------+||
|  |  | Callback  | | UI Thread  |||
|  |  | System    | | Management |||
|  |  +-----------+ +------------+||
|  |                              ||
|  +------------------------------+|
|                                  |
+----------------------------------+
              |
              | CGo
              v
+----------------------------------+
| C GTK4 Libraries                 |
+----------------------------------+
```

## Implementation Details

### Callback System

The callback system is one of the most critical components of GTK4Go. It handles the mapping between GTK signals and Go functions:

1. When a signal is connected, it's registered in a global callback manager
2. A unique ID is generated for the callback
3. The callback and its ID are stored in thread-safe maps
4. When a signal is emitted from C, the callback is looked up by ID and executed on the UI thread

This system supports various callback signatures and allows for safe disconnection of signals.

### UI Thread Management

GTK4 requires all UI operations to happen on the main thread. GTK4Go handles this through the `uithread` package:

1. At initialization, the main UI thread's ID is captured
2. `RunOnUIThread()` function ensures callbacks run on the UI thread
3. Function calls are queued and executed on the main thread
4. Uses GTK's idle functions to schedule execution

This approach ensures thread-safety while maintaining responsiveness.

### Widget Implementation

Widgets follow a pattern:

1. Each widget type has a corresponding Go struct with embedded base types
2. Constructor functions use options pattern for configuration
3. CGo calls are wrapped with type conversion
4. Memory management handled through finalizers and cleanup methods

### Optimizations

GTK4Go includes several optimizations:

1. **CSS Rendering Optimization**: During window resizing, a lightweight CSS provider is temporarily used
2. **Window Resize Detection**: Uses property notifications to detect window resize operations
3. **Provider Caching**: CSS providers are cached to reduce redundant creation
4. **Hardware Acceleration**: Configuration for GPU-accelerated rendering when available

## Usage Example

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/justyntemme/gtk4go"
	"github.com/justyntemme/gtk4go/gtk4"
)

func main() {
	// Initialize GTK (done automatically on import)
	if err := gtk4go.Initialize(); err != nil {
		fmt.Printf("Failed to initialize GTK: %v\n", err)
		os.Exit(1)
	}

	// Create application
	app := gtk4.NewApplication("com.example.HelloWorld")

	// Create window
	win := gtk4.NewWindow("Hello GTK4 from Go!")
	win.SetDefaultSize(400, 300)

	// Create a vertical box
	box := gtk4.NewBox(gtk4.OrientationVertical, 10)

	// Create a label
	label := gtk4.NewLabel("Hello, World!")
	box.Append(label)

	// Create a button
	button := gtk4.NewButton("Click Me")
	button.ConnectClicked(func() {
		label.SetText("Button clicked at " + time.Now().Format(time.RFC3339))
	})
	box.Append(button)

	// Set the window's child
	win.SetChild(box)

	// Add window to application
	app.AddWindow(win)

	// Run the application
	os.Exit(app.Run())
}
```

## Installation

### Prerequisites

- Go 1.18 or later
- GTK4 development libraries (4.8 or later recommended)

### Installing GTK4 Development Libraries

#### Ubuntu/Debian
```bash
sudo apt-get install libgtk-4-dev
```

#### Fedora
```bash
sudo dnf install gtk4-devel
```

#### macOS
```bash
brew install gtk4
```

#### Windows
It's recommended to use MSYS2/MinGW:
```bash
pacman -S mingw-w64-x86_64-gtk4
```

### Installing GTK4Go

```bash
go get github.com/justyntemme/gtk4go
```

## Building Applications

When building applications with GTK4Go, make sure to set up CGo correctly:

```bash
CGO_ENABLED=1 go build
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
