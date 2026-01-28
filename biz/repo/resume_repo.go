package repo

import (
	"ai_interview/biz/entity"
	"context"
)

type ResumeRepo interface {
	//新建简历
	CreateResume(ctx context.Context, resume *entity.Resume) error
}
