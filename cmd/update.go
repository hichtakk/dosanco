package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdUpdate() *cobra.Command {
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
	}
	updateCmd.AddCommand(
		NewCmdCreateNetwork(),
	)

	return updateCmd
}