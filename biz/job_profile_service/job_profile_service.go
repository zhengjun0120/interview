package job_profile_service

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/pkg/error_msg"
	"ai_interview/util"
	"context"
	"errors"
)

type JobProfileService struct {
	jobProfileRepo repo.JobProfileRepo
}

func NewJobProfileService(jobProfileRepo repo.JobProfileRepo) *JobProfileService {
	return &JobProfileService{jobProfileRepo: jobProfileRepo}
}

func (j *JobProfileService) SaveJobProfile(ctx context.Context, req *types.SaveJobProfileParams) error {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return error_msg.GET_USER_ID_ERROR
	}

	jobProfile := entity.JobProfile{
		UserID:                 userID,
		JobTitle:               req.JobTitle,
		RedLineCondition:       req.RedLineCondition,
		AiAdjustmentSuggestion: req.AiAdjustmentSuggestion,
		Competencies:           req.Competencies,
	}

	err := j.jobProfileRepo.CheckJobTitle(ctx, req.JobTitle)
	if err != nil {
		//如果职位已存在，则更新职位信息
		if errors.Is(err, error_msg.JOB_TITLE_ALREADY_EXISTS) {
			err := j.jobProfileRepo.UpdateJobProfile(ctx, &jobProfile)
			if err != nil {
				return err
			}
		}
		return err
	} else {
		//如果职位不存在
		jobProfile.JobProfileID = util.GenerateStringID()
		err := j.jobProfileRepo.CreateJobProfile(ctx, &jobProfile)
		if err != nil {
			return err
		}
	}

	return nil

}

func (j *JobProfileService) GetJobProfile(ctx context.Context) ([]types.GetJobProfile, error) {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return nil, error_msg.GET_USER_ID_ERROR
	}

	jobProfiles, err := j.jobProfileRepo.GetJobProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := make([]types.GetJobProfile, len(jobProfiles))
	for i, jobProfile := range jobProfiles {
		resp[i] = types.GetJobProfile{
			JobProfileID:           jobProfile.JobProfileID,
			JobTitle:               jobProfile.JobTitle,
			RedLineCondition:       jobProfile.RedLineCondition,
			AiAdjustmentSuggestion: jobProfile.AiAdjustmentSuggestion,
			Competencies:           jobProfile.Competencies,
			CreatedAt:              jobProfile.CreatedAt,
		}
	}
	return resp, nil

}
