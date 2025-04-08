//go:build linux
// +build linux

package uithread

// Linux doesn't require special thread handling for GTK
// The init function is empty as the main thread.go file handles everything