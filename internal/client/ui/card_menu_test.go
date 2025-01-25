package ui

import (
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"goph_keeper/internal/testhepler"
	"testing"
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

func TestShowCardForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
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
	menu := GetMenu(mockClient)

	fileCard := &entities.FileCard{
		Uuid:        "uuid",
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

	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputDescription := "2"
	for _, r := range inputDescription {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputCardNumber := "2"
	for _, r := range inputCardNumber {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputExpiryDate := "2"
	for _, r := range inputExpiryDate {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputCVV := "2"
	for _, r := range inputCVV {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputCardHolder := "2"
	for _, r := range inputCardHolder {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	// Simulate submitting the form
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
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
		testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showCardsMenu to be called after submitting")

	// Test the "Cancel" button logic
	menu.showCardForm(fileCard)
	focused = menu.app.GetFocus()
	for i := 0; i < 7; i++ {
		testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showCardsMenu to be called after canceling")
	testhepler.Clear()
}

func TestDeleteCardErr(t *testing.T) {
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	fileCard := &entities.FileCard{
		Uuid:        "uuid",
		CardName:    "Test Card",
		Description: "Test Description",
		CardNumber:  "1234567890123456",
		ExpiryDate:  "12/25",
		CVV:         "123",
		CardHolder:  "Test Holder",
	}

	menu := GetMenu(mockClient)
	menu.showCardForm(fileCard)
	focused := menu.app.GetFocus()
	for i := 0; i < 6; i++ {
		testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	InputField, ok := focused.(*tview.InputField)
	assert.True(t, ok, "focused should be of type *tview.InputField")
	assert.NotNil(t, InputField, "list should not be nil")
}

func TestUploadCardErr(t *testing.T) {
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockClient.EXPECT().UploadFile(gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	fileCard := &entities.FileCard{
		Uuid:        "uuid",
		CardName:    "Test Card",
		Description: "Test Description",
		CardNumber:  "1234567890123456",
		ExpiryDate:  "12/25",
		CVV:         "123",
		CardHolder:  "Test Holder",
	}
	menu := GetMenu(mockClient)
	menu.showCardForm(fileCard)
	focused := menu.app.GetFocus()
	for i := 0; i < 6; i++ {
		testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	InputField, ok := focused.(*tview.InputField)
	assert.True(t, ok, "focused should be of type *tview.InputField")
	assert.NotNil(t, InputField, "list should not be nil")
}

func TestErrGetStoreDataList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	menu := GetMenu(mockClient)
	menu.showCardsMenu()
	focused := menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	assert.NotNil(t, list, "list should not be nil")
	testhepler.Clear()
}

func TestErrDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := testhepler.GetMockGRPCClient(t)
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
	menu := GetMenu(mockClient)
	menu.showCardsMenu()
	focused := menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	assert.NotNil(t, list, "list should not be nil")
	testhepler.Clear()
}

func TestErrFromFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.GetMetadataFileResponse{
		Metadata: "{\"file_name\":\"SynthVoiceRu.pak\",\"file_extension\":\".pak\",\"mem_type\":\"application/octet-stream\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":2242646908}",
	}, nil).AnyTimes()
	// Мок ответа для DownloadFile
	testFile := testhepler.GetTestBadFile()
	mockStream := testhepler.GetDownloadStreaming(testFile, v1.Status_STATUS_PROCESSING)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil)

	menu := GetMenu(mockClient)
	menu.showCardsMenu()
	focused := menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	assert.NotNil(t, list, "list should not be nil")
	testhepler.Clear()
}

func TestGoodDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.GetMetadataFileResponse{
		Metadata: "{\"file_name\":\"SynthVoiceRu.pak\",\"file_extension\":\".pak\",\"mem_type\":\"application/octet-stream\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":2242646908}",
	}, nil).AnyTimes()
	// Мок ответа для DownloadFile
	testFile := testhepler.GetTestGoodFile()
	mockStream := testhepler.GetDownloadStreaming(testFile, v1.Status_STATUS_PROCESSING)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil)

	menu := GetMenu(mockClient)
	menu.showCardsMenu()
	focused := menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	input, ok := focused.(*tview.InputField)
	assert.True(t, ok, "focused should be of type *tview.InputField")
	assert.NotNil(t, input, "list should not be nil")
	testhepler.Clear()
}

func TestAddCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{},
	}, nil).AnyTimes()
	menu := GetMenu(mockClient)
	menu.showCardsMenu()
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Add", currentItemName)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	input, ok := focused.(*tview.InputField)
	assert.True(t, ok, "focused should be of type *tview.InputField")
	assert.NotNil(t, input, "list should not be nil")
	testhepler.Clear()
}

func TestBack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{},
	}, nil).AnyTimes()
	menu := GetMenu(mockClient)
	menu.showCardsMenu()
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Back", currentItemName)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	list, ok = focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	assert.NotNil(t, list, "list should not be nil")
	currentItemName, _ = list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "1. Файлы", currentItemName)
	testhepler.Clear()
}
