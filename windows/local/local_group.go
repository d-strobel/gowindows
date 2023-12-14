package local

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

type groupType interface {
	Group | []Group
}

// GroupRead gets a group by a SID or Name and returns a Group object.
func (c *LocalClient) GroupRead(ctx context.Context, params GroupParams) (Group, error) {

	// Declare Group object
	var g Group

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return g, fmt.Errorf("windows.local.GroupRead: group parameter 'Name' or 'SID' must be set")
	}

	// We want to retrieve exactly one group.
	if strings.Contains(params.Name, "*") {
		return g, fmt.Errorf("windows.local.GroupRead: group parameter 'Name' does not allow wildcards")
	}

	// Base command
	cmds := []string{"Get-LocalGroup"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	// JSON Output
	cmds = append(cmds, "| ConvertTo-Json")
	cmd := strings.Join(cmds, " ")

	if err := groupRun(ctx, c, cmd, &g); err != nil {
		return g, fmt.Errorf("windows.local.GroupRead: %s", err)
	}

	return g, nil
}

// GroupList returns all groups.
func (c *LocalClient) GroupList(ctx context.Context) ([]Group, error) {

	// Declare slice of Group object
	var g []Group

	// Command
	cmd := "Get-LocalGroup | ConvertTo-Json"

	if err := groupRun(ctx, c, cmd, &g); err != nil {
		return g, fmt.Errorf("windows.local.GroupList: %s", err)
	}

	return g, nil
}

// GroupCreate creates a new group and returns the Group object.
func (c *LocalClient) GroupCreate(ctx context.Context, params GroupParams) (Group, error) {

	// Declare Group object
	var g Group

	// Assert needed parameters
	if params.Name == "" {
		return g, fmt.Errorf("windows.local.GroupCreate: group parameter 'Name' must be set")
	}

	// Base command
	cmds := []string{"New-LocalGroup"}

	// Add parameters
	cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))

	if params.Description != "" {
		cmds = append(cmds, fmt.Sprintf("-Description '%s'", params.Description))
	}

	// JSON Output
	cmds = append(cmds, "| ConvertTo-Json")
	cmd := strings.Join(cmds, " ")

	if err := groupRun(ctx, c, cmd, &g); err != nil {
		return g, fmt.Errorf("windows.local.GroupCreate: %s", err)
	}

	return g, nil
}

// GroupUpdate updates a group.
// Currently only the description parameter can be changed.
func (c *LocalClient) GroupUpdate(ctx context.Context, params GroupParams) error {

	// Satisfy groupType interface
	var g Group

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.GroupUpdate: group parameter 'Name' or 'SID' must be set")
	}

	if params.Description == "" {
		return fmt.Errorf("windows.local.GroupUpdate: group parameter 'Description' must be set")
	}

	// Base command
	cmds := []string{"Set-LocalGroup"}

	// Add parameters
	// Prefer SID over Name to identifiy group
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmds = append(cmds, fmt.Sprintf("-Description '%s'", params.Description))
	cmd := strings.Join(cmds, " ")

	if err := groupRun(ctx, c, cmd, &g); err != nil {
		return fmt.Errorf("windows.local.GroupUpdate: %s", err)
	}

	return nil
}

// GroupDelete removes a group by a SID or Name.
func (c *LocalClient) GroupDelete(ctx context.Context, params GroupParams) error {

	// Satisfy groupType interface
	var g Group

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.GroupDelete: group parameter 'Name' or 'SID' must be set")
	}

	// Base command
	cmds := []string{"Remove-LocalGroup"}

	// Add parameters
	// Prefer SID over Name to identifiy group
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmd := strings.Join(cmds, " ")

	if err := groupRun(ctx, c, cmd, &g); err != nil {
		return fmt.Errorf("windows.local.GroupDelete:\n%s", err)
	}

	return nil
}

// groupRun runs a powershell command against a system.
func groupRun[T groupType](ctx context.Context, c *LocalClient, cmd string, g *T) error {

	// Run the command
	result, err := c.Connection.Run(ctx, cmd)
	if err != nil {
		return err
	}

	// Handle stderr
	if result.StdErr != "" {
		errXML, err := c.parser.DecodeCLIXML(result.StdErr)
		if err != nil {
			return err
		}

		return fmt.Errorf("%s", errXML)
	}

	if result.StdOut == "" {
		return nil
	}

	// Unmarshal stdout
	if err = json.Unmarshal([]byte(result.StdOut), &g); err != nil {
		return err
	}

	return nil
}
