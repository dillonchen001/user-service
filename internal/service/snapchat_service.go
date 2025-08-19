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
)

type SnapchatService struct {
	cfg          *conf.Jwt
	log          *log.Helper
	userAuthCase *biz.UserAuthCase
	userCase     *biz.UserCase
	jwtGen       *jwt.Generator
	httpClient   *http.Client
}

func NewSnapchatService(cfg *conf.Jwt, logger log.Logger, userAuthCase *biz.UserAuthCase, userCase *biz.UserCase) *SnapchatService {
	return &SnapchatService{
		cfg:          cfg,
		log:          log.NewHelper(logger),
		userAuthCase: userAuthCase,
		userCase:     userCase,
		jwtGen:       jwt.NewGenerator(cfg.Secret, int(cfg.Expires)),
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// SnapchatUserResponse 用户信息响应结构
type SnapchatUserResponse struct {
	Data struct {
		ExternalID  string `json:"external_id"`
		DisplayName string `json:"display_name"`
		Bitmoji     struct {
			AvatarURL string `json:"avatar_url"`
		} `json:"bitmoji"`
	} `json:"data"`
}

// SnapchatTokenResponse 访问令牌验证响应结构
type SnapchatTokenResponse struct {
	Data struct {
		Valid  bool   `json:"valid"`
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (s *SnapchatService) Login(ctx context.Context, req *v1.LoginWithSnapchatRequest) (*v1.LoginResponse, error) {
	// 验证参数
	if req.AccessToken == "" {
		return nil, errors.New("access_token is required")
	}

	// 验证 access token
	tokenInfo, err := s.verifyAccessToken(ctx, req.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify snapchat access token: %w", err)
	}

	if !tokenInfo.Data.Valid {
		return nil, errors.New("invalid or expired access token")
	}

	// 获取用户信息
	userInfo, err := s.getUserInfo(ctx, req.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapchat user info: %w", err)
	}

	// 查找或创建用户
	u, isNew, err := s.userAuthCase.FindOrCreateBySnapchatID(ctx,
		tokenInfo.Data.UserID,
		userInfo.Data.DisplayName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// 更新用户头像（如果有）
	if userInfo.Data.Bitmoji.AvatarURL != "" && u.Avatar != userInfo.Data.Bitmoji.AvatarURL {
		u.Avatar = userInfo.Data.Bitmoji.AvatarURL

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
		},
	}, nil
}

func (s *SnapchatService) verifyAccessToken(ctx context.Context, accessToken string) (*SnapchatTokenResponse, error) {
	// 构建请求
	url := "https://kit.snapchat.com/oauth2/verify"
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

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
		return nil, fmt.Errorf("snapchat token verification api returned status code: %d", resp.StatusCode)
	}

	// 解析响应
	var tokenResponse SnapchatTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (s *SnapchatService) getUserInfo(ctx context.Context, accessToken string) (*SnapchatUserResponse, error) {
	// 构建请求
	url := "https://kit.snapchat.com/v1/me"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

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
		return nil, fmt.Errorf("snapchat user info api returned status code: %d", resp.StatusCode)
	}

	// 解析响应
	var userResponse SnapchatUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, err
	}

	// 验证响应
	if userResponse.Data.ExternalID == "" {
		return nil, errors.New("invalid snapchat user response")
	}

	return &userResponse, nil
}
