package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	UserType int    `json:"user_type"`
	jwt.RegisteredClaims
}

// 添加一些辅助函数
func NewNumericDate(t time.Time) *jwt.NumericDate {
	return jwt.NewNumericDate(t)
}

type RegisteredClaims = jwt.RegisteredClaims

// GenerateToken 生成JWT token
func GenerateToken(userID int, username string, userType int) (string, error) {
	// 从配置中获取 JWT 密钥和过期时间
	secret := []byte(viper.GetString("JWT_SECRET"))
	expireStr := viper.GetString("JWT_EXPIRE")

	// 解析过期时间
	expireDuration, err := time.ParseDuration(expireStr)
	if err != nil {
		expireDuration = 30 * 24 * time.Hour // 默认30天
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	secret := []byte(viper.GetString("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateTokenWithClaims 使用自定义Claims生成JWT token
func GenerateTokenWithClaims(claims Claims) (string, error) {
	secret := []byte(viper.GetString("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
