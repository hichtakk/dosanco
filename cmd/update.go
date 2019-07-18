package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdUpdate() *cobra.Command {
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a resource",
	}
	updateCmd.AddCommand(
		NewCmdCreateNetwork(),
	)

	return updateCmd
}
