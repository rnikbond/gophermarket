// Code generated by MockGen. DO NOT EDIT.
// Source: order.go

// Package order is a generated GoMock package.
package order

import (
	context "context"
	order "gophermarket/pkg/order"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockServiceOrder is a mock of ServiceOrder interface.
type MockServiceOrder struct {
	ctrl     *gomock.Controller
	recorder *MockServiceOrderMockRecorder
}

// MockServiceOrderMockRecorder is the mock recorder for MockServiceOrder.
type MockServiceOrderMockRecorder struct {
	mock *MockServiceOrder
}

// NewMockServiceOrder creates a new mock instance.
func NewMockServiceOrder(ctrl *gomock.Controller) *MockServiceOrder {
	mock := &MockServiceOrder{ctrl: ctrl}
	mock.recorder = &MockServiceOrderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceOrder) EXPECT() *MockServiceOrderMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockServiceOrder) Create(ctx context.Context, number int64, username string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, number, username)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockServiceOrderMockRecorder) Create(ctx, number, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockServiceOrder)(nil).Create), ctx, number, username)
}

// CreateWithPayment mocks base method.
func (m *MockServiceOrder) CreateWithPayment(ctx context.Context, number int64, username string, sum float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWithPayment", ctx, number, username, sum)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateWithPayment indicates an expected call of CreateWithPayment.
func (mr *MockServiceOrderMockRecorder) CreateWithPayment(ctx, number, username, sum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWithPayment", reflect.TypeOf((*MockServiceOrder)(nil).CreateWithPayment), ctx, number, username, sum)
}

// UserOrders mocks base method.
func (m *MockServiceOrder) UserOrders(ctx context.Context, username string) ([]order.InfoOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserOrders", ctx, username)
	ret0, _ := ret[0].([]order.InfoOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserOrders indicates an expected call of UserOrders.
func (mr *MockServiceOrderMockRecorder) UserOrders(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserOrders", reflect.TypeOf((*MockServiceOrder)(nil).UserOrders), ctx, username)
}
