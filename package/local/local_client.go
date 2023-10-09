package local

import (
	"github.com/d-strobel/gowindows/package/connection"
)

type Client struct {
	Connection *connection.Connection
}

// New returns a Client for the local package
func New(conn *connection.Connection) *Client {
	return &Client{Connection: conn}
}
