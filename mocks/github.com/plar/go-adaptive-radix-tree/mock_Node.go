// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_art

import (
	art "github.com/plar/go-adaptive-radix-tree"
	mock "github.com/stretchr/testify/mock"
)

// MockNode is an autogenerated mock type for the Node type
type MockNode struct {
	mock.Mock
}

type MockNode_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNode) EXPECT() *MockNode_Expecter {
	return &MockNode_Expecter{mock: &_m.Mock}
}

// Key provides a mock function with no fields
func (_m *MockNode) Key() art.Key {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Key")
	}

	var r0 art.Key
	if rf, ok := ret.Get(0).(func() art.Key); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(art.Key)
		}
	}

	return r0
}

// MockNode_Key_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Key'
type MockNode_Key_Call struct {
	*mock.Call
}

// Key is a helper method to define mock.On call
func (_e *MockNode_Expecter) Key() *MockNode_Key_Call {
	return &MockNode_Key_Call{Call: _e.mock.On("Key")}
}

func (_c *MockNode_Key_Call) Run(run func()) *MockNode_Key_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockNode_Key_Call) Return(_a0 art.Key) *MockNode_Key_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNode_Key_Call) RunAndReturn(run func() art.Key) *MockNode_Key_Call {
	_c.Call.Return(run)
	return _c
}

// Kind provides a mock function with no fields
func (_m *MockNode) Kind() art.Kind {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Kind")
	}

	var r0 art.Kind
	if rf, ok := ret.Get(0).(func() art.Kind); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(art.Kind)
	}

	return r0
}

// MockNode_Kind_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Kind'
type MockNode_Kind_Call struct {
	*mock.Call
}

// Kind is a helper method to define mock.On call
func (_e *MockNode_Expecter) Kind() *MockNode_Kind_Call {
	return &MockNode_Kind_Call{Call: _e.mock.On("Kind")}
}

func (_c *MockNode_Kind_Call) Run(run func()) *MockNode_Kind_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockNode_Kind_Call) Return(_a0 art.Kind) *MockNode_Kind_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNode_Kind_Call) RunAndReturn(run func() art.Kind) *MockNode_Kind_Call {
	_c.Call.Return(run)
	return _c
}

// Value provides a mock function with no fields
func (_m *MockNode) Value() art.Value {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Value")
	}

	var r0 art.Value
	if rf, ok := ret.Get(0).(func() art.Value); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(art.Value)
		}
	}

	return r0
}

// MockNode_Value_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Value'
type MockNode_Value_Call struct {
	*mock.Call
}

// Value is a helper method to define mock.On call
func (_e *MockNode_Expecter) Value() *MockNode_Value_Call {
	return &MockNode_Value_Call{Call: _e.mock.On("Value")}
}

func (_c *MockNode_Value_Call) Run(run func()) *MockNode_Value_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockNode_Value_Call) Return(_a0 art.Value) *MockNode_Value_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNode_Value_Call) RunAndReturn(run func() art.Value) *MockNode_Value_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockNode creates a new instance of MockNode. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNode(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNode {
	mock := &MockNode{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
