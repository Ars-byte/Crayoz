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

func (m Model) leftWidth() int {
	w := m.width * 35 / 100
	if w < 24 {
		w = 24
	}
	if w > 42 {
		w = 42
	}
	return w
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.adjustScroll()
		m.updatePreview()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab", "t":
			m.lang = m.lang.Next()
			m.updatePreview()

		case "up", "k":
			if m.fullScreen {
				if m.fileScroll > 0 {
					m.fileScroll--
				}
			} else if m.cursor > 0 {
				m.cursor--
				m.adjustScroll()
				m.updatePreview()
			}

		case "down", "j":
			if m.fullScreen {
				m.fileScroll++
			} else if m.cursor < len(m.items)-1 {
				m.cursor++
				m.adjustScroll()
				m.updatePreview()
			}

		case "pageup", "b":
			step := m.height - 6
			if step < 1 {
				step = 1
			}
			if m.fullScreen {
				m.fileScroll -= step
				if m.fileScroll < 0 {
					m.fileScroll = 0
				}
			} else {
				m.cursor -= step
				if m.cursor < 0 {
					m.cursor = 0
				}
				m.adjustScroll()
				m.updatePreview()
			}

		case "pagedown", "f", " ":
			step := m.height - 6
			if step < 1 {
				step = 1
			}
			if m.fullScreen {
				m.fileScroll += step
			} else {
				m.cursor += step
				if m.cursor >= len(m.items) {
					m.cursor = len(m.items) - 1
				}
				m.adjustScroll()
				m.updatePreview()
			}

		case "enter", "right", "l":
			if len(m.items) == 0 {
				break
			}
			item := m.items[m.cursor]
			if item.isDir {
				target := ""
				if item.name == ".." {
					target = filepath.Base(m.currDir)
				}
				m.currDir = item.path
				m.loadDir(target)
			} else {
				m.fullScreen = true
				m.fileScroll = 0
			}

		case "left", "h", "backspace":
			if m.fullScreen {
				m.fullScreen = false
			} else if m.currDir != "/" {
				prevBase := filepath.Base(m.currDir)
				m.currDir = filepath.Dir(m.currDir)
				m.loadDir(prevBase)
			}

		case "esc":
			if m.fullScreen {
				m.fullScreen = false
			}

		case "v":
			if len(m.items) > 0 && !m.items[m.cursor].isDir {
				m.fullScreen = !m.fullScreen
				m.fileScroll = 0
			}
		}
	}

	return m, nil
}
