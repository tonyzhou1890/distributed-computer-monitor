package config

// 配置类型
type Config struct {
	Mode     string   `yaml:"mode"`
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

// 服务配置
type Server struct {
	AppName string `yaml:"appName"`
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
	Env     string `yaml:"env"`
}

// 数据库配置
type Database struct {
	Sqlite Sqlite `yaml:"sqlite"`
}

// sqlite 配置
type Sqlite struct {
	Path   string `yaml:"path"`
	MaxDay int    `yaml:"maxDay"`
}
