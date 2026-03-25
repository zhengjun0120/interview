package types

import "context"

type ICodeService interface {
	CaptchaSend(ctx context.Context, way CaptchaWayType, key string) error
}

type Captcha struct {
	Way         CaptchaWayType `json:"way"`
	Key         string         `json:"captcha_key"`
	RedisKey    string         `json:"redis_key"`
	CaptchaCode string         `json:"captcha_code"`
}

// 验证码用途
type CaptchaWayType string

const (
	CaptchaWayTypeLogin    CaptchaWayType = "login"    // 登录
	CaptchaWayTypeRegister CaptchaWayType = "register" // 注册
	CaptchaWayTypeReset    CaptchaWayType = "reset"    // 重置
)
