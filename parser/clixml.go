package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

type clixml struct {
	XML []string `xml:"S"`
}

// unmarshal unmarshals a CLIXML string to an XML object.
func (x *clixml) unmarshal(clixml string) error {
	clixml = strings.ReplaceAll(clixml, "#< CLIXML", "")

	if err := xml.Unmarshal([]byte(clixml), &x); err != nil {
		return err
	}

	return nil
}

// stringSlice removes all unneccessary characters and whitespaces from the XML document
// and returns a new string slice.
func (x *clixml) stringSlice() []string {
	result := make([]string, len(x.XML))

	for i, v := range x.XML {
		s := strings.TrimSpace(v)
		s = strings.ReplaceAll(s, "_x000D__x000A_", "")

		if len(s) > 2 && s[0] == '+' {
			result[i] = fmt.Sprintf("\n%s", s[2:])
		} else {
			result[i] = s
		}
	}

	return result
}

// DecodeCLIXML converts a CLIXML string to a
// human readable powershell error message.
func DecodeCLIXML(xmldoc string) (string, error) {
	if strings.Contains(xmldoc, "#< CLIXML") {
		clixml := &clixml{}

		// Unmarshal to XML
		if err := clixml.unmarshal(xmldoc); err != nil {
			return "", err
		}

		// Convert to stringslice
		s := clixml.stringSlice()

		// Join new Stringslice
		return strings.Join(s, ""), nil
	}

	return "", errors.New("parser.DecodeCLIXML: the input string is not a clixml document")
}
