package storage

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/biz/user_service"
	"ai_interview/infra/database"
	"ai_interview/infra/storage/po"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var (
	HR     = "hr"
	Seeker = "seeker"
)

type UserStorage struct {
	db *gorm.DB
}

var us *UserStorage

func InitUserStorage() {
	db := database.GetDB()

	if err := db.AutoMigrate(&po.UserHR{}, &po.UserSeeker{}); err != nil {
		panic("user表自动迁移失败" + err.Error())
	}

	us = &UserStorage{db}
}

func GetUserStorage() repo.UserRepo {
	return us
}

func castUserHR2Entity(user *po.UserHR) *entity.User {
	return &entity.User{
		UserID:   user.UserID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Type:     HR,
	}
}

func castUserSeeker2Entity(user *po.UserSeeker) *entity.User {
	return &entity.User{
		UserID:   user.UserID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Type:     Seeker,
	}
}

func errorDB(err error) error {
	return fmt.Errorf("数据库错误:%v", err)
}

func (u *UserStorage) GetUserByID(ctx context.Context, userID string, userType string) (*entity.User, error) {
	if userType == HR {
		var user po.UserHR
		err := u.db.WithContext(ctx).Model(&po.UserHR{}).Where("user_id = ?", userID).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, user_service.USER_NOT_EXIST
			}
			return nil, errorDB(err)
		}
		return castUserHR2Entity(&user), nil
	} else if userType == Seeker {
		var user po.UserSeeker
		err := u.db.WithContext(ctx).Model(&po.UserSeeker{}).Where("user_id = ?", userID).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, user_service.USER_NOT_EXIST
			}
			return nil, errorDB(err)
		}
		return castUserSeeker2Entity(&user), nil
	} else {
		return nil, user_service.USER_TYPE_ERROR
	}
}

func (u *UserStorage) GetUserByEmail(ctx context.Context, email string, userType string) (*entity.User, error) {

	if userType == HR {
		var user po.UserHR
		err := u.db.WithContext(ctx).Model(&po.UserHR{}).Where("email = ?", email).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, user_service.USER_NOT_EXIST
			}
			return nil, errorDB(err)
		}

		return castUserHR2Entity(&user), nil
	} else if userType == Seeker {
		var user po.UserSeeker

		err := u.db.WithContext(ctx).Model(&po.UserSeeker{}).Where("email = ?", email).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, user_service.USER_NOT_EXIST
			}
			return nil, errorDB(err)
		}

		return castUserSeeker2Entity(&user), nil
	} else {
		return nil, user_service.USER_TYPE_ERROR
	}

}

func (u *UserStorage) CreateUser(ctx context.Context, user *entity.User, userType string) error {
	if userType == HR {
		var poUser po.UserHR
		poUser.UserID = user.UserID
		poUser.Username = user.Username
		poUser.Password = user.Password
		poUser.Email = user.Email

		err := u.db.WithContext(ctx).Model(&po.UserHR{}).Create(&poUser).Error
		if err != nil {
			return errorDB(err)
		}
		return nil
	} else if userType == Seeker {
		var poUser po.UserSeeker
		poUser.UserID = user.UserID
		poUser.Username = user.Username
		poUser.Password = user.Password
		poUser.Email = user.Email

		err := u.db.WithContext(ctx).Model(&po.UserSeeker{}).Create(&poUser).Error
		if err != nil {
			return errorDB(err)
		}
		return nil
	} else {
		return user_service.USER_TYPE_ERROR
	}
}

func (u *UserStorage) UpdatePassword(ctx context.Context, userID string, newPassword string, userType string) error {
	if newPassword == "" {
		return user_service.PASSWORD_NOT_NULL
	}
	upddate := map[string]interface{}{
		"password": newPassword,
	}
	if userType == HR {
		err := u.db.WithContext(ctx).Model(&po.UserHR{}).Where("user_id = ?", userID).Updates(upddate).Error
		if err != nil {
			return errorDB(err)
		}
		return nil
	} else if userType == Seeker {
		err := u.db.WithContext(ctx).Model(&po.UserSeeker{}).Where("user_id = ?", userID).Updates(upddate).Error
		if err != nil {
			return errorDB(err)
		}
		return nil
	} else {
		return user_service.USER_TYPE_ERROR
	}
}

func (u *UserStorage) UpdateUsername(ctx context.Context, userID string, newUsername string, userType string) error {
	if newUsername == "" {
		return user_service.USERNAME_NOT_NULL
	}

	update := map[string]interface{}{
		"username": newUsername,
	}

	if userType == HR {
		err := u.db.WithContext(ctx).Model(&po.UserHR{}).Where("user_id = ?", userID).Updates(update).Error
		if err != nil {
			return errorDB(err)
		}
		return nil
	} else if userType == Seeker {
		err := u.db.WithContext(ctx).Model(&po.UserSeeker{}).Where("user_id = ?", userID).Updates(update).Error
		if err != nil {
			return errorDB(err)
		}
		return nil
	} else {
		return user_service.USER_TYPE_ERROR
	}

}

func (u *UserStorage) CheckEmail(ctx context.Context, email string, userType string) (bool, error) {
	if userType == HR {
		var user po.UserHR
		err := u.db.WithContext(ctx).Model(&po.UserHR{}).Where("email = ?", email).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, errorDB(err)
		}
		return true, nil
	} else if userType == Seeker {
		var user po.UserSeeker
		err := u.db.WithContext(ctx).Model(&po.UserSeeker{}).Where("email = ?", email).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, errorDB(err)
		}
		return true, nil
	} else {
		return false, user_service.USER_TYPE_ERROR
	}
}
