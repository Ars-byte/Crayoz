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
