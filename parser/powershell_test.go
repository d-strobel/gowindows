package parser

import (
	"testing"
)

func TestNewPwshCommandWithJSONOutputTrue(t *testing.T) {
	opts := &PwshOpts{
		JSONOutput: true,
	}
	cmd := []string{"your_command"}

	result, err := NewPwshCommand(cmd, opts)

	expected := "your_command | ConvertTo-Json"

	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected result to be %s, but got %s", expected, result)
	}
}

func TestNewPwshCommandWithJSONOutputFalse(t *testing.T) {
	opts := &PwshOpts{
		JSONOutput: false,
	}
	cmd := []string{"your_command"}

	result, err := NewPwshCommand(cmd, opts)

	expected := "your_command"

	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected result to be %s, but got %s", expected, result)
	}
}
