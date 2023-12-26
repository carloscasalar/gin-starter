package app

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port string `default:"8080"`
}

func ReadConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("SERVICE", &cfg); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	return &cfg, nil
}
