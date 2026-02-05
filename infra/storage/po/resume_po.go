package po

import "time"

type Resume struct {
	ID int `gorm:"column:id;primary_key;autoIncrement"`

	UserID     string `gorm:"column:user_id;not null"`          //所属用户ID
	ResumeID   string `gorm:"column:resume_id;not null;unique"` //简历ID
	ResumeName string `gorm:"column:resume_name;not null"`      //简历名称
	ResumeUrl  string `gorm:"column:resume_url;not null"`       //简历URL

	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}

type TalentPool struct {
	ID int `gorm:"column:id;primary_key;autoIncrement"`

	UserID          string    `gorm:"column:user_id;not null"`          //所属用户ID
	ResumeID        string    `gorm:"column:resume_id;not null;unique"` //所属简历ID
	TalentID        string    `gorm:"column:talent_id;not null;unique"` //人才ID
	FullName        string    `gorm:"column:full_name;not null"`        //人才姓名
	TargetPosition  string    `gorm:"column:target_position;not null"`  //目标职位
	MatchScore      int       `gorm:"column:match_score;not null"`      //匹配度
	InterviewStatus string    `gorm:"column:interview_status;not null"` //面试状态
	CoreAdvantages  string    `gorm:"column:core_advantages;not null"`  //核心优势
	HireStatus      string    `gorm:"column:hire_status;not null"`      //录用状态
	CreatedAt       time.Time `gorm:"column:created_at;"`
	UpdatedAt       time.Time `gorm:"column:updated_at;"`
}
