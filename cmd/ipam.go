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

	"github.com/hichikaw/dosanco/handler"
	"github.com/hichikaw/dosanco/model"
)

// Flags
var (
	networkID int
	address   string
	hostFlag  bool
)

func showIPAllocation(cmd *cobra.Command, args []string) {
	// get network
	nBody, err := sendRequest("GET", Conf.APIServer.URL+"/network", []byte{})
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

	url := Conf.APIServer.URL + "/ipam"
	if hostFlag == true {
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
	hostname := cmd.Flag("name").Value.String()
	cidr := cmd.Flag("network").Value.String()
	// get network
	nBody, err := sendRequest("GET", Conf.APIServer.URL+"/network/cidr/"+strings.Replace(cidr, "/", "-", 1), []byte{})
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(nBody, errBody); err != nil {
			return fmt.Errorf("response parse error")
		}
		return fmt.Errorf(errBody.Error.Message)
	}
	data := new(model.IPv4Network)
	if err := json.Unmarshal(nBody, data); err != nil {
		return fmt.Errorf("json unmarshal err: %v", err)
	}
	// validation
	addr := args[0]
	url := Conf.APIServer.URL + "/ipam"
	reqModel := model.IPv4Allocation{Name: hostname, IPv4NetworkID: data.ID, Address: addr, Description: description}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	req, err := http.NewRequest(
		"POST",
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
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
	url := Conf.APIServer.URL + "/ip/v4/" + args[0]
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
	alloc := new(model.IPv4Allocation)
	if err := json.Unmarshal(body, alloc); err != nil {
		fmt.Println("json unmarshal error:", err)
		return err
	}
	url = Conf.APIServer.URL + "/ipam/" + strconv.Itoa(int(alloc.ID))
	alloc.Description = description
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
	url := Conf.APIServer.URL + "/ip/v4/" + args[0]
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
	alloc := new(model.IPv4Allocation)
	if err := json.Unmarshal(body, alloc); err != nil {
		fmt.Println("json unmarshal error:", err)
		return err
	}
	url = Conf.APIServer.URL + "/ipam/" + strconv.Itoa(int(alloc.ID))
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
