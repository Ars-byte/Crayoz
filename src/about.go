package src

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderAboutModal() string {
	d := T(m.lang)

	boxWidth := 46
	if m.width > 0 && boxWidth > m.width-4 {
		boxWidth = m.width - 4
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#63b3ed")).
		Padding(1, 2).
		Width(boxWidth).
		Background(lipgloss.Color("#1a202c"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#63b3ed"))

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#a0aec0")).
		Width(12)

	valStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#edf2f7"))

	closeHintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#718096"))

	var sb strings.Builder
	sb.WriteString(titleStyle.Render(d.AboutTitle) + "\n\n")
	sb.WriteString(labelStyle.Render(d.AboutProject) + valStyle.Render("Crayoz") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutCreator) + valStyle.Render("Ars-byte") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutLanguage) + valStyle.Render("Go") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutYear) + valStyle.Render("2026") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutLicense) + valStyle.Render("MIT License") + "\n\n")
	sb.WriteString(closeHintStyle.Render(d.AboutCloseHint))

	modal := boxStyle.Render(sb.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
}
