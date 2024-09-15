package ui

import "github.com/rivo/tview"

func (m *Menu) showPasswordForm() {
	var title, password string

	form := tview.NewForm().
		AddInputField("Title", "", 20, nil, func(text string) {
			title = text
		}).
		AddPasswordField("Password", "", 20, '*', func(text string) {
			password = text
		}).
		AddButton("Submit", func() {
			// Add your submission logic here
			m.app.QueueUpdateDraw(func() {
				// Update UI or show a message
			})
		}).
		AddButton("Cancel", func() {
			m.showAppMenu()
		})

	form.SetBorder(true).SetTitle("Пароли").SetTitleAlign(tview.AlignLeft)

	m.app.SetRoot(form, true)
}
