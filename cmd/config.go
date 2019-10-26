package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// APIServerConfig represents api server endpoint.
type APIServerConfig struct {
	URL string `json:"url"`
}

// Config is struct for main configuration.
type Config struct {
	APIServer APIServerConfig `json:"apiserver"`
}

func initConfig() {
	viper.SetConfigType("json")
	if flagConfig != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flagConfig)
	} else {
		// Search config in home directory
		homedir := os.Getenv("HOME")
		viper.AddConfigPath(homedir)
		viper.SetConfigName(".dosanco.json")
	}
	if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
		// set default configuration
		viper.Set("apiserver.url", "http://localhost:15187")
		viper.WriteConfig()
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("read config error")
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewCmdConfig() *cobra.Command {
	var cfgCmd = &cobra.Command{
		Use:   "config [SUBCOMMAND]",
		Short: "Modify dosanco client config file",
	}
	cfgCmd.AddCommand(
		NewCmdConfigView(),
		NewCmdConfigSetEndpoint(),
	)

	return cfgCmd
}

// NewCmdConfigShow is subcommand to display dosanco client configuration
func NewCmdConfigView() *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "view",
		Short: "set apiserver endpoint",
		Run:   configView,
	}

	return showCmd
}

// NewCmdConfigSetEndpoint is subcommand to set apiserver endpoint configuration
func NewCmdConfigSetEndpoint() *cobra.Command {
	var epCmd = &cobra.Command{
		Use:   "set-endpoint [ENDPOINT_URL]",
		Short: "set apiserver endpoint",
		Args:  cobra.ExactArgs(1),
		Run:   configSetEndpoint,
	}

	return epCmd
}

func configView(_ *cobra.Command, _ []string) {
	confPath := viper.ConfigFileUsed()
	file, _ := os.Open(confPath)
	defer file.Close()
	io.Copy(os.Stdout, file)
}

func configSetEndpoint(cmd *cobra.Command, args []string) {
	viper.Set("apiserver.url", args[0])
	if err := viper.WriteConfig(); err != nil {
		fmt.Println(err)
		os.Exit(255)
	}
	fmt.Println("endpoint was updated to " + args[0])
}
