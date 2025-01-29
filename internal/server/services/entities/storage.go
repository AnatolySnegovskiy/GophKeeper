package entities

import (
	"goph_keeper/internal/services/grpc/goph_keeper/v1"
	"time"
)

type StorageEntity struct {
	ID                 uint `gorm:"primaryKey"`
	CreatedAt          time.Time
	Uuid               string      `gorm:"unique;not null;size:255;index:idx_uuid"`
	UserID             uint        `gorm:"not null;name:user_id;type:bigint;index:idx_user_id"`
	SizeBytesPartition int64       `gorm:"column:size_bytes_partition;not null;default:0"`
	Path               string      `gorm:"unique;not null;size:255"`
	UserPath           string      `gorm:"not null"`
	DataType           v1.DataType `gorm:"column:data_type;not null;index:idx_data_type"`
	Metadata           string      `gorm:"column:metadata;not null"`
}

func (s *StorageEntity) TableName() string {
	return "storages"
}

type FileMetadata struct {
	FileName        string `json:"file_name"`
	FileExtension   string `json:"file_extension"`
	MemType         string `json:"mem_type"`
	IsCompressed    bool   `json:"is_compressed"`
	CompressionType string `json:"compression_type"`
	FileSize        int64  `json:"file_size"`
}
