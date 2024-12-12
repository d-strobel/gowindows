package dns

import (
	"context"
	"errors"

	"github.com/d-strobel/gowindows/connection"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
)

// Fixtures
const (
	zone     = `{"NotifyServers":null,"SecondaryServers":null,"AllowedDcForNsRecordsAutoCreation":null,"DistinguishedName":"DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local","IsAutoCreated":false,"IsDsIntegrated":true,"IsPaused":false,"IsReadOnly":false,"IsReverseLookupZone":false,"IsShutdown":false,"ZoneName":"test.local","ZoneType":"Primary","DirectoryPartitionName":"DomainDnsZones.test.local","DynamicUpdate":"Secure","IgnorePolicies":false,"IsSigned":false,"IsWinsEnabled":false,"Notify":"NotifyServers","ReplicationScope":"Domain","SecureSecondaries":"NoTransfer","ZoneFile":null,"PSComputerName":null}`
	zoneList = `[{"NotifyServers":null,"SecondaryServers":null,"AllowedDcForNsRecordsAutoCreation":null,"DistinguishedName":"DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local","IsAutoCreated":false,"IsDsIntegrated":true,"IsPaused":false,"IsReadOnly":false,"IsReverseLookupZone":false,"IsShutdown":false,"ZoneName":"test.local","ZoneType":"Primary","DirectoryPartitionName":"DomainDnsZones.test.local","DynamicUpdate":"Secure","IgnorePolicies":false,"IsSigned":false,"IsWinsEnabled":false,"Notify":"NotifyServers","ReplicationScope":"Domain","SecureSecondaries":"NoTransfer","ZoneFile":null,"PSComputerName":null},{"NotifyServers":null,"SecondaryServers":null,"AllowedDcForNsRecordsAutoCreation":null,"DistinguishedName":"DC=test2.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test2,DC=local","IsAutoCreated":false,"IsDsIntegrated":true,"IsPaused":false,"IsReadOnly":false,"IsReverseLookupZone":false,"IsShutdown":false,"ZoneName":"test2.local","ZoneType":"Primary","DirectoryPartitionName":"DomainDnsZones.test2.local","DynamicUpdate":"Secure","IgnorePolicies":false,"IsSigned":false,"IsWinsEnabled":false,"Notify":"NotifyServers","ReplicationScope":"Domain","SecureSecondaries":"NoTransfer","ZoneFile":null,"PSComputerName":null}]`
)

var (
	expectedZone = Zone{
		NotifyServers:                     "",
		SecondaryServers:                  "",
		AllowedDcForNsRecordsAutoCreation: "",
		DistinguishedName:                 "DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
		IsAutoCreated:                     false,
		IsDsIntegrated:                    true,
		IsPaused:                          false,
		IsReadOnly:                        false,
		IsReverseLookupZone:               false,
		IsShutdown:                        false,
		ZoneName:                          "test.local",
		ZoneType:                          "Primary",
		DirectoryPartitionName:            "DomainDnsZones.test.local",
		DynamicUpdate:                     "Secure",
		IgnorePolicies:                    false,
		IsSigned:                          false,
		IsWinsEnabled:                     false,
		Notify:                            "NotifyServers",
		ReplicationScope:                  "Domain",
		SecureSecondaries:                 "NoTransfer",
		ZoneFile:                          "",
	}
	expectedZoneList = []Zone{
		{
			NotifyServers:                     "",
			SecondaryServers:                  "",
			AllowedDcForNsRecordsAutoCreation: "",
			DistinguishedName:                 "DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
			IsAutoCreated:                     false,
			IsDsIntegrated:                    true,
			IsPaused:                          false,
			IsReadOnly:                        false,
			IsReverseLookupZone:               false,
			IsShutdown:                        false,
			ZoneName:                          "test.local",
			ZoneType:                          "Primary",
			DirectoryPartitionName:            "DomainDnsZones.test.local",
			DynamicUpdate:                     "Secure",
			IgnorePolicies:                    false,
			IsSigned:                          false,
			IsWinsEnabled:                     false,
			Notify:                            "NotifyServers",
			ReplicationScope:                  "Domain",
			SecureSecondaries:                 "NoTransfer",
			ZoneFile:                          "",
		},
		{
			NotifyServers:                     "",
			SecondaryServers:                  "",
			AllowedDcForNsRecordsAutoCreation: "",
			DistinguishedName:                 "DC=test2.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test2,DC=local",
			IsAutoCreated:                     false,
			IsDsIntegrated:                    true,
			IsPaused:                          false,
			IsReadOnly:                        false,
			IsReverseLookupZone:               false,
			IsShutdown:                        false,
			ZoneName:                          "test2.local",
			ZoneType:                          "Primary",
			DirectoryPartitionName:            "DomainDnsZones.test2.local",
			DynamicUpdate:                     "Secure",
			IgnorePolicies:                    false,
			IsSigned:                          false,
			IsWinsEnabled:                     false,
			Notify:                            "NotifyServers",
			ReplicationScope:                  "Domain",
			SecureSecondaries:                 "NoTransfer",
			ZoneFile:                          "",
		},
	}
)

// Test ZoneRead related methods.
func (suite *DnsServerUnitTestSuite) TestZoneReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters ZoneReadParams
			expectedCmd     string
		}{
			{
				"assert zone read by name",
				ZoneReadParams{Name: "test.local"},
				"Get-DnsServerZone -Name 'test.local' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestZoneRead() {
	suite.Run("should return the correct zone", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-DnsServerZone -Name 'test.local' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: zone}, nil)
		actualZone, err := c.ZoneRead(ctx, ZoneReadParams{Name: "test.local"})
		suite.NoError(err)
		suite.Equal(expectedZone, actualZone)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters ZoneReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				ZoneReadParams{},
				"windows.dns.server.ZoneRead: zone parameter 'Name' must be set",
			},
			{
				"assert error when name contains wildcard",
				ZoneReadParams{Name: "*.local"},
				"windows.dns.server.ZoneRead: zone parameter 'Name' does not allow wildcards",
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
			_, err := c.ZoneRead(ctx, tc.inputParameters)
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
			RunWithPowershell(ctx, "Get-DnsServerZone -Name 'test.local' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{}, errors.New("test-error"))
		_, err := c.ZoneRead(ctx, ZoneReadParams{Name: "test.local"})
		suite.EqualError(err, "windows.dns.server.ZoneRead: test-error")
	})
}

// Test ZoneList related methods.
func (suite *DnsServerUnitTestSuite) TestZoneList() {
	suite.Run("should return the correct list of zones", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Get-DnsServerZone | ConvertTo-Json -Compress"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{StdOut: zoneList}, nil)
		actualZoneList, err := c.ZoneList(ctx)
		suite.NoError(err)
		suite.Equal(expectedZoneList, actualZoneList)
	})

	suite.Run("should return error if run fails", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		cmd := "Get-DnsServerZone | ConvertTo-Json -Compress"
		mockConn.EXPECT().
			RunWithPowershell(ctx, cmd).
			Return(connection.CmdResult{}, errors.New("test-error"))
		_, err := c.ZoneList(ctx)
		suite.EqualError(err, "windows.dns.server.ZoneList: test-error")
	})
}
