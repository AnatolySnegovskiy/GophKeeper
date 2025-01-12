package ui

import (
	"errors"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goph_keeper/internal/client"
	"log/slog"
	"os"
	"runtime"
	"testing"
)

func TestNewMenu(t *testing.T) {
	app := tview.NewApplication()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	grpcClient := &client.GrpcClient{}

	menu := NewMenu(app, logger, grpcClient)

	require.NotNil(t, menu)
	assert.Equal(t, app, menu.app)
	assert.Equal(t, logger, menu.logger)
	assert.Equal(t, grpcClient, menu.grpcClient)
	assert.Equal(t, "Goph Keeper Client", menu.title)
}

func TestErrorHandler(t *testing.T) {
	app := tview.NewApplication()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	grpcClient := &client.GrpcClient{}
	menu := NewMenu(app, logger, grpcClient)

	err := errors.New("test error")
	callback := func() {
		assert.True(t, true)
	}

	menu.errorHandler(err, callback)
}

func TestIsRootPath(t *testing.T) {
	testCases := []struct {
		os     string
		path   string
		result bool
	}{
		{os: "windows", path: "C:\\", result: true},
		{os: "windows", path: "C:", result: true},
		{os: "windows", path: "C:\\folder", result: false},
		{os: "linux", path: "/", result: true},
		{os: "linux", path: "/home", result: false},
	}

	for _, tc := range testCases {
		currentOS := runtime.GOOS
		if currentOS != tc.os {
			continue
		}
		assert.Equal(t, tc.result, isRootPath(tc.path))
	}
}

func TestIsDriveRoot(t *testing.T) {
	testCases := []struct {
		os     string
		path   string
		result bool
	}{
		{os: "windows", path: "C:\\", result: true},
		{os: "windows", path: "C:", result: false},
		{os: "windows", path: "C:\\folder", result: false},
		{os: "linux", path: "/home", result: false},
	}

	for _, tc := range testCases {
		currentOS := runtime.GOOS
		if currentOS != tc.os {
			continue
		}
		assert.Equal(t, tc.result, isDriveRoot(tc.path), tc.path)
	}
}
