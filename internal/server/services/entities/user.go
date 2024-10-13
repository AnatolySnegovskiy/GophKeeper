package entities

import (
	"gorm.io/gorm"
)

type UserEntity struct {
	gorm.Model
	Username  string          `gorm:"unique;not null;size:255"`
	Password  string          `gorm:"not null;size:255"`
	SshPubKey string          `gorm:"not null"`
	Storages  []StorageEntity `gorm:"foreignKey:UserID"`
}

func (u *UserEntity) TableName() string {
	return "users"
}
