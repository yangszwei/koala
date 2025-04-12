package config

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/spf13/viper"
)

//go:embed default.yaml
var defaultConfig []byte

// Config is the top-level application configuration structure.
type Config struct {
	Http HttpConfig `mapstructure:"http"`
}

// HttpConfig holds HTTP server configuration parameters such as address binding.
type HttpConfig struct {
	Addr string `mapstructure:"address"`
}

// Load loads configuration from a YAML file.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// Load embedded default config
	if err := v.ReadConfig(bytes.NewBuffer(defaultConfig)); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Optional YAML config file (overrides defaults)
	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to merge config: %w", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &c, nil
}
