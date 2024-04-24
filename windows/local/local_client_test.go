package local

import (
	"context"
	"errors"
	"testing"

	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	"github.com/d-strobel/gowindows/parser"
	mockParser "github.com/d-strobel/gowindows/parser/mocks"
	"github.com/stretchr/testify/suite"
)

// Unit test suite for all local functions
type LocalUnitTestSuite struct {
	suite.Suite
}

// Run all local unit tests
func TestLocalUnitTestSuite(t *testing.T) {
	suite.Run(t, &LocalUnitTestSuite{})
}

func (suite *LocalUnitTestSuite) TestNewClient() {
	mockConn := mockConnection.NewMockConnection(suite.T())
	mockParser := &parser.Parser{}
	actualLocalClient := NewClientWithParser(mockConn, mockParser)
	expectedLocalClient := &LocalClient{Connection: mockConn, parser: mockParser}

	suite.Equal(expectedLocalClient, actualLocalClient)
}

func (suite *LocalUnitTestSuite) TestLocalRun() {

	suite.Run("should return the user group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup -Name Users | ConvertTo-Json -Compress"
		mockConn.On("RunWithPowershell", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: usersGroup,
		}, nil)
		var g Group
		err := localRun(ctx, c, expectedCMD, &g)
		suite.NoError(err)
		suite.Equal(expectedUsersGroup, g)
		mockConn.AssertCalled(suite.T(), "RunWithPowershell", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should return a slice of groups", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup | ConvertTo-Json"
		mockConn.On("RunWithPowershell", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: groupList,
		}, nil)
		var g []Group
		err := localRun(ctx, c, expectedCMD, &g)
		suite.NoError(err)
		suite.Equal(expectedGroupList, g)
		mockConn.AssertCalled(suite.T(), "RunWithPowershell", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should return the user group with compressed json", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup -Name Users | ConvertTo-Json -compress"
		mockConn.On("RunWithPowershell", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: usersGroup,
		}, nil)
		var g Group
		err := localRun(ctx, c, expectedCMD, &g)
		suite.NoError(err)
		suite.Equal(expectedUsersGroup, g)
		mockConn.AssertCalled(suite.T(), "RunWithPowershell", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should not error when no stdout is empty string", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Remove-LocalGroup -Name Test"
		mockConn.On("RunWithPowershell", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: "",
		}, nil)
		var g Group
		var expectedGroup Group
		err := localRun(ctx, c, expectedCMD, &g)
		suite.NoError(err)
		suite.Equal(expectedGroup, g)
		mockConn.AssertCalled(suite.T(), "RunWithPowershell", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should error when connection run errors", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Remove-LocalGroup -Name Test"
		mockConn.On("RunWithPowershell", ctx, expectedCMD).Return(connection.CMDResult{}, errors.New("test-error"))
		var g Group
		expectedErr := errors.New("test-error")
		err := localRun(ctx, c, expectedCMD, &g)
		suite.Error(err)
		suite.Equal(expectedErr, err)
		mockConn.AssertCalled(suite.T(), "RunWithPowershell", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})

	suite.Run("should return powershell error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup -name Userrs"
		mockConn.On("RunWithPowershell", ctx, expectedCMD).Return(connection.CMDResult{
			StdErr: "clixml-error",
		}, nil)
		mockParser.On("DecodeCLIXML", "clixml-error").Return("powershell-error", nil)
		var g Group
		expectedErr := errors.New("Command:\nGet-LocalGroup -name Userrs\n\nPowershell error:\npowershell-error")
		err := localRun(ctx, c, expectedCMD, &g)
		suite.Error(err)
		suite.Equal(expectedErr, err)
		mockConn.AssertCalled(suite.T(), "RunWithPowershell", ctx, expectedCMD)
		mockParser.AssertCalled(suite.T(), "DecodeCLIXML", "clixml-error")
	})

	suite.Run("should return an error from parser", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup -name Userrs"
		mockConn.On("RunWithPowershell", ctx, expectedCMD).Return(connection.CMDResult{
			StdErr: "incorrect-clixml-error",
		}, nil)
		mockParser.On("DecodeCLIXML", "incorrect-clixml-error").Return("", errors.New("parser-error"))
		var g Group
		expectedErr := errors.New("parser-error")
		err := localRun(ctx, c, expectedCMD, &g)
		suite.Error(err)
		suite.Equal(expectedErr, err)
		mockConn.AssertCalled(suite.T(), "RunWithPowershell", ctx, expectedCMD)
		mockParser.AssertCalled(suite.T(), "DecodeCLIXML", "incorrect-clixml-error")
	})

	suite.Run("should return error from json unmarshal with incorrect json", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockParser := mockParser.NewMockParserInterface(suite.T())
		c := &LocalClient{
			Connection: mockConn,
			parser:     mockParser,
		}
		expectedCMD := "Get-LocalGroup -name Users"
		mockConn.On("RunWithPowershell", ctx, expectedCMD).Return(connection.CMDResult{
			StdOut: groupList,
		}, nil)
		var g Group
		err := localRun(ctx, c, expectedCMD, &g)
		suite.Error(err)
		mockConn.AssertCalled(suite.T(), "RunWithPowershell", ctx, expectedCMD)
		mockParser.AssertNotCalled(suite.T(), "DecodeCLIXML")
	})
}
