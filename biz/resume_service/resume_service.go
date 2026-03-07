package resume_service

import (
	"ai_interview/biz/ai_chat"
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/biz/resume_service/cos_service"
	"ai_interview/biz/types"
	"ai_interview/pkg/error_msg"
	"ai_interview/util"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"math"
	"strconv"
	"strings"
	"time"
)

type ResumeService struct {
	//注入需要的依赖

	resumeRepo     repo.ResumeRepo
	chatRepo       repo.ChatRepo
	jobProfileRepo repo.JobProfileRepo
	chatService    ai_chat.IChatService
	cosService     cos_service.ICosService
}

func NewResumeService(resumeRepo repo.ResumeRepo, chatRepo repo.ChatRepo, jobProfileRepo repo.JobProfileRepo, chatService ai_chat.IChatService, cosService cos_service.ICosService) *ResumeService {
	return &ResumeService{resumeRepo: resumeRepo, chatRepo: chatRepo, jobProfileRepo: jobProfileRepo, chatService: chatService, cosService: cosService}
}

func (h *ResumeService) jobProfileEntity2String(jobProfiles []entity.JobProfile) string {
	if len(jobProfiles) <= 0 {
		return ""
	}

	var builder strings.Builder

	builder.WriteString("以下是公司岗位画像库，匹配简历时需：1. 先匹配岗位名称；2. 红线条件不满足则匹配度直接为0；3. 核心指标按权重加权评分\n[岗位画像库]\n")

	for i, v := range jobProfiles {
		builder.WriteString("=== 岗位")
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(" ===\n")
		builder.WriteString("岗位名称：")
		builder.WriteString(v.JobTitle)
		builder.WriteString("\n")

		if len(v.RedLineCondition) > 0 {
			builder.WriteString("红线条件（必须满足）：\n")
			for condIdx, condition := range v.RedLineCondition {
				builder.WriteString(strconv.Itoa(condIdx))
				builder.WriteString(". ")
				builder.WriteString(condition)
				builder.WriteString("\n")
			}
		} else {
			builder.WriteString("红线条件（必须满足）：无\n")
		}

		if len(v.Competencies) > 0 {
			builder.WriteString("核心匹配指标（含权重）：\n")
			for compIdx, competency := range v.Competencies {
				builder.WriteString(strconv.Itoa(compIdx))
				builder.WriteString(". 指标名称：")
				builder.WriteString(competency.Name)
				builder.WriteString(" - 指标类型：")
				builder.WriteString(competency.Type)
				builder.WriteString(" - 权重：")
				builder.WriteString(strconv.Itoa(competency.Weight))
				builder.WriteString("\n")
			}
		} else {
			builder.WriteString("核心匹配指标（含权重）：无\n")
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// 上传简历
func (h *ResumeService) UploadResume(ctx context.Context, req *types.UploadResumeParams) (*types.UploadResumeResponse, error) {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return nil, error_msg.GET_USER_ID_ERROR
	}

	resumeID := util.GenerateStringID()

	// 上传简历到cos
	resumeUrl, err := h.cosService.UploadResume(req.File, req.FileHeader, resumeID)
	if err != nil {
		return nil, err
	}

	// 保存简历信息
	resume := entity.Resume{
		ResumeID:   resumeID,
		ResumeName: req.FileHeader.Filename,
		ResumeUrl:  resumeUrl,
		UserID:     userID,
		CreatedAt:  time.Now(),
	}

	// 获取用户岗位画像
	jobProfiles, err := h.jobProfileRepo.GetJobProfileByUserID(ctx, userID)
	if err != nil {
		h.cosService.DeleteFile(ctx, resumeUrl)
		return nil, err
	}

	// 将岗位画像转换为字符串
	jobProfileStr := h.jobProfileEntity2String(jobProfiles)

	// 解析简历核心信息
	messageStr, err := h.chatService.DocxToTalentDataChat(ctx, resumeUrl, jobProfileStr)
	if err != nil {
		h.cosService.DeleteFile(ctx, resumeUrl)
		return nil, err
	}
	var talentJson entity.TalentJson
	err = json.Unmarshal([]byte(messageStr), &talentJson)
	if err != nil {
		h.cosService.DeleteFile(ctx, resumeUrl)
		return nil, fmt.Errorf("解析ai消息失败: %v\n 原始消息: %s", err, messageStr)
	}

	// 获取简历文本内容
	resumeStr, err := h.chatService.DocxToResumeStrChat(ctx, resumeUrl)
	if err != nil {
		h.cosService.DeleteFile(ctx, resumeUrl)
		return nil, err
	}
	// 保存简历文本内容
	resume.ResumeStr = resumeStr
	// 创建简历和人才信息
	var talentEntity = entity.Talent{
		UserID:          userID,
		ResumeID:        resumeID,
		TalentID:        util.GenerateStringID(),
		FullName:        talentJson.FullName,
		TargetPosition:  talentJson.TargetPosition,
		MatchScore:      talentJson.MatchScore,
		InterviewStatus: entity.InterviewStatusUninterviewed, //未面试
		CoreAdvantages:  talentJson.CoreAdvantages,
		HireStatus:      entity.HireStatusNotHired, //未录用
		CreatedAt:       time.Now(),
	}

	// 保存简历和人才信息
	err = h.resumeRepo.CreateResumeAndTalent(ctx, &resume, &talentEntity)
	if err != nil {
		h.cosService.DeleteFile(ctx, resumeUrl)
		return nil, err
	}

	return &types.UploadResumeResponse{
		TalentID:        talentEntity.TalentID,
		FullName:        talentEntity.FullName,
		TargetPosition:  talentEntity.TargetPosition,
		MatchScore:      talentEntity.MatchScore,
		InterviewStatus: talentEntity.InterviewStatus,
		CoreAdvantages:  talentEntity.CoreAdvantages,
		HireStatus:      talentEntity.HireStatus,
		CreatedAt:       talentEntity.CreatedAt,
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
		if v.InterviewEndTime == nil {
			timeMap[v.TalentID] = 0
		} else {
			timeMap[v.TalentID] = h.getInterviewTimeSeconds(v.InterviewStartTime, *v.InterviewEndTime)
		}

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

func (h *ResumeService) GetResumeUrl(ctx context.Context, req *types.GetResumeUrlRequest) (*types.GetResumeUrlResponse, error) {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return nil, error_msg.GET_USER_ID_ERROR
	}

	resume, err := h.resumeRepo.GetResumeByTalentID(ctx, req.TalentID)
	if err != nil {
		return nil, err
	}
	if resume.UserID != userID {
		return nil, error_msg.RESUME_NOT_EXIST
	}

	return &types.GetResumeUrlResponse{
		ResumeUrl: resume.ResumeUrl,
	}, nil
}

// 获取面试报告
func (h *ResumeService) GetTalentReport(ctx context.Context, req *types.GetTalentReportRequest) (*types.GetTalentReportResponse, error) {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return nil, error_msg.GET_USER_ID_ERROR
	}

	interview, err := h.chatRepo.GetInterviewByTalentID(ctx, req.TalentID, userID)
	if err != nil {
		return nil, err
	}

	if interview.TalentReport != nil {
		var getTalentReportQuestionJson []types.GetTalentReportQuestionArrJson
		err := json.Unmarshal(interview.TalentReport, &getTalentReportQuestionJson)
		if err != nil {
			return nil, fmt.Errorf("解析ai消息失败: %v\n 原始消息: %s", err, interview.TalentReport)
		}
		return &types.GetTalentReportResponse{
			List: getTalentReportQuestionJson,
		}, nil
	}

	message, err := h.chatRepo.GetInterviewMessageByRoomID(ctx, interview.RoomID)
	if err != nil {
		return nil, err
	}

	aiResp, err := h.chatService.InterviewMessageAnalyseChat(ctx, message)
	if err != nil {
		return nil, err
	}

	var getTalentReportQuestionJson []types.GetTalentReportQuestionArrJson
	err = json.Unmarshal([]byte(aiResp), &getTalentReportQuestionJson)
	if err != nil {
		return nil, fmt.Errorf("解析ai消息失败: %v\n 原始消息: %s", err, aiResp)
	}

	err = h.chatRepo.SaveTalentReport(ctx, datatypes.JSON(aiResp), userID, req.TalentID)
	if err != nil {
		return nil, err
	}

	return &types.GetTalentReportResponse{
		List: getTalentReportQuestionJson,
	}, nil
}

func (h *ResumeService) MarkTalentHireStatus(ctx context.Context, req *types.MarkTalentHireStatusRequest) error {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return error_msg.GET_USER_ID_ERROR
	}
	err := h.resumeRepo.MarkTalentHireStatus(ctx, req.TalentID, userID, req.HireStatus)
	if err != nil {
		return err
	}
	return nil
}
