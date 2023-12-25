package local

import (
	"context"
	"encoding/json"
	"errors"
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

// GroupParams represents parameters for interacting with local users, including creation, updating, and deletion.
type UserParams struct {
	Name        string
	Description string
	SID         string
}

// userType is an interface for either a single User or a slice of User.
type userType interface {
	User | []User
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
	if err := userRun[User](ctx, c, cmd, &u); err != nil {
		return u, fmt.Errorf("windows.local.UserRead: %s", err)
	}

	return u, nil
}

// userRun runs a PowerShell command against a Windows system, handles the command results,
// and unmarshals the output into a User object or a slice of User objects.
func userRun[T userType](ctx context.Context, c *LocalClient, cmd string, u *T) error {

	// Run the command
	result, err := c.Connection.Run(ctx, cmd)
	if err != nil {
		return err
	}

	// Handle stderr
	if result.StdErr != "" {
		stderr, err := c.parser.DecodeCLIXML(result.StdErr)
		if err != nil {
			return err
		}

		return errors.New(stderr)
	}

	if result.StdOut == "" {
		return nil
	}

	// Unmarshal stdout
	if err = json.Unmarshal([]byte(result.StdOut), &u); err != nil {
		return err
	}

	return nil
}
