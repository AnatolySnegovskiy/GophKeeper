// Code generated by MockGen. DO NOT EDIT.
// Source: services/entities/file_helper.go
//
// Generated by this command:
//
//	mockgen -source=services/entities/file_helper.go -destination=mocks/file_helper_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	"github.com/golang/mock/gomock"
	os "os"
	reflect "reflect"
)

// MockFileEntity is a mock of FileEntity interface.
type MockFileEntity struct {
	ctrl     *gomock.Controller
	recorder *MockFileEntityMockRecorder
}

// MockFileEntityMockRecorder is the mock recorder for MockFileEntity.
type MockFileEntityMockRecorder struct {
	mock *MockFileEntity
}

// NewMockFileEntity creates a new mock instance.
func NewMockFileEntity(ctrl *gomock.Controller) *MockFileEntity {
	mock := &MockFileEntity{ctrl: ctrl}
	mock.recorder = &MockFileEntityMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileEntity) EXPECT() *MockFileEntityMockRecorder {
	return m.recorder
}

// FromFile mocks base method.
func (m *MockFileEntity) FromFile(file *os.File) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FromFile", file)
	ret0, _ := ret[0].(error)
	return ret0
}

// FromFile indicates an expected call of FromFile.
func (mr *MockFileEntityMockRecorder) FromFile(file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromFile", reflect.TypeOf((*MockFileEntity)(nil).FromFile), file)
}

// GetName mocks base method.
func (m *MockFileEntity) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockFileEntityMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockFileEntity)(nil).GetName))
}

// ToFile mocks base method.
func (m *MockFileEntity) ToFile() (*os.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToFile")
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToFile indicates an expected call of ToFile.
func (mr *MockFileEntityMockRecorder) ToFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToFile", reflect.TypeOf((*MockFileEntity)(nil).ToFile))
}
