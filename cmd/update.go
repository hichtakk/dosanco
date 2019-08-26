package cmd

import (
	"github.com/spf13/cobra"
)

// NewCmdUpdate is subcommand to update resources.
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

// NewCmdUpdateNetwork is subcommand represents update network resource.
func NewCmdUpdateNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network [CIDR]",
		Aliases: []string{"net", "nw"},
		Short:   "update network description",
		Args:    cobra.ExactArgs(1),
		RunE:    updateNetwork,
	}
	networkCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested network")

	return networkCmd
}

// NewCmdUpdateIPAllocation is subcommand represents update ip allocation resource.
func NewCmdUpdateIPAllocation() *cobra.Command {
	//var name string
	var ipamCmd = &cobra.Command{
		Use:     "ipam [ADDRESS]",
		Aliases: []string{"ip"},
		Short:   "update ip allocation data",
		Args:    cobra.ExactArgs(1),
		RunE:    updateIPAllocation,
	}
	ipamCmd.Flags().StringVarP(&description, "description", "d", "", "new description of the requested ip allocation")
	//ipamCmd.Flags().StringVarP(&name, "name", "", "", "new hostname of the requested allocation")

	return ipamCmd
}

// NewCmdUpdateVlan is subcommand represents update vlan resource.
func NewCmdUpdateVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan [VLAN_ID]",
		Aliases: []string{"vlan"},
		Short:   "update vlan description",
		Args:    cobra.ExactArgs(1),
		RunE:    updateVlan,
	}
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")

	return vlanCmd
}

// NewCmdUpdateDataCenter is subcommand represents update datacenter resource.
func NewCmdUpdateDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "update datacenter address",
		Args:    cobra.ExactArgs(1),
		RunE:    updateDataCenter,
	}
	dcCmd.Flags().StringVarP(&address, "address", "a", "", "address of the datacenter")
	dcCmd.MarkFlagRequired("address")

	return dcCmd
}
