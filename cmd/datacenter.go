package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/model"
)

func getDataCenter(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/datacenter"
	if len(args) > 0 {
		// show specified datacenter
		body, err := sendRequest("GET", url+"/name/"+args[0], []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		dc := new(model.DataCenter)
		if err = json.Unmarshal(body, dc); err != nil {
			fmt.Println("response parse error", err)
		}
		url = url + "/" + strconv.Itoa(int(dc.ID))
		body, err = sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.DataCenter)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("response parse error", err)
			return
		}
		data.Write(cmd.Flag("output").Value.String())
	} else {
		// show all datacenters
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.DataCenters)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("json unmarshall error:", err)
			return
		}
		data.Write(cmd.Flag("output").Value.String())
	}
}

func getDataCenterFloor(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/datacenter/floor"
	dcName := cmd.Flag("dc").Value.String()
	if len(args) > 0 {
		// show specified datacenter floor
		if dcName == "" {
			fmt.Println("datacenter name is required for showing specific floor")
			return
		}
		body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		dc := new(model.DataCenter)
		if err = json.Unmarshal(body, dc); err != nil {
			fmt.Println("response parse error", err)
		}
		dcID := strconv.Itoa(int(dc.ID))
		body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.Floors)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("response parse error:", err)
			return
		}
		floor := model.Floor{}
		for _, flr := range *data {
			if flr.Name == args[0] {
				floor = flr
			}
		}
		if floor.ID != 0 {
			floor.Write(cmd.Flag("output").Value.String())
		} else {
			fmt.Println("floor not found")
		}
	} else {
		if dcName != "" {
			// show all datacenter floors
			url = Conf.APIServer.URL + "/datacenter/name/" + dcName
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			dc := new(model.DataCenter)
			if err = json.Unmarshal(body, dc); err != nil {
				fmt.Println("response parse error", err)
			}
			dcID := strconv.Itoa(int(dc.ID))
			body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			data := new(model.Floors)
			if err := json.Unmarshal(body, data); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			data.Write(cmd.Flag("output").Value.String())
		} else {
			// show all datacenter floors
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			data := new(model.Floors)
			if err := json.Unmarshal(body, data); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			data.Write(cmd.Flag("output").Value.String())
		}
	}
}

func getDataCenterHall(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/datacenter/hall"
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	if len(args) > 0 {
		// show specified datacenter floor
		if dcName == "" {
			fmt.Println("datacenter name is required for showing specific hall")
			return
		}
		if floorName == "" {
			fmt.Println("floor name is required for showing specific hall")
			return
		}
		body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		dc := new(model.DataCenter)
		if err = json.Unmarshal(body, dc); err != nil {
			fmt.Println("response parse error", err)
		}
		dcID := strconv.Itoa(int(dc.ID))
		body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.Floors)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("response parse error:", err)
			return
		}
		floor := model.Floor{}
		for _, flr := range *data {
			if flr.Name == floorName {
				floor = flr
			}
		}
		if floor.ID == 0 {
			fmt.Println("floor not found")
		}

		fID := strconv.Itoa(int(floor.ID))
		f := new(model.Floor)
		body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/floor/"+fID, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = json.Unmarshal(body, f); err != nil {
			fmt.Println(err)
			return
		}
		for _, h := range f.Halls {
			if h.Name == args[0] {
				h.Write(cmd.Flag("output").Value.String())
			}
		}
	} else {
		if dcName != "" {
			// show all floors of specified datacenter
			url = Conf.APIServer.URL + "/datacenter/name/" + dcName
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			dc := new(model.DataCenter)
			if err = json.Unmarshal(body, dc); err != nil {
				fmt.Println("response parse error", err)
			}
			dcID := strconv.Itoa(int(dc.ID))
			body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			data := new(model.Floors)
			if err := json.Unmarshal(body, data); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			if floorName == "" {
				halls := model.Halls{}
				for _, f := range *data {
					floor := new(model.Floor)
					fID := strconv.Itoa(int(f.ID))
					body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/floor/"+fID, []byte{})
					if err != nil {
						fmt.Println(err)
						return
					}
					if err = json.Unmarshal(body, floor); err != nil {
						fmt.Println(err)
						return
					}
					if len(floor.Halls) > 0 {
						for _, h := range floor.Halls {
							halls = append(halls, h)
						}
					}
				}
				halls.Write(cmd.Flag("output").Value.String())
			} else {
				floor := new(model.Floor)
				for _, f := range *data {
					if f.Name == floorName {
						floor = &f
						break
					}
				}
				fID := strconv.Itoa(int(floor.ID))
				body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/floor/"+fID, []byte{})
				if err != nil {
					fmt.Println(err)
					return
				}
				if err = json.Unmarshal(body, floor); err != nil {
					fmt.Println(err)
					return
				}
				floor.Halls.Write(cmd.Flag("output").Value.String())
			}
		} else {
			// show all datacenter halls
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			data := new(model.Halls)
			if err := json.Unmarshal(body, data); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			data.Write(cmd.Flag("output").Value.String())
		}
	}
}

func getRackRow(cmd *cobra.Command, args []string) {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	if len(args) > 0 {
		// show specified row
		url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&name=" + args[0]
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.RackRows)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		if len(*data) > 1 {
			fmt.Println("multiple row found")
			return
		}
		for _, r := range *data {
			r.Write(cmd.Flag("output").Value.String())
			break
		}
	} else {
		// show list of rows
		url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName
		if floorName != "" {
			url = url + "&floor=" + floorName
		}
		if hallName != "" {
			url = url + "&hall=" + hallName
		}
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.RackRows)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		data.Write(cmd.Flag("output").Value.String())
	}
}

func getRack(cmd *cobra.Command, args []string) {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()
	if len(args) > 0 {
		/*
			// show specified rack
			url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&name=" + args[0]
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			data := new(model.RackRows)
			if err := json.Unmarshal(body, data); err != nil {
				fmt.Println("parse response error")
				return
			}
			if len(*data) > 1 {
				fmt.Println("multiple row found")
				return
			}
			for _, r := range *data {
				r.Write(cmd.Flag("output").Value.String())
				break
			}
		*/
	} else {
		// show list of racks
		url := Conf.APIServer.URL + "/datacenter/rack?dc=" + dcName
		if floorName != "" {
			url = url + "&floor=" + floorName
		}
		if hallName != "" {
			url = url + "&hall=" + hallName
		}
		if rowName != "" {
			url = url + "&row=" + rowName
		}
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.Racks)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		data.Write(cmd.Flag("output").Value.String())
	}
}

func createDataCenter(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	reqModel := model.DataCenter{Name: args[0], Address: cmd.Flag("address").Value.String()}
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

func createDataCenterFloor(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	// get datacenter
	dcName := cmd.Flag("dc").Value.String()
	body, err := sendRequest("GET", url+"/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error")
	}
	// prepare request floor model
	reqModel := model.Floor{Name: args[0], DataCenterID: dc.ID}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url+"/floor", reqJSON)
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

func createDataCenterHall(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	// get datacenter
	dcName := cmd.Flag("dc").Value.String()
	body, err := sendRequest("GET", url+"/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error")
	}
	// get datacenter floor
	dcID := strconv.Itoa(int(dc.ID))
	floorName := cmd.Flag("floor").Value.String()
	body, err = sendRequest("GET", url+"/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	floors := new(model.Floors)
	if err = json.Unmarshal(body, floors); err != nil {
		return fmt.Errorf("response parse error")
	}
	var floorID uint
	for _, floor := range *floors {
		if floor.Name == floorName {
			floorID = floor.ID
			break
		}
	}
	if floorID == 0 {
		return fmt.Errorf("floor not found")
	}
	// hall type
	hallType := cmd.Flag("type").Value.String()
	if (hallType != "generic") && (hallType != "network") {
		return fmt.Errorf("type must be 'network' or 'generic'")
	}

	// prepare request hall model
	reqModel := model.Hall{Name: args[0], FloorID: floorID, Type: hallType}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url+"/hall", reqJSON)
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

func createRackRow(cmd *cobra.Command, args []string) error {
	// get data hall
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	url := Conf.APIServer.URL + "/datacenter/hall?dc=" + dcName + "&floor=" + floorName + "&name=" + hallName
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	halls := new(model.Halls)
	if err = json.Unmarshal(body, halls); err != nil {
		return fmt.Errorf("response parse error")
	}
	if len(*halls) == 0 {
		return fmt.Errorf("hall not found")
	} else if len(*halls) > 1 {
		return fmt.Errorf("multiple hall found")
	}
	hall := model.Hall{}
	for _, h := range *halls {
		hall = h
	}
	// prepare request hall model
	reqModel := model.RackRow{Name: args[0], HallID: hall.ID}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", Conf.APIServer.URL+"/datacenter/row", reqJSON)
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

func createRack(cmd *cobra.Command, args []string) error {
	// get data hall
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()
	url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&name=" + rowName
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	rows := new(model.RackRows)
	if err = json.Unmarshal(body, rows); err != nil {
		return fmt.Errorf("response parse error")
	}
	if len(*rows) == 0 {
		return fmt.Errorf("row not found")
	} else if len(*rows) > 1 {
		return fmt.Errorf("multiple row found")
	}
	row := model.RackRow{}
	for _, r := range *rows {
		row = r
		break
	}
	// prepare request rack model
	reqModel := model.Rack{Name: args[0], RowID: row.ID}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", Conf.APIServer.URL+"/datacenter/rack", reqJSON)
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

func updateDataCenter(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	url = url + "/" + args[0]
	reqModel := model.DataCenter{Address: cmd.Flag("address").Value.String()}
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
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateDataCenterFloor(cmd *cobra.Command, args []string) error {
	floorName := cmd.Flag("name").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	dcID := strconv.Itoa(int(dc.ID))
	body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	data := new(model.Floors)
	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("response parse error:" + err.Error())
	}
	floor := model.Floor{}
	for _, flr := range *data {
		if flr.Name == args[0] {
			floor = flr
		}
	}
	if floor.ID == 0 {
		return fmt.Errorf("floor not found")
	}
	if floorName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	floor.Name = floorName
	reqJSON, _ := json.Marshal(floor)
	floorID := strconv.Itoa(int(floor.ID))
	url := Conf.APIServer.URL + "/datacenter/floor/" + floorID
	body, err = sendRequest("PUT", url, reqJSON)
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

func updateDataCenterHall(cmd *cobra.Command, args []string) error {
	hallName := cmd.Flag("name").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	if hallName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
	if err != nil {
		return err
	}

	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	dcID := strconv.Itoa(int(dc.ID))
	body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	data := new(model.Floors)
	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("response parse error:" + err.Error())
	}
	floor := model.Floor{}
	for _, flr := range *data {
		if flr.Name == floorName {
			floor = flr
			break
		}
	}
	if floor.ID == 0 {
		return fmt.Errorf("floor not found")
	}
	floorID := strconv.Itoa(int(floor.ID))
	url := Conf.APIServer.URL + "/datacenter/floor/" + floorID
	body, err = sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	f := new(model.Floor)
	if err := json.Unmarshal(body, f); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	hall := model.Hall{}
	for _, h := range f.Halls {
		if h.Name == args[0] {
			hall = h
			break
		}
	}
	if hall.ID == 0 {
		return fmt.Errorf("hall not found")
	}
	hall.Name = hallName
	reqJSON, _ := json.Marshal(hall)
	hallID := strconv.Itoa(int(hall.ID))
	url = Conf.APIServer.URL + "/datacenter/hall/" + hallID
	body, err = sendRequest("PUT", url, reqJSON)
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

func updateRackRow(cmd *cobra.Command, args []string) error {
	rowName := cmd.Flag("name").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	if rowName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	rows := new([]model.RackRow)
	if err := json.Unmarshal(body, rows); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*rows) > 1 {
		return fmt.Errorf("multiple row found")
	}
	row := new(model.RackRow)
	for _, r := range *rows {
		row = &r
		break
	}
	row.Name = rowName
	reqJSON, _ := json.Marshal(row)
	rowID := strconv.Itoa(int(row.ID))
	url = Conf.APIServer.URL + "/datacenter/row/" + rowID
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateRack(cmd *cobra.Command, args []string) error {
	rackName := cmd.Flag("name").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()
	if rackName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	url := Conf.APIServer.URL + "/datacenter/rack?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&row=" + rowName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	racks := new(model.Racks)
	if err := json.Unmarshal(body, racks); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*racks) > 1 {
		return fmt.Errorf("multiple rack found")
	}
	rack := new(model.Rack)
	for _, r := range *racks {
		rack = &r
		break
	}
	rack.Name = rackName
	reqJSON, _ := json.Marshal(rack)
	rackID := strconv.Itoa(int(rack.ID))
	url = Conf.APIServer.URL + "/datacenter/rack/" + rackID
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteDataCenter(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter/" + args[0]
	body, reqErr := sendRequest("DELETE", url, []byte{})
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

func deleteDataCenterFloor(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	dcID := strconv.Itoa(int(dc.ID))
	body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	data := new(model.Floors)
	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("response parse error:" + err.Error())
	}
	floor := model.Floor{}
	for _, flr := range *data {
		if flr.Name == args[0] {
			floor = flr
		}
	}
	if floor.ID == 0 {
		return fmt.Errorf("floor not found")
	}
	floorID := strconv.Itoa(int(floor.ID))
	url := Conf.APIServer.URL + "/datacenter/floor/" + floorID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteDataCenterHall(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	dcID := strconv.Itoa(int(dc.ID))
	body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	data := new(model.Floors)
	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("response parse error:" + err.Error())
	}
	floor := model.Floor{}
	for _, flr := range *data {
		if flr.Name == floorName {
			floor = flr
		}
	}
	if floor.ID == 0 {
		return fmt.Errorf("floor not found")
	}
	floorID := strconv.Itoa(int(floor.ID))
	url := Conf.APIServer.URL + "/datacenter/floor/" + floorID

	f := new(model.Floor)
	body, err = sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, f); err != nil {
		return err
	}

	hall := model.Hall{}
	for _, h := range f.Halls {
		if h.Name == args[0] {
			hall = h
			break
		}
	}
	if hall.ID == 0 {
		return fmt.Errorf("hall not found")
	}
	hallID := strconv.Itoa(int(hall.ID))
	url = Conf.APIServer.URL + "/datacenter/hall/" + hallID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteRackRow(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()

	url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	rows := new([]model.RackRow)
	if err := json.Unmarshal(body, rows); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*rows) > 1 {
		return fmt.Errorf("multiple row found")
	}
	row := new(model.RackRow)
	for _, r := range *rows {
		row = &r
		break
	}
	rowID := strconv.Itoa(int(row.ID))
	url = Conf.APIServer.URL + "/datacenter/row/" + rowID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteRack(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()

	url := Conf.APIServer.URL + "/datacenter/rack?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&row=" + rowName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	racks := new(model.Racks)
	if err := json.Unmarshal(body, racks); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*racks) > 1 {
		return fmt.Errorf("multiple rack found")
	}
	rack := new(model.Rack)
	for _, r := range *racks {
		rack = &r
		break
	}
	rackID := strconv.Itoa(int(rack.ID))
	url = Conf.APIServer.URL + "/datacenter/rack/" + rackID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}
