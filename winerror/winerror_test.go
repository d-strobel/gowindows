package winerror

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

// Unit test suite for all WinError functions
type WinErrorUnitTestSuite struct {
	suite.Suite
}

func TestWinErrorUnitTestSuite(t *testing.T) {
	suite.Run(t, &WinErrorUnitTestSuite{})
}

func (suite *WinErrorUnitTestSuite) TestError() {
	suite.Run("should return the error message from the WinError object", func() {
		err := &WinError{
			Command: "test-command",
			Err:     errors.New("error-message"),
		}
		suite.Error(err)
		suite.Equal(err.Error(), "error-message")
	})
}

func (suite *WinErrorUnitTestSuite) TestUnwrap() {
	suite.Run("should return the unwraped error message from the WinError object", func() {
		err := &WinError{
			Command: "test-command",
			Err:     errors.New("error-message"),
		}
		suite.Equal(err.Unwrap(), errors.New("error-message"))
	})
}

func (suite *WinErrorUnitTestSuite) TestNew() {
	suite.Run("should create a new WinError object", func() {
		expectedErr := &WinError{
			Command: "test-command",
			Err:     errors.New("error-message"),
		}
		err := New("test-command", errors.New("error-message"))
		suite.Error(err)
		suite.Equal(expectedErr, err)
	})
}

func (suite *WinErrorUnitTestSuite) TestErrorf() {
	suite.T().Parallel()

	suite.Run("should create a new WinError with default string", func() {
		expectedErr := &WinError{
			Command: "test-command",
			Err:     errors.New("error-message"),
		}
		err := Errorf("test-command", "error-message")
		suite.Error(err)
		suite.Equal(expectedErr, err)
	})

	suite.Run("should create a new WinError with format string", func() {
		expectedErr := &WinError{
			Command: "test-command",
			Err:     errors.New("error-message test"),
		}
		err := Errorf("test-command", "error-message %s", "test")
		suite.Error(err)
		suite.Equal(expectedErr, err)
	})

	suite.Run("should create a new WinError with format string and multiple values", func() {
		expectedErr := &WinError{
			Command: "test-command",
			Err:     errors.New("error-message test test2"),
		}
		err := Errorf("test-command", "error-message %s %s", "test", "test2")
		suite.Error(err)
		suite.Equal(expectedErr, err)
	})
}

func (suite *WinErrorUnitTestSuite) TestUnwrapCommand() {
	suite.T().Parallel()

	suite.Run("should return the unwrapped command from the WinError object", func() {
		err := &WinError{
			Command: "test-command",
			Err:     errors.New("error-message"),
		}
		suite.Equal(UnwrapCommand(err), "test-command")
	})

	suite.Run("should return empty string when error is not a WinError object", func() {
		err := errors.New("error-message")
		suite.Equal(UnwrapCommand(err), "")
	})
}
