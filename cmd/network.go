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

	nw.Write()
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

	nw.Write()
}

func getNetworks(cmd *cobra.Command, url string, query string) {
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
	if cmd.Flag("tree").Value.String() == "true" {
		printNetworkTree(data, 0)
		return
	}

	data.Write()
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
		return fmt.Errorf(resMsg.Message)
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateNetwork(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/network"
	url = url + "/cidr/" + strings.Replace(args[0], "/", "-", 1)
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return err
		}
		fmt.Println(errBody.Error.Message)
		return err
	}
	nw := new(model.IPv4Network)
	if err := json.Unmarshal(body, nw); err != nil {
		fmt.Println("json unmarshal error:", err)
		return err
	}

	url = Conf.APIServer.URL + "/network/" + strconv.FormatUint(uint64(nw.ID), 10)
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
	url := Conf.APIServer.URL + "/network"
	url = url + "/cidr/" + strings.Replace(args[0], "/", "-", 1)
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return err
		}
		fmt.Println(errBody.Error.Message)
		return err
	}
	nw := new(model.IPv4Network)
	if err := json.Unmarshal(body, nw); err != nil {
		fmt.Println("json unmarshal error:", err)
		return err
	}

	url = Conf.APIServer.URL + "/network/" + strconv.FormatUint(uint64(nw.ID), 10)
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
	url := Conf.APIServer.URL + "/vlan"
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	vlanID, _ := strconv.ParseUint(cmd.Flag("network-id").Value.String(), 10, 32)
	reqModel := model.Vlan{Description: description, IPv4NetworkID: uint(vlanID)}
	reqModel.ID = uint(id)
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
		return fmt.Errorf(resMsg.Message)
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
