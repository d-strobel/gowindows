package accounts_test

import (
	"context"

	"github.com/d-strobel/gowindows/windows/local/accounts"
)

var groupMemberTestCases = []string{
	"Guest",
	"DefaultAccount",
}

// We insert numbers into the function names to ensure that
// the test functions for each local_* file run in a specific order.
func (suite *LocalAccTestSuite) TestGroupMember1Read() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, c := range suite.clients {
		params := accounts.GroupMemberReadParams{
			Name:   "Administrators",
			Member: "Administrator",
		}
		u, err := c.GroupMemberRead(ctx, params)
		suite.Require().NoError(err)
		suite.Equal(accounts.GroupMember{
			Name: "WIN2022SC\\Administrator",
			SID: accounts.SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-500",
			},
			ObjectClass: "User",
		}, u)
	}
}

func (suite *LocalAccTestSuite) TestGroupMember2List() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, c := range suite.clients {
		u, err := c.GroupMemberList(ctx, accounts.GroupMemberListParams{
			Name: "Administrators",
		})
		suite.Require().NoError(err)
		suite.Contains(u, accounts.GroupMember{
			Name: "WIN2022SC\\Administrator",
			SID: accounts.SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-500",
			},
			ObjectClass: "User",
		})
		suite.Contains(u, accounts.GroupMember{
			Name: "WIN2022SC\\vagrant",
			SID: accounts.SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-1000",
			},
			ObjectClass: "User",
		})
	}
}

func (suite *LocalAccTestSuite) TestGroupMember3Create() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, c := range suite.clients {
		params := accounts.GroupMemberCreateParams{
			Name:   "Administrators",
			Member: groupMemberTestCases[i],
		}
		err := c.GroupMemberCreate(ctx, params)
		suite.NoError(err)
	}
}

func (suite *LocalAccTestSuite) TestGroupMember4Delete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, c := range suite.clients {
		params := accounts.GroupMemberDeleteParams{
			Name:   "Administrators",
			Member: groupMemberTestCases[i],
		}
		err := c.GroupMemberDelete(ctx, params)
		suite.NoError(err)
	}
}
