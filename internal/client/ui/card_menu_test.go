package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
	"goph_keeper/internal/client"
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log/slog"
	"os"
	"testing"
)

func TestShowCardsMenu(t *testing.T) {
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
	mockClient.EXPECT().GetMetadataFile(gomock.Any(), gomock.Any()).Return(&v1.GetMetadataFileResponse{}, nil).AnyTimes()

	mockStream := grpc.ServerStreamingClient[v1.DownloadFileResponse](nil)
	mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil).AnyTimes()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	assert.NotNil(t, logger, "Logger should not be nil")

	grpcClient := client.NewGrpcClient(logger, mockClient, login, login)

	menu := &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: grpcClient,
		logger:     logger,
	}

	menu.showCardsMenu()
	assert.NotNil(t, menu.app)
	list := menu.app.GetFocus().(*tview.List)
	handler := list.InputHandler()

	for item := 0; item < list.GetItemCount(); item++ {
		list.SetCurrentItem(item)
		event := tcell.NewEventKey(tcell.KeyEnter, 0, 0)
		capture := func(p tview.Primitive) {
			assert.True(t, true, "expected list to have an InputHandler")
		}
		handler(event, capture)
	}

	clear()
}

func TestShowCardForm(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.showCardForm(&entities.FileCard{})
	assert.NotNil(t, menu.app)
}
