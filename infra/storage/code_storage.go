package storage

import (
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/infra/database"
	"ai_interview/pkg/error_msg"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

var cds *CodeStorage

type CodeStorage struct {
	client *redis.Client
}

func InitCodeStorage() {
	client := database.GetRedis()

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("Redis 连接失败" + err.Error())
	}

	cds = &CodeStorage{client: client}
}

func GetCodeStorage() repo.CodeRepo {
	return cds
}

func (cds *CodeStorage) CaptchaStash(ctx context.Context, captcha *types.Captcha) error {
	captcha.RedisKey = "Captcha:" + string(captcha.Way) + ":" + captcha.Key

	num := rand.Intn(900000) + 100000
	captcha.CaptchaCode = fmt.Sprintf("%06d", num)

	err := cds.client.Set(ctx, captcha.RedisKey, captcha.CaptchaCode, time.Minute*5).Err()
	if err != nil {
		return error_msg.ErrorCode2Redis
	}
	return nil
}

func (cds *CodeStorage) CaptchaCheck(ctx context.Context, way types.CaptchaWayType, key string, code string) error {
	captcha := &types.Captcha{
		Way: way,
		Key: key,
	}

	captcha.RedisKey = "Captcha:" + string(captcha.Way) + ":" + captcha.Key

	ret, err := cds.client.Get(ctx, captcha.RedisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 验证码过期
			return error_msg.ErrorCaptchaExpire
		}
		return err
	}
	if ret != code {
		// 验证码错误
		return error_msg.ErrorCaptchaCheck
	}

	return nil
}
