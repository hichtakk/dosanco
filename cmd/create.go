package cmd

import (
	"github.com/spf13/cobra"
)

// NewCmdCreate is subcommand to create resources.
func NewCmdCreate() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "create [RESOURCE]",
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

// NewCmdCreateNetwork is subcommand represents network resource.
func NewCmdCreateNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network [CIDR]",
		Aliases: []string{"net", "nw"},
		Short:   "create new network",
		Long:    "create new network",
		Args:    cobra.ExactArgs(1),
		RunE:    createNetwork,
	}
	networkCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested network")

	return networkCmd
}

// NewCmdCreateIPAllocation is subcommand represents ip allocation resource.
func NewCmdCreateIPAllocation() *cobra.Command {
	var name string
	var network string
	var ipamCmd = &cobra.Command{
		Use:     "ipam [ADDRESS]",
		Aliases: []string{"ip"},
		Short:   "create new ip allocation",
		Args:    cobra.ExactArgs(1),
		RunE:    createIPAllocation,
	}
	ipamCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested ip allocation")
	ipamCmd.Flags().StringVarP(&name, "name", "", "", "hostname for the ip address")
	ipamCmd.Flags().StringVarP(&network, "network", "", "", "network CIDR for the ip address")
	ipamCmd.MarkFlagRequired("network")

	return ipamCmd
}

// NewCmdCreateVlan is subcommand represents vlan resource.
func NewCmdCreateVlan() *cobra.Command {
	var networkID int
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vlan"},
		Short:   "create new vlan",
		Args:    cobra.ExactArgs(1),
		RunE:    createVlan,
	}
	vlanCmd.Flags().IntVarP(&networkID, "network-id", "n", 0, "network id of the requested ip allocation")
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")
	vlanCmd.MarkFlagRequired("network-id")

	return vlanCmd
}

// NewCmdCreateDataCenter is subcommand represents datacenter resource.
func NewCmdCreateDataCenter() *cobra.Command {
	var address string
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "create new datacenter",
		Args:    cobra.ExactArgs(1),
		RunE:    createDataCenter,
	}
	dcCmd.Flags().StringVarP(&address, "address", "a", "", "address of data center")
	dcCmd.MarkFlagRequired("address")

	return dcCmd
}
