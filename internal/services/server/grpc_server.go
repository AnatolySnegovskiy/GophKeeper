package server

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	pb "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"goph_keeper/internal/services/jwt"
	"gorm.io/gorm"
	"log/slog"
	"net"
)

type GrpcServer struct {
	pb.UnimplementedGophKeeperV1ServiceServer
	logger *slog.Logger
	jwt    *jwt.Jwt
	redis  *redis.Client
	db     *gorm.DB
}

func NewGrpcServer(logger *slog.Logger, jwt *jwt.Jwt, redis *redis.Client, db *gorm.DB) *GrpcServer {
	server := &GrpcServer{}
	server.logger = logger
	server.jwt = jwt
	server.redis = redis
	server.db = db
	return server
}

func (s *GrpcServer) Run(lis net.Listener) error {
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(JWTStreamInterceptor))
	pb.RegisterGophKeeperV1ServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// RegisterUser handles user registration.
func (s *GrpcServer) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	s.logger.Info("RegisterUser")
	return &pb.RegisterUserResponse{}, nil
}

// AuthenticateUser handles user authentication.
func (s *GrpcServer) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	return &pb.AuthenticateUserResponse{}, nil
}

// StorePrivateData handles storing private data.
func (s *GrpcServer) StorePrivateData(grpc.ClientStreamingServer[pb.StorePrivateDataRequest, pb.StorePrivateDataResponse]) error {
	return nil
}

// SyncData handles data synchronization.
func (s *GrpcServer) SyncData(grpc.ClientStreamingServer[pb.SyncDataRequest, pb.SyncDataResponse]) error {
	return nil
}

// RequestPrivateData handles requesting private data.
func (s *GrpcServer) RequestPrivateData(grpc.BidiStreamingServer[pb.RequestPrivateDataRequest, pb.RequestPrivateDataResponse]) error {
	return nil
}

func (s *GrpcServer) Stop() {

}
