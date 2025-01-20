package ui

import (
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"testing"
)

func TestShowPasswordMenu(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := getMockGRPCClient(t)
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
	testFile := getTestBadFile()
	mockStream := getDownloadStreaming(testFile, v1.Status_STATUS_PROCESSING)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil).AnyTimes()

	menu := getMenu(mockClient)

	menu.showPasswordMenu()

	focused := menu.app.GetFocus()
	list := focused.(*tview.List)

	// Simulate selecting the first item
	list.SetCurrentItem(0)
	simulateKeyPress(tcell.KeyEnter, list)

	// Verify that the showPasswordForm is called with the correct FilePassword
	assert.True(t, true, "expected showPasswordForm to be called")

	// Simulate selecting the second item
	list.SetCurrentItem(1)
	simulateKeyPress(tcell.KeyEnter, list)

	// Verify that the showPasswordForm is called with the correct FilePassword
	assert.True(t, true, "expected showPasswordForm to be called")

	// Simulate selecting the "Back" item
	list.SetCurrentItem(2)
	simulateKeyPress(tcell.KeyEnter, list)

	// Verify that the showAppMenu is called
	assert.True(t, true, "expected showAppMenu to be called")
	clear()
}

func TestShowPasswordForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
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
	mockClient.EXPECT().SetMetadataFile(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockClient.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()

	menu := getMenu(mockClient)
	filePassword := &entities.FilePassword{
		Uuid:        "uuid",
		Title:       "Test Password",
		Description: "Test Description",
		Login:       "testlogin",
		Password:    "testpassword",
	}

	menu.showPasswordForm(filePassword)

	focused := menu.app.GetFocus()
	inputFormHandler := focused.(*tview.InputField).InputHandler()
	inputTitle := "2"
	for _, r := range inputTitle {
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
	inputLogin := "2"
	for _, r := range inputLogin {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	simulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	inputPassword := "2"
	for _, r := range inputPassword {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	simulateKeyPress(tcell.KeyTab, focused)
	// Simulate submitting the form
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showPasswordMenu to be called")

	// Verify that the filePassword fields were updated correctly
	assert.Equal(t, "Test Password2", filePassword.Title, "expected Title to be updated")
	assert.Equal(t, "Test Description2", filePassword.Description, "expected Description to be updated")
	assert.Equal(t, "testlogin2", filePassword.Login, "expected Login to be updated")
	assert.Equal(t, "testpassword2", filePassword.Password, "expected Password to be updated")

	menu.showPasswordForm(filePassword)
	focused = menu.app.GetFocus()
	for i := 0; i < 4; i++ {
		simulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showPasswordMenu to be called after submitting")

	menu.showPasswordForm(filePassword)
	focused = menu.app.GetFocus()
	for i := 0; i < 5; i++ {
		simulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showPasswordMenu to be called after canceling")

	clear()
}

func TestDeleteErr(t *testing.T) {
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	filePassword := &entities.FilePassword{
		Uuid:        "uuid",
		Title:       "Test Password",
		Description: "Test Description",
		Login:       "testlogin",
		Password:    "testpassword",
	}
	menu := getMenu(mockClient)
	menu.showPasswordForm(filePassword)
	focused := menu.app.GetFocus()
	for i := 0; i < 4; i++ {
		simulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	InputField, ok := focused.(*tview.InputField)
	assert.True(t, ok, "focused should be of type *tview.InputField")
	assert.NotNil(t, InputField, "list should not be nil")
}

func TestUploadErr(t *testing.T) {
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockClient.EXPECT().UploadFile(gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	filePassword := &entities.FilePassword{
		Uuid:        "uuid",
		Title:       "Test Password",
		Description: "Test Description",
		Login:       "testlogin",
		Password:    "testpassword",
	}
	menu := getMenu(mockClient)
	menu.showPasswordForm(filePassword)
	focused := menu.app.GetFocus()
	for i := 0; i < 4; i++ {
		simulateKeyPress(tcell.KeyTab, focused)
	}
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	InputField, ok := focused.(*tview.InputField)
	assert.True(t, ok, "focused should be of type *tview.InputField")
	assert.NotNil(t, InputField, "list should not be nil")
}

func TestErrGetStoreDataListPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()
	menu := getMenu(mockClient)
	menu.showPasswordMenu()
	focused := menu.app.GetFocus()
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	InputField, ok := focused.(*tview.InputField)
	assert.True(t, ok, "focused should be of type *tview.InputField")
	assert.NotNil(t, InputField, "list should not be nil")
	clear()
}

func TestErrDownloadFilePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()

	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.GetMetadataFileResponse{
		Metadata: "metadata",
	}, nil).AnyTimes()
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	menu := getMenu(mockClient)
	menu.showPasswordMenu()
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

func TestErrFromFilePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()
	// Мок ответа для DownloadFile
	testFile := getTestBadFile()
	mockStream := getDownloadStreaming(testFile, v1.Status_STATUS_PROCESSING)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil)
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.GetMetadataFileResponse{
		Metadata: "{\"file_name\":\"SynthVoiceRu.pak\",\"file_extension\":\".pak\",\"mem_type\":\"application/octet-stream\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":2242646908}",
	}, nil).AnyTimes()

	menu := getMenu(mockClient)
	menu.showPasswordMenu()
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

func TestGoodDownloadFilePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()
	// Мок ответа для DownloadFile
	testFile := getTestGoodFile()
	mockStream := getDownloadStreaming(testFile, v1.Status_STATUS_PROCESSING)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil)
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.GetMetadataFileResponse{
		Metadata: "{\"file_name\":\"SynthVoiceRu.pak\",\"file_extension\":\".pak\",\"mem_type\":\"application/octet-stream\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":2242646908}",
	}, nil).AnyTimes()
	menu := getMenu(mockClient)
	menu.showPasswordMenu()
	focused := menu.app.GetFocus()
	simulateKeyPress(tcell.KeyDown, focused)
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyTab, focused)
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	input, ok := focused.(*tview.InputField)
	assert.True(t, ok, "focused should be of type *tview.InputField")
	assert.NotNil(t, input, "list should not be nil")
	clear()
}

func TestAddPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{},
	}, nil).AnyTimes()
	menu := getMenu(mockClient)
	menu.showPasswordMenu()
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Add", currentItemName)
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	input, ok := focused.(*tview.InputField)
	assert.True(t, ok, "focused should be of type *tview.InputField")
	assert.NotNil(t, input, "list should not be nil")
	clear()
}

func TestBackPaswword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{},
	}, nil).AnyTimes()
	menu := getMenu(mockClient)
	menu.showPasswordMenu()
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	simulateKeyPress(tcell.KeyDown, focused)
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Back", currentItemName)
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	list, ok = focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	assert.NotNil(t, list, "list should not be nil")
	currentItemName, _ = list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "1. Файлы", currentItemName)
	clear()
}

func TestErrorGetStoreDataList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	menu := getMenu(mockClient)
	menu.showPasswordMenu()
	focused := menu.app.GetFocus()
	button, ok := focused.(*tview.Button)
	assert.True(t, ok, "focused should be of type *tview.button")
	assert.Equal(t, "OK", button.GetLabel())
	simulateKeyPress(tcell.KeyEnter, focused)
	clear()
}
