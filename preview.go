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

import (
	"bytes"
	"encoding/csv"
)

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

import (
	"encoding/json"
)

func renderJSON(data []byte, maxLines int, lang Language) string {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		return renderPlainText(data, maxLines, lang)
	}
	return renderPlainText(pretty.Bytes(), maxLines, lang)
}

import (
	"archive/zip"
)

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

import (
	"compress/gzip"
	"io"
)

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
