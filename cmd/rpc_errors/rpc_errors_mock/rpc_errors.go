// Code generated by MockGen. DO NOT EDIT.
// Source: rpc_errors.pb.go

// Package rpc_errors_mock is a generated GoMock package.
package rpc_errors_mock

import (
	x "."
	context "context"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockGreeterClient is a mock of GreeterClient interface
type MockGreeterClient struct {
	ctrl     *gomock.Controller
	recorder *MockGreeterClientMockRecorder
}

// MockGreeterClientMockRecorder is the mock recorder for MockGreeterClient
type MockGreeterClientMockRecorder struct {
	mock *MockGreeterClient
}

// NewMockGreeterClient creates a new mock instance
func NewMockGreeterClient(ctrl *gomock.Controller) *MockGreeterClient {
	mock := &MockGreeterClient{ctrl: ctrl}
	mock.recorder = &MockGreeterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGreeterClient) EXPECT() *MockGreeterClientMockRecorder {
	return m.recorder
}

// SayHello mocks base method
func (m *MockGreeterClient) SayHello(ctx context.Context, in *x.HelloRequest, opts ...grpc.CallOption) (*x.HelloReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SayHello", varargs...)
	ret0, _ := ret[0].(*x.HelloReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SayHello indicates an expected call of SayHello
func (mr *MockGreeterClientMockRecorder) SayHello(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockGreeterClient)(nil).SayHello), varargs...)
}

// MockGreeterServer is a mock of GreeterServer interface
type MockGreeterServer struct {
	ctrl     *gomock.Controller
	recorder *MockGreeterServerMockRecorder
}

// MockGreeterServerMockRecorder is the mock recorder for MockGreeterServer
type MockGreeterServerMockRecorder struct {
	mock *MockGreeterServer
}

// NewMockGreeterServer creates a new mock instance
func NewMockGreeterServer(ctrl *gomock.Controller) *MockGreeterServer {
	mock := &MockGreeterServer{ctrl: ctrl}
	mock.recorder = &MockGreeterServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGreeterServer) EXPECT() *MockGreeterServerMockRecorder {
	return m.recorder
}

// SayHello mocks base method
func (m *MockGreeterServer) SayHello(arg0 context.Context, arg1 *x.HelloRequest) (*x.HelloReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SayHello", arg0, arg1)
	ret0, _ := ret[0].(*x.HelloReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SayHello indicates an expected call of SayHello
func (mr *MockGreeterServerMockRecorder) SayHello(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockGreeterServer)(nil).SayHello), arg0, arg1)
}
