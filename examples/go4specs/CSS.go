package main

import (
	"../../gtk4/"
	"fmt"
)

// loadAppStyles loads the CSS styles for the application
func loadAppStyles() error {
	cssProvider, err := gtk4.LoadCSS(`
		/* Base application styling */
		window {
			background-color: #f5f5f5;
		}
		
		/* Header bar styling with blue color in all states */
		headerbar,
		headerbar:backdrop,
		window .titlebar,
		window:backdrop .titlebar,
		window:active .titlebar,
		window .titlebar headerbar,
		window:backdrop .titlebar headerbar,
		window:active .titlebar headerbar {
			background-color: #3584e4;
			background-image: none;
			color: white;
		}
		
		/* Make all headerbar buttons consistent */
		headerbar button, 
		headerbar button.image-button,
		headerbar button.titlebutton,
		.headerbar-refresh-button {
			background: none;
			background-color: transparent;
			background-image: none;
			border: none;
			border-radius: 50%;
			box-shadow: none;
			outline: none;
			min-width: 16px;
			min-height: 16px;
			padding: 8px;
			margin: 0 2px;
		}
		
		/* Ensure consistent hover effects */
		headerbar button:hover, 
		headerbar button.image-button:hover,
		headerbar button.titlebutton:hover,
		.headerbar-refresh-button:hover {
			background-color: rgba(255, 255, 255, 0.1);
			border: none;
			box-shadow: none;
		}
		
		/* Consistent active/pressed effects */
		headerbar button:active, 
		headerbar button.image-button:active,
		headerbar button.titlebutton:active,
		.headerbar-refresh-button:active {
			background-color: rgba(255, 255, 255, 0.2);
			border: none;
			box-shadow: none;
		}
		
		/* Ensure black icon for our refresh button */
		.headerbar-refresh-button image {
			color: white;
		}

    headerbar button image ,
    headerbar button.image-button image,
    headerbar button.titlebutton image{
      color: white;
    }
		
		/* Clean up focus styles */
		headerbar button:focus, 
		headerbar button.image-button:focus,
		headerbar button.titlebutton:focus,
		.headerbar-refresh-button:focus {
			border: none;
			box-shadow: none;
			outline: none;
		}
		
		/* ======== DISK INFO GRID STYLING ======== */
		.disk-info-grid {
			background-color: #f7f7f7;
			border-radius: 4px;
			padding: 8px;
			margin: 4px 0;
		}
		
		.disk-header {
			font-weight: bold;
			color: #303030;
			padding: 4px 0;
		}
		
		.disk-separator {
			color: #777777;
		}
		
		.disk-device, .info-key {
			font-family: monospace;
			font-weight: bold;
			padding: 3px 0;
		}
		
		.info-value {
			font-family: monospace;
			padding: 3px 0;
		}
		
		/* Add overflow handling for long text values */
		.info-value, .disk-mount {
			text-overflow: ellipsis;
			overflow: hidden;
			white-space: nowrap;
			max-width: 300px;
		}		
		
		/* ======== BUTTON STYLING ======== */
		/* Default button styling - light background with dark text */
		.default-btn {
			background-color: #e8e8e8; 
			color: #303030;
			border-radius: 4px;
			padding: 6px 12px;
			border: none;
			box-shadow: none;
		}
		
		.default-btn label {
			color: #303030;
		}
		
		.default-btn:hover {
			background-color: #f0f0f0;
		}
		
		/* Dark area button - for buttons in dark backgrounds */
		.dark-area-btn {
			background-color: rgba(255, 255, 255, 0.15);
			border-radius: 4px;
			padding: 6px 12px;
			border: none;
		}
		
		.dark-area-btn label {
			color: white;
		}
		
		.dark-area-btn:hover {
			background-color: rgba(255, 255, 255, 0.25);
		}
		
		/* Special button classes */
		.square-button {
			background-color: #3584e4;
			border-radius: 4px;
			padding: 8px 16px;
			font-weight: bold;
			border: none;
			box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
		}
		
		.square-button label {
			color: white;
		}
		
		.square-button:hover {
			background-color: #4a94ea;
		}
		
		.toggle-button {
			background-color: rgba(0, 0, 0, 0.3);
			border-radius: 4px;
			color: black;
			padding: 4px 8px;
			font-size: 12px;
			border: none;
		}
		
		.toggle-button label {
			color: black;
		}
		
		.toggle-button:hover {
			background-color: rgba(0, 0, 0, 0.4);
		}
		
		/* ======== SIDEBAR STYLING ======== */
		.sidebar {
			background-color: #323232;
			min-width: 200px;
			padding: 0;
		}
		
		.sidebar-button {
			background-color: transparent;
			border-radius: 0;
			border-left: 4px solid transparent;
			padding: 16px;
			margin: 0;
			text-align: left;
		}
		
		.sidebar-button label {
			color: #eeeeee;
		}
		
		.sidebar-button:hover {
			background-color: rgba(255, 255, 255, 0.1);
		}
		
		.sidebar-button-selected {
			background-color: rgba(255, 255, 255, 0.15);
			border-left: 4px solid #3584e4;
		}
		
		.sidebar-button-selected label {
			font-weight: bold;
		}
		
		/* ======== CONTENT PANEL STYLING ======== */
		.content-panel {
			padding: 24px;
			background-color: #fafafa;
		}
		
		.panel-title {
			font-size: 22px;
			font-weight: bold;
			margin-bottom: 16px;
			color: #303030;
		}
		
		.info-card {
			background-color: white;
			border-radius: 8px;
			padding: 16px;
			margin-bottom: 16px;
			box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		}
		
		.card-title {
			font-size: 16px;
			font-weight: bold;
			margin-bottom: 8px;
			color: #303030;
		}
		
		.info-grid {
			margin: 8px 0;
		}
		
		.info-key {
			font-weight: normal;
			color: #707070;
			padding-right: 16px;
		}
		
		.info-value {
			font-weight: bold;
			color: #303030;
		}
		
		.disk-info {
			font-family: monospace;
			padding: 12px;
			border-radius: 4px;
			background-color: #f5f5f5;
		}
		
		/* ======== DISK INFO GRID STYLING ======== */
		.disk-info-grid {
			background-color: #f7f7f7;
			border-radius: 4px;
			padding: 8px;
			margin: 4px 0;
		}

		.disk-header {
			font-weight: bold;
			color: #303030;
			padding: 4px 0;
		}

		.disk-separator {
			color: #777777;
		}

		.disk-device {
			font-family: monospace;
			font-weight: bold;
			padding: 3px 0;
		}

		.disk-size, .disk-used, .disk-avail, .disk-percent, .disk-mount {
			font-family: monospace;
			padding: 3px 0;
		}

		.disk-mount {
			text-align: left;
		}

		.disk-usage-normal {
			color: #287c37;
		}

		.disk-usage-warning {
			color: #a85913;
		}

		.disk-usage-critical {
			color: #b00020;
			font-weight: bold;
		}

		.disk-info-error, .disk-info-message {
			font-style: italic;
			color: #707070;
			padding: 8px 4px;
		}
		
		/* ======== STATUS BAR STYLING ======== */
		.status-bar {
			background-color: #323232;
			padding: 8px 16px;
			border-top: 1px solid #444444;
		}
		
		.status-label {
			color: #eeeeee;
		}
		
		.update-time {
			color: #bbbbbb;
			font-size: 12px;
		}
	`)

	if err != nil {
		return fmt.Errorf("failed to load CSS: %v", err.Error())
	} else {
		// Apply CSS provider to the entire application
		gtk4.AddProviderForDisplay(cssProvider, 600)
	}
	return nil
}
