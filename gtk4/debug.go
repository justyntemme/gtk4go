// Package gtk4 provides debugging utilities for GTK4 components
// File: gtk4go/gtk4/debug.go
package gtk4

import (
	"fmt"
	"log"
	"sync"
)

// Debug logging levels
const (
	DebugLevelNone    = 0 // No debug output
	DebugLevelError   = 1 // Error messages only
	DebugLevelWarning = 2 // Warnings and errors
	DebugLevelInfo    = 3 // Informational messages
	DebugLevelVerbose = 4 // Verbose output
)

// Component identifiers for filtering debug output
const (
	DebugComponentGeneral     = "general"
	DebugComponentCallback    = "callback"
	DebugComponentDialog      = "dialog"
	DebugComponentListView    = "listview"
	DebugComponentListFactory = "listfactory"
	DebugComponentAction      = "action"
	DebugComponentSelection   = "selection"
)

// Global debug configuration
var (
	debugLevel         = DebugLevelNone
	debugFilter        = make(map[string]bool)
	debugLogPrefix     = "[GTK4Go] "
	debugMutex         sync.RWMutex
	debugToStdErr      = false
)

// SetDebugLevel sets the global debug level
func SetDebugLevel(level int) {
	debugMutex.Lock()
	defer debugMutex.Unlock()
	debugLevel = level
}

// GetDebugLevel gets the current debug level
func GetDebugLevel() int {
	debugMutex.RLock()
	defer debugMutex.RUnlock()
	return debugLevel
}

// EnableComponent enables debug output for a specific component
func EnableDebugComponent(component string) {
	debugMutex.Lock()
	defer debugMutex.Unlock()
	debugFilter[component] = true
}

// DisableComponent disables debug output for a specific component
func DisableDebugComponent(component string) {
	debugMutex.Lock()
	defer debugMutex.Unlock()
	debugFilter[component] = false
}

// EnableAllComponents enables debug output for all components
func EnableAllDebugComponents() {
	debugMutex.Lock()
	defer debugMutex.Unlock()
	// Add all known components
	debugFilter[DebugComponentGeneral] = true
	debugFilter[DebugComponentCallback] = true
	debugFilter[DebugComponentDialog] = true
	debugFilter[DebugComponentListView] = true
	debugFilter[DebugComponentListFactory] = true
	debugFilter[DebugComponentAction] = true
	debugFilter[DebugComponentSelection] = true
}

// SetDebugToStdErr sets whether debug output should go to stderr
func SetDebugToStdErr(useStdErr bool) {
	debugMutex.Lock()
	defer debugMutex.Unlock()
	debugToStdErr = useStdErr
}

// DebugLog logs a debug message if the current level is high enough
// and the component is enabled for debugging
func DebugLog(level int, component string, format string, args ...interface{}) {
	debugMutex.RLock()
	currentLevel := debugLevel
	isComponentEnabled, exists := debugFilter[component]
	useStdErr := debugToStdErr
	debugMutex.RUnlock()

	// Only log if the level is appropriate and component is enabled (or not explicitly disabled)
	if level <= currentLevel && (isComponentEnabled || !exists) {
		message := fmt.Sprintf(format, args...)
		logMessage := fmt.Sprintf("%s[%s] %s", debugLogPrefix, component, message)
		
		if useStdErr {
			log.Printf("%s\n", logMessage)
		} else {
			fmt.Printf("%s\n", logMessage)
		}
	}
}