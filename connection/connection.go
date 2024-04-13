// Package connection provides a Go library for handling connections to Windows-based systems using WinRM and SSH protocols.
package connection

import (
	"context"
)

// Config defines the interface for a connection configuration.
type Config interface {
	validate() error
	defaults()
	NewClient() (Connection, error)
}

// Connection defines the interface for a connection.
// Every connection type must implement this interface.
type Connection interface {
	Run(ctx context.Context, cmd string) (CMDResult, error)
	Close() error
}

// CMDResult represents the result of executing a command, including stdout and stderr.
type CMDResult struct {
	StdOut string
	StdErr string
}
