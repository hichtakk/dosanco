package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Cobra flags
var (
	rootCmd    *cobra.Command
	verbose    bool
	output     string
	flagConfig string
)

// Conf is global variable for configuration data.
var Conf Config

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
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
	)
	homedir := os.Getenv("HOME")
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", homedir+"/.dosanco.toml", "configuration file")
	//rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "default", "output style [default,json]")
}

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dosanco",
		Short: "dosanco command-line client",
		Long:  "dosanco controls Dosanco infrastructure database",
	}
}
