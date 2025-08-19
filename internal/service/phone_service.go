package service

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	v1 "user-service/api/auth/v1"
	"user-service/internal/biz"
	"user-service/internal/conf"
	"user-service/third_party/jwt"
	"user-service/third_party/sms"
)

// 修改PhoneService结构体

type PhoneService struct {
	log        *log.Helper
	cfg        *conf.Jwt
	userCase   *biz.UserCase
	jwtGen     *jwt.Generator
	smsService sms.Service // 添加SMS服务
}

// 修改NewPhoneService函数

func NewPhoneService(cfg *conf.Jwt, logger log.Logger, userCase *biz.UserCase) *PhoneService {
	// 创建SMS服务
	smsConfig := sms.DefaultConfig()
	smsService := sms.NewService(smsConfig)

	return &PhoneService{
		cfg:        cfg,
		log:        log.NewHelper(logger),
		userCase:   userCase,
		jwtGen:     jwt.NewGenerator(cfg.Secret, int(cfg.Expires)),
		smsService: smsService,
	}
}

func (s *PhoneService) Login(ctx context.Context, req *v1.LoginWithPhoneRequest) (*v1.LoginResponse, error) {
	// 验证手机号格式
	if !isValidPhone(req.PhoneNumber) {
		return nil, errors.New("invalid phone number")
	}

	// 测试不收验证码
	code, _ := s.smsService.SendVerificationCode(ctx, req.PhoneNumber)

	// 验证验证码
	if err := s.smsService.VerifyCode(ctx, req.PhoneNumber, code); err != nil {
		return nil, err
	}

	// 查找或创建用户
	u, isNew, err := s.userCase.FindOrCreateByPhone(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	// 生成JWT token
	token, err := s.jwtGen.GenerateToken(u.UserID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	return &v1.LoginResponse{
		Token:     token,
		IsNewUser: isNew,
		UserInfo: &v1.UserInfo{
			UserId: u.UserID,
			Name:   u.Name,
			Phone:  u.Phone,
		},
	}, nil
}

// 简单的手机号格式验证
func isValidPhone(phone string) bool {
	// 实际应用中应该使用更严格的验证
	return len(phone) >= 10 && len(phone) <= 15
}
