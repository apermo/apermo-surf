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
	"github.com/apermo/apermo-surf/internal/userconfig"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:               "open [name] [ticket]",
	Short:             "Open a project link by fuzzy name",
	Args:              cobra.RangeArgs(0, 2),
	RunE:              runOpen,
	ValidArgsFunction: completeOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func completeOpen(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Only complete the first argument (link name)
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	path, err := config.Find(cwd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	cfg, err := config.Load(path)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	allLinks := cfg.AllLinks()
	var completions []string
	for name, link := range allLinks {
		completions = append(completions, fmt.Sprintf("%s\t%s", name, link.URL))
	}
	sort.Strings(completions)

	return completions, cobra.ShellCompDirectiveNoFileComp
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
	compoundMatched := false

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
		// Two args: try compound name first (e.g. "admin staging"),
		// then fall back to name + ticket semantics
		if len(args) == 2 {
			compound := args[0] + " " + args[1]
			match, _ = fuzzy.BestMatch(compound, names)
			compoundMatched = match != ""
		}

		if match == "" {
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
	}

	link := allLinks[match]
	configDir := filepath.Dir(path)

	var explicitArg string
	if len(args) == 2 && !compoundMatched {
		explicitArg = args[1]
	}
	result := resolve.Resolve(link, configDir, explicitArg)

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}

	fmt.Printf("opening %s → %s\n", match, result.URL)
	return browser.OpenWith(result.URL, browserFlag, userconfig.Load())
}
