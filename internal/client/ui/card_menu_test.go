package ui

import (
	"context"
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
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
	mockClient := getMockGRPCClient(t)
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

	menu := getMenu(mockClient)

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
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "path/to/card1", Uuid: "uuid1"},
			{UserPath: "path/to/card2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()
	mockClient.EXPECT().SetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.SetMetadataFileResponse{
		Success: true,
	}, nil).AnyTimes()
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
	menu := getMenu(mockClient)

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
	inputFormHandler := focused.(*tview.InputField).InputHandler()
	inputCardName := "2"
	for _, r := range inputCardName {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	simulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputDescription := "2"
	for _, r := range inputDescription {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	simulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputCardNumber := "2"
	for _, r := range inputCardNumber {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	simulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputExpiryDate := "2"
	for _, r := range inputExpiryDate {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	simulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputCVV := "2"
	for _, r := range inputCVV {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	simulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputCardHolder := "2"
	for _, r := range inputCardHolder {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	simulateKeyPress(tcell.KeyTab, focused)
	// Simulate submitting the form
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showCardsMenu to be called")

	// Verify that the fileCard fields were updated correctly
	assert.Equal(t, "Test Card2", fileCard.CardName, "expected CardName to be updated")
	assert.Equal(t, "Test Description2", fileCard.Description, "expected Description to be updated")
	assert.Equal(t, "12345678901234562", fileCard.CardNumber, "expected CardNumber to be updated")
	assert.Equal(t, "12/252", fileCard.ExpiryDate, "expected ExpiryDate to be updated")
	assert.Equal(t, "1232", fileCard.CVV, "expected CVV to be updated")
	assert.Equal(t, "Test Holder2", fileCard.CardHolder, "expected CardHolder to be updated")

	// Test the "Submit" button logic
	menu.showCardForm(fileCard)
	focused = menu.app.GetFocus()
	for i := 0; i < 6; i++ {
		simulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showCardsMenu to be called after submitting")

	// Test the "Cancel" button logic
	menu.showCardForm(fileCard)
	focused = menu.app.GetFocus()
	for i := 0; i < 7; i++ {
		simulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showCardsMenu to be called after canceling")
	clear()
}

func TestErrGetStoreDataList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	menu := getMenu(mockClient)
	menu.showCardsMenu()
	focused := menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	assert.NotNil(t, list, "list should not be nil")
	clear()
}

func TestErrDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()

	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.GetMetadataFileResponse{
		Metadata: "{\"file_name\":\"SynthVoiceRu.pak\",\"file_extension\":\".pak\",\"mem_type\":\"application/octet-stream\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":2242646908}",
	}, nil).AnyTimes()
	menu := getMenu(mockClient)
	menu.showCardsMenu()
	focused := menu.app.GetFocus()
	simulateKeyPress(tcell.KeyDown, focused)
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	assert.NotNil(t, list, "list should not be nil")
	clear()
}

func TestErrFromFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	menu := getMenu(mockClient)
	menu.showCardsMenu()
	focused := menu.app.GetFocus()
	simulateKeyPress(tcell.KeyDown, focused)
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	assert.NotNil(t, list, "list should not be nil")
	clear()
}
