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

	fmt.Printf("# Network Data\n")
	fmt.Printf(" ID:             %-3d\n", nw.ID)
	fmt.Printf(" CIDR:           %v\n", nw.CIDR)
	fmt.Printf(" Description:    %v\n", nw.Description)
	fmt.Printf(" SupernetworkID: %d\n\n", nw.SupernetworkID)
	if len(nw.Subnetworks) > 0 {
		fmt.Println("# Subnetworks")
		for _, s := range nw.Subnetworks {
			fmt.Printf(" %-15v %v\n", s.CIDR, s.Description)
		}
	}
	if len(nw.Allocations) > 0 {
		fmt.Println("# IP Allocations")
		for _, a := range nw.Allocations {
			fmt.Printf(" %-15v %v, %v\n", a.Address, a.Name, a.Description)
		}
	}
}

func getNetworkByCIDR(url string, cidr string) {
	url = url + "/cidr/" + strings.Replace(cidr, "/", "-", 1)
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

	fmt.Printf("# Network Data\n")
	fmt.Printf(" ID:             %-3d\n", nw.ID)
	fmt.Printf(" CIDR:           %v\n", nw.CIDR)
	fmt.Printf(" Description:    %v\n", nw.Description)
	fmt.Printf(" SupernetworkID: %d\n\n", nw.SupernetworkID)
	if len(nw.Subnetworks) > 0 {
		fmt.Println("# Subnetworks")
		for _, s := range nw.Subnetworks {
			fmt.Printf(" %-15v %v\n", s.CIDR, s.Description)
		}
	}
	if len(nw.Allocations) > 0 {
		fmt.Printf("# IP Allocations\n")
		for _, a := range nw.Allocations {
			fmt.Printf(" %-15v %v, %v\n", a.Address, a.Name, a.Description)
		}
	}
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
	data := new(model.IPv4Networks)
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

func printNetworkTree(networks *model.IPv4Networks, depth int) {
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
