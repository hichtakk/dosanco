package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdCreate() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a resource",
	}
	createCmd.AddCommand(
		NewCmdCreateNetwork(),
		NewCmdCreateIPAllocation(),
		NewCmdCreateVlan(),
		NewCmdCreateDataCenter(),
	)

	return createCmd
}
