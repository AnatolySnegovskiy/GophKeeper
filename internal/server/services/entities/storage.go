package entities

import (
	"goph_keeper/internal/services/grpc/goph_keeper/v1"
	"gorm.io/gorm"
)

type StorageEntity struct {
	gorm.Model         `json:"gorm.Model"`
	Uuid               string      `gorm:"unique;not null;size:255;index:idx_uuid"`
	UserID             uint        `gorm:"not null;name:user_id;type:bigint;index:idx_user_id"`
	SizeBytesPartition int64       `gorm:"column:size_bytes_partition;not null;default:0"`
	Path               string      `gorm:"unique;not null;size:255"`
	UserPath           string      `gorm:"unique;not null"`
	DataType           v1.DataType `gorm:"column:data_type;not null;index:idx_data_type"`
	Metadata           string      `gorm:"column:metadata;not null"`
	IsDeleted          bool        `gorm:"column:is_deleted;not null;default:false"`
}

func (s *StorageEntity) TableName() string {
	return "storages"
}
