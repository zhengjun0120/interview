package repo

import (
	"ai_interview/biz/entity"
	"context"
)

type ResumeRepo interface {
	//新建简历
	CreateResume(ctx context.Context, resume *entity.Resume) error

	//删除简历
	DeleteResume(ctx context.Context, resumeID, userID string) error

	//获取简历
	GetResume(ctx context.Context, resumeID, userID string) (*entity.Resume, error)

	//更新简历
	UpdateResume(ctx context.Context, newResume *entity.Resume) error

	//新建人才
	CreateTalent(ctx context.Context, talent *entity.Talent) error

	//根据用户ID获取人才
	GetTalentByUserID(ctx context.Context, userID string) ([]entity.Talent, error)

	//根据人才ID获取人才
	GetTalentByID(ctx context.Context, talentID string) (*entity.Talent, error)

	//根据用户ID获取已面试的人才
	GetTalentInterviewByUserID(ctx context.Context, userID string) ([]entity.Talent, error)

	//更新人才
	UpdateTalent(ctx context.Context, newTalent *entity.Talent) error

	//删除人才
	DeleteTalent(ctx context.Context, talentID, userID string) error

	//Hr端 一键上传简历+人才
	CreateResumeAndTalent(ctx context.Context, resume *entity.Resume, talent *entity.Talent) error

	//Hr端 一键删除简历+人才
	DeleteResumeAndTalent(ctx context.Context, resumeID, talentID, userID string) error
}
