package local

import (
	"context"
	"errors"

	"github.com/d-strobel/gowindows/connection"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	mockParser "github.com/d-strobel/gowindows/parser/mocks"
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

func (suite *LocalUnitTestSuite) TestGroupMemberRead() {

	suite.Run("should return the correct GroupMember", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroupMember -Name 'Administrators' -Member 'Administrator' | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: groupMemberRead,
		}, nil)
		actualGroupMemberRead, err := c.GroupMemberRead(ctx, GroupMemberParams{Name: "Administrators", Member: "Administrator"})
		suite.Require().NoError(err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		suite.Equal(expectedGroupMemberRead, actualGroupMemberRead)
	})

	suite.Run("should run the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberParams
			expectedCMD     string
		}{
			{
				"assert user by name",
				GroupMemberParams{Name: "Administrators", Member: "Administrator"},
				"Get-LocalGroupMember -Name 'Administrators' -Member 'Administrator' | ConvertTo-Json -Compress",
			},
			{
				"assert users by sid",
				GroupMemberParams{SID: "123456789", Member: "Test"},
				"Get-LocalGroupMember -SID 123456789 -Member 'Test' | ConvertTo-Json -Compress",
			},
			{
				"assert users by name and sid",
				GroupMemberParams{Name: "Users", SID: "123456789", Member: "Test"},
				"Get-LocalGroupMember -SID 123456789 -Member 'Test' | ConvertTo-Json -Compress",
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
			_, err := c.GroupMemberRead(ctx, tc.inputParameters)
			suite.Require().NoError(err)
			mockConn.AssertCalled(suite.T(), "Run", ctx, tc.expectedCMD)
			mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		}
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				GroupMemberParams{},
				"windows.local.GroupMemberRead: group member parameter 'Name' or 'SID' must be set",
			},
			{
				"assert no member",
				GroupMemberParams{Name: "Administrators"},
				"windows.local.GroupMemberRead: group member parameter 'Member' must be set",
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
			_, err := c.GroupMemberRead(ctx, tc.inputParameters)
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
		expectedCMD := "Get-LocalGroupMember -Name 'Administrator' -Member 'Test' | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		_, err := c.GroupMemberRead(ctx, GroupMemberParams{Name: "Administrator", Member: "Test"})
		suite.EqualError(err, "windows.local.GroupMemberRead: test-error")
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}

func (suite *LocalUnitTestSuite) TestGroupMemberList() {

	suite.Run("should return the correct list of group member", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "$gm=Get-LocalGroupMember -Name 'Administrators' ;if($gm.Count -eq 1){ConvertTo-Json @($gm) -Compress}else{ConvertTo-Json $gm -Compress}"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: groupMemberList,
		}, nil)
		actualGroupMemberList, err := c.GroupMemberList(ctx, GroupMemberParams{Name: "Administrators"})
		suite.Require().NoError(err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		suite.Equal(expectedGroupMemberList, actualGroupMemberList)
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
		expectedCMD := "$gm=Get-LocalGroupMember -Name 'Administrators' ;if($gm.Count -eq 1){ConvertTo-Json @($gm) -Compress}else{ConvertTo-Json $gm -Compress}"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		_, err := c.GroupMemberList(ctx, GroupMemberParams{Name: "Administrators"})
		suite.EqualError(err, "windows.local.GroupMemberList: test-error")
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}

func (suite *LocalUnitTestSuite) TestGroupMemberCreate() {
	suite.Run("should run the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberParams
			expectedCMD     string
		}{
			{
				"assert user with Name + Member",
				GroupMemberParams{Name: "Administrators", Member: "TestUser"},
				"Add-LocalGroupMember -Name 'Administrators' -Member 'TestUser'",
			},
			{
				"assert user with Name + SID + Member",
				GroupMemberParams{Name: "Administrators", SID: "123456", Member: "TestUser"},
				"Add-LocalGroupMember -SID 123456 -Member 'TestUser'",
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
			err := c.GroupMemberCreate(ctx, tc.inputParameters)
			suite.Require().NoError(err)
			mockConn.AssertCalled(suite.T(), "Run", ctx, tc.expectedCMD)
			mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		}
	})
}

func (suite *LocalUnitTestSuite) TestGroupMemberDelete() {
	suite.Run("should run the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupMemberParams
			expectedCMD     string
		}{
			{
				"assert user with Name",
				GroupMemberParams{Name: "Administrators", Member: "TestUser"},
				"Remove-LocalGroupMember -Name 'Administrators' -Member 'TestUser'",
			},
			{
				"assert user with Name + SID + Member",
				GroupMemberParams{Name: "Administrators", SID: "123456", Member: "TestUser"},
				"Remove-LocalGroupMember -SID 123456 -Member 'TestUser'",
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
			err := c.GroupMemberDelete(ctx, tc.inputParameters)
			suite.Require().NoError(err)
			mockConn.AssertCalled(suite.T(), "Run", ctx, tc.expectedCMD)
			mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		}
	})
}
