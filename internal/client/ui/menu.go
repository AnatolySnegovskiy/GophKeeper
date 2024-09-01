package ui

import (
	"github.com/rivo/tview"
	"goph_keeper/internal/client"
	"log/slog"
	"runtime"
)

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

func (m *Menu) errorHandler(err error) {
	if err != nil {
		m.logger.Error(err.Error())
		modal := tview.NewModal().AddButtons([]string{"OK"})
		modal.SetText(err.Error())
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			m.ShowMainMenu()
		})

		m.app.SetRoot(modal, true)
	}
}

func isRootPath(path string) bool {
	if runtime.GOOS == "windows" {
		return len(path) <= 3 && path[1] == ':'
	}
	return path == "/"
}

func isDriveRoot(path string) bool {
	if runtime.GOOS == "windows" {
		return len(path) == 3 && path[1] == ':'
	}
	return false
}
