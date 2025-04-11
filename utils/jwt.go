package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("your-secret-key") // 在实际应用中应该从环境变量获取

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, role string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour) // token有效期24小时

	claims := Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 解析JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
