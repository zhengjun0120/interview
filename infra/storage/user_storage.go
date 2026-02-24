package storage

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/infra/database"
	"ai_interview/infra/storage/po"
	"ai_interview/pkg/error_msg"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	HR     = "hr"
	Seeker = "seeker"
)

var us *UserStorage

type UserStorage struct {
	db     *gorm.DB
	client *redis.Client
}

func InitUserStorage() {
	db := database.GetDB()
	client := database.GetRedis()

	if err := db.AutoMigrate(&po.User{}); err != nil {
		panic("user表自动迁移失败" + err.Error())
	}

	us = &UserStorage{
		db,
		client,
	}
}

func GetUserStorage() repo.UserRepo {
	return us
}

func errorDB(err error) error {
	return fmt.Errorf("数据库错误:%v", err)
}

func (u *UserStorage) GetUserByID(ctx context.Context, userID string, userType string) (*entity.User, error) {

	if userType != HR && userType != Seeker {
		return nil, error_msg.USER_TYPE_ERROR
	}

	var user po.User
	err := u.db.WithContext(ctx).Model(&po.User{}).Where("type = ? AND user_id = ?", userType, userID).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error_msg.USER_NOT_EXIST
		}
		return nil, errorDB(err)
	}

	return &entity.User{
		UserID:   user.UserID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Type:     user.Type,
	}, nil
}

func (u *UserStorage) GetUserByEmail(ctx context.Context, email string, userType string) (*entity.User, error) {

	if userType != HR && userType != Seeker {
		return nil, error_msg.USER_TYPE_ERROR
	}

	var user po.User
	err := u.db.WithContext(ctx).Model(&po.User{}).Where("type = ? AND email = ?", userType, email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error_msg.USER_NOT_EXIST
		}
		return nil, errorDB(err)
	}

	return &entity.User{
		UserID:   user.UserID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Type:     user.Type,
	}, nil

}

func (u *UserStorage) CreateUser(ctx context.Context, user *entity.User, userType string) error {
	if userType != HR && userType != Seeker {
		return error_msg.USER_TYPE_ERROR
	}

	var poUser po.User
	poUser.UserID = user.UserID
	poUser.Username = user.Username
	poUser.Password = user.Password
	poUser.Email = user.Email
	poUser.Type = userType

	err := u.db.WithContext(ctx).Model(&po.User{}).Create(&poUser).Error
	if err != nil {
		return errorDB(err)
	}

	return nil
}

func (u *UserStorage) UpdatePassword(ctx context.Context, userID string, newPassword string, userType string) error {
	if newPassword == "" {
		return error_msg.PASSWORD_NOT_NULL
	} else if userType != HR && userType != Seeker {
		return error_msg.USER_TYPE_ERROR
	}
	update := map[string]interface{}{
		"password": newPassword,
	}

	err := u.db.WithContext(ctx).Model(&po.User{}).Where("type = ? AND user_id = ?", userType, userID).Updates(update).Error
	if err != nil {
		return errorDB(err)
	}

	return nil

}

func (u *UserStorage) UpdateUsername(ctx context.Context, userID string, newUsername string, userType string) error {
	if newUsername == "" {
		return error_msg.USERNAME_NOT_NULL
	} else if userType != HR && userType != Seeker {
		return error_msg.USER_TYPE_ERROR
	}

	update := map[string]interface{}{
		"username": newUsername,
	}

	err := u.db.WithContext(ctx).Model(&po.User{}).Where("type = ? AND user_id = ?", userType, userID).Updates(update).Error
	if err != nil {
		return errorDB(err)
	}

	return nil

}

func (u *UserStorage) CheckEmail(ctx context.Context, email string, userType string) (bool, error) {
	var user po.User
	err := u.db.WithContext(ctx).Model(&po.User{}).Where("type = ? AND email = ?", userType, email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errorDB(err)
	}
	return true, nil
}

func (rds *UserStorage) Get(c context.Context, key string) (interface{}, error) {
	return rds.client.Get(c, key).Result()
}

func (rds *UserStorage) Set(c context.Context, key string, value any, expiration time.Duration) error {
	return rds.client.Set(c, key, value, expiration).Err()
}

func (rds *UserStorage) HGet(c context.Context, key string, receiver any) error {
	val, err := rds.client.Get(c, key).Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(val, receiver)
	if err != nil {
		return err
	}
	return nil
}

func (rds *UserStorage) HSet(c context.Context, key string, value any, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rds.client.Set(c, key, val, expiration).Err()
}

func (rds *UserStorage) Del(c context.Context, key string) error {
	_, err := rds.client.Del(c, key).Result()
	if err != nil {
		return err
	}
	return nil
}
