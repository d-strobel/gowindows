// Package gowindows provides a Go library for interacting remotely with Windows-based systems.
package gowindows

import (
	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
	"github.com/d-strobel/gowindows/windows/local"
)

// Client represents a client object for interacting with Windows systems.
type Client struct {
	Connection connection.Connection
	parser     *parser.Parser
	Local      *local.LocalClient
}

// NewClient returns a new instance of the Client object, initialized with the provided configuration.
// Use this client to execute functions within the Windows subpackages.
func NewClient(conf connection.Config) (*Client, error) {
	var err error

	// Initialize a new client
	c := &Client{}

	// Store a new connection in the client
	c.Connection, err = conf.NewClient()
	if err != nil {
		return nil, err
	}

	// Store a parser in the client
	parser := parser.NewParser()

	// Build the client with the subpackages
	c.Local = local.NewLocalClient(c.Connection, parser)

	return c, nil
}

// Close closes any open connection.
// Currently, only SSH connections will be terminated.
// To avoid surprises in the future, this should always be called using a defer statement.
func (c *Client) Close() error {
	return c.Connection.Close()
}
