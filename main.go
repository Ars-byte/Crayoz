package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	name  string
	path  string
	isDir bool
	size  int64
}

type Model struct {
	currDir    string
	items      []Item
	cursor     int
	scroll     int
	preview    string
	fullScreen bool
	fileScroll int
	width      int
	height     int
	err        string
	lang       Language
}

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

	var m Model
	if info.IsDir() {
		m = NewModel(target)
	} else {
		dir := filepath.Dir(target)
		m = NewModel(dir)
		base := filepath.Base(target)
		for i, item := range m.items {
			if item.name == base {
				m.cursor = i
				m.adjustScroll()
				m.updatePreview()
				break
			}
		}
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
