package local_test

import (
	"context"
	"time"

	"github.com/d-strobel/gowindows/parser"
	"github.com/d-strobel/gowindows/windows/local"
)

// We insert numbers into the function names to ensure that
// the test functions for each local_* file run in a specific order.
func (suite *LocalAccTestSuite) TestUser1Read() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, c := range suite.clients {
		params := local.UserParams{
			Name: "Administrator",
		}
		u, err := c.UserRead(ctx, params)
		suite.Require().NoError(err)
		suite.Equal(local.User{
			AccountExpires:         parser.WinTime{},
			Description:            "Built-in account for administering the computer/domain",
			Enabled:                true,
			FullName:               "",
			PasswordChangeableDate: parser.WinTime{Time: time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)},
			PasswordExpires:        parser.WinTime{},
			UserMayChangePassword:  true,
			PasswordRequired:       true,
			PasswordLastSet:        parser.WinTime{Time: time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)},
			LastLogon:              parser.WinTime{},
			Name:                   "Administrator",
			SID: local.SID{
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
		suite.Contains(u, local.User{
			AccountExpires:         parser.WinTime{},
			Description:            "Built-in account for administering the computer/domain",
			Enabled:                true,
			FullName:               "",
			PasswordChangeableDate: parser.WinTime{Time: time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)},
			PasswordExpires:        parser.WinTime{},
			UserMayChangePassword:  true,
			PasswordRequired:       true,
			PasswordLastSet:        parser.WinTime{Time: time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)},
			LastLogon:              parser.WinTime{},
			Name:                   "Administrator",
			SID: local.SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-500",
			},
		})
		suite.Contains(u, local.User{
			AccountExpires:         parser.WinTime{},
			Description:            "Built-in account for guest access to the computer/domain",
			Enabled:                false,
			FullName:               "",
			PasswordChangeableDate: parser.WinTime{},
			PasswordExpires:        parser.WinTime{},
			UserMayChangePassword:  false,
			PasswordRequired:       false,
			PasswordLastSet:        parser.WinTime{},
			LastLogon:              parser.WinTime{},
			Name:                   "Guest",
			SID: local.SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-501",
			},
		})
	}
}
