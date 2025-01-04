package dhcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/d-strobel/gowindows/parsing"
	"github.com/d-strobel/gowindows/winerror"
)

// FailoverV4 represents an IPv4 failover object.
type FailoverV4 struct {
	Name                string                  `json:"Name"`
	ScopeId             scopeIdVal              `json:"ScopeId"`
	PrimaryServerIp     addressString           `json:"PrimaryServerIP"`
	PrimaryServerName   string                  `json:"PrimaryServerName"`
	SecondaryServerIp   addressString           `json:"SecondaryServerIP"`
	SecondaryServerName string                  `json:"SecondaryServerName"`
	AutoStateTransition bool                    `json:"AutoStateTransition"`
	EnableAuth          bool                    `json:"EnableAuth"`
	LoadBalancePercent  uint32                  `json:"LoadBalancePercent"`
	MaxClientLeadTime   parsing.CimTimeDuration `json:"MaxClientLeadTime"`
	Mode                string                  `json:"Mode"`
	ReservePercent      uint32                  `json:"ReservePercent"`
	ServerRole          string                  `json:"ServerRole"`
	ServerType          string                  `json:"ServerType"`
	State               string                  `json:"State"`
	StateSwitchInterval parsing.CimTimeDuration `json:"StateSwitchInterval"`
}

// FailoverV4ReadParams represents the parameters for the FailoverV4Read function.
type FailoverV4ReadParams struct {
	// Specifies the name of a failover relationship for which the properties are returned.
	Name string
}

// pwshCommand returns the PowerShell command to create an IPv4 Failover.
func (params FailoverV4ReadParams) pwshCommand() string {
	return fmt.Sprintf(
		"Get-DhcpServerv4Failover -Name '%s' | ConvertTo-Json -Compress",
		params.Name,
	)
}

// FailoverV4Read returns a FailoverV4 object.
func (c *Client) FailoverV4Read(ctx context.Context, params FailoverV4ReadParams) (FailoverV4, error) {
	var f FailoverV4

	// Assert needed parameters.
	if params.Name == "" {
		return f, errors.New("windows.dhcp.FailoverV4Read: failover parameter 'Name' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &f); err != nil {
		return f, winerror.Errorf(cmd, "windows.dhcp.FailoverV4Read: %s", err)
	}

	return f, nil
}
