package service

import (
	"context"

	v1 "user-service/api/auth/v1"
	"user-service/internal/biz"
	"user-service/internal/conf"
	"user-service/third_party/jwt"
	"user-service/third_party/snowflake"

	"github.com/go-kratos/kratos/v2/log"
)

type LoginService struct {
	v1.UnimplementedAuthServiceServer
	log             *log.Helper
	uidGen          *snowflake.Node
	phoneService    *PhoneService
	facebookService *FacebookService
	appleService    *AppleService
	googleService   *GoogleService
	snapchatService *SnapchatService
	jwtGenerator    *jwt.Generator
}

func NewLoginService(cfg *conf.Jwt, logger log.Logger, uidGen *snowflake.Node, userAuthCase *biz.UserAuthCase, userCase *biz.UserCase) *LoginService {
	jwtGenerator := jwt.NewGenerator(cfg.Secret, int(cfg.Expires))

	return &LoginService{
		log:             log.NewHelper(logger),
		uidGen:          uidGen,
		phoneService:    NewPhoneService(cfg, logger, userCase),
		facebookService: NewFacebookService(cfg, logger, userAuthCase, userCase),
		appleService:    NewAppleService(cfg, logger, userAuthCase),
		googleService:   NewGoogleService(cfg, logger, userAuthCase, userCase),
		snapchatService: NewSnapchatService(cfg, logger, userAuthCase, userCase),
		jwtGenerator:    jwtGenerator,
	}
}

// LoginWithPhone 手机号登录
func (s *LoginService) LoginWithPhone(ctx context.Context, req *v1.LoginWithPhoneRequest) (*v1.LoginResponse, error) {
	return s.phoneService.Login(ctx, req)
}

// LoginWithFacebook Facebook登录
func (s *LoginService) LoginWithFacebook(ctx context.Context, req *v1.LoginWithFacebookRequest) (*v1.LoginResponse, error) {
	return s.facebookService.Login(ctx, req)
}

// LoginWithApple Apple登录
func (s *LoginService) LoginWithApple(ctx context.Context, req *v1.LoginWithAppleRequest) (*v1.LoginResponse, error) {
	return s.appleService.Login(ctx, req)
}

// LoginWithGoogle Google登录
func (s *LoginService) LoginWithGoogle(ctx context.Context, req *v1.LoginWithGoogleRequest) (*v1.LoginResponse, error) {
	return s.googleService.Login(ctx, req)
}

// LoginWithSnapchat Snapchat登录
func (s *LoginService) LoginWithSnapchat(ctx context.Context, req *v1.LoginWithSnapchatRequest) (*v1.LoginResponse, error) {
	return s.snapchatService.Login(ctx, req)
}
