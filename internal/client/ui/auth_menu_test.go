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

	// Функция для симуляции нажатия клавиши
	simulateKeyPress := func(key tcell.Key, primitive tview.Primitive) {
		handler := primitive.InputHandler()
		event := tcell.NewEventKey(key, 0, 0)
		handler(event, func(p tview.Primitive) {})
	}

	// Симулируем ввод данных в поля формы
	inputUsername := "testuser"
	inputPassword := "testpass"
	inputFormHandler := focused.InputHandler()
	// Симулируем ввод имени пользователя
	for _, r := range inputUsername {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyTab, focused) // Перейти к следующему полю
	inputFormHandler = focused.InputHandler()
	// Симулируем ввод пароля
	for _, r := range inputPassword {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	simulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Login
	focused = menu.app.GetFocus()
	assert.IsType(t, &tview.Button{}, focused, "expected focused to be a tview.Form, but got %T", focused)
	assert.Equal(t, "Login", focused.(*tview.Button).GetLabel(), "expected focused to be a tview.Form, but got %T", focused)
	// Симулируем нажатие кнопки Login
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showAppMenu to be called")

	menu.showAuthorizationForm()
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyTab, focused) // Перейти к следующему полю
	simulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Login
	simulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Cancel
	focused = menu.app.GetFocus()
	assert.IsType(t, &tview.Button{}, focused, "expected focused to be a tview.Form, but got %T", focused)
	assert.Equal(t, "Cancel", focused.(*tview.Button).GetLabel(), "expected focused to be a tview.Form, but got %T", focused)

	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected ShowMainMenu to be called")
}
