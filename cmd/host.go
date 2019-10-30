package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/model"
)

func showHost(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
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

		if host.GroupID != 0 {
			group, err := getHostGroup(host.GroupID)
			if err != nil {
				fmt.Println("group not found")
			}
			host.Group = group
		}

		host.Write(cmd.Flag("output").Value.String())
	} else {
		hosts, err := getHosts(map[string]string{"group": cmd.Flag("group").Value.String()})
		if err != nil {
			fmt.Println(err)
			return
		}
		hosts.Write(cmd.Flag("output").Value.String())
	}
}

func showHostGroup(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/host/group"
	if len(args) > 0 {

	} else {
		// show all groups
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.HostGroups)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("json unmarshall error:", err)
			return
		}
		data.Write(cmd.Flag("output").Value.String())
	}
}

func createHost(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/host"
	// get options
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()
	rackName := cmd.Flag("rack").Value.String()
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

func createHostGroup(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/host/group"
	name := args[0]
	description := cmd.Flag("description").Value.String()
	reqModel := model.HostGroup{Name: name, Description: description}
	reqJSON, _ := json.Marshal(reqModel)
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
	group := cmd.Flag("group").Value.String()
	if name == "-" && description == "-" && location == "-" && group == "-" {
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
	if group != "-" {
		query := map[string]string{"name": group}
		groups, err := getHostGroups(query)
		if err != nil {
			fmt.Println("get group error")
		}
		if len(*groups) > 1 {
			fmt.Println("multiple group found")
		}
		for _, grp := range *groups {
			host.GroupID = grp.ID
		}
	}
	// update location
	if location != "" {
		//query := map[string]string{"name": location}

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

func updateHostGroup(cmd *cobra.Command, args []string) error {
	groups, err := getHostGroups(map[string]string{"name": args[0]})
	if err != nil {
		return err
	}
	group := new(model.HostGroup)
	for _, g := range *groups {
		group = &g
		break
	}
	id := strconv.Itoa(int(group.ID))
	name := cmd.Flag("name").Value.String()
	description := cmd.Flag("description").Value.String()
	if name == "-" && description == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	if name != "-" {
		group.Name = name
		// ensure the new name is not already exists in database
		check, _ := getHostGroups(map[string]string{"name": name})
		if len(*check) != 0 {
			return fmt.Errorf(fmt.Sprintf("'%v' is already exists", name))
		}
	}
	if description != "-" {
		group.Description = description
	}
	reqJSON, _ := json.Marshal(group)
	url := Conf.APIServer.URL + "/host/group/" + id
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

func deleteHostGroup(cmd *cobra.Command, args []string) error {
	group := new(model.HostGroup)
	groups, err := getHostGroups(map[string]string{"name": args[0]})
	for _, grp := range *groups {
		group = &grp
		break
	}
	url := Conf.APIServer.URL + "/host/group/" + strconv.Itoa(int(group.ID))
	body, err := sendRequest("DELETE", url, []byte{})
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

func getHostGroup(id uint) (*model.HostGroup, error) {
	group := new(model.HostGroup)
	idStr := strconv.Itoa(int(id))
	body, err := sendRequest("GET", Conf.APIServer.URL+"/host/group/"+idStr, []byte{})
	if err != nil {
		return group, err
	}
	if err := json.Unmarshal(body, group); err != nil {
		return group, fmt.Errorf("response parse error")
	}

	return group, nil
}

func getHosts(query map[string]string) (*model.Hosts, error) {
	hosts := new(model.Hosts)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/host?"+queryString, []byte{})
	if err != nil {
		return hosts, err
	}
	if err := json.Unmarshal(body, hosts); err != nil {
		return hosts, fmt.Errorf("response parse error")
	}

	return hosts, nil
}

func getHostGroups(query map[string]string) (*model.HostGroups, error) {
	groups := new(model.HostGroups)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/host/group?"+queryString, []byte{})
	if err != nil {
		return groups, err
	}
	if err := json.Unmarshal(body, groups); err != nil {
		return groups, fmt.Errorf("response parse error")
	}
	if len(*groups) == 0 {
		return groups, fmt.Errorf("group not found")
	}

	return groups, nil
}
