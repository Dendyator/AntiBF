package config

import (
	"github.com/Dendyator/AntiBF/pkg/logger"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Logger      LoggerConfig      `mapstructure:"logger"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	RateLimiter RateLimiterConfig `mapstructure:"rate_limiter"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Address string `mapstructure:"address"`
}

type RateLimiterConfig struct {
	LoginLimit    int `mapstructure:"login_limit"`
	PasswordLimit int `mapstructure:"password_limit"`
	IPLimit       int `mapstructure:"ip_limit"`
}

func LoadConfig(configPath string, appLogger *logger.Logger) Config {
	var config Config

	appLogger.Infof("Loading configuration from %s", configPath)

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		appLogger.Fatalf("Error reading config file: %s", err)
	}

	appLogger.Infof("RateLimiter values: login_limit=%d, password_limit=%d, ip_limit=%d",
		viper.GetInt("rate_limiter.login_limit"),
		viper.GetInt("rate_limiter.password_limit"),
		viper.GetInt("rate_limiter.ip_limit"))

	err := viper.Unmarshal(&config)
	if err != nil {
		appLogger.Fatalf("Unable to decode into struct, %v", err)
	}

	appLogger.Infof("Configuration unmarshaled successfully: %+v", config)
	return config
}
