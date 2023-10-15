package local

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/d-strobel/gowindows/parser"
)

type Group struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	SID         SID    `json:"SID"`
}

type GroupParams struct {
	Name        string
	Description string
	SID         string
	Context     context.Context
}

var g Group

// GroupRead gets a group by a SID or Name and returns a Group object
func (c *Client) GroupRead(params GroupParams) (*Group, error) {

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return nil, errors.New("Name or SID must be set")
	}
	if params.Context == nil {
		return nil, errors.New("Context must be set")
	}

	// Base command
	cmds := []string{"Get-LocalGroup"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name %s", params.Name))
	}

	cmd := strings.Join(cmds, " ")

	// Optional parameters
	opts := &parser.PwshOpts{
		JSONOutput: true,
	}

	// Powershell command object
	pwshCmd, err := parser.NewPwshCommand([]string{cmd}, opts)

	// Run the comand
	result, err := c.Connection.Run(params.Context, pwshCmd)
	if err != nil {
		return nil, err
	}

	// Handle stderr
	if result.StdErr != "" {
		errXML, err := parser.DecodeCLIXML(result.StdErr)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(errXML)
	}

	// Unmarshal result
	err = json.Unmarshal([]byte(result.StdOut), &g)
	if err != nil {
		return nil, err
	}

	return &g, nil
}
