package parser

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type psString string

// Normalize error messages
func (s *psString) UnmarshalText(text []byte) error {
	str := string(text)
	str = strings.TrimSpace(str)
	if str[0] == '+' && len(str) > 2 {
		*s = psString(fmt.Sprintf("\n%s", str[2:]))
	} else {
		*s = psString(str)
	}

	return nil
}

// PSOutput is used to unmarshall CLIXML output
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
func DecodeCLIXML(xmlDoc string) (string, error) {

	if strings.Contains(xmlDoc, "#< CLIXML") {

		var v PSOutput

		xmlDoc = strings.Replace(xmlDoc, "#< CLIXML", "", -1)

		err := xml.Unmarshal([]byte(xmlDoc), &v)
		if err != nil {
			return "", fmt.Errorf("while unmarshalling CLIXML document: %s", err)
		}

		xmlDoc = strings.TrimSpace(v.String())
	}

	return xmlDoc, nil
}
