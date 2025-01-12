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
			if err != nil {
				m.errorHandler(err, func() {
					m.showAuthorizationForm()
				})
			}

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
	var username, password string
	form := tview.NewForm()
	form.AddInputField("Username", "", 20, nil, func(text string) { username = text }).
		AddPasswordField("Password", "", 20, '*', func(text string) { password = text }).
		AddButton("Login", func() {
			authenticate, err := m.grpcClient.Authenticate(context.Background(), username, password)
			if err != nil {
				m.errorHandler(err, func() {
					m.showAuthorizationForm()
				})
				return
			}

			if !authenticate.Success {
				form.AddTextArea("Error", "Неверный логин или пароль", 20, 4, 255, nil)
				return
			}

			m.showAppMenu()
		}).
		AddButton("Cancel", func() {
			m.ShowMainMenu()
		})

	form.SetBorder(true).SetTitle("Авторизация").SetTitleAlign(tview.AlignLeft)

	m.app.SetRoot(form, true).SetFocus(form)
}
