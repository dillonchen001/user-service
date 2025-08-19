package biz

import (
	"context"

	"user-service/internal/data/ent"

	"github.com/go-kratos/kratos/v2/log"
)

// AuthProviderRepo 定义用户授权接口
type AuthProviderRepo interface {
	FindByAppleID(ctx context.Context, appleID string) (*User, error)
	// FindByGoogleID 根据Google ID查找用户
	FindByGoogleID(ctx context.Context, googleID string) (*User, error)
	// FindByFacebookID 根据Facebook ID查找用户
	FindByFacebookID(ctx context.Context, facebookID string) (*User, error)
	// FindBySnapchatID 根据Snapchat ID查找用户
	FindBySnapchatID(ctx context.Context, snapchatID string) (*User, error)
	// Create 创建授权用户
	Create(ctx context.Context, proType, id string, userInfo *ent.User) error
}

// AuthProviderCase 用户登陆授权实例的使用
type AuthProviderCase struct {
	repo AuthProviderRepo
	log  *log.Helper
}

// NewAuthProviderCase 创建新的用户登陆授权实例
func NewAuthProviderCase(repo AuthProviderRepo, logger log.Logger) *AuthProviderCase {
	return &AuthProviderCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// FindByGoogleID 根据Google ID查找用户
func (uc *AuthProviderCase) FindByGoogleID(ctx context.Context, googleID string) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindByGoogleID: %v", googleID)
	return uc.repo.FindByGoogleID(ctx, googleID)
}

// FindByAppleID 根据Apple ID查找用户
func (uc *AuthProviderCase) FindByAppleID(ctx context.Context, appleID string) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindByAppleID: %v", appleID)
	return uc.repo.FindByAppleID(ctx, appleID)
}

// FindByFacebookID 根据Facebook ID查找用户
func (uc *AuthProviderCase) FindByFacebookID(ctx context.Context, facebookID string) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindByFacebookID: %v", facebookID)
	return uc.repo.FindByFacebookID(ctx, facebookID)
}

// FindBySnapchatID 根据Snapchat ID查找用户
func (uc *AuthProviderCase) FindBySnapchatID(ctx context.Context, snapchatID string) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindBySnapchatID: %v", snapchatID)
	return uc.repo.FindBySnapchatID(ctx, snapchatID)
}
