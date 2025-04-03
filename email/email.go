package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/rebooe/pkgo"
)

type Email struct {
	smtpHost string    // 邮件服务地址
	from     string    // 发件人邮箱
	fromName string    // 发送人名称
	auth     smtp.Auth // 邮件授权信息
}

type EmailConfig struct {
	SmtpHost string `yaml:"SmtpHost"` // 邮件服务器地址
	From     string `yaml:"From"`     // 发件人邮箱
	Password string `yaml:"Password"` // 授权码/密码
	FromName string `yaml:"FromName"` // 发送人名称
}

func New(config EmailConfig) *Email {
	host, _, _ := strings.Cut(config.SmtpHost, ":")
	auth := smtp.PlainAuth("", config.From, config.Password, host)

	return &Email{
		smtpHost: config.SmtpHost,
		from:     config.From,
		fromName: config.FromName,
		auth:     auth,
	}
}

// SendMail 发送邮件
//
//	to: 收件人邮箱列表
//	subject: 邮件主题
//	body: 邮件内容
func (em *Email) SendMail(to []string, subject string, body string) error {
	message := fmt.Sprintf("From: %s <%s>\r\nSubject: %s\r\n\r\n%s",
		em.fromName, em.from, subject, body)
	return smtp.SendMail(em.smtpHost, em.auth, em.from, to, pkgo.StringToBytes(message))
}
