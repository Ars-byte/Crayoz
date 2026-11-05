package src

type Language string

const (
	LangEN Language = "en"
	LangES Language = "es"
)

func (l Language) Next() Language {
	if l == LangEN {
		return LangES
	}
	return LangEN
}

type Dict struct {
	Starting          string
	FullPreviewHeader string
	NoFileSelected    string
	ExplorerHeader    string
	EmptyDirectory    string
	GoUpLevel         string
	NoItems           string
	ItemStatus        string
	ExplorerFooter    string
	FullscreenHelp    string

	// About
	AboutTitle      string
	AboutSubtitle   string
	AboutVersion    string
	AboutProject    string
	AboutCreator    string
	AboutRepo       string
	AboutLanguage   string
	AboutYear       string
	AboutLicense    string
	AboutLicenseVal string
	AboutBuiltWith  string
	AboutReturnHint string
	AboutCloseHint  string

	// Errors
	ErrAccess         string
	ErrRunApp         string
	ErrGeneral        string
	ErrReadDir        string
	ErrOpenZip        string
	ErrOpenGzip       string
	ErrDecompressGzip string
	ErrReadGzip       string
	ErrReadFile       string

	// Preview formatting
	DirHeader        string
	MoreItems        string
	ZipHeader        string
	MoreFiles        string
	GzipHeader       string
	EmptyFile        string
	MoreRows         string
	RemainingLines   string
	BinaryDetected   string
	MimeType         string
	Size             string
	Permissions      string
	HexPreviewHeader string
	RemainingBytes   string
	MarkdownHeader   string
}

var dictionaries = map[Language]Dict{
	LangEN: {
		Starting:          "Starting...",
		FullPreviewHeader: "[FULL VIEW] %s (%s)",
		NoFileSelected:    "No file selected",
		ExplorerHeader:    "[EXPLORER] %s",
		EmptyDirectory:    "(Empty directory)",
		GoUpLevel:         "Go up one level to:\n%s",
		NoItems:           "No items",
		ItemStatus:        "[%d/%d] %s (%s)",
		ExplorerFooter:    "%s  │  [↑/↓] Move  [Enter/→] Open  [←] Up  [v] Fullscreen  [Tab/t] Lang: EN  [?] About  [q] Quit",
		FullscreenHelp:    "[↑/↓] Scroll  [PgUp/PgDn] Page  [Esc/←] Back  [Tab/t] Lang: EN  [?] About  [q] Quit",

		AboutTitle:      "crayoz",
		AboutSubtitle:   "Simple file reader, written in Go",
		AboutVersion:    "Version:",
		AboutProject:    "Project:",
		AboutCreator:    "Creator:",
		AboutRepo:       "Repository:",
		AboutLanguage:   "Language:",
		AboutYear:       "Year:",
		AboutLicense:    "License:",
		AboutLicenseVal: "MIT License",
		AboutBuiltWith:  "Made with:",
		AboutReturnHint: "Press any key to return",
		AboutCloseHint:  "[Esc / ? / a / Enter] Close",

		ErrAccess:         "Error: could not access %q: %v\n",
		ErrRunApp:         "Error running application: %v\n",
		ErrGeneral:        "Error: %v",
		ErrReadDir:        "Error reading directory: %v",
		ErrOpenZip:        "Error opening ZIP: %v",
		ErrOpenGzip:       "Error opening GZIP: %v",
		ErrDecompressGzip: "Error decompressing GZIP: %v",
		ErrReadGzip:       "Error reading GZIP: %v",
		ErrReadFile:       "Error reading file: %v",

		DirHeader:        "Directory: %s (%d items)\n",
		MoreItems:        "  ... and %d more item(s)\n",
		ZipHeader:        "ZIP Archive: %s (%d items)\n",
		MoreFiles:        "  ... and %d more file(s)\n",
		GzipHeader:       "GZIP Archive: %s (uncompressed: %s)\n",
		EmptyFile:        "Empty file (0 bytes)",
		MoreRows:         "... and %d more row(s)\n",
		RemainingLines:   "... [%d remaining lines]\n",
		BinaryDetected:   "Binary File Detected\n",
		MimeType:         "  - MIME Type  : %s\n",
		Size:             "  - Size       : %s (%d bytes)\n",
		Permissions:      "  - Permissions: %s\n",
		HexPreviewHeader: "\n--- Hexadecimal Preview ---\n",
		RemainingBytes:   "... [%d remaining bytes]\n",
		MarkdownHeader:   "Markdown Document: %s (%d lines)\n",
	},
	LangES: {
		Starting:          "Iniciando...",
		FullPreviewHeader: "[LECTURA COMPLETA] %s (%s)",
		NoFileSelected:    "No hay archivo seleccionado",
		ExplorerHeader:    "[EXPLORADOR] %s",
		EmptyDirectory:    "(Directorio vacío)",
		GoUpLevel:         "Subir un nivel hacia:\n%s",
		NoItems:           "Sin elementos",
		ItemStatus:        "[%d/%d] %s (%s)",
		ExplorerFooter:    "%s  │  [↑/↓] Mover  [Enter/→] Abrir  [←] Subir nivel  [v] Pantalla completa  [Tab/t] Idioma: ES  [?] Acerca de  [q] Salir",
		FullscreenHelp:    "[↑/↓] Desplazar  [PgUp/PgDn] Página  [Esc/←] Volver  [Tab/t] Idioma: ES  [?] Acerca de  [q] Salir",

		AboutTitle:      "crayoz",
		AboutSubtitle:   "Lector de archivos simple, escrito en Go",
		AboutVersion:    "Versión:",
		AboutProject:    "Proyecto:",
		AboutCreator:    "Creador:",
		AboutRepo:       "Repositorio:",
		AboutLanguage:   "Lenguaje:",
		AboutYear:       "Año:",
		AboutLicense:    "Licencia:",
		AboutLicenseVal: "Licencia MIT",
		AboutBuiltWith:  "Hecho con:",
		AboutReturnHint: "Pulsa cualquier tecla para volver",
		AboutCloseHint:  "[Esc / ? / a / Enter] Cerrar",

		ErrAccess:         "Error: no se pudo acceder a %q: %v\n",
		ErrRunApp:         "Error al ejecutar la aplicación: %v\n",
		ErrGeneral:        "Error: %v",
		ErrReadDir:        "Error al leer directorio: %v",
		ErrOpenZip:        "Error al abrir ZIP: %v",
		ErrOpenGzip:       "Error al abrir GZIP: %v",
		ErrDecompressGzip: "Error al descomprimir GZIP: %v",
		ErrReadGzip:       "Error al leer GZIP: %v",
		ErrReadFile:       "Error al leer archivo: %v",

		DirHeader:        "Directorio: %s (%d elementos)\n",
		MoreItems:        "  ... y %d elemento(s) más\n",
		ZipHeader:        "Archivo ZIP: %s (%d elementos)\n",
		MoreFiles:        "  ... y %d archivo(s) más\n",
		GzipHeader:       "Archivo GZIP: %s (descomprimido: %s)\n",
		EmptyFile:        "Archivo vacío (0 bytes)",
		MoreRows:         "... y %d fila(s) más\n",
		RemainingLines:   "... [%d líneas restantes]\n",
		BinaryDetected:   "Archivo Binario Detectado\n",
		MimeType:         "  - Tipo MIME : %s\n",
		Size:             "  - Tamaño    : %s (%d bytes)\n",
		Permissions:      "  - Permisos  : %s\n",
		HexPreviewHeader: "\n--- Vista Previa Hexadecimal ---\n",
		RemainingBytes:   "... [%d bytes restantes]\n",
		MarkdownHeader:   "Documento Markdown: %s (%d líneas)\n",
	},
}

func T(lang Language) Dict {
	if d, ok := dictionaries[lang]; ok {
		return d
	}
	return dictionaries[LangEN]
}
