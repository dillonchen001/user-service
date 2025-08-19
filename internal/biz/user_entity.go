package biz

import (
	"time"

	"github.com/google/uuid"
)

// User 用户实体
type User struct {
	ID        string    `gorm:"primary_key;size:36" json:"id"`
	Name      string    `gorm:"size:100" json:"name"`
	Email     string    `gorm:"size:255;unique" json:"email"`
	Phone     string    `gorm:"size:20;unique" json:"phone"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate 钩子函数，自动生成UUID
func (u *User) BeforeCreate() error {
	u.ID = uuid.New().String()
	return nil
}

// AuthProvider 认证提供者
type AuthProvider struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	UserID       string    `gorm:"size:36;index" json:"user_id"`
	ProviderType string    `gorm:"size:20;index:idx_provider_type_id" json:"provider_type"` // phone, facebook, apple, google, snapchat
	ProviderID   string    `gorm:"size:100;index:idx_provider_type_id" json:"provider_id"`
	CreatedAt    time.Time `json:"created_at"`
}
