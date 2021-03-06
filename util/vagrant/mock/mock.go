// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/nanobox-io/nanobox/util/vagrant (interfaces: Vagrant)

package mock_vagrant

import gomock "github.com/golang/mock/gomock"

// Mock of Vagrant interface
type MockVagrant struct {
	ctrl     *gomock.Controller
	recorder *_MockVagrantRecorder
}

// Recorder for MockVagrant (not exported)
type _MockVagrantRecorder struct {
	mock *MockVagrant
}

func NewMockVagrant(ctrl *gomock.Controller) *MockVagrant {
	mock := &MockVagrant{ctrl: ctrl}
	mock.recorder = &_MockVagrantRecorder{mock}
	return mock
}

func (_m *MockVagrant) EXPECT() *_MockVagrantRecorder {
	return _m.recorder
}

func (_m *MockVagrant) Destroy() error {
	ret := _m.ctrl.Call(_m, "Destroy")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVagrantRecorder) Destroy() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Destroy")
}

func (_m *MockVagrant) Exists() bool {
	ret := _m.ctrl.Call(_m, "Exists")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockVagrantRecorder) Exists() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Exists")
}

func (_m *MockVagrant) Init() {
	_m.ctrl.Call(_m, "Init")
}

func (_mr *_MockVagrantRecorder) Init() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Init")
}

func (_m *MockVagrant) Install() error {
	ret := _m.ctrl.Call(_m, "Install")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVagrantRecorder) Install() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Install")
}

func (_m *MockVagrant) NewLogger(_param0 string) {
	_m.ctrl.Call(_m, "NewLogger", _param0)
}

func (_mr *_MockVagrantRecorder) NewLogger(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewLogger", arg0)
}

func (_m *MockVagrant) Reload() error {
	ret := _m.ctrl.Call(_m, "Reload")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVagrantRecorder) Reload() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Reload")
}

func (_m *MockVagrant) Resume() error {
	ret := _m.ctrl.Call(_m, "Resume")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVagrantRecorder) Resume() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Resume")
}

func (_m *MockVagrant) SSH() error {
	ret := _m.ctrl.Call(_m, "SSH")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVagrantRecorder) SSH() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SSH")
}

func (_m *MockVagrant) Status() string {
	ret := _m.ctrl.Call(_m, "Status")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockVagrantRecorder) Status() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Status")
}

func (_m *MockVagrant) Suspend() error {
	ret := _m.ctrl.Call(_m, "Suspend")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVagrantRecorder) Suspend() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Suspend")
}

func (_m *MockVagrant) Up() error {
	ret := _m.ctrl.Call(_m, "Up")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVagrantRecorder) Up() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Up")
}

func (_m *MockVagrant) Update() error {
	ret := _m.ctrl.Call(_m, "Update")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVagrantRecorder) Update() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Update")
}
