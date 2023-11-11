package local

import (
	"testing"

	"github.com/d-strobel/gowindows/connection"
)

func TestNew(t *testing.T) {
	// Create a mock connection for testing
	mockConn := &connection.Connection{}

	// Call the New function with the mock connection
	client := New(mockConn)

	// Check if the Connection field in the returned Client is the same as the mockConn
	if client.Connection != mockConn {
		t.Errorf("New function did not set Connection field correctly")
	}
}
