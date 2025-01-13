package ui

import (
	"context"
	"encoding/hex"
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"goph_keeper/internal/client"
	"goph_keeper/internal/mocks"
	"goph_keeper/internal/services"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"io"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

// MockDownloadFileClient - мок для интерфейса grpc.ServerStreamingClient[v1.DownloadFileResponse]
// MockDownloadFileClient - мок для интерфейса grpc.ServerStreamingClient[v1.DownloadFileResponse]
type MockDownloadFileClient struct {
	grpc.ServerStreamingClient[v1.DownloadFileResponse]
	recvFunc  func() (*v1.DownloadFileResponse, error)
	ctx       context.Context
	callCount int
}

// Recv - мок метода Recv
func (m *MockDownloadFileClient) Recv() (*v1.DownloadFileResponse, error) {
	m.callCount++
	return m.recvFunc()
}

var login = "TEST"

func getMockGRPCClient(t *testing.T) *mocks.MockGophKeeperV1ServiceClient {
	ssh := services.NewSshKeyGen()
	publicKey, _ := ssh.Generate(login)
	randomToken := strconv.Itoa(10)
	tokenByte, _ := ssh.EncryptMessage([]byte(randomToken), publicKey)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mocks.NewMockGophKeeperV1ServiceClient(ctrl)
	mockClient.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(&v1.AuthenticateUserResponse{
		Success: true,
		Token:   string(tokenByte),
	}, nil).AnyTimes()
	mockClient.EXPECT().Verify2FA(gomock.Any(), gomock.Any()).Return(&v1.Verify2FAResponse{
		Success:  true,
		JwtToken: randomToken,
	}, nil).AnyTimes()
	return mockClient
}

func getGrpcClient(mockClient v1.GophKeeperV1ServiceClient, logger *slog.Logger) *client.GrpcClient {
	grpcClient := client.NewGrpcClient(logger, mockClient)
	_, _ = grpcClient.Authenticate(context.Background(), login, "123")

	return grpcClient
}

func getMenu(mockClient v1.GophKeeperV1ServiceClient) *Menu {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: getGrpcClient(mockClient, logger),
		logger:     logger,
	}
}

func clear() {
	os.RemoveAll("./.ssh/")
}

func simulateKeyPress(key tcell.Key, primitive tview.Primitive) {
	handler := primitive.InputHandler()
	event := tcell.NewEventKey(key, 0, 0)
	handler(event, func(p tview.Primitive) {})
}

func getDownloadStreaming(content []byte, status v1.Status) *MockDownloadFileClient {
	chunkSize := 1024 * 1024 // Размер чанка
	chunks := make([][]byte, 0, len(content)/chunkSize+1)
	for i := 12; i < len(content); i += chunkSize {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}
		chunks = append(chunks, content[i:end])
	}

	callCount := 0
	return &MockDownloadFileClient{
		recvFunc: func() (*v1.DownloadFileResponse, error) {
			callCount++
			if callCount == 1 {
				// Первый вызов возвращает nonce (первые 12 байт из файла)
				return &v1.DownloadFileResponse{
					Status: status,
					Data:   content[:12],
				}, nil
			}
			if callCount-2 < len(chunks) {
				// Возвращаем чанки данных
				return &v1.DownloadFileResponse{
					Status: v1.Status_STATUS_PROCESSING,
					Data:   chunks[callCount-2],
				}, nil
			}
			// Последний вызов возвращает успешный статус
			return &v1.DownloadFileResponse{
				Status: v1.Status_STATUS_SUCCESS,
			}, io.EOF
		},
		ctx: context.Background(),
	}
}

func getTestGoodFile() []byte {
	hexString := "73c069125bf16669dd6be8ca8b87d24c4acb4b396467cb3e06a98e81ed5c0924f624d01a330d9cde23b64abecc368540529760a09f674295d49ab4128b97ea7f8446afbe6b9e96ab23c702cf84e124e42bad6e003dae7f6d939dd8407584"
	byteSlice, _ := hex.DecodeString(hexString)
	return byteSlice
}

func getTestBadFile() []byte {
	hexString := "73c069125bf16669dd6be8ca8b87d24c4acb4b396467cb3e06a98e81ed5c0924f624d01a330d9cde23b64abecc368540529760a09f674295d49ab4128b97ea7f8446afbe6b9e96ab23c702cf84e124e42bad6e003dae7f6d939dd84075"
	byteSlice, _ := hex.DecodeString(hexString)
	return byteSlice
}
