package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/lansdownian/goalans/internal/ui/design"
)

var (
	menuItemStyle = lipgloss.NewStyle().Foreground(textColor)
	menuSelected  = lipgloss.NewStyle().Foreground(highlightColor).Bold(true)
	menuHelp      = lipgloss.NewStyle().Foreground(dimColor).Align(lipgloss.Center)
)

func RenderMenu(width, height, selected int, mock bool, loading bool, sp spinner.Model) string {
	items := []string{"Live Matches", "Finished Today", "World Cup 2026"}
	var menu strings.Builder
	for i, item := range items {
		if i == selected {
			menu.WriteString(menuSelected.Render("▸ " + item) + "\n")
		} else {
			menu.WriteString(menuItemStyle.Render("  "+item) + "\n")
		}
	}

	title := design.ApplyGradientToMultilineText("⚽ Goalans")
	subtitle := neonDimStyle.Render("Football scores in your terminal")

	var spinLine string
	if loading {
		spinLine = neonDimStyle.Render("Loading… ") + sp.View()
	}

	help := menuHelp.Render("↑/↓ or j/k · enter · q quit")
	if mock {
		help += neonDimStyle.Render(" · mock mode")
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		"", title, subtitle, "", spinLine, menu.String(), "", help,
	)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

func AccentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(neonCyan)
}

func DimStyle() lipgloss.Style {
	return neonDimStyle
}

func TitleStyle() lipgloss.Style {
	return neonHeaderStyle
}

func PanelStyle() lipgloss.Style {
	return neonPanelStyle
}

func ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(neonRed)
}

func LiveStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(neonRed).Bold(true)
}
