package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"goph_keeper/internal/client"
	"log/slog"
	"os"
	"path/filepath"
)

func handleError(logger *slog.Logger, err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	handleError(logger, err)
	defer conn.Close()

	c := client.NewGrpcClient(logger, conn, "User1", "test")

	cwd, err := os.Getwd()
	handleError(logger, err)

	filePath := filepath.Join(cwd, "cmd/client/storage/Kingdom.of.the.Planet.of.the.Apes.2024.D.WEBRip.1O8Op.mkv")
	res, err := c.StoreData(context.Background(), filePath)
	handleError(logger, err)
	logger.Info(fmt.Sprintf("Success: %t, Message: %s", res.Success, res.Message))
}
