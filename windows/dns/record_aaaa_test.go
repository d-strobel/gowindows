package dns

import (
	"context"
	"time"

	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	"github.com/d-strobel/gowindows/parsing"
)

// Fixtures
const (
	recordAAAAJson = `[{"DistinguishedName":"DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local","Hostname":"test","RecordType":"A","Timestamp":null,"timetolive":{"Ticks":36000000000,"Days":0,"Hours":1,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":0.041666666666666664,"TotalHours":1,"TotalMilliseconds":3600000,"TotalMinutes":60,"TotalSeconds":3600},"RecordData":{"CimClass":"root/Microsoft/Windows/DNS:DnsServerResourceRecordAAAA","CimInstanceProperties":"IPv6Address = \"2001:db8:0000::1\"","CimSystemProperties":"Microsoft.Management.Infrastructure.CimSystemProperties"},"Type":1}]`

	recordAAAAExistsErr = `Add-DnsServerResourceRecordAAAA : Fehler beim Erstellen des Ressourcendatensatzes "terratest" in der Zone "test.local" auf dem Server "DC-01".In Zeile:1 Zeichen:46
        ... ntinue'; $r=Add-DnsServerResourceRecordAAAA -AllowUpdateAny:$false -C ...
                        ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
        CategoryInfo          : ResourceExists: (testo2:root/Microsoft/...ourceRecordAAAA) [Add-DnsServerResourceRecordA AAA], CimException
        FullyQualifiedErrorId : WIN32 9711,Add-DnsServerResourceRecordAAAA
	`
)

var (
	expectedRecordAAAA = RecordAAAA{
		DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
		Name:              "test",
		Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		TimeToLive:        time.Second * 3600,
		Addresses:         []string{"2001:db8:0000::1"},
	}
)

// Test the convertOutput method.
func (suite *DnsServerUnitTestSuite) TestRecordAAAAConvertOutput() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description           string
			expectedRecordAAAA    RecordAAAA
			inputRecordAAAAObject []recordObject
		}{
			{
				"should return the correct RecordAAAA object",
				RecordAAAA{
					Name:              "test",
					DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
					Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					TimeToLive:        time.Second * 3600,
					Addresses:         []string{"2001:db8:0000::1"},
				},
				[]recordObject{
					{
						DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
						Name:              "test",
						RecordType:        "AAAA",
						Timestamp:         parsing.DotnetTime{},
						TimeToLive:        parsing.CimTimeDuration{Duration: time.Second * 3600},
						RecordData: recordRecordData{
							CimInstanceProperties: parsing.CimClassKeyVal{
								"IPv6Address": "2001:db8:0000::1",
							},
						},
					},
				},
			},
			{
				"should return the correct RecordAAAA object with multiple addresses and the lowest TTL",
				RecordAAAA{
					Name:              "test",
					DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
					Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					TimeToLive:        time.Second * 60,
					Addresses:         []string{"2001:db8:0000::1", "2001:db8:0000::2"},
				},
				[]recordObject{
					{
						DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
						Name:              "test",
						RecordType:        "AAAA",
						Timestamp:         parsing.DotnetTime{},
						TimeToLive:        parsing.CimTimeDuration{Duration: time.Second * 3600},
						RecordData: recordRecordData{
							CimInstanceProperties: parsing.CimClassKeyVal{
								"IPv6Address": "2001:db8:0000::1",
							},
						},
					},
					{
						DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
						Name:              "test",
						RecordType:        "AAAA",
						Timestamp:         parsing.DotnetTime{},
						TimeToLive:        parsing.CimTimeDuration{Duration: time.Second * 60},
						RecordData: recordRecordData{
							CimInstanceProperties: parsing.CimClassKeyVal{
								"IPv6Address": "2001:db8:0000::2",
							},
						},
					},
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			r := RecordAAAA{}
			r.convertOutput(tc.inputRecordAAAAObject)
			suite.Equal(tc.expectedRecordAAAA, r)
		}
	})
}

// Test RecordAAAARead related methods.
func (suite *DnsServerUnitTestSuite) TestRecordAAAAReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordAAAAReadParams
			expectedCmd     string
		}{
			{
				"assert correct command A-Record read by name and zone",
				RecordAAAAReadParams{Name: "test", Zone: "test.local"},
				"$r=Get-DnsServerResourceRecord -RRType 'AAAA' -Node -Name 'test' -ZoneName 'test.local' ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordAAAARead() {
	suite.Run("should return the correct A-Record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "$r=Get-DnsServerResourceRecord -RRType 'AAAA' -Node -Name 'test' -ZoneName 'test.local' ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}").
			Return(connection.CmdResult{StdOut: recordAAAAJson}, nil)
		actualRecordAAAA, err := c.RecordAAAARead(ctx, RecordAAAAReadParams{Name: "test", Zone: "test.local"})
		suite.NoError(err)
		suite.Equal(expectedRecordAAAA, actualRecordAAAA)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters RecordAAAAReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				RecordAAAAReadParams{},
				"windows.dns.server.RecordAAAARead: record parameters 'Name' and 'Zone' must be set",
			},
			{
				"assert error with Name only parameters",
				RecordAAAAReadParams{Name: "tester"},
				"windows.dns.server.RecordAAAARead: record parameters 'Name' and 'Zone' must be set",
			},
			{
				"assert error with Zone only parameters",
				RecordAAAAReadParams{Zone: "test.local"},
				"windows.dns.server.RecordAAAARead: record parameters 'Name' and 'Zone' must be set",
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
			_, err := c.RecordAAAARead(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

// Test RecordAAAACreate related methods.
func (suite *DnsServerUnitTestSuite) TestRecordAAAACreatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordAAAACreateParams
			expectedCmd     string
		}{
			{
				"assert without ttl parameter",
				RecordAAAACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"2001:db8:0000::1"}},
				"$r=Add-DnsServerResourceRecordAAAA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 86400) -IPv6Address @('2001:db8:0000::1') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}",
			},
			{
				"assert with multiple ip addresses",
				RecordAAAACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"2001:db8:0000::1", "2001:db8:0000::2", "2001:db8:0000::3"}},
				"$r=Add-DnsServerResourceRecordAAAA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 86400) -IPv6Address @('2001:db8:0000::1','2001:db8:0000::2','2001:db8:0000::3') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}",
			},
			{
				"assert with ttl parameter",
				RecordAAAACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"2001:db8:0000::1"}, TimeToLive: time.Second * 3600},
				"$r=Add-DnsServerResourceRecordAAAA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 3600) -IPv6Address @('2001:db8:0000::1') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordAAAACreate() {
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
			RunWithPowershell(ctx, "$r=Add-DnsServerResourceRecordAAAA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 3600) -IPv6Address @('2001:db8:0000::1') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}").
			Return(connection.CmdResult{StdOut: recordAAAAJson}, nil)
		actualRecord, err := c.RecordAAAACreate(ctx, RecordAAAACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"2001:db8:0000::1"}, TimeToLive: time.Second * 3600})
		suite.NoError(err)
		suite.Equal(expectedRecordAAAA, actualRecord)
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
			RunWithPowershell(ctx, "$r=Add-DnsServerResourceRecordAAAA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -TimeToLive $(New-TimeSpan -Seconds 3600) -IPv6Address @('2001:db8:0000::1') ;if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}").
			Return(connection.CmdResult{StdErr: recordAAAAExistsErr}, nil)

		_, err := c.RecordAAAACreate(ctx, RecordAAAACreateParams{Name: "test", Zone: "test.local", Addresses: []string{"2001:db8:0000::1"}, TimeToLive: time.Second * 3600})
		suite.EqualError(err, "windows.dns.server.RecordAAAACreate: the specified record already exists.")
	})
}

// Test RecordAAAAUpdate related methods.
func (suite *DnsServerUnitTestSuite) TestRecordAAAAUpdatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordAAAAUpdateParams
			expectedCmd     string
		}{
			{
				"assert without ttl parameter",
				RecordAAAAUpdateParams{Name: "test", Zone: "test.local", TimeToLive: time.Second * 3600},
				"$nr=@();Get-DnsServerResourceRecord -RRType 'AAAA' -Node -Name 'test' -ZoneName 'test.local' | ForEach-Object{$r=$_;$n=[ciminstance]::new($r);$n.TimeToLive=New-TimeSpan -Seconds 3600 ;$nr+=Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName 'test.local' -PassThru} ;if($nr.Count -ge 2){ConvertTo-Json $nr -Compress}else{ConvertTo-Json @($nr) -Compress}",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordAAAAUpdate() {
	suite.Run("should return the correct record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "$nr=@();Get-DnsServerResourceRecord -RRType 'AAAA' -Node -Name 'test' -ZoneName 'test.local' | ForEach-Object{$r=$_;$n=[ciminstance]::new($r);$n.TimeToLive=New-TimeSpan -Seconds 3600 ;$nr+=Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName 'test.local' -PassThru} ;if($nr.Count -ge 2){ConvertTo-Json $nr -Compress}else{ConvertTo-Json @($nr) -Compress}").
			Return(connection.CmdResult{StdOut: recordAAAAJson}, nil)
		actualRecord, err := c.RecordAAAAUpdate(ctx, RecordAAAAUpdateParams{Name: "test", Zone: "test.local", TimeToLive: time.Second * 3600})
		suite.NoError(err)
		suite.Equal(expectedRecordAAAA, actualRecord)
	})
}

// Test RecordAAAADelete related methods.
func (suite *DnsServerUnitTestSuite) TestRecordAAAADeletePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordAAAADeleteParams
			expectedCmd     string
		}{
			{
				"assert with name and zone",
				RecordAAAADeleteParams{Name: "test", Zone: "test.local"},
				"Remove-DnsServerResourceRecord -RRType 'AAAA' -Force -Name 'test' -ZoneName 'test.local'",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordAAAADelete() {
	suite.Run("should return the correct record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Remove-DnsServerResourceRecord -RRType 'AAAA' -Force -Name 'test' -ZoneName 'test.local'").
			Return(connection.CmdResult{StdOut: recordAAAAJson}, nil)
		err := c.RecordAAAADelete(ctx, RecordAAAADeleteParams{Name: "test", Zone: "test.local"})
		suite.NoError(err)
	})
}
