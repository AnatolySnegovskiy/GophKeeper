package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"goph_keeper/internal/client"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log/slog"
	"os"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestShowAppMenu(t *testing.T) {
	login := "TEST"
	mockClient := getMockGRPCClient(t, login)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{
				UserPath: "Test",
				Uuid:     "Test",
			},
		},
	}, nil).AnyTimes()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	grpcClient := client.NewGrpcClient(logger, mockClient, login, login)

	// Create a test menu
	menu := &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: grpcClient,
	}

	// Call the showAppMenu function
	menu.showAppMenu()

	// Check that the app is not nil
	assert.NotNil(t, menu.app)

	// Check that the focus is set to the list
	focused := menu.app.GetFocus()
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "expected focus to be on a tview.List, but got %T", focused)

	for item := 0; item < list.GetItemCount(); item++ {
		list.SetCurrentItem(item)
		handler := list.InputHandler()
		assert.NotNil(t, handler, "expected list to have an InputHandler")
		event := tcell.NewEventKey(tcell.KeyEnter, 0, 0)
		capture := func(p tview.Primitive) {}
		handler(event, capture)
	}

	// Check that the input capture is set
	assert.NotNil(t, menu.app.GetInputCapture())

	// Check that the Esc key stops the app
	event := tcell.NewEventKey(tcell.KeyEsc, 0, 0)
	capturedEvent := menu.app.GetInputCapture()(event)
	assert.Nil(t, capturedEvent, "expected capturedEvent to be nil")

	event = tcell.NewEventKey(tcell.KeyTab, 0, 0)
	capturedEvent = menu.app.GetInputCapture()(event)
	assert.NotNil(t, capturedEvent, "expected capturedEvent to be nil")

	clear()
}
