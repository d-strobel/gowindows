package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

// clixml represents the structure for unmarshaling CLIXML.
type clixml struct {
	XML []string `xml:"S"`
}

// unmarshal unmarshals a CLIXML string to an XML object.
func (x *clixml) unmarshal(clixml string) error {
	// Remove CLIXML identifier
	clixml = strings.ReplaceAll(clixml, "#< CLIXML", "")

	// Unmarshal to XML
	if err := xml.Unmarshal([]byte(clixml), &x); err != nil {
		return err
	}

	return nil
}

// stringSlice removes all unnecessary characters and whitespaces from the XML document
// and returns a new string slice.
func (x *clixml) stringSlice() []string {
	result := make([]string, len(x.XML))

	for i, v := range x.XML {
		// Trim whitespaces
		s := strings.TrimSpace(v)

		// Remove specific characters
		s = strings.ReplaceAll(s, "_x000D__x000A_", "")

		// Handle line continuation
		if len(s) > 2 && s[0] == '+' {
			result[i] = fmt.Sprintf("\n%s", s[2:])
		} else {
			result[i] = s
		}
	}

	return result
}

// DecodeCLIXML converts a CLIXML string to a
// human-readable PowerShell error message.
func (p *Parser) DecodeCLIXML(xmldoc string) (string, error) {
	if strings.Contains(xmldoc, "#< CLIXML") {
		clixml := &clixml{}

		// Unmarshal to XML
		if err := clixml.unmarshal(xmldoc); err != nil {
			return "", err
		}

		// Convert to string slice
		s := clixml.stringSlice()

		// Join new string slice
		return strings.Join(s, ""), nil
	}

	return "", errors.New("parser.DecodeCLIXML: the input string is not a CLIXML document")
}
