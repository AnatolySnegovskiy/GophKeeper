package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	config "goph_keeper/config/server"
	"goph_keeper/internal/server"
	"goph_keeper/internal/server/services/db"
	"goph_keeper/internal/server/services/jwt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func handleError(logger *slog.Logger, err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file")
		os.Exit(1)
	}
	multiWriter := io.MultiWriter(file, os.Stdout)
	logger := slog.New(slog.NewJSONHandler(multiWriter, nil))
	conf, err := config.NewConfig()
	handleError(logger, err)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":" + conf.Redis.Port,
		Password: conf.Redis.Password,
	})
	defer redisClient.Close()
	handleError(logger, redisClient.Ping(context.Background()).Err())

	gorm, err := db.NewGormPostgres(conf.DB.DSN)
	handleError(logger, err)

	serv := server.NewGrpcServer(logger, jwt.NewJwt(), redisClient, gorm)
	lis, err := net.Listen("tcp", conf.Server.Host+":"+conf.Server.Port)
	handleError(logger, err)

	go func() { handleError(logger, serv.Run(lis)) }()

	logger.Info(fmt.Sprintf("server listening at: %s", lis.Addr()))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-quit
	logger.Info("Received signal, exiting...")
	serv.Stop()
	os.Exit(0)
}
