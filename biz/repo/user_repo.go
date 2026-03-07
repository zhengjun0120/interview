package repo

import (
	"ai_interview/biz/entity"
	"context"
	"time"
)

type UserRepo interface {
	// 根据用户ID获取用户信息
	GetUserByID(ctx context.Context, userID string, userType string) (*entity.User, error)

	// 根据用户邮箱获取用户信息
	GetUserByEmail(ctx context.Context, email string, userType string) (*entity.User, error)

	// 创建新用户
	CreateUser(ctx context.Context, user *entity.User, userType string) error

	// 修改密码
	UpdatePassword(ctx context.Context, userID string, newPassword string, userType string) error

	//修改昵称
	UpdateUsername(ctx context.Context, userID string, newUsername string, userType string) error

	//检查邮箱是否注册
	CheckEmail(ctx context.Context, email string, userType string) (bool, error)

	//检查用户是否存在
	CheckUser(ctx context.Context, userID string) (bool, error)

	//redis操作
	//根据key获取用户信息
	Get(c context.Context, key string) (interface{}, error)

	//储存用户信息
	Set(c context.Context, key string, value any, expiration time.Duration) error

	//哈希存储和获取
	HGet(c context.Context, key string, receiver any) error
	HSet(c context.Context, key string, value any, expiration time.Duration) error

	//根据key删用户
	Del(c context.Context, key string) error
}
