// Package winerror provides a custom error type for Windows client errors.
package winerror

import (
	"fmt"
)

// WinError represents a custom error type for Windows client errors.
type WinError struct {
	Err     error  // Error message
	Command string // Executed command
}

// Error implements the error interface.
// It returns the error message.
func (e *WinError) Error() string {
	return e.Err.Error()
}

// Unwrap returns the wrapped error.
func (e *WinError) Unwrap() error {
	return e.Err
}

// New creates a new WinError.
func New(cmd string, err error) *WinError {
	return &WinError{
		Err:     err,
		Command: cmd,
	}
}

// Errorf creates a new WinError object from a formatted string.
func Errorf(cmd string, format string, a ...any) *WinError {
	return New(cmd, fmt.Errorf(format, a...))
}

// ExtractCommand extracts the Command field from the error.
func UnwrapCommand(err error) string {
	if e, ok := err.(*WinError); ok {
		return e.Command
	}
	return ""
}
