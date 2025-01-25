package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"goph_keeper/internal/testhepler"
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
	testhepler.SimulateKeyPress(tcell.KeyUp, focused)
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "Back", currentItemName)

	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
}

func TestExploreUnix(t *testing.T) {
	temp := GetGOOS
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
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	GetGOOS = temp
}

func TestExploreDir(t *testing.T) {
	temp := GetGOOS
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
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.SimulateKeyPress(tcell.KeyDown, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.SimulateKeyPress(tcell.KeyLeft, focused)
	testhepler.SimulateKeyPress(tcell.KeyRight, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	testhepler.SimulateKeyPress(tcell.KeyPgUp, focused)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)

	for i := 0; i < 50; i++ {
		testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	}
	GetGOOS = temp
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
