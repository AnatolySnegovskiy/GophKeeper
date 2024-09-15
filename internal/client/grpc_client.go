package client

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	entities2 "goph_keeper/internal/server/services/entities"
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
		DataType:   v1.DataType_DATA_TYPE_BINARY,
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

func (c *GrpcClient) DownloadFile(ctx context.Context, uuid string, path string, progressChan chan<- int) error {
	authCTX, err := c.getAuthCTX(ctx)
	if err != nil {
		return err
	}
	metadataResponse, err := c.grpcClient.GetMetadataFile(authCTX, &v1.GetMetadataFileRequest{
		Uuid: uuid,
	})

	if err != nil {
		return err
	}

	stream, err := c.grpcClient.DownloadFile(authCTX, &v1.DownloadFileRequest{
		Uuid: uuid,
	})

	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	metadataStruct := entities2.FileMetadata{}
	err = json.Unmarshal([]byte(metadataResponse.Metadata), &metadataStruct)

	file, err := os.Create(path + "/" + metadataStruct.FileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	defer file.Close()
	downloadedBytes := 0
	for {
		resp, err := stream.Recv()
		if resp == nil {
			return fmt.Errorf("failed to receive response: %v", err)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to receive response: %v", err)
		}

		if resp.Status == v1.Status_STATUS_PROCESSING {
			writingBytes, err := file.Write(resp.Data)
			if err != nil {
				return fmt.Errorf("failed to write data to file: %v", err)
			}

			downloadedBytes += writingBytes
			progressChan <- int(float32(downloadedBytes) / float32(metadataStruct.FileSize) * 100)
		}

		if resp.Status == v1.Status_STATUS_SUCCESS {
			progressChan <- 100
			return nil
		}

		if resp.Status == v1.Status_STATUS_FAIL {
			progressChan <- 0
			return fmt.Errorf("failed to receive response: %v", err)
		}

		if resp.Status == v1.Status_STATUS_CANCELLED {
			progressChan <- 0
			return fmt.Errorf("cancelled file download")
		}
	}

	return nil
}
