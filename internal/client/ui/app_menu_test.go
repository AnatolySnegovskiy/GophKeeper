package ui

import (
	"github.com/gdamore/tcell/v2"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestShowAppMenu(t *testing.T) {
	// Create a test menu
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	// Call the showAppMenu function
	menu.showAppMenu()

	// Check that the app is not nil
	assert.NotNil(t, menu.app)

	list := tview.NewList().
		AddItem("1. Файлы", "", '1', func() {
			menu.showFilesMenu()
		}).
		AddItem("2. Пароли", "", '2', func() {
			menu.showPasswordMenu()
		}).
		AddItem("3. Карты", "", '2', func() {
			menu.showCardsMenu()
		}).
		SetSelectedFocusOnly(true)
	// Check that the focus is set to the specific list created in showAppMenu
	focused := menu.app.GetFocus()
	assert.Equal(t, list, focused, "expected focus to be on the menu.list")

	// Check that the input capture is set
	assert.NotNil(t, menu.app.GetInputCapture())

	// Check that the Esc key stops the app
	event := tcell.NewEventKey(tcell.KeyEsc, 0, 0)
	capturedEvent := menu.app.GetInputCapture()(event)
	assert.Nil(t, capturedEvent, "expected capturedEvent to be nil")
}
