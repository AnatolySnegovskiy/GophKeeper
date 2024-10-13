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

	// Check that the focus is set to the list
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	testTitleList := [3]string{
		"1. Файлы",
		"2. Пароли",
		"3. Карты",
	}

	for item := 0; item < list.GetItemCount(); item++ {
		mainText, _ := list.GetItemText(item)
		assert.Equal(t, mainText, testTitleList[item], focused)
		list.SetCurrentItem(item)
	}

	assert.True(t, ok, "expected focus to be on a tview.List, but got %T", focused)

	// Check that the input capture is set
	assert.NotNil(t, menu.app.GetInputCapture())

	// Check that the Esc key stops the app
	event := tcell.NewEventKey(tcell.KeyEsc, 0, 0)
	capturedEvent := menu.app.GetInputCapture()(event)
	assert.Nil(t, capturedEvent, "expected capturedEvent to be nil")
}
