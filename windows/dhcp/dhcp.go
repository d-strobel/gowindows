// Package dhcp provides a Go library for handling Windows DHCP Server resources.
// The functions are related to the Powershell DhcpServer cmdlets provided by Windows.
// https://learn.microsoft.com/en-us/powershell/module/dhcpserver/?view=windowsserver2022-ps
package dhcp

import (
	"context"
	"encoding/json"
	"net/netip"

	"errors"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parsing"
)

// dhcp is a type constraint for the run function, ensuring it works with specific types.
type dhcp interface {
	scopeObject
}

// scopeObject is used to unmarshal the JSON output of a scope object.
type scopeObject struct {
	Name             string                  `json:"Name"`
	Description      string                  `json:"Description"`
	ScopeId          scopeId                 `json:"ScopeId"`
	StartRange       startRange              `json:"StartRange"`
	EndRange         endRange                `json:"EndRange"`
	SubnetMask       subnetMask              `json:"SubnetMask"`
	State            string                  `json:"State"`
	MaxBootpClients  uint32                  `json:"MaxBootpClients"`
	ActivatePolicies bool                    `json:"ActivatePolicies"`
	NapEnable        bool                    `json:"NapEnable"`
	NapProfile       string                  `json:"NapProfile"`
	Delay            uint16                  `json:"Delay"`
	LeaseDuration    parsing.CimTimeDuration `json:"LeaseDuration"`
}
type scopeId struct {
	Address netip.Addr `json:"IPAddressToString"`
}
type startRange struct {
	Address netip.Addr `json:"IPAddressToString"`
}
type endRange struct {
	Address netip.Addr `json:"IPAddressToString"`
}
type subnetMask struct {
	Address netip.Addr `json:"IPAddressToString"`
}

// Client represents a client for handling DHCP server functions.
type Client struct {
	// Connection represents a connection.Connection object.
	Connection connection.Connection

	// decodeCliXmlErr represents a function that decodes a CLIXML error and returns aa  human readable string.
	decodeCliXmlErr func(string) (string, error)
}

// NewClient returns a new instance of the Client.
func NewClient(conn connection.Connection) *Client {
	return NewClientWithParser(conn, parsing.DecodeCliXmlErr)
}

// NewClientWithParser returns a new instance of the Client.
// It requires a connection and parsing as input parameters.
func NewClientWithParser(conn connection.Connection, parsing func(string) (string, error)) *Client {
	return &Client{Connection: conn, decodeCliXmlErr: parsing}
}

// run runs a PowerShell command against a Windows system, handles the command results,
// and unmarshals the output into a local object type.
func run[T dhcp](ctx context.Context, c *Client, cmd string, t *T) error {
	// Run the command
	result, err := c.Connection.RunWithPowershell(ctx, cmd)
	if err != nil {
		return err
	}

	// Handle stderr
	if result.StdErr != "" {
		stderr, err := c.decodeCliXmlErr(result.StdErr)
		if err != nil {
			return err
		}

		return errors.New(stderr)
	}

	if result.StdOut == "" {
		return nil
	}

	// Unmarshal stdout
	if err = json.Unmarshal([]byte(result.StdOut), &t); err != nil {
		return err
	}

	return nil
}
