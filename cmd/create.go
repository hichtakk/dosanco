package cmd

import (
	"fmt"

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
		Use:     "network ${CIDR}",
		Aliases: []string{"net", "nw"},
		Short:   "create new network",
		Long:    "create new network",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("cidr style network is required 'XXX.XXX.XXX.XXX/XX'")
			}
			return nil
		},
		RunE: createNetwork,
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
		Use:     "ip ${IP_ADDRESS} --name ${NAME} --network ${CIDR}",
		Aliases: []string{"ip-alloc"},
		Short:   "create new ip allocation",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("ip address is required 'XXX.XXX.XXX.XXX'")
			}
			return nil
		},
		RunE: createIPAllocation,
	}
	ipCmd.Flags().StringVarP(&description, "description", "d", "", "description for the requested ip allocation")
	ipCmd.Flags().StringVarP(&name, "name", "", "", "hostname for the ip address. [REQUIRED] in case allocation type is generic.")
	ipCmd.Flags().StringVarP(&network, "network", "", "", "network CIDR for the ip address. [REQUIRED]")
	ipCmd.Flags().StringVarP(&allocType, "type", "t", "generic", "type of address. use 'reserved' or 'generic'")
	ipCmd.MarkFlagRequired("network")

	return ipCmd
}

// NewCmdCreateVlan is subcommand represents vlan resource.
func NewCmdCreateVlan() *cobra.Command {
	var network string
	var vlanCmd = &cobra.Command{
		Use:   "vlan ${VLAN_ID} --network XXX.XXX.XXX.XXX/XX",
		Short: "create new vlan",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("vlan id is required")
			}
			if network == "" {
				cmd.Help()
				return fmt.Errorf("network cidr for the vlan is required")
			}
			return nil
		},
		RunE: createVlan,
	}
	vlanCmd.Flags().StringVarP(&network, "network", "", "", "network cidr of the requested vlan [REQUIRED]")
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")
	vlanCmd.MarkFlagRequired("cidr")

	return vlanCmd
}

// NewCmdCreateHost is subcommand represents vlan resource.
func NewCmdCreateHost() *cobra.Command {
	var location string
	var group string
	var hostCmd = &cobra.Command{
		Use:     "host ${NAME} --group ${GROUP}",
		Aliases: []string{"server", "node"},
		Short:   "create new host",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("hostname is required")
			}
			return nil
		},
		RunE: createHost,
	}
	hostCmd.Flags().StringVarP(&location, "location", "l", "", "location of host installed. use format '{DC}/{FLOOR}/{HALL}/{ROW}/{RACK}'")
	hostCmd.Flags().StringVarP(&description, "description", "d", "", "description of the host")
	hostCmd.Flags().StringVarP(&group, "group", "g", "", "group of the host [REQUIRED]")
	hostCmd.MarkFlagRequired("group")

	return hostCmd
}

// NewCmdCreateHostGroup is subcommand represents vlan resource.
func NewCmdCreateHostGroup() *cobra.Command {
	var groupCmd = &cobra.Command{
		Use:   "group ${NAME}",
		Short: "create new host group",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("group is required")
			}
			return nil
		},
		RunE: createHostGroup,
	}
	groupCmd.Flags().StringVarP(&description, "description", "d", "", "name of host group")

	return groupCmd
}

// NewCmdCreateDataCenter is subcommand represents datacenter resource.
func NewCmdCreateDataCenter() *cobra.Command {
	var address string
	var dcCmd = &cobra.Command{
		Use:     "datacenter ${NAME}",
		Aliases: []string{"dc"},
		Short:   "create new datacenter",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			return nil
		},
		RunE: createDataCenter,
	}
	dcCmd.Flags().StringVarP(&address, "address", "a", "", "address of data center")
	dcCmd.MarkFlagRequired("address")

	return dcCmd
}

// NewCmdCreateDataCenterFloor is subcommand represents datacenter resource.
func NewCmdCreateDataCenterFloor() *cobra.Command {
	var dc string
	var flrCmd = &cobra.Command{
		Use:     "floor ${NAME} --dc ${DC_NAME}",
		Aliases: []string{"dc-floor", "area"},
		Short:   "create new floor to datacenter",
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
		RunE: createDataCenterFloor,
	}
	flrCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	flrCmd.MarkFlagRequired("dc")

	return flrCmd
}

// NewCmdCreateDataCenterHall is subcommand represents datacenter resource.
func NewCmdCreateDataCenterHall() *cobra.Command {
	var dc string
	var floor string
	var hallCmd = &cobra.Command{
		Use:     "hall ${NAME} --dc ${DC_NAME} --floor ${FLOOR_NAME}",
		Aliases: []string{"dc-hall"},
		Short:   "create new hall to datacenter floor",
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
		RunE: createDataCenterHall,
	}
	hallCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	hallCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor [REQUIRED]")
	hallCmd.MarkFlagRequired("dc")
	hallCmd.MarkFlagRequired("floor")

	return hallCmd
}

// NewCmdCreateRackRow is subcommand represents datacenter resource.
func NewCmdCreateRackRow() *cobra.Command {
	var dc string
	var floor string
	var hall string
	var rowCmd = &cobra.Command{
		Use:     "row ${NAME} --dc ${DC_NAME} --floor ${FLOOR_NAME} --hall ${HALL_NAME}",
		Aliases: []string{"rack-row"},
		Short:   "create new hall to datacenter floor",
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
		RunE: createRackRow,
	}
	rowCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	rowCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor [REQUIRED]")
	rowCmd.Flags().StringVarP(&hall, "hall", "", "", "name of data hall [REQUIRED]")
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
		Use:     "rack ${RACK_NAME} --dc ${DC_NAME} --floor ${FLOOR_NAME} --hall ${HALL_NAME} --row ${ROW_NAME}",
		Aliases: []string{""},
		Short:   "create new rack to row",
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
		RunE: createRack,
	}
	rackCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	rackCmd.Flags().StringVarP(&floor, "floor", "", "", "name of datacenter floor [REQUIRED]")
	rackCmd.Flags().StringVarP(&hall, "hall", "", "", "name of data hall [REQUIRED]")
	rackCmd.Flags().StringVarP(&row, "row", "", "", "name of row [REQUIRED]")
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
		Use:   "ups ${UPS_NAME} --dc ${DC_NAME}",
		Short: "create new ups",
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
		RunE: createUPS,
	}
	upsCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	upsCmd.MarkFlagRequired("dc")

	return upsCmd
}

// NewCmdCreatePDU is subcommand represents pdu resource.
func NewCmdCreatePDU() *cobra.Command {
	var dc string
	var primary string
	var secondary string
	var pduCmd = &cobra.Command{
		Use:   "row-pdu ${ROW_PDU_NAME} --dc ${DC_NAME} --primary ${UPS_NAME} [--secondary ${UPS_NAME}]",
		Short: "create new row-pdu",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("row-pdu name is required")
			}
			if dc == "" {
				cmd.Help()
				return fmt.Errorf("datacenter name is required")
			}
			if primary == "" {
				cmd.Help()
				return fmt.Errorf("primary ups name is required")
			}
			return nil
		},
		RunE: createPDU,
	}
	pduCmd.Flags().StringVarP(&dc, "dc", "", "", "name of datacenter [REQUIRED]")
	pduCmd.Flags().StringVarP(&primary, "primary", "", "", "name of primary power source [REQUIRED]")
	pduCmd.Flags().StringVarP(&secondary, "secondary", "", "", "name of secondary power source")
	pduCmd.MarkFlagRequired("dc")
	pduCmd.MarkFlagRequired("primary")

	return pduCmd
}

// NewCmdCreateRackPDU is subcommand represents pdu resource.
func NewCmdCreateRackPDU() *cobra.Command {
	var location string
	var primary string
	var secondary string
	var group string
	var pduCmd = &cobra.Command{
		Use:   "rack-pdu ${RACK_PDU_NAME} --location ${LOCATION_PATH} --group ${GROUP_NAME} --primary ${ROW_PDU_NAME} [--secondary ${ROW_PDU_NAME}]",
		Short: "create new rack-pdu",
		Long:  `create new rack-pdu dosanco create rack-pdu --dc DC1 --primary ROW-PDU-1 rack-pdu01.dosanco`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("rack-pdu name is required")
			}
			if location == "" {
				cmd.Help()
				return fmt.Errorf("location path is required")
			}
			if group == "" {
				cmd.Help()
				return fmt.Errorf("group name is required")
			}
			if primary == "" {
				cmd.Help()
				return fmt.Errorf("primary row-pdu name is required")
			}
			return nil
		},
		RunE: createRackPDU,
	}
	pduCmd.Flags().StringVarP(&location, "location", "l", "", "location of host installed. use format '{DC}/{FLOOR}/{HALL}/{ROW}/{RACK}'")
	pduCmd.Flags().StringVarP(&primary, "primary", "", "", "name of primary power source [REQUIRED]")
	pduCmd.Flags().StringVarP(&secondary, "secondary", "", "", "name of secondary power source")
	pduCmd.Flags().StringVarP(&group, "group", "", "", "name of host group [REQUIRED]")
	pduCmd.MarkFlagRequired("dc")
	pduCmd.MarkFlagRequired("floor")
	pduCmd.MarkFlagRequired("hall")
	pduCmd.MarkFlagRequired("row")
	pduCmd.MarkFlagRequired("rack")
	pduCmd.MarkFlagRequired("group")
	pduCmd.MarkFlagRequired("primary")

	return pduCmd
}
