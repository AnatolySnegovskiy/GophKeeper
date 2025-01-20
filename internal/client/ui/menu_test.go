package ui

import (
	"errors"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goph_keeper/internal/client"
	"log/slog"
	"os"
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

	currentOS := GetGOOS()
	for _, tc := range testCases {
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

	currentOS := GetGOOS()
	for _, tc := range testCases {
		if currentOS != tc.os {
			continue
		}
		assert.Equal(t, tc.result, isDriveRoot(tc.path), tc.path)
	}
}

func TestIsRootPathGOOS(t *testing.T) {
	// Сохраняем оригинальную функцию getGOOS
	originalGetGOOS := GetGOOS
	defer func() { GetGOOS = originalGetGOOS }()

	// Тест для Windows
	GetGOOS = func() string { return "windows" }
	if !isRootPath("C:\\") {
		t.Errorf("isRootPath(\"C:\\\") = false; want true")
	}
	if isRootPath("C:\\Windows") {
		t.Errorf("isRootPath(\"C:\\Windows\") = true; want false")
	}

	// Тест для Unix
	GetGOOS = func() string { return "linux" }
	if !isRootPath("/") {
		t.Errorf("isRootPath(\"/\") = false; want true")
	}
	if isRootPath("/home") {
		t.Errorf("isRootPath(\"/home\") = true; want false")
	}
}

func TestIsDriveRootGOOS(t *testing.T) {
	// Сохраняем оригинальную функцию getGOOS
	originalGetGOOS := GetGOOS
	defer func() { GetGOOS = originalGetGOOS }()

	// Тест для Windows
	GetGOOS = func() string { return "windows" }
	if !isDriveRoot("C:\\") {
		t.Errorf("isDriveRoot(\"C:\\\") = false; want true")
	}
	if isDriveRoot("C:\\Windows") {
		t.Errorf("isDriveRoot(\"C:\\Windows\") = true; want false")
	}

	// Тест для Unix
	GetGOOS = func() string { return "linux" }
	if isDriveRoot("/") {
		t.Errorf("isDriveRoot(\"/\") = true; want false")
	}
	if isDriveRoot("/home") {
		t.Errorf("isDriveRoot(\"/home\") = true; want false")
	}
}
