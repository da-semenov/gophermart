// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/da-semenov/gophermart/internal/app/handlers (interfaces: BalanceService)

// Package mock_handlers is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/da-semenov/gophermart/internal/app/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockBalanceService is a mock of BalanceService interface.
type MockBalanceService struct {
	ctrl     *gomock.Controller
	recorder *MockBalanceServiceMockRecorder
}

// MockBalanceServiceMockRecorder is the mock recorder for MockBalanceService.
type MockBalanceServiceMockRecorder struct {
	mock *MockBalanceService
}

// NewMockBalanceService creates a new mock instance.
func NewMockBalanceService(ctrl *gomock.Controller) *MockBalanceService {
	mock := &MockBalanceService{ctrl: ctrl}
	mock.recorder = &MockBalanceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalanceService) EXPECT() *MockBalanceServiceMockRecorder {
	return m.recorder
}

// GetCurrentBalance mocks base method.
func (m *MockBalanceService) GetCurrentBalance(arg0 context.Context, arg1 int) (*domain.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentBalance", arg0, arg1)
	ret0, _ := ret[0].(*domain.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentBalance indicates an expected call of GetCurrentBalance.
func (mr *MockBalanceServiceMockRecorder) GetCurrentBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentBalance", reflect.TypeOf((*MockBalanceService)(nil).GetCurrentBalance), arg0, arg1)
}

// GetWithdrawalsList mocks base method.
func (m *MockBalanceService) GetWithdrawalsList(arg0 context.Context, arg1 int) ([]domain.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawalsList", arg0, arg1)
	ret0, _ := ret[0].([]domain.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawalsList indicates an expected call of GetWithdrawalsList.
func (mr *MockBalanceServiceMockRecorder) GetWithdrawalsList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawalsList", reflect.TypeOf((*MockBalanceService)(nil).GetWithdrawalsList), arg0, arg1)
}

// Withdraw mocks base method.
func (m *MockBalanceService) Withdraw(arg0 context.Context, arg1 *domain.Withdraw, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockBalanceServiceMockRecorder) Withdraw(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockBalanceService)(nil).Withdraw), arg0, arg1, arg2)
}

