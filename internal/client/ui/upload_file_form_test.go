package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"goph_keeper/internal/testhepler"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShowSendFileForm(t *testing.T) {
	mockClient := testhepler.GetMockGRPCClient(t)
	mockStream := &MockUploadFileClient{
		sendFunc: func(req *v1.UploadFileRequest) error {
			return nil
		},
		recvFunc: func() (*v1.UploadFileResponse, error) {
			return &v1.UploadFileResponse{}, nil
		},
		closeFunc: func() error {
			return nil
		},
	}
	mockClient.EXPECT().UploadFile(gomock.Any()).Return(mockStream, nil).AnyTimes()

	menu := GetMenu(mockClient)
	file, _ := os.CreateTemp("", "testfile_*.json")
	defer os.Remove(file.Name())
	menu.showSendFileForm()
	assert.NotNil(t, menu.app)

	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	if !ok {
		t.Fatalf("Expected focused element to be a tview.List, got %T", focused)
	}

	// Разбиваем путь файла на части
	targetPathParts := strings.Split(file.Name(), string(os.PathSeparator))
	// Рекурсивная функция для навигации по списку
	navigateToFile(t, list, targetPathParts)

	// Проверяем, что текущий элемент списка соответствует пути файла
	currentItemName, _ := list.GetItemText(list.GetCurrentItem())
	assert.Equal(t, filepath.Base(file.Name()), currentItemName)

	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)
	focused = menu.app.GetFocus()
	inputFormHandler := focused.(*tview.InputField).InputHandler()
	inputTitle := "test"
	for _, r := range inputTitle {
		inputFormHandler(tcell.NewEventKey(tcell.KeyRune, r, 0), nil)
	}

	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	testhepler.SimulateKeyPress(tcell.KeyTab, focused)
	focused = menu.app.GetFocus()
	testhepler.SimulateKeyPress(tcell.KeyEnter, focused)

	testhepler.Clear()
}

// Рекурсивная функция для навигации по списку
func navigateToFile(t *testing.T, list *tview.List, targetPathParts []string) {
	if len(targetPathParts) == 0 {
		return
	}

	currentPart := targetPathParts[0]
	for i := 0; i < list.GetItemCount(); i++ {
		currentItemName, _ := list.GetItemText(i)
		currentItemName = strings.Trim(currentItemName, "\\/ ")
		println(currentItemName)
		if currentItemName == currentPart {
			// Симулируем нажатие клавиши Enter на текущем элементе
			list.SetCurrentItem(i)
			testhepler.SimulateKeyPress(tcell.KeyEnter, list)
			// Рекурсивно продолжаем навигацию по оставшимся частям пути
			navigateToFile(t, list, targetPathParts[1:])
			return
		}
	}

	t.Fatalf("Path part '%s' not found in the list", currentPart)
}
