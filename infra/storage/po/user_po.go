package po

import "time"

type User struct {
	ID int `gorm:"column:id;primary_key;autoIncrement"`

	UserID   string `gorm:"column:user_id;not null;unique"`
	Username string `gorm:"column:username;not null"`
	Email    string `gorm:"column:email;not null;"`
	Password string `gorm:"column:password;not null"`

	Type string `gorm:"column:type;not null"`

	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}
