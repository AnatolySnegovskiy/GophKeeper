package server

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	passwordhash "github.com/vzglad-smerti/password_hash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"goph_keeper/internal/server/services/jwt"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log/slog"
	"net"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestGrpcServer_Run(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err, "an error was not expected when opening a stub database connection")

	lis, err := net.Listen("tcp", ":0")
	assert.NoError(t, err, "an error was not expected when listening on a TCP port")

	mockRedis, _ := redismock.NewClientMock()
	server := getServer(db, mockRedis)
	assert.NotNil(t, server, "server should not be nil")

	go func() {
		err := server.Run(lis)
		assert.NoError(t, err, "an error was not expected when running the server")
	}()

	// Wait for the server to start
	time.Sleep(1000 * time.Millisecond)

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
	mockRedis, _ := redismock.NewClientMock()
	server := getServer(db, mockRedis)
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

func TestGrpcServer_AuthenticateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	cleanup, _ := setupTestFiles("test_user")

	password, _ := passwordhash.Hash("test_password")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE username = $1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT $2")).
		WithArgs("test_user", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "ssh_pub_key"}).AddRow(1, "test_user", password, "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA3mo7QRfA8cFnKhPfQz2P\nsVFKbjI4F+KQY74skglPN3B3lfE73/k16me46R4EryjTkBo91H0hi0v1rQ6Fuy6A\nG/o9PyNhGSRLWxnwg84ltry3+CVQcrA4UXQBoRTSsv+tjobF5X+QZl3u63ZbVeUH\n12OfOMQhJcwTcJ3TCA2z++XFIFMCgUPa6E3Uy7XxA3Vz2Pk1MXmatjYRJxrdf4U6\nONdS92xbea0E49LS/ckTwDqSeWo/2Jd5KtYBFbiOBNZpsWDA7//mB8GNx1w+UBbo\nLuAJG9k2mATQIirbb1MSqMiWJrQqZIBf3trhgt7Zo3VoYaVvfrvGBU3yj6FugScf\n2bTtBsVnYQkTCutZn7vnVVaNx5MJyLug6o7/nPiyXMpZv4mcQBFwyJB35gUqbqp3\njx5yvsXi0Pi+8nNlNFdpN1Vrr66BYJ4QrV2NeaCvylmi0lvxdqwEJKlw0O3IEGlQ\nhFbgU8pSX9E10bbt7CUX4HYFIVOdXBVvoNig6PmWPORpLYQAZnOaWn0BuxwKl+LT\nX3Acj+zTSm3mJIqjG2R6skDnZX8akQWmAJhMo8Kw3qC6wn5ggF3FPwg+/ontNnIu\nhc2HYebtmgU3DzSeFz/kkL2SNaV5JRBgJb4/Q+mh3q1YbZVJMetvBikE/soXEmzi\nSfKk5jdKtLL3P9PndPiS+jECAwEAAQ==\n-----END PUBLIC KEY-----\n"))
	mockRedis, _ := redismock.NewClientMock()
	server := getServer(db, mockRedis)
	req := &v1.AuthenticateUserRequest{
		Username: "test_user",
		Password: "test_password",
	}

	ctx := context.Background()
	resp, err := server.AuthenticateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, true, resp.Success)
	cleanup()
}

func TestGrpcServer_Verify2FA(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	cleanup, _ := setupTestFiles("test_user")

	password, _ := passwordhash.Hash("test_password")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE username = $1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT $2")).
		WithArgs("test_user", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "ssh_pub_key"}).AddRow(1, "test_user", password, "test_ssh_pub_key"))

	mockRedis, rMock := redismock.NewClientMock()
	rMock.ExpectGet("test_token").SetVal("123") // Set the expected value for the Get method
	server := getServer(db, mockRedis)
	req := &v1.Verify2FARequest{
		Token: "test_token",
	}

	ctx := context.Background()
	resp, err := server.Verify2FA(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.JwtToken)
	assert.Equal(t, true, resp.Success)
	cleanup()
}

func getServer(db gorm.ConnPool, mockRedis *redis.Client) *GrpcServer {
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)

	return NewGrpcServer(
		slog.Default(),
		jwt.NewJwt(),
		mockRedis,
		gdb,
	)
}

func setupTestFiles(login string) (cleanup func(), err error) {
	dir := "./.ssh/" + login
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	publicKeyPath := dir + "/public_key.pem"
	privateKeyPath := dir + "/private_key.pem"

	if err := os.WriteFile(publicKeyPath, []byte("public_key_content"), 0600); err != nil {
		return nil, err
	}
	if err := os.WriteFile(privateKeyPath, []byte("private_key_content"), 0600); err != nil {
		return nil, err
	}

	cleanup = func() {
		os.RemoveAll("./.ssh/" + login)
	}

	return cleanup, nil
}

func TestGrpcServer_SetMetadataFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRedis, _ := redismock.NewClientMock()
	server := getServer(db, mockRedis)

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", "1"))

	req := &v1.SetMetadataFileRequest{
		Uuid:       "test_uuid",
		DataType:   v1.DataType_DATA_TYPE_BINARY,
		Metadata:   "test_metadata",
		UserPath:   "test_user_path",
		SizeChunks: 1024,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"storages\" WHERE uuid = $1 ORDER BY \"storages\".\"id\" LIMIT $2")).
		WithArgs("test_uuid", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "user_path", "metadata"}).AddRow(1, "test_uuid", "test_user_path", "test_metadata"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE \"storages\" SET \"data_type\"=$1,\"metadata\"=$2,\"size_bytes_partition\"=$3,\"user_path\"=$4 WHERE \"id\" = $5")).
		WithArgs(3, "test_metadata", 1024, "test_user_path", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	resp, err := server.SetMetadataFile(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
}

func TestGrpcServer_GetStoreDataList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRedis, _ := redismock.NewClientMock()
	server := getServer(db, mockRedis)

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", "1"))

	req := &v1.GetStoreDataListRequest{
		DataType: v1.DataType_DATA_TYPE_BINARY,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"storages\" WHERE user_id = $1 AND data_type = $2")).
		WithArgs(1, v1.DataType_DATA_TYPE_BINARY).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "user_path", "metadata"}).
			AddRow("test_uuid", "test_user_path", "{\"FileName\":\"test_file.txt\"}"))

	resp, err := server.GetStoreDataList(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Entries))
	assert.Equal(t, "test_uuid", resp.Entries[0].Uuid)
	assert.Equal(t, "test_user_path/", resp.Entries[0].UserPath)
}

func TestGrpcServer_DownloadFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRedis, _ := redismock.NewClientMock()
	server := getServer(db, mockRedis)

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", "1"))

	req := &v1.DownloadFileRequest{
		Uuid: "test_uuid",
	}
	file, _ := os.CreateTemp("", "test_file")
	_, err = file.Write([]byte("test_content"))
	assert.NoError(t, err)
	file.Close()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"storages\" WHERE user_id = $1 AND uuid = $2 ORDER BY \"storages\".\"id\" LIMIT $3")).
		WithArgs(1, "test_uuid", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "user_path", "metadata", "path"}).AddRow(1, "test_uuid", "test_user_path", "test_metadata", file.Name()))

	stream := &mockDownloadFileServer{}
	stream.ctx = ctx
	err = server.DownloadFile(req, stream)
	assert.NoError(t, err)

	resp := stream.resp
	assert.NotNil(t, resp)
	assert.Equal(t, v1.Status_STATUS_SUCCESS, resp.Status)
}

type mockDownloadFileServer struct {
	grpc.ServerStream
	ctx  context.Context
	resp *v1.DownloadFileResponse
}

func (m *mockDownloadFileServer) Context() context.Context {
	return m.ctx
}

func (m *mockDownloadFileServer) Send(resp *v1.DownloadFileResponse) error {
	m.resp = resp
	return nil
}

func (m *mockDownloadFileServer) Recv() (*v1.DownloadFileRequest, error) {
	return nil, io.EOF
}

func TestGrpcServer_GetMetadataFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRedis, _ := redismock.NewClientMock()
	server := getServer(db, mockRedis)

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", "1"))

	req := &v1.GetMetadataFileRequest{
		Uuid: "test_uuid",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"storages\" WHERE user_id = $1 AND uuid = $2 ORDER BY \"storages\".\"id\" LIMIT $3")).
		WithArgs(1, "test_uuid", 1).
		WillReturnRows(sqlmock.NewRows([]string{"metadata"}).
			AddRow("test_metadata"))

	resp, err := server.GetMetadataFile(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "test_metadata", resp.Metadata)
}

func TestGrpcServer_DeleteFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRedis, _ := redismock.NewClientMock()
	server := getServer(db, mockRedis)

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", "1"))

	req := &v1.DeleteFileRequest{
		Uuid: "test_uuid",
	}
	file, _ := os.CreateTemp("", "test_file")
	file.Close()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"storages\" WHERE user_id = $1 AND uuid = $2 ORDER BY \"storages\".\"id\" LIMIT $3")).
		WithArgs(1, "test_uuid", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "user_path", "metadata", "path"}).AddRow(1, "test_uuid", "test_user_path", "test_metadata", file.Name()))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM \"storages\" WHERE user_id = $1 AND uuid = $2")).
		WithArgs(1, "test_uuid").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	resp, err := server.DeleteFile(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
}
