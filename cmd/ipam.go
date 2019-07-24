package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	//"strings"

	"github.com/spf13/cobra"

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
	hostname  string
)

func NewCmdShowIPAM() *cobra.Command {
	var ipamCmd = &cobra.Command{
		Use:   "ipam",
		Short: "show ip allocation",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires one argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			url := Conf.APIServer.Url + "/ipam"
			_, err := strconv.Atoi(args[0])
			if err != nil {
			}
			url = url + "/network/" + args[0]
			resp, err := http.Get(url)
			if err != nil {
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
			}
			allocs := new([]model.IPv4Allocation)
			if err := json.Unmarshal(body, allocs); err != nil {
			}
			fmt.Printf("%-15v  %-16v  %-v\n", "Address", "Name", "Description")
			for _, alloc := range *allocs {
				fmt.Printf("%-15v  %-16v  %-v\n", alloc.Address, alloc.Name, alloc.Description)
			}
		},
	}

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
		fmt.Errorf(err.Error())
	}
	if resp.StatusCode != 200 {
		return errors.New(resMsg.Message)
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateIPAllocation(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/ipam"
	aid, err := strconv.Atoi(args[0])
	if err != nil {
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
		fmt.Errorf(err.Error())
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
		fmt.Errorf(err.Error())
	}
	if resp.StatusCode != 200 {
		return errors.New(resMsg.Message)
	}
	fmt.Println(resMsg.Message)

	return nil
}
