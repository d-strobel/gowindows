package accounts

import (
	"context"
	"errors"

	"github.com/d-strobel/gowindows/connection"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
)

// Fixtures
const (
	usersGroup = `{"Description":"Users are prevented from making accidental or intentional system-wide changes and can run most applications","Name":"Users","SID":{"BinaryLength":16,"AccountDomainSid":null,"Value":"S-1-5-32-545"},"PrincipalSource":1,"ObjectClass":"Group"}`
	groupList  = `[{"Description":"Users are prevented from making accidental or intentional system-wide changes and can run most applications","Name":"Users","SID":{"BinaryLength":16,"AccountDomainSid":null,"Value":"S-1-5-32-545"},"PrincipalSource":1,"ObjectClass":"Group"},{"Description":"Administrators have complete and unrestricted access to the computer/domain","Name":"Administrators","SID":{"BinaryLength":16,"AccountDomainSid":null,"Value":"S-1-5-32-544"},"PrincipalSource":1,"ObjectClass":"Group"}]`
	testGroup  = `{"Description":"Test group","Name":"Test","SID":{"BinaryLength":16,"AccountDomainSid":null,"Value":"S-123456789"},"PrincipalSource":1,"ObjectClass":"Group"}`
)

var (
	expectedUsersGroup = Group{
		Name:        "Users",
		Description: "Users are prevented from making accidental or intentional system-wide changes and can run most applications",
		SID: SID{
			Value: "S-1-5-32-545",
		},
	}
	expectedGroupList = []Group{
		{
			Name:        "Users",
			Description: "Users are prevented from making accidental or intentional system-wide changes and can run most applications",
			SID: SID{
				Value: "S-1-5-32-545",
			},
		},
		{
			Name:        "Administrators",
			Description: "Administrators have complete and unrestricted access to the computer/domain",
			SID: SID{
				Value: "S-1-5-32-544",
			},
		},
	}
	expectedTestGroup = Group{
		Name:        "Test",
		Description: "Test group",
		SID: SID{
			Value: "S-123456789",
		},
	}
)

// Test GroupRead related methods.
func (suite *LocalUnitTestSuite) TestGroupReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupReadParams
			expectedCmd     string
		}{
			{
				"assert users group by name",
				GroupReadParams{Name: "Users"},
				"Get-LocalGroup -Name 'Users' | ConvertTo-Json -Compress",
			},
			{
				"assert users group by sid",
				GroupReadParams{SID: "123456789"},
				"Get-LocalGroup -SID 123456789 | ConvertTo-Json -Compress",
			},
			{
				"assert users group by name and sid",
				GroupReadParams{Name: "Users", SID: "123456789"},
				"Get-LocalGroup -SID 123456789 | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestGroupRead() {
	suite.Run("should return the correct group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-LocalGroup -Name 'Users' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: usersGroup}, nil)
		actualUsersGroup, err := c.GroupRead(ctx, GroupReadParams{Name: "Users"})
		suite.NoError(err)
		suite.Equal(expectedUsersGroup, actualUsersGroup)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters GroupReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				GroupReadParams{},
				"windows.local.accounts.GroupRead: group parameter 'Name' or 'SID' must be set",
			},
			{
				"assert error when name contains wildcard",
				GroupReadParams{Name: "Remote*"},
				"windows.local.accounts.GroupRead: group parameter 'Name' does not allow wildcards",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mockConn := mockConnection.NewMockConnection(suite.T())
			c := &Client{
				Connection:      mockConn,
				decodeCliXmlErr: func(s string) (string, error) { return "", nil },
			}
			_, err := c.GroupRead(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})

	suite.Run("should return error if run fails", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-LocalGroup -Name 'Users' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{}, errors.New("test-error"))
		_, err := c.GroupRead(ctx, GroupReadParams{Name: "Users"})
		suite.EqualError(err, "windows.local.accounts.GroupRead: test-error")
	})
}

// Test GroupList related methods.
func (suite *LocalUnitTestSuite) TestGroupList() {
	suite.Run("should return the correct list of groups", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Get-LocalGroup | ConvertTo-Json -Compress"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{StdOut: groupList}, nil)
		actualGroupList, err := c.GroupList(ctx)
		suite.NoError(err)
		suite.Equal(expectedGroupList, actualGroupList)
	})

	suite.Run("should return error if run fails", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Get-LocalGroup | ConvertTo-Json -Compress"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{}, errors.New("test-error"))
		_, err := c.GroupList(ctx)
		suite.EqualError(err, "windows.local.accounts.GroupList: test-error")
	})
}

// Test GroupCreate related methods.
func (suite *LocalUnitTestSuite) TestGroupCreatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupCreateParams
			expectedCmd     string
		}{
			{
				"assert without description parameter",
				GroupCreateParams{Name: "Test"},
				"New-LocalGroup -Name 'Test' | ConvertTo-Json -Compress",
			},
			{
				"assert with name and description parameter",
				GroupCreateParams{Name: "Test", Description: "Test group"},
				"New-LocalGroup -Name 'Test' -Description 'Test group' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestGroupCreate() {
	suite.Run("should return the correct group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "New-LocalGroup -Name 'Test' -Description 'Test group' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: testGroup}, nil)
		actualTestGroup, err := c.GroupCreate(ctx, GroupCreateParams{Name: "Test", Description: "Test group"})
		suite.NoError(err)
		suite.Equal(expectedTestGroup, actualTestGroup)
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
			RunWithPowershell(ctx, "New-LocalGroup -Name 'Test' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{}, errors.New("test-error"))
		_, err := c.GroupCreate(ctx, GroupCreateParams{Name: "Test"})
		suite.EqualError(err, "windows.local.accounts.GroupCreate: test-error")
	})
}

// Test GroupUpdate related methods.
func (suite *LocalUnitTestSuite) TestGroupUpdatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupUpdateParams
			expectedCmd     string
		}{
			{
				"assert with Name and Desctiption parameter",
				GroupUpdateParams{Name: "Test", Description: "Testing"},
				"Set-LocalGroup -Name 'Test' -Description 'Testing'",
			},
			{
				"assert with SID and Desctiption parameter",
				GroupUpdateParams{SID: "S-12345", Description: "Testing"},
				"Set-LocalGroup -SID S-12345 -Description 'Testing'",
			},
			{
				"assert with Name, SID and Desctiption parameter",
				GroupUpdateParams{Name: "Test", SID: "S-12345", Description: "Testing"},
				"Set-LocalGroup -SID S-12345 -Description 'Testing'",
			},
			{
				"assert with Name parameter",
				GroupUpdateParams{Name: "Test"},
				"Set-LocalGroup -Name 'Test' -Description ' '",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestGroupUpdate() {
	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters GroupUpdateParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				GroupUpdateParams{},
				"windows.local.accounts.GroupUpdate: group parameter 'Name' or 'SID' must be set",
			},
			{
				"assert error with just the description parameter",
				GroupUpdateParams{Description: "test"},
				"windows.local.accounts.GroupUpdate: group parameter 'Name' or 'SID' must be set",
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
			err := c.GroupUpdate(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
			mockConn.AssertNotCalled(suite.T(), "RunWithPowershell")
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
			RunWithPowershell(ctx, "Set-LocalGroup -Name 'Test' -Description 'Test'").
			Return(connection.CmdResult{}, errors.New("test-error"))
		err := c.GroupUpdate(ctx, GroupUpdateParams{Name: "Test", Description: "Test"})
		suite.EqualError(err, "windows.local.accounts.GroupUpdate: test-error")
	})
}

// Test GroupDelete related methods.
func (suite *LocalUnitTestSuite) TestGroupDeletePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupDeleteParams
			expectedCmd     string
		}{
			{
				"assert with Name parameter",
				GroupDeleteParams{Name: "Test"},
				"Remove-LocalGroup -Name 'Test'",
			},
			{
				"assert with SID parameter",
				GroupDeleteParams{SID: "S-12345"},
				"Remove-LocalGroup -SID S-12345",
			},
			{
				"assert with Name and SID parameter",
				GroupDeleteParams{Name: "Test", SID: "S-12345"},
				"Remove-LocalGroup -SID S-12345",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestGroupDelete() {
	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters GroupDeleteParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				GroupDeleteParams{},
				"windows.local.accounts.GroupDelete: group parameter 'Name' or 'SID' must be set",
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
			err := c.GroupDelete(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
			mockConn.AssertNotCalled(suite.T(), "RunWithPowershell")
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
			RunWithPowershell(ctx, "Remove-LocalGroup -Name 'Test'").
			Return(connection.CmdResult{}, errors.New("test-error"))
		err := c.GroupDelete(ctx, GroupDeleteParams{Name: "Test"})
		suite.EqualError(err, "windows.local.accounts.GroupDelete: test-error")
	})
}
