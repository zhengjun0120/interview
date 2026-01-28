package po

import "time"

type Resume struct {
	ID int `gorm:"column:id;primary_key;autoIncrement"`

	UserID     string `gorm:"column:user_id;not null"`
	ResumeID   string `gorm:"column:resume_id;not null;unique"`
	ResumeName string `gorm:"column:resume_name;not null"`
	ResumeUrl  string `gorm:"column:resume_url;not null"`

	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}
