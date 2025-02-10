package main

import (
	"context"
	"fmt"
	config "goph_keeper/config/server"
	"goph_keeper/internal/server"
	"goph_keeper/internal/server/services/db"
	"goph_keeper/internal/server/services/jwt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var gormPostgres = db.NewGormPostgres
var redisClient = redis.NewClient

type app struct {
	logger  *slog.Logger
	config  config.Config
	redis   *redis.Client
	logFile *os.File
	gorm    *gorm.DB
	server  *server.GrpcServer
}

func (a *app) makeApp() net.Listener {
	a.upLogger()
	a.upConfig()
	a.upRedisClient()
	a.upGorm()

	return a.upListener()
}

func (a *app) handleError(err error) {
	if err != nil {
		a.logger.Error(err.Error())
		a.stop()
		os.Exit(1)
	}
}

func (a *app) upListener() net.Listener {
	a.server = server.NewGrpcServer(a.logger, jwt.NewJwt(), a.redis, a.gorm)
	lis, err := net.Listen("tcp", a.config.Server.Host+":"+a.config.Server.Port)
	a.handleError(err)
	return lis
}

func (a *app) upGorm() {
	var err error
	a.gorm, err = gormPostgres(a.config.DB.DSN)
	a.handleError(err)
}

func (a *app) upConfig() {
	var err error
	a.config, err = config.NewConfig()
	a.handleError(err)
}

func (a *app) upRedisClient() {
	a.redis = redisClient(&redis.Options{
		Addr:     a.config.Redis.Host + ":" + a.config.Redis.Port,
		Password: a.config.Redis.Password,
	})
	a.handleError(a.redis.Ping(context.Background()).Err())
}

func (a *app) upLogger() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("logs/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file")
		os.Exit(1)
	}
	multiWriter := io.MultiWriter(file, os.Stdout)
	a.logger = slog.New(slog.NewJSONHandler(multiWriter, nil))
	a.logFile = file
}

func (a *app) stop() {
	err := a.redis.Close()
	if err != nil {
		a.logger.Error(err.Error())
	}
	err = a.logFile.Close()
	if err != nil {
		a.logger.Error(err.Error())
	}
}

func main() {
	a := &app{}
	lis := a.makeApp()
	go func() { a.handleError(a.server.Run(lis)) }()
	a.logger.Info(fmt.Sprintf("server listening at: %s", lis.Addr()))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit

	a.logger.Info("Received signal, exiting...")
	a.stop()
	os.Exit(0)
}
