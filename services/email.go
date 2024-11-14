package services

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
	"pichub.api/config"
	"pichub.api/models"
	"pichub.api/pkg/jwt"
)

type EmailServiceImpl struct {
	dialer *gomail.Dialer
	once   sync.Once
}

var EmailService = &EmailServiceImpl{}

type EmailRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}

// getDialer 获取邮件发送器
func (s *EmailServiceImpl) getDialer() *gomail.Dialer {
	s.once.Do(func() {
		s.dialer = gomail.NewDialer(
			config.Config.Email.Host,
			config.Config.Email.Port,
			config.Config.Email.Username,
			config.Config.Email.Password,
		)
	})
	return s.dialer
}

// 发送邮件
func (s *EmailServiceImpl) SendEmail(req *EmailRequest) error {

	// 创建邮件
	m := gomail.NewMessage()
	m.SetHeader("From", config.Config.Email.FromAddress)
	m.SetAddressHeader("From", config.Config.Email.FromAddress, config.Config.Email.FromName)
	m.SetHeader("To", req.To)
	m.SetHeader("Subject", req.Subject)

	m.SetBody("text/html", req.Body)

	// 使用 getDialer() 替代直接创建 dialer
	if err := s.getDialer().DialAndSend(m); err != nil {
		return fmt.Errorf("send email failed (host=%s, port=%d): %w",
			config.Config.Email.Host,
			config.Config.Email.Port,
			err,
		)
	}

	return nil
}

// 发送激活邮件
func (s *EmailServiceImpl) SendActivationEmail(user *models.User) error {
	// 生成激活token，有效期24小时
	claims := jwt.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.GenerateTokenWithClaims(claims)
	if err != nil {
		return fmt.Errorf("generate activation token failed: %w", err)
	}

	// 构建激活链接
	// TODO 暂时使用api地址，后续需要使用web地址
	activationLink := fmt.Sprintf("%s/api/v1/auth/activate?token=%s",
		config.Config.Server.GetFrontendUrl(),
		token,
	)

	// 邮件内容
	body := fmt.Sprintf(`
		<h2>欢迎加入 PicHub！</h2>
		<p>亲爱的 %s：</p>
		<p>请点击下面的链接激活您的账户：</p>
		<p><a href="%s">激活账户</a></p>
		<p>或者复制以下链接到浏览器打开：</p>
		<p>%s</p>
		<p>此链接24小时内有效。</p>
		<p>如果这不是您的操作，请忽略此邮件。</p>
		<br>
		<p>PicHub 团队</p>
	`, user.Nickname, activationLink, activationLink)

	return s.SendEmail(&EmailRequest{
		Subject: "激活您的账户",
		Body:    body,
		To:      user.Email,
	})
}
