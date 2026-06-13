package api

import "time"

type MatchStatus string

const (
	MatchStatusNotStarted MatchStatus = "not_started"
	MatchStatusLive       MatchStatus = "live"
	MatchStatusFinished   MatchStatus = "finished"
)

type League struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type Match struct {
	ID        int         `json:"id"`
	League    League      `json:"league"`
	HomeTeam  Team        `json:"home_team"`
	AwayTeam  Team        `json:"away_team"`
	Status    MatchStatus `json:"status"`
	HomeScore *int        `json:"home_score,omitempty"`
	AwayScore *int        `json:"away_score,omitempty"`
	MatchTime *time.Time  `json:"match_time,omitempty"`
	LiveTime  *string     `json:"live_time,omitempty"`
	Round     string      `json:"round,omitempty"`
	PageURL   string      `json:"page_url,omitempty"`
}

type MatchEvent struct {
	Minute        int    `json:"minute"`
	DisplayMinute string `json:"display_minute,omitempty"`
	Type          string `json:"type"`
	Team          Team   `json:"team"`
	Player        string `json:"player,omitempty"`
	Assist        string `json:"assist,omitempty"`
}

type MatchDetails struct {
	Match
	Events     []MatchEvent `json:"events"`
	Venue      string       `json:"venue,omitempty"`
	Referee    string       `json:"referee,omitempty"`
	Statistics []Stat       `json:"statistics,omitempty"`
}

type Stat struct {
	Label     string `json:"label"`
	HomeValue string `json:"home_value"`
	AwayValue string `json:"away_value"`
}
