package server

import (
	"context"
	"time"

	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	"github.com/d-strobel/gowindows/parsing"
)

// Fixtures
const (
	recordAJson = `[{"DistinguishedName":"DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local","Hostname":"test","RecordType":"A","Timestamp":null,"timetolive":{"Ticks":36000000000,"Days":0,"Hours":1,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":0.041666666666666664,"TotalHours":1,"TotalMilliseconds":3600000,"TotalMinutes":60,"TotalSeconds":3600},"RecordData":{"CimClass":"root/Microsoft/Windows/DNS:DnsServerResourceRecordA","CimInstanceProperties":"IPv4Address = \"2.2.2.2\"","CimSystemProperties":"Microsoft.Management.Infrastructure.CimSystemProperties"},"Type":1}]`

	recordExistsErr = `Fehler beim Erstellen des Ressourcendatensatzes "terratest" in der Zone "test.local" auf dem Server "DC-01".In Zeile:1 Zeichen:43
        ... yContinue'; Add-DnsServerResourceRecordA -AllowUpdateAny:$false -Crea ...
                        ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
        CategoryInfo          : ResourceExists: (terratest:root/Microsoft/...ResourceRecordA) [Add-DnsServerResourceReco rdA], CimException
        FullyQualifiedErrorId : WIN32 9711,Add-DnsServerResourceRecordA)
	`
)

var (
	expectedRecordA = RecordA{
		DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
		Name:              "test",
		Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		TimeToLive:        3600,
		Addresses:         []string{"2.2.2.2"},
	}
)

// Test the convertOutput method.
func (suite *DnsServerUnitTestSuite) TestRecordAConvertOutput() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description        string
			expectedRecordA    RecordA
			inputRecordAObject []recordObject
		}{
			{
				"should return the correct RecordA object",
				RecordA{
					Name:              "test",
					DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
					Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					TimeToLive:        3600,
					Addresses:         []string{"2.2.2.2"},
				},
				[]recordObject{
					{
						DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
						Name:              "test",
						RecordType:        "A",
						Timestamp:         parsing.DotnetTime{},
						TimeToLive:        timeToLive{Seconds: 3600},
						RecordData: recordRecordData{
							CimInstanceProperties: parsing.CimClassKeyVal{
								"IPv4Address": "2.2.2.2",
							},
						},
					},
				},
			},
			{
				"should return the correct RecordA object with multiple addresses and the lowest TTL",
				RecordA{
					Name:              "test",
					DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
					Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					TimeToLive:        60,
					Addresses:         []string{"2.2.2.2", "3.3.3.3"},
				},
				[]recordObject{
					{
						DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
						Name:              "test",
						RecordType:        "A",
						Timestamp:         parsing.DotnetTime{},
						TimeToLive:        timeToLive{Seconds: 3600},
						RecordData: recordRecordData{
							CimInstanceProperties: parsing.CimClassKeyVal{
								"IPv4Address": "2.2.2.2",
							},
						},
					},
					{
						DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
						Name:              "test",
						RecordType:        "A",
						Timestamp:         parsing.DotnetTime{},
						TimeToLive:        timeToLive{Seconds: 60},
						RecordData: recordRecordData{
							CimInstanceProperties: parsing.CimClassKeyVal{
								"IPv4Address": "3.3.3.3",
							},
						},
					},
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			r := RecordA{}
			r.convertOutput(tc.inputRecordAObject)
			suite.Equal(tc.expectedRecordA, r)
		}
	})
}

// Test RecordARead related methods.
func (suite *DnsServerUnitTestSuite) TestRecordAReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordAReadParams
			expectedCmd     string
		}{
			{
				"assert correct command A-Record read by name and zone",
				RecordAReadParams{Name: "test", Zone: "test.local"},
				"$r=Get-DnsServerResourceRecord -RRType 'A' -Node -Name 'test' -ZoneName 'test.local' ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordARead() {
	suite.Run("should return the correct A-Record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "$r=Get-DnsServerResourceRecord -RRType 'A' -Node -Name 'test' -ZoneName 'test.local' ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}").
			Return(connection.CmdResult{StdOut: recordAJson}, nil)
		actualRecordA, err := c.RecordARead(ctx, RecordAReadParams{Name: "test", Zone: "test.local"})
		suite.NoError(err)
		suite.Equal(expectedRecordA, actualRecordA)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters RecordAReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				RecordAReadParams{},
				"windows.dns.server.RecordARead: record parameters 'Name' and 'Zone' must be set",
			},
			{
				"assert error with Name only parameters",
				RecordAReadParams{Name: "tester"},
				"windows.dns.server.RecordARead: record parameters 'Name' and 'Zone' must be set",
			},
			{
				"assert error with Zone only parameters",
				RecordAReadParams{Zone: "test.local"},
				"windows.dns.server.RecordARead: record parameters 'Name' and 'Zone' must be set",
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
			_, err := c.RecordARead(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

// Test RecordACreate related methods.
func (suite *DnsServerUnitTestSuite) TestRecordACreatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordACreateParams
			expectedCmd     string
		}{
			{
				"assert without ttl parameter",
				RecordACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"1.1.1.1"}},
				"$r=Add-DnsServerResourceRecordA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 86400) -IPv4Address @('1.1.1.1') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}",
			},
			{
				"assert with multiple ip addresses",
				RecordACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}},
				"$r=Add-DnsServerResourceRecordA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 86400) -IPv4Address @('1.1.1.1','2.2.2.2','3.3.3.3') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}",
			},
			{
				"assert with ttl parameter",
				RecordACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"1.1.1.1"}, TimeToLive: 3600},
				"$r=Add-DnsServerResourceRecordA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 3600) -IPv4Address @('1.1.1.1') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordACreate() {
	suite.T().Parallel()

	suite.Run("should return the correct record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "$r=Add-DnsServerResourceRecordA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 3600) -IPv4Address @('1.1.1.1') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}").
			Return(connection.CmdResult{StdOut: recordAJson}, nil)
		actualRecord, err := c.RecordACreate(ctx, RecordACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"1.1.1.1"}, TimeToLive: 3600})
		suite.NoError(err)
		suite.Equal(expectedRecordA, actualRecord)
	})

	suite.Run("should return 'record already exists' error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "$r=Add-DnsServerResourceRecordA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 3600) -IPv4Address @('1.1.1.1') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}").
			Return(connection.CmdResult{StdErr: recordExistsErr}, nil)

		_, err := c.RecordACreate(ctx, RecordACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"1.1.1.1"}, TimeToLive: 3600})
		suite.EqualError(err, "windows.dns.server.RecordACreate: the specified record already exists.")
	})
}

// Test RecordAUpdate related methods.
func (suite *DnsServerUnitTestSuite) TestRecordAUpdatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordAUpdateParams
			expectedCmd     string
		}{
			{
				"assert without ttl parameter",
				RecordAUpdateParams{Name: "test", Zone: "test.local", TimeToLive: 3600},
				"$nr=@();Get-DnsServerResourceRecord -RRType 'A' -Node -Name 'test' -ZoneName 'test.local' | ForEach-Object{$r=$_;$n=[ciminstance]::new($r);$n.TimeToLive=New-TimeSpan -Seconds 3600 ;$nr+=Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName 'test.local' -PassThru} ;if($nr.Count -ge 2){ConvertTo-Json $nr -Compress}else{ConvertTo-Json @($nr) -Compress}",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordAUpdate() {
	suite.T().Parallel()

	suite.Run("should return the correct record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "$nr=@();Get-DnsServerResourceRecord -RRType 'A' -Node -Name 'test' -ZoneName 'test.local' | ForEach-Object{$r=$_;$n=[ciminstance]::new($r);$n.TimeToLive=New-TimeSpan -Seconds 3600 ;$nr+=Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName 'test.local' -PassThru} ;if($nr.Count -ge 2){ConvertTo-Json $nr -Compress}else{ConvertTo-Json @($nr) -Compress}").
			Return(connection.CmdResult{StdOut: recordAJson}, nil)
		actualRecord, err := c.RecordAUpdate(ctx, RecordAUpdateParams{Name: "test", Zone: "test.local", TimeToLive: 3600})
		suite.NoError(err)
		suite.Equal(expectedRecordA, actualRecord)
	})
}

// Test RecordADelete related methods.
func (suite *DnsServerUnitTestSuite) TestRecordADeletePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordADeleteParams
			expectedCmd     string
		}{
			{
				"assert with name and zone",
				RecordADeleteParams{Name: "test", Zone: "test.local"},
				"Remove-DnsServerResourceRecord -RRType 'A' -Force -Name 'test' -ZoneName 'test.local'",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordADelete() {
	suite.Run("should return the correct record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Remove-DnsServerResourceRecord -RRType 'A' -Force -Name 'test' -ZoneName 'test.local'").
			Return(connection.CmdResult{StdOut: recordAJson}, nil)
		err := c.RecordADelete(ctx, RecordADeleteParams{Name: "test", Zone: "test.local"})
		suite.NoError(err)
	})
}
