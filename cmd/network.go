package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/model"
)

func NewCmdShowNetwork() *cobra.Command {
	var tree bool
	var depth int
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net"},
		Short:   "Print the version number of Hugo",
		//Long:    `All software has versions. This is Hugo's`,
		//Args: func(cmd *cobra.Command, args []string) error {
		//	if len(args) < 1 {
		//		// show list of network
		//		return errors.New("requires network resource name")
		//	}
		//	// show specified network
		//	return nil
		//},
		Run: func(cmd *cobra.Command, args []string) {
			url := "http://localhost:8080/network?"
			if tree == true {
				url = url + "tree=true"
			}
			if depth > 0 {
				url = url + "&depth=" + strconv.Itoa(depth)
			}
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
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
		},
	}
	networkCmd.Flags().BoolVarP(&tree, "tree", "t", false, "get network tree")
	networkCmd.Flags().IntVarP(&depth, "depth", "d", 0, "depth for tree network. this option work only with --tree,-t option")

	return networkCmd
}

func NewCmdCreateNetwork() *cobra.Command {
	var networkCmd = &cobra.Command{
		Use:     "network",
		Aliases: []string{"net"},
		Short:   "create new network",
		Long:    "create new network",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("create new network")
		},
	}

	return networkCmd
}

func printNetworkTree(networks *[]model.IPv4Network, depth int) {
	for _, network := range *networks {
		fmt.Printf("%v%v:%v\n", strings.Repeat("   ", depth), network.ID, network.CIDR)
		if len(network.Subnetworks) > 0 {
			printNetworkTree(&network.Subnetworks, depth+1)
		}
	}
}
