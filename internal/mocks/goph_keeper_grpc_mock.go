// Code generated by MockGen. DO NOT EDIT.
// Source: services/grpc/goph_keeper/v1/goph_keeper_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=services/grpc/goph_keeper/v1/goph_keeper_grpc.pb.go -destination=mocks/goph_keeper_grpc_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	reflect "reflect"

	"github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockGophKeeperV1ServiceClient is a mock of GophKeeperV1ServiceClient interface.
type MockGophKeeperV1ServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockGophKeeperV1ServiceClientMockRecorder
}

// MockGophKeeperV1ServiceClientMockRecorder is the mock recorder for MockGophKeeperV1ServiceClient.
type MockGophKeeperV1ServiceClientMockRecorder struct {
	mock *MockGophKeeperV1ServiceClient
}

// NewMockGophKeeperV1ServiceClient creates a new mock instance.
func NewMockGophKeeperV1ServiceClient(ctrl *gomock.Controller) *MockGophKeeperV1ServiceClient {
	mock := &MockGophKeeperV1ServiceClient{ctrl: ctrl}
	mock.recorder = &MockGophKeeperV1ServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGophKeeperV1ServiceClient) EXPECT() *MockGophKeeperV1ServiceClientMockRecorder {
	return m.recorder
}

// AuthenticateUser mocks base method.
func (m *MockGophKeeperV1ServiceClient) AuthenticateUser(ctx context.Context, in *v1.AuthenticateUserRequest, opts ...grpc.CallOption) (*v1.AuthenticateUserResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AuthenticateUser", varargs...)
	ret0, _ := ret[0].(*v1.AuthenticateUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateUser indicates an expected call of AuthenticateUser.
func (mr *MockGophKeeperV1ServiceClientMockRecorder) AuthenticateUser(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateUser", reflect.TypeOf((*MockGophKeeperV1ServiceClient)(nil).AuthenticateUser), varargs...)
}

// DeleteFile mocks base method.
func (m *MockGophKeeperV1ServiceClient) DeleteFile(ctx context.Context, in *v1.DeleteFileRequest, opts ...grpc.CallOption) (*v1.DeleteFileResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteFile", varargs...)
	ret0, _ := ret[0].(*v1.DeleteFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockGophKeeperV1ServiceClientMockRecorder) DeleteFile(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockGophKeeperV1ServiceClient)(nil).DeleteFile), varargs...)
}

// DownloadFile mocks base method.
func (m *MockGophKeeperV1ServiceClient) DownloadFile(ctx context.Context, in *v1.DownloadFileRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[v1.DownloadFileResponse], error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DownloadFile", varargs...)
	ret0, _ := ret[0].(grpc.ServerStreamingClient[v1.DownloadFileResponse])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadFile indicates an expected call of DownloadFile.
func (mr *MockGophKeeperV1ServiceClientMockRecorder) DownloadFile(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockGophKeeperV1ServiceClient)(nil).DownloadFile), varargs...)
}

// GetMetadataFile mocks base method.
func (m *MockGophKeeperV1ServiceClient) GetMetadataFile(ctx context.Context, in *v1.GetMetadataFileRequest, opts ...grpc.CallOption) (*v1.GetMetadataFileResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMetadataFile", varargs...)
	ret0, _ := ret[0].(*v1.GetMetadataFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadataFile indicates an expected call of GetMetadataFile.
func (mr *MockGophKeeperV1ServiceClientMockRecorder) GetMetadataFile(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadataFile", reflect.TypeOf((*MockGophKeeperV1ServiceClient)(nil).GetMetadataFile), varargs...)
}

// GetStoreDataList mocks base method.
func (m *MockGophKeeperV1ServiceClient) GetStoreDataList(ctx context.Context, in *v1.GetStoreDataListRequest, opts ...grpc.CallOption) (*v1.GetStoreDataListResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetStoreDataList", varargs...)
	ret0, _ := ret[0].(*v1.GetStoreDataListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStoreDataList indicates an expected call of GetStoreDataList.
func (mr *MockGophKeeperV1ServiceClientMockRecorder) GetStoreDataList(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStoreDataList", reflect.TypeOf((*MockGophKeeperV1ServiceClient)(nil).GetStoreDataList), varargs...)
}

// RegisterUser mocks base method.
func (m *MockGophKeeperV1ServiceClient) RegisterUser(ctx context.Context, in *v1.RegisterUserRequest, opts ...grpc.CallOption) (*v1.RegisterUserResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RegisterUser", varargs...)
	ret0, _ := ret[0].(*v1.RegisterUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockGophKeeperV1ServiceClientMockRecorder) RegisterUser(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockGophKeeperV1ServiceClient)(nil).RegisterUser), varargs...)
}

// SetMetadataFile mocks base method.
func (m *MockGophKeeperV1ServiceClient) SetMetadataFile(ctx context.Context, in *v1.SetMetadataFileRequest, opts ...grpc.CallOption) (*v1.SetMetadataFileResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetMetadataFile", varargs...)
	ret0, _ := ret[0].(*v1.SetMetadataFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetMetadataFile indicates an expected call of SetMetadataFile.
func (mr *MockGophKeeperV1ServiceClientMockRecorder) SetMetadataFile(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMetadataFile", reflect.TypeOf((*MockGophKeeperV1ServiceClient)(nil).SetMetadataFile), varargs...)
}

// UploadFile mocks base method.
func (m *MockGophKeeperV1ServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[v1.UploadFileRequest, v1.UploadFileResponse], error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UploadFile", varargs...)
	ret0, _ := ret[0].(grpc.ClientStreamingClient[v1.UploadFileRequest, v1.UploadFileResponse])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockGophKeeperV1ServiceClientMockRecorder) UploadFile(ctx any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockGophKeeperV1ServiceClient)(nil).UploadFile), varargs...)
}

// Verify2FA mocks base method.
func (m *MockGophKeeperV1ServiceClient) Verify2FA(ctx context.Context, in *v1.Verify2FARequest, opts ...grpc.CallOption) (*v1.Verify2FAResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Verify2FA", varargs...)
	ret0, _ := ret[0].(*v1.Verify2FAResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify2FA indicates an expected call of Verify2FA.
func (mr *MockGophKeeperV1ServiceClientMockRecorder) Verify2FA(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify2FA", reflect.TypeOf((*MockGophKeeperV1ServiceClient)(nil).Verify2FA), varargs...)
}

// MockGophKeeperV1ServiceServer is a mock of GophKeeperV1ServiceServer interface.
type MockGophKeeperV1ServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockGophKeeperV1ServiceServerMockRecorder
}

// MockGophKeeperV1ServiceServerMockRecorder is the mock recorder for MockGophKeeperV1ServiceServer.
type MockGophKeeperV1ServiceServerMockRecorder struct {
	mock *MockGophKeeperV1ServiceServer
}

// NewMockGophKeeperV1ServiceServer creates a new mock instance.
func NewMockGophKeeperV1ServiceServer(ctrl *gomock.Controller) *MockGophKeeperV1ServiceServer {
	mock := &MockGophKeeperV1ServiceServer{ctrl: ctrl}
	mock.recorder = &MockGophKeeperV1ServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGophKeeperV1ServiceServer) EXPECT() *MockGophKeeperV1ServiceServerMockRecorder {
	return m.recorder
}

// AuthenticateUser mocks base method.
func (m *MockGophKeeperV1ServiceServer) AuthenticateUser(arg0 context.Context, arg1 *v1.AuthenticateUserRequest) (*v1.AuthenticateUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateUser", arg0, arg1)
	ret0, _ := ret[0].(*v1.AuthenticateUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateUser indicates an expected call of AuthenticateUser.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) AuthenticateUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateUser", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).AuthenticateUser), arg0, arg1)
}

// DeleteFile mocks base method.
func (m *MockGophKeeperV1ServiceServer) DeleteFile(arg0 context.Context, arg1 *v1.DeleteFileRequest) (*v1.DeleteFileResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", arg0, arg1)
	ret0, _ := ret[0].(*v1.DeleteFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) DeleteFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).DeleteFile), arg0, arg1)
}

// DownloadFile mocks base method.
func (m *MockGophKeeperV1ServiceServer) DownloadFile(arg0 *v1.DownloadFileRequest, arg1 grpc.ServerStreamingServer[v1.DownloadFileResponse]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadFile indicates an expected call of DownloadFile.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) DownloadFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).DownloadFile), arg0, arg1)
}

// GetMetadataFile mocks base method.
func (m *MockGophKeeperV1ServiceServer) GetMetadataFile(arg0 context.Context, arg1 *v1.GetMetadataFileRequest) (*v1.GetMetadataFileResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadataFile", arg0, arg1)
	ret0, _ := ret[0].(*v1.GetMetadataFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadataFile indicates an expected call of GetMetadataFile.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) GetMetadataFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadataFile", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).GetMetadataFile), arg0, arg1)
}

// GetStoreDataList mocks base method.
func (m *MockGophKeeperV1ServiceServer) GetStoreDataList(arg0 context.Context, arg1 *v1.GetStoreDataListRequest) (*v1.GetStoreDataListResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStoreDataList", arg0, arg1)
	ret0, _ := ret[0].(*v1.GetStoreDataListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStoreDataList indicates an expected call of GetStoreDataList.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) GetStoreDataList(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStoreDataList", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).GetStoreDataList), arg0, arg1)
}

// RegisterUser mocks base method.
func (m *MockGophKeeperV1ServiceServer) RegisterUser(arg0 context.Context, arg1 *v1.RegisterUserRequest) (*v1.RegisterUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0, arg1)
	ret0, _ := ret[0].(*v1.RegisterUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) RegisterUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).RegisterUser), arg0, arg1)
}

// SetMetadataFile mocks base method.
func (m *MockGophKeeperV1ServiceServer) SetMetadataFile(arg0 context.Context, arg1 *v1.SetMetadataFileRequest) (*v1.SetMetadataFileResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMetadataFile", arg0, arg1)
	ret0, _ := ret[0].(*v1.SetMetadataFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetMetadataFile indicates an expected call of SetMetadataFile.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) SetMetadataFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMetadataFile", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).SetMetadataFile), arg0, arg1)
}

// UploadFile mocks base method.
func (m *MockGophKeeperV1ServiceServer) UploadFile(arg0 grpc.ClientStreamingServer[v1.UploadFileRequest, v1.UploadFileResponse]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) UploadFile(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).UploadFile), arg0)
}

// Verify2FA mocks base method.
func (m *MockGophKeeperV1ServiceServer) Verify2FA(arg0 context.Context, arg1 *v1.Verify2FARequest) (*v1.Verify2FAResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify2FA", arg0, arg1)
	ret0, _ := ret[0].(*v1.Verify2FAResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify2FA indicates an expected call of Verify2FA.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) Verify2FA(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify2FA", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).Verify2FA), arg0, arg1)
}

// mustEmbedUnimplementedGophKeeperV1ServiceServer mocks base method.
func (m *MockGophKeeperV1ServiceServer) mustEmbedUnimplementedGophKeeperV1ServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedGophKeeperV1ServiceServer")
}

// mustEmbedUnimplementedGophKeeperV1ServiceServer indicates an expected call of mustEmbedUnimplementedGophKeeperV1ServiceServer.
func (mr *MockGophKeeperV1ServiceServerMockRecorder) mustEmbedUnimplementedGophKeeperV1ServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedGophKeeperV1ServiceServer", reflect.TypeOf((*MockGophKeeperV1ServiceServer)(nil).mustEmbedUnimplementedGophKeeperV1ServiceServer))
}

// MockUnsafeGophKeeperV1ServiceServer is a mock of UnsafeGophKeeperV1ServiceServer interface.
type MockUnsafeGophKeeperV1ServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeGophKeeperV1ServiceServerMockRecorder
}

// MockUnsafeGophKeeperV1ServiceServerMockRecorder is the mock recorder for MockUnsafeGophKeeperV1ServiceServer.
type MockUnsafeGophKeeperV1ServiceServerMockRecorder struct {
	mock *MockUnsafeGophKeeperV1ServiceServer
}

// NewMockUnsafeGophKeeperV1ServiceServer creates a new mock instance.
func NewMockUnsafeGophKeeperV1ServiceServer(ctrl *gomock.Controller) *MockUnsafeGophKeeperV1ServiceServer {
	mock := &MockUnsafeGophKeeperV1ServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeGophKeeperV1ServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeGophKeeperV1ServiceServer) EXPECT() *MockUnsafeGophKeeperV1ServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedGophKeeperV1ServiceServer mocks base method.
func (m *MockUnsafeGophKeeperV1ServiceServer) mustEmbedUnimplementedGophKeeperV1ServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedGophKeeperV1ServiceServer")
}

// mustEmbedUnimplementedGophKeeperV1ServiceServer indicates an expected call of mustEmbedUnimplementedGophKeeperV1ServiceServer.
func (mr *MockUnsafeGophKeeperV1ServiceServerMockRecorder) mustEmbedUnimplementedGophKeeperV1ServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedGophKeeperV1ServiceServer", reflect.TypeOf((*MockUnsafeGophKeeperV1ServiceServer)(nil).mustEmbedUnimplementedGophKeeperV1ServiceServer))
}
