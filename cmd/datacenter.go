package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/handler"
	"github.com/hichikaw/dosanco/model"
)

func NewCmdShowDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "show datacenter",
		Run: getDataCenter,
	}

	return dcCmd
}

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
	fmt.Printf("ID	Name	Address\n")
	for _, dc := range *data {
		fmt.Printf("%2d	%-10s	%s\n", dc.ID, dc.Name, dc.Address)
	}
}

func NewCmdCreateDataCenter() *cobra.Command {
	var dcCmd = &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "create new datacenter",
		//Long:    "create new network",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires data center name")
			}
			return nil
		},
		RunE: createDataCenter,
	}
	dcCmd.Flags().StringVarP(&address, "address", "a", "", "address of data center")
	dcCmd.MarkFlagRequired("address")

	return dcCmd
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