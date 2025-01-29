package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"goph_keeper/internal/testhepler"
	"log/slog"
	"os"
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
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
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

	testhepler.SimulateKeyPress(tcell.KeyDown, focused)

	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, "2. Авторизация", currentItemName)
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
}

func TestEscape(t *testing.T) {
	menu := &Menu{
		app:    tview.NewApplication(),
		title:  "Test Title",
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
	done := make(chan struct{})
	go func() {
		menu.ShowMainMenu()
		_ = menu.app.Run()
		close(done)
	}()
	menu.app.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
	menu.app.QueueEvent(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
	<-done
	assert.True(t, true, "Application should have stopped")
}
