package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/model"
)

func showHost(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/host/name/" + args[0]
	resJSON, err := sendRequest("GET", url, []byte{})
	if err != nil {
		fmt.Println(err)
		return
	}
	host := new(model.Host)
	if err := json.Unmarshal(resJSON, host); err != nil {
		fmt.Println("unmarshal host error:", err)
		return
	}
	rack, err := getRack(host.RackID)
	if err != nil {
		fmt.Println("rack not found")
	}
	row, err := getRow(rack.RowID)
	if err != nil {
		fmt.Println("row not found")
	}
	hall, err := getHall(row.HallID)
	if err != nil {
		fmt.Println("hall not found")
	}
	floor, err := getFloor(hall.FloorID)
	if err != nil {
		fmt.Println("floor not found")
	}
	dc, err := getDataCenter(floor.DataCenterID)
	if err != nil {
		fmt.Println("datacenter not found")
	}
	floor.DataCenter = *dc
	hall.Floor = *floor
	row.Hall = *hall
	rack.RackRow = *row
	host.Rack = *rack

	host.Write(cmd.Flag("output").Value.String())
}

func createHost(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/host"
	// get options
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()
	rackName := cmd.Flag("rack").Value.String()
	//rackID := cmd.Flag("rack-id").Value.String()
	description := cmd.Flag("description").Value.String()
	name := args[0]
	racks, err := getRacks(dcName, floorName, hallName, rowName, rackName)
	if err != nil {
		return fmt.Errorf("rack not found for specified location")
	}
	rack := new(model.Rack)
	for _, r := range *racks {
		rack = &r
		break
	}
	reqModel := model.Host{Name: name, Description: description, RackID: rack.ID}
	reqJSON, _ := json.Marshal(reqModel)
	body, err := sendRequest("POST", url, reqJSON)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateHost(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/host/name/" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	host := new(model.Host)
	if err := json.Unmarshal(body, host); err != nil {
		fmt.Println("json unmarshal error:", err)
		return err
	}
	id := strconv.Itoa(int(host.ID))
	name := cmd.Flag("name").Value.String()
	location := cmd.Flag("location").Value.String()
	description := cmd.Flag("description").Value.String()
	if name == "-" && description == "-" && location == "-" {
		fmt.Println("nothing to be updated")
		return fmt.Errorf("nothing to be updated")
	}
	if name != "-" {
		host.Name = name
		// ensure the new name is not already exists in database
		url = Conf.APIServer.URL + "/host/name/" + name
		body, err = sendRequest("GET", url, []byte{})
		if err == nil {
			fmt.Printf("host '%v' is already exist\n", name)
			return fmt.Errorf("host '%v' is already exist", name)
		}
	}
	if description != "-" {
		host.Description = description
	}
	reqJSON, _ := json.Marshal(host)
	url = Conf.APIServer.URL + "/host/" + id
	body, err = sendRequest("PUT", url, reqJSON)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteHost(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/host/name/" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	host := new(model.Host)
	if err := json.Unmarshal(body, host); err != nil {
		fmt.Println("json unmarshal error:", err)
		return err
	}
	url = Conf.APIServer.URL + "/host/" + strconv.Itoa(int(host.ID))
	body, err = sendRequest("DELETE", url, []byte{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func getHostByName(name string) (*model.Host, error) {
	host := new(model.Host)
	body, err := sendRequest("GET", Conf.APIServer.URL+"/host/name/"+name, []byte{})
	if err != nil {
		return host, err
	}
	if err := json.Unmarshal(body, host); err != nil {
		return host, fmt.Errorf("response parse error")
	}

	return host, nil
}
