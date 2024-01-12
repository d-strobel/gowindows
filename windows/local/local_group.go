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

// GroupReadParams represents parameters for the GroupRead function.
type GroupReadParams struct {
	// Specifies the name of the group.
	Name string

	// Specifies the security ID (SID) of the security group.
	SID string
}

// GroupRead gets a local group by SID or Name and returns a Group object.
func (c *LocalClient) GroupRead(ctx context.Context, params GroupReadParams) (Group, error) {

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

// GroupCreateParams represents parameters for the GroupCreate function.
type GroupCreateParams struct {
	// Specifies a name for the group.
	// The maximum length is 256 characters.
	Name string

	// Specifies a comment for the group.
	// The maximum length is 48 characters.
	Description string
}

// GroupCreate creates a new local group and returns the Group object.
func (c *LocalClient) GroupCreate(ctx context.Context, params GroupCreateParams) (Group, error) {

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

// GroupUpdateParams represents parameters for the GroupUpdate function.
type GroupUpdateParams struct {
	// Specifies the name of the group.
	Name string

	// Specifies a comment for the group.
	Description string

	// Specifies the security ID (SID) of the security group.
	SID string
}

// GroupUpdate updates a local group.
func (c *LocalClient) GroupUpdate(ctx context.Context, params GroupUpdateParams) error {

	// Satisfy groupType interface
	var g Group

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.GroupUpdate: group parameter 'Name' or 'SID' must be set")
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

	if params.Description == "" {
		cmds = append(cmds, "-Description ' '")
	} else {
		cmds = append(cmds, fmt.Sprintf("-Description '%s'", params.Description))
	}

	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[Group](ctx, c, cmd, &g); err != nil {
		return fmt.Errorf("windows.local.GroupUpdate: %s", err)
	}

	return nil
}

// GroupDeleteParams represents parameters for the GroupDelete function.
type GroupDeleteParams struct {
	// Specifies the name of the group.
	Name string

	// Specifies the security ID (SID) of the security group.
	SID string
}

// GroupDelete removes a local group by SID or Name.
func (c *LocalClient) GroupDelete(ctx context.Context, params GroupDeleteParams) error {

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
