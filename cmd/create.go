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
		NewCmdCreateHostGroup(),
		NewCmdCreateDataCenter(),
		NewCmdCreateDataCenterFloor(),
		NewCmdCreateDataCenterHall(),
		NewCmdCreateRackRow(),
		NewCmdCreateRack(),
		NewCmdCreateUPS(),
		NewCmdCreatePDU(),
		NewCmdCreateRackPDU(),
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
	var cidr string
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vl"},
		Short:   "create new vlan",
		Args:    cobra.ExactArgs(1),
		RunE:    createVlan,
	}
	vlanCmd.Flags().StringVarP(&cidr, "cidr", "", "", "network cidr of the requested vlan")
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")
	vlanCmd.MarkFlagRequired("cidr")

	return vlanCmd
}

// NewCmdCreateHost is subcommand represents vlan resource.
func NewCmdCreateHost() *cobra.Command {
	var location string
	var group string
	var hostCmd = &cobra.Command{
		Use:     "host",
		Aliases: []string{"server"},
		Short:   "create new host",
		Args:    cobra.ExactArgs(1),
		RunE:    createHost,
	}
	hostCmd.Flags().StringVarP(&location, "location", "l", "", "location of host installed. use format '{DC}/{FLOOR}/{HALL}/{ROW}/{RACK}'")
	hostCmd.Flags().StringVarP(&description, "description", "d", "", "description of the host")
	hostCmd.Flags().StringVarP(&group, "group", "g", "", "group of the host")
	hostCmd.MarkFlagRequired("group")

	return hostCmd
}

// NewCmdCreateHostGroup is subcommand represents vlan resource.
func NewCmdCreateHostGroup() *cobra.Command {
	var groupCmd = &cobra.Command{
		Use:   "group [NAME]",
		Short: "create new host group",
		Args:  cobra.ExactArgs(1),
		RunE:  createHostGroup,
	}
	groupCmd.Flags().StringVarP(&description, "description", "d", "", "name of host group")

	return groupCmd
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
	//var hallType string
	var hallCmd = &cobra.Command{
		Use:     "hall",
		Aliases: []string{"dc-hall"},
		Short:   "create new hall to datacenter floor",
		Args:    cobra.ExactArgs(1),
		RunE:    createDataCenterHall,
	}
	hallCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	hallCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	//hallCmd.Flags().StringVarP(&hallType, "type", "", "", "type of data hall (network/generic)")
	hallCmd.MarkFlagRequired("dc")
	hallCmd.MarkFlagRequired("floor")
	//hallCmd.MarkFlagRequired("type")

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

// NewCmdCreateUPS is subcommand represents ups resource.
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

// NewCmdCreatePDU is subcommand represents pdu resource.
func NewCmdCreatePDU() *cobra.Command {
	var dc string
	var primary string
	var secondary string
	var pduCmd = &cobra.Command{
		Use:   "row-pdu",
		Short: "create new row-pdu",
		Args:  cobra.ExactArgs(1),
		RunE:  createPDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	pduCmd.Flags().StringVarP(&primary, "primary", "", "", "name of primary power source")
	pduCmd.Flags().StringVarP(&secondary, "secondary", "", "", "name of secondary power source")
	pduCmd.MarkFlagRequired("dc")
	pduCmd.MarkFlagRequired("primary")

	return pduCmd
}

// NewCmdCreateRackPDU is subcommand represents pdu resource.
func NewCmdCreateRackPDU() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var row string
	var rack string
	var primary string
	var secondary string
	var group string
	var pduCmd = &cobra.Command{
		Use:   "rack-pdu",
		Short: "create new rack-pdu",
		Long:  `create new rack-pdu dosanco create rack-pdu --dc DC1 --primary ROW-PDU-1 rack-pdu01.dosanco`,
		Args:  cobra.ExactArgs(1),
		RunE:  createRackPDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	pduCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	pduCmd.Flags().StringVarP(&hall, "hall", "", "", "name of data hall")
	pduCmd.Flags().StringVarP(&row, "row", "", "", "name of rack row")
	pduCmd.Flags().StringVarP(&rack, "rack", "", "", "name of rack")
	pduCmd.Flags().StringVarP(&primary, "primary", "", "", "name of primary power source")
	pduCmd.Flags().StringVarP(&secondary, "secondary", "", "", "name of secondary power source")
	pduCmd.Flags().StringVarP(&group, "group", "", "", "name of host group")
	pduCmd.MarkFlagRequired("dc")
	pduCmd.MarkFlagRequired("floor")
	pduCmd.MarkFlagRequired("hall")
	pduCmd.MarkFlagRequired("row")
	pduCmd.MarkFlagRequired("rack")
	pduCmd.MarkFlagRequired("group")
	pduCmd.MarkFlagRequired("primary")

	return pduCmd
}
