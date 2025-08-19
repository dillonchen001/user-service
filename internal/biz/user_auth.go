package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// UserAuthCase 用户关联授权登陆实例的使用
type UserAuthCase struct {
	userRepo UserRepo
	authRepo AuthProviderRepo
	log      *log.Helper
}

// NewUserAuthCase 创建新的用户关联授权登陆实例
func NewUserAuthCase(userRepo UserRepo, authRepo AuthProviderRepo, logger log.Logger) *UserAuthCase {
	return &UserAuthCase{
		userRepo: userRepo,
		authRepo: authRepo,
		log:      log.NewHelper(logger),
	}
}

// LinkAuthProvider 关联认证提供者
func (uc *UserAuthCase) LinkAuthProvider(ctx context.Context, userID int64, providerType, providerID string) error {
	uc.log.WithContext(ctx).Infof("LinkAuthProvider: %v %v %v", userID, providerType, providerID)
	// 首先确认用户存在
	userInfo, err := uc.userRepo.FindByIDOrigin(ctx, userID)
	if err != nil {
		return err
	}

	// 创建认证提供者记录
	return uc.authRepo.Create(ctx, providerType, providerID, userInfo)
}

// FindOrCreateByGoogleID 根据Google ID查找或创建用户
func (uc *UserAuthCase) FindOrCreateByGoogleID(ctx context.Context, googleID, name, email string) (*User, bool, error) {
	uc.log.WithContext(ctx).Infof("FindOrCreateByGoogleID: %v %v %v", googleID, name, email)
	found, err := uc.authRepo.FindByGoogleID(ctx, googleID)
	if err == nil {
		return found, false, nil
	}

	u := &User{Name: name, Email: email}
	createdUser, err := uc.userRepo.Create(ctx, u)
	if err != nil {
		return nil, true, err
	}

	err = uc.LinkAuthProvider(ctx, createdUser.UserID, "google", googleID)
	if err != nil {
		return nil, true, err
	}

	return createdUser, true, nil
}

// FindOrCreateByAppleID 根据Apple ID查找或创建用户
func (uc *UserAuthCase) FindOrCreateByAppleID(ctx context.Context, appleID string, email string) (*User, bool, error) {
	uc.log.WithContext(ctx).Infof("FindOrCreateByAppleID: %v %v", appleID, email)
	// 尝试根据Apple ID查找用户
	found, err := uc.authRepo.FindByAppleID(ctx, appleID)
	if err == nil {
		// 用户已存在，返回找到的用户
		return found, false, nil
	}

	// 如果用户不存在，创建新用户
	u := &User{Email: email}

	createdUser, err := uc.userRepo.Create(ctx, u)
	if err != nil {
		return nil, true, err
	}

	// 关联Apple ID和新创建的用户
	err = uc.LinkAuthProvider(ctx, createdUser.UserID, "apple", appleID)
	if err != nil {
		return nil, true, err
	}

	return createdUser, true, nil
}

// FindOrCreateByFacebookID 根据Facebook ID查找或创建用户
func (uc *UserAuthCase) FindOrCreateByFacebookID(ctx context.Context, facebookID, name, email string) (*User, bool, error) {
	uc.log.WithContext(ctx).Infof("FindOrCreateByFacebookID: %v %v %v", facebookID, name, email)
	// 根据Facebook ID查找用户
	found, err := uc.authRepo.FindByFacebookID(ctx, facebookID)
	if err == nil {
		// 用户已存在，返回找到的用户
		return found, false, nil
	}

	// 如果用户不存在，创建新用户
	u := &User{Name: name, Email: email}
	createdUser, err := uc.userRepo.Create(ctx, u)
	if err != nil {
		return nil, true, err
	}

	// 关联Facebook ID和新创建的用户
	err = uc.LinkAuthProvider(ctx, createdUser.UserID, "facebook", facebookID)
	if err != nil {
		return nil, true, err
	}

	return createdUser, true, nil
}

// FindOrCreateBySnapchatID 根据snap ID查找或创建用户
func (uc *UserAuthCase) FindOrCreateBySnapchatID(ctx context.Context, snapchatID, name string) (*User, bool, error) {
	uc.log.WithContext(ctx).Infof("FindOrCreateBySnapchatID: %v %v", snapchatID, name)
	// 根据Facebook ID查找用户
	found, err := uc.authRepo.FindBySnapchatID(ctx, snapchatID)
	if err == nil {
		// 用户已存在，返回找到的用户
		return found, false, nil
	}

	// 如果用户不存在，创建新用户
	u := &User{Name: name}
	createdUser, err := uc.userRepo.Create(ctx, u)
	if err != nil {
		return nil, true, err
	}

	// 关联Facebook ID和新创建的用户
	err = uc.LinkAuthProvider(ctx, createdUser.UserID, "snapchat", snapchatID)
	if err != nil {
		return nil, true, err
	}

	return createdUser, true, nil
}
