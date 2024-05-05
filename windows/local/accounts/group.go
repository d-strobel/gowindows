package accounts

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

// pwshCommand returns the PowerShell command to read a local group by SID or Name.
func (params GroupReadParams) pwshCommand() string {
	// Base command
	cmd := []string{"Get-LocalGroup"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmd = append(cmd, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmd = append(cmd, "| ConvertTo-Json -Compress")
	return strings.Join(cmd, " ")
}

// GroupRead gets a local group by SID or Name and returns a Group object.
func (c *Client) GroupRead(ctx context.Context, params GroupReadParams) (Group, error) {
	var g Group

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return g, fmt.Errorf("windows.local.accounts.GroupRead: group parameter 'Name' or 'SID' must be set")
	}

	// We want to retrieve exactly one group.
	if strings.Contains(params.Name, "*") {
		return g, fmt.Errorf("windows.local.accounts.GroupRead: group parameter 'Name' does not allow wildcards")
	}

	// Run command
	if err := run(ctx, c, params.pwshCommand(), &g); err != nil {
		return g, fmt.Errorf("windows.local.accounts.GroupRead: %s", err)
	}
	return g, nil
}

// GroupList returns a list of all local groups.
func (c *Client) GroupList(ctx context.Context) ([]Group, error) {
	var g []Group

	// Command
	cmd := "Get-LocalGroup | ConvertTo-Json -Compress"

	// Run command
	if err := run(ctx, c, cmd, &g); err != nil {
		return g, fmt.Errorf("windows.local.accounts.GroupList: %s", err)
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

// pwshCommand returns the PowerShell command to create a local group.
func (params GroupCreateParams) pwshCommand() string {
	// Base command
	cmd := []string{"New-LocalGroup"}

	// Add parameters
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))

	if params.Description != "" {
		cmd = append(cmd, fmt.Sprintf("-Description '%s'", params.Description))
	}

	cmd = append(cmd, "| ConvertTo-Json -Compress")
	return strings.Join(cmd, " ")
}

// GroupCreate creates a new local group and returns the Group object.
func (c *Client) GroupCreate(ctx context.Context, params GroupCreateParams) (Group, error) {
	var g Group

	// Assert needed parameters
	if params.Name == "" {
		return g, fmt.Errorf("windows.local.accounts.GroupCreate: group parameter 'Name' must be set")
	}

	// Run command
	if err := run(ctx, c, params.pwshCommand(), &g); err != nil {
		return g, fmt.Errorf("windows.local.accounts.GroupCreate: %s", err)
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

// pwshCommand returns the PowerShell command to update a local group.
func (params GroupUpdateParams) pwshCommand() string {
	// Base command
	cmd := []string{"Set-LocalGroup"}

	// Add parameters
	// Prefer SID over Name to identifiy group
	if params.SID != "" {
		cmd = append(cmd, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	}

	if params.Description == "" {
		cmd = append(cmd, "-Description ' '")
	} else {
		cmd = append(cmd, fmt.Sprintf("-Description '%s'", params.Description))
	}

	return strings.Join(cmd, " ")
}

// GroupUpdate updates a local group.
func (c *Client) GroupUpdate(ctx context.Context, params GroupUpdateParams) error {
	var g Group

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.accounts.GroupUpdate: group parameter 'Name' or 'SID' must be set")
	}

	// Run command
	if err := run(ctx, c, params.pwshCommand(), &g); err != nil {
		return fmt.Errorf("windows.local.accounts.GroupUpdate: %s", err)
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

// pwshCommand returns the PowerShell command to delete a local group by SID or Name.
func (params GroupDeleteParams) pwshCommand() string {
	// Base command
	cmd := []string{"Remove-LocalGroup"}

	// Add parameters
	// Prefer SID over Name to identifiy group
	if params.SID != "" {
		cmd = append(cmd, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	}

	return strings.Join(cmd, " ")
}

// GroupDelete removes a local group by SID or Name.
func (c *Client) GroupDelete(ctx context.Context, params GroupDeleteParams) error {
	var g Group

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.accounts.GroupDelete: group parameter 'Name' or 'SID' must be set")
	}

	// Run command
	if err := run(ctx, c, params.pwshCommand(), &g); err != nil {
		return fmt.Errorf("windows.local.accounts.GroupDelete: %s", err)
	}

	return nil
}
