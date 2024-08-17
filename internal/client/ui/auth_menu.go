package ui

import (
	"github.com/rivo/tview"
	"log/slog"
)

func (m *Menu) showRegistrationForm(app *tview.Application, logger *slog.Logger) {
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Register", func() {
			// Здесь нужно добавить логику регистрации
			// После успешной регистрации можно перейти к следующему меню
			m.showAppMenu(app, logger)
		}).
		AddButton("Cancel", func() {
			m.ShowMainMenu(app, logger)
		})

	form.SetBorder(true).SetTitle("Регистрация").SetTitleAlign(tview.AlignLeft)

	app.SetRoot(form, true)
}

func (m *Menu) showAuthorizationForm(app *tview.Application, logger *slog.Logger) {
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			// Здесь нужно добавить логику авторизации
			// После успешной авторизации можно перейти к следующему меню
			m.showAppMenu(app, logger)
		}).
		AddButton("Cancel", func() {
			m.ShowMainMenu(app, logger)
		})

	form.SetBorder(true).SetTitle("Авторизация").SetTitleAlign(tview.AlignLeft)

	app.SetRoot(form, true)
}
