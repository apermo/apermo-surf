package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/apermo/apermo-surf/internal/browser"
	"github.com/apermo/apermo-surf/internal/config"
	"github.com/apermo/apermo-surf/internal/fuzzy"
	"github.com/apermo/apermo-surf/internal/picker"
	"github.com/apermo/apermo-surf/internal/resolve"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [name] [ticket]",
	Short: "Open a project link by fuzzy name",
	Args:  cobra.RangeArgs(0, 2),
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

	var match string

	if len(args) == 0 {
		// Interactive picker mode
		urls := make([]string, len(names))
		for i, name := range names {
			urls[i] = allLinks[name].URL
		}
		idx, err := picker.Pick(names, urls)
		if err != nil {
			return err
		}
		match = names[idx]
	} else {
		pattern := args[0]
		var candidates []string
		match, candidates = fuzzy.BestMatch(pattern, names)

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
	}

	link := allLinks[match]
	configDir := filepath.Dir(path)

	var explicitArg string
	if len(args) == 2 {
		explicitArg = args[1]
	}
	result := resolve.Resolve(link, configDir, explicitArg)

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}

	fmt.Printf("opening %s → %s\n", match, result.URL)
	return browser.Open(result.URL)
}
