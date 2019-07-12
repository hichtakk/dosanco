package config

import (
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

	return conf, nil
}