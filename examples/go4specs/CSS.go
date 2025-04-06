package main

import (
	"../../gtk4"
	"fmt"
)

func testing() {
	fmt.Println("hello")
}

func loadAppStyles() error {
	cssProvider, err := gtk4.LoadCSS(`
		window {
			background-color: #f5f5f5;
		}
		
		.header-bar {
			background-color: #3584e4;
			color: white;
			padding: 8px 16px;
			min-height: 48px;
		}
		
		.header-title {
			font-size: 18px;
			font-weight: bold;
			color: white;
		}
		
		.refresh-button {
			padding: 8px 16px;
			background-color: rgba(255, 255, 255, 0.1);
			color: white;
			border-radius: 4px;
		}
		
		.sidebar {
			background-color: #323232;
			min-width: 200px;
			padding: 0;
		}
		
		.sidebar-button {
			background-color: transparent;
			color: #eeeeee;
			border-radius: 0;
			border-left: 4px solid transparent;
			padding: 16px;
			margin: 0;
		}
		
		.sidebar-button:hover {
			background-color: rgba(255, 255, 255, 0.1);
		}
		
		.sidebar-button-selected {
			background-color: rgba(255, 255, 255, 0.15);
			border-left: 4px solid #3584e4;
			font-weight: bold;
		}
		
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
		
		.status-bar {
			background-color: #323232;
			color: #eeeeee;
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
		
		.toggle-button {
			background-color: rgba(255, 255, 255, 0.1);
			border-radius: 4px;
			padding: 4px 8px;
			color: #eeeeee;
			font-size: 12px;
		}
		
		.toggle-button:hover {
			background-color: rgba(255, 255, 255, 0.2);
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
