package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/handler"
	"github.com/hichikaw/dosanco/model"
)

// Flags
var (
	//tree        bool
	//depth       int
	//rfc         bool
	//supernetID  int
	//description string
	networkID int
	address   string
	hostFlag  bool
)

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

func NewCmdCreateIPAllocation() *cobra.Command {
	var ipamCmd = &cobra.Command{
		Use:   "ipam",
		Short: "create new ip allocation",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires hostname")
			}
			return nil
		},
		RunE: createIPAllocation,
	}
	ipamCmd.Flags().IntVarP(&networkID, "network-id", "n", 0, "network id of the requested ip allocation")
	ipamCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested ip allocation")
	ipamCmd.Flags().StringVarP(&address, "address", "a", "", "ip address of the requested allocation")
	ipamCmd.MarkFlagRequired("network-id")
	ipamCmd.MarkFlagRequired("address")

	return ipamCmd
}

func NewCmdUpdateIPAllocation() *cobra.Command {
	var ipamCmd = &cobra.Command{
		Use:   "ipam",
		Short: "update ip allocation data",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires allocation ID")
			}
			return nil
		},
		RunE: updateIPAllocation,
	}
	ipamCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested ip allocation")
	//ipamCmd.Flags().StringVarP(&hostname, "hostname", "name", "", "ip address of the requested allocation")

	return ipamCmd
}

func NewCmdDeleteIPAllocation() *cobra.Command {
	var ipamCmd = &cobra.Command{
		Use:   "ipam",
		Short: "delete new ip allocation",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires allocation ID")
			}
			return nil
		},
		RunE: deleteIPAllocation,
	}

	return ipamCmd
}

func showIPAllocation(cmd *cobra.Command, args []string) {
	// get network
	nBody, err := sendRequest("GET", Conf.APIServer.Url+"/network", []byte{})
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(nBody, errBody); err != nil {
			fmt.Println("response parse error")
			return
		}
		fmt.Println(errBody.Error.Message)
		return
	}
	data := new([]model.IPv4Network)
	if err := json.Unmarshal(nBody, data); err != nil {
		fmt.Println("json unmarshal err:", err)
		return
	}

	url := Conf.APIServer.Url + "/ipam"
	//_, err := strconv.Atoi(args[0])
	//if err != nil {
	//}
	if hostFlag == true {
		url = url + "/host/" + args[0]
	} else {
		url = url + "/network/" + args[0]
	}
	resp, err := http.Get(url)
	if err != nil {
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	allocs := new([]model.IPv4Allocation)
	if err := json.Unmarshal(body, allocs); err != nil {
	}
	fmt.Printf("%-5s %-15v  %-15v  %-16v  %-v\n", "ID", "Address", "Network", "Name", "Description")
	for _, alloc := range *allocs {
		cidr := getNetworkCIDRfromID(data, alloc.IPv4NetworkID)
		fmt.Printf("%-5d %-15v  %-15v  %-16v  %-v\n", alloc.ID, alloc.Address, cidr, alloc.Name, alloc.Description)
	}
}

func createIPAllocation(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/ipam"
	reqModel := model.IPv4Allocation{Name: args[0], IPv4NetworkID: uint(networkID), Address: address, Description: description}
	reqJson, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(reqJson),
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

	body, err := ioutil.ReadAll(resp.Body)
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

func updateIPAllocation(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/ipam/" + args[0]
	aid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid id. integer is required for argument")
	}
	reqModel := model.IPv4Allocation{Description: description}
	reqModel.ID = uint(aid)
	reqJson, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	req, err := http.NewRequest(
		"PUT",
		url,
		bytes.NewBuffer(reqJson),
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

	body, err := ioutil.ReadAll(resp.Body)
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
	url := Conf.APIServer.Url + "/ipam/" + args[0]
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

	body, err := ioutil.ReadAll(resp.Body)
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
