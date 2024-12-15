package dhcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/d-strobel/gowindows/winerror"
)

// ScopeV4 represents an IPv4 DHCP scope.
type ScopeV4 struct {
	Name             string
	Description      string
	StartRange       netip.Addr
	EndRange         netip.Addr
	SubnetMask       netip.Addr
	Enabled          bool
	MaxBootpClients  uint32
	ActivatePolicies bool
	NapEnable        bool
	NapProfile       string
	Delay            uint16
	LeaseDuration    time.Duration
}

// convertOutput converts the unmarshaled JSON output from the scopeObject to a ScopeV4 object.
func (s *ScopeV4) convertOutput(o scopeObject) {
	s.Name = o.Name
	s.Description = o.Description
	s.StartRange = o.StartRange.Address
	s.EndRange = o.EndRange.Address
	s.SubnetMask = o.SubnetMask.Address
	s.Enabled = o.State == "Active"
	s.MaxBootpClients = o.MaxBootpClients
	s.ActivatePolicies = o.ActivatePolicies
	s.NapEnable = o.NapEnable
	s.NapProfile = o.NapProfile
	s.Delay = o.Delay
	s.LeaseDuration = o.LeaseDuration.Duration
}

// ScopeV4ReadParams represents parameters for the scope read function.
type ScopeV4ReadParams struct {
	// Specify the ID of the scope.
	// This is the scopes network address, e.g. 192.168.10.0.
	ScopeId netip.Addr
}

// pwshCommand returns the PowerShell command to read an A-Record.
func (params ScopeV4ReadParams) pwshCommand() string {
	// Base command
	return fmt.Sprintf("Get-DhcpServerv4Scope -ScopeId '%s' | ConvertTo-Json -Compress", params.ScopeId)
}

// ScopeV4Read gets a DHCP scope. It returns a ScopeV4 object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) ScopeV4Read(ctx context.Context, params ScopeV4ReadParams) (ScopeV4, error) {
	var s ScopeV4
	var o scopeObject

	// Assert needed parameters
	if !params.ScopeId.Is4() {
		return s, errors.New("windows.dhcp.ScopeV4Read: scope parameter 'ScopeId' must be an IPv4 network address")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return s, winerror.Errorf(cmd, "windows.dhcp.ScopeV4Read: %s", err)
	}

	// Convert the output to a ScopeV4 object.
	s.convertOutput(o)

	return s, nil
}
