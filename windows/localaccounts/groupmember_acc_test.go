package localaccounts_test

import (
	"context"

	"github.com/d-strobel/gowindows/windows/localaccounts"
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
		params := localaccounts.GroupMemberReadParams{
			Name:   "Administrators",
			Member: "Administrator",
		}
		u, err := c.GroupMemberRead(ctx, params)
		suite.Require().NoError(err)
		suite.Equal(localaccounts.GroupMember{
			Name: "WIN2022SC\\Administrator",
			SID: localaccounts.SID{
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
		u, err := c.GroupMemberList(ctx, localaccounts.GroupMemberListParams{
			Name: "Administrators",
		})
		suite.Require().NoError(err)
		suite.Contains(u, localaccounts.GroupMember{
			Name: "WIN2022SC\\Administrator",
			SID: localaccounts.SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-500",
			},
			ObjectClass: "User",
		})
		suite.Contains(u, localaccounts.GroupMember{
			Name: "WIN2022SC\\vagrant",
			SID: localaccounts.SID{
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
		params := localaccounts.GroupMemberCreateParams{
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
		params := localaccounts.GroupMemberDeleteParams{
			Name:   "Administrators",
			Member: groupMemberTestCases[i],
		}
		err := c.GroupMemberDelete(ctx, params)
		suite.NoError(err)
	}
}
