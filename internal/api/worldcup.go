package api

// WCFotMobLeagueID is FotMob's league ID for the FIFA World Cup.
const WCFotMobLeagueID = 77

// WCSeason2026 is the season query param for the 2026 tournament.
const WCSeason2026 = "2026"

type WCGroup struct {
	ID     int
	Letter string
	Name   string
	Teams  []LeagueTableEntry
}

type WCMatchup struct {
	HomeTeam    string
	HomeTeamID  int
	HomeShort   string
	AwayTeam    string
	AwayTeamID  int
	AwayShort   string
	HomeScore   *int
	AwayScore   *int
	WinnerID    *int
	IsPenalties bool
	TBDHome     bool
	TBDAway     bool
}

type WCKnockoutRound struct {
	Stage    string
	Label    string
	Matchups []WCMatchup
}

type WorldCupData struct {
	Season         string
	Name           string
	Groups         []WCGroup
	KnockoutRounds []WCKnockoutRound
	BronzeFinal    *WCMatchup
	Champion       *Team
	RunnerUp       *Team
}

func (d *WorldCupData) DeriveFinalists() (*Team, *Team) {
	for _, r := range d.KnockoutRounds {
		if r.Stage != "final" || len(r.Matchups) == 0 {
			continue
		}
		mu := r.Matchups[0]
		if mu.WinnerID == nil {
			return nil, nil
		}
		home := Team{ID: mu.HomeTeamID, Name: mu.HomeTeam, ShortName: mu.HomeShort}
		away := Team{ID: mu.AwayTeamID, Name: mu.AwayTeam, ShortName: mu.AwayShort}
		if *mu.WinnerID == mu.HomeTeamID {
			return &home, &away
		}
		return &away, &home
	}
	return nil, nil
}

// BracketLineCount returns content line count for bracket scrolling.
func (d *WorldCupData) BracketLineCount() int {
	if d == nil {
		return 0
	}
	count := 0
	for _, round := range d.KnockoutRounds {
		count += 2 + len(round.Matchups) + 1
	}
	if d.BronzeFinal != nil {
		count += 4
	}
	if d.Champion != nil {
		count += 2
	}
	return count
}
