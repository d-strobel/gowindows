// Package local provides a Go library for handling local Windows functions.
package local

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
)

// LocalClient represents a client for handling local Windows functions.
type LocalClient struct {
	Connection connection.ConnectionInterface
	parser     parser.ParserInterface
}

// localType is an interface for local types.
type localType interface {
	Group | []Group | User | []User | GroupMember | []GroupMember
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

// localRun runs a PowerShell command against a Windows system, handles the command results,
// and unmarshals the output into a local object type.
func localRun[T localType](ctx context.Context, c *LocalClient, cmd string, l *T) error {

	// Run the command
	result, err := c.Connection.Run(ctx, cmd)
	if err != nil {
		return err
	}

	// Handle stderr
	if result.StdErr != "" {
		stderr, err := c.parser.DecodeCLIXML(result.StdErr)
		if err != nil {
			return err
		}

		return errors.New(stderr)
	}

	if result.StdOut == "" {
		return nil
	}

	// Unmarshal stdout
	if err = json.Unmarshal([]byte(result.StdOut), &l); err != nil {
		return err
	}

	return nil
}
