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

// UserReadParams represents parameters for the UserRead function.
type UserReadParams struct {
	// Specifies the user name of the user account.
	Name string

	// Specifies a security ID (SID) of user account.
	SID string
}

// UserRead gets a local user by SID or Name and returns a User object.
func (c *LocalClient) UserRead(ctx context.Context, params UserReadParams) (User, error) {

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

// UserCreate creates a local user and returns a User object.
func (c *LocalClient) UserCreate(ctx context.Context, params UserCreateParams) (User, error) {

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

	if params.Enabled {
		cmds = append(cmds, "-Disabled:$false")
	} else {
		cmds = append(cmds, "-Disabled")
	}

	if params.FullName != "" {
		cmds = append(cmds, fmt.Sprintf("-FullName '%s'", params.FullName))
	}

	if params.Password != "" {
		cmds = append(cmds, fmt.Sprintf("-Password $(ConvertTo-SecureString -String '%s' -AsPlainText -Force)", params.Password))
		cmds = append(cmds, fmt.Sprintf("-PasswordNeverExpires:$%t", params.PasswordNeverExpires))
	} else {
		cmds = append(cmds, "-NoPassword")
	}

	if params.UserMayChangePassword {
		cmds = append(cmds, "-UserMayNotChangePassword:$false")
	} else {
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

// UserUpdate updates a local user.
func (c *LocalClient) UserUpdate(ctx context.Context, params UserUpdateParams) error {

	// Satisfy localType interface
	var u User

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.UserUpdate: user parameter 'Name' or 'SID' must be set")
	}

	// Base command
	cmds := []string{"Set-LocalUser"}
	cmds2 := []string{}

	if params.Enabled {
		cmds2 = append(cmds2, "Enable-LocalUser")
	} else {
		cmds2 = append(cmds2, "Disable-LocalUser")
	}

	// Add parameters
	// Prefer SID over Name to identify group
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
		cmds2 = append(cmds2, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
		cmds2 = append(cmds2, fmt.Sprintf("-Name '%s'", params.Name))
	}

	if params.AccountExpires.Compare(time.Now()) == 1 {
		accountExpires := params.AccountExpires.Format(time.DateTime)
		cmds = append(cmds, fmt.Sprintf("-AccountExpires $(Get-Date '%s')", accountExpires))
	} else {
		cmds = append(cmds, "-AccountNeverExpires")
	}

	// Always set Description and FullName to allow removal of these parameters
	cmds = append(cmds, fmt.Sprintf("-Description '%s'", params.Description))
	cmds = append(cmds, fmt.Sprintf("-FullName '%s'", params.FullName))

	if params.Password != "" {
		cmds = append(cmds, fmt.Sprintf("-Password $(ConvertTo-SecureString -String '%s' -AsPlainText -Force)", params.Password))
	}

	cmds = append(cmds, fmt.Sprintf("-PasswordNeverExpires:$%t", params.PasswordNeverExpires))
	cmds = append(cmds, fmt.Sprintf("-UserMayChangePassword:$%t", params.UserMayChangePassword))

	// Append second command with a semicolon
	cmds = append(cmds, fmt.Sprintf(";%s", strings.Join(cmds2, " ")))
	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[User](ctx, c, cmd, &u); err != nil {
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

// UserDelete removes a local user by SID or Name.
func (c *LocalClient) UserDelete(ctx context.Context, params UserDeleteParams) error {

	// Satisfy localType interface
	var u User

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return fmt.Errorf("windows.local.UserDelete: user parameter 'Name' or 'SID' must be set")
	}

	// Base command
	cmds := []string{"Remove-LocalUser"}

	// Add parameters
	// Prefer SID over Name to identifiy group
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmd := strings.Join(cmds, " ")

	// Run command
	if err := localRun[User](ctx, c, cmd, &u); err != nil {
		return fmt.Errorf("windows.local.UserDelete: %s", err)
	}

	return nil
}
