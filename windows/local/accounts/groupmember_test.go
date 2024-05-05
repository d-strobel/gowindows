package accounts

import (
	"context"
	"errors"

	"github.com/d-strobel/gowindows/connection"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
)

// Fixtures
const (
	groupMemberRead = `{"Name":"WIN2022SC\\Administrator","SID":{"BinaryLength":28,"AccountDomainSid":{"BinaryLength":24,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138"},"Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"}`
	groupMemberList = `[{"Name":"WIN2022SC\\Administrator","SID":{"BinaryLength":28,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138-500"},"PrincipalSource":1,"ObjectClass":"User"},{"Name":"WIN2022SC\\vagrant","SID":{"BinaryLength":28,"AccountDomainSid":"S-1-5-21-153895498-367353507-3704405138","Value":"S-1-5-21-153895498-367353507-3704405138-1000"},"PrincipalSource":1,"ObjectClass":"User"}]`
)

var (
	expectedGroupMemberRead = GroupMember{
		Name: "WIN2022SC\\Administrator",
		SID: SID{
			Value: "S-1-5-21-153895498-367353507-3704405138-500",
		},
		ObjectClass: "User",
	}
	expectedGroupMemberList = []GroupMember{
		{
			Name: "WIN2022SC\\Administrator",
			SID: SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-500",
			},
			ObjectClass: "User",
		},
		{
			Name: "WIN2022SC\\vagrant",
			SID: SID{
				Value: "S-1-5-21-153895498-367353507-3704405138-1000",
			},
			ObjectClass: "User",
		},
	}
)

// Test GroupMemberRead related methods.
func (suite *LocalUnitTestSuite) TestGroupMemberReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberReadParams
			expectedCmd     string
		}{
			{
				"assert user by name",
				GroupMemberReadParams{Name: "Administrators", Member: "Administrator"},
				"Get-LocalGroupMember -Name 'Administrators' -Member 'Administrator' | ConvertTo-Json -Compress",
			},
			{
				"assert users by sid",
				GroupMemberReadParams{SID: "123456789", Member: "Test"},
				"Get-LocalGroupMember -SID 123456789 -Member 'Test' | ConvertTo-Json -Compress",
			},
			{
				"assert users by name and sid",
				GroupMemberReadParams{Name: "Users", SID: "123456789", Member: "Test"},
				"Get-LocalGroupMember -SID 123456789 -Member 'Test' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestGroupMemberRead() {
	suite.Run("should return the correct GroupMember", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-LocalGroupMember -Name 'Administrators' -Member 'Administrator' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: groupMemberRead}, nil)
		actualGroupMemberRead, err := c.GroupMemberRead(ctx, GroupMemberReadParams{Name: "Administrators", Member: "Administrator"})
		suite.NoError(err)
		suite.Equal(expectedGroupMemberRead, actualGroupMemberRead)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				GroupMemberReadParams{},
				"windows.local.accounts.GroupMemberRead: group member parameter 'Name' or 'SID' must be set",
			},
			{
				"assert no member",
				GroupMemberReadParams{Name: "Administrators"},
				"windows.local.accounts.GroupMemberRead: group member parameter 'Member' must be set",
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
			_, err := c.GroupMemberRead(ctx, tc.inputParameters)
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
			RunWithPowershell(ctx, "Get-LocalGroupMember -Name 'Administrator' -Member 'Test' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{}, errors.New("test-error"))
		_, err := c.GroupMemberRead(ctx, GroupMemberReadParams{Name: "Administrator", Member: "Test"})
		suite.EqualError(err, "windows.local.accounts.GroupMemberRead: test-error")
	})
}

// Test GroupMemberList related methods.
func (suite *LocalUnitTestSuite) TestGroupMemberListPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberListParams
			expectedCmd     string
		}{
			{
				"assert group member list by Name",
				GroupMemberListParams{Name: "Users"},
				"$gm=Get-LocalGroupMember -Name 'Users' ;if($gm.Count -eq 1){ConvertTo-Json @($gm) -Compress}else{ConvertTo-Json $gm -Compress}",
			},
			{
				"assert group member list by SID",
				GroupMemberListParams{SID: "123456789"},
				"$gm=Get-LocalGroupMember -SID 123456789 ;if($gm.Count -eq 1){ConvertTo-Json @($gm) -Compress}else{ConvertTo-Json $gm -Compress}",
			},
			{
				"assert group member list by SID and Name",
				GroupMemberListParams{Name: "Users", SID: "123456789"},
				"$gm=Get-LocalGroupMember -SID 123456789 ;if($gm.Count -eq 1){ConvertTo-Json @($gm) -Compress}else{ConvertTo-Json $gm -Compress}",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestGroupMemberList() {
	suite.Run("should return the correct list of group member", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "$gm=Get-LocalGroupMember -Name 'Administrators' ;if($gm.Count -eq 1){ConvertTo-Json @($gm) -Compress}else{ConvertTo-Json $gm -Compress}").
			Return(connection.CmdResult{StdOut: groupMemberList}, nil)
		actualGroupMemberList, err := c.GroupMemberList(ctx, GroupMemberListParams{Name: "Administrators"})
		suite.NoError(err)
		suite.Equal(expectedGroupMemberList, actualGroupMemberList)
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
			RunWithPowershell(ctx, "$gm=Get-LocalGroupMember -Name 'Administrators' ;if($gm.Count -eq 1){ConvertTo-Json @($gm) -Compress}else{ConvertTo-Json $gm -Compress}").
			Return(connection.CmdResult{}, errors.New("test-error"))
		_, err := c.GroupMemberList(ctx, GroupMemberListParams{Name: "Administrators"})
		suite.EqualError(err, "windows.local.accounts.GroupMemberList: test-error")
	})
}

// Test GroupMemberCreate related methods.
func (suite *LocalUnitTestSuite) TestGroupMemberCreatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberCreateParams
			expectedCmd     string
		}{
			{
				"assert user with Name + Member",
				GroupMemberCreateParams{Name: "Administrators", Member: "TestUser"},
				"Add-LocalGroupMember -Name 'Administrators' -Member 'TestUser'",
			},
			{
				"assert user with Name + SID + Member",
				GroupMemberCreateParams{Name: "Administrators", SID: "123456", Member: "TestUser"},
				"Add-LocalGroupMember -SID 123456 -Member 'TestUser'",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestGroupMemberCreate() {
	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberCreateParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				GroupMemberCreateParams{},
				"windows.local.accounts.GroupMemberCreate: group member parameter 'Name' or 'SID' must be set",
			},
			{
				"assert no member",
				GroupMemberCreateParams{Name: "Administrators"},
				"windows.local.accounts.GroupMemberCreate: group member parameter 'Member' must be set",
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
			err := c.GroupMemberCreate(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

// Test GroupMemberDelete related methods.
func (suite *LocalUnitTestSuite) TestGroupMemberDeletePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberDeleteParams
			expectedCmd     string
		}{
			{
				"assert user with Name",
				GroupMemberDeleteParams{Name: "Administrators", Member: "TestUser"},
				"Remove-LocalGroupMember -Name 'Administrators' -Member 'TestUser'",
			},
			{
				"assert user with Name + SID + Member",
				GroupMemberDeleteParams{Name: "Administrators", SID: "123456", Member: "TestUser"},
				"Remove-LocalGroupMember -SID 123456 -Member 'TestUser'",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *LocalUnitTestSuite) TestGroupMemberDelete() {
	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberDeleteParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				GroupMemberDeleteParams{},
				"windows.local.accounts.GroupMemberDelete: group member parameter 'Name' or 'SID' must be set",
			},
			{
				"assert no member",
				GroupMemberDeleteParams{Name: "Administrators"},
				"windows.local.accounts.GroupMemberDelete: group member parameter 'Member' must be set",
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
			err := c.GroupMemberDelete(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}
