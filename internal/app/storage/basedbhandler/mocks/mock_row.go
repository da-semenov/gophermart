// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/da-semenov/gophermart/internal/app/storage/basedbhandler (interfaces: Row)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRow is a mock of Row interface.
type MockRow struct {
	ctrl     *gomock.Controller
	recorder *MockRowMockRecorder
}

// MockRowMockRecorder is the mock recorder for MockRow.
type MockRowMockRecorder struct {
	mock *MockRow
}

// NewMockRow creates a new mock instance.
func NewMockRow(ctrl *gomock.Controller) *MockRow {
	mock := &MockRow{ctrl: ctrl}
	mock.recorder = &MockRowMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRow) EXPECT() *MockRowMockRecorder {
	return m.recorder
}

// Scan mocks base method.
func (m *MockRow) Scan(arg0 ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockRowMockRecorder) Scan(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockRow)(nil).Scan), arg0...)
}

