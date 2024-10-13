package ui

import (
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"goph_keeper/internal/client"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log/slog"
	"os"
	"testing"
)

func TestShowServerFilesMenu(t *testing.T) {
	entries := []*v1.ListDataEntry{
		{
			UserPath: "Test",
			Uuid:     "Test",
		},
		{
			UserPath: "Test2/Test3",
			Uuid:     "Test2",
		},
	}

	mockClient := getMockGRPCClient(t, "Test")
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: entries,
	}, nil).AnyTimes()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	grpcClient := client.NewGrpcClient(logger, mockClient, "Test", "Test")

	menu := &Menu{
		app:        tview.NewApplication(),
		title:      "Test Title",
		grpcClient: grpcClient,
		logger:     logger,
	}

	menu.showServerFilesMenu("Test")
	assert.NotNil(t, menu.app)

	clear()
}

func TestBuildVirtualDirectories(t *testing.T) {
	entries := []*v1.ListDataEntry{
		{
			UserPath: "Test",
			Uuid:     "Test",
		},
	}

	vDirs := buildVirtualDirectories(entries)
	assert.Equal(t, 1, len(vDirs))
	assert.Equal(t, "Test", vDirs["Test"].UserPath)
}

func TestShowVirtualDirectoryContents(t *testing.T) {

	entries := []*v1.ListDataEntry{
		{
			UserPath: "Test",
			Uuid:     "Test",
		},
	}

	vDirs := buildVirtualDirectories(entries)
	assert.Equal(t, 1, len(vDirs))
	assert.Equal(t, "Test", vDirs["Test"].UserPath)
	assert.Equal(t, "Test", vDirs["Test"].Uuid)
	assert.Equal(t, "Test", vDirs["Test"].UserPath)
	assert.Equal(t, "Test", vDirs["Test"].Uuid)
	assert.Equal(t, "Test", vDirs["Test"].UserPath)
	assert.Equal(t, "Test", vDirs["Test"].Uuid)
}
