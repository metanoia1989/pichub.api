package services

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"pichub.api/config"
	"pichub.api/models"
	"pichub.api/utils/jwt"
)

type EmailServiceImpl struct {
	dialer *gomail.Dialer
}

var EmailService = &EmailServiceImpl{
	dialer: gomail.NewDialer(
		config.Config.Email.Host,
		config.Config.Email.Port,
		config.Config.Email.Username,
		config.Config.Email.Password,
	),
}

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
	activationLink := fmt.Sprintf("%s/activate?token=%s",
		viper.GetString("FRONTEND_URL"),
		token,
	)

	// 创建邮件
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("SMTP_FROM"))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "激活您的账户")

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

	m.SetBody("text/html", body)

	// 发送邮件
	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("send email failed: %w", err)
	}

	return nil
}
