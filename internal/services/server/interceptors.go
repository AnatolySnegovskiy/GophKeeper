package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func JWTStreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	token := md["authorization"]
	if len(token) == 0 {
		return status.Errorf(codes.Unauthenticated, "missing authorization token")
	}
	// Validate the token here
	// If the token is invalid, return an error
	// if !isValidToken(token[0]) {
	//     return grpc.Errorf(grpc.Code(grpc.Unauthenticated), "invalid token")
	// }

	return handler(srv, ss)
}
