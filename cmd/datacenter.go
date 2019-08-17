package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/handler"
	"github.com/hichikaw/dosanco/model"
)

func getDataCenter(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.Url + "/datacenter"
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return
		}
		fmt.Println(errBody.Error.Message)
		return
	}
	data := new([]model.DataCenter)
	if err := json.Unmarshal(body, data); err != nil {
		fmt.Println("json unmarshall error:", err)
		return
	}
	fmt.Printf("%2s	%-10s	%s\n", "ID", "Name", "Address")
	for _, dc := range *data {
		fmt.Printf("%2d	%-10s	%s\n", dc.ID, dc.Name, dc.Address)
	}
}

func createDataCenter(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/datacenter"
	reqModel := model.DataCenter{Name: args[0], Address: address}
	reqJson, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url, reqJson)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	} else {
		fmt.Println(resMsg.Message)
	}

	return nil
}

func updateDataCenter(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/datacenter"
	url = url + "/" + args[0]
	reqModel := model.DataCenter{Address: address}
	reqJson, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("PUT", url, reqJson)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	} else {
		fmt.Println(resMsg.Message)
	}

	return nil
}

func deleteDataCenter(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/datacenter/" + args[0]
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return fmt.Errorf(resMsg.Message)
	} else {
		fmt.Println(resMsg.Message)
	}

	return nil
}
