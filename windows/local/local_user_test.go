package local

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
	"github.com/stretchr/testify/suite"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	mockParser "github.com/d-strobel/gowindows/parser/mocks"
)

// Unit test suite for all User functions
type UserUnitTestSuite struct {
	suite.Suite
	// Fixtures
	adminUser         string
	expectedAdminUser User
	userList          string
	expectedUserList  []User
	testUser          string
	expectedTestUser  User
}

func (suite *UserUnitTestSuite) SetupSuite() {
	// Fixtures
	suite.adminUser = `{"AccountExpires":null,"Description":"Built-in account for administering the computer/domain","Enabled":true,"FullName":"","PasswordChangeableDate":"\/Date(1701379505092)\/","PasswordExpires":null,"UserMayChangePassword":true,"PasswordRequired":true,"PasswordLastSet":"\/Date(1701379505092)\/","LastLogon":null,"Name":"Administrator","SID":{"BinaryLength":28,"AccountDomainSid":{"BinaryLength":24,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138"},"Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"}`
	suite.userList = `[{"AccountExpires":null,"Description":"Built-in account for administering the computer/domain","Enabled":true,"FullName":"","PasswordChangeableDate":"\/Date(1701379505092)\/","PasswordExpires":null,"UserMayChangePassword":true,"PasswordRequired":true,"PasswordLastSet":"\/Date(1701379505092)\/","LastLogon":null,"Name":"Administrator","SID":{"BinaryLength":28,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"},{"AccountExpires":null,"Description":"Built-in account for guest access to the computer/domain","Enabled":false,"FullName":"","PasswordChangeableDate":null,"PasswordExpires":null,"UserMayChangePassword":false,"PasswordRequired":false,"PasswordLastSet":null,"LastLogon":null,"Name":"Guest","SID":{"BinaryLength":28,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138-501"},"PrincipalSource":1,"ObjectClass":"User"}]`
	suite.testUser = `{"AccountExpires":null,"Description":"Test user","Enabled":true,"FullName":"","PasswordChangeableDate":"\/Date(1701379505092)\/","PasswordExpires":null,"UserMayChangePassword":true,"PasswordRequired":true,"PasswordLastSet":"\/Date(1701379505092)\/","LastLogon":null,"Name":"Test-User","SID":{"BinaryLength":28,"AccountDomainSid":{"BinaryLength":24,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405139","Value":"S-1-5-21-153895498-367353507-3704405139"},"Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"}`

	suite.expectedAdminUser = User{
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
		SID: SID{
			Value: "S-1-5-21-153895498-367353507-3704405138-500",
		},
	}
	suite.expectedUserList = []User{
		{
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
			SID: SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-500",
			},
		},
		{
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
			SID: SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-501",
			},
		},
	}
	suite.expectedTestUser = User{
		AccountExpires:         parser.WinTime{},
		Description:            "Test user",
		Enabled:                true,
		FullName:               "",
		PasswordChangeableDate: parser.WinTime{Time: time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)},
		PasswordExpires:        parser.WinTime{},
		UserMayChangePassword:  true,
		PasswordRequired:       true,
		PasswordLastSet:        parser.WinTime{Time: time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)},
		LastLogon:              parser.WinTime{},
		Name:                   "Test-User",
		SID: SID{
			Value: "S-1-5-21-153895498-367353507-3704405138",
		},
	}
}

func TestUserUnitTestSuite(t *testing.T) {
	suite.Run(t, &UserUnitTestSuite{})
}

func (suite *UserUnitTestSuite) TestUserRun() {

	suite.Run("should return the administrator user", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalUser -Name 'Administrator' | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: suite.adminUser,
		}, nil)
		var u User
		err := userRun[User](ctx, c, expectedCMD, &u)
		suite.NoError(err)
		suite.Equal(suite.expectedAdminUser, u)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should return a slice of user", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalUser | ConvertTo-Json"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: suite.userList,
		}, nil)
		var u []User
		err := userRun[[]User](ctx, c, expectedCMD, &u)
		suite.NoError(err)
		suite.Equal(suite.expectedUserList, u)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
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
		expectedCMD := "Remove-LocalUser -Name 'Test'"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: "",
		}, nil)
		var u User
		var expectedUser User
		err := userRun[User](ctx, c, expectedCMD, &u)
		suite.NoError(err)
		suite.Equal(expectedUser, u)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
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
		expectedCMD := "Remove-LocalUser -Name 'Test'"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		var u User
		expectedErr := errors.New("test-error")
		err := userRun[User](ctx, c, expectedCMD, &u)
		suite.Error(err)
		suite.Equal(expectedErr, err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
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
		expectedCMD := "Get-LocalUser -name 'Administratr'"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdErr: "clixml-error",
		}, nil)
		mockParser.On("DecodeCLIXML", "clixml-error").Return("powershell-error", nil)
		var u User
		expectedErr := errors.New("powershell-error")
		err := userRun[User](ctx, c, expectedCMD, &u)
		suite.Error(err)
		suite.Equal(expectedErr, err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
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
		expectedCMD := "Get-LocalUser -name 'Administrar'"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdErr: "incorrect-clixml-error",
		}, nil)
		mockParser.On("DecodeCLIXML", "incorrect-clixml-error").Return("", errors.New("parser-error"))
		var u User
		expectedErr := errors.New("parser-error")
		err := userRun[User](ctx, c, expectedCMD, &u)
		suite.Error(err)
		suite.Equal(expectedErr, err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
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
		expectedCMD := "Get-LocalUser -name 'Administrator'"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: suite.userList,
		}, nil)
		var u User
		err := userRun[User](ctx, c, expectedCMD, &u)
		suite.Error(err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}
