package gowindows

import (
	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/windows/local"
)

type Client struct {
	Connection *connection.Connection
	Local      *local.LocalClient
}

// New returns a Client object that contains the Connection and the Windows package.
// Use this Client to run the functions inside the Windows subpackages.
func New(conf *connection.Config) (*Client, error) {

	var err error

	// Init new Client
	c := &Client{}

	// Store new connection to the Client
	c.Connection, err = connection.New(conf)
	if err != nil {
		return nil, err
	}

	// Build the client with the subpackages
	c.Local = local.New(c.Connection)

	return c, nil
}

// Close closes any open connection.
// For now only ssh connections will be terminated.
// To avoid surprises in the future, this should always be called in a defer statement.
func (c *Client) Close() error {
	return c.Connection.Close()
}
