/*
Copyright © 2025 lixw
*/
package config

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	Logging  *LoggingConfig
}

type ServerConfig struct {
	Addr            string
	StartTimeout    time.Duration
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxHeaderBytes  int
	BasePath        string
	Locale          string
}

type DatabaseConfig struct {
	Url             string
	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
	SlowThreshold   time.Duration
}

type LoggingConfig struct {
	Level      string
	Path       string
	MaxAge     int
	MaxSize    int
	MaxBackups int
	Compress   bool
	Format     string
}

func New(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("APP")

	setDefaultConfig()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	configFileContent, err := os.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		return nil, fmt.Errorf("error reading config file content: %w", err)
	}

	// 替换${ENV_VAR}格式的环境变量引用
	re := regexp.MustCompile(`\${([^}]+)}`)
	result := re.ReplaceAllStringFunc(string(configFileContent), func(match string) string {
		// 提取环境变量名称（去掉${}部分）
		envVar := match[2 : len(match)-1]
		// 获取环境变量值，如果不存在则保持原样
		if value := os.Getenv(envVar); value != "" {
			return value
		}
		return match
	})

	// 使用处理后的配置内容
	_ = viper.ReadConfig(strings.NewReader(result))
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	slog.Info("using config file", "path", viper.ConfigFileUsed())
	return &cfg, nil
}

func setDefaultConfig() {
	// server
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("server.startTimeout", 15*time.Second)
	viper.SetDefault("server.shutdownTimeout", 15*time.Second)
	viper.SetDefault("server.readTimeout", 5*time.Second)
	viper.SetDefault("server.writeTimeout", 10*time.Second)
	viper.SetDefault("server.idleTimeout", 30*time.Second)
	viper.SetDefault("server.maxHeaderBytes", 1<<20) // 1MB
	viper.SetDefault("server.basePath", "/")
	viper.SetDefault("server.locale", "zh-CN")

	// database
	viper.SetDefault("database.connMaxIdleTime", 5*time.Minute)
	viper.SetDefault("database.connMaxLifeTime", 30*time.Minute)
	viper.SetDefault("database.maxIdleConns", 5)
	viper.SetDefault("database.maxOpenConns", 10)
	viper.SetDefault("database.slowThreshold", 500*time.Millisecond)

	// logging
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.path", "logs/app.log")
	viper.SetDefault("logging.maxAge", 7)    // 7天
	viper.SetDefault("logging.maxSize", 100) // 100MB
	viper.SetDefault("logging.maxBackups", 10)
	viper.SetDefault("logging.compress", true)
	viper.SetDefault("logging.format", "text")
}

func GetString(key string) string {
	return viper.GetString(key)
}
