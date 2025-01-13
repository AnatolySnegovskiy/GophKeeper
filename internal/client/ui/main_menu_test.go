package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowMainMenu(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.ShowMainMenu()
	assert.NotNil(t, menu.app)
}

func TestSelectRegistration(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.ShowMainMenu()
	assert.NotNil(t, menu.app)
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")

	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "1. Регистрация", currentItemName)
	simulateKeyPress(tcell.KeyEnter, focused)
}

func TestSelectAuthorization(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.ShowMainMenu()
	assert.NotNil(t, menu.app)
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")

	simulateKeyPress(tcell.KeyDown, focused)

	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "2. Авторизация", currentItemName)
	simulateKeyPress(tcell.KeyEnter, focused)
}

func TestEscape(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.ShowMainMenu()
	assert.NotNil(t, menu.app)
	focused := menu.app.GetFocus()
	_, ok := focused.(*tview.List)
	assert.True(t, ok, "focused should be of type *tview.List")
	simulateKeyPress(tcell.KeyEscape, focused)
}
