package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmdVersion is subcommand to show version information.
func NewCmdVersion() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of dosanco client",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dosanco command-line client v0.0.1")
		},
	}

	return versionCmd
}
