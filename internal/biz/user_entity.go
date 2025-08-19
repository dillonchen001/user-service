package biz

import (
	"time"
)

// User 用户实体
type User struct {
	ID        int64     `gorm:"primary_key;unique" json:"id"`
	UserID    int64     `gorm:"unique" json:"user_id"`
	Name      string    `gorm:"size:100" json:"name"`
	Email     string    `gorm:"size:255;unique" json:"email"`
	Phone     string    `gorm:"size:20;unique" json:"phone"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthProvider 认证提供者
type AuthProvider struct {
	ID           int64     `gorm:"primary_key;unique" json:"id"`
	UserID       int64     `gorm:"unique" json:"user_id"`
	ProviderType string    `gorm:"size:20;index:idx_provider_type_id" json:"provider_type"` // phone, facebook, apple, google, snapchat
	ProviderID   string    `gorm:"size:100;index:idx_provider_type_id" json:"provider_id"`
	CreatedAt    time.Time `json:"created_at"`
}
