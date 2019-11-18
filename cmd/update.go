package cmd

import (
	"fmt"

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
		NewCmdUpdateHostGroup(),
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
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("cidr style is required 'XXX.XXX.XXX.XXX/XX'")
			}
			return nil
		},
		RunE: updateNetwork,
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
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("ip address is required 'XXX.XXX.XXX.XXX'")
			}
			return nil
		},
		RunE: updateIPAllocation,
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
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("vlan id is required")
			}
			return nil
		},
		RunE: updateVlan,
	}
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")

	return vlanCmd
}

// NewCmdUpdateHost is subcommand represents update host resource.
func NewCmdUpdateHost() *cobra.Command {
	var name string
	var location string
	var group string
	var hostCmd = &cobra.Command{
		Use:     "host [NAME]",
		Aliases: []string{"server"},
		Short:   "update host information",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("hostname is required")
			}
			return nil
		},
		RunE: updateHost,
	}
	hostCmd.Flags().StringVarP(&name, "name", "n", "-", "name of the requested host")
	hostCmd.Flags().StringVarP(&location, "location", "l", "-", "location of the requested host")
	hostCmd.Flags().StringVarP(&group, "group", "", "-", "group of the host")
	hostCmd.Flags().StringVarP(&description, "description", "d", "-", "description of the requested host")

	return hostCmd
}

// NewCmdUpdateHostGroup is subcommand represents update host resource.
func NewCmdUpdateHostGroup() *cobra.Command {
	var name string
	var description string
	var groupCmd = &cobra.Command{
		Use:   "group [NAME]",
		Short: "update host group information",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("group is required")
			}
			return nil
		},
		RunE: updateHostGroup,
	}
	groupCmd.Flags().StringVarP(&name, "name", "n", "-", "name of the requested host group")
	groupCmd.Flags().StringVarP(&description, "description", "d", "-", "description of the requested host group")

	return groupCmd
}

// NewCmdUpdateDataCenter is subcommand represents update datacenter resource.
func NewCmdUpdateDataCenter() *cobra.Command {
	var address string
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "update datacenter address",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			return nil
		},
		RunE: updateDataCenter,
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
		RunE: updateDataCenterFloor,
	}
	flrCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
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
		RunE: updateDataCenterHall,
	}
	hallCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	hallCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor [REQUIRED]")
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
		RunE: updateRackRow,
	}
	rowCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	rowCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor [REQUIRED]")
	rowCmd.Flags().StringVarP(&hall, "hall", "", "", "name of datacenter hall [REQUIRED]")
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
		RunE: updateRack,
	}
	rackCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	rackCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor [REQUIRED]")
	rackCmd.Flags().StringVarP(&hall, "hall", "", "", "name of datacenter hall [REQUIRED]")
	rackCmd.Flags().StringVarP(&row, "row", "", "", "name of rack row [REQUIRED]")
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
	var description string
	var upsCmd = &cobra.Command{
		Use:   "ups [UPS_NAME]",
		Short: "update ups information",
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
		RunE: updateUPS,
	}
	upsCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	upsCmd.Flags().StringVarP(&name, "name", "n", "-", "new name of ups")
	upsCmd.Flags().StringVarP(&description, "description", "d", "-", "new description for the ups")
	upsCmd.MarkFlagRequired("dc")

	return upsCmd
}

// NewCmdUpdatePDU is subcommand represents update rack resource.
func NewCmdUpdatePDU() *cobra.Command {
	var dc string
	var name string
	var description string
	var pduCmd = &cobra.Command{
		Use:   "row-pdu [PDU_NAME]",
		Short: "update row-pdu name",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("row-pdu name is required")
			}
			if dc == "" {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			/*
				if primary == "" {
					cmd.Help()
					return fmt.Errorf("primary ups name is required")
				}
			*/
			return nil
		},
		RunE: updatePDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	pduCmd.Flags().StringVarP(&name, "name", "n", "-", "new name of row-pdu")
	pduCmd.Flags().StringVarP(&description, "description", "d", "-", "new description of the requested row-pdu")
	pduCmd.MarkFlagRequired("dc")

	return pduCmd
}

// NewCmdUpdateRackPDU is subcommand represents update rack resource.
func NewCmdUpdateRackPDU() *cobra.Command {
	var dc string
	var name string
	var primary string
	var secondary string
	var pduCmd = &cobra.Command{
		Use:   "rack-pdu [RACK_PDU_NAME]",
		Short: "update rack-pdu name",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("target rack-pdu name is required")
			}
			/*
				if location == "" {
					cmd.Help()
					return fmt.Errorf("location path is required")
				}
				if group == "" {
					cmd.Help()
					return fmt.Errorf("group name is required")
				}
			*/
			/*
				if primary == "" {
					cmd.Help()
					return fmt.Errorf("primary row-pdu name is required")
				}
			*/
			return nil
		},
		RunE: updateRackPDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	pduCmd.Flags().StringVarP(&name, "name", "n", "-", "new name of rack")
	pduCmd.Flags().StringVarP(&primary, "primary", "", "-", "new primary row-pdu name")
	pduCmd.Flags().StringVarP(&secondary, "secondary", "", "-", "new secondary row-pdu name")
	pduCmd.MarkFlagRequired("dc")

	return pduCmd
}
