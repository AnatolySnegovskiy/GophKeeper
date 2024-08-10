package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "goph_keeper/internal/services/grpc/goph_keeper/v1"
)

func main() {
	conn, _ := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	grpcClient := pb.NewGophKeeperV1ServiceClient(conn)
	grpcClient.RegisterUser(context.Background(), &pb.RegisterUserRequest{})
}
