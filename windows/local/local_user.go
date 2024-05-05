package local

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/d-strobel/gowindows/parsing"
)

// User represents a Windows local user with its properties.
type User struct {
	AccountExpires         parsing.DotnetTime `json:"AccountExpires"`
	Description            string             `json:"Description"`
	Enabled                bool               `json:"Enabled"`
	FullName               string             `json:"FullName"`
	PasswordChangeableDate parsing.DotnetTime `json:"PasswordChangeableDate"`
	PasswordExpires        parsing.DotnetTime `json:"PasswordExpires"`
	UserMayChangePassword  bool               `json:"UserMayChangePassword"`
	PasswordRequired       bool               `json:"PasswordRequired"`
	PasswordLastSet        parsing.DotnetTime `json:"PasswordLastSet"`
	LastLogon              parsing.DotnetTime `json:"LastLogon"`
	Name                   string             `json:"Name"`
	SID                    SID                `json:"SID"`
}

// UserReadParams represents parameters for the UserRead function.
type UserReadParams struct {
	// Specifies the user name of the user account.
	Name string

	// Specifies a security ID (SID) of user account.
	SID string
}

// pwshCommand returns a PowerShell command for retrieving a local user.
func (params UserReadParams) pwshCommand() string {
	// Base command
	cmd := []string{"Get-LocalUser"}

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

// UserRead gets a local user by SID or Name and returns a User object.
func (c *LocalClient) UserRead(ctx context.Context, params UserReadParams) (User, error) {
	var u User

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return u, fmt.Errorf("windows.local.UserRead: user parameter 'Name' or 'SID' must be set")
	}

	// We want to retrieve exactly one user.
	if strings.Contains(params.Name, "*") {
		return u, fmt.Errorf("windows.local.UserRead: user parameter 'Name' does not allow wildcards")
	}

	// Run command
	if err := localRun(ctx, c, params.pwshCommand(), &u); err != nil {
		return u, fmt.Errorf("windows.local.UserRead: %s", err)
	}

	return u, nil
}

// UserList returns a list of all local user.
func (c *LocalClient) UserList(ctx context.Context) ([]User, error) {
	var u []User

	// Command
	cmd := "Get-LocalUser | ConvertTo-Json -Compress"

	// Run command
	if err := localRun(ctx, c, cmd, &u); err != nil {
		return u, fmt.Errorf("windows.local.UserList: %s", err)
	}

	return u, nil
}

// UserCreateParams represents parameters for the UserCreate function.
type UserCreateParams struct {
	// Specifies the user name for the user account.
	Name string

	// Specifies a comment for the user account.
	// The maximum length is 48 characters.
	Description string

	// Specifies when the user account expires.
	// If you don't specify this parameter, the account doesn't expire.
	AccountExpires time.Time

	// Indicates wheter the account is enabled.
	Enabled bool

	// Specifies the full name for the user account.
	// The full name differs from the user name of the user account.
	FullName string

	// Specifies a password for the user account.
	Password string

	// Indicates whether the new user's password expires.
	PasswordNeverExpires bool

	// Indicates that the user can change the password on the user account.
	UserMayChangePassword bool
}

// pwshCommand returns a PowerShell command for creating a local user.
func (params UserCreateParams) pwshCommand() string {
	// Base command
	cmd := []string{"New-LocalUser"}

	// Add parameters
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))

	if params.Description != "" {
		cmd = append(cmd, fmt.Sprintf("-Description '%s'", params.Description))
	}

	if params.AccountExpires.Compare(time.Now()) == 1 {
		accountExpires := params.AccountExpires.Format(time.DateTime)
		cmd = append(cmd, fmt.Sprintf("-AccountExpires $(Get-Date '%s')", accountExpires))
	} else {
		cmd = append(cmd, "-AccountNeverExpires")
	}

	if params.Enabled {
		cmd = append(cmd, "-Disabled:$false")
	} else {
		cmd = append(cmd, "-Disabled")
	}

	if params.FullName != "" {
		cmd = append(cmd, fmt.Sprintf("-FullName '%s'", params.FullName))
	}

	if params.Password != "" {
		cmd = append(cmd, fmt.Sprintf("-Password $(ConvertTo-SecureString -String '%s' -AsPlainText -Force)", params.Password))
		cmd = append(cmd, fmt.Sprintf("-PasswordNeverExpires:$%t", params.PasswordNeverExpires))
	} else {
		cmd = append(cmd, "-NoPassword")
	}

	if params.UserMayChangePassword {
		cmd = append(cmd, "-UserMayNotChangePassword:$false")
	} else {
		cmd = append(cmd, "-UserMayNotChangePassword")
	}

	cmd = append(cmd, "| ConvertTo-Json -Compress")
	return strings.Join(cmd, " ")
}

// UserCreate creates a local user and returns a User object.
func (c *LocalClient) UserCreate(ctx context.Context, params UserCreateParams) (User, error) {
	var u User

	// Assert needed parameters
	if params.Name == "" {
		return u, fmt.Errorf("windows.local.UserCreate: user parameter 'Name' must be set")
	}

	// Run command
	if err := localRun(ctx, c, params.pwshCommand(), &u); err != nil {
		return u, fmt.Errorf("windows.local.UserCreate: %s", err)
	}

	return u, nil
}

// UserUpdateParams represents parameters for the UserUpdate function.
type UserUpdateParams struct {
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

	// Indicates whether the account is enabled.
	Enabled bool

	// Specifies the full name for the user account.
	// The full name differs from the user name of the user account.
	FullName string

	// Specifies a password for the user account.
	Password string

	// Indicates whether the new user's password expires.
	PasswordNeverExpires bool

	// Indicates that the user can change the password on the user account.
	UserMayChangePassword bool
}

// pwshCommand returns a PowerShell command for updating a local user.
func (params UserUpdateParams) pwshCommand() string {
	// Base commands
	cmd1 := []string{"Set-LocalUser"}
	cmd2 := []string{}

	if params.Enabled {
		cmd2 = append(cmd2, "Enable-LocalUser")
	} else {
		cmd2 = append(cmd2, "Disable-LocalUser")
	}

	// Add parameters
	// Prefer SID over Name to identify group
	if params.SID != "" {
		cmd1 = append(cmd1, fmt.Sprintf("-SID %s", params.SID))
		cmd2 = append(cmd2, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmd1 = append(cmd1, fmt.Sprintf("-Name '%s'", params.Name))
		cmd2 = append(cmd2, fmt.Sprintf("-Name '%s'", params.Name))
	}

	if params.AccountExpires.Compare(time.Now()) == 1 {
		accountExpires := params.AccountExpires.Format(time.DateTime)
		cmd1 = append(cmd1, fmt.Sprintf("-AccountExpires $(Get-Date '%s')", accountExpires))
	} else {
		cmd1 = append(cmd1, "-AccountNeverExpires")
	}

	// Always set Description and FullName to allow removal of these parameters
	cmd1 = append(cmd1, fmt.Sprintf("-Description '%s'", params.Description))
	cmd1 = append(cmd1, fmt.Sprintf("-FullName '%s'", params.FullName))

	if params.Password != "" {
		cmd1 = append(cmd1, fmt.Sprintf("-Password $(ConvertTo-SecureString -String '%s' -AsPlainText -Force)", params.Password))
	}

	cmd1 = append(cmd1, fmt.Sprintf("-PasswordNeverExpires:$%t", params.PasswordNeverExpires))
	cmd1 = append(cmd1, fmt.Sprintf("-UserMayChangePassword:$%t", params.UserMayChangePassword))

	// Append second command with a semicolon
	cmd1 = append(cmd1, fmt.Sprintf(";%s", strings.Join(cmd2, " ")))
	return strings.Join(cmd1, " ")
}

// UserUpdate updates a local user.
func (c *LocalClient) UserUpdate(ctx context.Context, params UserUpdateParams) error {
	var u User

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.UserUpdate: user parameter 'Name' or 'SID' must be set")
	}

	// Run command
	if err := localRun(ctx, c, params.pwshCommand(), &u); err != nil {
		return fmt.Errorf("windows.local.UserUpdate: %s", err)
	}

	return nil
}

// UserDeleteParams represents parameters for the UserDelete function.
type UserDeleteParams struct {
	// Specifies the user name of the user account.
	Name string

	// Specifies a security ID (SID) of user account.
	SID string
}

// pwshCommand returns a PowerShell command for deleting a local user.
func (params UserDeleteParams) pwshCommand() string {
	// Base command
	cmd := []string{"Remove-LocalUser"}

	// Add parameters
	// Prefer SID over Name to identifiy group
	if params.SID != "" {
		cmd = append(cmd, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	}

	return strings.Join(cmd, " ")
}

// UserDelete removes a local user by SID or Name.
func (c *LocalClient) UserDelete(ctx context.Context, params UserDeleteParams) error {
	var u User

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.UserDelete: user parameter 'Name' or 'SID' must be set")
	}

	// Run command
	if err := localRun(ctx, c, params.pwshCommand(), &u); err != nil {
		return fmt.Errorf("windows.local.UserDelete: %s", err)
	}

	return nil
}
