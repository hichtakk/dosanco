package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdShow() *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "show dosanco resources",
		//Long:  `All software has versions. This is Hugo's`,
	}
	showCmd.AddCommand(
		NewCmdShowNetwork(),
	)

	return showCmd
}

/*
func getOpts(args []string) {
	if len(args) == 0 {
		fmt.Println("no args")
		return
	}

	key := args[0]
	if len(args) > 1 {

	}
}
*/
