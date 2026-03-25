package sms

import (
	"ai_interview/conf"
	"ai_interview/pkg/error_msg"
	"crypto/tls"
	"fmt"

	"gopkg.in/mail.v2"
)

type SmsService struct {
}

func SendCaptcha(email string, code string) error {
	m := mail.NewMessage()
	m.SetHeader("From", conf.GetConfig().Smtp.SmtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", conf.GetConfig().Smtp.EncodedName)
	m.SetBody("text/html", fmt.Sprintf(`
		<div style="font-size:14px;">
			您好！您的验证码是：<span style="color:red; font-weight:bold;">%s</span>
			<br>
			该验证码 5 分钟内有效，请及时使用，请勿泄露给他人。
		</div>
	`, code))
	d := mail.NewDialer(conf.GetConfig().Smtp.SmtpHost, conf.GetConfig().Smtp.SmtpPort, conf.GetConfig().Smtp.SmtpUser, conf.GetConfig().Smtp.SmtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return error_msg.ErrorCaptchaSend
	}
	return nil
}
