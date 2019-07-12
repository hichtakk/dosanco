package config

import (
	"os"
	"github.com/BurntSushi/toml"
)

type DBConfig struct {
	Path	string	`toml:"path"`
}

type Config struct {
	DB	DBConfig `toml:"database"`
}

func NewConfig(path string) (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return conf, err
	}

	// override configuration with environmental variables
	if envPath := os.Getenv("DOSANCO_DB"); envPath != "" {
		conf.DB.Path = os.Getenv("DOSANCO_DB")
	}

	return conf, nil
}