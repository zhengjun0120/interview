package entity

import "time"

type Resume struct {
	UserID     string
	ResumeID   string
	ResumeName string
	ResumeUrl  string
	CreatedAt  time.Time
}

type Talent struct {
	UserID          string    //所属用户ID
	ResumeID        string    //所属简历ID
	TalentID        string    //人才ID
	FullName        string    //人才姓名
	TargetPosition  string    //目标职位
	MatchScore      int       //匹配度
	InterviewStatus string    //面试状态
	CoreAdvantages  string    //核心优势
	HireStatus      string    //录用状态
	CreatedAt       time.Time //创建时间
}

type TalentJson struct {
	FullName       string `json:"full_name"`
	TargetPosition string `json:"target_position"`
	MatchScore     int    `json:"match_score"`
	CoreAdvantages string `json:"core_advantages"`
}
