package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "surf",
	Short: "Project-contextual link navigation",
	Long:  "Open project URLs (environments, tools, docs) from the terminal with fuzzy matching.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
