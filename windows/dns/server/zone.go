package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/d-strobel/gowindows/winerror"
)

// Zone represents a Windows DNS server zone with its properties.
type Zone struct {
	NotifyServers                     string `json:"NotifyServers"`
	SecondaryServers                  string `json:"SecondaryServers"`
	AllowedDcForNsRecordsAutoCreation string `json:"AllowedDcForNsRecordsAutoCreation"`
	DistinguishedName                 string `json:"DistinguishedName"`
	IsAutoCreated                     bool   `json:"IsAutoCreated"`
	IsDsIntegrated                    bool   `json:"IsDsIntegrated"`
	IsPaused                          bool   `json:"IsPaused"`
	IsReadOnly                        bool   `json:"IsReadOnly"`
	IsReverseLookupZone               bool   `json:"IsReverseLookupZone"`
	IsShutdown                        bool   `json:"IsShutdown"`
	ZoneName                          string `json:"ZoneName"`
	ZoneType                          string `json:"ZoneType"`
	DirectoryPartitionName            string `json:"DirectoryPartitionName"`
	DynamicUpdate                     string `json:"DynamicUpdate"`
	IgnorePolicies                    bool   `json:"IgnorePolicies"`
	IsSigned                          bool   `json:"IsSigned"`
	IsWinsEnabled                     bool   `json:"IsWinsEnabled"`
	Notify                            string `json:"Notify"`
	ReplicationScope                  string `json:"ReplicationScope"`
	SecureSecondaries                 string `json:"SecureSecondaries"`
	ZoneFile                          string `json:"ZoneFile"`
}

// ZoneReadParams represents parameters for the ZoneRead function.
type ZoneReadParams struct {
	// Specifies the name of the zone.
	Name string
}

// pwshCommand returns the PowerShell command to read a local group by SID or Name.
func (params ZoneReadParams) pwshCommand() string {
	// Base command
	cmd := []string{"Get-DnsServerZone"}

	// Add parameters
	cmd = append(cmd, fmt.Sprintf("-Name '%s'", params.Name))

	cmd = append(cmd, "| ConvertTo-Json -Compress")
	return strings.Join(cmd, " ")
}

// ZoneRead gets a DNS server zone by Name and returns a Zone object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) ZoneRead(ctx context.Context, params ZoneReadParams) (Zone, error) {
	var z Zone

	// Assert needed parameters
	if params.Name == "" {
		return z, errors.New("windows.dns.server.ZoneRead: zone parameter 'Name' must be set")
	}

	// We want to retrieve exactly one zone.
	if strings.Contains(params.Name, "*") {
		return z, errors.New("windows.dns.server.ZoneRead: zone parameter 'Name' does not allow wildcards")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &z); err != nil {
		return z, winerror.Errorf(cmd, "windows.dns.server.ZoneRead: %s", err)
	}
	return z, nil
}

// ZoneList gets all DNS server zones returns a list of Zone objects.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) ZoneList(ctx context.Context) ([]Zone, error) {
	var z []Zone

	// Run command
	cmd := "Get-DnsServerZone | ConvertTo-Json -Compress"
	if err := run(ctx, c, cmd, &z); err != nil {
		return z, winerror.Errorf(cmd, "windows.dns.server.ZoneList: %s", err)
	}
	return z, nil
}
