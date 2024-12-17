package dhcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/d-strobel/gowindows/parsing"
	"github.com/d-strobel/gowindows/winerror"
)

// ScopeV4 represents an IPv4 DHCP scope.
type ScopeV4 struct {
	Name             string
	Description      string
	ScopeId          netip.Addr
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
	s.ScopeId = o.ScopeId.Address
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

// pwshCommand returns the PowerShell command to read a DHCP scope.
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

// ScopeV4CreateParams represents parameters for the scope create function.
type ScopeV4CreateParams struct {
	// Specifies the enabled state of the policy enforcement on the scope that is added.
	ActivatePolicies bool

	// Specifies the number of milliseconds by which the DHCP server service should wait before
	// responding to the client requests. Specify this parameter if the scope is part of a split
	// scope deployment and this DHCP server service should act as a secondary DHCP server service
	// for the scope being added.
	Delay uint16

	// Specifies the description string for the IPv4 scope that is added.
	Description string

	// Specify if the scope is enabled.
	Enabled bool

	// Specifies the ending IP address of the range in the subnet from which IP addresses should
	// be leased by the DHCP server service.
	EndRange netip.Addr

	// Specifies the time interval for which an IP address should be leased to a client in this scope.
	LeaseDuration time.Duration

	// Specifies, if the scope type is specified as Both to allow for both DHCP and BootP clients,
	// the maximum number of BootP clients which should be leased an IP address from this scope.
	MaxBootpClients uint32

	// Specifies the name of the IPv4 scope that is added.
	Name string

	// Specifies the enabled state of Network Access Protection (NAP) for this scope.
	// If NAP is enabled, then the DHCP server service passes the statement of health (SoH) received
	// from the client to the network policy server (NPS). Based on the NAP profile set,
	// the NPS determines the network access to grant to the client.
	NapEnable bool

	// Specifies that the NAP profile should be set only if NAP is enabled on the scope.
	// The NAP profile refers to the MS Service Class which is a condition used in network policies on NPS.
	NapProfile string

	// Specifies the starting IP address of the range in the subnet from which IP addresses should be leased
	// by the DHCP server service.
	StartRange netip.Addr

	// Specifies the subnet mask for the scope specified in IP address format. For example: 255.255.255.0.
	SubnetMask netip.Addr

	// Specifies the name of the superscope to which the scope is added.
	Superscope string

	// Specifies the type of clients to be serviced by the scope.
	// The type of the scope determines whether the DHCP server service responds to only DHCP client requests,
	// only BootP client requests, or Both types of clients.
	//
	// The acceptable values for this parameter are:
	// "Dhcp", "Bootp", "Both".
	Type string
}

// pwshCommand returns the PowerShell command to create a DHCP scope.
func (params ScopeV4CreateParams) pwshCommand() string {
	// Base command
	cmd := []string{
		fmt.Sprintf("Add-DhcpServerv4Scope -PassThru -Confirm:$false -Name '%s' -StartRange '%s' -EndRange '%s' -SubnetMask '%s'",
			params.Name,
			params.StartRange,
			params.EndRange,
			params.SubnetMask,
		),
	}

	// Add optional parameters
	if params.Description != "" {
		cmd = append(cmd, fmt.Sprintf("-Description '%s'", params.Description))
	}

	if params.Enabled {
		cmd = append(cmd, "-State 'Active'")
	} else {
		cmd = append(cmd, "-State 'InActive'")
	}

	if params.MaxBootpClients != 0 {
		cmd = append(cmd, fmt.Sprintf("-MaxBootpClients %d", params.MaxBootpClients))
	}

	if params.ActivatePolicies {
		cmd = append(cmd, "-ActivatePolicies")
	}

	if params.NapEnable {
		cmd = append(cmd, "-NapEnable")
	}

	if params.NapProfile != "" {
		cmd = append(cmd, fmt.Sprintf("-NapProfile '%s'", params.NapProfile))
	}

	if params.Delay != 0 {
		cmd = append(cmd, fmt.Sprintf("-Delay %d", params.Delay))
	}

	if params.LeaseDuration != 0 {
		cmd = append(cmd, fmt.Sprintf("-LeaseDuration %s", parsing.PwshTimespanString(params.LeaseDuration)))
	}

	if params.Type != "" {
		cmd = append(cmd, fmt.Sprintf("-Type '%s'", params.Type))
	}

	if params.Superscope != "" {
		cmd = append(cmd, fmt.Sprintf("-SuperscopeName '%s'", params.Superscope))
	}

	// Convert output to json
	cmd = append(cmd, "| ConvertTo-Json -Compress")

	// Return the full command
	return strings.Join(cmd, " ")
}

// ScopeV4Create creates a new DHCP IPv4 scope. It returns a ScopeV4 object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) ScopeV4Create(ctx context.Context, params ScopeV4CreateParams) (ScopeV4, error) {
	var s ScopeV4
	var o scopeObject

	// Assert needed parameters
	if params.Name == "" {
		return s, errors.New("windows.dhcp.ScopeV4Create: scope parameter 'Name' must be set")
	}

	if !params.StartRange.Is4() || !params.EndRange.Is4() || !params.SubnetMask.Is4() {
		return s, errors.New("windows.dhcp.ScopeV4Create: scope parameter 'StartRange', 'EndRange' and 'SubnetMask' must be a valid IPv4 address")
	}

	if params.Type != "" {
		if params.Type != "Dhcp" && params.Type != "Bootp" && params.Type != "Both" {
			return s, errors.New("windows.dhcp.ScopeV4Create: scope parameter 'Type' must be one of the following values: 'Dhcp', 'Bootp', 'Both'")
		}
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &o); err != nil {
		return s, winerror.Errorf(cmd, "windows.dhcp.ScopeV4Create: %s", err)
	}

	// Convert the output to a ScopeV4 object.
	s.convertOutput(o)

	return s, nil
}
