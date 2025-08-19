package data

import (
	"context"
	"user-service/internal/data/ent"

	"user-service/internal/biz"
	"user-service/internal/data/ent/user"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo 创建新的用户仓储
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// FindByIDOrigin 根据ID查找用户, 原始输出
func (r *userRepo) FindByIDOrigin(ctx context.Context, id string) (*ent.User, error) {
	return r.data.db.User.Query().
		Where(user.ID(id)).
		First(ctx)
}

// FindByID 根据ID查找用户
func (r *userRepo) FindByID(ctx context.Context, id string) (*biz.User, error) {
	u, err := r.data.db.User.Query().
		Where(user.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

// FindByPhone 根据手机号查找用户
func (r *userRepo) FindByPhone(ctx context.Context, phone string) (*biz.User, error) {
	u, err := r.data.db.User.Query().
		Where(user.Phone(phone)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*biz.User, error) {
	u, err := r.data.db.User.Query().
		Where(user.Email(email)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

// Create 创建用户
func (r *userRepo) Create(ctx context.Context, u *biz.User) (*biz.User, error) {
	created, err := r.data.db.User.Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPhone(u.Phone).
		SetAvatar(u.Avatar).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:        created.ID,
		Name:      created.Name,
		Email:     created.Email,
		Phone:     created.Phone,
		Avatar:    created.Avatar,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}

// Update 更新用户
func (r *userRepo) Update(ctx context.Context, u *biz.User) (*biz.User, error) {
	updated, err := r.data.db.User.UpdateOneID(u.ID).
		SetName(u.Name).
		SetEmail(u.Email).
		SetPhone(u.Phone).
		SetAvatar(u.Avatar).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:        updated.ID,
		Name:      updated.Name,
		Email:     updated.Email,
		Phone:     updated.Phone,
		Avatar:    updated.Avatar,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}

// FindOrCreate 查找或创建用户
func (r *userRepo) FindOrCreate(ctx context.Context, u *biz.User) (*biz.User, error) {
	// 尝试根据邮箱查找用户
	found, err := r.data.db.User.Query().
		Where(user.Email(u.Email)).
		First(ctx)
	if err == nil {
		// 用户已存在，返回找到的用户
		return &biz.User{
			ID:        found.ID,
			Name:      found.Name,
			Email:     found.Email,
			Phone:     found.Phone,
			Avatar:    found.Avatar,
			CreatedAt: found.CreatedAt,
			UpdatedAt: found.UpdatedAt,
		}, nil
	}

	// 如果用户不存在，创建新用户
	return r.Create(ctx, u)
}

// FindOrCreateByPhone 根据Phone 查找或创建用户
func (r *userRepo) FindOrCreateByPhone(ctx context.Context, phone string) (*biz.User, bool, error) {
	// 尝试根据邮箱查找用户
	found, err := r.data.db.User.Query().
		Where(user.Phone(phone)).
		First(ctx)
	if err == nil {
		// 用户已存在，返回找到的用户
		return &biz.User{
			ID:        found.ID,
			Name:      found.Name,
			Email:     found.Email,
			Phone:     found.Phone,
			Avatar:    found.Avatar,
			CreatedAt: found.CreatedAt,
			UpdatedAt: found.UpdatedAt,
		}, false, nil
	}

	// 如果用户不存在，创建新用户
	u := &biz.User{Phone: phone}
	userInfo, errCre := r.Create(ctx, u)
	return userInfo, true, errCre
}
