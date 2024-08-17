package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
)

func (m *Menu) ShowMainMenu(app *tview.Application, logger *slog.Logger) error {
	title := tview.NewTextView().
		SetText(m.title).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	list := tview.NewList().
		AddItem("1. Регистрация", "", '1', func() {
			m.showRegistrationForm(app, logger)
		}).
		AddItem("2. Авторизация", "", '2', func() {
			m.showAuthorizationForm(app, logger)
		}).
		SetSelectedFocusOnly(true)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 3, 1, false).
		AddItem(list, 0, 1, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			app.Stop()
			return nil
		}
		return event
	})

	app.SetRoot(mainLayout, true).SetFocus(list)
	return app.Run()
}
