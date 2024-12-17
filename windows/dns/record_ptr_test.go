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
	recordPTRJson = `{"DistinguishedName":"DC=1,DC=10.168.192.in-addr.arpa,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local","Hostname":"1","RecordType":"PTR","Timestamp":null,"timetolive":{"Ticks":36000000000,"Days":0,"Hours":1,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":0.041666666666666664,"TotalHours":1,"TotalMilliseconds":3600000,"TotalMinutes":60,"TotalSeconds":3600},"RecordData":{"CimClass":"root/Microsoft/Windows/DNS:DnsServerResourceRecordPTR","CimInstanceProperties":["PtrDomainName = \"testptr.test.local.\""],"CimSystemProperties":"Microsoft.Management.Infrastructure.CimSystemProperties"},"Type":1}`

	recordPTRExistsErr = `Fehler beim Erstellen des Ressourcendatensatzes "terratest" in der Zone "test.local" auf dem Server "DC-01".In Zeile:1 Zeichen:43
        ... yContinue'; Add-DnsServerResourceRecordPTR -AllowUpdateAny:$false -Crea ...
                        ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
        CategoryInfo          : ResourceExists: (terratest:root/Microsoft/...ResourceRecordPTR) [Add-DnsServerResourceReco rdA], CimException
        FullyQualifiedErrorId : WIN32 9711,Add-DnsServerResourceRecordPTR)
	`
)

var (
	expectedRecordPTR = RecordPTR{
		DistinguishedName: "DC=1,DC=10.168.192.in-addr.arpa,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
		Name:              "1",
		Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		TimeToLive:        time.Second * 3600,
		PTR:               "testptr.test.local.",
	}
)

// Test the convertOutput method.
func (suite *DnsServerUnitTestSuite) TestRecordPTRConvertOutput() {
	suite.Run("should return the correct object", func() {
		tcs := []struct {
			description          string
			expectedRecordPTR    RecordPTR
			inputRecordPTRObject recordObject
		}{
			{
				"should return the correct RecordPTR object",
				RecordPTR{
					Name:              "1",
					DistinguishedName: "DC=1,DC=10.168.192.in-addr.arpa,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
					Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					TimeToLive:        time.Second * 3600,
					PTR:               "testptr.test.local.",
				},
				recordObject{
					DistinguishedName: "DC=1,DC=10.168.192.in-addr.arpa,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
					Name:              "1",
					RecordType:        "PTR",
					Timestamp:         parsing.DotnetTime{},
					TimeToLive:        parsing.CimTimeDuration{Duration: time.Second * 3600},
					RecordData: recordRecordData{
						CimInstanceProperties: parsing.CimClassKeyVal{
							"PtrDomainName": "testptr.test.local.",
						},
					},
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			r := RecordPTR{}
			r.convertOutput(tc.inputRecordPTRObject)
			suite.Equal(tc.expectedRecordPTR, r)
		}
	})
}

// Test RecordPTRRead related methods.
func (suite *DnsServerUnitTestSuite) TestRecordPTRReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordPTRReadParams
			expectedCmd     string
		}{
			{
				"assert correct command by name and zone",
				RecordPTRReadParams{Name: "1", Zone: "10.168.192.in-addr.arpa"},
				"Get-DnsServerResourceRecord -RRType 'PTR' -Node -Name '1' -ZoneName '10.168.192.in-addr.arpa' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordPTRRead() {
	suite.Run("should return the correct PTR-Record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-DnsServerResourceRecord -RRType 'PTR' -Node -Name '1' -ZoneName '10.168.192.in-addr.arpa' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: recordPTRJson}, nil)
		actualRecordPTR, err := c.RecordPTRRead(ctx, RecordPTRReadParams{Name: "1", Zone: "10.168.192.in-addr.arpa"})
		suite.NoError(err)
		suite.Equal(expectedRecordPTR, actualRecordPTR)
	})
}

// Test RecordPTRCreate related methods.
func (suite *DnsServerUnitTestSuite) TestRecordPTRCreatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordPTRCreateParams
			expectedCmd     string
		}{
			{
				"assert with default ttl parameter",
				RecordPTRCreateParams{Name: "1", Zone: "10.168.192.in-addr.arpa", PTR: "testptr.test.local."},
				"Add-DnsServerResourceRecordPTR -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru -Name '1' -ZoneName '10.168.192.in-addr.arpa' -PtrDomainName 'testptr.test.local.' -TimeToLive $(New-TimeSpan -Seconds 86400) | ConvertTo-Json -Compress",
			},
			{
				"assert with ttl parameter",
				RecordPTRCreateParams{Name: "1", Zone: "10.168.192.in-addr.arpa", TimeToLive: time.Second * 3600, PTR: "testptr.test.local."},
				"Add-DnsServerResourceRecordPTR -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru -Name '1' -ZoneName '10.168.192.in-addr.arpa' -PtrDomainName 'testptr.test.local.' -TimeToLive $(New-TimeSpan -Seconds 3600) | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordPTRCreate() {
	suite.Run("should return the correct record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Add-DnsServerResourceRecordPTR -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru -Name '1' -ZoneName '10.168.192.in-addr.arpa' -PtrDomainName 'testptr.test.local.' -TimeToLive $(New-TimeSpan -Seconds 3600) | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: recordPTRJson}, nil)
		actualRecord, err := c.RecordPTRCreate(ctx, RecordPTRCreateParams{Name: "1", Zone: "10.168.192.in-addr.arpa", PTR: "testptr.test.local.", TimeToLive: time.Second * 3600})
		suite.NoError(err)
		suite.Equal(expectedRecordPTR, actualRecord)
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
			RunWithPowershell(ctx, "Add-DnsServerResourceRecordPTR -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru -Name '1' -ZoneName '10.168.192.in-addr.arpa' -PtrDomainName 'testptr.test.local.' -TimeToLive $(New-TimeSpan -Seconds 3600) | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdErr: recordExistsErr}, nil)

		_, err := c.RecordPTRCreate(ctx, RecordPTRCreateParams{Name: "1", Zone: "10.168.192.in-addr.arpa", PTR: "testptr.test.local.", TimeToLive: time.Second * 3600})
		suite.EqualError(err, "windows.dns.RecordPTRCreate: the specified record already exists")
	})
}

// Test RecordPTRUpdate related methods.
func (suite *DnsServerUnitTestSuite) TestRecordPTRUpdatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordPTRUpdateParams
			expectedCmd     string
		}{
			{
				"assert with ttl parameter",
				RecordPTRUpdateParams{Name: "1", Zone: "10.168.192.in-addr.arpa", TimeToLive: time.Second * 3600, PTR: "testptr.test.local."},
				"$r=Get-DnsServerResourceRecord -RRType 'PTR' -Node -Name '1' -ZoneName '10.168.192.in-addr.arpa' ;$n=[ciminstance]::new($r) ;$n.TimeToLive=New-TimeSpan -Seconds 3600 ;$n.RecordData.PtrDomainName='testptr.test.local.' ;Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName '10.168.192.in-addr.arpa' -PassThru | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordPTRUpdate() {
	suite.Run("should return the correct updated record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "$r=Get-DnsServerResourceRecord -RRType 'PTR' -Node -Name '1' -ZoneName '10.168.192.in-addr.arpa' ;$n=[ciminstance]::new($r) ;$n.TimeToLive=New-TimeSpan -Seconds 3600 ;$n.RecordData.PtrDomainName='testptr.test.local.' ;Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName '10.168.192.in-addr.arpa' -PassThru | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: recordPTRJson}, nil)
		actualRecord, err := c.RecordPTRUpdate(ctx, RecordPTRUpdateParams{Name: "1", Zone: "10.168.192.in-addr.arpa", TimeToLive: time.Second * 3600, PTR: "testptr.test.local."})
		suite.NoError(err)
		suite.Equal(expectedRecordPTR, actualRecord)
	})
}

// Test RecordPTRDelete related methods.
func (suite *DnsServerUnitTestSuite) TestRecordPTRDeletePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordPTRDeleteParams
			expectedCmd     string
		}{
			{
				"assert with name and zone",
				RecordPTRDeleteParams{Name: "1", Zone: "10.168.192.in-addr.arpa"},
				"Remove-DnsServerResourceRecord -RRType 'PTR' -Force -Name '1' -ZoneName '10.168.192.in-addr.arpa'",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordPTRDelete() {
	suite.Run("should return no error after deletion", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Remove-DnsServerResourceRecord -RRType 'PTR' -Force -Name '1' -ZoneName '10.168.192.in-addr.arpa'").
			Return(connection.CmdResult{}, nil)
		err := c.RecordPTRDelete(ctx, RecordPTRDeleteParams{Name: "1", Zone: "10.168.192.in-addr.arpa"})
		suite.NoError(err)
	})
}
