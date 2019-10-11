package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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
	rfc, _ := strconv.ParseBool(cmd.Flag("show-rfc-reserved").Value.String())
	if len(args) > 0 {
		url = url + "/cidr/" + strings.Replace(args[0], "/", "-", 1)
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		nw := new(model.IPv4Network)
		if err := json.Unmarshal(body, nw); err != nil {
			fmt.Println("json unmarshal error:", err)
			return
		}
		nw.Write(cmd.Flag("output").Value.String())
	} else {
		query := "?"
		if tree == true {
			query = query + "&tree=true"
			if depth > 0 {
				query = query + "&depth=" + cmd.Flag("depth").Value.String()
			}
		}
		if rfc == true {
			query = query + "&show-rfc-reserved=true"
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
	url := Conf.APIServer.URL + "/network"
	url = url + "/cidr/" + strings.Replace(args[0], "/", "-", 1)
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return err
	}
	nw := new(model.IPv4Network)
	if err := json.Unmarshal(body, nw); err != nil {
		fmt.Println("json unmarshal error:", err)
		return err
	}

	url = Conf.APIServer.URL + "/network/" + strconv.FormatUint(uint64(nw.ID), 10)
	body, err = sendRequest("DELETE", url, []byte{})
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

func showIPAllocation(cmd *cobra.Command, args []string) {
	// get network
	nBody, err := sendRequest("GET", Conf.APIServer.URL+"/network", []byte{})
	if err != nil {
		fmt.Println(err)
		return
	}
	data := new([]model.IPv4Network)
	if err := json.Unmarshal(nBody, data); err != nil {
		fmt.Println("json unmarshal err:", err)
		return
	}

	url := Conf.APIServer.URL + "/ip/v4"
	if cmd.Flag("host").Value.String() == "true" {
		url = url + "/host/" + args[0]
	} else {
		targetNW := model.IPv4Network{}
		for _, n := range *data {
			if n.CIDR == args[0] {
				targetNW = n
			}
		}
		if targetNW.ID == 0 {
			fmt.Printf("network %v not found\n", args[0])
			return
		}
		url = url + "/network/" + strconv.Itoa(int(targetNW.ID))
	}
	resp, err := http.Get(url)
	if err != nil {
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	allocs := new(model.IPv4Allocations)
	if err := json.Unmarshal(body, allocs); err != nil {
	}

	output := model.IPv4Allocations{}
	netMap := make(map[uint]*model.IPv4Network)
	for _, alloc := range *allocs {
		network, ok := netMap[alloc.IPv4NetworkID]
		if ok != true {
			network, _ = getNetwork(alloc.IPv4NetworkID)
			netMap[alloc.IPv4NetworkID] = network
		}
		alloc.IPv4Network = netMap[alloc.IPv4NetworkID]
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
	nBody, err := sendRequest("GET", Conf.APIServer.URL+"/network/cidr/"+strings.Replace(cidr, "/", "-", 1), []byte{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	data := new(model.IPv4Network)
	if err := json.Unmarshal(nBody, data); err != nil {
		return fmt.Errorf("json unmarshal err: %v", err)
	}
	addr := args[0]
	url := Conf.APIServer.URL + "/ip/v4"
	reqModel := model.IPv4Allocation{Name: hostname, IPv4NetworkID: data.ID, Address: addr, Type: allocType, Description: description}
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
	url := Conf.APIServer.URL + "/ip/v4/addr/" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	alloc := new(model.IPv4Allocation)
	if err := json.Unmarshal(body, alloc); err != nil {
		fmt.Println("json unmarshal error:", err)
		return err
	}
	url = Conf.APIServer.URL + "/ip/v4/" + strconv.Itoa(int(alloc.ID))
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
	req, err := http.NewRequest(
		"PUT",
		url,
		bytes.NewBuffer(reqJSON),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Errorf(err.Error())
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return fmt.Errorf(err.Error())
	}
	if resp.StatusCode != 200 {
		return errors.New(resMsg.Message)
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteIPAllocation(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/ip/v4/addr/" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	alloc := new(model.IPv4Allocation)
	if err := json.Unmarshal(body, alloc); err != nil {
		fmt.Println("json unmarshal error:", err)
		return err
	}
	url = Conf.APIServer.URL + "/ip/v4/" + strconv.Itoa(int(alloc.ID))
	req, err := http.NewRequest(
		"DELETE",
		url,
		bytes.NewBuffer([]byte{}),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Errorf(err.Error())
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return fmt.Errorf(err.Error())
	}
	if resp.StatusCode != 200 {
		return errors.New(resMsg.Message)
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
