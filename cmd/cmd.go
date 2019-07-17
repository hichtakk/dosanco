package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

/*
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hello",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//    Run: func(cmd *cobra.Command, args []string) { },
}
*/
var rootCmd *cobra.Command
var verbose bool
var output string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd = newRootCmd()
	rootCmd.AddCommand(
		NewCmdVersion(),
		NewCmdCreate(),
		NewCmdDelete(),
		NewCmdShow(),
		NewCmdUpdate(),
	)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "default", "output style [default,json]")
}

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dcctl",
		Short: "This is dcctl command",
		Long:  `DCTL: dosanco command-line client`,
	}
}

/*
func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
*/
