package ui

import (
	"github.com/gdamore/tcell/v2"
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
	clear()
}
