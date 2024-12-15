package dhcp

import (
	"context"
	"testing"

	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	"github.com/stretchr/testify/suite"
)

// Unit test suite for all dhcp functions
type DhcpServerUnitTestSuite struct {
	suite.Suite
}

// Run all dhcp unit tests
func TestDhcpServerUnitTestSuite(t *testing.T) {
	suite.Run(t, &DhcpServerUnitTestSuite{})
}

func (suite *DhcpServerUnitTestSuite) TestNewClient() {
	suite.Run("should return a new dhcp client", func() {
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockDecodeCliXmlErr := func(s string) (string, error) { return "", nil }
		actualClient := NewClientWithParser(mockConn, mockDecodeCliXmlErr)
		expectedClient := &Client{Connection: mockConn, decodeCliXmlErr: mockDecodeCliXmlErr}
		suite.IsType(expectedClient, actualClient)
		suite.Equal(expectedClient.Connection, actualClient.Connection)
	})
}

func (suite *DhcpServerUnitTestSuite) TestDhcpRun() {
	suite.T().Parallel()

	suite.Run("should return an unmarshalled scope object", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Get-DhcpServerv4Scope -ScopeId '192.168.10.0' | ConvertTo-Json -Compress"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{StdOut: scopeV4Json}, nil)
		var o scopeObject
		err := run(ctx, c, cmd, &o)
		suite.NoError(err)
		suite.Equal(expectedScopeObject, o)
	})
}
