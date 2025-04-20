package main

// GetCSS returns the application's CSS styles
func GetCSS() string {
	return `
.heading {
    font-weight: bold;
    font-size: 16px;
    margin-bottom: 10px;
    padding: 5px;
    border-bottom: 1px solid #d3d7cf;
}

.list-item-box {
    padding: 8px 10px;
    border-radius: 2px;
}

.list-item-box:hover {
    background-color: rgba(53, 132, 228, 0.1);
}

.list-item-box.selected {
    background-color: #3584e4;
    color: #ffffff;
}

.list-item-icon {
    color: #2a76c6;
    font-weight: bold;
}

.list-item-label {
    font-family: monospace;
    padding-left: 5px;
}

.header-box {
    padding: 8px;
    border-bottom: 1px solid #d3d7cf;
    background-color: #f6f5f4;
}

.status-bar {
    padding: 5px;
    background-color: #f6f5f4;
    border-top: 1px solid #d3d7cf;
}

.warning-text {
    color: #ff7800;
}

.critical-text {
    color: #e01b24;
}

.performance-tab {
    padding: 10px;
}

.cpu-high {
    color: #e01b24;
}

.memory-high {
    color: #e01b24;
}

.usage-value {
    font-size: 14px;
    padding: 10px;
    margin: 10px;
    background-color: #f9f9f9;
    border-radius: 4px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.usage-normal {
    color: #2e3436;
}

.usage-high {
    color: #e01b24;
    font-weight: bold;
}

.refresh-button {
    margin-left: 5px;
    margin-right: 5px;
    background-color: #33d17a;
    color: white;
}

.end-process-button {
    background-color: #e01b24;
    color: white;
    margin-left: 5px;
    margin-right: 5px;
}

button {
    padding: 5px 10px;
    border-radius: 4px;
}

entry {
    min-width: 200px;
}
`
}
