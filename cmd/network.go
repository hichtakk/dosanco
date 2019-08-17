package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/handler"
	"github.com/hichikaw/dosanco/model"
)

// Flags
var (
	tree        bool
	depth       int
	rfc         bool
	supernetID  int
	description string
)

func NewCmdShowNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net", "nw"},
		Short:   "show network",
		Run: func(cmd *cobra.Command, args []string) {
			url := Conf.APIServer.Url + "/network"
			query := ""
			if tree == true {
				query = query + "?tree=true"
				if depth > 0 {
					query = query + "&depth=" + strconv.Itoa(depth)
				}
			}
			if len(args) > 0 {
				getNetwork(url, args[0])
			} else {
				getNetworks(url, query)
			}
		},
	}
	networkCmd.Flags().BoolVarP(&tree, "tree", "t", false, "get network tree")
	networkCmd.Flags().BoolVarP(&rfc, "show-rfc-defined", "", false, "show networks defined and reserved in RFC")
	networkCmd.Flags().IntVarP(&depth, "depth", "d", 0, "depth for network tree. this option only work with --tree,-t option")

	return networkCmd
}

func NewCmdCreateNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net", "nw"},
		Short:   "create new network",
		Long:    "create new network",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires one network cidr")
			}
			return nil
		},
		RunE: createNetwork,
	}
	networkCmd.Flags().IntVarP(&supernetID, "supernet-id", "s", 0, "supernetwork id of the requested network")
	networkCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested network")
	networkCmd.MarkFlagRequired("supernet-id")

	return networkCmd
}

func NewCmdUpdateNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net", "nw"},
		Short:   "update network description",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires one network id")
			}
			return nil
		},
		RunE: updateNetwork,
	}
	networkCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested network")

	return networkCmd
}

func NewCmdDeleteNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net", "nw"},
		Short:   "delete network description",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires one network id")
			}
			return nil
		},
		RunE: deleteNetwork,
	}

	return networkCmd
}

func getNetwork(url string, id string) {
	url = url + "/" + id
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return
		}
		fmt.Println(errBody.Error.Message)
		return
	}
	nw := new(model.IPv4Network)
	if err := json.Unmarshal(body, nw); err != nil {
		fmt.Println("json unmarshal error:", err)
		return
	}

	fmt.Printf("%2v %-20s	%s\n", nw.ID, nw.CIDR, nw.Description)
}

func getNetworks(url string, query string) {
	if query != "" {
		url = url + query
	}
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return
		}
		fmt.Println(errBody.Error.Message)
		return
	}
	data := new([]model.IPv4Network)
	if err := json.Unmarshal(body, data); err != nil {
		fmt.Println("json unmarshal err:", err)
		return
	}

	// write output
	if output == "json" {
		fmt.Println(string(body))
		return
	}
	if tree == true {
		printNetworkTree(data, 0)
		return
	}
	fmt.Printf("ID	CIDR			Description\n")
	for _, network := range *data {
		fmt.Printf("%2d	%-20s	%s\n", network.ID, network.CIDR, network.Description)
	}
}

func printNetworkTree(networks *[]model.IPv4Network, depth int) {
	for _, network := range *networks {
		fmt.Printf("%v%v:%v\n", strings.Repeat("   ", depth), network.ID, network.CIDR)
		if len(network.Subnetworks) > 0 {
			printNetworkTree(&network.Subnetworks, depth+1)
		}
	}
}

func createNetwork(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/network"
	reqModel := model.IPv4Network{CIDR: args[0], SupernetworkID: uint(supernetID), Description: description}
	reqJson, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url, reqJson)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	} else {
		fmt.Println(resMsg.Message)
	}

	return nil
}

func updateNetwork(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/network"
	url = url + "/" + args[0]
	reqModel := model.IPv4Network{Description: description}
	reqJson, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("PUT", url, reqJson)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	} else {
		fmt.Println(resMsg.Message)
	}

	return nil
}

func deleteNetwork(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/network"
	url = url + "/" + args[0]
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	} else {
		fmt.Println(resMsg.Message)
	}

	return nil
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

func getVlans(url string) {
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return
		}
		fmt.Println(errBody.Error.Message)
		return
	}
	data := new([]model.Vlan)
	if err := json.Unmarshal(body, data); err != nil {
		fmt.Println("json unmarshal err:", err)
		return
	}

	// write output
	if output == "json" {
		fmt.Println(string(body))
		return
	}
	fmt.Printf("%-4s %-20s %-20s\n", "ID", "NetworkID", "Description")
	for _, vlan := range *data {
		fmt.Printf("%-4d %-20d %-20s\n", vlan.ID, vlan.IPv4NetworkID, vlan.Description)
	}
}

func NewCmdCreateVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vlan"},
		Short:   "create new vlan",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires vlan id")
			}
			return nil
		},
		RunE: createVlan,
	}
	vlanCmd.Flags().IntVarP(&networkID, "network-id", "n", 0, "network id of the requested ip allocation")
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")
	vlanCmd.MarkFlagRequired("network-id")

	return vlanCmd
}

func createVlan(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/vlan"
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	reqModel := model.Vlan{Description: description, IPv4NetworkID: uint(networkID)}
	reqModel.ID = uint(id)
	reqJson, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url, reqJson)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	} else {
		fmt.Println(resMsg.Message)
	}

	return nil
}

func NewCmdUpdateVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vlan"},
		Short:   "update vlan description",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires vlan id")
			}
			return nil
		},
		RunE: updateVlan,
	}
	vlanCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested vlan")

	return vlanCmd
}

func updateVlan(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/vlan"
	url = url + "/" + args[0]
	reqModel := model.Vlan{Description: description}
	reqJson, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("PUT", url, reqJson)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	} else {
		fmt.Println(resMsg.Message)
	}

	return nil
}

func NewCmdDeleteVlan() *cobra.Command {
	var vlanCmd = &cobra.Command{
		Use:     "vlan",
		Aliases: []string{"vlan"},
		Short:   "delete vlan description",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires vlan id")
			}
			return nil
		},
		RunE: deleteVlan,
	}

	return vlanCmd
}

func deleteVlan(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/vlan"
	url = url + "/" + args[0]
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	} else {
		fmt.Println(resMsg.Message)
	}

	return nil
}
