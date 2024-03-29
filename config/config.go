package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/restlesswhy/eth-balance-searcher/pkg/constants"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Search microservice config path")
}

type Config struct {
	ServiceName string     `mapstructure:"serviceName"`
	Logger      LogConfig  `mapstructure:"logger"`
	Http        HttpConfig `mapstructure:"http"`
	APIToken    string
}

type HttpConfig struct {
	Port string `mapstructure:"port" validate:"required"`
}

type LogConfig struct {
	LogLevel string `mapstructure:"level"`
	DevMode  bool   `mapstructure:"devMode"`
	Encoder  string `mapstructure:"encoder"`
}

func Load() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/config/config.yaml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := godotenv.Load(); err != nil {
		return nil, errors.Wrap(err, "load env error")
	}
	cfg.APIToken = os.Getenv("API_TOKEN")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "cannot read cofiguration")
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "environment cant be loaded")
	}

	return cfg, nil
}
