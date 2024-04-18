package connection_test

import (
	"github.com/d-strobel/gowindows/connection"
)

func (suite *ConnectionAccTestSuite) TestNewConnectionWithWinRM() {
	suite.T().Parallel()

	suite.Run("should establish a connection via password", func() {
		winRMConfig := connection.WinRMConfig{
			WinRMHost:     suite.host,
			WinRMPort:     suite.winRMPort,
			WinRMUsername: suite.username,
			WinRMPassword: suite.password,
		}

		conn, err := connection.NewConnectionWithWinRM(&winRMConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})
}
