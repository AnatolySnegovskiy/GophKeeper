package models

import (
	"fmt"
	passwordhash "github.com/vzglad-smerti/password_hash"
	"goph_keeper/internal/server/services/entities"
	"gorm.io/gorm"
	"log/slog"
)

type UserModel struct {
	gorm.DB
	BaseModel
}

func NewUserModel(db *gorm.DB, logger *slog.Logger) *UserModel {
	return &UserModel{DB: *db, BaseModel: BaseModel{Logger: logger}}
}

func (u *UserModel) Create(userName string, password string, sshPubKey string) error {
	password, err := passwordhash.Hash(password)

	if u.ifErrorLog(err) != nil {
		return fmt.Errorf("failed to hash password")
	}

	if user, _ := u.GetByUsername(userName); user.ID != 0 {
		return fmt.Errorf("user %s already exists", userName)
	}

	user := &entities.UserEntity{Username: userName, Password: password, SshPubKey: sshPubKey}
	return u.DB.Create(user).Error
}

func (u *UserModel) Get(userId int) (*entities.UserEntity, error) {
	user := &entities.UserEntity{}
	return user, u.ifErrorLog(u.DB.Where("id = ?", userId).First(user).Error)
}

func (u *UserModel) GetByUsername(username string) (*entities.UserEntity, error) {
	user := &entities.UserEntity{}
	return user, u.ifErrorLog(u.DB.Where("username = ?", username).First(user).Error)
}

func (u *UserModel) Auth(username string, password string) (*entities.UserEntity, error) {
	user, err := u.GetByUsername(username)

	if u.ifErrorLog(err) != nil {
		return nil, fmt.Errorf("user not found")
	}

	if ok, err := passwordhash.Verify(user.Password, password); err != nil || !ok {
		_ = u.ifErrorLog(err)
		return nil, fmt.Errorf("incorrect password")
	}

	return user, nil
}
