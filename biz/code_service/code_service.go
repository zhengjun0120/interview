package code_service

import (
	"ai_interview/biz/code_service/sms"
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/pkg/error_msg"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

type CodeService struct {
	c   context.Context
	rds repo.UserRepo
	sms sms.SmsService
}

func NewCodeService(c context.Context, rds repo.UserRepo, sms sms.SmsService) *CodeService {
	return &CodeService{c, rds, sms}
}

func (cs *CodeService) CaptchaSend(way types.CaptchaWayType, key string) error {
	// 先存再发送
	capcha := &types.Captcha{
		Way: way,
		Key: key,
	}

	// 存
	if err := cs.captchaStash(capcha); err != nil {
		return err
	}

	// 发送
	//if err := cs.sms.SendCaptcha(cs.c, key, capcha.CaptchaCode); err != nil {
	//	return error_msg.ErrorCaptchaSend
	//}
	return nil
}

// 储存验证码
func (cs *CodeService) captchaStash(captcha *types.Captcha) error {
	// 整合
	GenerateCaptcha(captcha)
	GenerateCaptchaCode(captcha)
	// 存到redis
	err := cs.rds.Set(cs.c, captcha.RedisKey, captcha.CaptchaCode, time.Minute*1)
	if err != nil {
		return error_msg.ErrorCode2Redis
	}
	return nil
}

// 验证码检查
func (cs *CodeService) CaptchaCheck(way types.CaptchaWayType, code string) error {
	// 拿出来
	capcha := &types.Captcha{
		Way: way,
	}
	GenerateCaptcha(capcha)
	ret, err := cs.rds.Get(cs.c, capcha.RedisKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 验证码过期
			return error_msg.ErrorCaptchaExpire
		}
		return err
	}
	str, ok := ret.(string)
	if !ok {
		// 安全断言失败
		return error_msg.ErrorSecurityAssertion
	}

	if str != code {
		// 验证码错误
		return error_msg.ErrorCaptchaCheck
	}

	return nil
}

func GenerateCaptcha(captcha *types.Captcha) {
	captcha.RedisKey = "Captcha:" + string(captcha.Way) + ":" + captcha.Key
}

func GenerateCaptchaCode(captcha *types.Captcha) {
	num := rand.Intn(1000000)
	captcha.CaptchaCode = fmt.Sprintf("%06d", num)
}
