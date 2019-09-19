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
		NewCmdUpdateHost(),
		NewCmdUpdateDataCenter(),
		NewCmdUpdateDataCenterFloor(),
		NewCmdUpdateDataCenterHall(),
		NewCmdUpdateRackRow(),
		NewCmdUpdateRack(),
		NewCmdUpdateUPS(),
		NewCmdUpdatePDU(),
		NewCmdUpdateRackPDU(),
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
	var name string
	var description string
	var ipCmd = &cobra.Command{
		Use:     "ip [ADDRESS]",
		Aliases: []string{"ip-alloc"},
		Short:   "update ip allocation data",
		Args:    cobra.ExactArgs(1),
		RunE:    updateIPAllocation,
	}
	ipCmd.Flags().StringVarP(&description, "description", "d", "-", "new description of the requested ip allocation")
	ipCmd.Flags().StringVarP(&name, "name", "n", "-", "new hostname of the allocation")

	return ipCmd
}

// NewCmdUpdateVlan is subcommand represents update vlan resource.
func NewCmdUpdateVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:   "vlan [VLAN_ID]",
		Short: "update vlan description",
		Args:  cobra.ExactArgs(1),
		RunE:  updateVlan,
	}
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")

	return vlanCmd
}

// NewCmdUpdateHost is subcommand represents update host resource.
func NewCmdUpdateHost() *cobra.Command {
	var name string
	var location string
	var hostCmd = &cobra.Command{
		Use:     "host [NAME]",
		Aliases: []string{"server"},
		Short:   "update host information",
		Args:    cobra.ExactArgs(1),
		RunE:    updateHost,
	}
	hostCmd.Flags().StringVarP(&name, "name", "n", "-", "name of the requested host")
	hostCmd.Flags().StringVarP(&location, "location", "l", "-", "location of the requested host")
	hostCmd.Flags().StringVarP(&description, "description", "d", "-", "description of the requested vlan")

	return hostCmd
}

// NewCmdUpdateDataCenter is subcommand represents update datacenter resource.
func NewCmdUpdateDataCenter() *cobra.Command {
	var address string
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

// NewCmdUpdateDataCenterFloor is subcommand represents update datacenter floor/area resource.
func NewCmdUpdateDataCenterFloor() *cobra.Command {
	var dc string
	var name string
	var flrCmd = &cobra.Command{
		Use:     "floor",
		Aliases: []string{"dc-floor", "area"},
		Short:   "update datacenter address",
		Args:    cobra.ExactArgs(1),
		RunE:    updateDataCenterFloor,
	}
	flrCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	flrCmd.Flags().StringVarP(&name, "name", "n", "-", "name of datacenter floor")
	flrCmd.MarkFlagRequired("dc")
	flrCmd.MarkFlagRequired("name")

	return flrCmd
}

// NewCmdUpdateDataCenterHall is subcommand represents update datacenter hall resource.
func NewCmdUpdateDataCenterHall() *cobra.Command {
	var dc string
	var floor string
	var name string
	var hallCmd = &cobra.Command{
		Use:     "hall",
		Aliases: []string{"dc-hall"},
		Short:   "update data hall name",
		Args:    cobra.ExactArgs(1),
		RunE:    updateDataCenterHall,
	}
	hallCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	hallCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	hallCmd.Flags().StringVarP(&name, "name", "n", "-", "new name of data hall")
	hallCmd.MarkFlagRequired("dc")
	hallCmd.MarkFlagRequired("floor")
	hallCmd.MarkFlagRequired("name")

	return hallCmd
}

// NewCmdUpdateRackRow is subcommand represents update rack row resource.
func NewCmdUpdateRackRow() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var name string
	var rowCmd = &cobra.Command{
		Use:     "row",
		Aliases: []string{"rack-row"},
		Short:   "update rack row name",
		Args:    cobra.ExactArgs(1),
		RunE:    updateRackRow,
	}
	rowCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	rowCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	rowCmd.Flags().StringVarP(&hall, "hall", "", "", "name of datacenter hall")
	rowCmd.Flags().StringVarP(&name, "name", "n", "-", "new name of rack row")
	rowCmd.MarkFlagRequired("dc")
	rowCmd.MarkFlagRequired("floor")
	rowCmd.MarkFlagRequired("hall")
	rowCmd.MarkFlagRequired("name")

	return rowCmd
}

// NewCmdUpdateRack is subcommand represents update rack resource.
func NewCmdUpdateRack() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var row string
	var name string
	var rackCmd = &cobra.Command{
		Use:   "rack [RACK_NAME]",
		Short: "update rack name",
		Args:  cobra.ExactArgs(1),
		RunE:  updateRack,
	}
	rackCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	rackCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	rackCmd.Flags().StringVarP(&hall, "hall", "", "", "name of datacenter hall")
	rackCmd.Flags().StringVarP(&row, "row", "", "", "name of rack row")
	rackCmd.Flags().StringVarP(&name, "name", "n", "-", "new name of rack")
	rackCmd.MarkFlagRequired("dc")
	rackCmd.MarkFlagRequired("floor")
	rackCmd.MarkFlagRequired("hall")
	rackCmd.MarkFlagRequired("row")
	rackCmd.MarkFlagRequired("name")

	return rackCmd
}

// NewCmdUpdateUPS is subcommand represents update rack resource.
func NewCmdUpdateUPS() *cobra.Command {
	var dc string
	var name string
	var upsCmd = &cobra.Command{
		Use:   "ups [UPS_NAME]",
		Short: "update ups name",
		Args:  cobra.ExactArgs(1),
		RunE:  updateUPS,
	}
	upsCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	upsCmd.Flags().StringVarP(&name, "name", "n", "-", "new name of rack")
	upsCmd.MarkFlagRequired("dc")
	upsCmd.MarkFlagRequired("name")

	return upsCmd
}

// NewCmdUpdatePDU is subcommand represents update rack resource.
func NewCmdUpdatePDU() *cobra.Command {
	var dc string
	var name string
	var pduCmd = &cobra.Command{
		Use:   "dc-pdu [PDU_NAME]",
		Short: "update pdu name",
		Args:  cobra.ExactArgs(1),
		RunE:  updatePDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	pduCmd.Flags().StringVarP(&name, "name", "n", "-", "new name of rack")
	pduCmd.MarkFlagRequired("dc")
	pduCmd.MarkFlagRequired("name")

	return pduCmd
}

// NewCmdUpdateRackPDU is subcommand represents update rack resource.
func NewCmdUpdateRackPDU() *cobra.Command {
	var dc string
	var name string
	var pduCmd = &cobra.Command{
		Use:   "rack-pdu [RACK_PDU_NAME]",
		Short: "update rack pdu name",
		Args:  cobra.ExactArgs(1),
		RunE:  updateRackPDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	pduCmd.Flags().StringVarP(&name, "name", "n", "-", "new name of rack")
	pduCmd.MarkFlagRequired("dc")
	pduCmd.MarkFlagRequired("name")

	return pduCmd
}
