package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	entities2 "goph_keeper/internal/server/services/entities"
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
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(s.JWTStreamInterceptor), grpc.UnaryInterceptor(s.JWTUnaryInterceptor))
	v1.RegisterGophKeeperV1ServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *GrpcServer) RegisterUser(ctx context.Context, req *v1.RegisterUserRequest) (*v1.RegisterUserResponse, error) {
	s.db.WithContext(ctx)
	userModel := models.NewUserModel(s.db, s.logger)
	err := userModel.Create(req.Username, req.Password, req.SshPubKey)

	return &v1.RegisterUserResponse{
		Success: err == nil,
	}, err
}

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

func (s *GrpcServer) UploadFile(srv v1.GophKeeperV1Service_UploadFileServer) error {
	userId := s.getUserId(srv.Context())

	if userId == 0 {
		return status.Error(codes.Unauthenticated, "invalid token")
	}
	uuidFile := uuid.New().String()
	filename := s.storagePath + "/" + uuidFile
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		data, err := srv.Recv()

		if err != nil && data != nil && data.Status == v1.Status_STATUS_CANCELLED {
			return srv.SendAndClose(&v1.UploadFileResponse{
				Success: false,
				Uuid:    "",
			})
		}

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

	err = storageModel.Create(uint(userId), uuidFile, filename)
	if err != nil {
		return status.Error(codes.Internal, "failed to store data")
	}

	return srv.SendAndClose(&v1.UploadFileResponse{
		Success: true,
		Uuid:    uuidFile,
	})
}

func (s *GrpcServer) SetMetadataFile(ctx context.Context, req *v1.SetMetadataFileRequest) (*v1.SetMetadataFileResponse, error) {
	userId := s.getUserId(ctx)
	if userId == 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	storageModel := models.NewStorageModel(s.db, s.logger)
	err := storageModel.UpdateMetadata(req.Uuid, req.DataType, req.Metadata, req.UserPath, req.SizeChunks)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update metadata")
	}

	return &v1.SetMetadataFileResponse{
		Success: true,
	}, nil
}

func (s *GrpcServer) GetStoreDataList(ctx context.Context, req *v1.GetStoreDataListRequest) (*v1.GetStoreDataListResponse, error) {
	userId := s.getUserId(ctx)
	if userId == 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	storageModel := models.NewStorageModel(s.db, s.logger)
	entities, err := storageModel.GetListByDataType(uint(userId), req.DataType)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get data")
	}

	var listResponce = make([]*v1.ListDataEntry, 0)

	for _, entity := range entities {
		fileMetadata := entities2.FileMetadata{}
		err = json.Unmarshal([]byte(entity.Metadata), &fileMetadata)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to get data")
		}

		listResponce = append(listResponce, &v1.ListDataEntry{
			Uuid:     entity.Uuid,
			UserPath: entity.UserPath + "/" + fileMetadata.FileName,
		})
	}

	return &v1.GetStoreDataListResponse{
		Entries: listResponce,
	}, nil
}

func (s *GrpcServer) DownloadFile(req *v1.DownloadFileRequest, serv v1.GophKeeperV1Service_DownloadFileServer) error {
	userId := s.getUserId(serv.Context())
	if userId == 0 {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	storageModel := models.NewStorageModel(s.db, s.logger)
	data, err := storageModel.GetByUuid(uint(userId), req.Uuid)
	if err != nil {
		return status.Error(codes.Internal, "failed to get data")
	}

	file, err := os.Open(data.Path)
	if err != nil {
		return status.Error(codes.Internal, "failed to open data")
	}
	defer file.Close()

	buf := make([]byte, data.SizeBytesPartition)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return status.Error(codes.Internal, "failed to read data")
		}

		if n == 0 {
			break
		}

		err = serv.Send(&v1.DownloadFileResponse{
			Status: v1.Status_STATUS_PROCESSING,
			Data:   buf[:n],
		})

		if err != nil {
			return status.Error(codes.Internal, "failed to send data")
		}
	}

	return serv.Send(&v1.DownloadFileResponse{
		Status: v1.Status_STATUS_SUCCESS,
		Data:   nil,
	})
}

func (s *GrpcServer) GetMetadataFile(ctx context.Context, req *v1.GetMetadataFileRequest) (*v1.GetMetadataFileResponse, error) {
	userId := s.getUserId(ctx)
	if userId == 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	storageModel := models.NewStorageModel(s.db, s.logger)

	data, err := storageModel.GetByUuid(uint(userId), req.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get data")
	}

	return &v1.GetMetadataFileResponse{
		Metadata: data.Metadata,
	}, nil
}

func (s *GrpcServer) DeleteFile(ctx context.Context, req *v1.DeleteFileRequest) (*v1.DeleteFileResponse, error) {
	userId := s.getUserId(ctx)
	if userId == 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	storageModel := models.NewStorageModel(s.db, s.logger)
	data, err := storageModel.GetByUuid(uint(userId), req.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete data")
	}

	err = os.Remove(data.Path)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete data")
	}

	err = storageModel.Delete(uint(userId), req.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete data")
	}

	return &v1.DeleteFileResponse{
		Success: true,
	}, nil
}

func (s *GrpcServer) getUserId(ctx context.Context) int {
	md, ok := metadata.FromIncomingContext(ctx)
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
