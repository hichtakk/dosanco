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
