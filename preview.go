package main

import (
	"fmt"
)

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

import (
	"encoding/hex"
	"net/http"
	"os"
	"unicode/utf8"
)

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

func isCode(ext string) bool {
	switch ext {
	case ".go", ".py", ".js", ".ts", ".jsx", ".tsx", ".c", ".cpp", ".h", ".hpp",
		".rs", ".java", ".rb", ".php", ".sh", ".bash", ".zsh", ".html", ".css",
		".scss", ".xml", ".yaml", ".yml", ".toml", ".ini", ".conf", ".sql",
		".dockerfile", ".makefile", ".env", ".lua", ".zig", ".dart", ".md":
		return true
	default:
		return false
	}
}

func renderCode(data []byte, ext string, maxLines int, lang Language) string {
	d := T(lang)
	lines := strings.Split(string(data), "\n")
	var sb strings.Builder

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
