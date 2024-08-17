package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log/slog"
)

type GrpcClient struct {
	grpcClient v1.GophKeeperV1ServiceClient
	login      string
	password   string
	logger     *slog.Logger
}

func NewGrpcClient(logger *slog.Logger, conn *grpc.ClientConn, login string, password string) *GrpcClient {
	client := &GrpcClient{
		login:    login,
		password: password,
		logger:   logger,
	}

	client.grpcClient = v1.NewGophKeeperV1ServiceClient(conn)
	logger.Info("Connected to gRPC server")

	return client
}

func (c *GrpcClient) Authenticate(ctx context.Context, login string, password string) (*v1.AuthenticateUserResponse, error) {
	token, err := c.grpcClient.AuthenticateUser(ctx, &v1.AuthenticateUserRequest{
		Username: login,
		Password: password,
	})

	return token, err
}

func (c *GrpcClient) StoreData(ctx context.Context, filePath string, progressChan chan<- int) (*v1.StorePrivateDataResponse, error) {
	authCTX, err := c.getAuthCTX(ctx)
	if err != nil {
		return nil, err
	}

	stream, err := c.grpcClient.StorePrivateData(authCTX)
	if err != nil {
		return nil, err
	}

	err = c.sendFile(stream, filePath, progressChan)
	if err != nil {
		return nil, err
	}

	res, err := stream.CloseAndRecv()
	return res, err
}

func (c *GrpcClient) getAuthCTX(ctx context.Context) (context.Context, error) {
	token, err := c.Authenticate(ctx, c.login, c.password)

	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{"authorization": token.JwtToken})
	authCTX := metadata.NewOutgoingContext(ctx, md)
	return authCTX, nil
}
