package ui

import (
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log"
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
	clear()
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

	assert.Equal(t, 100, progressBar.current, "Expected progress to be 100")
}
