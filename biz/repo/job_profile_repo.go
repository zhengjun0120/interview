package repo

import (
	"ai_interview/biz/entity"
	"context"
)

type JobProfileRepo interface {
	//新建岗位画像
	CreateJobProfile(ctx context.Context, jobProfile *entity.JobProfile) error

	//删除岗位画像
	DeleteJobProfile(ctx context.Context, jobProfileID string) error

	//更新岗位画像
	UpdateJobProfile(ctx context.Context, jobProfile *entity.JobProfile) error

	//获取用户的所有岗位画像
	GetJobProfileByUserID(ctx context.Context, userID string) ([]entity.JobProfile, error)

	//获取单个岗位画像
	GetJobProfileByJobProfileID(ctx context.Context, jobProfileID string) (*entity.JobProfile, error)

	//检查职位名称是否已存在
	CheckJobTitle(ctx context.Context, jobTitle string) (string, error)
}
