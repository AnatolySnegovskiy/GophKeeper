package ui

import (
	"context"
	"github.com/rivo/tview"
)

func (m *Menu) showRegistrationForm() {
	var username, password string

	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, func(text string) {
			username = text
		}).
		AddPasswordField("Password", "", 20, '*', func(text string) {
			password = text
		}).
		AddButton("Register", func() {
			res, err := m.grpcClient.RegisterUser(context.Background(), username, password)
			m.errorHandler(err)

			if res.Success {
				m.showAuthorizationForm()
			}
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
