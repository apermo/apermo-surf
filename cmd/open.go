package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/apermo/apermo-surf/internal/browser"
	"github.com/apermo/apermo-surf/internal/config"
	"github.com/apermo/apermo-surf/internal/fuzzy"
	"github.com/apermo/apermo-surf/internal/resolve"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <name>",
	Short: "Open a project link by fuzzy name",
	Args:  cobra.ExactArgs(1),
	RunE:  runOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func runOpen(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path, err := config.Find(cwd)
	if err != nil {
		return err
	}

	cfg, err := config.Load(path)
	if err != nil {
		return err
	}

	allLinks := cfg.AllLinks()
	names := make([]string, 0, len(allLinks))
	for name := range allLinks {
		names = append(names, name)
	}
	sort.Strings(names)

	pattern := args[0]
	match, candidates := fuzzy.BestMatch(pattern, names)

	if match == "" && candidates == nil {
		return fmt.Errorf("no link matching %q — run surf links to see available links", pattern)
	}

	if match == "" {
		fmt.Fprintf(os.Stderr, "ambiguous match for %q:\n", pattern)
		for _, c := range candidates {
			fmt.Fprintf(os.Stderr, "  %s  %s\n", c, allLinks[c].URL)
		}
		return fmt.Errorf("be more specific or use the full name")
	}

	link := allLinks[match]
	configDir := filepath.Dir(path)
	result := resolve.Resolve(link, configDir)

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}

	fmt.Printf("opening %s → %s\n", match, result.URL)
	return browser.Open(result.URL)
}
