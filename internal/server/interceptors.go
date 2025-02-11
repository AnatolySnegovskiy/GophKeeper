package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) JWTStreamInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	token := md["authorization"]
	if len(token) == 0 {
		return status.Errorf(codes.Unauthenticated, "missing authorization token")
	}

	if err := s.jwt.CheckToken(token[0]); err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	userID, err := s.redis.Get(ss.Context(), token[0]).Int()

	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	md.Append("user_id", fmt.Sprintf("%d", userID))
	newCtx := metadata.NewIncomingContext(ss.Context(), md)

	return handler(srv, &serverStream{ss, newCtx})
}

func (s *GrpcServer) JWTUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	skipMethods := map[string]bool{
		"/grpc.goph_keeper.v1.GophKeeperV1Service/RegisterUser":     true,
		"/grpc.goph_keeper.v1.GophKeeperV1Service/AuthenticateUser": true,
		"/grpc.goph_keeper.v1.GophKeeperV1Service/Verify2FA":        true,
	}

	if skipMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	token := md["authorization"]
	if len(token) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
	}

	if err := s.jwt.CheckToken(token[0]); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	userID, err := s.redis.Get(ctx, token[0]).Int()
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	md.Append("user_id", fmt.Sprintf("%d", userID))
	newCtx := metadata.NewIncomingContext(ctx, md)

	return handler(newCtx, req)
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}
