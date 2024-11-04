package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"goph_keeper/internal/mocks"
	"goph_keeper/internal/services"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"os"
	"strconv"
	"testing"
)

func getMockGRPCClient(t *testing.T, login string) *mocks.MockGophKeeperV1ServiceClient {
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

func clear() {
	os.RemoveAll("./.ssh/")
}

func simulateKeyPress(key tcell.Key, primitive tview.Primitive) {
	handler := primitive.InputHandler()
	event := tcell.NewEventKey(key, 0, 0)
	handler(event, func(p tview.Primitive) {})
}
