package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// ServerConfig 对应 yaml 里的 server 部分
// mapstructure 是 viper 用来把 yaml 字段映射到结构体的标签
type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// LogConfig 对应 yaml 里的 log 部分
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// Config 这是一个总的结构体，把上面的小结构体都包进来
// 这样我们在代码里就能用 config.Server.Port 这样点出来，很方便
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
}

// LoadConfig 是一个“工厂函数”，负责读取配置文件
// path 是配置文件的路径，比如 "./configs"
func LoadConfig(path string) (*Config, error) {
	// SetConfigName 告诉 viper 我们的配置文件叫什么名字不需要写 .yaml 后缀
	viper.SetConfigName("config")
	// 告诉 viper 配置文件类型是 yaml
	viper.SetConfigType("yaml")
	// 告诉 viper 去哪里找文件
	viper.AddConfigPath(path)

	// 开始读取文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("哎呀，读取配置文件失败了: %w", err)
	}
	// 创建一个空的结构体变量
	var config Config
	// 把读取到的 yaml 内容，自动填充到 config 结构体里
	// Unmarshal反序列化到config中
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("格式不对，解析配置文件失败: %w", err)
	}
	//返回配置对象
	return &config, nil
}
