package sms

import "context"

type SMSService interface {
	SendCaptcha(c context.Context, code string, email string) error
}
