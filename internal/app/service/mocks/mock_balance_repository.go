// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/da-semenov/gophermart/internal/app/models (interfaces: BalanceRepository)

// Package mock_models is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/da-semenov/gophermart/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockBalanceRepository is a mock of BalanceRepository interface.
type MockBalanceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBalanceRepositoryMockRecorder
}

// MockBalanceRepositoryMockRecorder is the mock recorder for MockBalanceRepository.
type MockBalanceRepositoryMockRecorder struct {
	mock *MockBalanceRepository
}

// NewMockBalanceRepository creates a new mock instance.
func NewMockBalanceRepository(ctrl *gomock.Controller) *MockBalanceRepository {
	mock := &MockBalanceRepository{ctrl: ctrl}
	mock.recorder = &MockBalanceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalanceRepository) EXPECT() *MockBalanceRepositoryMockRecorder {
	return m.recorder
}

// CreateOperation mocks base method.
func (m *MockBalanceRepository) CreateOperation(arg0 context.Context, arg1 *models.Operation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOperation", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOperation indicates an expected call of CreateOperation.
func (mr *MockBalanceRepositoryMockRecorder) CreateOperation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOperation", reflect.TypeOf((*MockBalanceRepository)(nil).CreateOperation), arg0, arg1)
}

// FindWithdrawalByUser mocks base method.
func (m *MockBalanceRepository) FindWithdrawalByUser(arg0 context.Context, arg1 int) ([]models.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindWithdrawalByUser", arg0, arg1)
	ret0, _ := ret[0].([]models.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindWithdrawalByUser indicates an expected call of FindWithdrawalByUser.
func (mr *MockBalanceRepositoryMockRecorder) FindWithdrawalByUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindWithdrawalByUser", reflect.TypeOf((*MockBalanceRepository)(nil).FindWithdrawalByUser), arg0, arg1)
}

// GetAccount mocks base method.
func (m *MockBalanceRepository) GetAccount(arg0 context.Context, arg1 int) (*models.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1)
	ret0, _ := ret[0].(*models.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockBalanceRepositoryMockRecorder) GetAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockBalanceRepository)(nil).GetAccount), arg0, arg1)
}

// LockAccount mocks base method.
func (m *MockBalanceRepository) LockAccount(arg0 context.Context, arg1 int) (*models.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockAccount", arg0, arg1)
	ret0, _ := ret[0].(*models.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LockAccount indicates an expected call of LockAccount.
func (mr *MockBalanceRepositoryMockRecorder) LockAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockAccount", reflect.TypeOf((*MockBalanceRepository)(nil).LockAccount), arg0, arg1)
}

// SaveAccount mocks base method.
func (m *MockBalanceRepository) SaveAccount(arg0 context.Context, arg1 *models.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAccount indicates an expected call of SaveAccount.
func (mr *MockBalanceRepositoryMockRecorder) SaveAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAccount", reflect.TypeOf((*MockBalanceRepository)(nil).SaveAccount), arg0, arg1)
}
