package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// DBConfig is struct for database configuration.
type DBConfig struct {
	URL  string `toml:"url"`
	Type string `toml:"type"`
	Host string `toml:"host"`
	//User     string `toml:"user"`
	//Password string `toml:"password"`
	Port int    `toml:"port"`
	Name string `toml:"name"`
}

// Feature is struct for enabled dosanco features.
type Feature struct {
	Network    bool `toml:"network"`
	Host       bool `toml:"host"`
	DataCenter bool `toml:"datacenter"`
}

// Config is struct for main configuration.
type Config struct {
	DB      DBConfig `toml:"database"`
	Feature Feature  `toml:"feature"`
}

// NewConfig returns configuration instance.
func NewConfig(path string) (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return conf, err
	}

	// override configuration with environmental variables
	if envPath := os.Getenv("DOSANCO_DB"); envPath != "" {
		conf.DB.URL = envPath
	}

	return conf, nil
}
