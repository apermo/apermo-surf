package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "surf",
	Short: "Project-contextual link navigation",
	Long: `Open project URLs (environments, tools, docs) from the terminal with fuzzy matching.

Enable tab completion for link names:
  source <(surf completion zsh)    # add to ~/.zshrc
  source <(surf completion bash)   # add to ~/.bashrc
  surf completion fish | source    # or save to completions dir`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
