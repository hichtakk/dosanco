package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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

func NewCmdDeleteNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net", "nw"},
		Short:   "delete network description",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires one network id")
			}
			return nil
		},
		RunE: deleteNetwork,
	}

	return networkCmd
}

func NewCmdDeleteIPAllocation() *cobra.Command {
	var ipamCmd = &cobra.Command{
		Use:   "ipam",
		Short: "delete new ip allocation",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires allocation ID")
			}
			return nil
		},
		RunE: deleteIPAllocation,
	}

	return ipamCmd
}

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