package ssh

import (
	"os/user"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Unit test suite for all SSH functions.
type SSHUnitTestSuite struct {
	suite.Suite

	// Fixtures
	currentUserHomeDir string
}

func TestSSHUnitTestSuite(t *testing.T) {
	suite.Run(t, &SSHUnitTestSuite{})
}

// SetupSUite setups all neccessary fixtures for running the unit tests.
func (suite *SSHUnitTestSuite) SetupSuite() {
	// Get current user
	user, err := user.Current()
	suite.Require().NoError(err)
	suite.currentUserHomeDir = user.HomeDir
}
