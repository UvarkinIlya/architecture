package config

import (
	"github.com/spf13/viper"
)

const configPath = "/home/iuvarkin/unn/architecture/serverApp/config"

type (
	Config struct {
		HTTP            HTTP            `mapstructure:"http"`
		TCPSocket       TCPSocket       `mapstructure:"tcp_socket"`
		DistributedLock DistributedLock `mapstructure:"distributed_lock"`
		Logger          Logger          `mapstructure:"logger"`
		Storage         Storage         `mapstructure:"storage"`
		Syncer          Syncer          `mapstructure:"syncer"`
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

	Storage struct {
		MessageFilePath string `mapstructure:"messages_file_path"`
	}

	Syncer struct {
		HTTP `mapstructure:"http"`
	}

	Neighbour struct {
		HTTP      `mapstructure:"http"`
		TCPSocket `mapstructure:"tcp_socket"`
		Syncer    `mapstructure:"syncer"`
	}
)

func NewConfig(configName string) (*Config, error) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
