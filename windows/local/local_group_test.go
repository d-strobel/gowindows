package local

import (
	"context"
	"errors"
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
	usersGroup           string
	usersGroupCompressed string
	expectedUsersGroup   Group
	groupList            string
	expectedGroupList    []Group
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

	suite.usersGroupCompressed = `{"Description":"Users are prevented from making accidental or intentional system-wide changes and can run most applications","Name":"Users","SID":{"BinaryLength":16,"AccountDomainSid":null,"Value":"S-1-5-32-545"},"PrincipalSource":1,"ObjectClass":"Group"}`

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

	suite.Run("should return the user group with compressed json", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		mockConn.On("Run", ctx, "Get-LocalGroup -Name Users | ConvertTo-Json -compress").Return(connection.CMDResult{
			StdOut: suite.usersGroupCompressed,
		}, nil)
		var g Group
		err := groupRun[Group](ctx, c, "Get-LocalGroup -Name Users | ConvertTo-Json -compress", &g)
		suite.NoError(err)
		suite.Equal(suite.expectedUsersGroup, g)
		mockConn.AssertCalled(suite.T(), "Run", ctx, "Get-LocalGroup -Name Users | ConvertTo-Json -compress")
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should not error when no stdout is empty string", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		mockConn.On("Run", ctx, "Remove-LocalGroup -Name Test").Return(connection.CMDResult{
			StdOut: "",
		}, nil)
		var g Group
		var expectedGroup Group
		err := groupRun[Group](ctx, c, "Remove-LocalGroup -Name Test", &g)
		suite.NoError(err)
		suite.Equal(expectedGroup, g)
		mockConn.AssertCalled(suite.T(), "Run", ctx, "Remove-LocalGroup -Name Test")
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should error when connection run errors", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		mockConn.On("Run", ctx, "Remove-LocalGroup -Name Test").Return(connection.CMDResult{}, errors.New("test-error"))
		var g Group
		expectedErr := errors.New("test-error")
		err := groupRun[Group](ctx, c, "Remove-LocalGroup -Name Test", &g)
		suite.Error(err)
		suite.Equal(expectedErr, err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, "Remove-LocalGroup -Name Test")
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should return powershell error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		mockConn.On("Run", ctx, "Get-LocalGroup -name Userrs").Return(connection.CMDResult{
			StdErr: "clixml-error",
		}, nil)
		mockParser.On("DecodeCLIXML", "clixml-error").Return("powershell-error", nil)
		var g Group
		expectedErr := errors.New("powershell-error")
		err := groupRun[Group](ctx, c, "Get-LocalGroup -name Userrs", &g)
		suite.Error(err)
		suite.Equal(expectedErr, err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, "Get-LocalGroup -name Userrs")
		mockParser.AssertCalled(suite.T(), "DecodeCLIXML", "clixml-error")
	})

	suite.Run("should return an error from parser", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		mockConn.On("Run", ctx, "Get-LocalGroup -name Userrs").Return(connection.CMDResult{
			StdErr: "incorrect-clixml-error",
		}, nil)
		mockParser.On("DecodeCLIXML", "incorrect-clixml-error").Return("", errors.New("parser-error"))
		var g Group
		expectedErr := errors.New("parser-error")
		err := groupRun[Group](ctx, c, "Get-LocalGroup -name Userrs", &g)
		suite.Error(err)
		suite.Equal(expectedErr, err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, "Get-LocalGroup -name Userrs")
		mockParser.AssertCalled(suite.T(), "DecodeCLIXML", "incorrect-clixml-error")
	})

	suite.Run("should return error from json unmarshal with incorrect json", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		mockConn.On("Run", ctx, "Get-LocalGroup -name Users").Return(connection.CMDResult{
			StdOut: suite.groupList,
		}, nil)
		var g Group
		err := groupRun[Group](ctx, c, "Get-LocalGroup -name Users", &g)
		suite.Error(err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, "Get-LocalGroup -name Users")
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}
