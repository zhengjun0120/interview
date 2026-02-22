package resume_service

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/infra/cos"
	"ai_interview/pkg/error_msg"
	"ai_interview/util"
	"context"
	"math"
	"time"
)

type ResumeService struct {
	//注入需要的依赖

	resumeRepo repo.ResumeRepo
	chatRepo   repo.ChatRepo
}

func NewResumeService(resumeRepo repo.ResumeRepo, chatRepo repo.ChatRepo) *ResumeService {
	return &ResumeService{resumeRepo: resumeRepo, chatRepo: chatRepo}
}

// 上传简历
func (h *ResumeService) UploadResume(ctx context.Context, req *types.UploadResumeParams) (*types.UploadResumeResponse, error) {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return nil, error_msg.GET_USER_ID_ERROR
	}

	resumeID := util.GenerateStringID()

	resumeUrl, err := cos.UploadResume(req.File, req.FileHeader, resumeID)
	if err != nil {
		return nil, err
	}

	resume := entity.Resume{
		ResumeID:   resumeID,
		ResumeName: req.FileHeader.Filename,
		ResumeUrl:  resumeUrl,
		UserID:     userID,
	}

	err = h.resumeRepo.CreateResume(ctx, &resume)
	if err != nil {
		return nil, err
	}

	return &types.UploadResumeResponse{
		ResumeUrl: resumeUrl,
	}, nil

}

func (h *ResumeService) GetTalentAll(ctx context.Context) ([]types.GetTalentAllResponse, error) {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return nil, error_msg.GET_USER_ID_ERROR
	}

	talentList, err := h.resumeRepo.GetTalentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := make([]types.GetTalentAllResponse, len(talentList))

	for i, v := range talentList {
		resp[i] = types.GetTalentAllResponse{
			TalentID:        v.TalentID,
			FullName:        v.FullName,
			TargetPosition:  v.TargetPosition,
			MatchScore:      v.MatchScore,
			InterviewStatus: v.InterviewStatus,
			CoreAdvantages:  v.CoreAdvantages,
			HireStatus:      v.HireStatus,
			CreatedAt:       v.CreatedAt,
		}
	}

	return resp, nil
}

func (h *ResumeService) getInterviewTimeSeconds(startTime, endTime time.Time) int {
	//未结束
	if endTime.IsZero() {
		return 0
	}
	//时间异常 结束时间早于开始时间
	if endTime.Before(startTime) {
		return 0
	}

	duration := endTime.Sub(startTime)
	return int(math.Round(duration.Seconds()))
}

func (h *ResumeService) GetTalentInterview(ctx context.Context) ([]types.GetTalentInterviewResponse, error) {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return nil, error_msg.GET_USER_ID_ERROR
	}

	talents, err := h.resumeRepo.GetTalentInterviewByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	interviews, err := h.chatRepo.GetInterviewsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	timeMap := make(map[string]int)
	for _, v := range interviews {
		timeMap[v.TalentID] = h.getInterviewTimeSeconds(v.InterviewStartTime, v.InterviewEndTime)
	}

	resp := make([]types.GetTalentInterviewResponse, len(talents))
	for i, v := range talents {
		resp[i] = types.GetTalentInterviewResponse{
			TalentID:       v.TalentID,
			FullName:       v.FullName,
			TargetPosition: v.TargetPosition,
			CoreAdvantages: v.CoreAdvantages,
			HireStatus:     v.HireStatus,
			InterviewTime:  timeMap[v.TalentID],
			CreatedAt:      v.CreatedAt,
		}
	}

	return resp, nil

}
