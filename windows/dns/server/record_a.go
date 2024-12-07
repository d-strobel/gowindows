package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/d-strobel/gowindows/winerror"
)

// RecordA represents a DNS A-Record.
type RecordA struct {
	DistinguishedName string
	Name              string
	Addresses         []string
	Timestamp         time.Time
	TimeToLive        int32
}

// convertOutput converts the unmarshaled JSON output from the recordAObject to a RecordA object.
func (r *RecordA) convertOutput(o []recordObject) {
	// Set the values of the first object to the RecordA object.
	r.DistinguishedName = o[0].DistinguishedName
	r.Name = o[0].Name
	r.Timestamp = o[0].Timestamp.Time
	r.TimeToLive = o[0].TimeToLive.Seconds

	// Set the addresses and the lowest TTL.
	if len(o) == 1 {
		r.Addresses = []string{o[0].RecordData.CimInstanceProperties["IPv4Address"]}
	} else {
		for _, record := range o {
			r.Addresses = append(r.Addresses, record.RecordData.CimInstanceProperties["IPv4Address"])

			// Set the lowest TTL to be RFC2181 compliant.
			// https://www.rfc-editor.org/rfc/rfc2181#section-5.2
			if record.TimeToLive.Seconds < r.TimeToLive {
				r.TimeToLive = record.TimeToLive.Seconds
			}
		}
	}
}

// RecordAReadParams represents parameters for the A-Record function.
type RecordAReadParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string
}

// pwshCommand returns the PowerShell command to read a local group by SID or Name.
func (params RecordAReadParams) pwshCommand() string {
	// Base command
	cmd := []string{"$r=Get-DnsServerResourceRecord -RRType 'A' -Node"}

	// Add parameters
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))

	// Ensure output is always an array.
	cmd = append(cmd, ";if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}")
	return strings.Join(cmd, " ")
}

// RecordARead gets a DNS A-Record by Name and Zone. It returns a RecordA object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordARead(ctx context.Context, params RecordAReadParams) (RecordA, error) {
	var r RecordA
	var o []recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" {
		return r, errors.New("windows.dns.server.RecordARead: record parameters 'Name' and 'Zone' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return r, winerror.Errorf(cmd, "windows.dns.server.RecordARead: %s", err)
	}

	// Convert the output to a RecordA object.
	r.convertOutput(o)

	return r, nil
}

// RecordACreateParams represents parameters for the A-Record function.
type RecordACreateParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string

	// Specifies the IPv4 addresses of the record.
	Addresses []string

	// Specifies the time to live (TTL) of the record in seconds.
	// If not provided, the default is 86400 seconds.
	// A TTL of 0 is not allowed.
	TimeToLive int32
}

// pwshCommand returns the PowerShell command to create a new A-Record.
func (params RecordACreateParams) pwshCommand() string {
	addressList := []string{}

	// Base command
	cmd := []string{"$r=Add-DnsServerResourceRecordA -AllowUpdateAny:$false -CreatePtr:$false -AgeRecord:$false -Confirm:$false -PassThru"}

	// Add parameters
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))

	// Set default TTL if not provided.
	// New-TimeSpan only allows int32 values.
	// https://learn.microsoft.com/de-de/powershell/module/microsoft.powershell.utility/new-timespan?view=powershell-7.4
	if params.TimeToLive == 0 {
		params.TimeToLive = defaultTimeToLive
	}
	cmd = append(cmd, fmt.Sprintf("-TimeToLive %s", fmt.Sprintf("$(New-TimeSpan -Seconds %d)", params.TimeToLive)))

	// Add addresses with single quotes and join them with commas.
	for _, address := range params.Addresses {
		addressList = append(addressList, fmt.Sprintf("'%s'", address))
	}
	cmd = append(cmd, fmt.Sprintf("-IPv4Address @(%s)", strings.Join(addressList, ",")))

	cmd = append(cmd, ";if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}")
	return strings.Join(cmd, " ")
}

// RecordACreate creates a DNS A-Record. It returns a RecordA object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordACreate(ctx context.Context, params RecordACreateParams) (RecordA, error) {
	var r RecordA
	var o []recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" || len(params.Addresses) == 0 {
		return r, errors.New("windows.dns.server.RecordACreate: record parameters 'Name', 'Zone' and 'Addresses' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		// Handle record already exists error.
		if strings.Contains(err.Error(), "ResourceExists") {
			return r, winerror.Errorf(cmd, "windows.dns.server.RecordACreate: the specified record already exists.")
		}

		return r, winerror.Errorf(cmd, "windows.dns.server.RecordACreate: %s", err)
	}

	// Convert the output to a RecordA object.
	r.convertOutput(o)

	return r, nil
}

// RecordAUpdateParams represents parameters for the A-Record function.
// Only the TimeToLive can be updated.
type RecordAUpdateParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string

	// Specifies the time to live (TTL) of the record in seconds.
	// If not provided, the default TTL is 86400 seconds.
	// A TTL of 0 is not allowed.
	TimeToLive int32
}

// pwshCommand returns the PowerShell command to create a new A-Record.
func (params RecordAUpdateParams) pwshCommand() string {
	// Update to default TTL if not provided.
	// New-TimeSpan only allows int32 values.
	// https://learn.microsoft.com/de-de/powershell/module/microsoft.powershell.utility/new-timespan?view=powershell-7.4
	if params.TimeToLive == 0 {
		params.TimeToLive = defaultTimeToLive
	}

	// Base command
	cmd := []string{"$nr=@();Get-DnsServerResourceRecord -RRType 'A' -Node"}

	// Add parameters and logic for handling the TTL update.
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))
	cmd = append(cmd, fmt.Sprintf("| ForEach-Object{$r=$_;$n=[ciminstance]::new($r);$n.TimeToLive=New-TimeSpan -Seconds %d", params.TimeToLive))
	cmd = append(cmd, fmt.Sprintf(";$nr+=Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName '%s' -PassThru}", params.Zone))
	cmd = append(cmd, ";if($nr.Count -ge 2){ConvertTo-Json $nr -Compress}else{ConvertTo-Json @($nr) -Compress}")

	// Return the full command.
	return strings.Join(cmd, " ")
}

// RecordAUpdate updates a DNS A-Record. It returns a RecordA object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordAUpdate(ctx context.Context, params RecordAUpdateParams) (RecordA, error) {
	var r RecordA
	var o []recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" || params.TimeToLive == 0 {
		return r, errors.New("windows.dns.server.RecordAUpdate: record parameters 'Name', 'Zone' and 'TimeToLive' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return r, winerror.Errorf(cmd, "windows.dns.server.RecordAUpdate: %s", err)
	}

	// Convert the output to a RecordA object.
	r.convertOutput(o)

	return r, nil
}

// RecordADeleteParams represents parameters for the A-Record function.
type RecordADeleteParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string
}

// pwshCommand returns the PowerShell command to create a new A-Record.
func (params RecordADeleteParams) pwshCommand() string {
	// Base command
	return fmt.Sprintf("Remove-DnsServerResourceRecord -RRType 'A' -Force -Name '%s' -ZoneName '%s'", params.Name, params.Zone)
}

// RecordADelete deletes a DNS A-Record.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordADelete(ctx context.Context, params RecordADeleteParams) error {
	var o []recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" {
		return errors.New("windows.dns.server.RecordADelete: record parameters 'Name' and 'Zone' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return winerror.Errorf(cmd, "windows.dns.server.RecordADelete: %s", err)
	}

	return nil
}
