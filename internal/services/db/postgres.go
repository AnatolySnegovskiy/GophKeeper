package db

import (
	"goph_keeper/internal/services/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entities.UserEntity{}, &entities.StorageEntity{})

	return db, err
}
