package main

import (
	"fmt"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"goph_keeper/internal/client"
	"goph_keeper/internal/client/ui"
	"log/slog"
	"os"
)

func handleError(logger *slog.Logger, err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	file, err := os.OpenFile("log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file")
		os.Exit(1)
	}
	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, nil))
	conn, _ := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	c := client.NewGrpcClient(logger, conn, "Test", "Test")

	app := tview.NewApplication()
	menu := ui.NewMenu(app, logger, c)
	menu.ShowMainMenu()
	handleError(logger, app.Run())
}
