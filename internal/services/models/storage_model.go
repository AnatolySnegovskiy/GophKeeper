package models

import (
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"gorm.io/gorm"
	"log/slog"
)

type StorageModel struct {
	gorm.DB
	BaseModel
}

func NewStorageModel(db *gorm.DB, logger *slog.Logger) *StorageModel {
	return &StorageModel{DB: *db, BaseModel: BaseModel{Logger: logger}}
}

func (s *StorageModel) Create(userID uint, path string, metaData string, dataType v1.DataType) error {
	storage := &entities.StorageEntity{UserID: userID, Path: path, Metadata: metaData, DataType: dataType}
	return s.ifErrorLog(s.DB.Create(storage).Error)
}

func (s *StorageModel) GetListByDataType(u uint, dataType v1.DataType) ([]*entities.StorageEntity, error) {
	storages := make([]*entities.StorageEntity, 0)
	return storages, s.ifErrorLog(s.DB.Where("user_id = ? AND data_type = ?", u, dataType).Find(&storages).Error)
}
