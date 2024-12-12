package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/d-strobel/gowindows/winerror"
)

// RecordCName represents a DNS CName-Record.
type RecordCName struct {
	DistinguishedName string
	Name              string
	CName             string
	Timestamp         time.Time
	TimeToLive        time.Duration
}

// convertOutput converts the unmarshaled JSON output from the recordObject to a RecordCName object.
func (r *RecordCName) convertOutput(o recordObject) {
	r.DistinguishedName = o.DistinguishedName
	r.Name = o.Name
	r.Timestamp = o.Timestamp.Time
	r.TimeToLive = o.TimeToLive.Duration
	r.CName = o.RecordData.CimInstanceProperties["HostNameAlias"]
}

// RecordCNameReadParams represents parameters for the CName-Record read function.
type RecordCNameReadParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string
}

// pwshCommand returns the PowerShell command to read a CName-Record.
func (params RecordCNameReadParams) pwshCommand() string {
	// Base command
	cmd := []string{"Get-DnsServerResourceRecord -RRType 'CName' -Node"}

	// Add parameters
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))

	// Ensure Json Output
	cmd = append(cmd, "| ConvertTo-Json -Compress")
	return strings.Join(cmd, " ")
}

// RecordCNameRead gets a CName-Record. It returns a RecordCName object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordCNameRead(ctx context.Context, params RecordCNameReadParams) (RecordCName, error) {
	var r RecordCName
	var o recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" {
		return r, errors.New("windows.dns.server.RecordCNameRead: record parameters 'Name' and 'Zone' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return r, winerror.Errorf(cmd, "windows.dns.server.RecordCNameRead: %s", err)
	}

	// Convert the output to a RecordCName object.
	r.convertOutput(o)

	return r, nil
}

// RecordCNameCreateParams represents parameters for the CName-Record create function.
type RecordCNameCreateParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string

	// Specifies the CName of the record.
	CName string

	// Specifies the time to live (TTL) of the record in seconds.
	// If not provided, the default is 86400 seconds.
	// A TTL of 0 is not allowed.
	TimeToLive time.Duration
}

// pwshCommand returns the PowerShell command to create a new CName-Record.
func (params RecordCNameCreateParams) pwshCommand() string {
	// Base command
	cmd := []string{"Add-DnsServerResourceRecordCName -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru"}

	// Add parameters
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))
	cmd = append(cmd, fmt.Sprintf("-HostNameAlias '%s'", params.CName))

	// Set default TTL if not provided.
	// New-TimeSpan only allows int32 values.
	// https://learn.microsoft.com/de-de/powershell/module/microsoft.powershell.utility/new-timespan?view=powershell-7.4
	if params.TimeToLive == 0 {
		params.TimeToLive = defaultTimeToLive
	}
	seconds := int32(params.TimeToLive.Round(time.Second).Seconds())
	cmd = append(cmd, fmt.Sprintf("-TimeToLive %s", fmt.Sprintf("$(New-TimeSpan -Seconds %d)", seconds)))

	// Join the command and ensure Json Output
	cmd = append(cmd, "| ConvertTo-Json -Compress")

	return strings.Join(cmd, " ")
}

// RecordCNameCreate creates a CName-Record. It returns a RecordCName object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordCNameCreate(ctx context.Context, params RecordCNameCreateParams) (RecordCName, error) {
	var r RecordCName
	var o recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" || params.CName == "" {
		return r, errors.New("windows.dns.server.RecordCNameCreate: record parameters 'Name', 'Zone' and 'CName' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		// Handle record already exists error.
		if strings.Contains(err.Error(), "ResourceExists") {
			return r, winerror.Errorf(cmd, "windows.dns.server.RecordCNameCreate: the specified record already exists.")
		}

		return r, winerror.Errorf(cmd, "windows.dns.server.RecordCNameCreate: %s", err)
	}

	// Convert the output to a RecordCName object.
	r.convertOutput(o)

	return r, nil
}

// RecordCNameUpdateParams represents parameters for the CName-Record update function.
// The CName and TimeToLive can be updated.
type RecordCNameUpdateParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string

	// Specifies the CName of the record.
	CName string

	// Specifies the time to live (TTL) of the record in seconds.
	// If not provided, the default TTL is 86400 seconds.
	// A TTL of 0 is not allowed.
	TimeToLive time.Duration
}

// pwshCommand returns the PowerShell command to update a CName-Record.
func (params RecordCNameUpdateParams) pwshCommand() string {
	// Update to default TTL if not provided.
	// New-TimeSpan only allows int32 values.
	// https://learn.microsoft.com/de-de/powershell/module/microsoft.powershell.utility/new-timespan?view=powershell-7.4
	if params.TimeToLive == 0 {
		params.TimeToLive = defaultTimeToLive
	}
	seconds := int32(params.TimeToLive.Round(time.Second).Seconds())

	// Get command
	cmd := []string{"$r=Get-DnsServerResourceRecord -RRType 'CName' -Node"}
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))

	// Add logic for handling TTL and CName update.
	cmd = append(cmd, ";$n=[ciminstance]::new($r)")
	cmd = append(cmd, fmt.Sprintf(";$n.TimeToLive=New-TimeSpan -Seconds %d", seconds))
	cmd = append(cmd, fmt.Sprintf(";$n.RecordData.HostNameAlias='%s'", params.CName))
	cmd = append(cmd, fmt.Sprintf(";Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName '%s' -PassThru", params.Zone))

	// Ensure Json Output
	cmd = append(cmd, "| ConvertTo-Json -Compress")

	// Return the full command.
	return strings.Join(cmd, " ")
}

// RecordCNameUpdate updates a CName-Record. It returns a RecordCName object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordCNameUpdate(ctx context.Context, params RecordCNameUpdateParams) (RecordCName, error) {
	var r RecordCName
	var o recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" || params.CName == "" {
		return r, errors.New("windows.dns.server.RecordCNameUpdate: record parameters 'Name', 'Zone' and 'CName' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return r, winerror.Errorf(cmd, "windows.dns.server.RecordCNameUpdate: %s", err)
	}

	// Convert the output to a RecordCName object.
	r.convertOutput(o)

	return r, nil
}

// RecordCNameDeleteParams represents parameters for the CName-Record delete function.
type RecordCNameDeleteParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string
}

// pwshCommand returns the PowerShell command to delete a CName-Record.
func (params RecordCNameDeleteParams) pwshCommand() string {
	// Base command
	return fmt.Sprintf("Remove-DnsServerResourceRecord -RRType 'CName' -Force -Name '%s' -ZoneName '%s'", params.Name, params.Zone)
}

// RecordCNameDelete deletes a CName-Record.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordCNameDelete(ctx context.Context, params RecordCNameDeleteParams) error {
	var o recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" {
		return errors.New("windows.dns.server.RecordCNameDelete: record parameters 'Name' and 'Zone' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return winerror.Errorf(cmd, "windows.dns.server.RecordCNameDelete: %s", err)
	}

	return nil
}
