package src

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAboutViewToggle(t *testing.T) {
	m := NewModel(".", "")

	// Initially showAbout should be false
	if m.showAbout {
		t.Fatalf("expected showAbout to be false initially")
	}

	// Pressing '?' should open about view
	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if cmd != nil {
		t.Fatalf("expected nil cmd")
	}
	m = updatedModel.(Model)
	if !m.showAbout {
		t.Fatalf("expected showAbout to be true after pressing '?'")
	}

	view := m.View()
	if !strings.Contains(view, "crayoz") {
		t.Fatalf("expected view to contain 'crayoz', got: %s", view)
	}

	// Pressing any key should close about view
	updatedModel, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	if cmd != nil {
		t.Fatalf("expected nil cmd")
	}
	m = updatedModel.(Model)
	if m.showAbout {
		t.Fatalf("expected showAbout to be false after pressing 'x'")
	}

	// Pressing 'a' should also open about view
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updatedModel.(Model)
	if !m.showAbout {
		t.Fatalf("expected showAbout to be true after pressing 'a'")
	}

	// Pressing Ctrl+C while in about view should quit
	_, quitCmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if quitCmd == nil {
		t.Fatalf("expected quitCmd when pressing Ctrl+C in about view")
	}
}
