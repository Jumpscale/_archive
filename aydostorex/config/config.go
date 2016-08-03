package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	ListenAddr       string
	StoreRoot        string
	Authentification string
	Tls              TLS
	Sources          []string
}

type TLS struct {
	Certificate string
	Key         string
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := toml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
