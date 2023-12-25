package local

import (
	"context"
	"errors"
	"time"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	mockParser "github.com/d-strobel/gowindows/parser/mocks"
)

// Fixtures
const (
	adminUser = `{"AccountExpires":null,"Description":"Built-in account for administering the computer/domain","Enabled":true,"FullName":"","PasswordChangeableDate":"\/Date(1701379505092)\/","PasswordExpires":null,"UserMayChangePassword":true,"PasswordRequired":true,"PasswordLastSet":"\/Date(1701379505092)\/","LastLogon":null,"Name":"Administrator","SID":{"BinaryLength":28,"AccountDomainSid":{"BinaryLength":24,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138"},"Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"}`
	userList  = `[{"AccountExpires":null,"Description":"Built-in account for administering the computer/domain","Enabled":true,"FullName":"","PasswordChangeableDate":"\/Date(1701379505092)\/","PasswordExpires":null,"UserMayChangePassword":true,"PasswordRequired":true,"PasswordLastSet":"\/Date(1701379505092)\/","LastLogon":null,"Name":"Administrator","SID":{"BinaryLength":28,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"},{"AccountExpires":null,"Description":"Built-in account for guest access to the computer/domain","Enabled":false,"FullName":"","PasswordChangeableDate":null,"PasswordExpires":null,"UserMayChangePassword":false,"PasswordRequired":false,"PasswordLastSet":null,"LastLogon":null,"Name":"Guest","SID":{"BinaryLength":28,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138-501"},"PrincipalSource":1,"ObjectClass":"User"}]`
	testUser  = `{"AccountExpires":null,"Description":"Test user","Enabled":true,"FullName":"","PasswordChangeableDate":"\/Date(1701379505092)\/","PasswordExpires":null,"UserMayChangePassword":true,"PasswordRequired":true,"PasswordLastSet":"\/Date(1701379505092)\/","LastLogon":null,"Name":"Test-User","SID":{"BinaryLength":28,"AccountDomainSid":{"BinaryLength":24,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405139","Value":"S-1-5-21-153895498-367353507-3704405139"},"Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"}`
)

var (
	expectedAdminUser = User{
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
	expectedUserList = []User{
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
	expectedTestUser = User{
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
)

func (suite *LocalUnitTestSuite) TestUserRead() {

	suite.Run("should return the correct group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup -Name 'Users' | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: usersGroup,
		}, nil)
		actualUsersGroup, err := c.GroupRead(ctx, GroupParams{Name: "Users"})
		suite.Require().NoError(err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		suite.Equal(expectedUsersGroup, actualUsersGroup)
	})

	suite.Run("should run the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupParams
			expectedCMD     string
		}{
			{
				"assert users group by name",
				GroupParams{Name: "Users"},
				"Get-LocalGroup -Name 'Users' | ConvertTo-Json -Compress",
			},
			{
				"assert users group by sid",
				GroupParams{SID: "123456789"},
				"Get-LocalGroup -SID 123456789 | ConvertTo-Json -Compress",
			},
			{
				"assert users group by name and sid",
				GroupParams{Name: "Users", SID: "123456789"},
				"Get-LocalGroup -SID 123456789 | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mockConn := mockConnection.NewMockConnectionInterface(suite.T())
			mockParser := mockParser.NewMockParserInterface(suite.T())
			c := &LocalClient{
				Connection: mockConn,
				parser:     mockParser,
			}
			mockConn.On("Run", ctx, tc.expectedCMD).Return(connection.CMDResult{}, nil)
			_, err := c.GroupRead(ctx, tc.inputParameters)
			suite.Require().NoError(err)
			mockConn.AssertCalled(suite.T(), "Run", ctx, tc.expectedCMD)
			mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		}
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters GroupParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				GroupParams{},
				"windows.local.GroupRead: group parameter 'Name' or 'SID' must be set",
			},
			{
				"assert error with just the description parameter",
				GroupParams{Description: "test"},
				"windows.local.GroupRead: group parameter 'Name' or 'SID' must be set",
			},
			{
				"assert error when name contains wildcard",
				GroupParams{Name: "Remote*"},
				"windows.local.GroupRead: group parameter 'Name' does not allow wildcards",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mockConn := mockConnection.NewMockConnectionInterface(suite.T())
			mockParser := mockParser.NewMockParserInterface(suite.T())
			c := &LocalClient{
				Connection: mockConn,
				parser:     mockParser,
			}
			_, err := c.GroupRead(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
			mockConn.AssertNotCalled(suite.T(), "Run")
			mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		}
	})

	suite.Run("should return error if run fails", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup -Name 'Users' | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		_, err := c.GroupRead(ctx, GroupParams{Name: "Users"})
		suite.EqualError(err, "windows.local.GroupRead: test-error")
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}
