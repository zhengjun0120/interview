package util

import (
	"ai_interview/pkg/error_msg"
	"regexp"
)

// 检查邮箱格式
func CheckEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	ok := re.MatchString(email)

	if !ok {
		return error_msg.Email_FORMAT_INVALID
	}
	return nil
}
