package email

import (
	"fmt"
	"net/smtp"
	"strings"

	pkgo "github.com/rebooe/pkg-go"
)

type Email struct {
	smtpHost string    // 邮件服务地址
	from     string    // 发送放账户
	fromName string    // 发送方名称
	auth     smtp.Auth // 邮件授权信息
}

func New(smtpHost string, from, pwd, fromName string) *Email {
	host, _, _ := strings.Cut(smtpHost, ":")
	auth := smtp.PlainAuth("", from, pwd, host)

	return &Email{
		smtpHost: smtpHost,
		from:     from,
		fromName: fromName,
		auth:     auth,
	}
}

func (em *Email) SendMail(to []string, subject string, body string) error {
	message := fmt.Sprintf("From: %s <%s>\r\nSubject: %s\r\n\r\n%s",
		em.fromName, em.from, subject, body)
	return smtp.SendMail(em.smtpHost, em.auth, em.from, to, pkgo.StringToBytes(message))
}
