package ui

import (
	"fmt"
	"time"

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
	desc := m.Match.League.Name
	if m.Match.MatchTime != nil {
		desc += " · " + formatMatchDate(*m.Match.MatchTime)
	}
	return DimStyle().Render(desc)
}

func formatMatchDate(t time.Time) string {
	now := time.Now()
	local := t.In(now.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	matchDay := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, now.Location())
	switch {
	case matchDay.Equal(today):
		return "Today"
	case matchDay.Equal(today.AddDate(0, 0, -1)):
		return "Yesterday"
	default:
		return local.Format("Mon 2 Jan")
	}
}

func ToListItems(matches []api.Match) []list.Item {
	items := make([]list.Item, len(matches))
	for i, m := range matches {
		items[i] = MatchItem{Match: m, ID: m.ID}
	}
	return items
}
