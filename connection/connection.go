// Package connection provides a Go library for handling connections to Windows-based systems using WinRM and SSH protocols.
package connection

import (
	"context"
)

// Configuration defines the interface for a connection configuration.
type Configuration interface {
	Validate() error
	Defaults()
}

// Connection defines the interface for a connection, specifying methods like Run and Close.
type Connection interface {
	Run(ctx context.Context, cmd string) (CMDResult, error)
	Close() error
}

// CMDResult represents the result of executing a command, including stdout and stderr.
type CMDResult struct {
	StdOut string
	StdErr string
}
