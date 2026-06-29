package utils_jwt

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SecretKey string `envconfig:"SECRET_KEY" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("processing envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get jwt config: %w", err)
		panic(err)
	}

	return config
}
