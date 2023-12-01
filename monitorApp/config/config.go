package config

import (
	"github.com/spf13/viper"
)

const configPath = "/home/iuvarkin/unn/architecture/monitorApp/config"

type (
	Config struct {
		Logger   Logger   `mapstructure:"logger"`
		Watchdog Watchdog `mapstructure:"watchdog"`
		Server   Server   `mapstructure:"server"`
	}

	Logger struct {
		Filename string `mapstructure:"filename"`
		Level    string `mapstructure:"level"`
	}

	Watchdog struct {
		Filename string `mapstructure:"filename"`
		StartURL string `mapstructure:"start_url"`
		Interval int    `mapstructure:"interval"`
		MaxWait  int    `mapstructure:"max_wait"`
	}

	Server struct {
		ConfigPath string `mapstructure:"config"`
		BinPath    string `mapstructure:"bin_path"`
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
