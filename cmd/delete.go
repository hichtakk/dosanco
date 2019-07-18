package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDelete() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource",
	}
	deleteCmd.AddCommand(
		NewCmdCreateNetwork(),
	)

	return deleteCmd
}
