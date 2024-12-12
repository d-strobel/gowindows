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
	recordCNameJson = `{"DistinguishedName":"DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local","Hostname":"test","RecordType":"CName","Timestamp":null,"timetolive":{"Ticks":36000000000,"Days":0,"Hours":1,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":0.041666666666666664,"TotalHours":1,"TotalMilliseconds":3600000,"TotalMinutes":60,"TotalSeconds":3600},"RecordData":{"CimClass":"root/Microsoft/Windows/DNS:DnsServerResourceRecordCName","CimInstanceProperties":["HostNameAlias = \"testalias\""],"CimSystemProperties":"Microsoft.Management.Infrastructure.CimSystemProperties"},"Type":1}`

	recordCNameExistsErr = `Fehler beim Erstellen des Ressourcendatensatzes "terratest" in der Zone "test.local" auf dem Server "DC-01".In Zeile:1 Zeichen:43
        ... yContinue'; Add-DnsServerResourceRecordCName -AllowUpdateAny:$false -Crea ...
                        ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
        CategoryInfo          : ResourceExists: (terratest:root/Microsoft/...ResourceRecordCName) [Add-DnsServerResourceReco rdA], CimException
        FullyQualifiedErrorId : WIN32 9711,Add-DnsServerResourceRecordCName)
	`
)

var (
	expectedRecordCName = RecordCName{
		DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
		Name:              "test",
		Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		TimeToLive:        time.Second * 3600,
		CName:             "testalias",
	}
)

// Test the convertOutput method.
func (suite *DnsServerUnitTestSuite) TestRecordCNameConvertOutput() {
	suite.Run("should return the correct object", func() {
		tcs := []struct {
			description            string
			expectedRecordCName    RecordCName
			inputRecordCNameObject recordObject
		}{
			{
				"should return the correct RecordCName object",
				RecordCName{
					Name:              "test",
					DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
					Timestamp:         time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					TimeToLive:        time.Second * 3600,
					CName:             "testalias",
				},
				recordObject{
					DistinguishedName: "DC=test,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local",
					Name:              "test",
					RecordType:        "CName",
					Timestamp:         parsing.DotnetTime{},
					TimeToLive:        parsing.CimTimeDuration{Duration: time.Second * 3600},
					RecordData: recordRecordData{
						CimInstanceProperties: parsing.CimClassKeyVal{
							"HostNameAlias": "testalias",
						},
					},
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			r := RecordCName{}
			r.convertOutput(tc.inputRecordCNameObject)
			suite.Equal(tc.expectedRecordCName, r)
		}
	})
}

// Test RecordCNameRead related methods.
func (suite *DnsServerUnitTestSuite) TestRecordCNameReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordCNameReadParams
			expectedCmd     string
		}{
			{
				"assert correct command A-Record read by name and zone",
				RecordCNameReadParams{Name: "test", Zone: "test.local"},
				"Get-DnsServerResourceRecord -RRType 'CName' -Node -Name 'test' -ZoneName 'test.local' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordCNameRead() {
	suite.Run("should return the correct CName-Record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-DnsServerResourceRecord -RRType 'CName' -Node -Name 'test' -ZoneName 'test.local' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: recordCNameJson}, nil)
		actualRecordCName, err := c.RecordCNameRead(ctx, RecordCNameReadParams{Name: "test", Zone: "test.local"})
		suite.NoError(err)
		suite.Equal(expectedRecordCName, actualRecordCName)
	})
}

// Test RecordCNameCreate related methods.
func (suite *DnsServerUnitTestSuite) TestRecordCNameCreatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordCNameCreateParams
			expectedCmd     string
		}{
			{
				"assert with default ttl parameter",
				RecordCNameCreateParams{Name: "test", Zone: "test.local", CName: "testalias"},
				"Add-DnsServerResourceRecordCName -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -HostNameAlias 'testalias' -TimeToLive $(New-TimeSpan -Seconds 86400) | ConvertTo-Json -Compress",
			},
			{
				"assert with ttl parameter",
				RecordCNameCreateParams{Name: "test", Zone: "test.local", TimeToLive: time.Second * 3600, CName: "testalias"},
				"Add-DnsServerResourceRecordCName -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -HostNameAlias 'testalias' -TimeToLive $(New-TimeSpan -Seconds 3600) | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordCNameCreate() {
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
			RunWithPowershell(ctx, "Add-DnsServerResourceRecordCName -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -HostNameAlias 'testalias' -TimeToLive $(New-TimeSpan -Seconds 3600) | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: recordCNameJson}, nil)
		actualRecord, err := c.RecordCNameCreate(ctx, RecordCNameCreateParams{Name: "test", Zone: "test.local", CName: "testalias", TimeToLive: time.Second * 3600})
		suite.NoError(err)
		suite.Equal(expectedRecordCName, actualRecord)
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
			RunWithPowershell(ctx, "Add-DnsServerResourceRecordCName -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru -Name 'test' -ZoneName 'test.local' -HostNameAlias 'testalias' -TimeToLive $(New-TimeSpan -Seconds 3600) | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdErr: recordExistsErr}, nil)

		_, err := c.RecordCNameCreate(ctx, RecordCNameCreateParams{Name: "test", Zone: "test.local", CName: "testalias", TimeToLive: time.Second * 3600})
		suite.EqualError(err, "windows.dns.server.RecordCNameCreate: the specified record already exists.")
	})
}

// Test RecordCNameUpdate related methods.
func (suite *DnsServerUnitTestSuite) TestRecordCNameUpdatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordCNameUpdateParams
			expectedCmd     string
		}{
			{
				"assert without ttl parameter",
				RecordCNameUpdateParams{Name: "test", Zone: "test.local", TimeToLive: time.Second * 3600, CName: "testalias"},
				"$r=Get-DnsServerResourceRecord -RRType 'CName' -Node -Name 'test' -ZoneName 'test.local' ;$n=[ciminstance]::new($r) ;$n.TimeToLive=New-TimeSpan -Seconds 3600 ;$n.RecordData.HostNameAlias='testalias' ;Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName 'test.local' -PassThru | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordCNameUpdate() {
	suite.Run("should return the correct updated record", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "$r=Get-DnsServerResourceRecord -RRType 'CName' -Node -Name 'test' -ZoneName 'test.local' ;$n=[ciminstance]::new($r) ;$n.TimeToLive=New-TimeSpan -Seconds 3600 ;$n.RecordData.HostNameAlias='testalias' ;Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName 'test.local' -PassThru | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: recordCNameJson}, nil)
		actualRecord, err := c.RecordCNameUpdate(ctx, RecordCNameUpdateParams{Name: "test", Zone: "test.local", CName: "testalias", TimeToLive: time.Second * 3600})
		suite.NoError(err)
		suite.Equal(expectedRecordCName, actualRecord)
	})
}

// Test RecordCNameDelete related methods.
func (suite *DnsServerUnitTestSuite) TestRecordCNameDeletePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters RecordCNameDeleteParams
			expectedCmd     string
		}{
			{
				"assert with name and zone",
				RecordCNameDeleteParams{Name: "test", Zone: "test.local"},
				"Remove-DnsServerResourceRecord -RRType 'CName' -Force -Name 'test' -ZoneName 'test.local'",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DnsServerUnitTestSuite) TestRecordCNameDelete() {
	suite.Run("should return no error after deletion", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return s, nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Remove-DnsServerResourceRecord -RRType 'CName' -Force -Name 'test' -ZoneName 'test.local'").
			Return(connection.CmdResult{}, nil)
		err := c.RecordCNameDelete(ctx, RecordCNameDeleteParams{Name: "test", Zone: "test.local"})
		suite.NoError(err)
	})
}
