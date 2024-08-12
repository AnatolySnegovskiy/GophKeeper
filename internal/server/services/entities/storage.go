package entities

import (
	"goph_keeper/internal/server/services/grpc/goph_keeper/v1"
	"gorm.io/gorm"
)

type StorageEntity struct {
	gorm.Model `json:"gorm.Model"`
	Path       string      `gorm:"unique;not null;size:255"`
	DataType   v1.DataType `gorm:"column:data_type;not null;index:idx_data_type"`
	Metadata   string      `gorm:"column:metadata;not null"`
	IsDeleted  bool        `gorm:"column:is_deleted;not null;default:false"`
	UserID     uint        `gorm:"not null;name:user_id;type:bigint;index:idx_user_id"`
}

func (s *StorageEntity) TableName() string {
	return "storages"
}
