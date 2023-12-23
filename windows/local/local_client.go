// Package local provides a Go library for handling local Windows functions.
package local

import (
	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
)

// LocalClient represents a client for handling local Windows functions.
type LocalClient struct {
	Connection connection.ConnectionInterface
	parser     parser.ParserInterface
}

// NewLocalClient returns a new instance of the LocalClient.
// It requires a connection and parser as input parameters.
func NewLocalClient(conn *connection.Connection, parser *parser.Parser) *LocalClient {
	return &LocalClient{Connection: conn, parser: parser}
}

// SID represents the Security Identifier (SID) of a security principal.
// The Value field contains the actual SID value.
type SID struct {
	Value string `json:"Value"`
}
