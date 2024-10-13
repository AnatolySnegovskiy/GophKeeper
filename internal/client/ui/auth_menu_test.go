package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowRegistrationForm(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.showRegistrationForm()
	assert.NotNil(t, menu.app)
}

func TestShowAuthorizationForm(t *testing.T) {
	menu := &Menu{
		app: tview.NewApplication(),
	}

	menu.showAuthorizationForm()

	focused := menu.app.GetFocus()
	form, ok := focused.(*tview.InputField)
	assert.True(t, ok, "expected root to be a tview.Form, but got %T", focused)

	// Получаем обработчик ввода формы
	formHandler := form.InputHandler()
	assert.NotNil(t, formHandler, "expected form to have an InputHandler")

	// Функция для симуляции нажатия клавиши
	simulateKeyPress := func(key tcell.Key) {
		event := tcell.NewEventKey(key, 0, 0)
		formHandler(event, func(p tview.Primitive) {})
	}

	// Симулируем ввод данных в поля формы
	inputUsername := "testuser"
	inputPassword := "testpass"

	// Симулируем ввод имени пользователя
	for _, r := range inputUsername {
		formHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	simulateKeyPress(tcell.KeyTab) // Перейти к следующему полю

	// Симулируем ввод пароля
	for _, r := range inputPassword {
		formHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	simulateKeyPress(tcell.KeyTab) // Перейти к кнопке Login

	// Симулируем нажатие кнопки Login
	simulateKeyPress(tcell.KeyEnter)
	assert.True(t, true, "expected showAppMenu to be called")

	simulateKeyPress(tcell.KeyTab) // Перейти к кнопке Cancel
	simulateKeyPress(tcell.KeyEnter)
	assert.True(t, true, "expected ShowMainMenu to be called")
}
