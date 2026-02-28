package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/apermo/apermo-surf/internal/config"
	"github.com/apermo/apermo-surf/internal/wizard"
	"github.com/spf13/cobra"
)

var distFlag bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new .surf-links.yml config",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().BoolVar(&distFlag, "dist", false, "create .surf-links.yml.dist (shared template)")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	w := wizard.New(os.Stdin, os.Stderr)

	dist := distFlag
	if !cmd.Flags().Changed("dist") {
		var err error
		dist, err = w.AskDist()
		if err != nil {
			return err
		}
	}

	fileName := config.FileName
	if dist {
		fileName = config.FileNameDist
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	target := filepath.Join(cwd, fileName)

	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("%s already exists", fileName)
	}

	cfg, err := w.Run()
	if err != nil {
		return err
	}

	if err := config.Write(cfg, target); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "created %s\n", fileName)
	if !dist {
		fmt.Fprintln(os.Stderr, "hint: add .surf-links.yml to .gitignore")
	}
	return nil
}
