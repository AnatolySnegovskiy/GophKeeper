package client_test

import (
	"context"
	"goph_keeper/internal/client"
	"goph_keeper/internal/mocks"
	"goph_keeper/internal/services"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"goph_keeper/internal/testhepler"
	"log/slog"
	"os"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type MockUploadFileClient struct {
	grpc.ClientStream
	sendFunc  func(*v1.UploadFileRequest) error
	recvFunc  func() (*v1.UploadFileResponse, error)
	closeFunc func() error
}

func (m *MockUploadFileClient) Send(req *v1.UploadFileRequest) error {
	return m.sendFunc(req)
}

func (m *MockUploadFileClient) CloseAndRecv() (*v1.UploadFileResponse, error) {
	return m.recvFunc()
}

func (m *MockUploadFileClient) CloseSend() error {
	return m.closeFunc()
}

func TestAuthenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	login := "testuser"
	password := "testpass"

	mockClient := mocks.NewMockGophKeeperV1ServiceClient(ctrl)
	ssh := services.NewSshKeyGen()
	publicKey, _ := ssh.Generate(login)
	randomToken := strconv.Itoa(10)
	tokenByte, _ := ssh.EncryptMessage([]byte(randomToken), publicKey)
	mockClient.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(&v1.AuthenticateUserResponse{
		Success: true,
		Token:   string(tokenByte),
	}, nil).AnyTimes()
	mockClient.EXPECT().Verify2FA(gomock.Any(), gomock.Any()).Return(&v1.Verify2FAResponse{
		Success:  true,
		JwtToken: randomToken,
	}, nil).AnyTimes()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := client.NewGrpcClient(logger, mockClient)
	ctx := context.Background()

	response, err := client.Authenticate(ctx, login, password)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "10", response.JwtToken)
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockGophKeeperV1ServiceClient(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := client.NewGrpcClient(logger, mockClient)

	ctx := context.Background()
	login := "testuser"
	password := "testpass"

	// Mock RegisterUser response
	mockClient.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(&v1.RegisterUserResponse{
		Success: true,
	}, nil)

	response, err := client.RegisterUser(ctx, login, password)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestUploadFile(t *testing.T) {
	mockClient := testhepler.GetMockGRPCClient(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := testhepler.GetGrpcClient(mockClient, logger)

	ctx := context.Background()
	file, _ := os.CreateTemp("./.ssh/", "testfile.txt")
	_, err := file.Write([]byte("Hello, World!"))
	assert.NoError(t, err)
	file.Close()
	filePath := file.Name()
	userPath := "/testuser"
	fileType := v1.DataType_DATA_TYPE_BINARY

	// Mock UploadFile response
	mockStream := &MockUploadFileClient{
		sendFunc: func(req *v1.UploadFileRequest) error {
			return nil
		},
		recvFunc: func() (*v1.UploadFileResponse, error) {
			return &v1.UploadFileResponse{}, nil
		},
		closeFunc: func() error {
			return nil
		},
	} // Создаем канал для прогресса
	mockClient.EXPECT().UploadFile(gomock.Any()).Return(mockStream, nil)
	// Mock SetMetadataFile response
	mockClient.EXPECT().SetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.SetMetadataFileResponse{
		Success: true,
	}, nil)
	done := make(chan struct{})
	progressChan := make(chan int, 1)

	go func() {
		defer close(done)
		response, err := client.UploadFile(ctx, filePath, userPath, fileType, progressChan)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		close(progressChan)
	}()

	for progress := range progressChan {
		if progress == 100 {
			break
		}
	}

	<-done
}

func TestDownloadFile(t *testing.T) {
	mockClient := testhepler.GetMockGRPCClient(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := testhepler.GetGrpcClient(mockClient, logger)

	ctx := context.Background()
	uuid := "file_uuid"
	path := os.TempDir()
	progressChan := make(chan int, 100)

	// Mock getAuthCTX response
	mockClient.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(&v1.AuthenticateUserResponse{
		Success: true,
		Token:   "encrypted_token",
	}, nil)

	mockClient.EXPECT().Verify2FA(gomock.Any(), gomock.Any()).Return(&v1.Verify2FAResponse{
		Success:  true,
		JwtToken: "jwt_token",
	}, nil)

	// Mock GetMetadataFile response
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.GetMetadataFileResponse{
		Metadata: `{"file_name":"eteas707606254","file_extension":"","mem_type":"application/octet-stream","is_compressed":false,"compression_type":"","file_size":82}`,
	}, nil)

	// Mock DownloadFile response
	testFile := testhepler.GetTestGoodFile()
	mockStream := testhepler.GetDownloadStreaming(testFile, v1.Status_STATUS_SUCCESS)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil)
	done := make(chan struct{})
	go func() {
		defer close(done)
		file, err := client.DownloadFile(ctx, uuid, path, progressChan)
		assert.NoError(t, err)
		assert.NotNil(t, file)
		close(progressChan)
		file.Close()
	}()
	for progress := range progressChan {
		if progress == 100 {
			break
		}
	}
	<-done
}

func TestDeleteFile(t *testing.T) {
	mockClient := testhepler.GetMockGRPCClient(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := testhepler.GetGrpcClient(mockClient, logger)

	ctx := context.Background()
	uuid := "file_uuid"

	// Mock getAuthCTX response
	mockClient.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(&v1.AuthenticateUserResponse{
		Success: true,
		Token:   "encrypted_token",
	}, nil)
	mockClient.EXPECT().Verify2FA(gomock.Any(), gomock.Any()).Return(&v1.Verify2FAResponse{
		Success:  true,
		JwtToken: "jwt_token",
	}, nil)

	// Mock DeleteFile response
	mockClient.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).Return(&v1.DeleteFileResponse{
		Success: true,
	}, nil)

	err := client.DeleteFile(ctx, uuid)
	assert.NoError(t, err)
}
