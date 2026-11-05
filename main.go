package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"Crayoz/src"
)

func main() {
	target := "."
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	target = strings.Trim(target, "\"'`")
	if target == "~" || strings.HasPrefix(target, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			if target == "~" {
				target = home
			} else {
				target = filepath.Join(home, target[2:])
			}
		}
	}

	info, err := os.Stat(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not access %q: %v\n", target, err)
		os.Exit(1)
	}

	var m src.Model
	if info.IsDir() {
		m = src.NewModel(target, "")
	} else {
		dir := filepath.Dir(target)
		base := filepath.Base(target)
		m = src.NewModel(dir, base)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
