package browser

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/apermo/apermo-surf/internal/userconfig"
)

// builtinBrowsers maps well-known browser names to macOS/Linux/Windows commands.
var builtinBrowsers = map[string]map[string][]string{
	"chrome": {
		"darwin":  {"open", "-a", "Google Chrome"},
		"linux":   {"google-chrome"},
		"windows": {"cmd", "/c", "start", "chrome"},
	},
	"firefox": {
		"darwin":  {"open", "-a", "Firefox"},
		"linux":   {"firefox"},
		"windows": {"cmd", "/c", "start", "firefox"},
	},
	"safari": {
		"darwin": {"open", "-a", "Safari"},
	},
	"edge": {
		"darwin":  {"open", "-a", "Microsoft Edge"},
		"linux":   {"microsoft-edge"},
		"windows": {"cmd", "/c", "start", "msedge"},
	},
}

// Open opens the given URL in the system default browser.
func Open(url string) error {
	return OpenWith(url, "", userconfig.Config{})
}

// OpenWith opens a URL in a specific browser.
// Resolution order: flag (name) → config default → system default.
func OpenWith(url, name string, cfg userconfig.Config) error {
	browser := name
	if browser == "" {
		browser = cfg.Browser
	}

	if browser == "" {
		return openDefault(url)
	}

	// Check user-defined browsers first
	if bc, ok := cfg.Browsers[browser]; ok {
		args := append(bc.Args, url)
		return exec.Command(bc.Command, args...).Start()
	}

	// Check built-in browser map
	if platforms, ok := builtinBrowsers[browser]; ok {
		if cmdArgs, ok := platforms[runtime.GOOS]; ok {
			args := append(cmdArgs[1:], url)
			return exec.Command(cmdArgs[0], args...).Start()
		}
		return fmt.Errorf("browser %q is not available on %s", browser, runtime.GOOS)
	}

	return fmt.Errorf("unknown browser %q", browser)
}

func openDefault(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
