package ui

import (
	"errors"
	"goph_keeper/internal/client"
	"goph_keeper/internal/mocks"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"goph_keeper/internal/testhepler"
	"log/slog"
	"os"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestShowRegistrationForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mocks.NewMockGophKeeperV1ServiceClient(ctrl)
	mockClient.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(&v1.RegisterUserResponse{
		Success: true,
	}, nil).AnyTimes()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	menu := &Menu{
		app:        tview.NewApplication(),
		logger:     logger,
		grpcClient: client.NewGrpcClient(logger, mockClient),
	}

	menu.showRegistrationForm()

	focused := menu.app.GetFocus()

	// Симулируем ввод данных в поля формы
	inputUsername := "testuser"
	inputPassword := "testpass"
	inputFormHandler := focused.InputHandler()
	// Симулируем ввод имени пользователя
	for _, r := range inputUsername {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyTab, focused) // Перейти к следующему полю
	focused = menu.app.GetFocus()
	inputFormHandler = focused.InputHandler()
	// Симулируем ввод пароля
	for _, r := range inputPassword {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	testhepler.SimulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Register
	focused = menu.app.GetFocus()
	assert.IsType(t, &tview.Button{}, focused, "expected focused to be a tview.Form, but got %T", focused)
	assert.Equal(t, "Register", focused.(*tview.Button).GetLabel(), "expected focused to be a tview.Form, but got %T", focused)
	// Симулируем нажатие кнопки Register
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected showAuthorizationForm to be called")

	menu.showRegistrationForm()
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyTab, focused) // Перейти к следующему полю
	testhepler.SimulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Register
	testhepler.SimulateKeyPress(tcell.KeyTab, focused) // Перейти к кнопке Cancel
	focused = menu.app.GetFocus()
	assert.IsType(t, &tview.Button{}, focused, "expected focused to be a tview.Form, but got %T", focused)
	assert.Equal(t, "Cancel", focused.(*tview.Button).GetLabel(), "expected focused to be a tview.Form, but got %T", focused)

	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	assert.True(t, true, "expected ShowMainMenu to be called")
	testhepler.Clear()
}

func TestShowRegistrationFormFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mocks.NewMockGophKeeperV1ServiceClient(ctrl)
	mockClient.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(&v1.RegisterUserResponse{
		Success: false,
	}, errors.New("error")).AnyTimes()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	menu := &Menu{
		app:        tview.NewApplication(),
		logger:     logger,
		grpcClient: client.NewGrpcClient(logger, mockClient),
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
	focused = menu.app.GetFocus()
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

	focused = menu.app.GetFocus()
	assert.Equal(t, "OK", focused.(*tview.Button).GetLabel(), "expected focused to be a tview.Form, but got %T", focused)
	simulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	assert.Equal(t, "Username", focused.(*tview.InputField).GetLabel(), "expected focused to be a tview.Form, but got %T", focused)

	testhepler.Clear()
}

func TestShowAuthorizationForm(t *testing.T) {
	mockClient := testhepler.GetMockGRPCClient(t)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	menu := &Menu{
		app:        tview.NewApplication(),
		logger:     logger,
		grpcClient: client.NewGrpcClient(logger, mockClient),
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
	testhepler.Clear()
}
