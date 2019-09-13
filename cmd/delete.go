package cmd

import (
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
		NewCmdDeleteHost(),
		NewCmdDeleteDataCenter(),
		NewCmdDeleteDataCenterFloor(),
		NewCmdDeleteDataCenterHall(),
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
		Use:     "ip [ADDRESS]",
		Aliases: []string{"ip-alloc"},
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
		Short:   "delete vlan",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteVlan,
	}

	return vlanCmd
}

// NewCmdDeleteHost is subcommand represents delete vlan resource.
func NewCmdDeleteHost() *cobra.Command {
	var hostCmd = &cobra.Command{
		Use:     "host [NAME]",
		Aliases: []string{"server"},
		Short:   "delete host",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteHost,
	}

	return hostCmd
}

// NewCmdDeleteDataCenter is subcommand represents delete datacenter resource.
func NewCmdDeleteDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "delete datacenter",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteDataCenter,
	}

	return dcCmd
}

// NewCmdDeleteDataCenterFloor is subcommand represents delete datacenter floor resource.
func NewCmdDeleteDataCenterFloor() *cobra.Command {
	var dc string
	var flrCmd = &cobra.Command{
		Use:     "floor",
		Aliases: []string{"dc-floor"},
		Short:   "delete datacenter floor",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteDataCenterFloor,
	}
	flrCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	flrCmd.MarkFlagRequired("dc")

	return flrCmd
}

// NewCmdDeleteDataCenterHall is subcommand represents delete datacenter hall resource.
func NewCmdDeleteDataCenterHall() *cobra.Command {
	var dc string
	var floor string
	var hallCmd = &cobra.Command{
		Use:     "hall",
		Aliases: []string{"dc-hall"},
		Short:   "delete datacenter hall",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteDataCenterHall,
	}
	hallCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	hallCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	hallCmd.MarkFlagRequired("dc")
	hallCmd.MarkFlagRequired("floor")

	return hallCmd
}
