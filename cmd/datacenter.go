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
