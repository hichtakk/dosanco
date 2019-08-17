package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdUpdate() *cobra.Command {
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a resource",
	}
	updateCmd.AddCommand(
		NewCmdUpdateNetwork(),
		NewCmdUpdateIPAllocation(),
		NewCmdUpdateVlan(),
		NewCmdUpdateDataCenter(),
	)

	return updateCmd
}

func NewCmdUpdateNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net", "nw"},
		Short:   "update network description",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires one network id")
			}
			return nil
		},
		RunE: updateNetwork,
	}
	networkCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested network")

	return networkCmd
}

func NewCmdUpdateIPAllocation() *cobra.Command {
	var ipamCmd = &cobra.Command{
		Use:   "ipam",
		Short: "update ip allocation data",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires allocation ID")
			}
			return nil
		},
		RunE: updateIPAllocation,
	}
	ipamCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested ip allocation")
	//ipamCmd.Flags().StringVarP(&hostname, "hostname", "name", "", "ip address of the requested allocation")

	return ipamCmd
}

func NewCmdUpdateVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vlan"},
		Short:   "update vlan description",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires vlan id")
			}
			return nil
		},
		RunE: updateVlan,
	}
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")

	return vlanCmd
}

func NewCmdUpdateDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "update datacenter address",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires datacenter id")
			}
			return nil
		},
		RunE: updateDataCenter,
	}
	dcCmd.Flags().StringVarP(&address, "address", "a", "", "address of the datacenter")
	dcCmd.MarkFlagRequired("address")

	return dcCmd
}