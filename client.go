// Package gowindows provides a Go library for interacting remotely with Windows-based systems.
package gowindows

import (
	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/windows/local"
)

// Client represents a client object for interacting with Windows systems.
type Client struct {
	Connection connection.Connection
	Local      *local.LocalClient
}

// NewClient returns a new instance of the Client object, initialized with the provided configuration.
// Use this client to execute functions within the Windows subpackages.
func NewClient(conn connection.Connection) *Client {

	// Initialize a new client with the provided connection.
	c := &Client{
		Connection: conn,
	}

	// Build the client with the subpackages.
	c.Local = local.NewClient(c.Connection)

	return c
}

// Close closes any open connection.
func (c *Client) Close() error {
	return c.Connection.Close()
}
