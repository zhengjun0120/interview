package repo

import "ai_interview/biz/types"

type ICodeService interface {
	// 验证码发送
	CaptchaSend(way types.CaptchaWayType, key string) error
	// 验证码校验
	CaptchaCheck(way types.CaptchaWayType, key, code string) error
}
