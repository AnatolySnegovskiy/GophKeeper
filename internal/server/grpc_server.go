package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"goph_keeper/internal/server/services/jwt"
	"goph_keeper/internal/server/services/models"
	"goph_keeper/internal/services"
	"goph_keeper/internal/services/grpc/goph_keeper/v1"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type GrpcServer struct {
	v1.UnimplementedGophKeeperV1ServiceServer
	logger      *slog.Logger
	jwt         *jwt.Jwt
	redis       *redis.Client
	db          *gorm.DB
	storagePath string
}

func NewGrpcServer(logger *slog.Logger, jwt *jwt.Jwt, redis *redis.Client, db *gorm.DB) *GrpcServer {
	server := &GrpcServer{}
	server.logger = logger
	server.jwt = jwt
	server.redis = redis
	server.db = db
	server.storagePath = "./cmd/server/storage"
	_ = os.Mkdir(server.storagePath, os.ModePerm)
	return server
}

func (s *GrpcServer) Run(lis net.Listener) error {
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(s.JWTStreamInterceptor))
	v1.RegisterGophKeeperV1ServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// RegisterUser handles user registration.
func (s *GrpcServer) RegisterUser(ctx context.Context, req *v1.RegisterUserRequest) (*v1.RegisterUserResponse, error) {
	s.db.WithContext(ctx)
	userModel := models.NewUserModel(s.db, s.logger)
	err := userModel.Create(req.Username, req.Password, req.SshPubKey)
	message := ""

	if err != nil {
		message = "Failed to create user"
	} else {
		message = "User created successfully"
	}

	return &v1.RegisterUserResponse{
		Success: err == nil,
		Message: message,
	}, err
}

// AuthenticateUser handles user authentication.
func (s *GrpcServer) AuthenticateUser(ctx context.Context, req *v1.AuthenticateUserRequest) (*v1.AuthenticateUserResponse, error) {
	s.db.WithContext(ctx)
	userModel := models.NewUserModel(s.db, s.logger)
	user, err := userModel.Auth(req.Username, req.Password)
	randomToken := ""
	token := ""

	if err == nil {
		ssh := services.NewSshKeyGen()
		randomToken = strconv.Itoa(int(user.ID * uint(rand.Uint32())))
		token, err = ssh.EncryptMessage(randomToken, user.SshPubKey)
	}

	if token != "" {
		s.redis.Set(ctx, randomToken, user.ID, 1*time.Minute)
	}

	return &v1.AuthenticateUserResponse{
		Success: err == nil,
		Token:   token,
	}, err
}

func (s *GrpcServer) Verify2FA(ctx context.Context, req *v1.Verify2FARequest) (*v1.Verify2FAResponse, error) {
	s.db.WithContext(ctx)
	userId, err := s.redis.Get(ctx, req.Token).Int()
	if err != nil {
		return &v1.Verify2FAResponse{
			Success:  false,
			JwtToken: "",
		}, err
	}

	token, err := s.jwt.CreateToken()

	if userId != 0 && token != "" {
		s.redis.Set(ctx, token, userId, s.jwt.ExpiredAt)
	}

	return &v1.Verify2FAResponse{
		Success:  err == nil,
		JwtToken: token,
	}, err
}

func (s *GrpcServer) UploadData(srv v1.GophKeeperV1Service_UploadDataServer) error {
	userId := s.getUserId(srv)

	if userId == 0 {
		return status.Error(codes.Unauthenticated, "invalid token")
	}
	filename := s.storagePath + "/" + uuid.New().String()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	data := &v1.UploadDataRequest{}
	for {
		data, err := srv.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = file.Write(data.GetData())
		if err != nil {
			return status.Error(codes.Internal, "failed to store data")
		}
	}

	storageModel := models.NewStorageModel(s.db, s.logger)

	err = storageModel.Create(uint(userId), filename, data.Metadata, data.DataType)
	if err != nil {
		return status.Error(codes.Internal, "failed to store data")
	}

	return srv.SendAndClose(&v1.UploadDataResponse{
		Success: true,
		Message: "Data stored successfully",
	})
}

// RequestPrivateData handles requesting private data.
func (s *GrpcServer) DownloadData(req *v1.DownloadDataRequest, serv v1.GophKeeperV1Service_DownloadDataServer) error {
	userId := s.getUserId(serv)
	if userId == 0 {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	storageModel := models.NewStorageModel(s.db, s.logger)
	data, err := storageModel.GetListByUuid(uint(userId), req.Uuid)
	if err != nil {
		return status.Error(codes.Internal, "failed to get data")
	}

	file, err := os.Open(data.Path)
	if err != nil {
		return status.Error(codes.Internal, "failed to open data")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return status.Error(codes.Internal, "failed to read data")
	}
	//TODO: send Partial data
	err = serv.Send(&v1.DownloadDataResponse{
		Success:  true,
		Message:  "Data sent successfully",
		Metadata: data.Metadata,
		Data:     fileBytes,
	})
	if err != nil {
		return status.Error(codes.Internal, "failed to send data")
	}

	return serv.Send(&v1.DownloadDataResponse{
		Success: true,
		Message: "Data sent successfully",
	})
}

func (s *GrpcServer) getUserId(srv grpc.ServerStream) int {
	md, ok := metadata.FromIncomingContext(srv.Context())
	if !ok {
		return 0
	}

	userIDs := md["user_id"]
	if len(userIDs) == 0 {
		return 0
	}

	userID, err := strconv.Atoi(userIDs[0])
	if err != nil {
		return 0
	}

	return userID
}

func (s *GrpcServer) Stop() {

}
