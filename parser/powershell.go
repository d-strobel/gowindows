package parser

import (
	"strings"
)

// Optional parameters
type PwshOpts struct {
	JSONOutput bool
}

// NewPwshCommand returns a ready to run powershell command
func NewPwshCommand(cmd []string, opts *PwshOpts) (string, error) {

	// Convert output to json
	if opts.JSONOutput {
		cmd = append(cmd, "| ConvertTo-Json")
	}

	pwshCmd := strings.Join(cmd, " ")

	return pwshCmd, nil
}
