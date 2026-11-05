# Crayoz

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/Ars-byte/-File-reader-in-go?color=brightgreen)](https://github.com/Ars-byte/-File-reader-in-go/releases)

> **Language / Idioma**: **English** | [Español](README.es.md)

**Crayoz** is a fast, lightweight, terminal-based file explorer and previewer built in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

---

## Features

- **Dual-Panel Explorer**: Browse directories with real-time split preview.
- **Fullscreen Reader**: Distraction-free file reading mode with line-by-line and page scrolling.
- **Multi-Format Previews**:
  - **Source Code**: Line-numbered preview for Go, Python, JavaScript, TypeScript, Rust, C/C++, Shell, Markdown, and more.
  - **Structured Data**: Pretty-printed JSON and column-aligned CSV/TSV tables.
  - **Archives**: Inspect `.zip` file tables of contents and preview `.gz` / `.gzip` compressed data.
  - **Binary Files**: Automatic MIME type detection, file metadata, permissions, and hexadecimal byte dump.
- **Instant Language Switching**: Toggle between English and Spanish instantly inside the TUI with `Tab` or `t`.
- **Zero External Runtime Dependencies**: Single self-contained binary.
- **Automated CI/CD Workflow**: Every update compiles a ready-to-run `.bin` binary and packages a `.tar.gz` distribution into GitHub Releases.

---

## Keybindings

| Key | Action |
| --- | --- |
| `↑` / `k` | Move cursor up / scroll up |
| `↓` / `j` | Move cursor down / scroll down |
| `Enter` / `→` / `l` | Enter folder / open file in fullscreen |
| `←` / `h` / `Backspace` | Go up one directory / exit fullscreen |
| `PgUp` / `b` | Page up scroll |
| `PgDn` / `f` / `Space` | Page down scroll |
| `v` | Toggle fullscreen reading view |
| `Tab` / `t` | **Toggle language (English ↔ Español)** |
| `Esc` | Return to explorer view from fullscreen |
| `q` / `Ctrl+C` | Quit application |

---

## Installation

### Pre-built Binaries (GitHub Releases)

Download the latest pre-compiled binary (`Crayoz.bin`) or the compressed tarball (`Crayoz-linux-amd64.tar.gz`) from the [Releases](https://github.com/Ars-byte/-File-reader-in-go/releases) page:

```bash
# Download and make executable
curl -L -o Crayoz https://github.com/Ars-byte/-File-reader-in-go/releases/latest/download/Crayoz.bin
chmod +x Crayoz

# Move to your PATH (optional)
sudo mv Crayoz /usr/local/bin/
```

Or extract the tarball:

```bash
tar -xzf Crayoz-linux-amd64.tar.gz
chmod +x Crayoz.bin
./Crayoz.bin
```

### Build from Source

Requirements: [Go](https://go.dev/dl/) 1.22 or higher.

```bash
# Clone the repository
git clone https://github.com/Ars-byte/-File-reader-in-go.git Crayoz
cd Crayoz

# Download dependencies
go mod download

# Build binary
go build -ldflags="-s -w" -o Crayoz .

# Run
./Crayoz
```

---

## Usage

```bash
# Open file explorer in current directory
./Crayoz

# Open in a specific directory
./Crayoz ~/Documents

# Open and focus on a specific file
./Crayoz ./hello.txt
```

---

## Automated Releases (CI/CD)

This repository includes a GitHub Actions workflow (`.github/workflows/release.yml`) that automatically:
1. Triggers on pushes to `main` and version tags (`v*`).
2. Builds an optimized, statically linked `.bin` executable.
3. Packages the binary, documentation, and license into a `.tar.gz` archive.
4. Generates SHA256 checksums.
5. Publishes the assets to GitHub Releases.

---

## License

This project is licensed under the [MIT License](LICENSE).
