// Package connection defines an interface and utilities for establishing and managing generic connections to remote systems.
// It provides an abstraction layer for executing commands and managing the lifecycle of connections, facilitating interoperability with various connection types such as WinRM and SSH.
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
	// StdOut contains the standard output of the command.
	StdOut string

	// StdErr contains the standard error output of the command.
	StdErr string
}
