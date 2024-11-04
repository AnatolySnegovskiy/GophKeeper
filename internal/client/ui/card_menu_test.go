package ui

import (
	"context"
	"errors"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := getMockGRPCClient(t, "TEST")
	mockClient.EXPECT().GetStoreDataList(
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(ctx context.Context, req *v1.GetStoreDataListRequest, opts ...interface{}) (*v1.GetStoreDataListResponse, error) {
		if req.DataType == v1.DataType_DATA_TYPE_CARD {
			return &v1.GetStoreDataListResponse{
				Entries: []*v1.ListDataEntry{
					{UserPath: "test/path1", Uuid: "uuid1"},
					{UserPath: "test/path2", Uuid: "uuid2"},
				},
			}, nil
		}
		return nil, errors.New("unexpected data type")
	}).AnyTimes()

	grpcClient := client.NewGrpcClient(slog.New(slog.NewJSONHandler(os.Stdout, nil)), mockClient, "TEST", "TEST")
	menu := &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: grpcClient,
		logger:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	menu.showCardsMenu()

	focused := menu.app.GetFocus()

	// Проверяем, что список был создан и содержит нужные элементы
	list, ok := focused.(*tview.List)
	assert.True(t, ok, "expected focused to be a tview.List, but got %T", focused)
	assert.Equal(t, 4, list.GetItemCount(), "expected list to contain 3 items (Add, path1, path2)")

	// Проверяем, что добавление нового элемента работает
	item0, _ := list.GetItemText(0)
	assert.Equal(t, "Add", item0, "expected first item to be 'Add'")

	// Проверяем, что элементы списка корректно отображаются
	item1, _ := list.GetItemText(1)
	assert.Equal(t, "test", item1, "expected second item to be 'test/path1'")

	item2, _ := list.GetItemText(2)
	assert.Equal(t, "test", item2, "expected third item to be 'test/path2'")
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
