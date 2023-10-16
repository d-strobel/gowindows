package gowindows

import (
	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/windows/local"
)

type Client struct {
	Connection *connection.Connection
	Local      *local.Client
}

// NewClient returns a Client object that contains the Connection and the subpackages.
// Use this Client to run the functions inside the package directories.
func New(conf *connection.Config) (*Client, error) {

	var err error

	// Allocate a new Client
	c := new(Client)

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
// Only ssh connection will be terminated here.
// To avoid surprises in the future, this should always be called in a defer statement.
func (c *Client) Close() error {
	return c.Connection.Close()
}
