package cmd

import (
	//"bytes"
	"encoding/json"
	//"errors"
	"fmt"
	//"io/ioutil"
	//"net/http"
	"strconv"
	//"strings"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/handler"
	"github.com/hichikaw/dosanco/model"
)

func showHost(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/host/name/" + args[0]
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
	host := new(model.Host)
	if err := json.Unmarshal(body, host); err != nil {
		fmt.Println("json unmarshal error:", err)
		return
	}
	host.Write(cmd.Flag("output").Value.String())
}

func createHost(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/host"
	name := args[0]
	location := cmd.Flag("location").Value.String()
	description := cmd.Flag("description").Value.String()
	reqModel := model.Host{Name: name, Location: location, Description: description}
	reqJSON, _ := json.Marshal(reqModel)
	body, err := sendRequest("POST", url, reqJSON)
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return err
		}
		fmt.Println(errBody.Error.Message)
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
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return err
		}
		fmt.Println(errBody.Error.Message)
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
	//reqModel := model.Host{}
	if name == "-" && description == "-" && location == "-" {
		fmt.Println("nothing to be updated")
		return fmt.Errorf("nothing to be updated")
	}
	if name != "-" {
		host.Name = name
	}
	// ensure the new name is not already exists in database
	url = Conf.APIServer.URL + "/host/name/" + name
	body, err = sendRequest("GET", url, []byte{})
	if err == nil {
		fmt.Printf("host '%v' is already exist\n", name)
		return fmt.Errorf("host '%v' is already exist", name)
	}
	if location != "-" {
		host.Location = location
	}
	if description != "-" {
		host.Description = description
	}
	reqJSON, _ := json.Marshal(host)
	url = Conf.APIServer.URL + "/host/" + id
	body, err = sendRequest("PUT", url, reqJSON)
	if err != nil {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return err
		}
		fmt.Println(errBody.Error.Message)
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
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return err
		}
		fmt.Println(errBody.Error.Message)
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
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			fmt.Println("response parse error")
			return err
		}
		fmt.Println(errBody.Error.Message)
		return err
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}
