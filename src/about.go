package src

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderAboutModal() string {
	d := T(m.lang)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff"))

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#a0aec0")).
		Width(14)

	valStyle := lipgloss.NewStyle().
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

	padLeft := 4
	padTop := 2
	if m.width < 50 {
		padLeft = 2
		padTop = 1
	}
	containerStyle := lipgloss.NewStyle().Padding(padTop, padLeft)
	return containerStyle.Render(sb.String())
}
