package src

import (
	"strings"
	"testing"
)

func TestRenderAboutModalSpanish(t *testing.T) {
	m := Model{
		lang:   LangES,
		width:  80,
		height: 24,
	}

	out := m.renderAboutModal()

	expectedSubstrings := []string{
		"crayoz",
		"Lector de archivos simple, escrito en Go",
		"Versión:",
		"2.0.0",
		"Creador:",
		"Ars-byte",
		"Repositorio:",
		"github.com/Ars-byte/Crayoz",
		"Licencia:",
		"2026",
		"Hecho con:",
		"Go",
		"Pulsa cualquier tecla para volver",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(out, s) {
			t.Errorf("expected rendered about modal to contain %q, but got:\n%s", s, out)
		}
	}
}

func TestRenderAboutModalEnglish(t *testing.T) {
	m := Model{
		lang:   LangEN,
		width:  80,
		height: 24,
	}

	out := m.renderAboutModal()

	expectedSubstrings := []string{
		"crayoz",
		"Simple file reader, written in Go",
		"Version:",
		"2.0.0",
		"Creator:",
		"Ars-byte",
		"Repository:",
		"github.com/Ars-byte/Crayoz",
		"License:",
		"2026",
		"Made with:",
		"Go",
		"Press any key to return",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(out, s) {
			t.Errorf("expected rendered about modal to contain %q, but got:\n%s", s, out)
		}
	}
}
