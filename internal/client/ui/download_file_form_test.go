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
		tempDir := os.TempDir()
		slog.Info("Using temporary directory", "tempDir", tempDir)
		handleFileDownload(tempDir, entry, progressChan, info, grpcClient, app)
	}()

	go func() {
		<-progressChan
	}()

	select {
	case <-done:
		app.Stop()
		assert.Equal(t, "[green]Success: true", info.GetText(false), "Expected success message")
	case <-time.After(30 * time.Second):
		t.Fatal("Test timed out")
	}
	clear()
}
