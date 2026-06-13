package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/lansdownian/goalans/internal/api"
	"github.com/lansdownian/goalans/internal/ui/design"
)

func RenderMatchView(width, height int, title string, matchList list.Model, details *api.MatchDetails, sp spinner.Model, loading bool, lastError string, polling bool) string {
	leftW := max(width*38/100, 28)
	rightW := max(width-leftW-4, 30)
	contentW := leftW - 6

	leftTitle := design.RenderHeader(title, contentW)
	if polling {
		leftTitle += " " + LiveStyle().Render("● LIVE")
	}

	listView := matchList.View()
	if len(matchList.Items()) == 0 && !loading {
		listView = neonEmptyStyle.Width(contentW).Render("No matches found")
	}

	leftBody := lipgloss.JoinVertical(lipgloss.Left, leftTitle, "", listView)
	left := neonPanelStyle.Width(leftW).Height(height - 4).Render(leftBody)

	rightContent := renderDetails(details, loading, lastError, sp, contentW)
	rightTitle := design.RenderHeader("Match Details", rightW-6)
	rightBody := lipgloss.JoinVertical(lipgloss.Left, rightTitle, "", rightContent)
	right := neonPanelCyanStyle.Width(rightW).Height(height - 4).Render(rightBody)

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)
	header := design.RenderHeaderCentered("Goalans", width-2)
	help := menuHelp.Render("↑/↓ j/k · enter select · / filter · esc back · q quit")

	return lipgloss.JoinVertical(lipgloss.Left, header, body, help)
}

func renderDetails(details *api.MatchDetails, loading bool, lastError string, sp spinner.Model, width int) string {
	if loading && details == nil {
		return neonDimStyle.Render("Loading…") + " " + sp.View()
	}
	if lastError != "" && details == nil {
		return ErrorStyle().Render("Error: "+lastError) + "\n\n" + neonDimStyle.Render("Try --mock for offline demo.")
	}
	if details == nil {
		return neonEmptyStyle.Render("Select a match")
	}

	var b strings.Builder
	status := string(details.Status)
	if details.Status == api.MatchStatusFinished {
		status = "FT"
	}
	if details.LiveTime != nil && *details.LiveTime != "" {
		status = *details.LiveTime
	}
	score := "vs"
	if details.HomeScore != nil && details.AwayScore != nil {
		score = fmt.Sprintf("%d - %d", *details.HomeScore, *details.AwayScore)
	}

	b.WriteString(neonTeamStyle.Render(details.HomeTeam.Name+" "+score+" "+details.AwayTeam.Name) + "\n")
	b.WriteString(neonDimStyle.Render(details.League.Name+" · "+status) + "\n\n")

	if details.Venue != "" {
		b.WriteString(neonDimStyle.Render("Venue: "+details.Venue) + "\n")
	}
	if details.Referee != "" {
		b.WriteString(neonDimStyle.Render("Referee: "+details.Referee) + "\n")
	}
	if details.Venue != "" || details.Referee != "" {
		b.WriteString("\n")
	}

	if len(details.Events) > 0 {
		b.WriteString(neonHeaderStyle.Render("Events") + "\n")
		for _, e := range details.Events {
			icon := "·"
			switch e.Type {
			case "goal":
				icon = "⚽"
			case "card":
				icon = "▪"
			}
			player := e.Player
			if player == "" {
				player = "—"
			}
			assist := ""
			if e.Assist != "" {
				assist = neonDimStyle.Render(" (" + e.Assist + ")")
			}
			b.WriteString(fmt.Sprintf("  %s %d' %s%s\n", icon, e.Minute, player, assist))
		}
		b.WriteString("\n")
	}

	if len(details.Statistics) > 0 {
		b.WriteString(neonHeaderStyle.Render("Statistics") + "\n")
		for _, s := range details.Statistics {
			if s.Label == "" {
				continue
			}
			b.WriteString(fmt.Sprintf("  %-18s %6s  %6s\n", s.Label, s.HomeValue, s.AwayValue))
		}
	}
	if loading {
		b.WriteString("\n" + neonDimStyle.Render("Updating…") + " " + sp.View())
	}
	if lastError != "" {
		b.WriteString("\n\n" + ErrorStyle().Render("Could not load full details: "+lastError))
	}
	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
