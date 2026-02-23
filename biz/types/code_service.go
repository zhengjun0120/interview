package types

import (
	"fmt"
	"math/rand"
)

type ICodeService interface {
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

type CaptchaMethodType string

const (
	CaptchaMethodTypeEmail CaptchaMethodType = "email"
	CaptchaMethodTypePhone CaptchaMethodType = "phone"
)

func (cs *Captcha) GenrateCaptchaCode() {
	num := rand.Intn(1000000)
	cs.CaptchaCode = fmt.Sprintf("%06d", num)
}
