package gowindows

import (
	"github.com/d-strobel/gowindows/package/connection"
	"github.com/d-strobel/gowindows/package/local"
)

type Client struct {
	Connection *connection.Connection
	Local      *local.Client
}

// NewClient returns a Client that contains either a WinRM or SSH client.
// Use this Client to run the functions inside the package directories.
func NewClient(conf *connection.Config) (*Client, error) {

	c := new(Client)
	var err error

	c.Connection, err = connection.NewConnection(conf)
	if err != nil {
		return nil, err
	}

	c.Local = local.NewClient(c.Connection)

	return c, nil
}
