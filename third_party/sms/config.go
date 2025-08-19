package sms

import (
    "time"
)

// Config 定义SMS服务配置

type Config struct {
    ExpireDuration time.Duration `json:"expire_duration"`
    // 可以添加其他配置，如API密钥、短信模板等
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
    return Config{
        ExpireDuration: 5 * time.Minute,
    }
}