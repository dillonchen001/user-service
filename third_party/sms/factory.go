package sms

// NewService 创建SMS服务
func NewService(config Config) Service {
	// 在实际应用中，可以根据配置选择不同的SMS服务实现
	// 例如: 如果配置了短信提供商的API密钥，则使用该提供商的实现
	// 否则，使用内存实现（仅用于开发和测试）

	return NewMemoryService(config.ExpireDuration)
}
