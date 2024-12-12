package dns

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/d-strobel/gowindows/winerror"
)

// RecordA represents a DNS A-Record.
type RecordA struct {
	DistinguishedName string
	Name              string
	Addresses         []netip.Addr
	Timestamp         time.Time
	TimeToLive        time.Duration
}

// convertOutput converts the unmarshaled JSON output from the recordObject to a RecordA object.
func (r *RecordA) convertOutput(o []recordObject) error {
	// Set the values of the first object to the RecordA object.
	r.DistinguishedName = o[0].DistinguishedName
	r.Name = o[0].Name
	r.Timestamp = o[0].Timestamp.Time
	r.TimeToLive = o[0].TimeToLive.Duration

	// Set the addresses and the lowest TTL.
	if len(o) == 1 {
		ip, err := netip.ParseAddr(o[0].RecordData.CimInstanceProperties["IPv4Address"])
		if err != nil {
			return err
		}
		r.Addresses = []netip.Addr{ip}
	} else {
		for _, record := range o {
			ip, err := netip.ParseAddr(record.RecordData.CimInstanceProperties["IPv4Address"])
			if err != nil {
				return err
			}
			r.Addresses = append(r.Addresses, ip)

			// Set the lowest TTL to be RFC2181 compliant.
			// https://www.rfc-editor.org/rfc/rfc2181#section-5.2
			if record.TimeToLive.Duration < r.TimeToLive {
				r.TimeToLive = record.TimeToLive.Duration
			}
		}
	}

	return nil
}

// RecordAReadParams represents parameters for the A-Record read function.
type RecordAReadParams struct {
	// Specifies the name of the record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string
}

// pwshCommand returns the PowerShell command to read an A-Record.
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

// RecordARead gets an A-Record by Name and Zone. It returns a RecordA object.
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
	if err := r.convertOutput(o); err != nil {
		return r, fmt.Errorf(cmd, "windows.dns.server.RecordARead: failed to convert output to RecordA object: %s", err)
	}

	return r, nil
}

// RecordACreateParams represents parameters for the A-Record create function.
type RecordACreateParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string

	// Specifies the IPv4 addresses of the record.
	Addresses []netip.Addr

	// Specifies the time to live (TTL) of the record in seconds.
	// If not provided, the default is 86400 seconds.
	// A TTL of 0 is not allowed.
	TimeToLive time.Duration
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
	if params.TimeToLive == 0 {
		params.TimeToLive = defaultTimeToLive
	}

	// New-TimeSpan only allows int32 values. So we round the duration to seconds.
	// https://learn.microsoft.com/de-de/powershell/module/microsoft.powershell.utility/new-timespan?view=powershell-7.4
	seconds := int32(params.TimeToLive.Round(time.Second).Seconds())
	cmd = append(cmd, fmt.Sprintf("-TimeToLive %s", fmt.Sprintf("$(New-TimeSpan -Seconds %d)", seconds)))

	// Add addresses with single quotes and join them with commas.
	for _, address := range params.Addresses {
		addressList = append(addressList, fmt.Sprintf("'%s'", address.String()))
	}
	cmd = append(cmd, fmt.Sprintf("-IPv4Address @(%s)", strings.Join(addressList, ",")))

	cmd = append(cmd, ";if($r.Count -ge 2){ConvertTo-Json $r -Compress}else{ConvertTo-Json @($r) -Compress}")
	return strings.Join(cmd, " ")
}

// RecordACreate creates a new A-Record. It returns a RecordA object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) RecordACreate(ctx context.Context, params RecordACreateParams) (RecordA, error) {
	var r RecordA
	var o []recordObject

	// Assert needed parameters
	if params.Name == "" || params.Zone == "" || len(params.Addresses) == 0 {
		return r, errors.New("windows.dns.server.RecordACreate: record parameters 'Name', 'Zone' and 'Addresses' must be set")
	}

	// Assert Ipv4 addresses
	for _, address := range params.Addresses {
		if !address.Is4() {
			return r, errors.New("windows.dns.server.RecordACreate: record parameter 'Addresses' must be a list of IPv4 addresses")
		}
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
	if err := r.convertOutput(o); err != nil {
		return r, fmt.Errorf(cmd, "windows.dns.server.RecordACreate: failed to convert output to RecordA object: %s", err)
	}

	return r, nil
}

// RecordAUpdateParams represents parameters for the A-Record update function.
// Only the TimeToLive can be updated.
type RecordAUpdateParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string

	// Specifies the time to live (TTL) of the record in seconds.
	// If not provided, the default TTL is 86400 seconds.
	// A TTL of 0 is not allowed.
	TimeToLive time.Duration
}

// pwshCommand returns the PowerShell command to update an A-Record.
func (params RecordAUpdateParams) pwshCommand() string {
	// Update to default TTL if not provided.
	// New-TimeSpan only allows int32 values.
	// https://learn.microsoft.com/de-de/powershell/module/microsoft.powershell.utility/new-timespan?view=powershell-7.4
	if params.TimeToLive == 0 {
		params.TimeToLive = defaultTimeToLive
	}
	seconds := int32(params.TimeToLive.Round(time.Second).Seconds())

	// Base command
	cmd := []string{"$nr=@();Get-DnsServerResourceRecord -RRType 'A' -Node"}

	// Add parameters and logic for handling the TTL update.
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))
	cmd = append(cmd, fmt.Sprintf("-ZoneName '%s'", params.Zone))
	cmd = append(cmd, fmt.Sprintf("| ForEach-Object{$r=$_;$n=[ciminstance]::new($r);$n.TimeToLive=New-TimeSpan -Seconds %d", seconds))
	cmd = append(cmd, fmt.Sprintf(";$nr+=Set-DnsServerResourceRecord -OldInputObject $r -NewInputObject $n -ZoneName '%s' -PassThru}", params.Zone))
	cmd = append(cmd, ";if($nr.Count -ge 2){ConvertTo-Json $nr -Compress}else{ConvertTo-Json @($nr) -Compress}")

	// Return the full command.
	return strings.Join(cmd, " ")
}

// RecordAUpdate updates an A-Record. It returns a RecordA object.
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
	if err := r.convertOutput(o); err != nil {
		return r, fmt.Errorf(cmd, "windows.dns.server.RecordAUpdate: failed to convert output to RecordA object: %s", err)
	}

	return r, nil
}

// RecordADeleteParams represents parameters for the A-Record delete function.
type RecordADeleteParams struct {
	// Specifies the name of the Record.
	Name string

	// Specifies the zone in which the record is located.
	Zone string
}

// pwshCommand returns the PowerShell command to delete an A-Record.
func (params RecordADeleteParams) pwshCommand() string {
	// Base command
	return fmt.Sprintf("Remove-DnsServerResourceRecord -RRType 'A' -Force -Name '%s' -ZoneName '%s'", params.Name, params.Zone)
}

// RecordADelete deletes an A-Record.
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
