package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"testing"
)

func TestShowFilesMenu(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.showFilesMenu()
	assert.NotNil(t, menu.app)
}

func TestSelectDownloadFile(t *testing.T) {
	client := getMockGRPCClient(t)
	client.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{UserPath: "file1", Uuid: "uuid1"},
			{UserPath: "file2", Uuid: "uuid2"},
		},
	}, nil).AnyTimes()
	menu := getMenu(client)
	menu.showFilesMenu()

	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Download", currentItemName)
	simulateKeyPress(tcell.KeyEnter, focused)
	clear()
}

func TestSelectUploadFile(t *testing.T) {
	client := getMockGRPCClient(t)
	menu := getMenu(client)
	menu.showFilesMenu()

	focused := menu.app.GetFocus()
	simulateKeyPress(tcell.KeyDown, focused)
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Upload", currentItemName)
	simulateKeyPress(tcell.KeyEnter, focused)
	clear()
}

func TestSelectBackFile(t *testing.T) {
	client := getMockGRPCClient(t)
	menu := getMenu(client)
	menu.showFilesMenu()

	focused := menu.app.GetFocus()
	simulateKeyPress(tcell.KeyDown, focused)
	simulateKeyPress(tcell.KeyDown, focused)
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Back", currentItemName)
	simulateKeyPress(tcell.KeyEnter, focused)
	clear()
}
