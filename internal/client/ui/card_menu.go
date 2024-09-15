package ui

import "github.com/rivo/tview"

func (m *Menu) showCardsMenu() {
	title := tview.NewTextView().
		SetText("Cards Menu").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	form := tview.NewForm().
		AddInputField("Card Number", "", 20, nil, nil).
		AddInputField("Expiry Date", "", 10, nil, nil).
		AddInputField("CVV", "", 3, nil, nil).
		AddButton("Save", func() {
			// Сохранение данных карты
		}).
		AddButton("Cancel", func() {
			m.showAppMenu()
		})

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 3, 1, false).
		AddItem(form, 0, 1, true)

	m.app.SetRoot(mainLayout, true).SetFocus(form)
}
