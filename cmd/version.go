package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdVersion() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
		  fmt.Println("dosanco command-line client v0.0.1")
		},
	}

	return versionCmd
}