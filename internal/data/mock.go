package data

import (
	"time"

	"github.com/lansdownian/goalans/internal/api"
)

// DefaultLeagues are fetched when no custom config exists.
var DefaultLeagues = []api.League{
	{ID: 47, Name: "Premier League", Country: "England"},
	{ID: 87, Name: "La Liga", Country: "Spain"},
	{ID: 54, Name: "Bundesliga", Country: "Germany"},
	{ID: 55, Name: "Serie A", Country: "Italy"},
	{ID: 53, Name: "Ligue 1", Country: "France"},
	{ID: 77, Name: "FIFA World Cup", Country: "International"},
}

func LeagueIDs() []int {
	ids := make([]int, len(DefaultLeagues))
	for i, l := range DefaultLeagues {
		ids[i] = l.ID
	}
	return ids
}

func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }

func MockLiveMatches() []api.Match {
	t := time.Now().UTC().Add(-45 * time.Minute)
	return []api.Match{
		{
			ID: 9001,
			League: api.League{ID: 47, Name: "Premier League", Country: "England"},
			HomeTeam: api.Team{ID: 1, Name: "Arsenal", ShortName: "Arsenal"},
			AwayTeam: api.Team{ID: 2, Name: "Chelsea", ShortName: "Chelsea"},
			Status: api.MatchStatusLive,
			HomeScore: intPtr(2), AwayScore: intPtr(1),
			MatchTime: &t, LiveTime: strPtr("67'"),
			PageURL: "/matches/arsenal-vs-chelsea/mock9001",
		},
		{
			ID: 9002,
			League: api.League{ID: 87, Name: "La Liga", Country: "Spain"},
			HomeTeam: api.Team{ID: 3, Name: "Barcelona", ShortName: "Barcelona"},
			AwayTeam: api.Team{ID: 4, Name: "Real Madrid", ShortName: "Real Madrid"},
			Status: api.MatchStatusLive,
			HomeScore: intPtr(0), AwayScore: intPtr(0),
			MatchTime: &t, LiveTime: strPtr("23'"),
			PageURL: "/matches/barcelona-vs-real-madrid/mock9002",
		},
	}
}

func MockFinishedMatches() []api.Match {
	today := time.Now().UTC().Add(-3 * time.Hour)
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Add(-2 * time.Hour)
	twoDaysAgo := time.Now().UTC().AddDate(0, 0, -2).Add(-4 * time.Hour)
	return []api.Match{
		{
			ID: 9003,
			League: api.League{ID: 54, Name: "Bundesliga", Country: "Germany"},
			HomeTeam: api.Team{ID: 5, Name: "Bayern Munich", ShortName: "Bayern"},
			AwayTeam: api.Team{ID: 6, Name: "Dortmund", ShortName: "Dortmund"},
			Status: api.MatchStatusFinished,
			HomeScore: intPtr(3), AwayScore: intPtr(2),
			MatchTime: &today,
			PageURL: "/matches/bayern-vs-dortmund/mock9003",
		},
		{
			ID: 9004,
			League: api.League{ID: 55, Name: "Serie A", Country: "Italy"},
			HomeTeam: api.Team{ID: 7, Name: "Inter", ShortName: "Inter"},
			AwayTeam: api.Team{ID: 8, Name: "Milan", ShortName: "Milan"},
			Status: api.MatchStatusFinished,
			HomeScore: intPtr(1), AwayScore: intPtr(1),
			MatchTime: &today,
			PageURL: "/matches/inter-vs-milan/mock9004",
		},
		{
			ID: 9005,
			League: api.League{ID: 77, Name: "FIFA World Cup", Country: "International"},
			HomeTeam: api.Team{ID: 9, Name: "Mexico", ShortName: "Mexico"},
			AwayTeam: api.Team{ID: 10, Name: "South Africa", ShortName: "South Africa"},
			Status: api.MatchStatusFinished,
			HomeScore: intPtr(2), AwayScore: intPtr(0),
			MatchTime: &twoDaysAgo,
			PageURL: "/matches/mexico-vs-south-africa/mock9005",
		},
		{
			ID: 9006,
			League: api.League{ID: 77, Name: "FIFA World Cup", Country: "International"},
			HomeTeam: api.Team{ID: 11, Name: "South Korea", ShortName: "South Korea"},
			AwayTeam: api.Team{ID: 12, Name: "Denmark", ShortName: "Denmark"},
			Status: api.MatchStatusFinished,
			HomeScore: intPtr(1), AwayScore: intPtr(0),
			MatchTime: &yesterday,
			PageURL: "/matches/south-korea-vs-denmark/mock9006",
		},
	}
}

func MockMatchDetails(id int) *api.MatchDetails {
	for _, m := range append(MockLiveMatches(), MockFinishedMatches()...) {
		if m.ID == id {
			d := &api.MatchDetails{Match: m}
			switch id {
			case 9001:
				d.Events = []api.MatchEvent{
					{Minute: 12, Type: "goal", Team: m.HomeTeam, Player: "Saka"},
					{Minute: 34, Type: "goal", Team: m.AwayTeam, Player: "Palmer"},
					{Minute: 58, Type: "goal", Team: m.HomeTeam, Player: "Ødegaard", Assist: "Martinelli"},
				}
				d.Venue = "Emirates Stadium"
				d.Referee = "Michael Oliver"
				d.Statistics = []api.Stat{
					{Label: "Possession", HomeValue: "58%", AwayValue: "42%"},
					{Label: "Shots", HomeValue: "14", AwayValue: "9"},
					{Label: "Shots on target", HomeValue: "6", AwayValue: "4"},
				}
			case 9002:
				d.Events = []api.MatchEvent{
					{Minute: 19, Type: "card", Team: m.AwayTeam, Player: "Modrić"},
				}
				d.Venue = "Camp Nou"
			case 9003:
				d.Events = []api.MatchEvent{
					{Minute: 11, Type: "goal", Team: m.HomeTeam, Player: "Kane"},
					{Minute: 44, Type: "goal", Team: m.AwayTeam, Player: "Reus"},
					{Minute: 72, Type: "goal", Team: m.HomeTeam, Player: "Musiala"},
					{Minute: 88, Type: "goal", Team: m.AwayTeam, Player: "Füllkrug"},
					{Minute: 90, Type: "goal", Team: m.HomeTeam, Player: "Coman"},
				}
				d.Venue = "Allianz Arena"
			case 9004:
				d.Events = []api.MatchEvent{
					{Minute: 55, Type: "goal", Team: m.HomeTeam, Player: "Lautaro"},
					{Minute: 79, Type: "goal", Team: m.AwayTeam, Player: "Leão"},
				}
				d.Venue = "San Siro"
			}
			return d
		}
	}
	return nil
}
