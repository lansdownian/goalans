package app

import (
	"github.com/lansdownian/goalans/internal/api"
)

type matchesMsg struct {
	matches []api.Match
	err     error
}

type matchDetailsMsg struct {
	details *api.MatchDetails
	err     error
}

type pollTickMsg struct {
	matchID int
	gen     int
}

type wcDataMsg struct {
	data *api.WorldCupData
	err  error
}
