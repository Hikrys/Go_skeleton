package config

import "github.com/google/wire"

// NewServerConfig 把大配置里的 Server 部分拆出来
// Wire 有 *Config，能搞出 ServerConfig
func NewServerConfig(c *Config) ServerConfig {
	return c.Server
}

// NewLogConfig 把大配置里的 Log 部分拆出来
func NewLogConfig(c *Config) LogConfig {
	return c.Log
}

// ProviderSet 告诉 Wire 使用上面这两个函数
var ProviderSet = wire.NewSet(
	NewServerConfig,
	NewLogConfig,
)
