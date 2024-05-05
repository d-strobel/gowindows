package localaccounts

import (
	"context"
	"fmt"
	"strings"
)

// GroupMember represents a member of a local Windows group.
type GroupMember struct {
	Name        string `json:"Name"`
	SID         SID    `json:"SID"`
	ObjectClass string `json:"ObjectClass"`
}

// GroupMemberReadParams represent parameters for the GroupMemberRead function.
type GroupMemberReadParams struct {
	// Specifies the name of the security group.
	Name string

	// Specifies the security ID of the security group.
	SID string

	// Specifies a user or group of the security group.
	Member string
}

// pwshCommad returns a PowerShell command for reading a local group member.
func (params GroupMemberReadParams) pwshCommand() string {
	// Base command
	cmd := []string{"Get-LocalGroupMember"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmd = append(cmd, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmd = append(cmd, fmt.Sprintf("-Member '%s'", params.Member))
	cmd = append(cmd, "| ConvertTo-Json -Compress")
	return strings.Join(cmd, " ")
}

// GroupMemberRead retrieves information about a specific member in a local Windows group.
func (c *Client) GroupMemberRead(ctx context.Context, params GroupMemberReadParams) (GroupMember, error) {
	var gm GroupMember

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return gm, fmt.Errorf("windows.localaccounts.GroupMemberRead: group member parameter 'Name' or 'SID' must be set")
	}

	if params.Member == "" {
		return gm, fmt.Errorf("windows.localaccounts.GroupMemberRead: group member parameter 'Member' must be set")
	}

	// Run command
	if err := run(ctx, c, params.pwshCommand(), &gm); err != nil {
		return gm, fmt.Errorf("windows.localaccounts.GroupMemberRead: %s", err)
	}

	return gm, nil
}

// GroupMemberListParams represent parameters for the GroupMemberList function.
type GroupMemberListParams struct {
	// Specifies the name of the security group.
	Name string

	// Specifies the security ID of the security group.
	SID string
}

// pwshCommand returns a PowerShell command for listing members of a local group.
func (params GroupMemberListParams) pwshCommand() string {
	// Base command
	cmd := []string{"$gm=Get-LocalGroupMember"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmd = append(cmd, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	}

	// Ensure that groups with a single group member is also printed as an array
	cmd = append(cmd, ";if($gm.Count -eq 1){ConvertTo-Json @($gm) -Compress}else{ConvertTo-Json $gm -Compress}")
	return strings.Join(cmd, " ")
}

// GroupMemberList returns a list of members for a specific local Windows group.
func (c *Client) GroupMemberList(ctx context.Context, params GroupMemberListParams) ([]GroupMember, error) {
	var gm []GroupMember

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return gm, fmt.Errorf("windows.localaccounts.GroupMemberList: group member parameter 'Name' or 'SID' must be set")
	}

	// Run command
	if err := run(ctx, c, params.pwshCommand(), &gm); err != nil {
		return gm, fmt.Errorf("windows.localaccounts.GroupMemberList: %s", err)
	}

	return gm, nil
}

// GroupMemberCreateParams represent parameters for the GroupMemberCreate function.
type GroupMemberCreateParams struct {
	// Specifies the name of the security group.
	Name string

	// Specifies the security ID of the security group.
	SID string

	// Specifies a new user or group for the security group.
	Member string
}

// pwshCommand returns a PowerShell command for adding a new member to a local group.
func (params GroupMemberCreateParams) pwshCommand() string {
	// Base command
	cmd := []string{"Add-LocalGroupMember"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmd = append(cmd, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmd = append(cmd, fmt.Sprintf("-Member '%s'", params.Member))
	return strings.Join(cmd, " ")
}

// GroupMemberCreate adds a new member to a local Windows group.
func (c *Client) GroupMemberCreate(ctx context.Context, params GroupMemberCreateParams) error {
	var gm GroupMember

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.localaccounts.GroupMemberCreate: group member parameter 'Name' or 'SID' must be set")
	}

	if params.Member == "" {
		return fmt.Errorf("windows.localaccounts.GroupMemberCreate: group member parameter 'Member' must be set")
	}

	// Run command
	if err := run(ctx, c, params.pwshCommand(), &gm); err != nil {
		return fmt.Errorf("windows.localaccounts.GroupMemberCreate: %s", err)
	}

	return nil
}

// GroupMemberDeleteParams represent parameters for the GroupMemberDelete function.
type GroupMemberDeleteParams struct {
	// Specifies the name of the security group.
	Name string

	// Specifies the security ID of the security group.
	SID string

	// Specifies a user or group of the security group.
	Member string
}

// pwshCommand returns a PowerShell command for removing a member from a local group.
func (params GroupMemberDeleteParams) pwshCommand() string {
	// Base command
	cmd := []string{"Remove-LocalGroupMember"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmd = append(cmd, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	}
	cmd = append(cmd, fmt.Sprintf("-Member '%s'", params.Member))

	return strings.Join(cmd, " ")
}

// GroupMemberDelete removes a member from a local Windows group.
func (c *Client) GroupMemberDelete(ctx context.Context, params GroupMemberDeleteParams) error {
	var gm GroupMember

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.localaccounts.GroupMemberDelete: group member parameter 'Name' or 'SID' must be set")
	}

	if params.Member == "" {
		return fmt.Errorf("windows.localaccounts.GroupMemberDelete: group member parameter 'Member' must be set")
	}

	// Run command
	if err := run(ctx, c, params.pwshCommand(), &gm); err != nil {
		return fmt.Errorf("windows.localaccounts.GroupMemberDelete: %s", err)
	}

	return nil
}
