package ui

import (
	"github.com/rivo/tview"
)

func (m *Menu) showRegistrationForm() {
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Register", func() {
			// Здесь нужно добавить логику регистрации
			// После успешной регистрации можно перейти к следующему меню
			m.showAppMenu()
		}).
		AddButton("Cancel", func() {
			m.ShowMainMenu()
		})

	form.SetBorder(true).SetTitle("Регистрация").SetTitleAlign(tview.AlignLeft)

	m.app.SetRoot(form, true)
}

func (m *Menu) showAuthorizationForm() {
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			// Здесь нужно добавить логику авторизации
			// После успешной авторизации можно перейти к следующему меню
			m.showAppMenu()
		}).
		AddButton("Cancel", func() {
			m.ShowMainMenu()
		})

	form.SetBorder(true).SetTitle("Авторизация").SetTitleAlign(tview.AlignLeft)

	m.app.SetRoot(form, true)
}
