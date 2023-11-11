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
}

var (
	g  Group
	gs []Group
)

// GroupRead gets a group by a SID or Name and returns a Group object
func (c *Client) GroupRead(ctx context.Context, params GroupParams) (*Group, error) {

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return nil, errors.New("Name or SID must be set")
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
	result, err := c.Connection.Run(ctx, pwshCmd)
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

// GroupList returns all groups
func (c *Client) GroupList(ctx context.Context) (*[]Group, error) {

	// Command
	cmd := "Get-LocalGroup"

	// Optional parameters
	opts := &parser.PwshOpts{
		JSONOutput: true,
	}

	// Powershell command object
	pwshCmd, err := parser.NewPwshCommand([]string{cmd}, opts)

	// Run the comand
	result, err := c.Connection.Run(ctx, pwshCmd)
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
	err = json.Unmarshal([]byte(result.StdOut), &gs)
	if err != nil {
		return nil, err
	}

	return &gs, nil
}

// GroupCreate creates a new group and returns the Group object
func (c *Client) GroupCreate(ctx context.Context, params GroupParams) (*Group, error) {

	// Assert needed parameters
	if params.Name == "" {
		return nil, errors.New("Name must be set")
	}

	// Base command
	cmds := []string{"New-LocalGroup"}

	// Add parameters
	cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))

	if params.Description != "" {
		cmds = append(cmds, fmt.Sprintf("-Description '%s'", params.Description))
	}

	cmd := strings.Join(cmds, " ")

	// Optional parameters
	opts := &parser.PwshOpts{
		JSONOutput: true,
	}

	// Powershell command object
	pwshCmd, err := parser.NewPwshCommand([]string{cmd}, opts)

	// Run the comand
	result, err := c.Connection.Run(ctx, pwshCmd)
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
