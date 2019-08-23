package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// APIServerConfig represents api server endpoint.
type APIServerConfig struct {
	URL string `toml:"url"`
}

// Config is struct for main configuration.
type Config struct {
	APIServer APIServerConfig `toml:"apiserver"`
}

func initConfig() {
	if flagConfig != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flagConfig)
	} else {
		// Search config in home directory
		homedir := os.Getenv("HOME")
		viper.AddConfigPath(homedir)
		viper.SetConfigName(".dosanco.toml")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
