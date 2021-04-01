// Code generated by MockGen. DO NOT EDIT.
// Source: ../contract.go

// Package mock_configuration is a generated GoMock package.
package mock_configuration

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockConfigurationContract is a mock of ConfigurationContract interface
type MockConfigurationContract struct {
	ctrl     *gomock.Controller
	recorder *MockConfigurationContractMockRecorder
}

// MockConfigurationContractMockRecorder is the mock recorder for MockConfigurationContract
type MockConfigurationContractMockRecorder struct {
	mock *MockConfigurationContract
}

// NewMockConfigurationContract creates a new mock instance
func NewMockConfigurationContract(ctrl *gomock.Controller) *MockConfigurationContract {
	mock := &MockConfigurationContract{ctrl: ctrl}
	mock.recorder = &MockConfigurationContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigurationContract) EXPECT() *MockConfigurationContractMockRecorder {
	return m.recorder
}

// GetHttpHost mocks base method
func (m *MockConfigurationContract) GetHttpHost() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHttpHost")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHttpHost indicates an expected call of GetHttpHost
func (mr *MockConfigurationContractMockRecorder) GetHttpHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHttpHost", reflect.TypeOf((*MockConfigurationContract)(nil).GetHttpHost))
}

// GetHttpPort mocks base method
func (m *MockConfigurationContract) GetHttpPort() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHttpPort")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHttpPort indicates an expected call of GetHttpPort
func (mr *MockConfigurationContractMockRecorder) GetHttpPort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHttpPort", reflect.TypeOf((*MockConfigurationContract)(nil).GetHttpPort))
}

// GetRunningNodeName mocks base method
func (m *MockConfigurationContract) GetRunningNodeName() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRunningNodeName")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRunningNodeName indicates an expected call of GetRunningNodeName
func (mr *MockConfigurationContractMockRecorder) GetRunningNodeName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunningNodeName", reflect.TypeOf((*MockConfigurationContract)(nil).GetRunningNodeName))
}

// GetGeolocationUpdaterCronSpec mocks base method
func (m *MockConfigurationContract) GetGeolocationUpdaterCronSpec() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGeolocationUpdaterCronSpec")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGeolocationUpdaterCronSpec indicates an expected call of GetGeolocationUpdaterCronSpec
func (mr *MockConfigurationContractMockRecorder) GetGeolocationUpdaterCronSpec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGeolocationUpdaterCronSpec", reflect.TypeOf((*MockConfigurationContract)(nil).GetGeolocationUpdaterCronSpec))
}

// GetIpinfoUrl mocks base method
func (m *MockConfigurationContract) GetIpinfoUrl() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIpinfoUrl")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIpinfoUrl indicates an expected call of GetIpinfoUrl
func (mr *MockConfigurationContractMockRecorder) GetIpinfoUrl() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIpinfoUrl", reflect.TypeOf((*MockConfigurationContract)(nil).GetIpinfoUrl))
}

// GetIpinfoAccessToken mocks base method
func (m *MockConfigurationContract) GetIpinfoAccessToken() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIpinfoAccessToken")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIpinfoAccessToken indicates an expected call of GetIpinfoAccessToken
func (mr *MockConfigurationContractMockRecorder) GetIpinfoAccessToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIpinfoAccessToken", reflect.TypeOf((*MockConfigurationContract)(nil).GetIpinfoAccessToken))
}