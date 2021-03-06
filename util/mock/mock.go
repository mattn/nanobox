// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/nanobox-io/nanobox/util (interfaces: Util)

package mock_util

import gomock "github.com/golang/mock/gomock"

// Mock of Util interface
type MockUtil struct {
	ctrl     *gomock.Controller
	recorder *_MockUtilRecorder
}

// Recorder for MockUtil (not exported)
type _MockUtilRecorder struct {
	mock *MockUtil
}

func NewMockUtil(ctrl *gomock.Controller) *MockUtil {
	mock := &MockUtil{ctrl: ctrl}
	mock.recorder = &_MockUtilRecorder{mock}
	return mock
}

func (_m *MockUtil) EXPECT() *_MockUtilRecorder {
	return _m.recorder
}

func (_m *MockUtil) MD5sMatch(_param0 string, _param1 string) (bool, error) {
	ret := _m.ctrl.Call(_m, "MD5sMatch", _param0, _param1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockUtilRecorder) MD5sMatch(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MD5sMatch", arg0, arg1)
}
