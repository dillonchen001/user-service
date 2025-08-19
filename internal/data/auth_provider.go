package data

import (
	"context"

	"user-service/internal/biz"
	"user-service/internal/data/ent"
	"user-service/internal/data/ent/authprovider"

	"github.com/go-kratos/kratos/v2/log"
)

// authProviderRepo 实现用户授权接口
type authProviderRepo struct {
	data *Data
	log  *log.Helper
}

// NewAuthProviderRepo 创建新的用户授权仓储
func NewAuthProviderRepo(data *Data, logger log.Logger) biz.AuthProviderRepo {
	return &authProviderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create 创建授权用户
func (r *authProviderRepo) Create(ctx context.Context, proType, id string, userInfo *ent.User) error {
	// 创建认证提供者记录
	_, err := r.data.db.AuthProvider.Create().
		SetProviderType(proType).
		SetProviderID(id).
		SetUser(userInfo).
		SetUserID(userInfo.UserID).
		Save(ctx)

	return err
}

// FindByGoogleID 根据Google ID查找用户
func (r *authProviderRepo) FindByGoogleID(ctx context.Context, googleID string) (*biz.User, error) {
	authProvider, err := r.data.db.AuthProvider.Query().
		Where(
			authprovider.ProviderType("google"),
			authprovider.ProviderID(googleID),
		).
		WithUser().
		First(ctx)
	if err != nil {
		return nil, err
	}

	userInfo := authProvider.Edges.User
	return &biz.User{
		ID:        userInfo.ID,
		UserID:    userInfo.UserID,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		Phone:     userInfo.Phone,
		Avatar:    userInfo.Avatar,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
	}, nil
}

// FindByAppleID 根据Apple ID查找用户
func (r *authProviderRepo) FindByAppleID(ctx context.Context, appleID string) (*biz.User, error) {
	authProvider, err := r.data.db.AuthProvider.Query().
		Where(
			authprovider.ProviderType("apple"),
			authprovider.ProviderID(appleID),
		).
		WithUser().
		First(ctx)
	if err != nil {
		return nil, err
	}

	userInfo := authProvider.Edges.User
	return &biz.User{
		ID:        userInfo.ID,
		UserID:    userInfo.UserID,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		Phone:     userInfo.Phone,
		Avatar:    userInfo.Avatar,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
	}, nil
}

// FindByFacebookID 根据Facebook ID查找用户
func (r *authProviderRepo) FindByFacebookID(ctx context.Context, facebookID string) (*biz.User, error) {
	authProvider, err := r.data.db.AuthProvider.Query().
		Where(
			authprovider.ProviderType("facebook"),
			authprovider.ProviderID(facebookID),
		).
		WithUser().
		First(ctx)
	if err != nil {
		return nil, err
	}

	userInfo := authProvider.Edges.User
	return &biz.User{
		ID:        userInfo.ID,
		UserID:    userInfo.UserID,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		Phone:     userInfo.Phone,
		Avatar:    userInfo.Avatar,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
	}, nil
}

// FindBySnapchatID 根据Snapchat ID查找用户
func (r *authProviderRepo) FindBySnapchatID(ctx context.Context, snapchatID string) (*biz.User, error) {
	authProvider, err := r.data.db.AuthProvider.Query().
		Where(
			authprovider.ProviderType("snapchat"),
			authprovider.ProviderID(snapchatID),
		).
		WithUser().
		First(ctx)
	if err != nil {
		return nil, err
	}

	userInfo := authProvider.Edges.User
	return &biz.User{
		ID:        userInfo.ID,
		UserID:    userInfo.UserID,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		Phone:     userInfo.Phone,
		Avatar:    userInfo.Avatar,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
	}, nil
}
