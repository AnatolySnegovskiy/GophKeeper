package ui

import (
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestShowDownloadFileForm(t *testing.T) {
	mockClient := getMockGRPCClient(t)
	menu := getMenu(mockClient)

	menu.showDownloadFileForm(&v1.ListDataEntry{
		UserPath: "Test",
		Uuid:     "Test",
	}, func() {})
	assert.NotNil(t, menu.app)
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
	app := tview.NewApplication()

	// Запускаем приложение
	go func() {
		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	entry := &v1.ListDataEntry{Uuid: "test-uuid"}

	mockClient := getMockGRPCClient(t)
	testFile := getTestGoodFile()
	mockStream := getDownloadStreaming(testFile, v1.Status_STATUS_SUCCESS)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil)
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(
		&v1.GetMetadataFileResponse{
			Metadata: "{\"file_name\":\"eteas707606254\",\"file_extension\":\"\",\"mem_type\":\"application/octet-stream\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":82}",
		},
		nil,
	).AnyTimes()
	grpcClient := getGrpcClient(mockClient, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	done := make(chan struct{})
	progressChan := make(chan int)

	go func() {
		defer close(done)
		handleFileDownload(os.TempDir(), entry, progressChan, info, grpcClient, app)
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
		app.Stop()
		t.Log("Download completed, stopping application")
		assert.Equal(t, "[green]Success: true", info.GetText(false), "Expected success message")
	case <-time.After(30 * time.Second):
		t.Fatal("Test timed out")
	}
	clear()
}

func TestHandleProgressUpdates(t *testing.T) {
	progressBar := NewProgressBar(100)
	form := tview.NewForm()
	app := tview.NewApplication()

	// Запускаем приложение
	go func() {
		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	progressChan := make(chan int)
	rollbackFilesMenu := func() {}

	go handleProgressUpdates(progressChan, progressBar, rollbackFilesMenu, form, app)

	// Simulate progress updates
	progressChan <- 50
	progressChan <- 100
	close(progressChan)

	// Wait for the goroutine to finish
	time.Sleep(100 * time.Millisecond)

	// Останавливаем приложение
	app.Stop()

	assert.Equal(t, 100, progressBar.current, "Expected progress to be 100", "Expected form to have a button with text 'OK'")
}
