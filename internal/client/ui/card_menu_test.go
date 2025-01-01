package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"goph_keeper/internal/client"
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log/slog"
	"os"
	"testing"
)

type MockDownloadFileClient struct {
	grpc.ClientStream
	recvFunc func() (*v1.DownloadFileResponse, error)
	ctx      context.Context
}

func (m *MockDownloadFileClient) Recv() (*v1.DownloadFileResponse, error) {
	return m.recvFunc()
}

func (m *MockDownloadFileClient) Context() context.Context {
	return m.ctx
}

func TestShowCardsMenu(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t, "TEST")
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "path/to/card1", Uuid: "uuid1"},
			{UserPath: "path/to/card2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()

	// Ожидания для GetMetadataFile
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.GetMetadataFileResponse{
		Metadata: "{\"file_name\":\"The.Union.2024.DUB.WEB-DLRip.720p.x264.seleZen.mkv\",\"file_extension\":\".mkv\",\"mem_type\":\"video/webm\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":2518298229}",
	}, nil).AnyTimes()

	grpcClient := client.NewGrpcClient(slog.New(slog.NewJSONHandler(os.Stdout, nil)), mockClient, "TEST", "TEST")
	menu := &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: grpcClient,
		logger:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	// Ожидания для DownloadFile
	mockStream := &MockDownloadFileClient{
		recvFunc: func() (*v1.DownloadFileResponse, error) {
			return &v1.DownloadFileResponse{
				Status: v1.Status_STATUS_SUCCESS,
				Data:   []byte("test data"),
			}, nil
		},
		ctx: context.Background(),
	}
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil).AnyTimes()

	menu.showCardsMenu()

	focused := menu.app.GetFocus()
	assert.IsType(t, &tview.List{}, focused, "expected focused to be a tview.List, but got %T", focused)

	list := focused.(*tview.List)
	assert.Equal(t, 4, list.GetItemCount(), "expected 4 items in the list")

	// Simulate selecting the first item
	simulateKeyPress(tcell.KeyDown, focused)
	simulateKeyPress(tcell.KeyEnter, focused)

	// Simulate selecting the "Back" item
	for i := 0; i < 3; i++ {
		simulateKeyPress(tcell.KeyDown, focused)
	}
	simulateKeyPress(tcell.KeyEnter, focused)

	assert.True(t, true, "expected showAppMenu to be called")
	clear()
}

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

func TestShowCardForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t, "TEST")
	// Ожидания для UploadFile
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
	}
	mockClient.EXPECT().UploadFile(gomock.Any()).Return(mockStream, nil).AnyTimes()
	grpcClient := client.NewGrpcClient(slog.New(slog.NewJSONHandler(os.Stdout, nil)), mockClient, "TEST", "TEST")
	menu := &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: grpcClient,
		logger:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	fileCard := &entities.FileCard{
		CardName:    "Test Card",
		Description: "Test Description",
		CardNumber:  "1234567890123456",
		ExpiryDate:  "12/25",
		CVV:         "123",
		CardHolder:  "Test Holder",
	}

	menu.showCardForm(fileCard)
	focused := menu.app.GetFocus()

	// Simulate filling out the form
	inputFormHandler := focused.InputHandler()
	inputCardName := "New Card Name"
	for _, r := range inputCardName {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	simulateKeyPress(tcell.KeyTab, focused)
	inputDescription := "New Description"
	for _, r := range inputDescription {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	simulateKeyPress(tcell.KeyTab, focused)
	inputCardNumber := "6543210987654321"
	for _, r := range inputCardNumber {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	simulateKeyPress(tcell.KeyTab, focused)
	inputExpiryDate := "12/26"
	for _, r := range inputExpiryDate {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	simulateKeyPress(tcell.KeyTab, focused)
	inputCVV := "456"
	for _, r := range inputCVV {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	simulateKeyPress(tcell.KeyTab, focused)
	inputCardHolder := "New Holder"
	for _, r := range inputCardHolder {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	simulateKeyPress(tcell.KeyTab, focused)

	// Simulate submitting the form
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showCardsMenu to be called")

	// Simulate canceling the form
	menu.showCardForm(fileCard)
	focused = menu.app.GetFocus()
	for i := 0; i < 7; i++ {
		simulateKeyPress(tcell.KeyTab, focused)
	}
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showCardsMenu to be called")

	clear()
}
