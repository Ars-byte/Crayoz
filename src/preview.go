package src

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

func GetPreview(path string, maxLines, maxWidth int, lang Language) string {
	d := T(lang)
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Sprintf(d.ErrGeneral, err)
	}

	if info.IsDir() {
		return renderDirectory(path, maxLines, lang)
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".zip":
		return renderZip(path, maxLines, lang)
	case ".gz", ".gzip":
		return renderGzip(path, maxLines, maxWidth, lang)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf(d.ErrReadFile, err)
	}

	if len(data) == 0 {
		return d.EmptyFile
	}

	switch ext {
	case ".json":
		return renderJSON(data, maxLines, lang)
	case ".csv":
		return renderCSV(data, ',', maxLines, lang)
	case ".tsv":
		return renderCSV(data, '\t', maxLines, lang)
	case ".md", ".markdown":
		return renderMarkdown(data, filepath.Base(path), maxLines, maxWidth, lang)
	}

	if isBinary(data) {
		return renderBinary(data, info, lang)
	}

	if isCode(ext) {
		return renderCode(data, ext, maxLines, lang)
	}

	return renderPlainText(data, maxLines, lang)
}

func renderDirectory(path string, maxLines int, lang Language) string {
	d := T(lang)
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Sprintf(d.ErrReadDir, err)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(d.DirHeader, filepath.Base(path), len(entries)))
	sb.WriteString(strings.Repeat("-", 45) + "\n")

	limit := len(entries)
	if limit > maxLines {
		limit = maxLines
	}

	for i := 0; i < limit; i++ {
		entry := entries[i]
		info, err := entry.Info()
		sizeStr := "<DIR>"
		tag := "[DIR] "
		if err == nil && !entry.IsDir() {
			sizeStr = formatSize(info.Size())
			tag = "[FILE]"
		}
		sb.WriteString(fmt.Sprintf("  %-6s %-9s %s\n", tag, sizeStr, entry.Name()))
	}

	if len(entries) > limit {
		sb.WriteString(fmt.Sprintf(d.MoreItems, len(entries)-limit))
	}

	return sb.String()
}

func renderZip(path string, maxLines int, lang Language) string {
	d := T(lang)
	r, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Sprintf(d.ErrOpenZip, err)
	}
	defer r.Close()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(d.ZipHeader, filepath.Base(path), len(r.File)))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	limit := len(r.File)
	if limit > maxLines {
		limit = maxLines
	}

	for i := 0; i < limit; i++ {
		f := r.File[i]
		name := f.Name
		if len(name) > 30 {
			name = "..." + name[len(name)-27:]
		}
		sb.WriteString(fmt.Sprintf("  %-30s %10s\n", name, formatSize(int64(f.UncompressedSize64))))
	}

	if len(r.File) > limit {
		sb.WriteString(fmt.Sprintf(d.MoreFiles, len(r.File)-limit))
	}

	return sb.String()
}

func renderGzip(path string, maxLines, maxWidth int, lang Language) string {
	d := T(lang)
	f, err := os.Open(path)
	if err != nil {
		return fmt.Sprintf(d.ErrOpenGzip, err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Sprintf(d.ErrDecompressGzip, err)
	}
	defer gz.Close()

	data, err := io.ReadAll(gz)
	if err != nil {
		return fmt.Sprintf(d.ErrReadGzip, err)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(d.GzipHeader, filepath.Base(path), formatSize(int64(len(data)))))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	if isBinary(data) {
		sb.WriteString(renderBinary(data, nil, lang))
	} else {
		sb.WriteString(renderPlainText(data, maxLines-2, lang))
	}
	return sb.String()
}

func renderJSON(data []byte, maxLines int, lang Language) string {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		return renderPlainText(data, maxLines, lang)
	}
	return renderPlainText(pretty.Bytes(), maxLines, lang)
}

func renderCSV(data []byte, delimiter rune, maxLines int, lang Language) string {
	d := T(lang)
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil || len(records) == 0 {
		return renderPlainText(data, maxLines, lang)
	}

	colWidths := make(map[int]int)
	for _, row := range records {
		for colIdx, val := range row {
			w := len(val)
			if w > 25 {
				w = 25
			}
			if w > colWidths[colIdx] {
				colWidths[colIdx] = w
			}
		}
	}

	var sb strings.Builder
	limit := len(records)
	if limit > maxLines {
		limit = maxLines
	}

	for rowIdx := 0; rowIdx < limit; rowIdx++ {
		row := records[rowIdx]
		sb.WriteString("| ")
		for colIdx, val := range row {
			w := colWidths[colIdx]
			displayVal := val
			if len(displayVal) > 25 {
				displayVal = displayVal[:22] + "..."
			}
			sb.WriteString(fmt.Sprintf("%-*s | ", w, displayVal))
		}
		sb.WriteString("\n")

		if rowIdx == 0 {
			sb.WriteString("+-")
			for colIdx := range row {
				w := colWidths[colIdx]
				sb.WriteString(strings.Repeat("-", w) + "-+-")
			}
			sb.WriteString("\n")
		}
	}

	if len(records) > limit {
		sb.WriteString(fmt.Sprintf(d.MoreRows, len(records)-limit))
	}

	return sb.String()
}

func renderMarkdown(data []byte, filename string, maxLines, maxWidth int, lang Language) string {
	d := T(lang)
	lines := strings.Split(string(data), "\n")
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(d.MarkdownHeader, filename, len(lines)))
	sb.WriteString(strings.Repeat("-", 45) + "\n")

	limit := len(lines)
	if limit > maxLines {
		limit = maxLines
	}

	h1Style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#63b3ed"))
	h2Style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#4fd1c5"))
	h3Style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f6ad55"))
	quoteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0aec0")).Italic(true)
	bulletStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#68d391"))
	ruleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#4a5568"))
	codeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f6e05e"))

	inCodeBlock := false

	for i := 0; i < limit; i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			sb.WriteString(codeStyle.Render("  "+line) + "\n")
			continue
		}

		if inCodeBlock {
			sb.WriteString(codeStyle.Render(fmt.Sprintf("  │ %s", line)) + "\n")
			continue
		}

		if strings.HasPrefix(trimmed, "# ") {
			sb.WriteString(h1Style.Render(line) + "\n")
		} else if strings.HasPrefix(trimmed, "## ") {
			sb.WriteString(h2Style.Render(line) + "\n")
		} else if strings.HasPrefix(trimmed, "### ") {
			sb.WriteString(h3Style.Render(line) + "\n")
		} else if strings.HasPrefix(trimmed, "> ") {
			sb.WriteString(quoteStyle.Render("  ▎ "+line[2:]) + "\n")
		} else if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			sb.WriteString(bulletStyle.Render("  • ") + line[2:] + "\n")
		} else if trimmed == "---" || trimmed == "***" {
			w := maxWidth - 4
			if w < 10 {
				w = 10
			}
			if w > 50 {
				w = 50
			}
			sb.WriteString(ruleStyle.Render(strings.Repeat("─", w)) + "\n")
		} else {
			sb.WriteString(line + "\n")
		}
	}

	if len(lines) > limit {
		sb.WriteString(fmt.Sprintf(d.RemainingLines, len(lines)-limit))
	}

	return sb.String()
}

func getLanguageName(ext string) string {
	switch ext {
	case ".go":
		return "Go"
	case ".c", ".h":
		return "C"
	case ".cpp", ".cc", ".cxx", ".hpp", ".hxx", ".hh":
		return "C++"
	case ".rs":
		return "Rust"
	case ".py":
		return "Python"
	case ".js", ".mjs", ".cjs":
		return "JavaScript"
	case ".ts", ".mts", ".cts":
		return "TypeScript"
	case ".jsx":
		return "React JSX"
	case ".tsx":
		return "React TSX"
	case ".java":
		return "Java"
	case ".cs":
		return "C#"
	case ".rb":
		return "Ruby"
	case ".php":
		return "PHP"
	case ".sh", ".bash", ".zsh":
		return "Shell"
	case ".html", ".htm":
		return "HTML"
	case ".css", ".scss", ".sass", ".less":
		return "CSS"
	case ".sql":
		return "SQL"
	case ".lua":
		return "Lua"
	case ".zig":
		return "Zig"
	case ".dart":
		return "Dart"
	case ".toml":
		return "TOML"
	case ".yaml", ".yml":
		return "YAML"
	case ".xml":
		return "XML"
	default:
		return "Source Code"
	}
}

func renderCode(data []byte, ext string, maxLines int, lang Language) string {
	d := T(lang)
	lines := strings.Split(string(data), "\n")
	var sb strings.Builder

	langName := getLanguageName(ext)
	sb.WriteString(fmt.Sprintf("[%s - %d lines]\n", langName, len(lines)))
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	limit := len(lines)
	if limit > maxLines {
		limit = maxLines
	}

	digitWidth := len(fmt.Sprintf("%d", len(lines)))
	if digitWidth < 3 {
		digitWidth = 3
	}

	for i := 0; i < limit; i++ {
		sb.WriteString(fmt.Sprintf("%*d | %s\n", digitWidth, i+1, lines[i]))
	}

	if len(lines) > limit {
		sb.WriteString(fmt.Sprintf(d.RemainingLines, len(lines)-limit))
	}

	return sb.String()
}

func renderPlainText(data []byte, maxLines int, lang Language) string {
	d := T(lang)
	lines := strings.Split(string(data), "\n")
	var sb strings.Builder

	limit := len(lines)
	if limit > maxLines {
		limit = maxLines
	}

	for i := 0; i < limit; i++ {
		sb.WriteString(lines[i] + "\n")
	}

	if len(lines) > limit {
		sb.WriteString(fmt.Sprintf(d.RemainingLines, len(lines)-limit))
	}

	return sb.String()
}

func renderBinary(data []byte, info os.FileInfo, lang Language) string {
	d := T(lang)
	mimeType := http.DetectContentType(data)
	var sb strings.Builder

	sb.WriteString(d.BinaryDetected)
	sb.WriteString(fmt.Sprintf(d.MimeType, mimeType))
	if info != nil {
		sb.WriteString(fmt.Sprintf(d.Size, formatSize(info.Size()), info.Size()))
		sb.WriteString(fmt.Sprintf(d.Permissions, info.Mode()))
	}
	sb.WriteString(d.HexPreviewHeader)

	limit := 192
	if len(data) < limit {
		limit = len(data)
	}
	sb.WriteString(hex.Dump(data[:limit]))
	if len(data) > limit {
		sb.WriteString(fmt.Sprintf(d.RemainingBytes, len(data)-limit))
	}

	return sb.String()
}

func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	sampleSize := 1024
	if len(data) < sampleSize {
		sampleSize = len(data)
	}
	for i := 0; i < sampleSize; i++ {
		if data[i] == 0 {
			return true
		}
	}
	return !utf8.Valid(data[:sampleSize])
}

func isCode(ext string) bool {
	switch ext {
	case ".go", ".c", ".h", ".cpp", ".cc", ".cxx", ".hpp", ".hxx", ".hh",
		".rs", ".py", ".js", ".mjs", ".cjs", ".ts", ".mts", ".cts", ".jsx", ".tsx",
		".java", ".cs", ".rb", ".php", ".sh", ".bash", ".zsh", ".fish", ".html",
		".css", ".scss", ".sass", ".less", ".xml", ".yaml", ".yml", ".toml",
		".ini", ".conf", ".sql", ".dockerfile", ".makefile", ".env", ".lua",
		".zig", ".dart", ".scala", ".asm", ".s", ".cmake":
		return true
	default:
		return false
	}
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
