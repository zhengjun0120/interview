package po

import (
	"gorm.io/datatypes"
	"time"
)

// JobProfile 岗位画像
type JobProfile struct {
	ID                     int            `gorm:"column:id;primary_key;autoIncrement"`
	JobProfileID           string         `gorm:"column:job_profile_id;not null;unique"` //岗位画像ID
	UserID                 string         `gorm:"column:user_id;not null"`               //所属用户ID
	JobTitle               string         `gorm:"column:job_title;not null;unique"`      //岗位名称
	Competencies           datatypes.JSON `gorm:"column:competencies;not null"`          //岗位能力要求
	RedLineCondition       datatypes.JSON `gorm:"column:red_line_condition;not null"`    //岗位红线要求(自动淘汰规则)
	AiAdjustmentSuggestion string         `gorm:"column:ai_adjustment_suggestion"`       //岗位调优建议

	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}
