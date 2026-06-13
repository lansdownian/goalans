package app

import (
	"context"
	"time"

	"github.com/lansdownian/goalans/internal/data"
	"github.com/lansdownian/goalans/internal/fotmob"
	tea "github.com/charmbracelet/bubbletea"
)

const pollInterval = 90 * time.Second

func fetchLiveMatches(client *fotmob.Client, useMock bool) tea.Cmd {
	return func() tea.Msg {
		if useMock {
			return matchesMsg{matches: data.MockLiveMatches()}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		matches, err := client.LiveMatches(ctx)
		return matchesMsg{matches: matches, err: err}
	}
}

func fetchFinishedMatches(client *fotmob.Client, useMock bool) tea.Cmd {
	return func() tea.Msg {
		if useMock {
			return matchesMsg{matches: data.MockFinishedMatches()}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		matches, err := client.FinishedMatches(ctx)
		return matchesMsg{matches: matches, err: err}
	}
}

func fetchMatchDetails(client *fotmob.Client, matchID int, useMock bool) tea.Cmd {
	return func() tea.Msg {
		if useMock {
			return matchDetailsMsg{details: data.MockMatchDetails(matchID)}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		details, err := client.MatchDetails(ctx, matchID)
		return matchDetailsMsg{details: details, err: err}
	}
}

func schedulePoll(matchID, gen int) tea.Cmd {
	return tea.Tick(pollInterval, func(time.Time) tea.Msg {
		return pollTickMsg{matchID: matchID, gen: gen}
	})
}
