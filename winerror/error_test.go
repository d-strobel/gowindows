package winerror

import (
	"testing"
)

func TestErrorf(t *testing.T) {

	// Test simple error message
	t.Run("SimpleErrorMessage", func(t *testing.T) {
		root := "RootError"
		msg := "Sample error message"
		expectedResult := "[RootError] Sample error message"

		winErr := Errorf(root, msg)

		if winErr.Error() != expectedResult {
			t.Errorf("Expected error message: '%s'\nGot: '%s'", expectedResult, winErr.Error())
		}
	})

	// Test format error message
	t.Run("FormatErrorMessage", func(t *testing.T) {
		root := "RootError"
		msg := "test"
		expectedResult := "[RootError] Sample error message test"

		winErr := Errorf(root, "Sample error message %s", msg)

		if winErr.Error() != expectedResult {
			t.Errorf("Expected error message: '%s'\nGot: '%s'", expectedResult, winErr.Error())
		}
	})

	// Test root error message
	t.Run("FormatErrorMessage", func(t *testing.T) {
		msg := "Sample error message"

		// ConfigError
		root := ConfigError
		expectedResult := "[configuration_error] Sample error message"
		winErr := Errorf(root, msg)

		if winErr.Error() != expectedResult {
			t.Errorf("Expected error message: '%s'\nGot: '%s'", expectedResult, winErr.Error())
		}

		// ConnectionError
		root = ConnectionError
		expectedResult = "[connection_error] Sample error message"
		winErr = Errorf(root, msg)

		if winErr.Error() != expectedResult {
			t.Errorf("Expected error message: '%s'\nGot: '%s'", expectedResult, winErr.Error())
		}

		// WindowsError
		root = WindowsError
		expectedResult = "[windows_error] Sample error message"
		winErr = Errorf(root, msg)

		if winErr.Error() != expectedResult {
			t.Errorf("Expected error message: '%s'\nGot: '%s'", expectedResult, winErr.Error())
		}
	})
}
