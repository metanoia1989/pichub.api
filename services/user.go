package services

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"pichub.api/infra/database"
	"pichub.api/models"
	"pichub.api/utils/jwt"
)

type userService struct{}

var UserService = new(userService)

func (s *userService) Register(req models.RegisterRequest) (*models.User, error) {
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
		Username:     req.Username,
		Nickname:     req.Nickname,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		UserType:     1, // 普通用户类型
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(req models.LoginRequest) (string, *models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return "", nil, errors.New("user not found")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *userService) ActivateAccount(token string) error {
	// TODO: 实现账户激活逻辑
	// 1. 验证激活token
	// 2. 更新用户状态为已激活
	return nil
}

func (s *userService) HasGithubToken(userID int) bool {
	// TODO: 检查用户是否已添加 GitHub token
	return false
}
