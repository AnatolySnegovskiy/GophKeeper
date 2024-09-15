package models

import (
	"goph_keeper/internal/server/services/entities"
	"goph_keeper/internal/services/grpc/goph_keeper/v1"
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

func (s *StorageModel) Create(userID uint, uuid string, path string) error {
	storage := &entities.StorageEntity{UserID: userID, Uuid: uuid, Path: path}
	return s.ifErrorLog(s.DB.Create(storage).Error)
}

func (s *StorageModel) GetListByDataType(u uint, dataType v1.DataType) ([]*entities.StorageEntity, error) {
	storages := make([]*entities.StorageEntity, 0)
	return storages, s.ifErrorLog(s.DB.Where("user_id = ? AND data_type = ?", u, dataType).Find(&storages).Error)
}

func (s *StorageModel) GetByUuid(u uint, uuid string) (*entities.StorageEntity, error) {
	storage := &entities.StorageEntity{}
	return storage, s.ifErrorLog(s.DB.Where("user_id = ? AND uuid = ?", u, uuid).First(&storage).Error)
}

func (s *StorageModel) UpdateMetadata(uuid string, dataType v1.DataType, metadata string, userPath string, chunks int32) error {
	storage := &entities.StorageEntity{}
	if err := s.DB.Where("uuid = ?", uuid).First(&storage).Error; err != nil {
		return s.ifErrorLog(err)
	}

	updates := map[string]interface{}{
		"DataType":           dataType,
		"Metadata":           metadata,
		"UserPath":           userPath,
		"SizeBytesPartition": int64(chunks),
	}

	return s.ifErrorLog(s.DB.Model(&storage).Updates(updates).Error)
}

func (s *StorageModel) Delete(u uint, uuid string) error {
	return s.ifErrorLog(s.DB.Where("user_id = ? AND uuid = ?", u, uuid).Delete(&entities.StorageEntity{}).Error)
}
