package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmdDelete is subcommand to delete resources.
func NewCmdDelete() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource",
	}
	deleteCmd.AddCommand(
		NewCmdDeleteNetwork(),
		NewCmdDeleteIPAllocation(),
		NewCmdDeleteVlan(),
		NewCmdDeleteDataCenter(),
	)

	return deleteCmd
}

// NewCmdDeleteNetwork is subcommand represents delete network resource.
func NewCmdDeleteNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network [CIDR]",
		Aliases: []string{"net", "nw"},
		Short:   "delete network description",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteNetwork,
	}

	return networkCmd
}

// NewCmdDeleteIPAllocation is subcommand represents delete ip allocation resource.
func NewCmdDeleteIPAllocation() *cobra.Command {
	var ipCmd = &cobra.Command{
		Use:     "ipam [ADDRESS]",
		Aliases: []string{"ip"},
		Short:   "delete new ip allocation",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteIPAllocation,
	}

	return ipCmd
}

// NewCmdDeleteVlan is subcommand represents delete vlan resource.
func NewCmdDeleteVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vlan"},
		Short:   "delete vlan description",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires vlan id")
			}
			return nil
		},
		RunE: deleteVlan,
	}

	return vlanCmd
}

// NewCmdDeleteDataCenter is subcommand represents delete datacenter resource.
func NewCmdDeleteDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "delete datacenter",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires datacenter id")
			}
			return nil
		},
		RunE: deleteDataCenter,
	}

	return dcCmd
}
