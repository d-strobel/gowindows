package parsing

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Unit test suite for all Powershell parsing functions
type PowershellUnitTestSuite struct {
	suite.Suite
}

func TestPowershellUnitTestSuite(t *testing.T) {
	suite.Run(t, &PowershellUnitTestSuite{})
}

func (suite *PowershellUnitTestSuite) TestEncodePwshCmd() {
	suite.Run("should return the correct encoded powershell string", func() {
		cmd := "Get-LocalUser"
		expectedPwshCmd := "powershell.exe -NoProfile -EncodedCommand JABQAHIAbwBnAHIAZQBzAHMAUAByAGUAZgBlAHIAZQBuAGMAZQAgAD0AIAAnAFMAaQBsAGUAbgB0AGwAeQBDAG8AbgB0AGkAbgB1AGUAJwA7ACAARwBlAHQALQBMAG8AYwBhAGwAVQBzAGUAcgA="
		actualPwshCmd, err := EncodePwshCmd(cmd)
		suite.NoError(err)
		suite.Equal(expectedPwshCmd, actualPwshCmd)
	})
}
