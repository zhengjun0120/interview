package resume_service

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/infra/cos"
	"ai_interview/pkg/error_msg"
	"ai_interview/util"
	"context"
)

type ResumeService struct {
	//注入需要的依赖

	resumeRepo repo.ResumeRepo
}

func NewResumeService(resumeRepo repo.ResumeRepo) *ResumeService {
	return &ResumeService{resumeRepo: resumeRepo}
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
