package src

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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
	baseName := strings.ToLower(filepath.Base(path))

	// Compressed archives
	switch {
	case ext == ".zip" || ext == ".jar" || ext == ".war":
		return renderZip(path, maxLines, lang)
	case ext == ".tar":
		return renderTar(path, "", maxLines, lang)
	case ext == ".tgz" || strings.HasSuffix(baseName, ".tar.gz"):
		return renderTar(path, "gzip", maxLines, lang)
	case ext == ".tbz2" || strings.HasSuffix(baseName, ".tar.bz2"):
		return renderTar(path, "bzip2", maxLines, lang)
	case ext == ".gz" || ext == ".gzip":
		return renderGzip(path, maxLines, maxWidth, lang)
	case ext == ".bz2":
		return renderBzip2(path, maxLines, maxWidth, lang)
	}

	// Read file contents for inspection
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf(d.ErrReadFile, err)
	}

	if len(data) == 0 {
		return d.EmptyFile
	}

	// Specific data formats
	switch {
	case ext == ".json" || ext == ".json5" || ext == ".jsonc":
		return renderJSON(data, maxLines, lang)
	case ext == ".csv":
		return renderCSV(data, ',', maxLines, lang)
	case ext == ".tsv":
		return renderCSV(data, '\t', maxLines, lang)
	case ext == ".md" || ext == ".markdown":
		return renderMarkdown(data, filepath.Base(path), maxLines, maxWidth, lang)
	case isImage(ext):
		return renderImage(data, filepath.Base(path), ext, info, maxLines, lang)
	case isMedia(ext):
		return renderMediaCard(filepath.Base(path), ext, info, lang)
	case isDocument(ext):
		return renderDocumentCard(data, filepath.Base(path), ext, info, lang)
	}

	// Binary check
	if isBinary(data) {
		return renderBinary(data, info, lang)
	}

	// Code or structured text
	if isCode(ext) || isSpecialTextFile(baseName) {
		return renderCode(data, ext, baseName, maxLines, lang)
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
		tagText, tagStyle := GetFileBadge(entry.Name(), entry.IsDir())
		if err == nil && !entry.IsDir() {
			sizeStr = formatSize(info.Size())
		}
		badge := tagStyle.Render(tagText)
		sb.WriteString(fmt.Sprintf("  %s %-9s %s\n", badge, sizeStr, entry.Name()))
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

func renderTar(path string, compression string, maxLines int, lang Language) string {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Sprintf("Error opening archive: %v", err)
	}
	defer f.Close()

	var r io.Reader = f
	if compression == "gzip" {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return fmt.Sprintf("Error reading gzip archive: %v", err)
		}
		defer gz.Close()
		r = gz
	} else if compression == "bzip2" {
		r = bzip2.NewReader(f)
	}

	tarReader := tar.NewReader(r)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TAR Archive: %s\n", filepath.Base(path)))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	count := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			sb.WriteString(fmt.Sprintf("  [Error reading archive entry: %v]\n", err))
			break
		}
		count++
		if count <= maxLines {
			name := header.Name
			if len(name) > 30 {
				name = "..." + name[len(name)-27:]
			}
			typeStr := "FILE"
			if header.Typeflag == tar.TypeDir {
				typeStr = "DIR "
			}
			sb.WriteString(fmt.Sprintf("  [%s] %-30s %10s\n", typeStr, name, formatSize(header.Size)))
		}
	}

	if count > maxLines {
		sb.WriteString(fmt.Sprintf("  ... and %d more item(s)\n", count-maxLines))
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

func renderBzip2(path string, maxLines, maxWidth int, lang Language) string {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Sprintf("Error opening BZIP2: %v", err)
	}
	defer f.Close()

	bz := bzip2.NewReader(f)
	data, err := io.ReadAll(bz)
	if err != nil {
		return fmt.Sprintf("Error reading BZIP2: %v", err)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("BZIP2 Archive: %s (uncompressed: %s)\n", filepath.Base(path), formatSize(int64(len(data)))))
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

func isImage(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".ico", ".svg", ".tiff", ".tif":
		return true
	default:
		return false
	}
}

func isMedia(ext string) bool {
	switch ext {
	case ".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a", ".wma", ".opus",
		".mp4", ".mkv", ".avi", ".mov", ".webm", ".flv", ".wmv", ".m4v":
		return true
	default:
		return false
	}
}

func isDocument(ext string) bool {
	switch ext {
	case ".pdf", ".epub", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".odt", ".ods":
		return true
	default:
		return false
	}
}

func renderImage(data []byte, filename, ext string, info os.FileInfo, maxLines int, lang Language) string {
	if ext == ".svg" {
		return renderCode(data, ext, filename, maxLines, lang)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Image File: %s\n", filename))
	sb.WriteString(strings.Repeat("-", 45) + "\n")

	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err == nil {
		aspect := "Custom"
		if cfg.Height > 0 {
			r := float64(cfg.Width) / float64(cfg.Height)
			if r > 1.7 && r < 1.8 {
				aspect = "16:9 (Widescreen)"
			} else if r > 1.3 && r < 1.4 {
				aspect = "4:3 (Standard)"
			} else if r > 0.95 && r < 1.05 {
				aspect = "1:1 (Square)"
			}
		}
		sb.WriteString(fmt.Sprintf("  - Format     : %s\n", strings.ToUpper(format)))
		sb.WriteString(fmt.Sprintf("  - Dimensions : %d × %d pixels\n", cfg.Width, cfg.Height))
		sb.WriteString(fmt.Sprintf("  - Aspect     : %s\n", aspect))
	} else {
		mimeType := http.DetectContentType(data)
		sb.WriteString(fmt.Sprintf("  - Type       : %s\n", mimeType))
	}

	if info != nil {
		sb.WriteString(fmt.Sprintf("  - Size       : %s (%d bytes)\n", formatSize(info.Size()), info.Size()))
		sb.WriteString(fmt.Sprintf("  - Permissions: %s\n", info.Mode()))
	}

	sb.WriteString("\n  ┌────────────────────────┐\n")
	sb.WriteString("  │      [IMAGE DATA]      │\n")
	sb.WriteString("  │  Binary visual asset   │\n")
	sb.WriteString("  └────────────────────────┘\n")

	return sb.String()
}

func renderMediaCard(filename, ext string, info os.FileInfo, lang Language) string {
	var sb strings.Builder
	mediaType := "Audio"
	if ext == ".mp4" || ext == ".mkv" || ext == ".avi" || ext == ".mov" || ext == ".webm" || ext == ".flv" || ext == ".wmv" {
		mediaType = "Video"
	}

	sb.WriteString(fmt.Sprintf("%s Media Asset: %s\n", mediaType, filename))
	sb.WriteString(strings.Repeat("-", 45) + "\n")
	sb.WriteString(fmt.Sprintf("  - Format     : %s\n", strings.ToUpper(strings.TrimPrefix(ext, "."))))
	if info != nil {
		sb.WriteString(fmt.Sprintf("  - Size       : %s (%d bytes)\n", formatSize(info.Size()), info.Size()))
		sb.WriteString(fmt.Sprintf("  - Permissions: %s\n", info.Mode()))
		sb.WriteString(fmt.Sprintf("  - Modified   : %s\n", info.ModTime().Format("2006-01-02 15:04:05")))
	}

	sb.WriteString("\n  ┌────────────────────────┐\n")
	sb.WriteString(fmt.Sprintf("  │      [%-5s ASSET]     │\n", strings.ToUpper(mediaType)))
	sb.WriteString("  │ Multimedia stream data │\n")
	sb.WriteString("  └────────────────────────┘\n")

	return sb.String()
}

func renderDocumentCard(data []byte, filename, ext string, info os.FileInfo, lang Language) string {
	var sb strings.Builder
	docType := strings.ToUpper(strings.TrimPrefix(ext, "."))
	sb.WriteString(fmt.Sprintf("Document: %s (%s)\n", filename, docType))
	sb.WriteString(strings.Repeat("-", 45) + "\n")

	if ext == ".pdf" && bytes.HasPrefix(data, []byte("%PDF-")) {
		version := "PDF"
		if len(data) >= 8 {
			version = string(data[:8])
		}
		sb.WriteString(fmt.Sprintf("  - Specification : %s\n", version))
	}

	if info != nil {
		sb.WriteString(fmt.Sprintf("  - Size          : %s (%d bytes)\n", formatSize(info.Size()), info.Size()))
		sb.WriteString(fmt.Sprintf("  - Permissions   : %s\n", info.Mode()))
		sb.WriteString(fmt.Sprintf("  - Modified      : %s\n", info.ModTime().Format("2006-01-02 15:04:05")))
	}

	sb.WriteString("\n  ┌────────────────────────┐\n")
	sb.WriteString(fmt.Sprintf("  │     [%-4s DOCUMENT]    │\n", docType))
	sb.WriteString("  │ Digital document file  │\n")
	sb.WriteString("  └────────────────────────┘\n")

	return sb.String()
}

func getLanguageName(ext, baseName string) string {
	switch baseName {
	case "dockerfile", "containerfile":
		return "Dockerfile"
	case "makefile", "gnumakefile":
		return "Makefile"
	case "cmakelists.txt":
		return "CMake"
	case ".gitignore", ".gitattributes", ".gitconfig":
		return "Git Config"
	}

	switch ext {
	case ".go":
		return "Go"
	case ".c", ".h":
		return "C"
	case ".cpp", ".cc", ".cxx", ".hpp", ".hxx", ".hh", ".tpp", ".inl":
		return "C++"
	case ".rs":
		return "Rust"
	case ".py", ".pyw", ".ipynb":
		return "Python"
	case ".js", ".mjs", ".cjs":
		return "JavaScript"
	case ".ts", ".mts", ".cts":
		return "TypeScript"
	case ".jsx":
		return "React JSX"
	case ".tsx":
		return "React TSX"
	case ".vue":
		return "Vue"
	case ".svelte":
		return "Svelte"
	case ".astro":
		return "Astro"
	case ".java":
		return "Java"
	case ".kt", ".kts":
		return "Kotlin"
	case ".scala":
		return "Scala"
	case ".groovy", ".gradle":
		return "Groovy / Gradle"
	case ".cs", ".csx":
		return "C#"
	case ".fs", ".fsx":
		return "F#"
	case ".vb":
		return "Visual Basic"
	case ".swift":
		return "Swift"
	case ".m", ".mm":
		return "Objective-C"
	case ".dart":
		return "Dart"
	case ".zig":
		return "Zig"
	case ".nim":
		return "Nim"
	case ".d":
		return "D"
	case ".odin":
		return "Odin"
	case ".v":
		return "V"
	case ".asm", ".s", ".nasm":
		return "Assembly"
	case ".rb", ".rake", ".gemspec":
		return "Ruby"
	case ".php", ".phtml":
		return "PHP"
	case ".lua":
		return "Lua"
	case ".pl", ".pm", ".t":
		return "Perl"
	case ".sh", ".bash", ".zsh", ".fish", ".ksh", ".csh":
		return "Shell"
	case ".bat", ".cmd":
		return "Batch"
	case ".ps1", ".psm1":
		return "PowerShell"
	case ".hs", ".lhs":
		return "Haskell"
	case ".ex", ".exs":
		return "Elixir"
	case ".erl", ".hrl":
		return "Erlang"
	case ".clj", ".cljs", ".cljc", ".edn":
		return "Clojure"
	case ".lisp", ".lsp", ".el", ".scm", ".rkt":
		return "Lisp / Scheme"
	case ".ml", ".mli":
		return "OCaml"
	case ".r", ".rmd":
		return "R"
	case ".jl":
		return "Julia"
	case ".html", ".htm", ".xhtml":
		return "HTML"
	case ".css", ".scss", ".sass", ".less", ".styl":
		return "CSS"
	case ".sql", ".sqlite", ".cql":
		return "SQL"
	case ".yaml", ".yml":
		return "YAML"
	case ".toml":
		return "TOML"
	case ".xml", ".plist":
		return "XML"
	case ".ini", ".cfg", ".conf", ".properties", ".env":
		return "Configuration"
	case ".proto":
		return "Protocol Buffers"
	case ".graphql", ".gql":
		return "GraphQL"
	case ".tex", ".latex":
		return "LaTeX"
	case ".rst":
		return "reStructuredText"
	default:
		return "Source Code"
	}
}

func isSpecialTextFile(baseName string) bool {
	switch baseName {
	case "dockerfile", "containerfile", ".dockerignore",
		"makefile", "gnumakefile", "cmakelists.txt",
		".gitignore", ".gitattributes", ".gitconfig",
		".editorconfig", "license", "licence", "gemfile",
		"rakefile", "vagrantfile", "brewfile", "procfile":
		return true
	default:
		return false
	}
}

func renderCode(data []byte, ext, baseName string, maxLines int, lang Language) string {
	d := T(lang)
	lines := strings.Split(string(data), "\n")
	var sb strings.Builder

	langName := getLanguageName(ext, baseName)
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

func detectMagicBinary(data []byte) string {
	if len(data) >= 4 {
		// ELF (Linux executable/shared object)
		if bytes.HasPrefix(data, []byte("\x7fELF")) {
			bits := "32-bit"
			if data[4] == 2 {
				bits = "64-bit"
			}
			return fmt.Sprintf("ELF %s Executable/Binary", bits)
		}
		// Windows PE / DOS
		if bytes.HasPrefix(data, []byte("MZ")) {
			return "Windows PE Executable / DLL"
		}
		// WebAssembly
		if bytes.HasPrefix(data, []byte("\x00asm")) {
			return "WebAssembly (WASM) Module"
		}
		// Java Class
		if bytes.HasPrefix(data, []byte("\xca\xfe\xba\xbe")) {
			return "Compiled Java Class Bytecode"
		}
		// Mach-O (macOS)
		if bytes.HasPrefix(data, []byte("\xfe\xed\xfa\xce")) || bytes.HasPrefix(data, []byte("\xce\xfa\xed\xfe")) ||
			bytes.HasPrefix(data, []byte("\xfe\xed\xfa\xcf")) || bytes.HasPrefix(data, []byte("\xcf\xfa\xed\xfe")) {
			return "Mach-O (macOS) Executable"
		}
		// PDF
		if bytes.HasPrefix(data, []byte("%PDF-")) {
			return "Adobe PDF Document"
		}
		// SQLite
		if bytes.HasPrefix(data, []byte("SQLite format 3\x00")) {
			return "SQLite Database File"
		}
	}
	return ""
}

func renderBinary(data []byte, info os.FileInfo, lang Language) string {
	d := T(lang)
	mimeType := http.DetectContentType(data)
	magicType := detectMagicBinary(data)
	var sb strings.Builder

	sb.WriteString(d.BinaryDetected)
	if magicType != "" {
		sb.WriteString(fmt.Sprintf("  - Format   : %s\n", magicType))
	}
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
	case ".go", ".c", ".h", ".cpp", ".cc", ".cxx", ".hpp", ".hxx", ".hh", ".tpp", ".inl",
		".rs", ".py", ".pyw", ".ipynb", ".js", ".mjs", ".cjs", ".ts", ".mts", ".cts", ".jsx", ".tsx",
		".vue", ".svelte", ".astro", ".java", ".kt", ".kts", ".scala", ".groovy", ".gradle",
		".cs", ".csx", ".fs", ".fsx", ".vb", ".swift", ".m", ".mm", ".dart",
		".zig", ".nim", ".d", ".odin", ".v", ".asm", ".s", ".nasm",
		".rb", ".rake", ".gemspec", ".php", ".phtml", ".lua", ".pl", ".pm", ".t",
		".sh", ".bash", ".zsh", ".fish", ".ksh", ".csh", ".bat", ".cmd", ".ps1", ".psm1",
		".hs", ".lhs", ".ex", ".exs", ".erl", ".hrl", ".clj", ".cljs", ".cljc", ".edn",
		".lisp", ".lsp", ".el", ".scm", ".rkt", ".ml", ".mli", ".r", ".rmd", ".jl",
		".html", ".htm", ".xhtml", ".css", ".scss", ".sass", ".less", ".styl",
		".xml", ".yaml", ".yml", ".toml", ".ini", ".cfg", ".conf", ".properties", ".env",
		".sql", ".sqlite", ".cql", ".proto", ".graphql", ".gql", ".tex", ".latex", ".rst":
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
