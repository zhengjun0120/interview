package storage

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/infra/database"
	"ai_interview/infra/storage/po"
	"ai_interview/pkg/error_msg"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JobProfileStorage struct {
	db *gorm.DB
}

var jps *JobProfileStorage

func InitJobProfileStorage() {
	db := database.GetDB()

	if err := db.AutoMigrate(&po.JobProfile{}); err != nil {
		panic("job_profile表自动迁移失败: " + err.Error())
	}
	jps = &JobProfileStorage{db}
}

func GetJobProfileStorage() repo.JobProfileRepo {
	return jps
}

func (j *JobProfileStorage) CreateJobProfile(ctx context.Context, jobProfile *entity.JobProfile) error {
	if jobProfile.JobProfileID == "" {
		return error_msg.JOB_PROFILE_ID_NOT_NULL
	} else if jobProfile.JobTitle == "" {
		return error_msg.JOB_TITLE_NOT_NULL
	} else if len(jobProfile.JobTitle) > 20 {
		return error_msg.JOB_TITLE_TOO_LONG
	}

	var competencyWeight int = 50

	competencies := []entity.CompetenciesJson{
		{Name: "技术深度", Type: "专业技能", Weight: competencyWeight},
		{Name: "尽责性", Type: "五大人格", Weight: competencyWeight},
		{Name: "开放性", Type: "五大人格", Weight: competencyWeight},
		{Name: "情绪稳定", Type: "五大人格", Weight: competencyWeight},
		{Name: "团队协助", Type: "社会特质", Weight: competencyWeight},
		{Name: "业务敏感", Type: "胜任力", Weight: competencyWeight},
	}

	competenciesJsonData, err := json.Marshal(competencies)
	if err != nil {
		return fmt.Errorf("序列化 competencies 失败: %v", err)
	}

	var jobProfilePO po.JobProfile
	jobProfilePO.JobProfileID = jobProfile.JobProfileID
	jobProfilePO.JobTitle = jobProfile.JobTitle
	jobProfilePO.UserID = jobProfile.UserID
	jobProfilePO.RedLineCondition = datatypes.JSON("[]")
	jobProfilePO.Competencies = datatypes.JSON(competenciesJsonData)

	err = j.db.Model(&po.JobProfile{}).WithContext(ctx).Create(&jobProfilePO).Error
	if err != nil {
		return errorDB(err)
	}

	return nil
}

func (j *JobProfileStorage) DeleteJobProfile(ctx context.Context, jobProfileID string) error {
	if jobProfileID == "" {
		return error_msg.JOB_PROFILE_ID_NOT_NULL
	}

	err := j.db.Model(&po.JobProfile{}).WithContext(ctx).Where("job_profile_id = ?", jobProfileID).Delete(&po.JobProfile{}).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (j *JobProfileStorage) UpdateJobProfile(ctx context.Context, jobProfile *entity.JobProfile) error {
	if jobProfile.JobTitle == "" {
		return error_msg.JOB_TITLE_NOT_NULL
	} else if len(jobProfile.JobTitle) > 20 {
		return error_msg.JOB_TITLE_TOO_LONG
	} else if len(jobProfile.Competencies) == 0 {
		return error_msg.JOB_PROFILE_COMPETENCIES_NOT_NULL
	}

	updates := map[string]interface{}{
		"job_title":          jobProfile.JobTitle,
		"competencies":       jobProfile.Competencies,
		"red_line_condition": jobProfile.RedLineCondition,
	}

	err := j.db.Model(&po.JobProfile{}).WithContext(ctx).Where("job_profile_id = ?", jobProfile.JobProfileID).Updates(updates).Error
	if err != nil {
		return errorDB(err)
	}
	return nil

}

func (j *JobProfileStorage) GetJobProfileByUserID(ctx context.Context, userID string) ([]entity.JobProfile, error) {

	var jobProfilePOs []po.JobProfile
	err := j.db.Model(&po.JobProfile{}).WithContext(ctx).Where("user_id = ?", userID).Find(&jobProfilePOs).Error
	if err != nil {
		return nil, errorDB(err)
	}

	resp := make([]entity.JobProfile, len(jobProfilePOs))

	for i, v := range jobProfilePOs {
		resp[i] = entity.JobProfile{
			JobProfileID: v.JobProfileID,
			JobTitle:     v.JobTitle,
			UserID:       v.UserID,
		}

		err := json.Unmarshal(v.Competencies, &resp[i].Competencies)
		if err != nil {
			return nil, fmt.Errorf("反序列化 competencies 失败: %v", err)
		}

		err = json.Unmarshal(v.RedLineCondition, &resp[i].RedLineCondition)
		if err != nil {
			return nil, fmt.Errorf("反序列化 red_line_condition 失败: %v", err)
		}
	}

	return resp, nil
}

func (j *JobProfileStorage) GetJobProfileByJobProfileID(ctx context.Context, jobProfileID string) (*entity.JobProfile, error) {
	if jobProfileID == "" {
		return nil, error_msg.JOB_PROFILE_ID_NOT_NULL
	}

	var jobProfilePO po.JobProfile
	err := j.db.Model(&po.JobProfile{}).WithContext(ctx).Where("job_profile_id = ?", jobProfileID).First(&jobProfilePO).Error
	if err != nil {
		return nil, errorDB(err)
	}

	var jobProfile entity.JobProfile
	jobProfile.JobProfileID = jobProfilePO.JobProfileID
	jobProfile.JobTitle = jobProfilePO.JobTitle
	jobProfile.UserID = jobProfilePO.UserID

	err = json.Unmarshal(jobProfilePO.Competencies, &jobProfile.Competencies)
	if err != nil {
		return nil, fmt.Errorf("反序列化 competencies 失败: %v", err)
	}
	err = json.Unmarshal(jobProfilePO.RedLineCondition, &jobProfile.RedLineCondition)
	if err != nil {
		return nil, fmt.Errorf("反序列化 red_line_condition 失败: %v", err)
	}
	return &jobProfile, nil
}
