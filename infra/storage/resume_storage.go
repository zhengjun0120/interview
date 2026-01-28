package storage

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/infra/database"
	"ai_interview/infra/storage/po"
	"ai_interview/pkg/error_msg"
	"context"
	"gorm.io/gorm"
)

var rs *ResumeStorage

type ResumeStorage struct {
	db *gorm.DB
}

func InitResumeStorage() {
	db := database.GetDB()

	if err := db.AutoMigrate(&po.Resume{}); err != nil {
		panic("resume表自动迁移失败: " + err.Error())
	}

	rs = &ResumeStorage{db}
}

func GetResumeStorage() repo.ResumeRepo {
	return rs
}

func (r *ResumeStorage) CreateResume(ctx context.Context, resume *entity.Resume) error {
	if resume.ResumeName != "" {
		return error_msg.RESUME_NAME_NOT_NULL
	}

	var resumePO po.Resume
	resumePO.ResumeID = resume.ResumeID
	resumePO.ResumeName = resume.ResumeName
	resumePO.ResumeUrl = resume.ResumeUrl
	resumePO.UserID = resume.UserID

	err := r.db.WithContext(ctx).Create(&resumePO).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}
