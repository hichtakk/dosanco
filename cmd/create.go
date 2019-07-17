package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdCreate() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Print the version number of Hugo",
		//Long:  `All software has versions. This is Hugo's`,
	}
	createCmd.AddCommand(
		NewCmdCreateNetwork(),
	)

	return createCmd
}
