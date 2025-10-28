package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	Port        string `mapstructure:"SERVER_PORT"`
	MetricsPort string `mapstructure:"METRICS_PORT"`
	GRPCPort    string `mapstructure:"GRPC_PORT"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"POSTGRES_PORT"`
	DBName      string `mapstructure:"POSTGRES_DB"`
	DBUser      string `mapstructure:"POSTGRES_USER"`
	DBPass      string `mapstructure:"POSTGRES_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	viper.AutomaticEnv()

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName)
}

func validateConfig(cfg *Config) error {
	required := map[string]string{
		"JWT_SECRET":        cfg.JWTSecret,
		"SERVER_PORT":       cfg.Port,
		"METRICS_PORT":      cfg.MetricsPort,
		"GRPC_PORT":         cfg.GRPCPort,
		"DB_HOST":           cfg.DBHost,
		"POSTGRES_PORT":     cfg.DBPort,
		"POSTGRES_DB":       cfg.DBName,
		"POSTGRES_USER":     cfg.DBUser,
		"POSTGRES_PASSWORD": cfg.DBPass,
	}

	for field, value := range required {
		if value == "" {
			return fmt.Errorf("%s is reqired", field)
		}
	}

	return nil
}
