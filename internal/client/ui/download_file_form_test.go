package ui

import (
	"errors"
	"fmt"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"goph_keeper/internal/testhepler"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestShowDownloadFileForm(t *testing.T) {
	mockClient := testhepler.GetMockGRPCClient(t)
	testFile := testhepler.GetTestGoodFile()
	mockStream := testhepler.GetDownloadStreaming(testFile, v1.Status_STATUS_SUCCESS)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil)
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(
		&v1.GetMetadataFileResponse{
			Metadata: "{\"file_name\":\"SynthVoiceRu.pak\",\"file_extension\":\".pak\",\"mem_type\":\"application/octet-stream\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":2242646908}",
		},
		nil,
	).AnyTimes()

	menu := GetMenu(mockClient)
	menu.showDownloadFileForm(&v1.ListDataEntry{
		UserPath: "",
		Uuid:     "1111",
	}, func() {})
	assert.NotNil(t, menu.app)

	focused := menu.app.GetFocus()
	_, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	button, ok := focused.(*tview.Button)
	assert.True(t, ok, "focused should be of type *tview.Button")
	assert.Equal(t, "Select Directory", button.GetLabel())
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()

	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	fmt.Printf("Focused widget type: %T\n", focused)

	testhepler.Clear()
}

func TestCreateDownloadForm(t *testing.T) {
	info := tview.NewTextView().SetDynamicColors(true).SetText("")
	progressBar := NewProgressBar(100)

	form := createDownloadForm(info, progressBar)

	assert.NotNil(t, form, "Expected form to be created, got nil")
	assert.Equal(t, "Download File", form.GetTitle(), "Expected form title to be 'Download File'")
}

func TestCreateFlexLayout(t *testing.T) {
	form := tview.NewForm()
	flex := createFlexLayout(form)

	assert.NotNil(t, flex, "Expected flex layout to be created, got nil")
}

func TestCleanDirectoryPath(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"/path/to/dir/..", "/path/to/dir"},
		{"/path/to/dir/.", "/path/to/dir"},
		{"/path/to/dir/../", "/path/to/dir"},
		{"/path/to/dir/./", "/path/to/dir"},
		{"/path/to/dir/../..", "/path/to/dir"},
		{"/path/to/dir/./.", "/path/to/dir"},
	}

	for _, tc := range testCases {
		result := cleanDirectoryPath(tc.input)
		assert.Equal(t, tc.expected, result, "Expected %s, got %s", tc.expected, result)
	}
}

func TestHandleFileDownload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	info := tview.NewTextView().SetDynamicColors(true).SetText("")
	mockApp := &MockApplication{
		Application: tview.NewApplication(),
	}

	entry := &v1.ListDataEntry{Uuid: "test-uuid"}

	mockClient := testhepler.GetMockGRPCClient(t)
	testFile := testhepler.GetTestGoodFile()
	mockStream := testhepler.GetDownloadStreaming(testFile, v1.Status_STATUS_SUCCESS)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil)
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(
		&v1.GetMetadataFileResponse{
			Metadata: "{\"file_name\":\"eteas707606254\",\"file_extension\":\"\",\"mem_type\":\"application/octet-stream\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":82}",
		},
		nil,
	).AnyTimes()
	grpcClient := testhepler.GetGrpcClient(mockClient, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	done := make(chan struct{})
	progressChan := make(chan int)

	go func() {
		defer close(done)
		file, _ := os.CreateTemp("", "test")
		handleFileDownload(file.Name(), entry, progressChan, info, grpcClient, mockApp)
	}()

	// Wait for the progress to reach 100%
	progress := 0
	for progress < 100 {
		select {
		case progress = <-progressChan:
			t.Logf("Progress updated to: %d%%", progress)
		case <-time.After(30 * time.Second):
			t.Fatal("Test timed out waiting for progress update")
		}
	}

	select {
	case <-done:
		t.Log("Download completed, stopping application")
		assert.Equal(t, "[green]Success: true", info.GetText(false), "Expected success message")
	case <-time.After(30 * time.Second):
		t.Fatal("Test timed out")
	}
	testhepler.Clear()
}

func TestHandleFileDownloadErrorPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	info := tview.NewTextView().SetDynamicColors(true).SetText("")
	mockApp := &MockApplication{
		Application: tview.NewApplication(),
	}

	entry := &v1.ListDataEntry{Uuid: "test-uuid"}
	mockClient := testhepler.GetMockGRPCClient(t)
	grpcClient := testhepler.GetGrpcClient(mockClient, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	progressChan := make(chan int)
	handleFileDownload(os.TempDir()+"/invalid_path", entry, progressChan, info, grpcClient, mockApp)
	mockApp.QueueUpdateDraw(func() {
		assert.Conditionf(t, func() bool { return strings.Contains(info.GetText(false), "[red]Error") }, "Expected error message")
	})
	testhepler.Clear()
}

func TestHandleFileDownloadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	info := tview.NewTextView().SetDynamicColors(true).SetText("")
	mockApp := &MockApplication{
		Application: tview.NewApplication(),
	}

	entry := &v1.ListDataEntry{Uuid: "test-uuid"}
	mockClient := testhepler.GetMockGRPCClient(t)
	grpcClient := testhepler.GetGrpcClient(mockClient, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(
		nil,
		errors.New("test error"),
	).AnyTimes()
	progressChan := make(chan int)
	handleFileDownload(os.TempDir(), entry, progressChan, info, grpcClient, mockApp)
	mockApp.QueueUpdateDraw(func() {
		assert.Equal(
			t,
			"[red]Error: test error",
			info.GetText(false),
			"Expected error message",
		)
	})
	testhepler.Clear()
}

func TestHandleProgressUpdates(t *testing.T) {
	progressBar := NewProgressBar(100)
	form := tview.NewForm()
	mockApp := &MockApplication{
		Application: tview.NewApplication(),
	}

	progressChan := make(chan int)
	rollbackFilesMenu := func() {}
	done := make(chan struct{})

	go func() {
		defer close(done)
		handleProgressUpdates(progressChan, progressBar, rollbackFilesMenu, form, mockApp)
	}()

	progressChan <- 50
	progressChan <- 100
	close(progressChan)
	<-done

	assert.Equal(t, 100, progressBar.current, "Expected progress to be 100")
	assert.Equal(t, 100, progressBar.current, "Expected progress to be 100", "Expected form to have a button with text 'OK'")
	button := form.GetButton(0)
	assert.NotNil(t, button, "Expected form to have a button with text 'OK'")
	assert.Equal(t, "OK", button.GetLabel(), "Expected form to have a button with text 'OK'")
	form.SetFocus(0)
	testhepler.SimulateKeyPress(tcell.KeyEnter, button)
}

type MockApplication struct {
	*tview.Application
}

func (m *MockApplication) QueueUpdateDraw(t func()) *tview.Application {
	t()
	return m.Application
}

func (m *MockApplication) Draw() *tview.Application {
	return m.Application
}
