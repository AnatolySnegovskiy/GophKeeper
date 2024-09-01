package client

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"goph_keeper/internal/services"
	"goph_keeper/internal/services/file_helper"
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
	sizeChunk  int32
}

func NewGrpcClient(logger *slog.Logger, conn *grpc.ClientConn, login string, password string) *GrpcClient {
	client := &GrpcClient{
		login:     login,
		password:  password,
		logger:    logger,
		sizeChunk: 1024 * 1024,
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

	file, err := os.Open("./.ssh/" + login + "/private_key.pem")
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
	sshPub, err := ssh.Generate(login)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.RegisterUser(ctx, &v1.RegisterUserRequest{
		Username:  login,
		Password:  password,
		SshPubKey: sshPub,
	})
}

func (c *GrpcClient) UploadFile(ctx context.Context, filePath string, userPath string, progressChan chan<- int) (*v1.SetMetadataFileResponse, error) {
	authCTX, err := c.getAuthCTX(ctx)
	if err != nil {
		return nil, err
	}

	stream, err := c.grpcClient.UploadFile(authCTX)
	if err != nil {
		return nil, err
	}

	fileMetadata, err := file_helper.GetFileMetadata(filePath)
	if err != nil {
		return nil, err
	}

	err = c.sendFile(stream, fileMetadata, filePath, progressChan)
	if err != nil {
		return nil, err
	}

	responceFileSender, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	metadataJson, err := json.Marshal(fileMetadata)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.SetMetadataFile(authCTX, &v1.SetMetadataFileRequest{
		Uuid:       responceFileSender.Uuid,
		UserPath:   userPath,
		SizeChunks: c.sizeChunk,
		Metadata:   string(metadataJson),
	})
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

func (c *GrpcClient) GetStoreDataList(ctx context.Context) (*v1.GetStoreDataListResponse, error) {
	authCTX, err := c.getAuthCTX(ctx)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.GetStoreDataList(authCTX, &v1.GetStoreDataListRequest{
		DataType: v1.DataType_DATA_TYPE_BINARY,
	})
}
