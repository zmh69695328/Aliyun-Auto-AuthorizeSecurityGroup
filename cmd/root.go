package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)
var rootCmd = &cobra.Command{
	Use:   "aaa",
	Short: "A command-line program for automatically adding Alibaba Cloud security group policy rules.	",
	// Long: ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println(`Use "aaa add [command] --help" for more information about a command.`)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
