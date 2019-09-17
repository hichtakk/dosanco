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
		NewCmdCreateHost(),
		NewCmdCreateDataCenter(),
		NewCmdCreateDataCenterFloor(),
		NewCmdCreateDataCenterHall(),
		NewCmdCreateRackRow(),
		NewCmdCreateRack(),
		NewCmdCreateUPS(),
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
	var allocType string
	var ipCmd = &cobra.Command{
		Use:     "ip [ADDRESS]",
		Aliases: []string{"ip-alloc"},
		Short:   "create new ip allocation",
		Args:    cobra.ExactArgs(1),
		RunE:    createIPAllocation,
	}
	ipCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested ip allocation")
	ipCmd.Flags().StringVarP(&name, "name", "", "", "hostname for the ip address")
	ipCmd.Flags().StringVarP(&network, "network", "", "", "network CIDR for the ip address")
	ipCmd.Flags().StringVarP(&allocType, "type", "t", "generic", "type of address. use 'reserved' or 'generic'")
	ipCmd.MarkFlagRequired("network")

	return ipCmd
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

// NewCmdCreateHost is subcommand represents vlan resource.
func NewCmdCreateHost() *cobra.Command {
	var location string
	var hostCmd = &cobra.Command{
		Use:     "host",
		Aliases: []string{"server"},
		Short:   "create new host",
		Args:    cobra.ExactArgs(1),
		RunE:    createHost,
	}
	hostCmd.Flags().StringVarP(&location, "location", "l", "", "location of the requested host")
	hostCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")

	return hostCmd
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

// NewCmdCreateDataCenterFloor is subcommand represents datacenter resource.
func NewCmdCreateDataCenterFloor() *cobra.Command {
	var dc string
	var flrCmd = &cobra.Command{
		Use:     "floor",
		Aliases: []string{"dc-floor", "area"},
		Short:   "create new floor to datacenter",
		Args:    cobra.ExactArgs(1),
		RunE:    createDataCenterFloor,
	}
	flrCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	flrCmd.MarkFlagRequired("dc")

	return flrCmd
}

// NewCmdCreateDataCenterHall is subcommand represents datacenter resource.
func NewCmdCreateDataCenterHall() *cobra.Command {
	var dc string
	var floor string
	var hallType string
	var hallCmd = &cobra.Command{
		Use:     "hall",
		Aliases: []string{"dc-hall"},
		Short:   "create new hall to datacenter floor",
		Args:    cobra.ExactArgs(1),
		RunE:    createDataCenterHall,
	}
	hallCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	hallCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	hallCmd.Flags().StringVarP(&hallType, "type", "", "", "type of data hall (network/generic)")
	hallCmd.MarkFlagRequired("dc")
	hallCmd.MarkFlagRequired("floor")
	hallCmd.MarkFlagRequired("type")

	return hallCmd
}

// NewCmdCreateRackRow is subcommand represents datacenter resource.
func NewCmdCreateRackRow() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var rowCmd = &cobra.Command{
		Use:     "row",
		Aliases: []string{"rack-row"},
		Short:   "create new hall to datacenter floor",
		Args:    cobra.ExactArgs(1),
		RunE:    createRackRow,
	}
	rowCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	rowCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	rowCmd.Flags().StringVarP(&hall, "hall", "", "", "name of data hall")
	rowCmd.MarkFlagRequired("dc")
	rowCmd.MarkFlagRequired("floor")
	rowCmd.MarkFlagRequired("hall")

	return rowCmd
}

// NewCmdCreateRack is subcommand represents datacenter resource.
func NewCmdCreateRack() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var row string
	var rackCmd = &cobra.Command{
		Use:     "rack [RACK_NAME]",
		Aliases: []string{""},
		Short:   "create new rack to row",
		Args:    cobra.ExactArgs(1),
		RunE:    createRack,
	}
	rackCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	rackCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	rackCmd.Flags().StringVarP(&hall, "hall", "", "", "name of data hall")
	rackCmd.Flags().StringVarP(&row, "row", "", "", "name of row")
	rackCmd.MarkFlagRequired("dc")
	rackCmd.MarkFlagRequired("floor")
	rackCmd.MarkFlagRequired("hall")
	rackCmd.MarkFlagRequired("row")

	return rackCmd
}

// NewCmdCreateUPS is subcommand represents vlan resource.
func NewCmdCreateUPS() *cobra.Command {
	var dc string
	var upsCmd = &cobra.Command{
		Use:   "ups",
		Short: "create new ups",
		Args:  cobra.ExactArgs(1),
		RunE:  createUPS,
	}
	upsCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	upsCmd.MarkFlagRequired("dc")

	return upsCmd
}
