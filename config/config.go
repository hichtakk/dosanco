package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Server is apiserver daemon settings
type Server struct {
	DatabaseURL string `toml:"db_url"`
}

// Feature is struct for enabled dosanco features.
type Feature struct {
	Network    bool `toml:"network"`
	Host       bool `toml:"host"`
	DataCenter bool `toml:"datacenter"`
}

// Config is struct for main configuration.
type Config struct {
	Server  Server  `toml:"server"`
	Feature Feature `toml:"feature"`
}

// NewConfig returns configuration instance.
func NewConfig(path string) (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return conf, err
	}

	// override configuration with environmental variables
	if envPath := os.Getenv("DOSANCO_DB"); envPath != "" {
		conf.Server.DatabaseURL = envPath
	}

	return conf, nil
}
