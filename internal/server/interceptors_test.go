package server

import (
	"context"
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"goph_keeper/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestJWTStreamInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockJWT := mocks.NewMockJWTInterface(ctrl)
	DB, mockRedis := redismock.NewClientMock()
	server := &GrpcServer{jwt: mockJWT, redis: DB}

	t.Run("valid token", func(t *testing.T) {
		mockJWT.EXPECT().CheckToken("valid_token").Return(nil)
		mockRedis.ExpectGet("valid_token").SetVal("123")

		ss := &mockServerStream{}
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return nil
		}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{"valid_token"}})
		ss.On("Context").Return(ctx)

		err := server.JWTStreamInterceptor(nil, ss, nil, handler)
		assert.NoError(t, err)
	})

	t.Run("no metadata", func(t *testing.T) {
		ss := &mockServerStream{}
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return nil
		}

		ctx := context.Background()
		ss.On("Context").Return(ctx)

		err := server.JWTStreamInterceptor(nil, ss, nil, handler)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = metadata is not provided", err.Error())
	})

	t.Run("no token", func(t *testing.T) {
		ss := &mockServerStream{}
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return nil
		}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
		ss.On("Context").Return(ctx)

		err := server.JWTStreamInterceptor(nil, ss, nil, handler)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = missing authorization token", err.Error())
	})

	t.Run("invalid token", func(t *testing.T) {
		mockJWT.EXPECT().CheckToken("invalid_token").Return(errors.New("invalid token"))

		ss := &mockServerStream{}
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return nil
		}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{"invalid_token"}})
		ss.On("Context").Return(ctx)

		err := server.JWTStreamInterceptor(nil, ss, nil, handler)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = invalid token", err.Error())
	})

	t.Run("redis error", func(t *testing.T) {
		mockJWT.EXPECT().CheckToken("valid_token").Return(nil)
		mockRedis.ExpectGet("valid_token").SetErr(errors.New("redis error"))

		ss := &mockServerStream{}
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return nil
		}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{"valid_token"}})
		ss.On("Context").Return(ctx)

		err := server.JWTStreamInterceptor(nil, ss, nil, handler)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = invalid token", err.Error())
	})
}

func TestJWTUnaryInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockJWT := mocks.NewMockJWTInterface(ctrl)
	DB, mockRedis := redismock.NewClientMock()
	server := &GrpcServer{jwt: mockJWT, redis: DB}

	t.Run("valid token", func(t *testing.T) {
		mockJWT.EXPECT().CheckToken("valid_token").Return(nil).AnyTimes()
		mockRedis.ExpectGet("valid_token").SetVal("123")

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		info := &grpc.UnaryServerInfo{
			Server:     server,
			FullMethod: "/grpc.goph_keeper.v1.GophKeeperV1Service/AuthenticateUser",
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{"valid_token"}})
		_, err := server.JWTUnaryInterceptor(ctx, nil, info, handler)
		assert.NoError(t, err)
	})

	t.Run("no metadata", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		info := &grpc.UnaryServerInfo{
			Server:     server,
			FullMethod: "/grpc.goph_keeper.v1.GophKeeperV1Service/DeleteFile",
		}
		ctx := context.Background()
		_, err := server.JWTUnaryInterceptor(ctx, nil, info, handler)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = metadata is not provided", err.Error())
	})

	t.Run("no token", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		info := &grpc.UnaryServerInfo{
			Server:     server,
			FullMethod: "/grpc.goph_keeper.v1.GophKeeperV1Service/DeleteFile",
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
		_, err := server.JWTUnaryInterceptor(ctx, nil, info, handler)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = missing authorization token", err.Error())
	})

	t.Run("invalid token", func(t *testing.T) {
		mockJWT.EXPECT().CheckToken("invalid_token").Return(errors.New("invalid token"))

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		info := &grpc.UnaryServerInfo{
			Server:     server,
			FullMethod: "/grpc.goph_keeper.v1.GophKeeperV1Service/DeleteFile",
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{"invalid_token"}})
		_, err := server.JWTUnaryInterceptor(ctx, nil, info, handler)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = invalid token", err.Error())
	})

	t.Run("redis error", func(t *testing.T) {
		DB, mockRedis := redismock.NewClientMock()
		server := &GrpcServer{jwt: mockJWT, redis: DB}
		mockJWT.EXPECT().CheckToken("valid_token").Return(nil).AnyTimes()
		mockRedis.ExpectGet("valid_token").SetErr(errors.New("redis error"))

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		info := &grpc.UnaryServerInfo{
			Server:     server,
			FullMethod: "/grpc.goph_keeper.v1.GophKeeperV1Service/DeleteFile",
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{"valid_token"}})
		_, err := server.JWTUnaryInterceptor(ctx, nil, info, handler)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = invalid token", err.Error())
	})

	t.Run("skip method", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		info := &grpc.UnaryServerInfo{
			Server:     server,
			FullMethod: "/grpc.goph_keeper.v1.GophKeeperV1Service/RegisterUser",
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{"valid_token"}})
		_, err := server.JWTUnaryInterceptor(ctx, nil, info, handler)
		assert.NoError(t, err)
	})
}

// mockServerStream is a mock implementation of grpc.ServerStream
type mockServerStream struct {
	mock.Mock
}

func (m *mockServerStream) SetHeader(metadata.MD) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServerStream) SendHeader(metadata.MD) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServerStream) SetTrailer(metadata.MD) {
	m.Called()
}

func (m *mockServerStream) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *mockServerStream) SendMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *mockServerStream) RecvMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}
