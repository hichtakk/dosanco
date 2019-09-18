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
		NewCmdDeleteRackRow(),
		NewCmdDeleteRack(),
		NewCmdDeleteUPS(),
		NewCmdDeletePDU(),
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

// NewCmdRackRow is subcommand represents delete rack row resource.
func NewCmdDeleteRackRow() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var rowCmd = &cobra.Command{
		Use:     "row",
		Aliases: []string{"rack-row"},
		Short:   "delete rack row",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteRackRow,
	}
	rowCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	rowCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	rowCmd.Flags().StringVarP(&hall, "hall", "", "", "name of datacenter hall")
	rowCmd.MarkFlagRequired("dc")
	rowCmd.MarkFlagRequired("floor")
	rowCmd.MarkFlagRequired("hall")

	return rowCmd
}

// NewCmdDeleteRack is subcommand represents delete rack row resource.
func NewCmdDeleteRack() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var row string
	var rackCmd = &cobra.Command{
		Use:     "rack [RACK_NAME]",
		Aliases: []string{"rack-row"},
		Short:   "delete rack row",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteRack,
	}
	rackCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	rackCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor")
	rackCmd.Flags().StringVarP(&hall, "hall", "", "", "name of datacenter hall")
	rackCmd.Flags().StringVarP(&row, "row", "", "", "name of rack row")
	rackCmd.MarkFlagRequired("dc")
	rackCmd.MarkFlagRequired("floor")
	rackCmd.MarkFlagRequired("hall")
	rackCmd.MarkFlagRequired("row")

	return rackCmd
}

// NewCmdDeleteUPS is subcommand represents delete rack row resource.
func NewCmdDeleteUPS() *cobra.Command {
	var dc string
	var upsCmd = &cobra.Command{
		Use:   "ups [RACK_NAME]",
		Short: "delete ups",
		Args:  cobra.ExactArgs(1),
		RunE:  deleteUPS,
	}
	upsCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	upsCmd.MarkFlagRequired("dc")

	return upsCmd
}

// NewCmdDeletePDU is subcommand represents delete rack row resource.
func NewCmdDeletePDU() *cobra.Command {
	var dc string
	var pduCmd = &cobra.Command{
		Use:   "dc-pdu [PDU_NAME]",
		Short: "delete pdu",
		Args:  cobra.ExactArgs(1),
		RunE:  deletePDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter")
	pduCmd.MarkFlagRequired("dc")

	return pduCmd
}
