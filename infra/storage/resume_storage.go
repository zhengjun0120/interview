package storage

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/infra/database"
	"ai_interview/infra/storage/po"
	"ai_interview/pkg/error_msg"
	"ai_interview/pkg/zlog"
	"context"
	"errors"
	"gorm.io/gorm"
)

var rs *ResumeStorage

type ResumeStorage struct {
	db *gorm.DB
}

func InitResumeStorage() {
	db := database.GetDB()

	if err := db.AutoMigrate(&po.Resume{}, &po.TalentPool{}); err != nil {
		panic("resume表或talent_pool表自动迁移失败: " + err.Error())
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
	resumePO.ResumeStr = resume.ResumeStr

	err := r.db.WithContext(ctx).Create(&resumePO).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (r *ResumeStorage) DeleteResume(ctx context.Context, resumeID, userID string) error {
	if resumeID == "" {
		return error_msg.RESUME_ID_NOT_NULL
	}

	err := r.db.Model(&po.Resume{}).WithContext(ctx).Where("resume_id = ? AND user_id = ?", resumeID, userID).Delete(&po.Resume{}).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (r *ResumeStorage) GetResume(ctx context.Context, resumeID, userID string) (*entity.Resume, error) {
	if resumeID == "" {
		return nil, error_msg.RESUME_ID_NOT_NULL
	}

	var resumePO po.Resume
	err := r.db.Model(&po.Resume{}).WithContext(ctx).Where("resume_id = ? AND user_id = ?", resumeID, userID).First(&resumePO).Error
	if err != nil {
		return nil, errorDB(err)
	}
	return &entity.Resume{
		ResumeID:   resumePO.ResumeID,
		ResumeName: resumePO.ResumeName,
		ResumeUrl:  resumePO.ResumeUrl,
		UserID:     resumePO.UserID,
		ResumeStr:  resumePO.ResumeStr,
	}, nil
}

func (r *ResumeStorage) UpdateResume(ctx context.Context, newResume *entity.Resume) error {
	if newResume.ResumeID == "" {
		return error_msg.RESUME_ID_NOT_NULL
	} else if newResume.ResumeName == "" {
		return error_msg.RESUME_NAME_NOT_NULL
	} else if len(newResume.ResumeName) >= 30 {
		return error_msg.RESUME_NAME_TOO_LONG
	}

	updates := map[string]interface{}{
		"resume_name": newResume.ResumeName,
		"resume_url":  newResume.ResumeUrl,
	}

	err := r.db.Model(&po.Resume{}).WithContext(ctx).Where("resume_id = ? AND user_id = ?", newResume.ResumeID, newResume.UserID).Updates(updates).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (r *ResumeStorage) CreateTalent(ctx context.Context, talent *entity.Talent) error {
	if talent.TalentID == "" {
		return error_msg.TALENT_ID_NOT_NULL
	} else if talent.FullName == "" {
		return error_msg.TALENT_NAME_NOT_NULL
	} else if len(talent.FullName) >= 30 {
		return error_msg.TALENT_NAME_TOO_LONG
	} else if talent.TargetPosition == "" {
		return error_msg.TALENT_TARGET_POSITION_NOT_NULL
	} else if talent.MatchScore < 0 || talent.MatchScore > 100 {
		return error_msg.TALENT_MATCH_SCORE_INVALID
	}

	var talentPO po.TalentPool
	talentPO.TalentID = talent.TalentID
	talentPO.FullName = talent.FullName
	talentPO.TargetPosition = talent.TargetPosition
	talentPO.MatchScore = talent.MatchScore
	talentPO.CoreAdvantages = talent.CoreAdvantages
	talentPO.HireStatus = talent.HireStatus
	talentPO.InterviewStatus = talent.InterviewStatus
	talentPO.UserID = talent.UserID

	err := r.db.Model(&po.TalentPool{}).WithContext(ctx).Create(&talentPO).Error
	if err != nil {
		return errorDB(err)
	}
	return nil

}

func (r *ResumeStorage) GetTalentByUserID(ctx context.Context, userID string) ([]entity.Talent, error) {

	var talentPOs []po.TalentPool
	err := r.db.Model(&po.TalentPool{}).WithContext(ctx).Where("user_id = ?", userID).Find(&talentPOs).Error
	if err != nil {
		return nil, errorDB(err)
	}

	talents := make([]entity.Talent, len(talentPOs))
	for i, talentPO := range talentPOs {
		talents[i] = entity.Talent{
			TalentID:        talentPO.TalentID,
			FullName:        talentPO.FullName,
			TargetPosition:  talentPO.TargetPosition,
			MatchScore:      talentPO.MatchScore,
			CoreAdvantages:  talentPO.CoreAdvantages,
			HireStatus:      talentPO.HireStatus,
			InterviewStatus: talentPO.InterviewStatus,
			UserID:          talentPO.UserID,
			CreatedAt:       talentPO.CreatedAt,
		}
	}
	return talents, nil
}

func (r *ResumeStorage) GetTalentByID(ctx context.Context, talentID string) (*entity.Talent, error) {
	if talentID == "" {
		return nil, error_msg.TALENT_ID_NOT_NULL
	}

	var talentPO po.TalentPool
	err := r.db.Model(&po.TalentPool{}).WithContext(ctx).Where("talent_id = ?", talentID).First(&talentPO).Error
	if err != nil {
		return nil, errorDB(err)
	}
	return &entity.Talent{
		TalentID:        talentPO.TalentID,
		FullName:        talentPO.FullName,
		TargetPosition:  talentPO.TargetPosition,
		MatchScore:      talentPO.MatchScore,
		CoreAdvantages:  talentPO.CoreAdvantages,
		HireStatus:      talentPO.HireStatus,
		InterviewStatus: talentPO.InterviewStatus,
		UserID:          talentPO.UserID,
		CreatedAt:       talentPO.CreatedAt,
	}, nil
}

func (r *ResumeStorage) GetTalentInterviewByUserID(ctx context.Context, userID string) ([]entity.Talent, error) {
	var talentPOs []po.TalentPool

	err := r.db.Model(&po.TalentPool{}).WithContext(ctx).Where("user_id = ? AND interview_status = ?", userID, "已面试").Find(&talentPOs).Error
	if err != nil {
		return nil, errorDB(err)
	}

	talents := make([]entity.Talent, len(talentPOs))
	for i, talentPO := range talentPOs {
		talents[i] = entity.Talent{
			TalentID:        talentPO.TalentID,
			FullName:        talentPO.FullName,
			TargetPosition:  talentPO.TargetPosition,
			MatchScore:      talentPO.MatchScore,
			CoreAdvantages:  talentPO.CoreAdvantages,
			HireStatus:      talentPO.HireStatus,
			InterviewStatus: talentPO.InterviewStatus,
			UserID:          talentPO.UserID,
			CreatedAt:       talentPO.CreatedAt,
		}
	}

	return talents, nil
}

func (r *ResumeStorage) UpdateTalent(ctx context.Context, newTalent *entity.Talent) error {

	updates := map[string]interface{}{
		"full_name":        newTalent.FullName,
		"target_position":  newTalent.TargetPosition,
		"match_score":      newTalent.MatchScore,
		"core_advantages":  newTalent.CoreAdvantages,
		"hire_status":      newTalent.HireStatus,
		"interview_status": newTalent.InterviewStatus,
	}

	err := r.db.Model(&po.TalentPool{}).WithContext(ctx).Where("talent_id = ? AND user_id = ?", newTalent.TalentID, newTalent.UserID).Updates(updates).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (r *ResumeStorage) DeleteTalent(ctx context.Context, talentID, userID string) error {
	if talentID == "" {
		return error_msg.TALENT_ID_NOT_NULL
	}

	err := r.db.Model(&po.TalentPool{}).WithContext(ctx).Where("talent_id = ? AND user_id = ?", talentID, userID).Delete(&po.TalentPool{}).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (r *ResumeStorage) CreateResumeAndTalent(ctx context.Context, resume *entity.Resume, talent *entity.Talent) error {
	if talent.TalentID == "" {
		return error_msg.TALENT_ID_NOT_NULL
	} else if talent.FullName == "" {
		return error_msg.TALENT_NAME_NOT_NULL
	} else if len(talent.FullName) >= 30 {
		return error_msg.TALENT_NAME_TOO_LONG
	} else if talent.TargetPosition == "" {
		return error_msg.TALENT_TARGET_POSITION_NOT_NULL
	} else if talent.MatchScore < 0 || talent.MatchScore > 100 {
		return error_msg.TALENT_MATCH_SCORE_INVALID
	} else if resume.ResumeName == "" {
		return error_msg.RESUME_NAME_NOT_NULL
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			zlog.Errorf("事务执行异常，已回滚")
		}
	}()

	var resumePO po.Resume
	var talentPO po.TalentPool

	resumePO.ResumeID = resume.ResumeID
	resumePO.ResumeName = resume.ResumeName
	resumePO.ResumeUrl = resume.ResumeUrl
	resumePO.UserID = resume.UserID
	resumePO.ResumeStr = resume.ResumeStr
	resumePO.CreatedAt = resume.CreatedAt

	talentPO.TalentID = talent.TalentID
	talentPO.FullName = talent.FullName
	talentPO.TargetPosition = talent.TargetPosition
	talentPO.MatchScore = talent.MatchScore
	talentPO.CoreAdvantages = talent.CoreAdvantages
	talentPO.HireStatus = talent.HireStatus
	talentPO.InterviewStatus = talent.InterviewStatus
	talentPO.UserID = talent.UserID
	talentPO.CreatedAt = talent.CreatedAt
	talentPO.ResumeID = talent.ResumeID

	err := tx.Model(&po.Resume{}).WithContext(ctx).Create(&resumePO).Error
	if err != nil {
		tx.Rollback()
		return errorDB(err)
	}

	err = tx.Model(&po.TalentPool{}).WithContext(ctx).Create(&talentPO).Error
	if err != nil {
		tx.Rollback()
		return errorDB(err)
	}

	err = tx.Commit().Error
	if err != nil {
		return errorDB(err)
	}
	return nil

}

func (r *ResumeStorage) DeleteResumeAndTalent(ctx context.Context, resumeID, talentID, userID string) error {
	if resumeID == "" {
		return error_msg.RESUME_ID_NOT_NULL
	} else if talentID == "" {
		return error_msg.TALENT_ID_NOT_NULL
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			zlog.Errorf("事务执行异常，已回滚")
		}
	}()

	err := tx.Model(&po.Resume{}).WithContext(ctx).Where("resume_id = ? AND user_id = ?", resumeID, userID).Delete(&po.Resume{}).Error
	if err != nil {
		tx.Rollback()
		return errorDB(err)
	}

	err = tx.Model(&po.TalentPool{}).WithContext(ctx).Where("talent_id = ? AND user_id = ?", talentID, userID).Delete(&po.TalentPool{}).Error
	if err != nil {
		tx.Rollback()
		return errorDB(err)
	}

	err = tx.Commit().Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

// 检查 talentID 是否存在
func (r *ResumeStorage) CheckTalentID(ctx context.Context, talentID string) (bool, error) {
	var talentPO po.TalentPool
	err := r.db.Model(&po.TalentPool{}).WithContext(ctx).Where("talent_id = ?", talentID).First(&talentPO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errorDB(err)
	}
	return true, nil
}
