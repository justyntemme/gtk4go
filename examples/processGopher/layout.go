// processGopher/layout.go
package main

import (
	"fmt"

	"github.com/justyntemme/gtk4go/gtk4"
)

// LayoutComponents holds references to all the UI components we need
type LayoutComponents struct {
	MainBox        *gtk4.Box
	ProcessList    *gtk4.ListView
	ProcessModel   *gtk4.StringList
	StatusLabel    *gtk4.Label
	TotalProcLabel *gtk4.Label
}

// CreateProcessLayout creates the main UI layout for the process viewer
// Returns both the layout and references to important components
func CreateProcessLayout() *LayoutComponents {
	// Main vertical box
	mainBox := gtk4.NewBox(gtk4.OrientationVertical, 10)
	mainBox.SetHExpand(true)
	mainBox.SetVExpand(true)

	// Create header
	header := createHeader()
	mainBox.Append(header)

	// Create process list container and get references
	listContainer, processList, processModel := createProcessListContainer()
	mainBox.Append(listContainer)

	// Create status bar and get label references
	statusBar, statusLabel, totalProcLabel := createStatusBar()
	mainBox.Append(statusBar)

	// Apply CSS styling
	applyCSSStyles()

	return &LayoutComponents{
		MainBox:        mainBox,
		ProcessList:    processList,
		ProcessModel:   processModel,
		StatusLabel:    statusLabel,
		TotalProcLabel: totalProcLabel,
	}
}

// createHeader creates the application header
func createHeader() *gtk4.Box {
	headerBox := gtk4.NewBox(gtk4.OrientationVertical, 5)
	headerBox.AddCssClass("header")

	// Title
	title := gtk4.NewLabel("Process Gopher")
	title.AddCssClass("title")
	headerBox.Append(title)

	// Subtitle
	subtitle := gtk4.NewLabel("Real-time System Process Monitor")
	subtitle.AddCssClass("subtitle")
	headerBox.Append(subtitle)

	// Column headers
	columnHeaders := createColumnHeaders()
	headerBox.Append(columnHeaders)

	return headerBox
}

// createColumnHeaders creates the column headers for the process list
func createColumnHeaders() *gtk4.Box {
	headersBox := gtk4.NewBox(gtk4.OrientationHorizontal, 10)
	headersBox.AddCssClass("column-headers")

	headers := []string{"PID", "Process Name", "User", "CPU %", "Memory (MB)"}
	for _, header := range headers {
		label := gtk4.NewLabel(header)
		label.AddCssClass("column-header")
		label.SetHExpand(true)
		headersBox.Append(label)
	}

	return headersBox
}

// createProcessListContainer creates the container for the process list
// Returns the scrolled window, list view, and string model
func createProcessListContainer() (*gtk4.ScrolledWindow, *gtk4.ListView, *gtk4.StringList) {
	// Create scrolled window
	scrollWindow := gtk4.NewScrolledWindow(
		gtk4.WithVScrollbarPolicy(gtk4.ScrollbarPolicyAutomatic),
		gtk4.WithHScrollbarPolicy(gtk4.ScrollbarPolicyNever),
		gtk4.WithVExpand(true),
		gtk4.WithHExpand(true),
	)

	// Create the process ListView and selection model
	stringModel := gtk4.NewStringList()
	selectionModel := gtk4.NewSingleSelection(stringModel,
		gtk4.WithAutoselect(false),
	)

	// Create list item factory
	factory := gtk4.NewSignalListItemFactory()

	// Setup list item creation
	factory.ConnectSetup(func(listItem *gtk4.ListItem) {
		label := gtk4.NewLabel("")
		label.AddCssClass("process-list-item")
		label.SetHExpand(true)
		listItem.SetChild(label)
	})

	// Bind data to list items
	factory.ConnectBind(func(listItem *gtk4.ListItem) {
		// Get the label from the list item
		label, ok := listItem.GetChild().(*gtk4.Label)
		if !ok {
			return
		}

		// Get the text from the model
		text := listItem.GetText()
		if text == "" {
			// Fallback for items without text
			text = fmt.Sprintf("Process %d", listItem.GetPosition()+1)
		}

		// Set the text on the label
		label.SetText(text)

		// Apply alternating row styling
		if listItem.GetPosition()%2 == 0 {
			label.AddCssClass("even-row")
		} else {
			label.AddCssClass("odd-row")
		}
	})

	// Create list view
	listView := gtk4.NewListView(selectionModel, factory,
		gtk4.WithShowSeparators(false),
		gtk4.WithSingleClickActivate(false),
	)
	listView.AddCssClass("process-list")

	// Add list view to scrolled window
	scrollWindow.SetChild(listView)

	return scrollWindow, listView, stringModel
}

// createStatusBar creates the bottom status bar
// Returns the status bar box and references to the labels
func createStatusBar() (*gtk4.Box, *gtk4.Label, *gtk4.Label) {
	statusBox := gtk4.NewBox(gtk4.OrientationHorizontal, 10)
	statusBox.AddCssClass("status-bar")

	// Status label
	statusLabel := gtk4.NewLabel("Status: Loading...")
	statusLabel.AddCssClass("status-text")
	statusBox.Append(statusLabel)

	// Spacer
	spacer := gtk4.NewLabel("")
	spacer.SetHExpand(true)
	statusBox.Append(spacer)

	// Total processes label
	totalProcLabel := gtk4.NewLabel("Total Processes: 0")
	totalProcLabel.AddCssClass("status-text")
	statusBox.Append(totalProcLabel)

	return statusBox, statusLabel, totalProcLabel
}

// applyCSSStyles applies custom CSS styles to the application
func applyCSSStyles() {
	cssContent := `
	.header {
		padding: 20px;
		background-color: #f5f5f5;
		border-bottom: 1px solid #ddd;
	}

	.title {
		font-size: 24px;
		font-weight: bold;
		color: #333;
		margin-bottom: 5px;
	}

	.subtitle {
		font-size: 14px;
		color: #666;
		margin-bottom: 15px;
	}

	.column-headers {
		padding: 10px;
		background-color: #e8e8e8;
		border-bottom: 1px solid #ccc;
	}

	.column-header {
		font-weight: bold;
		color: #444;
		font-size: 14px;
	}

	.process-list {
		font-family: monospace;
		font-size: 13px;
	}

	.process-list-item {
		padding: 5px 10px;
		font-family: monospace;
	}

	.even-row {
		background-color: #fafafa;
	}

	.odd-row {
		background-color: #ffffff;
	}

	.process-list-item:hover {
		background-color: #e8f4fd;
	}

	.status-bar {
		padding: 10px;
		background-color: #f5f5f5;
		border-top: 1px solid #ddd;
	}

	.status-text {
		font-size: 12px;
		color: #666;
	}

	.about-dialog {
		padding: 20px;
		min-width: 300px;
	}

	.about-title {
		font-size: 18px;
		font-weight: bold;
		margin-bottom: 10px;
	}

	.about-description {
		color: #666;
		margin-bottom: 20px;
	}

	.about-author {
		color: #888;
		font-size: 12px;
	}
	`

	// Load CSS
	cssProvider, err := gtk4.LoadCSS(cssContent)
	if err != nil {
		fmt.Printf("Failed to load CSS: %v\n", err)
		return
	}

	// Apply CSS to display
	gtk4.AddProviderForDisplay(cssProvider, 600)
}
