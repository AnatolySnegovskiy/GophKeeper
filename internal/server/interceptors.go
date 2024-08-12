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

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}
