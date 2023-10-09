package gowindows

import (
	"github.com/d-strobel/gowindows/package/connection"
	"github.com/d-strobel/gowindows/package/local"
)

type Client struct {
	Connection *connection.Connection
	Local      *local.Client
}

// NewClient returns a Client object that contains the Connection and the subpackages.
// Use this Client to run the functions inside the package directories.
func NewClient(conf *connection.Config) (*Client, error) {

	var err error

	// Allocate a new Client
	c := new(Client)

	// Store new connection to the Client
	c.Connection, err = connection.NewConnection(conf)
	if err != nil {
		return nil, err
	}

	// Build the client with the subpackages
	c.Local = local.NewClient(c.Connection)

	return c, nil
}

// Close closes all active connections
func (c *Client) Close() error {
	return c.Connection.Close()
}
