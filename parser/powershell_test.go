package parser

import (
	"testing"
)

func TestNewPwshCommandWithJSONOutputTrue(t *testing.T) {
	opts := &PwshOpts{
		JSONOutput: true,
	}
	cmd := []string{"your_command"}

	result := NewPwshCommand(cmd, opts)
	expected := "your_command | ConvertTo-Json"

	if result != expected {
		t.Errorf("Expected result to be %s, but got %s", expected, result)
	}
}

func TestNewPwshCommandWithJSONOutputFalse(t *testing.T) {
	opts := &PwshOpts{
		JSONOutput: false,
	}
	cmd := []string{"your_command"}

	result := NewPwshCommand(cmd, opts)
	expected := "your_command"

	if result != expected {
		t.Errorf("Expected result to be %s, but got %s", expected, result)
	}
}
