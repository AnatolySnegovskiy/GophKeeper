package db

import (
	entities2 "goph_keeper/internal/server/services/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entities2.UserEntity{}, &entities2.StorageEntity{})

	return db, err
}
