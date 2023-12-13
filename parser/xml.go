package parser

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type psString string

type PSOutput struct {
	PSStrings []psString `xml:"S"`
}

func (p *PSOutput) stringSlice() []string {
	out := make([]string, len(p.PSStrings))
	for idx, v := range p.PSStrings {
		out[idx] = string(v)
	}
	return out
}

// String() return a string containing the error message that was serialised in a CLIXML message
func (p *PSOutput) String() string {
	str := strings.Join(p.stringSlice(), "")
	replacer := strings.NewReplacer("_x000D_", "", "_x000A_", "")
	str = replacer.Replace(str)
	return str
}

// DecodeCLIXML extracts an error message if stderr is formatted in CLIXML
func DecodeCLIXML(xmlErr string) (string, error) {

	if strings.Contains(xmlErr, "#< CLIXML") {

		var v PSOutput

		xmlErr = strings.Replace(xmlErr, "#< CLIXML", "", -1)

		err := xml.Unmarshal([]byte(xmlErr), &v)
		if err != nil {
			return "", fmt.Errorf("DecodeCLIXML: Failed to unmarshal xml document: %s", err)
		}

		xmlErr = strings.TrimSpace(v.String())
	}

	return xmlErr, nil
}
