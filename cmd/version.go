package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v1.0.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the boatswain version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
