package winrm

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Unit test suite for all WinRM functions.
type WinRMUnitTestSuite struct {
	suite.Suite
}

func TestWinRMUnitTestSuite(t *testing.T) {
	suite.Run(t, &WinRMUnitTestSuite{})
}
