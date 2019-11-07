package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/model"
)

// Flags
var (
	description string
)

func showNetwork(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/network"
	tree, _ := strconv.ParseBool(cmd.Flag("tree").Value.String())
	depth, _ := strconv.Atoi(cmd.Flag("depth").Value.String())
	if len(args) > 0 {
		query := map[string]string{"cidr": args[0]}
		networks, err := getNetworks(query)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		network := new(model.IPv4Network)
		for _, nw := range *networks {
			network = &nw
		}
		network.Write(cmd.Flag("output").Value.String())
	} else {
		query := "?"
		if tree == true {
			query = query + "&tree=true"
			if depth > 0 {
				query = query + "&depth=" + cmd.Flag("depth").Value.String()
			}
		}
		url = url + query
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.IPv4Networks)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("json unmarshal err:", err)
			return
		}
		data.Write(cmd.Flag("output").Value.String())
	}
}

func createNetwork(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/network"
	reqModel := model.IPv4Network{CIDR: args[0], Description: description}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateNetwork(cmd *cobra.Command, args []string) error {
	query := map[string]string{"cidr": args[0]}
	networks, err := getNetworks(query)
	if err != nil {
		return err
	}
	network := new(model.IPv4Network)
	for _, nw := range *networks {
		network = &nw
	}
	url := Conf.APIServer.URL + "/network/" + strconv.FormatUint(uint64(network.ID), 10)
	reqModel := model.IPv4Network{Description: description}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteNetwork(cmd *cobra.Command, args []string) error {
	query := map[string]string{"cidr": args[0]}
	networks, err := getNetworks(query)
	if err != nil {
		return err
	}
	network := new(model.IPv4Network)
	for _, nw := range *networks {
		network = &nw
	}
	url := Conf.APIServer.URL + "/network/" + strconv.FormatUint(uint64(network.ID), 10)
	body, err := sendRequest("DELETE", url, []byte{})
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func getVlans(url string) {
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		fmt.Println(err)
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

func createVlan(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/vlan"
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	cidr := cmd.Flag("network").Value.String()
	networks, err := getNetworks(map[string]string{"cidr": cidr})
	if err != nil {
		return err
	}
	network := new(model.IPv4Network)
	for _, n := range *networks {
		network = &n
		break
	}
	reqModel := model.Vlan{Description: description, IPv4NetworkID: uint(network.ID)}
	reqModel.ID = uint(id)
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url, reqJSON)
	if reqErr != nil {
		return reqErr
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateVlan(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/vlan"
	url = url + "/" + args[0]
	reqModel := model.Vlan{Description: description}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteVlan(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/vlan"
	url = url + "/" + args[0]
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	}
	fmt.Println(resMsg.Message)

	return nil
}

func showIPAllocation(cmd *cobra.Command, args []string) {
	// get network
	query := map[string]string{"cidr": args[0]}
	networks, err := getNetworks(query)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	network := new(model.IPv4Network)
	for _, nw := range *networks {
		network = &nw
	}
	allocQuery := map[string]string{}
	allocQuery["cidr"] = network.CIDR

	allocs, _ := getIPv4Allocations(allocQuery)
	output := model.IPv4Allocations{}
	for _, alloc := range *allocs {
		alloc.IPv4Network = network
		output = append(output, alloc)
	}

	output.Write(cmd.Flag("output").Value.String())
}

func createIPAllocation(cmd *cobra.Command, args []string) error {
	hostname := cmd.Flag("name").Value.String()
	cidr := cmd.Flag("network").Value.String()
	allocType := cmd.Flag("type").Value.String()
	if (allocType != "reserved") && (allocType != "generic") {
		return fmt.Errorf("ip allocation type error")
	}
	// get network
	query := map[string]string{"cidr": cidr}
	networks, err := getNetworks(query)
	if err != nil {
		return err
	}
	network := new(model.IPv4Network)
	for _, nw := range *networks {
		network = &nw
	}

	addr := args[0]
	url := Conf.APIServer.URL + "/ip/v4"
	reqModel := model.IPv4Allocation{Name: hostname, IPv4NetworkID: network.ID, Address: addr, Type: allocType, Description: description}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, err := sendRequest("POST", url, reqJSON)
	if err != nil {
		return err
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateIPAllocation(cmd *cobra.Command, args []string) error {
	alloc := new(model.IPv4Allocation)
	allocs, _ := getIPv4Allocations(map[string]string{"address": args[0]})
	for _, a := range *allocs {
		alloc = &a
		break
	}
	name := cmd.Flag("name").Value.String()
	description := cmd.Flag("description").Value.String()
	if name != "-" && name != alloc.Name {
		alloc.Name = name
	}
	if description != "-" && description != alloc.Description {
		alloc.Description = description
	}
	if name == "-" && description == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	reqJSON, err := json.Marshal(alloc)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", alloc)
	}
	url := Conf.APIServer.URL + "/ip/v4/" + strconv.Itoa(int(alloc.ID))
	body, err := sendRequest("PUT", url, reqJSON)
	if err != nil {
		return err
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteIPAllocation(cmd *cobra.Command, args []string) error {
	alloc := new(model.IPv4Allocation)
	allocs, _ := getIPv4Allocations(map[string]string{"address": args[0]})
	for _, a := range *allocs {
		alloc = &a
		break
	}
	url := Conf.APIServer.URL + "/ip/v4/" + strconv.Itoa(int(alloc.ID))
	body, err := sendRequest("DELETE", url, []byte{})
	if err != nil {
		return err
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func getNetworkCIDRfromID(networks *[]model.IPv4Network, id uint) string {
	for _, network := range *networks {
		if network.ID == id {
			return network.CIDR
		}
	}

	return "?"
}

func getNetwork(id uint) (*model.IPv4Network, error) {
	network := new(model.IPv4Network)
	idStr := strconv.Itoa(int(id))
	body, err := sendRequest("GET", Conf.APIServer.URL+"/network/"+idStr, []byte{})
	if err != nil {
		return network, err
	}
	if err := json.Unmarshal(body, network); err != nil {
		return network, fmt.Errorf("response parse error")
	}

	return network, nil
}

func getNetworks(query map[string]string) (*model.IPv4Networks, error) {
	networks := new(model.IPv4Networks)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/network?"+queryString, []byte{})
	if err != nil {
		return networks, err
	}
	if err := json.Unmarshal(body, networks); err != nil {
		return networks, fmt.Errorf("response parse error")
	}
	if len(*networks) == 0 {
		return networks, fmt.Errorf("no network found")
	}

	return networks, nil
}

func getIPv4Allocation(id uint) (*model.IPv4Allocation, error) {
	alloc := new(model.IPv4Allocation)
	idStr := strconv.Itoa(int(id))
	body, err := sendRequest("GET", Conf.APIServer.URL+"/ip/v4/"+idStr, []byte{})
	if err != nil {
		return alloc, err
	}
	if err := json.Unmarshal(body, alloc); err != nil {
		return alloc, fmt.Errorf("response parse error")
	}

	return alloc, nil
}

func getIPv4Allocations(query map[string]string) (*model.IPv4Allocations, error) {
	allocs := new(model.IPv4Allocations)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/ip/v4?"+queryString, []byte{})
	if err != nil {
		return allocs, err
	}
	if err := json.Unmarshal(body, allocs); err != nil {
		return allocs, fmt.Errorf("response parse error")
	}
	//if len(*allocs) == 0 {
	//	return allocs, fmt.Errorf("no allocation found")
	//}

	return allocs, nil
}
