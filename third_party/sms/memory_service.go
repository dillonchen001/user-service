package sms

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// MemoryService 实现基于内存的SMS服务

type MemoryService struct {
	mu             sync.RWMutex
	codeStore      map[string]codeInfo
	expireDuration time.Duration
}

type codeInfo struct {
	code      string
	createdAt time.Time
}

// NewMemoryService 创建新的内存SMS服务
func NewMemoryService(expireDuration time.Duration) *MemoryService {
	return &MemoryService{
		codeStore:      make(map[string]codeInfo),
		expireDuration: expireDuration,
	}
}

// SendVerificationCode 发送验证码
func (s *MemoryService) SendVerificationCode(ctx context.Context, phone string) (string, error) {
	// 简单验证手机号格式
	if len(phone) < 10 {
		return "", ErrInvalidPhoneNumber
	}

	// 生成6位数验证码
	code := generateCode()

	// 存储验证码
	s.mu.Lock()
	s.codeStore[phone] = codeInfo{
		code:      code,
		createdAt: time.Now(),
	}
	s.mu.Unlock()

	// 在实际应用中，这里应该调用短信服务提供商的API
	// 例如: return code, s.provider.Send(phone, "Your verification code is: " + code)

	return code, nil
}

// VerifyCode 验证验证码
func (s *MemoryService) VerifyCode(ctx context.Context, phone string, code string) error {
	s.mu.RLock()
	info, exists := s.codeStore[phone]
	s.mu.RUnlock()

	if !exists {
		return ErrInvalidCode
	}

	// 检查验证码是否过期
	if time.Since(info.createdAt) > s.expireDuration {
		// 删除过期的验证码
		s.mu.Lock()
		delete(s.codeStore, phone)
		s.mu.Unlock()
		return ErrCodeExpired
	}

	// 检查验证码是否匹配
	if info.code != code {
		return ErrInvalidCode
	}

	// 验证成功后删除验证码
	s.mu.Lock()
	delete(s.codeStore, phone)
	s.mu.Unlock()

	return nil
}

// generateCode 生成6位数验证码
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
