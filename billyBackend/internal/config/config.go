package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Session SessionConfig `mapstructure:"session"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Tracing TracingConfig `mapstructure:"tracing"`
	Logging LoggingConfig `mapstructure:"logging"`
}

type ServerConfig struct {
	Port     int `mapstructure:"port"`
	Timeout  int `mapstructure:"timeout"`
	MaxItems int `mapstructure:"maxItems"`
}

type SessionConfig struct {
	Store string `mapstructure:"store"`
}

type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type TracingConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}

type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
