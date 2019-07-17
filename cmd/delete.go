package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDelete() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
	}
	deleteCmd.AddCommand(
		NewCmdCreateNetwork(),
	)

	return deleteCmd
}