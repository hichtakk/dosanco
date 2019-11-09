package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/model"
)

func showHost(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		host := new(model.Host)
		hosts, _ := getHosts(map[string]string{"name": args[0]})
		if len(*hosts) == 0 {
			fmt.Printf("host '%v' not found. It might be registered only ip allocation.\n", args[0])
			return
		}
		if len(*hosts) > 1 {
			fmt.Println("multiple hosts found")
			return
		}
		for _, h := range *hosts {
			host = &h
			break
		}
		if host.RackID != 0 {
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
			floor.DataCenter = dc
			hall.Floor = floor
			row.Hall = hall
			rack.RackRow = row
			host.Rack = *rack
		} else {
			host.Rack = model.Rack{}
		}
		if host.GroupID != 0 {
			group, err := getHostGroup(host.GroupID)
			if err != nil {
				fmt.Println("group not found")
			}
			host.Group = group
		}
		networks := map[uint]model.IPv4Network{}
		allocs, _ := getIPv4Allocations(map[string]string{"name": host.Name})
		resAllocs := new(model.IPv4Allocations)
		if len(*allocs) > 0 {
			for _, alloc := range *allocs {
				if network, ok := networks[alloc.IPv4NetworkID]; ok {
					alloc.IPv4Network = &network
				} else {
					nw, _ := getNetwork(alloc.IPv4NetworkID)
					networks[alloc.IPv4NetworkID] = *nw
					alloc.IPv4Network = nw
				}
				*resAllocs = append(*resAllocs, alloc)
			}
			host.IPv4Allocations = *resAllocs
		}

		host.Write(cmd.Flag("output").Value.String())
	} else {
		location := cmd.Flag("location").Value.String()
		group := cmd.Flag("group").Value.String()
		query := map[string]string{}
		if location != "" {
			query["location"] = url.QueryEscape(location)
		}
		if group != "" {
			query["group"] = group
		}
		if len(query) == 0 {
			fmt.Println("flag 'group' or 'location' is required")
			return
		}
		hosts, err := getHosts(query)
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
	location := cmd.Flag("location").Value.String()
	groupName := cmd.Flag("group").Value.String()
	description := cmd.Flag("description").Value.String()
	name := args[0]

	rack := new(model.Rack)
	if location != "" {
		locSlice := strings.Split(location, "/")
		if len(locSlice) != 5 {
			return fmt.Errorf("invalid location format. use '{DC}/{FLOOR}/{HALL}/{ROW}/{RACK}'")
		}
		dcName := locSlice[0]
		floorName := locSlice[1]
		hallName := locSlice[2]
		rowName := locSlice[3]
		rackName := locSlice[4]
		racks, err := getRacks(map[string]string{"dc": dcName, "floor": floorName, "hall": hallName, "row": rowName, "name": rackName})
		if err != nil {
			return fmt.Errorf("rack not found for specified location")
		}
		for _, r := range *racks {
			rack = &r
			break
		}
	}
	groups, err := getHostGroups(map[string]string{"name": groupName})
	if err != nil {
		return err
	}
	var group *model.HostGroup
	for _, g := range *groups {
		group = &g
		break
	}
	reqModel := model.Host{Name: name, Description: description, RackID: rack.ID, GroupID: group.ID}
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
	name := cmd.Flag("name").Value.String()
	location := cmd.Flag("location").Value.String()
	description := cmd.Flag("description").Value.String()
	group := cmd.Flag("group").Value.String()
	if name == "-" && description == "-" && location == "-" && group == "-" {
		return fmt.Errorf("nothing to be updated")
	}

	host := new(model.Host)
	hosts, _ := getHosts(map[string]string{"name": args[0]})
	if len(*hosts) == 0 {
		return fmt.Errorf("host not found")
	} else if len(*hosts) > 1 {
		return fmt.Errorf("multiple hosts are found")
	}
	for _, h := range *hosts {
		host = &h
	}
	if name != "-" {
		host.Name = name
		// ensure the new name is not already exists in database
		exist, _ := getHosts(map[string]string{"name": name})
		if len(*exist) != 0 {
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
	if location == "" {
		host.RackID = 0
	} else if location != "-" {
		locSlice := strings.Split(location, "/")
		if len(locSlice) != 5 {
			return fmt.Errorf("invalid location format. use '{DC}/{FLOOR}/{HALL}/{ROW}/{RACK}'")
		}
		dcName := locSlice[0]
		floorName := locSlice[1]
		hallName := locSlice[2]
		rowName := locSlice[3]
		rackName := locSlice[4]
		racks, err := getRacks(map[string]string{"dc": dcName, "floor": floorName, "hall": hallName, "row": rowName, "name": rackName})
		if err != nil {
			return err
		}
		if len(*racks) > 1 {
			return fmt.Errorf("multiple rack found")
		}
		rack := new(model.Rack)
		for _, r := range *racks {
			rack = &r
		}
		host.RackID = rack.ID
	}
	reqJSON, _ := json.Marshal(host)
	url := Conf.APIServer.URL + "/host/" + strconv.Itoa(int(host.ID))
	body, err := sendRequest("PUT", url, reqJSON)
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
	hosts, err := getHosts(map[string]string{"name": args[0]})
	if err != nil {
		return err
	}
	if len(*hosts) == 0 {
		return fmt.Errorf("host '%v' not found", args[0])
	}
	host := new(model.Host)
	for _, h := range *hosts {
		host = &h
		break
	}
	url := Conf.APIServer.URL + "/host/" + strconv.Itoa(int(host.ID))
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
