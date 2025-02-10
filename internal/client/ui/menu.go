package ui

import (
	"goph_keeper/internal/client"
	"log/slog"
	"runtime"

	"github.com/rivo/tview"
)

var GetGOOS = func() string {
	return runtime.GOOS
}

type Menu struct {
	app        *tview.Application
	logger     *slog.Logger
	grpcClient *client.GrpcClient
	title      string
}

func NewMenu(app *tview.Application, logger *slog.Logger, grpcClient *client.GrpcClient) *Menu {
	return &Menu{
		app:        app,
		logger:     logger,
		grpcClient: grpcClient,
		title:      "Goph Keeper Client",
	}
}

func (m *Menu) errorHandler(err error, callback func()) {
	m.logger.Error(err.Error())
	modal := tview.NewModal().AddButtons([]string{"OK"})
	modal.SetText(err.Error())
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		callback()
	})

	m.app.SetRoot(modal, true)
}

func isRootPath(path string) bool {
	if GetGOOS() == "windows" {
		return len(path) <= 3 && path[1] == ':'
	}
	return path == "/"
}

func isDriveRoot(path string) bool {
	if GetGOOS() == "windows" {
		return len(path) == 3 && path[1] == ':'
	}
	return false
}
