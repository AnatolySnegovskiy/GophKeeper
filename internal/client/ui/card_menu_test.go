package ui

import (
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	grpcClient := client.NewGrpcClient(logger, mockClient, login, login)

	menu := &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: grpcClient,
	}

	menu.showCardsMenu()
	assert.NotNil(t, menu.app)
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
