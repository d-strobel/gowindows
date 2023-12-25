package local

import (
	"context"
	"fmt"
	"strings"
)

// Group represents a Windows local group with its properties.
type Group struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	SID         SID    `json:"SID"`
}

// GroupParams represents parameters for interacting with local groups, including creation, updating, and deletion.
type GroupParams struct {
	Name        string
	Description string
	SID         string
}

// GroupRead gets a local group by SID or Name and returns a Group object.
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

	cmds = append(cmds, "| ConvertTo-Json -Compress")
	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[Group](ctx, c, cmd, &g); err != nil {
		return g, fmt.Errorf("windows.local.GroupRead: %s", err)
	}

	return g, nil
}

// GroupList returns a list of all local groups.
func (c *LocalClient) GroupList(ctx context.Context) ([]Group, error) {

	// Declare slice of Group object
	var g []Group

	// Command
	cmd := "Get-LocalGroup | ConvertTo-Json -Compress"

	// Run command
	if err := localRun[[]Group](ctx, c, cmd, &g); err != nil {
		return g, fmt.Errorf("windows.local.GroupList: %s", err)
	}

	return g, nil
}

// GroupCreate creates a new local group and returns the Group object.
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

	cmds = append(cmds, "| ConvertTo-Json -Compress")
	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[Group](ctx, c, cmd, &g); err != nil {
		return g, fmt.Errorf("windows.local.GroupCreate: %s", err)
	}

	return g, nil
}

// GroupUpdate updates a local group. Currently, only the description parameter can be changed.
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

	// Run command
	if err := localRun[Group](ctx, c, cmd, &g); err != nil {
		return fmt.Errorf("windows.local.GroupUpdate: %s", err)
	}

	return nil
}

// GroupDelete removes a local group by SID or Name.
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

	// Run command
	if err := localRun[Group](ctx, c, cmd, &g); err != nil {
		return fmt.Errorf("windows.local.GroupDelete: %s", err)
	}

	return nil
}
