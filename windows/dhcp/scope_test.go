package dhcp

import (
	"context"
	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	"net/netip"
	"time"

	"github.com/d-strobel/gowindows/parsing"
)

// Fixtures
const (
	scopeV4Json = `{"ScopeId":{"Address":698560,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.10.0"},"SubnetMask":{"Address":16777215,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"255.255.255.0"},"StartRange":{"Address":84584640,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.10.5"},"EndRange":{"Address":168470720,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.10.10"},"ActivatePolicies":true,"Delay":0,"Description":"Test description","LeaseDuration":{"Ticks":6912000000000,"Days":8,"Hours":0,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":8,"TotalHours":192,"TotalMilliseconds":691200000,"TotalMinutes":11520,"TotalSeconds":691200},"MaxBootpClients":4294967295,"Name":"test","NapEnable":false,"NapProfile":"","State":"Active","SuperscopeName":"","Type":"Dhcp","PSComputerName":null}`
)

var (
	expectedScopeObject = scopeObject{
		Name:        "test",
		Description: "Test description",
		StartRange: startRange{
			Address: netip.MustParseAddr("192.168.10.5"),
		},
		EndRange: endRange{
			netip.MustParseAddr("192.168.10.10"),
		},
		SubnetMask: subnetMask{
			netip.MustParseAddr("255.255.255.0"),
		},
		State:            "Active",
		MaxBootpClients:  4294967295,
		ActivatePolicies: true,
		NapEnable:        false,
		NapProfile:       "",
		Delay:            0,
		LeaseDuration: parsing.CimTimeDuration{
			Duration: time.Hour * 24 * 8,
		},
	}
	expectedScopeV4 = ScopeV4{
		Name:             "test",
		Description:      "Test description",
		StartRange:       netip.MustParseAddr("192.168.10.5"),
		EndRange:         netip.MustParseAddr("192.168.10.10"),
		SubnetMask:       netip.MustParseAddr("255.255.255.0"),
		Enabled:          true,
		MaxBootpClients:  4294967295,
		ActivatePolicies: true,
		NapEnable:        false,
		NapProfile:       "",
		Delay:            0,
		LeaseDuration:    time.Hour * 24 * 8,
	}
)

// Test the convertOutput method.
func (suite *DhcpServerUnitTestSuite) TestRecordAConvertOutput() {
	suite.Run("should return the correct command", func() {
		s := ScopeV4{}
		s.convertOutput(expectedScopeObject)
		suite.Equal(expectedScopeV4, s)
	})
}

// Test ScopeV4Read related methods.
func (suite *DhcpServerUnitTestSuite) TestScopeV4ReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4ReadParams
			expectedCmd     string
		}{
			{
				"assert correct command ScopeV4 read by ScopeId",
				ScopeV4ReadParams{ScopeId: netip.MustParseAddr("192.168.10.0")},
				"Get-DhcpServerv4Scope -ScopeId '192.168.10.0' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DhcpServerUnitTestSuite) TestScopeV4Read() {
	suite.Run("should return the correct ScopeV4", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-DhcpServerv4Scope -ScopeId '192.168.10.0' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: scopeV4Json}, nil)
		actualScopeV4, err := c.ScopeV4Read(ctx, ScopeV4ReadParams{
			ScopeId: netip.MustParseAddr("192.168.10.0"),
		})
		suite.NoError(err)
		suite.Equal(expectedScopeV4, actualScopeV4)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4ReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				ScopeV4ReadParams{},
				"windows.dhcp.ScopeV4Read: scope parameter 'ScopeId' must be an IPv4 network address",
			},
			{
				"assert error with empty parameters",
				ScopeV4ReadParams{ScopeId: netip.MustParseAddr("fe80:0010::")},
				"windows.dhcp.ScopeV4Read: scope parameter 'ScopeId' must be an IPv4 network address",
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
			_, err := c.ScopeV4Read(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}
