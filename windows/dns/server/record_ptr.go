package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/d-strobel/gowindows/winerror"
)

// RecordPTR represents a DNS PTR-Record.
type RecordPTR struct {
	DistinguishedName string
	Name              string
	PTR               string
	Timestamp         time.Time
	TimeToLive        int32
}

// convertOutput converts the unmarshaled JSON output from the recordObject to a RecordPTR object.
func (r *RecordPTR) convertOutput(o recordObject) {
	r.DistinguishedName = o.DistinguishedName
	r.Name = o.Name
	r.Timestamp = o.Timestamp.Time
	r.TimeToLive = o.TimeToLive.Seconds
	r.PTR = o.RecordData.CimInstanceProperties["PtrDomainName"]
}

// RecordPTRReadParams represents parameters for the PTR-Record read function.
type RecordPTRReadParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string
}

// pwshCommand returns the PowerShell command to read a PTR-Record.
func (params RecordPTRReadParams) pwshCommand() string {
	// Base command
	cmd := []string{"Get-DnsServerResourceRecord -RRType 'PTR' -Node"}

	// Add parameters
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))

	// Ensure Json Output
	cmd = append(cmd, "| ConvertTo-Json -Compress")
	return strings.Join(cmd, " ")
}

// RecordPTRRead gets a PTR-Record. It returns a RecordPTR object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordPTRRead(ctx context.Context, params RecordPTRReadParams) (RecordPTR, error) {
	var r RecordPTR
	var o recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" {
		return r, errors.New("windows.dns.server.RecordPTRRead: record parameters 'Name' and 'Zone' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return r, winerror.Errorf(cmd, "windows.dns.server.RecordPTRRead: %s", err)
	}

	// Convert the output to a RecordPTR object.
	r.convertOutput(o)

	return r, nil
}

// RecordPTRCreateParams represents parameters for the PTR-Record create function.
type RecordPTRCreateParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string

	// Specifies the canonical name this record will point to.
	PTR string

	// Specifies the time to live (TTL) of the record in seconds.
	// If not provided, the default is 86400 seconds.
	// A TTL of 0 is not allowed.
	TimeToLive int32
}

// pwshCommand returns the PowerShell command to create a new PTR-Record.
func (params RecordPTRCreateParams) pwshCommand() string {
	// Base command
	cmd := []string{"Add-DnsServerResourceRecordPTR -AllowUpdateAny:$false -AgeRecord:$false -Confirm:$false -PassThru"}

	// Add parameters
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))
	cmd = append(cmd, fmt.Sprintf("-PtrDomainName '%s'", params.PTR))

	// Set default TTL if not provided.
	// New-TimeSpan only allows int32 values.
	// https://learn.microsoft.com/de-de/powershell/module/microsoft.powershell.utility/new-timespan?view=powershell-7.4
	if params.TimeToLive == 0 {
		params.TimeToLive = defaultTimeToLive
	}
	cmd = append(cmd, fmt.Sprintf("-TimeToLive %s", fmt.Sprintf("$(New-TimeSpan -Seconds %d)", params.TimeToLive)))

	// Join the command and ensure Json Output
	cmd = append(cmd, "| ConvertTo-Json -Compress")

	return strings.Join(cmd, " ")
}

// RecordPTRCreate creates a PTR-Record. It returns a RecordPTR object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordPTRCreate(ctx context.Context, params RecordPTRCreateParams) (RecordPTR, error) {
	var r RecordPTR
	var o recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" || params.PTR == "" {
		return r, errors.New("windows.dns.server.RecordPTRCreate: record parameters 'Name', 'Zone' and 'PTR' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		// Handle record already exists error.
		if strings.Contains(err.Error(), "ResourceExists") {
			return r, winerror.Errorf(cmd, "windows.dns.server.RecordPTRCreate: the specified record already exists.")
		}

		return r, winerror.Errorf(cmd, "windows.dns.server.RecordPTRCreate: %s", err)
	}

	// Convert the output to a RecordPTR object.
	r.convertOutput(o)

	return r, nil
}

// RecordPTRUpdateParams represents parameters for the PTR-Record update function.
// The PTR and TimeToLive can be updated.
type RecordPTRUpdateParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string

	// Specifies the canonical name this record will point to.
	PTR string

	// Specifies the time to live (TTL) of the record in seconds.
	// If not provided, the default TTL is 86400 seconds.
	// A TTL of 0 is not allowed.
	TimeToLive int32
}

// pwshCommand returns the PowerShell command to update a PTR-Record.
func (params RecordPTRUpdateParams) pwshCommand() string {
	// Update to default TTL if not provided.
	// New-TimeSpan only allows int32 values.
	// https://learn.microsoft.com/de-de/powershell/module/microsoft.powershell.utility/new-timespan?view=powershell-7.4
	if params.TimeToLive == 0 {
		params.TimeToLive = defaultTimeToLive
	}

	// Get command
	cmd := []string{"$r=Get-DnsServerResourceRecord -RRType 'PTR' -Node"}
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))

	// Add logic for handling TTL and PTR update.
	cmd = append(cmd, ";$n=[ciminstance]::new($r)")
	cmd = append(cmd, fmt.Sprintf(";$n.TimeToLive=New-TimeSpan -Seconds %d", params.TimeToLive))
	cmd = append(cmd, fmt.Sprintf(";$n.RecordData.PtrDomainName='%s'", params.PTR))
	cmd = append(cmd, fmt.Sprintf(";Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName '%s' -PassThru", params.Zone))

	// Ensure Json Output
	cmd = append(cmd, "| ConvertTo-Json -Compress")

	// Return the full command.
	return strings.Join(cmd, " ")
}

// RecordPTRUpdate updates a PTR-Record. It returns a RecordPTR object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordPTRUpdate(ctx context.Context, params RecordPTRUpdateParams) (RecordPTR, error) {
	var r RecordPTR
	var o recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" || params.PTR == "" {
		return r, errors.New("windows.dns.server.RecordPTRUpdate: record parameters 'Name', 'Zone' and 'PTR' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return r, winerror.Errorf(cmd, "windows.dns.server.RecordPTRUpdate: %s", err)
	}

	// Convert the output to a RecordPTR object.
	r.convertOutput(o)

	return r, nil
}

// RecordPTRDeleteParams represents parameters for the PTR-Record delete function.
type RecordPTRDeleteParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string
}

// pwshCommand returns the PowerShell command to delete a PTR-Record.
func (params RecordPTRDeleteParams) pwshCommand() string {
	// Base command
	return fmt.Sprintf("Remove-DnsServerResourceRecord -RRType 'PTR' -Force -Name '%s' -ZoneName '%s'", params.Name, params.Zone)
}

// RecordPTRDelete deletes a PTR-Record.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordPTRDelete(ctx context.Context, params RecordPTRDeleteParams) error {
	var o recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" {
		return errors.New("windows.dns.server.RecordPTRDelete: record parameters 'Name' and 'Zone' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return winerror.Errorf(cmd, "windows.dns.server.RecordPTRDelete: %s", err)
	}

	return nil
}
