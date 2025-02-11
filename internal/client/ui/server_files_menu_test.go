package ui

import (
	"errors"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"goph_keeper/internal/testhepler"
	"path/filepath"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestShowServerFilesMenu(t *testing.T) {
	entries := []*v1.ListDataEntry{
		{
			UserPath: "Test",
			Uuid:     "Test",
		},
		{
			UserPath: "Test2/Test3",
			Uuid:     "Test2",
		},
	}

	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: entries,
	}, nil).AnyTimes()
	menu := GetMenu(mockClient)

	menu.showServerFilesMenu("Test")
	assert.NotNil(t, menu.app)

	focused := menu.app.GetFocus()
	list, _ := focused.(*tview.List)
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Back", currentItemName)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	list, _ = focused.(*tview.List)
	currentItemName, _ = list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Download", currentItemName)
	testhepler.Clear()
}

func TestErrGetStoreDataListServerFilesMenu(t *testing.T) {
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	menu := GetMenu(mockClient)
	menu.showServerFilesMenu("Test")
	assert.NotNil(t, menu.app)

	focused := menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.Clear()
}

func TestBuildVirtualDirectories(t *testing.T) {
	entries := []*v1.ListDataEntry{
		{
			UserPath: "Test",
			Uuid:     "Test",
		},
	}

	vDirs := buildVirtualDirectories(entries)
	assert.Equal(t, 1, len(vDirs))
	assert.Equal(t, "Test", vDirs["Test"].UserPath)
}

func TestBuildVirtualDirectories2(t *testing.T) {

	entries := []*v1.ListDataEntry{
		{
			UserPath: "Test",
			Uuid:     "Test",
		},
	}

	vDirs := buildVirtualDirectories(entries)
	assert.Equal(t, 1, len(vDirs))
	assert.Equal(t, "Test", vDirs["Test"].UserPath)
	assert.Equal(t, "Test", vDirs["Test"].Uuid)
}

func TestShowVirtualDirectoryContents(t *testing.T) {
	entries := []*v1.ListDataEntry{
		{
			UserPath: "/Test",
			Uuid:     "Test",
		},
		{
			UserPath: "/Test2/Test3",
			Uuid:     "Test2",
		},
	}
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: entries,
	}, nil).AnyTimes()
	menu := GetMenu(mockClient)
	menu.showServerFilesMenu(string(filepath.Separator))

	focused := menu.app.GetFocus()
	list, _ := focused.(*tview.List)

	// Simulate navigating through the menu

	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "..", currentItemName)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)

	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	focused = menu.app.GetFocus()
	list, _ = focused.(*tview.List)

	currentItemName, _ = list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Test2"+string(filepath.Separator), currentItemName)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.SimulateKeyPress(tcell.KeyUp, focused)
	focused = menu.app.GetFocus()
	list, _ = focused.(*tview.List)
	currentItemName, _ = list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Test3", currentItemName)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)

	// Simulate the modal dialog for file download

	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyRight, focused)
	focused = menu.app.GetFocus()
	button, _ := focused.(*tview.Button)
	assert.NotNil(t, button)
	assert.Equal(t, "No", button.GetLabel())

	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)

	// Verify that the menu returns to the previous state
	focused = menu.app.GetFocus()
	list, _ = focused.(*tview.List)
	currentItemName, _ = list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "..", currentItemName)

	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	button, _ = focused.(*tview.Button)
	assert.NotNil(t, button)
	assert.Equal(t, "Yes", button.GetLabel())
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.Clear()
}

func TestShowVirtualDirectoryContentsBack(t *testing.T) {
	entries := []*v1.ListDataEntry{
		{
			UserPath: "/Test",
			Uuid:     "Test",
		},
		{
			UserPath: "/Test2/Test3",
			Uuid:     "Test2",
		},
	}
	mockClient := testhepler.GetMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: entries,
	}, nil).AnyTimes()
	menu := GetMenu(mockClient)
	menu.showServerFilesMenu(string(filepath.Separator))

	focused := menu.app.GetFocus()
	list, _ := focused.(*tview.List)

	// Simulate navigating through the menu
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Back", currentItemName)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)

	focused = menu.app.GetFocus()
	list, _ = focused.(*tview.List)
	currentItemName, _ = list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Download", currentItemName)
	testhepler.Clear()
}
