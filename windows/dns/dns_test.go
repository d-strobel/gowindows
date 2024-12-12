package dns

import (
	"context"
	"testing"

	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	"github.com/stretchr/testify/suite"
)

// Unit test suite for all local functions
type DnsServerUnitTestSuite struct {
	suite.Suite
}

// Run all local unit tests
func TestDnsServerUnitTestSuite(t *testing.T) {
	suite.Run(t, &DnsServerUnitTestSuite{})
}

func (suite *DnsServerUnitTestSuite) TestNewClient() {
	suite.Run("should return a new dns server client", func() {
		mockConn := mockConnection.NewMockConnection(suite.T())
		mockDecodeCliXmlErr := func(s string) (string, error) { return "", nil }
		actualClient := NewClientWithParser(mockConn, mockDecodeCliXmlErr)
		expectedClient := &Client{Connection: mockConn, decodeCliXmlErr: mockDecodeCliXmlErr}
		suite.IsType(expectedClient, actualClient)
		suite.Equal(expectedClient.Connection, actualClient.Connection)
	})
}

func (suite *DnsServerUnitTestSuite) TestDnsServerRun() {
	suite.T().Parallel()

	suite.Run("should return a zone", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Get-DnsServerZone -Name Test | ConvertTo-Json -Compress"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{StdOut: zone, StdErr: ""}, nil)
		var z Zone
		err := run(ctx, c, cmd, &z)
		suite.NoError(err)
		suite.Equal(expectedZone, z)
	})
}
