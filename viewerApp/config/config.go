package config

import (
	"github.com/spf13/viper"
)

const configPath = "/home/iuvarkin/unn/architecture/viewerApp/config"

type (
	Config struct {
		HTTP   HTTP   `mapstructure:"http"`
		Logger Logger `mapstructure:"logger"`
		Server Server `mapstructure:"server"`
	}

	HTTP struct {
		Port int `mapstructure:"port"`
	}

	Logger struct {
		Filename string `mapstructure:"filename"`
		Level    string `mapstructure:"level"`
	}

	Server struct {
		HTTP `mapstructure:"http"`
	}
)

func NewConfig(configName string) (*Config, error) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
