package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of Aliyun-Auto-AuthorizeSecurityGroup",
  Long:  `All software has versions. This is Aliyun-Auto-AuthorizeSecurityGroup's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("v0.1")
  },
}