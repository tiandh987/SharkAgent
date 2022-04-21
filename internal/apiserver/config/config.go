package config

import "github.com/tiandh987/SharkAgent/internal/apiserver/options"

// Config 是 apiserver 运行时使用的配置
type Config struct {
	*options.Options
}

// CreateConfigFromOptions 基于命令行、配置文件 创建一个运行时使用的配置实例
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}