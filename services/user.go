package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"pichub.api/infra/database"
	"pichub.api/models"
	"pichub.api/pkg/jwt"
	"pichub.api/pkg/utils"
)

type UserServiceImpl struct{}

var UserService = &UserServiceImpl{}

// 用户状态常量
const (
	UserStatusInactive = 0
	UserStatusActive   = 1

	UserTypeNormal = 0
	UserTypeAdmin  = 1

	EmailVerificationCodeExpiration = 10 * time.Minute
	EmailVerificationCodePrefix     = "email_verification:"
)

func (s *UserServiceImpl) Register(req models.RegisterRequest) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if result := database.DB.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if result := database.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := &models.User{
		UserBaseInfo: models.UserBaseInfo{
			Username:  req.Username,
			Nickname:  req.Nickname,
			Email:     req.Email,
			UserType:  UserTypeNormal,     // 普通用户类型
			Status:    UserStatusInactive, // 设置初始状态为未激活
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PasswordHash: string(hashedPassword),
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) Login(req models.LoginRequest) (string, *models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ? or email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		return "", nil, errors.New("user not found")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	// 是否激活
	if user.Status == UserStatusInactive {
		return "", nil, errors.New("account not activated")
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *UserServiceImpl) ActivateAccount(token string) error {
	// 解析并验证 token
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return errors.New("invalid or expired activation token")
	}

	// 获取用户
	var user models.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		return errors.New("user not found")
	}

	// 检查用户是否已激活
	if user.Status == UserStatusActive {
		return errors.New("account already activated")
	}

	// 更新用户状态为已激活
	user.Status = UserStatusActive
	if err := database.DB.Model(&user).Update("status", UserStatusActive).Error; err != nil {
		return fmt.Errorf("failed to activate account: %w", err)
	}

	return nil
}

func (s *UserServiceImpl) HasGithubToken(userID int) (bool, error) {
	githubToken, err := ConfigService.Get("github", "token", userID)
	if err != nil {
		return false, err
	}

	// 验证 github token 是否有效
	if utils.IsEmpty(githubToken) {
		return false, errors.New("github token is not set")
	}

	// 验证 token 有效性
	valid, err := GithubService.ValidateToken(githubToken.(string))
	if err != nil {
		return false, fmt.Errorf("invalid github token: %v", err)
	}

	return valid, nil
}

func (s *UserServiceImpl) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserServiceImpl) UpdateProfile(userID int, req models.UpdateProfileRequest) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// 如果要更新用户名，检查是否存在重复
	if req.Username != "" && req.Username != user.Username {
		var existingUser models.User
		if result := database.DB.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser); result.Error == nil {
			return errors.New("username already exists")
		}
		user.Username = req.Username
	}

	// 更新昵称
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	// 更新密码
	if req.NewPassword != "" {
		// 如果提供了新密码，验证旧密码
		if req.OldPassword == "" {
			return errors.New("old password is required when updating password")
		}

		// 验证旧密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
			return errors.New("invalid old password")
		}

		// 加密新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}

	// 保存更新
	return database.DB.Save(&user).Error
}

func (s *UserServiceImpl) UpdateGithubToken(userID int, token string) error {
	// 先验证 token 是否有效
	valid, err := GithubService.ValidateToken(token)
	if err != nil || !valid {
		return fmt.Errorf("invalid github token: %v", err)
	}

	// 保存 token 到配置
	if err := ConfigService.Set("github", "token", token, userID); err != nil {
		return fmt.Errorf("failed to save github token: %v", err)
	}

	return nil
}

// SendEmailVerificationCode 发送邮箱验证码
func (s *UserServiceImpl) SendEmailVerificationCode(userID int, newEmail string) error {
	// 检查新邮箱是否已被使用
	var existingUser models.User
	if result := database.DB.Where("email = ? AND id != ?", newEmail, userID).First(&existingUser); result.Error == nil {
		return errors.New("email already exists")
	}

	// 生成6位随机验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 将验证码保存到 Redis，设置 10 分钟过期
	key := fmt.Sprintf("%s%d:%s", EmailVerificationCodePrefix, userID, newEmail)
	if err := RedisService.Set(context.Background(), key, code, EmailVerificationCodeExpiration).Err(); err != nil {
		return fmt.Errorf("failed to save verification code: %v", err)
	}

	// 发送验证码邮件
	body := fmt.Sprintf(`
		<h2>邮箱验证</h2>
		<p>您正在更换邮箱地址，验证码为：</p>
		<h3>%s</h3>
		<p>验证码有效期为10分钟。</p>
		<p>如果这不是您的操作，请忽略此邮件。</p>
		<br>
		<p>PicHub 团队</p>
	`, code)

	if err := EmailService.SendEmail(&EmailRequest{
		Subject: "更换邮箱验证码",
		Body:    body,
		To:      newEmail,
	}); err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	return nil
}

// UpdateEmail 验证码确认并更新邮箱
func (s *UserServiceImpl) UpdateEmail(userID int, req models.UpdateEmailRequest) error {
	// 获取存储的验证码
	key := fmt.Sprintf("%s%d:%s", EmailVerificationCodePrefix, userID, req.NewEmail)
	storedCode, err := RedisService.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("verification code expired or not found")
		}
		return fmt.Errorf("failed to get verification code: %v", err)
	}

	// 验证码匹配检查
	if storedCode != req.Code {
		return errors.New("invalid verification code")
	}

	// 更新邮箱
	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("email", req.NewEmail).Error; err != nil {
		return fmt.Errorf("failed to update email: %v", err)
	}

	// 删除验证码
	RedisService.Del(context.Background(), key)

	return nil
}
