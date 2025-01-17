package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExplore(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}
	callback := func(filePath string, rollback func()) {
	}
	menu.explore(File, callback)
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	simulateKeyPress(tcell.KeyUp, focused)
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Back", currentItemName)

	simulateKeyPress(tcell.KeyEnter, focused)
}

func TestExploreUnix(t *testing.T) {
	GetGOOS = func() string { return "linux" }
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}
	callback := func(filePath string, rollback func()) {
	}
	menu.explore(File, callback)
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")

	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "/", currentItemName)
	simulateKeyPress(tcell.KeyEnter, focused)
}

func TestExploreDir(t *testing.T) {
	GetGOOS = func() string { return "linux" }
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}
	callback := func(filePath string, rollback func()) {
		assert.True(t, true)
		rollback()
	}
	menu.explore(Dir, callback)
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")

	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "/", currentItemName)
	simulateKeyPress(tcell.KeyEnter, focused)
	simulateKeyPress(tcell.KeyLeft, focused)
	simulateKeyPress(tcell.KeyRight, focused)
	simulateKeyPress(tcell.KeyEnter, focused)
	simulateKeyPress(tcell.KeyEnter, focused)
}

func TestShowDirectoryContents(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}
	callback := func(filePath string, rollback func()) {
	}
	menu.showDirectoryContents(&explore{list: tview.NewList()}, callback)
}
