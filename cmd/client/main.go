package main

import (
	"fmt"
	config "goph_keeper/config/client"
	"goph_keeper/internal/client"
	"goph_keeper/internal/client/ui"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (a *app) handleError(err error) {
	if err != nil {
		a.logger.Error(err.Error())
		a.stop()
		os.Exit(1)
	}
}

type app struct {
	logger         *slog.Logger
	logFile        *os.File
	config         config.Config
	grpcClientConn *grpc.ClientConn
	tviewApp       *tview.Application
}

func (a *app) upLogger() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("logs/client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file")
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(file, nil))
	logger.Info("Log file opened")

	a.logger = logger
	a.logFile = file
}

func (a *app) upConfig() {
	var err error
	a.config, err = config.NewConfig()
	a.handleError(err)
}

func (a *app) upClientGrpc() *client.GrpcClient {
	a.grpcClientConn, _ = grpc.NewClient(a.config.Server.Host+":"+a.config.Server.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return client.NewGrpcClient(a.logger, v1.NewGophKeeperV1ServiceClient(a.grpcClientConn))
}

func (a *app) stop() {
	err := a.grpcClientConn.Close()
	if err != nil {
		a.logger.Error(err.Error())
	}
	err = a.logFile.Close()
	if err != nil {
		a.logger.Error(err.Error())
	}
}

func (a *app) makeApp() {
	a.upLogger()
	a.upConfig()
	a.tviewApp = tview.NewApplication()
	menu := ui.NewMenu(a.tviewApp, a.logger, a.upClientGrpc())
	menu.ShowMainMenu()
}

func main() {
	a := &app{}
	a.makeApp()
	a.handleError(a.tviewApp.Run())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit

	a.logger.Info("Received signal, exiting...")
	a.stop()
	os.Exit(0)
}
