package local

import (
	"testing"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
	"github.com/stretchr/testify/assert"
)

func TestNewLocalClient(t *testing.T) {
	mockConn := &connection.Connection{}
	mockParser := &parser.Parser{}
	actualLocalClient := NewLocalClient(mockConn, mockParser)
	expectedLocalClient := &LocalClient{Connection: mockConn, parser: mockParser}

	assert.Equal(t, expectedLocalClient, actualLocalClient)
}
