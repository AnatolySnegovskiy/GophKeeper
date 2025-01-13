package ui

import (
	"github.com/stretchr/testify/assert"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"testing"
)

func TestShowDownloadFileForm(t *testing.T) {
	mockClient := getMockGRPCClient(t)
	menu := getMenu(mockClient)

	menu.showDownloadFileForm(&v1.ListDataEntry{
		UserPath: "Test",
		Uuid:     "Test",
	}, func() {})
	assert.NotNil(t, menu.app)
	clear()
}
