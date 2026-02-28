package picker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Pick presents an interactive selection from link names and URLs.
// Uses fzf when available, otherwise falls back to a numbered list on the terminal.
func Pick(names []string, urls []string) (int, error) {
	if hasFzf() {
		return pickWithFzf(names, urls)
	}
	return pickWithList(names, urls, os.Stdin, os.Stdout)
}

func hasFzf() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

func pickWithFzf(names []string, urls []string) (int, error) {
	var input strings.Builder
	for i, name := range names {
		fmt.Fprintf(&input, "%s\t%s\n", name, urls[i])
	}

	cmd := exec.Command("fzf",
		"--height=50%",
		"--layout=reverse",
		"--delimiter=\t",
		"--with-nth=1",
		"--preview=echo {2}",
	)
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return -1, fmt.Errorf("picker cancelled")
	}

	selected := strings.Split(strings.TrimSpace(string(out)), "\t")[0]
	for i, name := range names {
		if name == selected {
			return i, nil
		}
	}
	return -1, fmt.Errorf("picker cancelled")
}

// pickWithList shows a numbered list and reads the user's choice.
// Accepts io.Reader/Writer for testability.
func pickWithList(names []string, urls []string, in io.Reader, out io.Writer) (int, error) {
	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	for i, name := range names {
		fmt.Fprintf(out, "  %2d  %-*s  %s\n", i+1, maxLen, name, urls[i])
	}
	fmt.Fprint(out, "Pick a link [1]: ")

	scanner := bufio.NewScanner(in)
	if !scanner.Scan() {
		return -1, fmt.Errorf("picker cancelled")
	}

	line := strings.TrimSpace(scanner.Text())
	if line == "" {
		return 0, nil
	}

	var choice int
	if _, err := fmt.Sscanf(line, "%d", &choice); err != nil {
		return -1, fmt.Errorf("invalid choice: %q", line)
	}

	if choice < 1 || choice > len(names) {
		return -1, fmt.Errorf("choice out of range: %d", choice)
	}

	return choice - 1, nil
}
