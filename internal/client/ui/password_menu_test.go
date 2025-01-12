package ui

import (
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"testing"
)

func TestShowPasswordMenu(t *testing.T) {
	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: []*v1.ListDataEntry{
			{
				UserPath: "Test",
				Uuid:     "Test",
			},
		},
	}, nil).AnyTimes()
	menu := getMenu(mockClient)
	menu.showPasswordMenu()
	assert.NotNil(t, menu.app)
}

func TestShowPasswordForm(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}
	menu.showPasswordForm(&entities.FilePassword{})
	assert.NotNil(t, menu.app)
}
