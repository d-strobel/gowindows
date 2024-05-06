// Package gowindows provides a Go library for interacting remotely with Windows-based systems.
package gowindows

import (
	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/windows/local/accounts"
)

// Mockery generates mocks for the gowindows interfaces.
// Use "go generate" to update the mocks.
//go:generate go run github.com/vektra/mockery/v2

// Client represents a client object for interacting with Windows systems.
type Client struct {
	Connection    connection.Connection
	LocalAccounts *accounts.Client
}

// NewClient returns a new instance of the Client object, initialized with the provided configuration.
// Use this client to execute functions within the Windows subpackages.
func NewClient(conn connection.Connection) *Client {

	// Initialize a new client with the provided connection.
	c := &Client{
		Connection: conn,
	}

	// Build the client with the subpackages.
	c.LocalAccounts = accounts.NewClient(c.Connection)

	return c
}

// Close closes any open connection.
func (c *Client) Close() error {
	return c.Connection.Close()
}
