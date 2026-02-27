package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Set via ldflags: -X github.com/apermo/apermo-surf/cmd.version=...
var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("surf " + version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
