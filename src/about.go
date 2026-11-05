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

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
		Width(14)

	valStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	var sb strings.Builder

	sb.WriteString(titleStyle.Render(d.AboutTitle) + "\n")
	sb.WriteString(subtitleStyle.Render(d.AboutSubtitle) + "\n")
	sb.WriteString(labelStyle.Render(d.AboutVersion) + valStyle.Render(Version) + "\n\n")

	sb.WriteString(labelStyle.Render(d.AboutCreator) + valStyle.Render("Ars-byte") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutRepo) + valStyle.Render("github.com/Ars-byte/Crayoz") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutLicense) + valStyle.Render(d.AboutLicenseVal) + "\n")
	sb.WriteString(labelStyle.Render(d.AboutYear) + valStyle.Render("2026") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutBuiltWith) + valStyle.Render("Go") + "\n\n")

	sb.WriteString(footerStyle.Render(d.AboutReturnHint))

	padLeft := 2
	padTop := 1
	if m.width >= 60 {
		padLeft = 4
		padTop = 2
	}

	containerStyle := lipgloss.NewStyle().Padding(padTop, padLeft)
	return containerStyle.Render(sb.String())
}
