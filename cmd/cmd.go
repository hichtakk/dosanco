package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "v0.2.0"

var revision = ""

// Cobra flags
var (
	rootCmd    *cobra.Command
	debug      bool
	output     string
	flagConfig string
)

// Conf is global variable for configuration data.
var Conf Config

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//fmt.Printf("\033[31m" + err.Error() + "\033[0m\n")
		fmt.Printf("\033[31m" + err.Error() + "\n")
		os.Exit(255)
		fmt.Printf("\033[0m\n")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd = newRootCmd()
	rootCmd.AddCommand(
		NewCmdVersion(),
		NewCmdCreate(),
		NewCmdDelete(),
		NewCmdShow(),
		NewCmdUpdate(),
		NewCmdConfig(),
	)
	homedir := os.Getenv("HOME")
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", homedir+"/.dosanco.json", "configuration file")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "default", "output style [default,json]")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "debug output")
}

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dosanco",
		Short: "dosanco command-line client",
		Long:  "Dosanco: Simple infrastructure database \U0001f434",
	}
}
