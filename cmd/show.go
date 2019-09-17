package cmd

import (
	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdShow() *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show dosanco resources",
	}
	showCmd.AddCommand(
		NewCmdShowNetwork(),
		NewCmdShowIPAM(),
		NewCmdShowVlan(),
		NewCmdShowHost(),
		NewCmdShowDataCenter(),
		NewCmdShowDataCenterFloor(),
		NewCmdShowDataCenterHall(),
		NewCmdShowRackRow(),
		NewCmdShowRack(),
	)

	return showCmd
}

// NewCmdShowNetwork is subcommand represents show network resource.
func NewCmdShowNetwork() *cobra.Command {
	var tree bool
	var depth int
	var rfc bool
	var networkCmd = &cobra.Command{
		Use:     "network [CIDR]",
		Aliases: []string{"net", "nw"},
		Short:   "show network",
		Args:    cobra.MaximumNArgs(1),
		Run:     showNetwork,
	}
	networkCmd.Flags().BoolVarP(&tree, "tree", "t", false, "get network tree")
	networkCmd.Flags().BoolVarP(&rfc, "show-rfc-reserved", "", false, "show networks defined and reserved in RFC")
	networkCmd.Flags().IntVarP(&depth, "depth", "d", 0, "depth for network tree. this option only work with --tree,-t option")

	return networkCmd
}

// NewCmdShowIPAM is subcommand represents show ip allocation resource.
func NewCmdShowIPAM() *cobra.Command {
	var host bool
	var ipCmd = &cobra.Command{
		Use:     "ip [CIDR]",
		Aliases: []string{"ip-alloc"},
		Short:   "show ip allocation",
		Args:    cobra.ExactArgs(1),
		Run:     showIPAllocation,
	}
	ipCmd.Flags().BoolVarP(&host, "host", "", false, "use host name to get ip allocation")

	return ipCmd
}

// NewCmdShowVlan is subcommand represents show vlan resource.
func NewCmdShowVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vl"},
		Short:   "show vlan",
		Run: func(cmd *cobra.Command, args []string) {
			url := Conf.APIServer.URL + "/vlan"
			getVlans(url)
		},
	}

	return vlanCmd
}

// NewCmdShowHost is subcommand represents show host resource.
func NewCmdShowHost() *cobra.Command {
	var hostCmd = &cobra.Command{
		Use:     "host [NAME]",
		Aliases: []string{"server"},
		Short:   "show host",
		Args:    cobra.ExactArgs(1),
		Run:     showHost,
	}

	return hostCmd
}

// NewCmdShowDataCenter is subcommand represents show datacenter resource.
func NewCmdShowDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "show datacenter",
		Args:    cobra.MaximumNArgs(1),
		Run:     getDataCenter,
	}

	return dcCmd
}

// NewCmdShowDataCenterFloor is subcommand represents show datacenter resource.
func NewCmdShowDataCenterFloor() *cobra.Command {
	var dc string
	var dcCmd = &cobra.Command{
		Use:     "floor",
		Aliases: []string{"dc-floor"},
		Short:   "show datacenter floor",
		Args:    cobra.MaximumNArgs(1),
		Run:     getDataCenterFloor,
	}
	dcCmd.Flags().StringVarP(&dc, "dc", "", "", "specify datacenter")

	return dcCmd
}

// NewCmdShowDataCenterHall is subcommand represents show datacenter resource.
func NewCmdShowDataCenterHall() *cobra.Command {
	var dc string
	var floor string
	var dcCmd = &cobra.Command{
		Use:     "hall",
		Aliases: []string{"dc-hall"},
		Short:   "show datacenter hall",
		Args:    cobra.MaximumNArgs(1),
		Run:     getDataCenterHall,
	}
	dcCmd.Flags().StringVarP(&dc, "dc", "", "", "specify datacenter")
	dcCmd.Flags().StringVarP(&floor, "floor", "", "", "specify datacenter floor")

	return dcCmd
}

// NewCmdShowRackRow is subcommand represents show rack row resource.
func NewCmdShowRackRow() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var rowCmd = &cobra.Command{
		Use:     "row [ROW_NAME]",
		Aliases: []string{"rack-row"},
		Short:   "show row of racks",
		Args:    cobra.MaximumNArgs(1),
		Run:     getRackRow,
	}
	rowCmd.Flags().StringVarP(&dc, "dc", "", "", "specify datacenter")
	rowCmd.Flags().StringVarP(&floor, "floor", "", "", "specify datacenter floor")
	rowCmd.Flags().StringVarP(&hall, "hall", "", "", "specify datacenter hall")
	rowCmd.MarkFlagRequired("dc")

	return rowCmd
}

// NewCmdShowRack is subcommand represents show rack row resource.
func NewCmdShowRack() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var row string
	var rowCmd = &cobra.Command{
		Use:     "rack [RACK_NAME]",
		Aliases: []string{""},
		Short:   "show rack",
		Args:    cobra.MaximumNArgs(1),
		Run:     getRack,
	}
	rowCmd.Flags().StringVarP(&dc, "dc", "", "", "specify datacenter")
	rowCmd.Flags().StringVarP(&floor, "floor", "", "", "specify datacenter floor")
	rowCmd.Flags().StringVarP(&hall, "hall", "", "", "specify datacenter hall")
	rowCmd.Flags().StringVarP(&row, "row", "", "", "specify rack row")
	rowCmd.MarkFlagRequired("dc")

	return rowCmd
}
