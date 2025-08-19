package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	v1 "user-service/api/auth/v1"
	"user-service/internal/biz"
	"user-service/internal/conf"
	"user-service/third_party/jwt"

	"github.com/go-kratos/kratos/v2/log"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

type GoogleService struct {
	cfg          *conf.Jwt
	log          *log.Helper
	userAuthCase *biz.UserAuthCase
	userCase     *biz.UserCase
	jwtGen       *jwt.Generator
	httpClient   *http.Client
}

func NewGoogleService(cfg *conf.Jwt, logger log.Logger, userAuthCase *biz.UserAuthCase, userCase *biz.UserCase) *GoogleService {
	return &GoogleService{
		cfg:          cfg,
		log:          log.NewHelper(logger),
		userAuthCase: userAuthCase,
		userCase:     userCase,
		jwtGen:       jwt.NewGenerator(cfg.Secret, int(cfg.Expires)),
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// GoogleClaims JWT 声明结构
type GoogleClaims struct {
	jwtv4.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"` // Google 用户唯一标识符
}

func (s *GoogleService) Login(ctx context.Context, req *v1.LoginWithGoogleRequest) (*v1.LoginResponse, error) {
	// 验证参数
	if req.IdToken == "" {
		return nil, errors.New("id_token is required")
	}

	// 验证 id_token
	claims, err := s.verifyIdToken(ctx, req.IdToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify google id_token: %w", err)
	}

	// 检查邮箱是否已验证
	if !claims.EmailVerified {
		return nil, errors.New("google email is not verified")
	}

	// 查找或创建用户
	u, isNew, err := s.userAuthCase.FindOrCreateByGoogleID(ctx, claims.Sub, claims.Name, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// 更新用户头像
	if claims.Picture != "" && u.Avatar != claims.Picture {
		u.Avatar = claims.Picture
		if _, err = s.userCase.Update(ctx, u); err != nil {
			s.log.Errorf("failed to update user avatar, error: %v", err)
		}
	}

	// 生成 JWT token
	token, err := s.jwtGen.GenerateToken(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 构建响应
	return &v1.LoginResponse{
		Token:     token,
		IsNewUser: isNew,
		UserInfo: &v1.UserInfo{
			Id:     u.ID,
			Name:   u.Name,
			Avatar: u.Avatar,
			Email:  u.Email,
		},
	}, nil
}

func (s *GoogleService) verifyIdToken(ctx context.Context, idToken string) (*GoogleClaims, error) {
	// 调用 Google API 验证 id_token
	// 注意：实际应用中可以使用 jwt.Parse 直接验证，并通过 Google 的公钥验证签名
	// 这里为了简化，使用 Google 的验证 API
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			s.log.Errorf("close io reader failed, error: %v", errClose)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google tokeninfo api returned status code: %d", resp.StatusCode)
	}

	var tokenInfo struct {
		Iss           string `json:"iss"`
		Sub           string `json:"sub"`
		Azp           string `json:"azp"`
		Aud           string `json:"aud"`
		Iat           int64  `json:"iat"`
		Exp           int64  `json:"exp"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, err
	}

	// 验证 issuer
	if tokenInfo.Iss != "https://accounts.google.com" && tokenInfo.Iss != "accounts.google.com" {
		return nil, errors.New("invalid issuer")
	}

	// 验证 audience
	// 应该验证 audience 是否为您的应用客户端 ID

	// 验证过期时间
	now := time.Now().Unix()
	if tokenInfo.Exp < now {
		return nil, errors.New("token has expired")
	}

	// 构建 claims
	claims := &GoogleClaims{
		RegisteredClaims: jwtv4.RegisteredClaims{
			Issuer:    tokenInfo.Iss,
			Subject:   tokenInfo.Sub,
			ExpiresAt: jwtv4.NewNumericDate(time.Unix(tokenInfo.Exp, 0)),
			IssuedAt:  jwtv4.NewNumericDate(time.Unix(tokenInfo.Iat, 0)),
		},
		Email:         tokenInfo.Email,
		EmailVerified: tokenInfo.EmailVerified,
		Name:          tokenInfo.Name,
		Picture:       tokenInfo.Picture,
		Sub:           tokenInfo.Sub,
	}

	return claims, nil
}
