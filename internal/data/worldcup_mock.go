package data

import "github.com/lansdownian/goalans/internal/api"

func wcInt(v int) *int { return &v }

// MockWorldCupData returns sample World Cup data for --mock mode.
// Based on the 2022 tournament structure for offline demo purposes.
func MockWorldCupData() *api.WorldCupData {
	arg := api.Team{ID: 6706, Name: "Argentina", ShortName: "ARG"}
	fra := api.Team{ID: 6723, Name: "France", ShortName: "FRA"}

	return &api.WorldCupData{
		Season: "2026",
		Name:   "FIFA World Cup 2026 (sample data)",
		Groups: []api.WCGroup{
			{
				Letter: "A", Name: "Group A",
				Teams: []api.LeagueTableEntry{
					{Position: 1, Team: api.Team{ID: 6713, Name: "USA", ShortName: "USA"}, Played: 3, Won: 2, Drawn: 1, Lost: 0, GoalsFor: 5, GoalsAgainst: 2, GoalDifference: 3, Points: 7},
					{Position: 2, Team: api.Team{ID: 6706, Name: "Argentina", ShortName: "ARG"}, Played: 3, Won: 2, Drawn: 0, Lost: 1, GoalsFor: 4, GoalsAgainst: 3, GoalDifference: 1, Points: 6},
					{Position: 3, Team: api.Team{ID: 6710, Name: "Mexico", ShortName: "MEX"}, Played: 3, Won: 1, Drawn: 1, Lost: 1, GoalsFor: 3, GoalsAgainst: 4, GoalDifference: -1, Points: 4},
					{Position: 4, Team: api.Team{ID: 5810, Name: "Canada", ShortName: "CAN"}, Played: 3, Won: 0, Drawn: 0, Lost: 3, GoalsFor: 2, GoalsAgainst: 5, GoalDifference: -3, Points: 0},
				},
			},
			{
				Letter: "B", Name: "Group B",
				Teams: []api.LeagueTableEntry{
					{Position: 1, Team: api.Team{ID: 8491, Name: "England", ShortName: "ENG"}, Played: 3, Won: 2, Drawn: 1, Lost: 0, GoalsFor: 6, GoalsAgainst: 2, GoalDifference: 4, Points: 7},
					{Position: 2, Team: api.Team{ID: 6723, Name: "France", ShortName: "FRA"}, Played: 3, Won: 2, Drawn: 0, Lost: 1, GoalsFor: 5, GoalsAgainst: 3, GoalDifference: 2, Points: 6},
					{Position: 3, Team: api.Team{ID: 6720, Name: "Spain", ShortName: "ESP"}, Played: 3, Won: 1, Drawn: 0, Lost: 2, GoalsFor: 3, GoalsAgainst: 5, GoalDifference: -2, Points: 3},
					{Position: 4, Team: api.Team{ID: 6711, Name: "Iran", ShortName: "IRN"}, Played: 3, Won: 0, Drawn: 1, Lost: 2, GoalsFor: 2, GoalsAgainst: 6, GoalDifference: -4, Points: 1},
				},
			},
			{
				Letter: "C", Name: "Group C",
				Teams: []api.LeagueTableEntry{
					{Position: 1, Team: api.Team{ID: 8256, Name: "Brazil", ShortName: "BRA"}, Played: 3, Won: 3, Drawn: 0, Lost: 0, GoalsFor: 7, GoalsAgainst: 2, GoalDifference: 5, Points: 9},
					{Position: 2, Team: api.Team{ID: 6715, Name: "Japan", ShortName: "JPN"}, Played: 3, Won: 2, Drawn: 0, Lost: 1, GoalsFor: 4, GoalsAgainst: 3, GoalDifference: 1, Points: 6},
					{Position: 3, Team: api.Team{ID: 6262, Name: "Morocco", ShortName: "MAR"}, Played: 3, Won: 1, Drawn: 0, Lost: 2, GoalsFor: 3, GoalsAgainst: 4, GoalDifference: -1, Points: 3},
					{Position: 4, Team: api.Team{ID: 6716, Name: "Australia", ShortName: "AUS"}, Played: 3, Won: 0, Drawn: 0, Lost: 3, GoalsFor: 1, GoalsAgainst: 6, GoalDifference: -5, Points: 0},
				},
			},
			{
				Letter: "D", Name: "Group D",
				Teams: []api.LeagueTableEntry{
					{Position: 1, Team: api.Team{ID: 8570, Name: "Germany", ShortName: "GER"}, Played: 3, Won: 2, Drawn: 1, Lost: 0, GoalsFor: 6, GoalsAgainst: 3, GoalDifference: 3, Points: 7},
					{Position: 2, Team: api.Team{ID: 10155, Name: "Croatia", ShortName: "CRO"}, Played: 3, Won: 2, Drawn: 0, Lost: 1, GoalsFor: 5, GoalsAgainst: 4, GoalDifference: 1, Points: 6},
					{Position: 3, Team: api.Team{ID: 8361, Name: "Portugal", ShortName: "POR"}, Played: 3, Won: 1, Drawn: 0, Lost: 2, GoalsFor: 4, GoalsAgainst: 5, GoalDifference: -1, Points: 3},
					{Position: 4, Team: api.Team{ID: 5796, Name: "Uruguay", ShortName: "URU"}, Played: 3, Won: 0, Drawn: 1, Lost: 2, GoalsFor: 2, GoalsAgainst: 5, GoalDifference: -3, Points: 1},
				},
			},
			{Letter: "E", Name: "Group E", Teams: mockGroupTeams("NED", "SEN", "ECU", "QAT")},
			{Letter: "F", Name: "Group F", Teams: mockGroupTeams("MAR", "CRO", "BEL", "CAN")},
			{Letter: "G", Name: "Group G", Teams: mockGroupTeams("BRA", "SUI", "CMR", "SRB")},
			{Letter: "H", Name: "Group H", Teams: mockGroupTeams("POR", "KOR", "URU", "GHA")},
			{Letter: "I", Name: "Group I", Teams: mockGroupTeams("ITA", "AUT", "SCO", "NGA")},
			{Letter: "J", Name: "Group J", Teams: mockGroupTeams("COL", "PER", "CHI", "PAR")},
			{Letter: "K", Name: "Group K", Teams: mockGroupTeams("TUR", "CZE", "UKR", "SWE")},
			{Letter: "L", Name: "Group L", Teams: mockGroupTeams("POL", "DEN", "SRB", "WAL")},
		},
		KnockoutRounds: []api.WCKnockoutRound{
			{
				Stage: "1/8", Label: "Round of 16",
				Matchups: []api.WCMatchup{
					{HomeTeam: "USA", HomeTeamID: 6713, HomeShort: "USA", AwayTeam: "Iran", AwayTeamID: 6711, AwayShort: "IRN", HomeScore: wcInt(2), AwayScore: wcInt(0), WinnerID: wcInt(6713)},
					{HomeTeam: "Argentina", HomeTeamID: 6706, HomeShort: "ARG", AwayTeam: "Australia", AwayTeamID: 6716, AwayShort: "AUS", HomeScore: wcInt(2), AwayScore: wcInt(1), WinnerID: wcInt(6706)},
					{HomeTeam: "England", HomeTeamID: 8491, HomeShort: "ENG", AwayTeam: "Senegal", AwayTeamID: 6395, AwayShort: "SEN", HomeScore: wcInt(3), AwayScore: wcInt(0), WinnerID: wcInt(8491)},
					{HomeTeam: "France", HomeTeamID: 6723, HomeShort: "FRA", AwayTeam: "Poland", AwayTeamID: 8568, AwayShort: "POL", HomeScore: wcInt(3), AwayScore: wcInt(1), WinnerID: wcInt(6723)},
				},
			},
			{
				Stage: "1/4", Label: "Quarterfinals",
				Matchups: []api.WCMatchup{
					{HomeTeam: "Argentina", HomeTeamID: 6706, HomeShort: "ARG", AwayTeam: "Netherlands", AwayTeamID: 6708, AwayShort: "NED", HomeScore: wcInt(2), AwayScore: wcInt(2), WinnerID: wcInt(6706), IsPenalties: true},
					{HomeTeam: "France", HomeTeamID: 6723, HomeShort: "FRA", AwayTeam: "England", AwayTeamID: 8491, AwayShort: "ENG", HomeScore: wcInt(2), AwayScore: wcInt(1), WinnerID: wcInt(6723)},
				},
			},
			{
				Stage: "final", Label: "Final",
				Matchups: []api.WCMatchup{
					{HomeTeam: "Argentina", HomeTeamID: 6706, HomeShort: "ARG", AwayTeam: "France", AwayTeamID: 6723, AwayShort: "FRA", HomeScore: wcInt(3), AwayScore: wcInt(3), WinnerID: wcInt(6706), IsPenalties: true},
				},
			},
		},
		BronzeFinal: &api.WCMatchup{
			HomeTeam: "Croatia", HomeTeamID: 10155, HomeShort: "CRO",
			AwayTeam: "Morocco", AwayTeamID: 6262, AwayShort: "MAR",
			HomeScore: wcInt(2), AwayScore: wcInt(1), WinnerID: wcInt(10155),
		},
		Champion: &arg,
		RunnerUp: &fra,
	}
}

func mockGroupTeams(a, b, c, d string) []api.LeagueTableEntry {
	codes := []string{a, b, c, d}
	pts := []int{9, 6, 3, 0}
	out := make([]api.LeagueTableEntry, 4)
	for i, code := range codes {
		out[i] = api.LeagueTableEntry{
			Position: i + 1,
			Team:     api.Team{ID: 8000 + i, Name: code, ShortName: code},
			Played: 3, Won: 3 - i, Drawn: 0, Lost: i,
			GoalsFor: 6 - i, GoalsAgainst: 2 + i, GoalDifference: 4 - 2*i, Points: pts[i],
		}
	}
	return out
}
