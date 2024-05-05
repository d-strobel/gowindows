package accounts_test

import (
	"context"
	"fmt"
	"time"

	"github.com/d-strobel/gowindows/parsing"
	"github.com/d-strobel/gowindows/windows/local/accounts"
)

// We insert numbers into the function names to ensure that
// the test functions for each local_* file run in a specific order.
func (suite *LocalAccTestSuite) TestUser1Read() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, c := range suite.clients {
		params := accounts.UserReadParams{
			Name: "Administrator",
		}
		u, err := c.UserRead(ctx, params)
		suite.Require().NoError(err)
		suite.Equal(accounts.User{
			AccountExpires:         parsing.DotnetTime{},
			Description:            "Built-in account for administering the computer/domain",
			Enabled:                true,
			FullName:               "",
			PasswordChangeableDate: parsing.DotnetTime(time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)),
			PasswordExpires:        parsing.DotnetTime{},
			UserMayChangePassword:  true,
			PasswordRequired:       true,
			PasswordLastSet:        parsing.DotnetTime(time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)),
			LastLogon:              parsing.DotnetTime{},
			Name:                   "Administrator",
			SID: accounts.SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-500",
			},
		}, u)
	}
}

func (suite *LocalAccTestSuite) TestUser2List() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, c := range suite.clients {
		u, err := c.UserList(ctx)
		suite.Require().NoError(err)
		suite.Contains(u, accounts.User{
			AccountExpires:         parsing.DotnetTime{},
			Description:            "Built-in account for administering the computer/domain",
			Enabled:                true,
			FullName:               "",
			PasswordChangeableDate: parsing.DotnetTime(time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)),
			PasswordExpires:        parsing.DotnetTime{},
			UserMayChangePassword:  true,
			PasswordRequired:       true,
			PasswordLastSet:        parsing.DotnetTime(time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)),
			LastLogon:              parsing.DotnetTime{},
			Name:                   "Administrator",
			SID: accounts.SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-500",
			},
		})
		suite.Contains(u, accounts.User{
			AccountExpires:         parsing.DotnetTime{},
			Description:            "Built-in account for guest access to the computer/domain",
			Enabled:                false,
			FullName:               "",
			PasswordChangeableDate: parsing.DotnetTime{},
			PasswordExpires:        parsing.DotnetTime{},
			UserMayChangePassword:  false,
			PasswordRequired:       false,
			PasswordLastSet:        parsing.DotnetTime{},
			LastLogon:              parsing.DotnetTime{},
			Name:                   "Guest",
			SID: accounts.SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-501",
			},
		})
	}
}

func (suite *LocalAccTestSuite) TestUser3Create() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, c := range suite.clients {
		params := accounts.UserCreateParams{
			Name:                 fmt.Sprintf("Test-User-%d", i),
			Description:          "This is a test user",
			FullName:             fmt.Sprintf("Full-Test-User-%d", i),
			Password:             "Start123!!!",
			PasswordNeverExpires: true,
			AccountExpires:       time.Date(2025, time.November, 10, 16, 0, 0, 0, time.UTC),
			Enabled:              true,
		}
		g, err := c.UserCreate(ctx, params)
		suite.NoError(err)
		suite.Equal(accounts.User{Name: fmt.Sprintf("Test-User-%d", i)}.Name, g.Name)
		suite.Equal(accounts.User{Description: "This is a test user"}.Description, g.Description)
		suite.Equal(accounts.User{FullName: fmt.Sprintf("Full-Test-User-%d", i)}.FullName, g.FullName)
		suite.Equal(accounts.User{PasswordExpires: parsing.DotnetTime{}}.PasswordExpires, g.PasswordExpires)
		suite.Equal(accounts.User{AccountExpires: parsing.DotnetTime(time.Date(2025, time.November, 10, 16, 0, 0, 0, time.UTC))}.AccountExpires, g.AccountExpires)
		suite.Equal(accounts.User{UserMayChangePassword: false}.UserMayChangePassword, g.UserMayChangePassword)
		suite.Equal(accounts.User{Enabled: true}.Enabled, g.Enabled)
	}
}

func (suite *LocalAccTestSuite) TestUser4Update() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, c := range suite.clients {
		params := accounts.UserUpdateParams{
			Name:           fmt.Sprintf("Test-User-%d", i),
			Description:    "Updated - This is a test user",
			FullName:       fmt.Sprintf("Updated-Full-Test-User-%d", i),
			Password:       "Start123!!!3",
			AccountExpires: time.Date(2026, time.November, 10, 16, 0, 0, 0, time.UTC),
			Enabled:        false,
		}
		err := c.UserUpdate(ctx, params)
		suite.NoError(err)
	}
}

func (suite *LocalAccTestSuite) TestUser5Delete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, c := range suite.clients {
		params := accounts.UserDeleteParams{
			Name: fmt.Sprintf("Test-User-%d", i),
		}
		err := c.UserDelete(ctx, params)
		suite.NoError(err)
	}
}
