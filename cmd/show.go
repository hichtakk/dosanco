package cmd

import (
	"strconv"
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdShow() *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show dosanco resources",
	}
	showCmd.AddCommand(
		NewCmdShowNetwork(),
		NewCmdShowIPAM(),
		NewCmdShowVlan(),
		NewCmdShowDataCenter(),
	)

	return showCmd
}

func NewCmdShowNetwork() *cobra.Command {
	var cidr bool
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net", "nw"},
		Short:   "show network",
		Run: func(cmd *cobra.Command, args []string) {
			url := Conf.APIServer.Url + "/network"
			query := "?"
			if tree == true {
				query = query + "&tree=true"
				if depth > 0 {
					query = query + "&depth=" + strconv.Itoa(depth)
				}
			}
			if rfc == true {
				query = query + "&show-rfc-reserved=true"
			}
			if len(args) > 0 {
				if cidr {
					getNetworkByCIDR(url, args[0])
				} else {
					getNetwork(url, args[0])
				}
			} else {
				getNetworks(url, query)
			}
		},
	}
	networkCmd.Flags().BoolVarP(&tree, "tree", "t", false, "get network tree")
	networkCmd.Flags().BoolVarP(&rfc, "show-rfc-reserved", "", false, "show networks defined and reserved in RFC")
	networkCmd.Flags().IntVarP(&depth, "depth", "d", 0, "depth for network tree. this option only work with --tree,-t option")
	networkCmd.Flags().BoolVarP(&cidr, "cidr", "", false, "get network by cidr")

	return networkCmd
}

func NewCmdShowIPAM() *cobra.Command {
	var ipamCmd = &cobra.Command{
		Use:   "ipam",
		Short: "show ip allocation",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires network id or hostname. In case of hostname, option '--host' is required.")
			}
			return nil
		},
		Run: showIPAllocation,
	}
	ipamCmd.Flags().BoolVarP(&hostFlag, "host", "", false, "use host name to get ip allocation")

	return ipamCmd
}

func NewCmdShowVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vlan"},
		Short:   "show vlan",
		Run: func(cmd *cobra.Command, args []string) {
			url := Conf.APIServer.Url + "/vlan"
			getVlans(url)
		},
	}

	return vlanCmd
}

func NewCmdShowDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "show datacenter",
		Run:     getDataCenter,
	}

	return dcCmd
}