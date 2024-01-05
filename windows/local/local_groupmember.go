package local

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

// GroupMemberRead retrieves information about a specific member in a local Windows group.
func (c *LocalClient) GroupMemberRead(ctx context.Context, params GroupMemberReadParams) (GroupMember, error) {
	// Declare GroupMember
	var gm GroupMember

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return gm, fmt.Errorf("windows.local.GroupMemberRead: group member parameter 'Name' or 'SID' must be set")
	}

	if params.Member == "" {
		return gm, fmt.Errorf("windows.local.GroupMemberRead: group member parameter 'Member' must be set")
	}

	// Base command
	cmds := []string{"Get-LocalGroupMember"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmds = append(cmds, fmt.Sprintf("-Member '%s'", params.Member))
	cmds = append(cmds, "| ConvertTo-Json -Compress")
	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[GroupMember](ctx, c, cmd, &gm); err != nil {
		return gm, fmt.Errorf("windows.local.GroupMemberRead: %s", err)
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

// GroupMemberList returns a list of members for a specific local Windows group.
func (c *LocalClient) GroupMemberList(ctx context.Context, params GroupMemberListParams) ([]GroupMember, error) {
	// Declare slice of GroupMember
	var gm []GroupMember

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return gm, fmt.Errorf("windows.local.GroupMemberList: group member parameter 'Name' or 'SID' must be set")
	}

	// Base command
	cmds := []string{"$gm=Get-LocalGroupMember"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	// Ensure that groups with a single group member is also printed as an array
	cmds = append(cmds, ";if($gm.Count -eq 1){ConvertTo-Json @($gm) -Compress}else{ConvertTo-Json $gm -Compress}")
	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[[]GroupMember](ctx, c, cmd, &gm); err != nil {
		return gm, fmt.Errorf("windows.local.GroupMemberList: %s", err)
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

// GroupMemberCreate adds a new member to a local Windows group.
func (c *LocalClient) GroupMemberCreate(ctx context.Context, params GroupMemberCreateParams) error {
	// Satisfy the localType interface
	var gm GroupMember

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.GroupMemberCreate: group member parameter 'Name' or 'SID' must be set")
	}

	if params.Member == "" {
		return fmt.Errorf("windows.local.GroupMemberCreate: group member parameter 'Member' must be set")
	}

	// Base command
	cmds := []string{"Add-LocalGroupMember"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmds = append(cmds, fmt.Sprintf("-Member '%s'", params.Member))
	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[GroupMember](ctx, c, cmd, &gm); err != nil {
		return fmt.Errorf("windows.local.GroupMemberCreate: %s", err)
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

// GroupMemberDelete removes a member from a local Windows group.
func (c *LocalClient) GroupMemberDelete(ctx context.Context, params GroupMemberDeleteParams) error {
	// Satisfy the localType interface
	var gm GroupMember

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.GroupMemberDelete: group member parameter 'Name' or 'SID' must be set")
	}

	if params.Member == "" {
		return fmt.Errorf("windows.local.GroupMemberDelete: group member parameter 'Member' must be set")
	}

	// Base command
	cmds := []string{"Remove-LocalGroupMember"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmds = append(cmds, fmt.Sprintf("-Member '%s'", params.Member))
	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[GroupMember](ctx, c, cmd, &gm); err != nil {
		return fmt.Errorf("windows.local.GroupMemberDelete: %s", err)
	}

	return nil
}
