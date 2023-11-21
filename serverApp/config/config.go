package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTP            HTTP            `mapstructure:"http"`
		TCPSocket       TCPSocket       `mapstructure:"tcp_socket"`
		DistributedLock DistributedLock `mapstructure:"distributed_lock"`
		Logger          Logger          `mapstructure:"logger"`
		Neighbour       Neighbour       `mapstructure:"neighbour"`
	}

	HTTP struct {
		Port int `mapstructure:"port"`
	}

	TCPSocket struct {
		Port int `mapstructure:"port"`
	}

	DistributedLock struct {
		Port int `mapstructure:"port"`
	}

	Logger struct {
		Filename string `mapstructure:"filename"`
		Level    string `mapstructure:"level"`
	}

	Neighbour struct {
		HTTP      `mapstructure:"http"`
		TCPSocket `mapstructure:"tcp_socket"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	viper.SetConfigName(configPath)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
