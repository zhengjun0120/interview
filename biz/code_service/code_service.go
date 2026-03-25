package code_service

import (
	"ai_interview/biz/code_service/sms"
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/pkg/error_msg"
	"context"
)

type CodeService struct {
	CodeRepo repo.CodeRepo
}

func NewCodeService(CodeRepo repo.CodeRepo) *CodeService {
	return &CodeService{CodeRepo: CodeRepo}
}

func (cds *CodeService) CaptchaSend(ctx context.Context, way types.CaptchaWayType, key string) error {
	capcha := &types.Captcha{
		Way: way,
		Key: key,
	}

	// 存
	if err := cds.CodeRepo.CaptchaStash(ctx, capcha); err != nil {
		return err
	}

	// 发送
	if err := sms.SendCaptcha(key, capcha.CaptchaCode); err != nil {
		return error_msg.ErrorCaptchaSend
	}
	return nil
}
