package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/model"
)

// Flags
var (
	tree        bool
	depth       int
	rfc         bool
	supernetID  int
	description string
)

func NewCmdShowNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net"},
		Short:   "show network",
		Run: func(cmd *cobra.Command, args []string) {
			url := Conf.APIServer.Url + "/network"
			query := ""
			if tree == true {
				query = query + "?tree=true"
				if depth > 0 {
					query = query + "&depth=" + strconv.Itoa(depth)
				}
			}
			if len(args) > 0 {
				getNetwork(url, args[0])
			} else {
				getNetworks(url, query)
			}
		},
	}
	networkCmd.Flags().BoolVarP(&tree, "tree", "t", false, "get network tree")
	networkCmd.Flags().BoolVarP(&rfc, "show-rfc-defined", "", false, "show networks defined and reserved in RFC")
	networkCmd.Flags().IntVarP(&depth, "depth", "d", 0, "depth for network tree. this option only work with --tree,-t option")

	return networkCmd
}

func NewCmdCreateNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net"},
		Short:   "create new network",
		Long:    "create new network",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires one network cidr")
			}
			return nil
		},
		RunE: createNetwork,
	}
	networkCmd.Flags().IntVarP(&supernetID, "supernet-id", "s", 0, "supernetwork id of the requested network")
	networkCmd.Flags().StringVarP(&description, "description", "d", "", "description of the requested network")
	networkCmd.MarkFlagRequired("supernet-id")

	return networkCmd
}

func getNetwork(url string, id string) {
	url = url + "/" + id
	body, err := sendRequest(url)
	if err != nil {
		fmt.Println(err.Error)
	}
	nw := new(model.IPv4Network)
	if err := json.Unmarshal(body, nw); err != nil {
		fmt.Println("json unmarshal error:", err)
		return
	}

	fmt.Printf("%2v %-20s	%s\n", nw.ID, nw.CIDR, nw.Description)
}

func getNetworks(url string, query string) {
	if query != "" {
		url = url + query
	}
	body, err := sendRequest(url)
	if err != nil {
		fmt.Println(err)
	}
	data := new([]model.IPv4Network)
	if err := json.Unmarshal(body, data); err != nil {
		fmt.Println("json unmarshal err:", err)
		return
	}

	if output == "json" {
		out := new(bytes.Buffer)
		json.Indent(out, body, "", "	")
		fmt.Println(out.String())
		return
	}
	if tree == true {
		printNetworkTree(data, 0)
		return
	}
	fmt.Printf("ID	CIDR			Description\n")
	for _, network := range *data {
		fmt.Printf("%2d	%-20s	%s\n", network.ID, network.CIDR, network.Description)
	}
}

func printNetworkTree(networks *[]model.IPv4Network, depth int) {
	for _, network := range *networks {
		fmt.Printf("%v%v:%v\n", strings.Repeat("   ", depth), network.ID, network.CIDR)
		if len(network.Subnetworks) > 0 {
			printNetworkTree(&network.Subnetworks, depth+1)
		}
	}
}

func createNetwork(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.Url + "/network"
	reqModel := model.IPv4Network{CIDR: args[0], SupernetworkID: uint(supernetID), Description: description}
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

func sendRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

type responseMessage struct {
	Message string `json:"message"`
}
