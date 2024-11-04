package ui

import (
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"goph_keeper/internal/client"
	"goph_keeper/internal/mocks"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log/slog"
	"os"
	"testing"
)

func TestShowRegistrationForm(t *testing.T) {
	login := "TEST2"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mocks.NewMockGophKeeperV1ServiceClient(ctrl)
	mockClient.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(&v1.RegisterUserResponse{
		Success: true,
	}, nil).AnyTimes()
	grpcClient := client.NewGrpcClient(slog.New(slog.NewJSONHandler(os.Stdout, nil)), mockClient, login, login)

	menu := &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: grpcClient,
	}

	menu.showRegistrationForm()

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

	simulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Register
	focused = menu.app.GetFocus()
	assert.IsType(t, &tview.Button{}, focused, "expected focused to be a tview.Form, but got %T", focused)
	assert.Equal(t, "Register", focused.(*tview.Button).GetLabel(), "expected focused to be a tview.Form, but got %T", focused)
	// Симулируем нажатие кнопки Register
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showAuthorizationForm to be called")

	menu.showRegistrationForm()
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyTab, focused) // Перейти к следующему полю
	simulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Register
	simulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Cancel
	focused = menu.app.GetFocus()
	assert.IsType(t, &tview.Button{}, focused, "expected focused to be a tview.Form, but got %T", focused)
	assert.Equal(t, "Cancel", focused.(*tview.Button).GetLabel(), "expected focused to be a tview.Form, but got %T", focused)

	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected ShowMainMenu to be called")

	mockClient.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(&v1.RegisterUserResponse{
		Success: false,
	}, errors.New("test error")).AnyTimes()
	grpcClient = client.NewGrpcClient(slog.New(slog.NewJSONHandler(os.Stdout, nil)), mockClient, login, login)

	menu = &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: grpcClient,
	}

	menu.showRegistrationForm()
	focused = menu.app.GetFocus()
	simulateKeyPress(tcell.KeyTab, focused) // Перейти к следующему полю
	simulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Register
	focused = menu.app.GetFocus()
	assert.IsType(t, &tview.Button{}, focused, "expected focused to be a tview.Form, but got %T", focused)
	assert.Equal(t, "Register", focused.(*tview.Button).GetLabel(), "expected focused to be a tview.Form, but got %T", focused)
	// Симулируем нажатие кнопки Register
	simulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showAuthorizationForm to be called")

	clear()
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
