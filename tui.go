package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewModel(initialDir string) Model {
	absDir, err := filepath.Abs(initialDir)
	if err != nil {
		absDir = initialDir
	}
	m := Model{
		currDir: absDir,
		width:   80,
		height:  24,
		lang:    LangEN,
	}
	m.loadDir("")
	return m
}

func (m *Model) loadDir(selectTarget string) {
	d := T(m.lang)
	m.err = ""
	m.items = nil

	if m.currDir != "/" {
		m.items = append(m.items, Item{
			name:  "..",
			path:  filepath.Dir(m.currDir),
			isDir: true,
		})
	}

	entries, err := os.ReadDir(m.currDir)
	if err != nil {
		m.err = fmt.Sprintf(d.ErrReadDir, err)
		m.preview = m.err
		return
	}

	var dirs, files []Item
	for _, entry := range entries {
		info, err := entry.Info()
		var size int64
		if err == nil {
			size = info.Size()
		}
		item := Item{
			name:  entry.Name(),
			path:  filepath.Join(m.currDir, entry.Name()),
			isDir: entry.IsDir(),
			size:  size,
		}
		if entry.IsDir() {
			dirs = append(dirs, item)
		} else {
			files = append(files, item)
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].name) < strings.ToLower(dirs[j].name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].name) < strings.ToLower(files[j].name)
	})

	m.items = append(m.items, dirs...)
	m.items = append(m.items, files...)

	m.cursor = 0
	if selectTarget != "" {
		for i, item := range m.items {
			if item.name == selectTarget {
				m.cursor = i
				break
			}
		}
	}

	m.adjustScroll()
	m.updatePreview()
}

func (m *Model) adjustScroll() {
	visibleRows := m.height - 4
	if visibleRows < 1 {
		visibleRows = 1
	}
	if m.cursor < m.scroll {
		m.scroll = m.cursor
	}
	if m.cursor >= m.scroll+visibleRows {
		m.scroll = m.cursor - visibleRows + 1
	}
	if m.scroll < 0 {
		m.scroll = 0
	}
}

func (m *Model) updatePreview() {
	d := T(m.lang)
	if len(m.items) == 0 {
		m.preview = d.EmptyDirectory
		return
	}
	if m.cursor < 0 || m.cursor >= len(m.items) {
		m.cursor = 0
	}

	item := m.items[m.cursor]
	if item.name == ".." {
		m.preview = fmt.Sprintf(d.GoUpLevel, item.path)
		return
	}

	maxLines := m.height - 5
	if maxLines < 5 {
		maxLines = 5
	}
	rightWidth := m.width - m.leftWidth() - 3
	if rightWidth < 10 {
		rightWidth = 10
	}

	m.preview = GetPreview(item.path, maxLines, rightWidth, m.lang)
}
