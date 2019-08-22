package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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

func NewCmdCreateNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network [CIDR]",
		Aliases: []string{"net", "nw"},
		Short:   "create new network",
		Long:    "create new network",
		Args:    cobra.ExactArgs(1),
		RunE:    createNetwork,
	}
	networkCmd.Flags().IntVarP(&supernetID, "supernet-id", "s", 0, "supernetwork id of the requested network")
	networkCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested network")
	networkCmd.MarkFlagRequired("supernet-id")

	return networkCmd
}

func NewCmdCreateIPAllocation() *cobra.Command {
	var ipamCmd = &cobra.Command{
		Use:   "ipam",
		Short: "create new ip allocation",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires hostname")
			}
			return nil
		},
		RunE: createIPAllocation,
	}
	ipamCmd.Flags().IntVarP(&networkID, "network-id", "n", 0, "network id of the requested ip allocation")
	ipamCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested ip allocation")
	ipamCmd.Flags().StringVarP(&address, "address", "a", "", "ip address of the requested allocation")
	ipamCmd.MarkFlagRequired("network-id")
	ipamCmd.MarkFlagRequired("address")

	return ipamCmd
}

func NewCmdCreateVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vlan"},
		Short:   "create new vlan",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires vlan id")
			}
			return nil
		},
		RunE: createVlan,
	}
	vlanCmd.Flags().IntVarP(&networkID, "network-id", "n", 0, "network id of the requested ip allocation")
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")
	vlanCmd.MarkFlagRequired("network-id")

	return vlanCmd
}

func NewCmdCreateDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "create new datacenter",
		//Long:    "create new network",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires data center name")
			}
			return nil
		},
		RunE: createDataCenter,
	}
	dcCmd.Flags().StringVarP(&address, "address", "a", "", "address of data center")
	dcCmd.MarkFlagRequired("address")

	return dcCmd
}
