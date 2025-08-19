package sms

import (
    "context"
    "errors"
)

// Service 定义SMS服务接口
type Service interface {
    // SendVerificationCode 发送验证码
    SendVerificationCode(ctx context.Context, phone string) (string, error)
    // VerifyCode 验证验证码
    VerifyCode(ctx context.Context, phone string, code string) error
}

// 定义错误类型
var (
    ErrInvalidPhoneNumber = errors.New("invalid phone number")
    ErrCodeExpired        = errors.New("verification code expired")
    ErrInvalidCode        = errors.New("invalid verification code")
    ErrSendFailed         = errors.New("failed to send verification code")
)