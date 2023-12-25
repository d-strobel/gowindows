package local

import (
	"context"
	"fmt"
	"strings"

	"github.com/d-strobel/gowindows/parser"
)

// User represents a Windows local user with its properties.
type User struct {
	AccountExpires         parser.WinTime `json:"AccountExpires"`
	Description            string         `json:"Description"`
	Enabled                bool           `json:"Enabled"`
	FullName               string         `json:"FullName"`
	PasswordChangeableDate parser.WinTime `json:"PasswordChangeableDate"`
	PasswordExpires        parser.WinTime `json:"PasswordExpires"`
	UserMayChangePassword  bool           `json:"UserMayChangePassword"`
	PasswordRequired       bool           `json:"PasswordRequired"`
	PasswordLastSet        parser.WinTime `json:"PasswordLastSet"`
	LastLogon              parser.WinTime `json:"LastLogon"`
	Name                   string         `json:"Name"`
	SID                    SID            `json:"SID"`
}

// UserParams represents parameters for interacting with local users, including creation, updating, and deletion.
type UserParams struct {
	Name        string
	Description string
	SID         string
}

// UserRead gets a local user by SID or Name and returns a User object.
func (c *LocalClient) UserRead(ctx context.Context, params UserParams) (User, error) {

	// Declare User object
	var u User

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return u, fmt.Errorf("windows.local.UserRead: user parameter 'Name' or 'SID' must be set")
	}

	// We want to retrieve exactly one user.
	if strings.Contains(params.Name, "*") {
		return u, fmt.Errorf("windows.local.UserRead: user parameter 'Name' does not allow wildcards")
	}

	// Base command
	cmds := []string{"Get-LocalUser"}

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
	if err := localRun[User](ctx, c, cmd, &u); err != nil {
		return u, fmt.Errorf("windows.local.UserRead: %s", err)
	}

	return u, nil
}

// UserList returns a list of all local user.
func (c *LocalClient) UserList(ctx context.Context) ([]User, error) {

	// Declare slice of User
	var u []User

	// Command
	cmd := "Get-LocalUser | ConvertTo-Json -Compress"

	// Run command
	if err := localRun[[]User](ctx, c, cmd, &u); err != nil {
		return u, fmt.Errorf("windows.local.UserList: %s", err)
	}

	return u, nil
}
