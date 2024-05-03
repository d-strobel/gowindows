package gowindows

import (
	"testing"

	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	"github.com/d-strobel/gowindows/windows/local"
	"github.com/stretchr/testify/suite"
)

// Unit test suite for gowindows
type GowindowsUnitTestSuite struct {
	suite.Suite
}

// Run all gowindows unit tests
func TestGowindowsUnitTestSuite(t *testing.T) {
	suite.Run(t, &GowindowsUnitTestSuite{})
}

func (suite *GowindowsUnitTestSuite) TestNewClient() {
	suite.Run("should return a new client", func() {
		mockConn := mockConnection.NewMockConnection(suite.T())

		expectedClient := &Client{
			Connection: mockConn,
			Local:      local.NewClient(mockConn),
		}

		actualClient, err := NewClient(mockConn)
		suite.NoError(err)
		suite.Equal(expectedClient, actualClient)
	})
}
