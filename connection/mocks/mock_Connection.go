// Code generated by mockery. DO NOT EDIT.

package connection

import (
	context "context"

	connection "github.com/d-strobel/gowindows/connection"

	mock "github.com/stretchr/testify/mock"
)

// MockConnection is an autogenerated mock type for the Connection type
type MockConnection struct {
	mock.Mock
}

type MockConnection_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConnection) EXPECT() *MockConnection_Expecter {
	return &MockConnection_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *MockConnection) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockConnection_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockConnection_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockConnection_Expecter) Close() *MockConnection_Close_Call {
	return &MockConnection_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockConnection_Close_Call) Run(run func()) *MockConnection_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConnection_Close_Call) Return(_a0 error) *MockConnection_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConnection_Close_Call) RunAndReturn(run func() error) *MockConnection_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Run provides a mock function with given fields: ctx, cmd
func (_m *MockConnection) Run(ctx context.Context, cmd string) (connection.CmdResult, error) {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for Run")
	}

	var r0 connection.CmdResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (connection.CmdResult, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) connection.CmdResult); ok {
		r0 = rf(ctx, cmd)
	} else {
		r0 = ret.Get(0).(connection.CmdResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockConnection_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type MockConnection_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
//   - ctx context.Context
//   - cmd string
func (_e *MockConnection_Expecter) Run(ctx interface{}, cmd interface{}) *MockConnection_Run_Call {
	return &MockConnection_Run_Call{Call: _e.mock.On("Run", ctx, cmd)}
}

func (_c *MockConnection_Run_Call) Run(run func(ctx context.Context, cmd string)) *MockConnection_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockConnection_Run_Call) Return(_a0 connection.CmdResult, _a1 error) *MockConnection_Run_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockConnection_Run_Call) RunAndReturn(run func(context.Context, string) (connection.CmdResult, error)) *MockConnection_Run_Call {
	_c.Call.Return(run)
	return _c
}

// RunWithPowershell provides a mock function with given fields: ctx, cmd
func (_m *MockConnection) RunWithPowershell(ctx context.Context, cmd string) (connection.CmdResult, error) {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for RunWithPowershell")
	}

	var r0 connection.CmdResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (connection.CmdResult, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) connection.CmdResult); ok {
		r0 = rf(ctx, cmd)
	} else {
		r0 = ret.Get(0).(connection.CmdResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockConnection_RunWithPowershell_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RunWithPowershell'
type MockConnection_RunWithPowershell_Call struct {
	*mock.Call
}

// RunWithPowershell is a helper method to define mock.On call
//   - ctx context.Context
//   - cmd string
func (_e *MockConnection_Expecter) RunWithPowershell(ctx interface{}, cmd interface{}) *MockConnection_RunWithPowershell_Call {
	return &MockConnection_RunWithPowershell_Call{Call: _e.mock.On("RunWithPowershell", ctx, cmd)}
}

func (_c *MockConnection_RunWithPowershell_Call) Run(run func(ctx context.Context, cmd string)) *MockConnection_RunWithPowershell_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockConnection_RunWithPowershell_Call) Return(_a0 connection.CmdResult, _a1 error) *MockConnection_RunWithPowershell_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockConnection_RunWithPowershell_Call) RunAndReturn(run func(context.Context, string) (connection.CmdResult, error)) *MockConnection_RunWithPowershell_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockConnection creates a new instance of MockConnection. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConnection(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConnection {
	mock := &MockConnection{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
