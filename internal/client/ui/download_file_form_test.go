package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"testing"
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

func TestShowDownloadFileFormErr(t *testing.T) {
	mockClient := getMockGRPCClient(t)
	testFile := getTestBadFile()
	mockStream := getDownloadStreaming(testFile, v1.Status_STATUS_PROCESSING)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil)
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(
		&v1.GetMetadataFileResponse{
			Metadata: "{\"file_name\":\"SynthVoiceRu.pak\",\"file_extension\":\".pak\",\"mem_type\":\"application/octet-stream\",\"is_compressed\":false,\"compression_type\":\"\",\"file_size\":2242646908}",
		},
		nil,
	).AnyTimes()
	menu := getMenu(mockClient)

	menu.showDownloadFileForm(&v1.ListDataEntry{
		UserPath: "Test",
		Uuid:     "Test",
	}, func() {})
	assert.NotNil(t, menu.app)

	focused := menu.app.GetFocus()
	_, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	simulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	button, ok := focused.(*tview.Button)
	assert.True(t, ok, "focused should be of type *tview.Button")
	assert.Equal(t, "Select Directory", button.GetLabel())
	simulateKeyPress(tcell.KeyEnter, focused)

	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	fmt.Printf("Focused widget type: %T\n", focused)
	_, ok = focused.(*ProgressBar)
	assert.True(t, ok, "focused should be of type *ProgressBar")
	clear()
}
