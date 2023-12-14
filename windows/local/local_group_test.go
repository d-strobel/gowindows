package local

import (
	"context"
	"testing"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
	"github.com/d-strobel/gowindows/windows/local/fixtures"
	"github.com/stretchr/testify/suite"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
)

// Setup unit test suite
type GroupUnitTestSuite struct {
	suite.Suite
	expectedUsersGroup Group
	expectedGroupList  []Group
}

func (suite *GroupUnitTestSuite) SetupTest() {
	suite.expectedUsersGroup = Group{
		Name:        "Users",
		Description: "Users are prevented from making accidental or intentional system-wide changes and can run most applications",
		SID: SID{
			Value: "S-1-5-32-545",
		},
	}

	suite.expectedGroupList = []Group{
		{
			Name:        "Administrators",
			Description: "Administrators have complete and unrestricted access to the computer/domain",
			SID: SID{
				Value: "S-1-5-32-544",
			},
		},
		{
			Name:        "Users",
			Description: "Users are prevented from making accidental or intentional system-wide changes and can run most applications",
			SID: SID{
				Value: "S-1-5-32-545",
			},
		},
	}
}

func TestGroupUnitTestSuite(t *testing.T) {
	suite.Run(t, &GroupUnitTestSuite{})
}

// Unit test functions
func (suite *GroupUnitTestSuite) TestGroupRun() {
	suite.T().Parallel()

	suite.Run("Should return the user Group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser: parser.Parser{
				DecodeCLIXML: func(xmlErr string) (string, error) {
					return xmlErr, nil
				},
			},
		}
		mockConn.On("Run", ctx, "Get-LocalGroup -Name Users | ConvertTo-Json").Return(connection.CMDResult{
			StdOut: fixtures.UsersGroup,
		}, nil)
		var g Group
		err := groupRun[Group](ctx, c, "Get-LocalGroup -Name Users | ConvertTo-Json", &g)
		suite.NoError(err)
		suite.Equal(suite.expectedUsersGroup, g)
	})

	suite.Run("Should return a slice of Groups", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser: parser.Parser{
				DecodeCLIXML: func(xmlErr string) (string, error) {
					return xmlErr, nil
				},
			},
		}
		mockConn.On("Run", ctx, "Get-LocalGroup | ConvertTo-Json").Return(connection.CMDResult{
			StdOut: fixtures.GroupList,
		}, nil)
		var g []Group
		err := groupRun[[]Group](ctx, c, "Get-LocalGroup | ConvertTo-Json", &g)
		suite.NoError(err)
		suite.Equal(suite.expectedGroupList, g)
	})
}
