package local

import (
	"context"
	"errors"
	"testing"

	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
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
	suite.Run("should return a new local client", func() {
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockDecodeCliXmlErr := func(s string) (string, error) { return "", nil }
		actualLocalClient := NewClientWithParser(mockConn, mockDecodeCliXmlErr)
		expectedLocalClient := &LocalClient{Connection: mockConn, decodeCliXmlErr: mockDecodeCliXmlErr}
		suite.IsType(expectedLocalClient, actualLocalClient)
		suite.Equal(expectedLocalClient.Connection, actualLocalClient.Connection)
	})
}

func (suite *LocalUnitTestSuite) TestLocalRun() {
	suite.T().Parallel()

	suite.Run("should return the user group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &LocalClient{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Get-LocalGroup -Name Users | ConvertTo-Json -Compress"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{StdOut: usersGroup, StdErr: ""}, nil)
		var g Group
		err := localRun(ctx, c, cmd, &g)
		suite.NoError(err)
		suite.Equal(expectedUsersGroup, g)
	})

	suite.Run("should return a slice of groups", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &LocalClient{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Get-LocalGroup | ConvertTo-Json"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{StdOut: groupList, StdErr: ""}, nil)
		var g []Group
		err := localRun(ctx, c, cmd, &g)
		suite.NoError(err)
		suite.Equal(expectedGroupList, g)
	})

	suite.Run("should not error when no stdout is empty string", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &LocalClient{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Remove-LocalGroup -Name Test"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{StdOut: "", StdErr: ""}, nil)
		var g Group
		var expectedGroup Group
		err := localRun(ctx, c, cmd, &g)
		suite.NoError(err)
		suite.Equal(expectedGroup, g)
	})

	suite.Run("should error when connection run errors", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &LocalClient{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Remove-LocalGroup -Name Test"
		expectedErr := errors.New("test-error")
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{}, expectedErr)
		var g Group
		err := localRun(ctx, c, cmd, &g)
		suite.Error(err)
		suite.Equal(expectedErr, err)
	})

	suite.Run("should return powershell error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &LocalClient{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		cmd := "Get-LocalGroup -name Userrs"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{StdOut: "", StdErr: "test-error"}, nil)
		var g Group
		expectedErr := errors.New("Command:\nGet-LocalGroup -name Userrs\n\nPowershell error:\ntest-error")
		err := localRun(ctx, c, cmd, &g)
		suite.Error(err)
		suite.Equal(expectedErr, err)
	})

	suite.Run("should return error from json unmarshal with incorrect json", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &LocalClient{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		cmd := "Get-LocalGroup -name Users"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{StdOut: groupList}, nil)
		var g Group
		err := localRun(ctx, c, cmd, &g)
		suite.Error(err)
	})
}
