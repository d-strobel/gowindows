package local

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	// Specifies the user name for the user account.
	Name string
	// Specifies a comment for the user account.
	// The maximum length is 48 characters.
	Description string
	// Specifies a security ID (SID) of user account.
	SID string
	// Specifies when the user account expires.
	// If you don't specify this parameter, the account doesn't expire.
	AccountExpires time.Time
	// Indicates wheter the account is disabled.
	Disabled bool
	// Specifies the full name for the user account.
	// The full name differs from the user name of the user account.
	FullName string
	// Specifies a password for the user account.
	Password string
	// Indicates whether the new user's password expires.
	PasswordNeverExpires bool
	// Indicates that the user can't change the password on the user account.
	UserMayNotChangePassword bool
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

// UserCreate creates a local user and returns a User object.
func (c *LocalClient) UserCreate(ctx context.Context, params UserParams) (User, error) {

	// Declare User object
	var u User

	// Assert needed parameters
	if params.Name == "" {
		return u, fmt.Errorf("windows.local.UserCreate: user parameter 'Name' must be set")
	}

	// Base command
	cmds := []string{"New-LocalUser"}

	// Add parameters
	cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))

	if params.Description != "" {
		cmds = append(cmds, fmt.Sprintf("-Description '%s'", params.Description))
	}

	if params.AccountExpires.Compare(time.Now()) == 1 {
		accountExpires := params.AccountExpires.Format(time.DateTime)
		cmds = append(cmds, fmt.Sprintf("-AccountExpires $(Get-Date '%s')", accountExpires))
	} else {
		cmds = append(cmds, "-AccountNeverExpires")
	}

	if params.Disabled {
		cmds = append(cmds, "-Disabled")
	}

	if params.FullName != "" {
		cmds = append(cmds, fmt.Sprintf("-FullName '%s'", params.FullName))
	}

	if params.Password != "" {
		cmds = append(cmds, fmt.Sprintf("-Password $(ConvertTo-SecureString -String '%s' -AsPlainText -Force)", params.Password))
	} else {
		cmds = append(cmds, "-NoPassword")
	}

	if params.PasswordNeverExpires {
		cmds = append(cmds, "-PasswordNeverExpires")
	}

	if params.UserMayNotChangePassword {
		cmds = append(cmds, "-UserMayNotChangePassword")
	}

	cmds = append(cmds, "| ConvertTo-Json -Compress")
	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[User](ctx, c, cmd, &u); err != nil {
		return u, fmt.Errorf("windows.local.UserCreate: %s", err)
	}

	return u, nil
}