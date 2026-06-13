package api

// LeagueTableEntry represents a team's standing in a table.
type LeagueTableEntry struct {
	Position       int  `json:"position"`
	Team           Team `json:"team"`
	Played         int  `json:"played"`
	Won            int  `json:"won"`
	Drawn          int  `json:"drawn"`
	Lost           int  `json:"lost"`
	GoalsFor       int  `json:"goals_for"`
	GoalsAgainst   int  `json:"goals_against"`
	GoalDifference int  `json:"goal_difference"`
	Points         int  `json:"points"`
}
