package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP       `yaml:"http"`
		DB         `yaml:"db"`
		SubReddits `yaml:"subreddits"`
	}

	HTTP struct {
		Addr string `env-required:"true" yaml:"address" env:"HTTP_ADDR"`
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	DB struct {
		DBFile string `env-required:"true" yaml:"db_file" env:"DB_FILE"`
	}

	SubReddits struct {
		SRFile string `env-required:"true" yaml:"sr_file" env:"SR_FILE"`
	}
)

// NewConfig returns app config.
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
