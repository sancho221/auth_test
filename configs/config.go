package configs

import "github.com/spf13/viper"

type Config struct {
	JWTSecret string
	Port      string
}

func LoadConfig() *Config {
	viper.SetDefault("JWT_SECRET", "default-secret-key")
	viper.SetDefault("SERVER_PORT", "8082")

	viper.AutomaticEnv()

	return &Config{
		JWTSecret: viper.GetString("JWT_SECRET"),
		Port:      viper.GetString("SERVER_PORT"),
	}
}
