# File_reader

[![Versión de Go](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Licencia: MIT](https://img.shields.io/badge/Licencia-MIT-yellow.svg)](LICENSE)
[![Releases](https://img.shields.io/github/v/release/Ars-byte/-File-reader-in-go?color=brightgreen)](https://github.com/Ars-byte/-File-reader-in-go/releases)

> **Idioma / Language**: **Español** | [English](README.md)

**File_reader** es un explorador y visor de archivos rápido, ligero y para la terminal desarrollado en Go utilizando [Bubble Tea](https://github.com/charmbracelet/bubbletea) y [Lip Gloss](https://github.com/charmbracelet/lipgloss).

---

## Características

- **Explorador de doble panel**: Navega por carpetas con panel lateral de previsualización en tiempo real.
- **Lector a pantalla completa**: Modo de lectura sin distracciones con desplazamiento línea por línea y por página.
- **Previsualización de múltiples formatos**:
  - **Código fuente**: Vista numerada de líneas para Go, Python, JavaScript, TypeScript, Rust, C/C++, Shell, Markdown y más.
  - **Datos estructurados**: Formateo automático de JSON y tablas alineadas por columnas para CSV y TSV.
  - **Archivos comprimidos**: Inspección de contenido en archivos `.zip` y descompresión al vuelo de archivos `.gz` / `.gzip`.
  - **Archivos binarios**: Detección de tipo MIME, metadatos, permisos y volcado hexadecimal interactivo.
- **Cambio instantáneo de idioma**: Alterna entre inglés y español al instante dentro de la TUI pulsando `Tab` o `t`.
- **Sin dependencias externas en tiempo de ejecución**: Binario único y autocontenido.
- **Compilación y Releases automáticas (CI/CD)**: Cada cambio compila un ejecutable `.bin` y genera un paquete `.tar.gz` publicado en GitHub Releases.

---

## Atajos de teclado

| Tecla | Acción |
| --- | --- |
| `↑` / `k` | Mover cursor hacia arriba / desplazarse arriba |
| `↓` / `j` | Mover cursor hacia abajo / desplazarse abajo |
| `Enter` / `→` / `l` | Abrir carpeta / ver archivo en pantalla completa |
| `←` / `h` / `Backspace` | Subir un nivel de directorio / salir de pantalla completa |
| `PgUp` / `b` | Desplazar una página hacia arriba |
| `PgDn` / `f` / `Espacio` | Desplazar una página hacia abajo |
| `v` | Alternar modo pantalla completa |
| `Tab` / `t` | **Cambiar idioma de la TUI (Inglés ↔ Español)** |
| `Esc` | Volver al explorador desde pantalla completa |
| `q` / `Ctrl+C` | Salir del programa |

---

## Instalación

### Binarios precompilados (GitHub Releases)

Descarga el binario listo para usar (`File_reader.bin`) o el tarball comprimido (`File_reader-linux-amd64.tar.gz`) desde la sección de [Releases](https://github.com/Ars-byte/-File-reader-in-go/releases):

```bash
# Descargar y dar permisos de ejecución
curl -L -o File_reader https://github.com/Ars-byte/-File-reader-in-go/releases/latest/download/File_reader.bin
chmod +x File_reader

# Mover a tu PATH (opcional)
sudo mv File_reader /usr/local/bin/
```

O descomprime el tarball:

```bash
tar -xzf File_reader-linux-amd64.tar.gz
chmod +x File_reader.bin
./File_reader.bin
```

### Compilar desde el código fuente

Requisitos: [Go](https://go.dev/dl/) 1.22 o superior.

```bash
# Clonar el repositorio
git clone https://github.com/Ars-byte/-File-reader-in-go.git File_reader
cd File_reader

# Descargar dependencias
go mod download

# Compilar binario
go build -ldflags="-s -w" -o File_reader .

# Ejecutar
./File_reader
```

---

## Uso

```bash
# Iniciar en el directorio actual
./File_reader

# Iniciar en un directorio específico
./File_reader ~/Documentos

# Iniciar abriendo un archivo concreto
./File_reader ./hello.txt
```

---

## Releases Automáticas (CI/CD)

Este repositorio cuenta con un workflow de GitHub Actions (`.github/workflows/release.yml`) que automáticamente:
1. Se ejecuta con cada push a `main` o al crear etiquetas de versión (`v*`).
2. Compila un binario optimizado `.bin` para Linux.
3. Empaqueta el ejecutable, la documentación y la licencia en un archivo `.tar.gz`.
4. Genera sumas de verificación SHA256 (`checksums.txt`).
5. Publica todos los archivos en GitHub Releases.

---

## Licencia

Este proyecto está bajo la [Licencia MIT](LICENSE).
