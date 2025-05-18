package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Application ApplicationConfig `mapstructure:"application"`
}

type ApplicationConfig struct {
	BatchSize         int    `mapstructure:"batch-size"`
	InputPath         string `mapstructure:"input-path"`
	OutputPath        string `mapstructure:"output-path"`
	TempFileDirectory string `mapstructure:"temporary-file-directory"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("config") // look for config in the config directory
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return config, nil
}
