package repo

import (
	"ai_interview/biz/types"
	"context"
)

type CodeRepo interface {
	//储存验证码
	CaptchaStash(ctx context.Context, captcha *types.Captcha) error

	//检验验证码
	CaptchaCheck(ctx context.Context, way types.CaptchaWayType, email string, code string) error
}
