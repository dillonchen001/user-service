package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Generator struct {
	secret  string
	expires int // 小时
}

func NewGenerator(secret string, expires int) *Generator {
	return &Generator{
		secret:  secret,
		expires: expires,
	}
}

func (g *Generator) GenerateToken(userID int64) (string, error) {
	// 创建claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(g.expires)).Unix(),
		"iat":     time.Now().Unix(),
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串token
	tokenString, err := token.SignedString([]byte(g.secret))

	return tokenString, err
}
