package cmd

import (
	"fmt"

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
		NewCmdDeleteHostGroup(),
		NewCmdDeleteDataCenter(),
		NewCmdDeleteDataCenterFloor(),
		NewCmdDeleteDataCenterHall(),
		NewCmdDeleteRackRow(),
		NewCmdDeleteRack(),
		NewCmdDeleteUPS(),
		NewCmdDeletePDU(),
		NewCmdDeleteRackPDU(),
	)

	return deleteCmd
}

// NewCmdDeleteNetwork is subcommand represents delete network resource.
func NewCmdDeleteNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network ${CIDR}",
		Aliases: []string{"net", "nw"},
		Short:   "delete network description",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("cidr style network is required 'XXX.XXX.XXX.XXX/XX'")
			}
			return nil
		},
		RunE: deleteNetwork,
	}

	return networkCmd
}

// NewCmdDeleteIPAllocation is subcommand represents delete ip allocation resource.
func NewCmdDeleteIPAllocation() *cobra.Command {
	var ipCmd = &cobra.Command{
		Use:     "ip ${ADDRESS}",
		Aliases: []string{"ip-alloc"},
		Short:   "delete new ip allocation",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("ip address is required 'XXX.XXX.XXX.XXX'")
			}
			return nil
		},
		RunE: deleteIPAllocation,
	}

	return ipCmd
}

// NewCmdDeleteVlan is subcommand represents delete vlan resource.
func NewCmdDeleteVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan ${VLAN_ID}",
		Aliases: []string{"vlan"},
		Short:   "delete vlan",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("vlan id is required")
			}
			return nil
		},
		RunE: deleteVlan,
	}

	return vlanCmd
}

// NewCmdDeleteHost is subcommand represents delete vlan resource.
func NewCmdDeleteHost() *cobra.Command {
	var hostCmd = &cobra.Command{
		Use:     "host ${NAME}",
		Aliases: []string{"server"},
		Short:   "delete host",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("hostname is required")
			}
			return nil
		},
		RunE: deleteHost,
	}

	return hostCmd
}

// NewCmdDeleteHostGroup is subcommand represents delete vlan resource.
func NewCmdDeleteHostGroup() *cobra.Command {
	var groupCmd = &cobra.Command{
		Use:   "group ${NAME}",
		Short: "delete host-group",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("group is required")
			}
			return nil
		},
		RunE: deleteHostGroup,
	}

	return groupCmd
}

// NewCmdDeleteDataCenter is subcommand represents delete datacenter resource.
func NewCmdDeleteDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter ${DC_NAME}",
		Aliases: []string{"dc"},
		Short:   "delete datacenter",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			return nil
		},
		RunE: deleteDataCenter,
	}

	return dcCmd
}

// NewCmdDeleteDataCenterFloor is subcommand represents delete datacenter floor resource.
func NewCmdDeleteDataCenterFloor() *cobra.Command {
	var dc string
	var flrCmd = &cobra.Command{
		Use:     "floor ${FLOOR_NAME} --dc ${DC_NAME}",
		Aliases: []string{"dc-floor"},
		Short:   "delete datacenter floor",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("floor name is required")
			}
			if dc == "" {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			return nil
		},
		RunE: deleteDataCenterFloor,
	}
	flrCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	flrCmd.MarkFlagRequired("dc")

	return flrCmd
}

// NewCmdDeleteDataCenterHall is subcommand represents delete datacenter hall resource.
func NewCmdDeleteDataCenterHall() *cobra.Command {
	var dc string
	var floor string
	var hallCmd = &cobra.Command{
		Use:     "hall ${HALL_NAME} --dc ${DC_NAME} --floor ${FLOOR_NAME}",
		Aliases: []string{"dc-hall"},
		Short:   "delete datacenter hall",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("hall name is required")
			}
			if dc == "" {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			if floor == "" {
				cmd.Help()
				return fmt.Errorf("floor name is required")
			}
			return nil
		},
		RunE: deleteDataCenterHall,
	}
	hallCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	hallCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor [REQUIRED]")
	hallCmd.MarkFlagRequired("dc")
	hallCmd.MarkFlagRequired("floor")

	return hallCmd
}

// NewCmdDeleteRackRow is subcommand represents delete rack row resource.
func NewCmdDeleteRackRow() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var rowCmd = &cobra.Command{
		Use:     "row ${ROW_NAME} --dc ${DC_NAME} --floor ${FLOOR_NAME} --hall ${HALL_NAME}",
		Aliases: []string{"rack-row"},
		Short:   "delete rack row",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("row name is required")
			}
			if dc == "" {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			if floor == "" {
				cmd.Help()
				return fmt.Errorf("floor name is required")
			}
			if hall == "" {
				cmd.Help()
				return fmt.Errorf("hall name is required")
			}
			return nil
		},
		RunE: deleteRackRow,
	}
	rowCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	rowCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor [REQUIRED]")
	rowCmd.Flags().StringVarP(&hall, "hall", "", "", "name of datacenter hall [REQUIRED]")
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
		Use:     "rack ${RACK_NAME} --dc ${DC_NAME} --floor ${FLOOR_NAME} --hall ${HALL_NAME} --row ${ROW_NAME}",
		Aliases: []string{"rack-row"},
		Short:   "delete rack row",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("rack name is required")
			}
			if dc == "" {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			if floor == "" {
				cmd.Help()
				return fmt.Errorf("floor name is required")
			}
			if hall == "" {
				cmd.Help()
				return fmt.Errorf("hall name is required")
			}
			if row == "" {
				cmd.Help()
				return fmt.Errorf("row name is required")
			}
			return nil
		},
		RunE: deleteRack,
	}
	rackCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	rackCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor [REQUIRED]")
	rackCmd.Flags().StringVarP(&hall, "hall", "", "", "name of datacenter hall [REQUIRED]")
	rackCmd.Flags().StringVarP(&row, "row", "", "", "name of rack row [REQUIRED]")
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
		Use:   "ups ${UPS_NAME} --dc ${DC_NAME}",
		Short: "delete ups",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("ups name is required")
			}
			if dc == "" {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			return nil
		},
		RunE: deleteUPS,
	}
	upsCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	upsCmd.MarkFlagRequired("dc")

	return upsCmd
}

// NewCmdDeletePDU is subcommand represents delete rack row resource.
func NewCmdDeletePDU() *cobra.Command {
	var dc string
	var pduCmd = &cobra.Command{
		Use:   "row-pdu ${PDU_NAME} --dc ${DC_NAME}",
		Short: "delete row-pdu",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("row-pdu name is required")
			}
			if dc == "" {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			return nil
		},
		RunE: deletePDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	pduCmd.MarkFlagRequired("dc")

	return pduCmd
}

// NewCmdDeleteRackPDU is subcommand represents delete rack row resource.
func NewCmdDeleteRackPDU() *cobra.Command {
	var dc string
	var pduCmd = &cobra.Command{
		Use:   "rack-pdu ${RACK_PDU_NAME} --dc ${DC_NAME}",
		Short: "delete rack-pdu",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("rack-pdu name is required")
			}
			return nil
		},
		RunE: deleteRackPDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	pduCmd.MarkFlagRequired("dc")

	return pduCmd
}
