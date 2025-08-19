package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	v1 "user-service/api/auth/v1"
	"user-service/internal/biz"
	"user-service/internal/conf"
	"user-service/third_party/jwt"
)

type FacebookService struct {
	cfg          *conf.Jwt
	log          *log.Helper
	userAuthCase *biz.UserAuthCase
	userCase     *biz.UserCase
	jwt          *jwt.Generator
	httpClient   *http.Client
}

func NewFacebookService(cfg *conf.Jwt, logger log.Logger, userAuthCase *biz.UserAuthCase, userCase *biz.UserCase) *FacebookService {
	return &FacebookService{
		cfg:          cfg,
		log:          log.NewHelper(logger),
		userAuthCase: userAuthCase,
		userCase:     userCase,
		jwt:          jwt.NewGenerator(cfg.Secret, int(cfg.Expires)),
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// FacebookUserResponse 用户信息响应结构
type FacebookUserResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture struct {
		Data struct {
			Url string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

func (s *FacebookService) Login(ctx context.Context, req *v1.LoginWithFacebookRequest) (*v1.LoginResponse, error) {
	// 验证 access token
	if req.AccessToken == "" {
		return nil, errors.New("access token is required")
	}

	// 调用 Facebook Graph API 获取用户信息
	userInfo, err := s.getUserInfo(ctx, req.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get facebook user info: %w", err)
	}

	// 查找或创建用户
	u, isNew, err := s.userAuthCase.FindOrCreateByFacebookID(ctx, userInfo.ID, userInfo.Name, userInfo.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// 更新用户头像（如果有）
	if userInfo.Picture.Data.Url != "" {
		u.Avatar = userInfo.Picture.Data.Url
		if _, err := s.userCase.Update(ctx, u); err != nil {
			s.log.Errorf("failed to update user avatar, error: %v", err)
		}
	}

	// 生成 JWT token
	token, err := s.jwt.GenerateToken(u.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 构建响应
	return &v1.LoginResponse{
		Token:     token,
		IsNewUser: isNew,
		UserInfo: &v1.UserInfo{
			UserId: u.UserID,
			Name:   u.Name,
			Avatar: u.Avatar,
			Phone:  u.Phone,
		},
	}, nil
}

func (s *FacebookService) getUserInfo(ctx context.Context, accessToken string) (*FacebookUserResponse, error) {
	// 构建请求 URL
	url := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email,picture&access_token=%s", accessToken)

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			s.log.Errorf("close io reader failed, error: %v", errClose)
		}
	}(resp.Body)

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("facebook api returned status code: %d", resp.StatusCode)
	}

	// 解析响应
	var userResponse FacebookUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, err
	}

	// 验证响应
	if userResponse.ID == "" {
		return nil, errors.New("invalid facebook user response")
	}

	return &userResponse, nil
}
