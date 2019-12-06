package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hichikaw/dosanco/model"
	"github.com/spf13/cobra"
)

// NewCmdVersion is subcommand to show version information.
func NewCmdVersion() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of dosanco client",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dosanco command-line client %s, revision %s\n", version, revision)
		},
	}

	return versionCmd
}

func checkServerVersion(cmd *cobra.Command, args []string) error {
	srvVersion := new(model.Version)
	body, err := sendRequest("GET", Conf.APIServer.URL+"/version", []byte{})
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, srvVersion); err != nil {
		return fmt.Errorf("response parse error")
	}
	if srvVersion.Version != version {
		srvVerSlice := strings.Split(srvVersion.Version, ".")
		verSlice := strings.Split(version, ".")
		for i, v := range srvVerSlice {
			if v != verSlice[i] {
				if i == 0 {
					return fmt.Errorf(fmt.Sprintf("major version mismatch between dosano api server and client. Server:'%v' Client:'%v'", srvVersion.Version, version))
				} else if i == 1 {
					fmt.Printf("\033[33m")
					fmt.Printf(fmt.Sprintf("minor version mismatch between dosano api server and client. Server:'%v' Client:'%v'. Some commands might not be able to use.\n", srvVersion.Version, version))
					fmt.Printf("\033[0m\n")
					break
				} else {
					fmt.Printf("\033[36m")
					fmt.Printf(fmt.Sprintf("patch version mismatch between dosano api server and client. Server:'%v' Client:'%v'.\n", srvVersion.Version, version))
					fmt.Printf("\033[0m\n")
					break
				}
			}
		}
	}

	return nil
}
