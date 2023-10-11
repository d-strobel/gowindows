package local

import (
	"github.com/d-strobel/gowindows/connection"
)

type Client struct {
	Connection *connection.Connection
}

// New returns a Client for the local package
func New(conn *connection.Connection) *Client {
	return &Client{Connection: conn}
}

// SID is a common struct by all security principals
// The structure we get from powershell contains more fields, but we're only interested in the Value.
type SID struct {
	Value string `json:"Value"`
}
