package gowindows

import (
	"testing"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/package/local"
	"github.com/stretchr/testify/mock"
)

// MockConnection is a mock implementation of connection.Connection
type MockConnection struct {
	mock.Mock
}

// MockNew is a mock implementation of connection.New
func (m *MockConnection) New(conf *connection.Config) (*connection.Connection, error) {
	args := m.Called(conf)
	return args.Get(0).(*connection.Connection), args.Error(1)
}

// MockClient is a mock implementation of local.Client
type MockLocalClient struct {
	mock.Mock
}

// MockNew is a mock implementation of local.New
func (m *MockLocalClient) New(conn *connection.Connection) *local.Client {
	args := m.Called(conn)
	return args.Get(0).(*local.Client)
}

func TestNew(t *testing.T) {

	// Create a mock connection
	mockConn := new(MockConnection)

	// Create a mock of packae clients
	mockLocalClient := new(MockLocalClient)

	// WinRM test config for mock
	winRMConfigTest := connection.WinRMConfig{
		WinRMUsername: "test",
		WinRMPassword: "test",
		WinRMHost:     "test",
	}

	// Configure the expected behavior of the mocks
	mockConn.On("New", mock.Anything).Return(&connection.Connection{}, nil)
	mockLocalClient.On("New", mock.Anything).Return(&local.Client{})

	// Call the New function with the mock objects
	conn, err := mockConn.New(&connection.Config{WinRM: &winRMConfigTest})
	_ = mockLocalClient.New(conn)

	// Check if the NewClient function returned an error
	if err != nil {
		t.Errorf("NewClient returned an error: %v", err)
	}

	// Assert that the mocks were called as expected
	mockConn.AssertExpectations(t)
	mockLocalClient.AssertExpectations(t)
}
