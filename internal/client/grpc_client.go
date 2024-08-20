package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"goph_keeper/internal/services"
	"goph_keeper/internal/services/grpc/goph_keeper/v1"
	"io"
	"log/slog"
	"os"
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

func (c *GrpcClient) Authenticate(ctx context.Context, login string, password string) (*v1.Verify2FAResponse, error) {
	tokenSsh, err := c.grpcClient.AuthenticateUser(ctx, &v1.AuthenticateUserRequest{
		Username: login,
		Password: password,
	})

	if err != nil || tokenSsh == nil || !tokenSsh.Success {
		return nil, err
	}

	file, err := os.Open("./.ssh/private_key.pem")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	ssh := services.NewSshKeyGen()
	token, err := ssh.DecryptionFunction(tokenSsh.Token, string(content))
	if err != nil {
		return nil, err
	}

	return c.grpcClient.Verify2FA(ctx, &v1.Verify2FARequest{
		Token: token,
	})
}

func (c *GrpcClient) RegisterUser(ctx context.Context, login string, password string) (*v1.RegisterUserResponse, error) {
	ssh := services.NewSshKeyGen()
	sshPub, err := ssh.Generate()
	if err != nil {
		return nil, err
	}

	return c.grpcClient.RegisterUser(ctx, &v1.RegisterUserRequest{
		Username:  login,
		Password:  password,
		SshPubKey: sshPub,
	})
}

func (c *GrpcClient) UploadFile(ctx context.Context, filePath string, progressChan chan<- int) (*v1.UploadDataResponse, error) {
	authCTX, err := c.getAuthCTX(ctx)
	if err != nil {
		return nil, err
	}

	stream, err := c.grpcClient.UploadData(authCTX)
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
