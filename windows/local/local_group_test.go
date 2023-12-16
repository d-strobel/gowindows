package local

import (
	"context"
	"testing"

	"github.com/d-strobel/gowindows/connection"
	"github.com/stretchr/testify/suite"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	mockParser "github.com/d-strobel/gowindows/parser/mocks"
)

// Unit test suite for all Group functions
type GroupUnitTestSuite struct {
	suite.Suite
	// Fixtures
	usersGroup         string
	expectedUsersGroup Group
	groupList          string
	expectedGroupList  []Group
}

func (suite *GroupUnitTestSuite) SetupTest() {
	// Fixtures
	suite.usersGroup = `{
    "Description":  "Users are prevented from making accidental or intentional system-wide changes and can run most applications",
    "Name":  "Users",
    "SID":  {
                "BinaryLength":  16,
                "AccountDomainSid":  null,
                "Value":  "S-1-5-32-545"
            },
    "PrincipalSource":  1,
    "ObjectClass":  "Group"
}`

	suite.groupList = `[
    {
        "Description":  "Administrators have complete and unrestricted access to the computer/domain",
        "Name":  "Administrators",
        "SID":  {
                    "BinaryLength":  16,
                    "AccountDomainSid":  null,
                    "Value":  "S-1-5-32-544"
                },
        "PrincipalSource":  1,
        "ObjectClass":  "Group"
    },
    {
        "Description":  "Users are prevented from making accidental or intentional system-wide changes and can run most applications",
        "Name":  "Users",
        "SID":  {
                    "BinaryLength":  16,
                    "AccountDomainSid":  null,
                    "Value":  "S-1-5-32-545"
                },
        "PrincipalSource":  1,
        "ObjectClass":  "Group"
    }
]`

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

func (suite *GroupUnitTestSuite) TestGroupRun() {
	suite.T().Parallel()

	suite.Run("should return the user group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		mockConn.On("Run", ctx, "Get-LocalGroup -Name Users | ConvertTo-Json").Return(connection.CMDResult{
			StdOut: suite.usersGroup,
		}, nil)
		var g Group
		err := groupRun[Group](ctx, c, "Get-LocalGroup -Name Users | ConvertTo-Json", &g)
		suite.NoError(err)
		suite.Equal(suite.expectedUsersGroup, g)
		mockConn.AssertCalled(suite.T(), "Run", ctx, "Get-LocalGroup -Name Users | ConvertTo-Json")
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should return a slice of groups", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		mockConn.On("Run", ctx, "Get-LocalGroup | ConvertTo-Json").Return(connection.CMDResult{
			StdOut: suite.groupList,
		}, nil)
		var g []Group
		err := groupRun[[]Group](ctx, c, "Get-LocalGroup | ConvertTo-Json", &g)
		suite.NoError(err)
		suite.Equal(suite.expectedGroupList, g)
		mockConn.AssertCalled(suite.T(), "Run", ctx, "Get-LocalGroup | ConvertTo-Json")
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}
