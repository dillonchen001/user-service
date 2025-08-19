package biz

import (
	"context"

	v1 "user-service/api/helloworld/v1"
	"user-service/internal/data/ent"
	"user-service/third_party/snowflake"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// UserRepo 定义用户仓储接口
type UserRepo interface {
	// FindByIDOrigin 根据ID查找用户, 原始输出
	FindByIDOrigin(ctx context.Context, userID int64) (*ent.User, error)
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, userID int64) (*User, error)
	// FindByPhone 根据手机号查找用户
	FindByPhone(ctx context.Context, phone string) (*User, error)
	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*User, error)
	// Create 创建用户
	Create(ctx context.Context, u *User) (*User, error)
	// Update 更新用户
	Update(ctx context.Context, u *User) (*User, error)
	// FindOrCreate 查找或创建用户
	FindOrCreate(ctx context.Context, u *User) (*User, error)
	// FindOrCreateByPhone 根据Phone 查找或创建用户
	FindOrCreateByPhone(ctx context.Context, phone string) (*User, bool, error)
}

// UserCase 用户实例的使用
type UserCase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserCase 创建新的用户实例
func NewUserCase(repo UserRepo, logger log.Logger) *UserCase {
	return &UserCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// FindByID find user by id
func (uc *UserCase) FindByID(ctx context.Context, userID int64) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindByID: %v", userID)
	return uc.repo.FindByID(ctx, userID)
}

// FindByPhone 根据手机号查找用户
func (uc *UserCase) FindByPhone(ctx context.Context, phone string) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindByPhone: %v", phone)
	return uc.repo.FindByPhone(ctx, phone)
}

// FindByEmail 根据邮箱查找用户
func (uc *UserCase) FindByEmail(ctx context.Context, email string) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindByEmail: %v", email)
	return uc.repo.FindByEmail(ctx, email)
}

// Create 创建用户
func (uc *UserCase) Create(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("Create: %v", u)
	return uc.repo.Create(ctx, u)
}

// Update 更新用户
func (uc *UserCase) Update(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("Update: %v", u)
	return uc.repo.Update(ctx, u)
}

// FindOrCreate 查找或创建用户
func (uc *UserCase) FindOrCreate(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindOrCreate: %v", u)
	return uc.repo.FindOrCreate(ctx, u)
}

// FindOrCreateByPhone 根据Phone 查找或创建用户
func (uc *UserCase) FindOrCreateByPhone(ctx context.Context, uidGen *snowflake.Node, phone string) (*User, bool, error) {
	uc.log.WithContext(ctx).Infof("FindOrCreateByPhone: %v", phone)
	return uc.repo.FindOrCreateByPhone(ctx, phone)
}
