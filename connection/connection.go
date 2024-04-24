// Package connection provides a Go library for handling connections to Windows-based systems using WinRM and SSH protocols.
package connection

import (
	"context"
)

// Connection defines the interface for a connection.
// Every connection type must implement this interface.
type Connection interface {

	// Run runs a command using the configured connection and context.
	// It returns the result of the command execution.
	Run(ctx context.Context, cmd string) (CMDResult, error)

	// RunWithPowershell runs a command using the configured connection and context via Powershell.
	// It returns the result of the command execution.
	RunWithPowershell(ctx context.Context, cmd string) (CMDResult, error)

	// Close closes any open connection.
	Close() error
}

// CMDResult represents the result of executing a command, including stdout and stderr.
type CMDResult struct {
	StdOut string
	StdErr string
}
