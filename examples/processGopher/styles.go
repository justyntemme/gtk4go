package main

// GetCSS returns the application's CSS styles
func GetCSS() string {
	return `
.heading {
    font-weight: bold;
    font-size: 16px;
    margin-bottom: 10px;
}

.process-list-row:selected {
    background-color: #3584e4;
    color: #ffffff;
}

.warning-text {
    color: #ff7800;
}

.critical-text {
    color: #e01b24;
}

.info-bar {
    padding: 5px;
    background-color: #f6f5f4;
    border-top: 1px solid #d3d7cf;
}

.sidebar {
    background-color: #f6f5f4;
    border-right: 1px solid #d3d7cf;
}

.cpu-high {
    color: #e01b24;
}

.memory-high {
    color: #e01b24;
}

.refresh-button {
    margin-left: 5px;
    margin-right: 5px;
}

.end-process-button {
    background-color: #e01b24;
    color: #ffffff;
    margin-left: 5px;
    margin-right: 5px;
}
`
}
