package server

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"goph_keeper/internal/server/services/jwt"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
	"net"
	"regexp"
	"testing"
	"time"
)

func TestGrpcServer_Run(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err, "an error was not expected when opening a stub database connection")

	lis, err := net.Listen("tcp", ":0")
	assert.NoError(t, err, "an error was not expected when listening on a TCP port")

	server := getServer(db)
	assert.NotNil(t, server, "server should not be nil")

	go func() {
		err := server.Run(lis)
		assert.NoError(t, err, "an error was not expected when running the server")
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Test that the server is listening
	cc, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "an error was not expected when dialing the server")
	defer cc.Close()
}

func TestGrpcServer_RegisterUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE username = $1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT $2")).
		WithArgs("test_user", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "ssh_pub_key"}))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password","ssh_pub_key") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "test_user", sqlmock.AnyArg(), "test_ssh_pub_key").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()
	server := getServer(db)
	req := &v1.RegisterUserRequest{
		Username:  "test_user",
		Password:  "test_password",
		SshPubKey: "test_ssh_pub_key",
	}

	ctx := context.Background()
	resp, err := server.RegisterUser(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
}

func getServer(db gorm.ConnPool) *GrpcServer {
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)

	return NewGrpcServer(
		slog.Default(),
		jwt.NewJwt(),
		&redis.Client{},
		gdb,
	)
}
