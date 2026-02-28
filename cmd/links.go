package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/apermo/apermo-surf/internal/config"
	"github.com/spf13/cobra"
)

var (
	flagEnv   bool
	flagTools bool
	flagDocs  bool
)

var linksCmd = &cobra.Command{
	Use:   "links",
	Short: "List project links",
	Long:  "List all or filtered links from .surf-links.yml.",
	RunE:  runLinks,
}

func init() {
	linksCmd.Flags().BoolVar(&flagEnv, "env", false, "show environments only")
	linksCmd.Flags().BoolVar(&flagTools, "tools", false, "show tools only")
	linksCmd.Flags().BoolVar(&flagDocs, "docs", false, "show docs only")
	rootCmd.AddCommand(linksCmd)
}

func runLinks(cmd *cobra.Command, args []string) error {
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

	if cfg.Name != "" {
		fmt.Printf("# %s\n\n", cfg.Name)
	}

	cats := cfg.Categories()
	filtered := filterCategories(cats)

	for i, cat := range filtered {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s:\n", cat.Name)
		printLinks(cat.Links)
	}

	return nil
}

func filterCategories(cats []config.Category) []config.Category {
	// No flags → show all
	if !flagEnv && !flagTools && !flagDocs {
		return cats
	}

	allowed := map[string]bool{
		"environments": flagEnv,
		"tools":        flagTools,
		"docs":         flagDocs,
	}

	var out []config.Category
	for _, cat := range cats {
		if allowed[cat.Name] {
			out = append(out, cat)
		}
	}
	return out
}

func printLinks(links map[string]config.Link) {
	names := make([]string, 0, len(links))
	maxLen := 0
	for name := range links {
		names = append(names, name)
		if len(name) > maxLen {
			maxLen = len(name)
		}
		for sub := range links[name].Links {
			if len(sub)+2 > maxLen {
				maxLen = len(sub) + 2
			}
		}
	}
	sort.Strings(names)

	for _, name := range names {
		link := links[name]
		fmt.Printf("  %-*s  %s\n", maxLen, name, link.URL)
		if len(link.Links) > 0 {
			subNames := make([]string, 0, len(link.Links))
			for sub := range link.Links {
				subNames = append(subNames, sub)
			}
			sort.Strings(subNames)
			for _, sub := range subNames {
				subURL := strings.TrimRight(link.URL, "/") + link.Links[sub]
				fmt.Printf("    %-*s  %s\n", maxLen-2, sub, subURL)
			}
		}
	}
}
