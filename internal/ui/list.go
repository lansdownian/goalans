package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/lansdownian/goalans/internal/api"
)

type MatchItem struct {
	Match api.Match
	ID    int
}

func (m MatchItem) FilterValue() string {
	return fmt.Sprintf("%s %s %s", m.Match.HomeTeam.Name, m.Match.AwayTeam.Name, m.Match.League.Name)
}

func (m MatchItem) Title() string {
	home := m.Match.HomeTeam.ShortName
	away := m.Match.AwayTeam.ShortName
	score := "vs"
	if m.Match.HomeScore != nil && m.Match.AwayScore != nil {
		score = fmt.Sprintf("%d-%d", *m.Match.HomeScore, *m.Match.AwayScore)
	}
	live := ""
	if m.Match.LiveTime != nil && *m.Match.LiveTime != "" {
		live = " " + LiveStyle().Render(*m.Match.LiveTime)
	}
	return fmt.Sprintf("%s %s %s%s", home, score, away, live)
}

func (m MatchItem) Description() string {
	return DimStyle().Render(m.Match.League.Name)
}

func ToListItems(matches []api.Match) []list.Item {
	items := make([]list.Item, len(matches))
	for i, m := range matches {
		items[i] = MatchItem{Match: m, ID: m.ID}
	}
	return items
}
