package accounts

import (
	"context"
	"errors"
	"time"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parsing"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
)

// Fixtures
const (
	adminUser = `{"AccountExpires":null,"Description":"Built-in account for administering the computer/domain","Enabled":true,"FullName":"","PasswordChangeableDate":"\/Date(1701379505092)\/","PasswordExpires":null,"UserMayChangePassword":true,"PasswordRequired":true,"PasswordLastSet":"\/Date(1701379505092)\/","LastLogon":null,"Name":"Administrator","SID":{"BinaryLength":28,"AccountDomainSid":{"BinaryLength":24,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138"},"Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"}`
	userList  = `[{"AccountExpires":null,"Description":"Built-in account for administering the computer/domain","Enabled":true,"FullName":"","PasswordChangeableDate":"\/Date(1701379505092)\/","PasswordExpires":null,"UserMayChangePassword":true,"PasswordRequired":true,"PasswordLastSet":"\/Date(1701379505092)\/","LastLogon":null,"Name":"Administrator","SID":{"BinaryLength":28,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"},{"AccountExpires":null,"Description":"Built-in account for guest access to the computer/domain","Enabled":false,"FullName":"","PasswordChangeableDate":null,"PasswordExpires":null,"UserMayChangePassword":false,"PasswordRequired":false,"PasswordLastSet":null,"LastLogon":null,"Name":"Guest","SID":{"BinaryLength":28,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138-501"},"PrincipalSource":1,"ObjectClass":"User"}]`
	testUser  = `{"AccountExpires":"\/Date(1762790400000)\/","Description":"This is a test user","Enabled":true,"FullName":"Full-Test-User","PasswordChangeableDate":null,"PasswordExpires":null,"UserMayChangePassword":true,"PasswordRequired":false,"PasswordLastSet":null,"LastLogon":null,"Name":"Test-User","SID":{"BinaryLength":28,"AccountDomainSid":{"BinaryLength":24,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138"},"Value":"S-1-5-21-153895498-367353507-3704405138-1016"},"PrincipalSource":1,"ObjectClass":"User"}`
)

var (
	expectedAdminUser = User{
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
		SID: SID{
			Value: "S-1-5-21-153895498-367353507-3704405138-500",
		},
	}
	expectedUserList = []User{
		{
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
			SID: SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-500",
			},
		},
		{
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
			SID: SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-501",
			},
		},
	}
	expectedTestUser = User{
		AccountExpires:         parsing.DotnetTime(time.Date(2025, time.November, 10, 16, 0, 0, 0, time.UTC)),
		Description:            "This is a test user",
		Enabled:                true,
		FullName:               "Full-Test-User",
		PasswordChangeableDate: parsing.DotnetTime{},
		PasswordExpires:        parsing.DotnetTime{},
		UserMayChangePassword:  true,
		PasswordRequired:       false,
		PasswordLastSet:        parsing.DotnetTime{},
		LastLogon:              parsing.DotnetTime{},
		Name:                   "Test-User",
		SID: SID{
			Value: "S-1-5-21-153895498-367353507-3704405138-1016",
		},
	}
)

// Test GroupRead related methods.
func (suite *LocalUnitTestSuite) TestUserReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters UserReadParams
			expectedCmd     string
		}{
			{
				"assert user by name",
				UserReadParams{Name: "Administrator"},
				"Get-LocalUser -Name 'Administrator' | ConvertTo-Json -Compress",
			},
			{
				"assert users by sid",
				UserReadParams{SID: "123456789"},
				"Get-LocalUser -SID 123456789 | ConvertTo-Json -Compress",
			},
			{
				"assert users by name and sid",
				UserReadParams{Name: "Users", SID: "123456789"},
				"Get-LocalUser -SID 123456789 | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestUserRead() {
	suite.Run("should return the correct user", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-LocalUser -Name 'Administrator' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: adminUser}, nil)
		actualAdminUser, err := c.UserRead(ctx, UserReadParams{Name: "Administrator"})
		suite.NoError(err)
		suite.Equal(expectedAdminUser, actualAdminUser)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters UserReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				UserReadParams{},
				"windows.local.accounts.UserRead: user parameter 'Name' or 'SID' must be set",
			},
			{
				"assert error when name contains wildcard",
				UserReadParams{Name: "Remote*"},
				"windows.local.accounts.UserRead: user parameter 'Name' does not allow wildcards",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mockConn := mockConnection.NewMockConnection(suite.T())
			c := &Client{
				Connection:      mockConn,
				decodeCliXmlErr: func(s string) (string, error) { return s, nil },
			}
			_, err := c.UserRead(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})

	suite.Run("should return error if run fails", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-LocalUser -Name 'Administrator' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{}, errors.New("test-error"))
		_, err := c.UserRead(ctx, UserReadParams{Name: "Administrator"})
		suite.EqualError(err, "windows.local.accounts.UserRead: test-error")
	})
}

// Test GroupList related methods.
func (suite *LocalUnitTestSuite) TestUserList() {
	suite.Run("should return the correct list of user", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-LocalUser | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: userList}, nil)
		actualUserList, err := c.UserList(ctx)
		suite.NoError(err)
		suite.Equal(expectedUserList, actualUserList)
	})

	suite.Run("should return error if run fails", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-LocalUser | ConvertTo-Json -Compress").
			Return(connection.CmdResult{}, errors.New("test-error"))
		_, err := c.UserList(ctx)
		suite.EqualError(err, "windows.local.accounts.UserList: test-error")
	})
}

// Test UserCreate related methods.
func (suite *LocalUnitTestSuite) TestUserCreatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters UserCreateParams
			expectedCmd     string
		}{
			{
				"assert user with Name",
				UserCreateParams{Name: "Tester"},
				"New-LocalUser -Name 'Tester' -AccountNeverExpires -Disabled -NoPassword -UserMayNotChangePassword | ConvertTo-Json -Compress",
			},
			{
				"assert user with Name + Description",
				UserCreateParams{Name: "Tester", Description: "This is a test user"},
				"New-LocalUser -Name 'Tester' -Description 'This is a test user' -AccountNeverExpires -Disabled -NoPassword -UserMayNotChangePassword | ConvertTo-Json -Compress",
			},
			{
				"assert user with Name + AccountExpires",
				UserCreateParams{Name: "Tester", AccountExpires: time.Date(3024, time.April, 10, 15, 0, 0, 0, time.UTC)},
				"New-LocalUser -Name 'Tester' -AccountExpires $(Get-Date '3024-04-10 15:00:00') -Disabled -NoPassword -UserMayNotChangePassword | ConvertTo-Json -Compress",
			},
			{
				"assert user with Name + Enabled",
				UserCreateParams{Name: "Tester", Enabled: true},
				"New-LocalUser -Name 'Tester' -AccountNeverExpires -Disabled:$false -NoPassword -UserMayNotChangePassword | ConvertTo-Json -Compress",
			},
			{
				"assert user with Name + FullName",
				UserCreateParams{Name: "Tester", FullName: "Tester1"},
				"New-LocalUser -Name 'Tester' -AccountNeverExpires -Disabled -FullName 'Tester1' -NoPassword -UserMayNotChangePassword | ConvertTo-Json -Compress",
			},
			{
				"assert user with Name + Password",
				UserCreateParams{Name: "Tester", Password: "Start123!!!"},
				"New-LocalUser -Name 'Tester' -AccountNeverExpires -Disabled -Password $(ConvertTo-SecureString -String 'Start123!!!' -AsPlainText -Force) -PasswordNeverExpires:$false -UserMayNotChangePassword | ConvertTo-Json -Compress",
			},
			{
				"assert user with Name + PasswordNeverExpires + UserMayNotChangePassword",
				UserCreateParams{Name: "Tester", PasswordNeverExpires: true, UserMayChangePassword: true},
				"New-LocalUser -Name 'Tester' -AccountNeverExpires -Disabled -NoPassword -UserMayNotChangePassword:$false | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestUserCreate() {
	suite.Run("should return the correct new user", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "New-LocalUser -Name 'Test-User' -Description 'This is a test user' -AccountExpires $(Get-Date '2025-11-10 16:00:00') -Disabled:$false -FullName 'Full-Test-User' -NoPassword -UserMayNotChangePassword:$false | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: testUser}, nil)
		actualTestUser, err := c.UserCreate(ctx, UserCreateParams{
			Name:                  "Test-User",
			Description:           "This is a test user",
			FullName:              "Full-Test-User",
			Enabled:               true,
			UserMayChangePassword: true,
			AccountExpires:        time.Date(2025, time.November, 10, 16, 0, 0, 0, time.UTC),
		})
		suite.NoError(err)
		suite.Equal(expectedTestUser, actualTestUser)
	})
}

// Test UserUpdate related methods.
func (suite *LocalUnitTestSuite) TestUserUpdatePwshCommand() {
	suite.Run("should run the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters UserUpdateParams
			expectedCmd     string
		}{
			{
				"assert user with Name",
				UserUpdateParams{Name: "Tester"},
				"Set-LocalUser -Name 'Tester' -AccountNeverExpires -Description '' -FullName '' -PasswordNeverExpires:$false -UserMayChangePassword:$false ;Disable-LocalUser -Name 'Tester'",
			},
			{
				"assert user with Name + Enabled",
				UserUpdateParams{Name: "Tester", Enabled: true},
				"Set-LocalUser -Name 'Tester' -AccountNeverExpires -Description '' -FullName '' -PasswordNeverExpires:$false -UserMayChangePassword:$false ;Enable-LocalUser -Name 'Tester'",
			},
			{
				"assert user with SID + Enabled",
				UserUpdateParams{SID: "S-1000", Enabled: true},
				"Set-LocalUser -SID S-1000 -AccountNeverExpires -Description '' -FullName '' -PasswordNeverExpires:$false -UserMayChangePassword:$false ;Enable-LocalUser -SID S-1000",
			},
			{
				"assert user with Name + AccountExpires",
				UserUpdateParams{Name: "Tester", AccountExpires: time.Date(3024, time.April, 10, 15, 0, 0, 0, time.UTC)},
				"Set-LocalUser -Name 'Tester' -AccountExpires $(Get-Date '3024-04-10 15:00:00') -Description '' -FullName '' -PasswordNeverExpires:$false -UserMayChangePassword:$false ;Disable-LocalUser -Name 'Tester'",
			},
			{
				"assert user with Name + Description + FullName",
				UserUpdateParams{Name: "Tester", Description: "test-description", FullName: "Full-Tester"},
				"Set-LocalUser -Name 'Tester' -AccountNeverExpires -Description 'test-description' -FullName 'Full-Tester' -PasswordNeverExpires:$false -UserMayChangePassword:$false ;Disable-LocalUser -Name 'Tester'",
			},
			{
				"assert user with Name + Password + PasswordNeverExpires + UserMayChangePassword",
				UserUpdateParams{Name: "Tester", Password: "Start123!!!", PasswordNeverExpires: true, UserMayChangePassword: true},
				"Set-LocalUser -Name 'Tester' -AccountNeverExpires -Description '' -FullName '' -Password $(ConvertTo-SecureString -String 'Start123!!!' -AsPlainText -Force) -PasswordNeverExpires:$true -UserMayChangePassword:$true ;Disable-LocalUser -Name 'Tester'",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestUserUpdate() {
	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters UserUpdateParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				UserUpdateParams{},
				"windows.local.accounts.UserUpdate: user parameter 'Name' or 'SID' must be set",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mockConn := mockConnection.NewMockConnection(suite.T())
			c := &Client{
				Connection:      mockConn,
				decodeCliXmlErr: func(s string) (string, error) { return s, nil },
			}
			err := c.UserUpdate(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

// Test UserDelete related methods.
func (suite *LocalUnitTestSuite) TestUserDeletePwshCommand() {
	suite.Run("should run the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters UserDeleteParams
			expectedCmd     string
		}{
			{
				"assert user with Name",
				UserDeleteParams{Name: "Tester"},
				"Remove-LocalUser -Name 'Tester'",
			},
			{
				"assert user with SID",
				UserDeleteParams{SID: "S-1000"},
				"Remove-LocalUser -SID S-1000",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestUserDelete() {
	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters UserDeleteParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				UserDeleteParams{},
				"windows.local.accounts.UserDelete: user parameter 'Name' or 'SID' must be set",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mockConn := mockConnection.NewMockConnection(suite.T())
			c := &Client{
				Connection:      mockConn,
				decodeCliXmlErr: func(s string) (string, error) { return s, nil },
			}
			err := c.UserDelete(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}
