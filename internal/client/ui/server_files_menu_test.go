package ui

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
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

	mockClient := getMockGRPCClient(t)
	mockClient.EXPECT().GetStoreDataList(gomock.Any(), gomock.Any()).Return(&v1.GetStoreDataListResponse{
		Entries: entries,
	}, nil).AnyTimes()
	menu := getMenu(mockClient)

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
