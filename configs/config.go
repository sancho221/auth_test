package configs

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	JWTSecret string `env:"JWT_SECRET"`
	Port      string `env:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
