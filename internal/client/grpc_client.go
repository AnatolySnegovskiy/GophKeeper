package client

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/chacha20"
	"google.golang.org/grpc/metadata"
	"goph_keeper/internal/services"
	"goph_keeper/internal/services/entities"
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

func NewGrpcClient(logger *slog.Logger, grpcClient v1.GophKeeperV1ServiceClient, login string, password string) *GrpcClient {
	client := &GrpcClient{
		login:     login,
		password:  password,
		logger:    logger,
		sizeChunk: 1024 * 1024,
	}

	client.grpcClient = grpcClient
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
	content, err := GetPrivateKey(login)
	if err != nil {
		return nil, err
	}

	ssh := services.NewSshKeyGen()
	token, err := ssh.DecryptionFunction([]byte(tokenSsh.Token), content)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.Verify2FA(ctx, &v1.Verify2FARequest{
		Token: string(token),
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

func (c *GrpcClient) UploadFile(ctx context.Context, filePath string, userPath string, fileType v1.DataType, progressChan chan<- int) (*v1.SetMetadataFileResponse, error) {
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
		DataType:   fileType,
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

func (c *GrpcClient) GetStoreDataList(ctx context.Context, fileType v1.DataType) (*v1.GetStoreDataListResponse, error) {
	authCTX, err := c.getAuthCTX(ctx)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.GetStoreDataList(authCTX, &v1.GetStoreDataListRequest{
		DataType: fileType,
	})
}

func (c *GrpcClient) DownloadFile(ctx context.Context, uuid string, path string, progressChan chan<- int) (*os.File, error) {
	authCTX, err := c.getAuthCTX(ctx)
	if err != nil {
		return nil, err
	}
	metadataResponse, err := c.grpcClient.GetMetadataFile(authCTX, &v1.GetMetadataFileRequest{
		Uuid: uuid,
	})

	if err != nil {
		return nil, err
	}

	stream, err := c.grpcClient.DownloadFile(authCTX, &v1.DownloadFileRequest{
		Uuid: uuid,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	metadataStruct := entities.FileMetadata{}
	err = json.Unmarshal([]byte(metadataResponse.Metadata), &metadataStruct)

	file, err := os.Create(path + "/" + metadataStruct.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}

	defer func(file *os.File, offset int64, whence int) {
		_, err := file.Seek(offset, whence)
		if err != nil {
			c.logger.Error(fmt.Sprintf("failed to seek file: %v", err))
			return
		}
	}(file, 0, io.SeekStart)
	downloadedBytes := 0

	// Получение nonce
	resp, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("failed to receive nonce: %v", err)
	}
	nonce := resp.Data

	// Генерация ключа для ChaCha20
	key := []byte("example key 1234example key 1234") // 32 байта

	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if resp == nil {
			return nil, fmt.Errorf("failed to receive response: %v", err)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive response: %v", err)
		}

		if resp.Status == v1.Status_STATUS_PROCESSING {
			encryptedChunk := resp.Data
			decryptedChunk := make([]byte, len(encryptedChunk))
			cipher.XORKeyStream(decryptedChunk, encryptedChunk)

			writingBytes, err := file.Write(decryptedChunk)
			if err != nil {
				return nil, fmt.Errorf("failed to write data to file: %v", err)
			}

			downloadedBytes += writingBytes

			if progressChan != nil {
				progressChan <- int(float32(downloadedBytes) / float32(metadataStruct.FileSize) * 100)
			}
		}

		if resp.Status == v1.Status_STATUS_SUCCESS {
			if progressChan != nil {
				progressChan <- 100
			}
			return file, nil
		}

		if resp.Status == v1.Status_STATUS_FAIL {
			if progressChan != nil {
				progressChan <- 0
			}
			return nil, fmt.Errorf("failed to receive response: %v", err)
		}

		if resp.Status == v1.Status_STATUS_CANCELLED {
			if progressChan != nil {
				progressChan <- 0
			}
			return nil, fmt.Errorf("cancelled file download")
		}
	}

	return file, nil
}

func (c *GrpcClient) DeleteFile(background context.Context, uuid string) error {

	authCTX, err := c.getAuthCTX(background)

	if err != nil {
		return err
	}

	_, err = c.grpcClient.DeleteFile(authCTX, &v1.DeleteFileRequest{
		Uuid: uuid,
	})

	return err
}
