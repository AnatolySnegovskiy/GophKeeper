package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (m *Menu) showAppMenu() {
	title := tview.NewTextView().
		SetText(m.title).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	list := tview.NewList().
		AddItem("1. Файлы", "", '1', func() {
			m.showFilesMenu()
		}).
		AddItem("2. Пароли", "", '2', func() {
		}).
		AddItem("3. Карты", "", '2', func() {
		}).
		SetSelectedFocusOnly(true)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 3, 1, false).
		AddItem(list, 0, 1, true)

	m.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			m.app.Stop()
			return nil
		}
		return event
	})

	m.app.SetRoot(mainLayout, true).SetFocus(list)
}
