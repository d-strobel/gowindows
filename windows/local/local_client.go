package local

import (
	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
)

type LocalClient struct {
	Connection connection.ConnectionInterface
	parser     parser.ParserInterface
}

// NewLocalClient returns a Client for the local package.
func NewLocalClient(conn *connection.Connection, parser *parser.Parser) *LocalClient {
	return &LocalClient{Connection: conn, parser: parser}
}

// SID is a common struct by all security principals.
// The structure we get from powershell contains more fields, but we're only interested in the Value.
type SID struct {
	Value string `json:"Value"`
}
