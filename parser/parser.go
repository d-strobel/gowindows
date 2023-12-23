// Package parser provides a Go library for decoding and parsing data.
package parser

// Parser represents a parser object for decoding and parsing data.
type Parser struct{}

// ParserInterface defines the interface for a parser, specifying methods like DecodeCLIXML.
type ParserInterface interface {
	DecodeCLIXML(clixml string) (string, error)
}

// NewParser returns a new instance of the Parser object.
func NewParser() *Parser {
	return &Parser{}
}
