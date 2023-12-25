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

func (suite *LocalUnitTestSuite) TestGroupRun() {

	suite.Run("should return the user group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup -Name Users | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: usersGroup,
		}, nil)
		var g Group
		err := groupRun[Group](ctx, c, expectedCMD, &g)
		suite.NoError(err)
		suite.Equal(expectedUsersGroup, g)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
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
		expectedCMD := "Get-LocalGroup | ConvertTo-Json"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: groupList,
		}, nil)
		var g []Group
		err := groupRun[[]Group](ctx, c, expectedCMD, &g)
		suite.NoError(err)
		suite.Equal(expectedGroupList, g)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
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
		expectedCMD := "Get-LocalGroup -Name Users | ConvertTo-Json -compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: usersGroup,
		}, nil)
		var g Group
		err := groupRun[Group](ctx, c, expectedCMD, &g)
		suite.NoError(err)
		suite.Equal(expectedUsersGroup, g)
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
		expectedCMD := "Remove-LocalGroup -Name Test"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: "",
		}, nil)
		var g Group
		var expectedGroup Group
		err := groupRun[Group](ctx, c, expectedCMD, &g)
		suite.NoError(err)
		suite.Equal(expectedGroup, g)
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
		expectedCMD := "Remove-LocalGroup -Name Test"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		var g Group
		expectedErr := errors.New("test-error")
		err := groupRun[Group](ctx, c, expectedCMD, &g)
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
		expectedCMD := "Get-LocalGroup -name Userrs"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdErr: "clixml-error",
		}, nil)
		mockParser.On("DecodeCLIXML", "clixml-error").Return("powershell-error", nil)
		var g Group
		expectedErr := errors.New("powershell-error")
		err := groupRun[Group](ctx, c, expectedCMD, &g)
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
		expectedCMD := "Get-LocalGroup -name Userrs"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdErr: "incorrect-clixml-error",
		}, nil)
		mockParser.On("DecodeCLIXML", "incorrect-clixml-error").Return("", errors.New("parser-error"))
		var g Group
		expectedErr := errors.New("parser-error")
		err := groupRun[Group](ctx, c, expectedCMD, &g)
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
		expectedCMD := "Get-LocalGroup -name Users"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: groupList,
		}, nil)
		var g Group
		err := groupRun[Group](ctx, c, expectedCMD, &g)
		suite.Error(err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}

func (suite *LocalUnitTestSuite) TestGroupRead() {

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

func (suite *LocalUnitTestSuite) TestGroupList() {

	suite.Run("should return the correct list of groups", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: groupList,
		}, nil)
		actualGroupList, err := c.GroupList(ctx)
		suite.Require().NoError(err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		suite.Equal(expectedGroupList, actualGroupList)
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
		expectedCMD := "Get-LocalGroup | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		_, err := c.GroupList(ctx)
		suite.EqualError(err, "windows.local.GroupList: test-error")
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}

func (suite *LocalUnitTestSuite) TestGroupCreate() {

	suite.Run("should return the correct group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnectionInterface(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "New-LocalGroup -Name 'Test' -Description 'Test group' | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: testGroup,
		}, nil)
		actualTestGroup, err := c.GroupCreate(ctx, GroupParams{Name: "Test", Description: "Test group"})
		suite.Require().NoError(err)
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
		suite.Equal(expectedTestGroup, actualTestGroup)
	})

	suite.Run("should run the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupParams
			expectedCMD     string
		}{
			{
				"assert without description parameter",
				GroupParams{Name: "Test"},
				"New-LocalGroup -Name 'Test' | ConvertTo-Json -Compress",
			},
			{
				"assert with name and description parameter",
				GroupParams{Name: "Test", Description: "Test group"},
				"New-LocalGroup -Name 'Test' -Description 'Test group' | ConvertTo-Json -Compress",
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
			_, err := c.GroupCreate(ctx, tc.inputParameters)
			suite.Require().NoError(err)
			mockConn.AssertCalled(suite.T(), "Run", ctx, tc.expectedCMD)
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
		expectedCMD := "New-LocalGroup -Name 'Test' | ConvertTo-Json -Compress"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		_, err := c.GroupCreate(ctx, GroupParams{Name: "Test"})
		suite.EqualError(err, "windows.local.GroupCreate: test-error")
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}

func (suite *LocalUnitTestSuite) TestGroupUpdate() {

	suite.Run("should run the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupParams
			expectedCMD     string
		}{
			{
				"assert with Name and Desctiption parameter",
				GroupParams{Name: "Test", Description: "Testing"},
				"Set-LocalGroup -Name 'Test' -Description 'Testing'",
			},
			{
				"assert with SID and Desctiption parameter",
				GroupParams{SID: "S-12345", Description: "Testing"},
				"Set-LocalGroup -SID S-12345 -Description 'Testing'",
			},
			{
				"assert with Name, SID and Desctiption parameter",
				GroupParams{Name: "Test", SID: "S-12345", Description: "Testing"},
				"Set-LocalGroup -SID S-12345 -Description 'Testing'",
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
			err := c.GroupUpdate(ctx, tc.inputParameters)
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
				"windows.local.GroupUpdate: group parameter 'Name' or 'SID' must be set",
			},
			{
				"assert error with just the description parameter",
				GroupParams{Description: "test"},
				"windows.local.GroupUpdate: group parameter 'Name' or 'SID' must be set",
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
			err := c.GroupUpdate(ctx, tc.inputParameters)
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
		expectedCMD := "Set-LocalGroup -Name 'Test' -Description 'Test'"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		err := c.GroupUpdate(ctx, GroupParams{Name: "Test", Description: "Test"})
		suite.EqualError(err, "windows.local.GroupUpdate: test-error")
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}

func (suite *LocalUnitTestSuite) TestGroupDelete() {

	suite.Run("should run the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters GroupParams
			expectedCMD     string
		}{
			{
				"assert with Name parameter",
				GroupParams{Name: "Test"},
				"Remove-LocalGroup -Name 'Test'",
			},
			{
				"assert with SID parameter",
				GroupParams{SID: "S-12345"},
				"Remove-LocalGroup -SID S-12345",
			},
			{
				"assert with Name and SID parameter",
				GroupParams{Name: "Test", SID: "S-12345"},
				"Remove-LocalGroup -SID S-12345",
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
			err := c.GroupDelete(ctx, tc.inputParameters)
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
				"windows.local.GroupDelete: group parameter 'Name' or 'SID' must be set",
			},
			{
				"assert error with just the description parameter",
				GroupParams{Description: "test"},
				"windows.local.GroupDelete: group parameter 'Name' or 'SID' must be set",
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
			err := c.GroupDelete(ctx, tc.inputParameters)
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
		expectedCMD := "Remove-LocalGroup -Name 'Test'"
		mockConn.On("Run", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		err := c.GroupDelete(ctx, GroupParams{Name: "Test"})
		suite.EqualError(err, "windows.local.GroupDelete: test-error")
		mockConn.AssertCalled(suite.T(), "Run", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}
