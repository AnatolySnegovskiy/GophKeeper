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

func isRootPath(path string) bool {
	if runtime.GOOS == "windows" {
		return len(path) <= 3 && path[1] == ':'
	}
	return path == "/"
}
