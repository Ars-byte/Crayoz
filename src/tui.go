package src

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Name  string
	Path  string
	IsDir bool
	Size  int64
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

func NewModel(initialDir, selectTarget string) Model {
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
	m.loadDir(selectTarget)
	return m
}

func (m *Model) loadDir(selectTarget string) {
	d := T(m.lang)
	m.err = ""
	m.items = nil

	if m.currDir != "/" {
		m.items = append(m.items, Item{
			Name:  "..",
			Path:  filepath.Dir(m.currDir),
			IsDir: true,
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
			Name:  entry.Name(),
			Path:  filepath.Join(m.currDir, entry.Name()),
			IsDir: entry.IsDir(),
			Size:  size,
		}
		if entry.IsDir() {
			dirs = append(dirs, item)
		} else {
			files = append(files, item)
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	m.items = append(m.items, dirs...)
	m.items = append(m.items, files...)

	m.cursor = 0
	if selectTarget != "" {
		for i, item := range m.items {
			if item.Name == selectTarget {
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
	if item.Name == ".." {
		m.preview = fmt.Sprintf(d.GoUpLevel, item.Path)
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

	m.preview = GetPreview(item.Path, maxLines, rightWidth, m.lang)
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

func GetFileBadge(name string, isDir bool) (string, lipgloss.Style) {
	if isDir {
		if name == ".." {
			return "[UP  ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#a0aec0")).Bold(true)
		}
		return "[DIR ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#63b3ed")).Bold(true)
	}

	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".go":
		return "[GO  ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8")).Bold(true)
	case ".c", ".h":
		return "[C   ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#5c95d4")).Bold(true)
	case ".cpp", ".cc", ".cxx", ".hpp", ".hxx", ".hh":
		return "[C++ ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#00599c")).Bold(true)
	case ".rs":
		return "[RUST]", lipgloss.NewStyle().Foreground(lipgloss.Color("#dea584")).Bold(true)
	case ".py":
		return "[PY  ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd343")).Bold(true)
	case ".js", ".mjs", ".cjs":
		return "[JS  ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#f7df1e")).Bold(true)
	case ".ts", ".mts", ".cts":
		return "[TS  ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#3178c6")).Bold(true)
	case ".jsx", ".tsx":
		return "[REACT]", lipgloss.NewStyle().Foreground(lipgloss.Color("#61dafb")).Bold(true)
	case ".java":
		return "[JAVA]", lipgloss.NewStyle().Foreground(lipgloss.Color("#b07219")).Bold(true)
	case ".cs":
		return "[C#  ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#178600")).Bold(true)
	case ".rb":
		return "[RUBY]", lipgloss.NewStyle().Foreground(lipgloss.Color("#e0115f")).Bold(true)
	case ".php":
		return "[PHP ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#4f5d95")).Bold(true)
	case ".sh", ".bash", ".zsh":
		return "[BASH]", lipgloss.NewStyle().Foreground(lipgloss.Color("#4ebd38")).Bold(true)
	case ".html", ".htm":
		return "[HTML]", lipgloss.NewStyle().Foreground(lipgloss.Color("#e34c26")).Bold(true)
	case ".css", ".scss", ".sass", ".less":
		return "[CSS ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#805ad5")).Bold(true)
	case ".json":
		return "[JSON]", lipgloss.NewStyle().Foreground(lipgloss.Color("#48bb78")).Bold(true)
	case ".yaml", ".yml":
		return "[YAML]", lipgloss.NewStyle().Foreground(lipgloss.Color("#e53e3e")).Bold(true)
	case ".toml":
		return "[TOML]", lipgloss.NewStyle().Foreground(lipgloss.Color("#ed8936")).Bold(true)
	case ".xml":
		return "[XML ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#319795")).Bold(true)
	case ".sql":
		return "[SQL ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#9f7aea")).Bold(true)
	case ".md", ".markdown":
		return "[MD  ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#4299e1")).Bold(true)
	case ".zip", ".tar", ".gz", ".gzip", ".7z", ".rar":
		return "[ARCH]", lipgloss.NewStyle().Foreground(lipgloss.Color("#ed64a6")).Bold(true)
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".ico":
		return "[IMG ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#d53f8c")).Bold(true)
	case ".txt", ".log":
		return "[TXT ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#a0aec0"))
	case ".bin", ".exe", ".so", ".dll", ".o":
		return "[BIN ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#f56565")).Bold(true)
	default:
		return "[FILE]", lipgloss.NewStyle().Foreground(lipgloss.Color("#cbd5e0"))
	}
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
			if item.IsDir {
				target := ""
				if item.Name == ".." {
					target = filepath.Base(m.currDir)
				}
				m.currDir = item.Path
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
			if len(m.items) > 0 && !m.items[m.cursor].IsDir {
				m.fullScreen = !m.fullScreen
				m.fileScroll = 0
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	d := T(m.lang)
	if m.width == 0 || m.height == 0 {
		return d.Starting
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#2d3748")).
		Padding(0, 1).
		Width(m.width)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a0aec0")).
		Background(lipgloss.Color("#1a202c")).
		Padding(0, 1).
		Width(m.width)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#718096")).
		Background(lipgloss.Color("#1a202c")).
		Padding(0, 1).
		Width(m.width)

	if m.fullScreen {
		if len(m.items) == 0 {
			return d.NoFileSelected
		}
		item := m.items[m.cursor]
		header := headerStyle.Render(fmt.Sprintf(d.FullPreviewHeader, item.Name, formatSize(item.Size)))

		lines := strings.Split(GetPreview(item.Path, 5000, m.width-2, m.lang), "\n")
		contentRows := m.height - 3
		if contentRows < 1 {
			contentRows = 1
		}

		if m.fileScroll >= len(lines) {
			m.fileScroll = len(lines) - 1
		}
		if m.fileScroll < 0 {
			m.fileScroll = 0
		}

		end := m.fileScroll + contentRows
		if end > len(lines) {
			end = len(lines)
		}

		content := strings.Join(lines[m.fileScroll:end], "\n")
		help := helpStyle.Render(d.FullscreenHelp)
		return fmt.Sprintf("%s\n%s\n%s", header, content, help)
	}

	header := headerStyle.Render(fmt.Sprintf(d.ExplorerHeader, m.currDir))

	contentHeight := m.height - 3
	if contentHeight < 1 {
		contentHeight = 1
	}

	leftW := m.leftWidth()
	rightW := m.width - leftW - 3
	if rightW < 10 {
		rightW = 10
	}

	var leftLines []string
	endIndex := m.scroll + contentHeight
	if endIndex > len(m.items) {
		endIndex = len(m.items)
	}

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#63b3ed")).
		Width(leftW)


	for i := m.scroll; i < endIndex; i++ {
		item := m.items[i]
		tagText, tagStyle := GetFileBadge(item.Name, item.IsDir)

		displayName := item.Name
		maxNameLen := leftW - 9
		if len(displayName) > maxNameLen && maxNameLen > 3 {
			displayName = displayName[:maxNameLen-3] + "..."
		}

		if i == m.cursor {
			lineText := fmt.Sprintf(" %s %s", tagText, displayName)
			leftLines = append(leftLines, selectedStyle.Render(lineText))
		} else {
			styledTag := tagStyle.Render(tagText)
			nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#edf2f7"))
			if item.IsDir {
				nameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#63b3ed"))
			}
			lineText := fmt.Sprintf(" %s %s", styledTag, nameStyle.Render(displayName))
			leftLines = append(leftLines, lineText)
		}
	}

	for len(leftLines) < contentHeight {
		leftLines = append(leftLines, "")
	}

	rawPreviewLines := strings.Split(m.preview, "\n")
	var rightLines []string
	for i := 0; i < contentHeight; i++ {
		if i < len(rawPreviewLines) {
			line := rawPreviewLines[i]
			if len(line) > rightW {
				line = line[:rightW]
			}
			rightLines = append(rightLines, line)
		} else {
			rightLines = append(rightLines, "")
		}
	}

	var bodyBuilder strings.Builder
	for i := 0; i < contentHeight; i++ {
		leftPart := lipgloss.NewStyle().Width(leftW).Render(leftLines[i])
		rightPart := rightLines[i]
		bodyBuilder.WriteString(fmt.Sprintf("%s │ %s\n", leftPart, rightPart))
	}

	statusText := d.NoItems
	if len(m.items) > 0 {
		currItem := m.items[m.cursor]
		sizeInfo := "<DIR>"
		if !currItem.IsDir {
			sizeInfo = formatSize(currItem.Size)
		}
		statusText = fmt.Sprintf(d.ItemStatus, m.cursor+1, len(m.items), currItem.Name, sizeInfo)
	}

	footer := footerStyle.Render(fmt.Sprintf(d.ExplorerFooter, statusText))
	return fmt.Sprintf("%s\n%s%s", header, bodyBuilder.String(), footer)
}
