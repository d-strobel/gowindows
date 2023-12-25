package local

import (
	"testing"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
	"github.com/stretchr/testify/suite"
)

// Unit test suite for all local functions
type LocalUnitTestSuite struct {
	suite.Suite
}

// Run all local unit tests
func TestLocalUnitTestSuite(t *testing.T) {
	suite.Run(t, &LocalUnitTestSuite{})
}

func (suite *LocalUnitTestSuite) TestNewLocalClient() {
	mockConn := &connection.Connection{}
	mockParser := &parser.Parser{}
	actualLocalClient := NewLocalClient(mockConn, mockParser)
	expectedLocalClient := &LocalClient{Connection: mockConn, parser: mockParser}

	suite.Equal(expectedLocalClient, actualLocalClient)
}
