package winerror

import (
	"fmt"
)

const (
	// Root errors
	ConfigError     = "configuration_error"
	ConnectionError = "connection_error"
	WindowsError    = "windows_error"
	ParserError     = "parser_error"
)

// Error type for gowindows
type Winerror struct {
	Root string
	Msg  string
}

// Satisfy the Error interface
func (e *Winerror) Error() string {
	return fmt.Sprintf("[%s] %s", e.Root, e.Msg)
}

// Errorf creates a new Winerror
func Errorf(root string, format string, a ...interface{}) *Winerror {

	msg := fmt.Sprintf(format, a...)

	return &Winerror{
		Root: root,
		Msg:  msg,
	}
}
