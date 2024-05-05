package parsing

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

// clixml represents the structure for unmarshaling CLIXML.
type clixml struct {
	Xml []string `xml:"S"`
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
	result := make([]string, len(x.Xml))

	for i, v := range x.Xml {
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

// DecodeCliXmlErr converts a CLIXML error string to a
// human-readable PowerShell error message.
func DecodeCliXmlErr(text string) (string, error) {
	// Check if input string is a valid CLIXML document
	if !strings.Contains(text, "#< CLIXML") {
		return "", errors.New("parsing.DecodeCliXmlErr: the input string is not a CLIXML error string")
	}

	clixml := &clixml{}

	// Unmarshal to XML
	if err := clixml.unmarshal(text); err != nil {
		return "", err
	}

	// Convert to string slice
	s := clixml.stringSlice()

	// Join new string slice
	return strings.Join(s, ""), nil
}
