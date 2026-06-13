package fotmob

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/lansdownian/goalans/internal/api"
)

type wcPageResponse struct {
	Table []struct {
		Data struct {
			Tables []struct {
				LeagueID   int    `json:"leagueId"`
				LeagueName string `json:"leagueName"`
				Table      struct {
					All []wcTableRow `json:"all"`
				} `json:"table"`
			} `json:"tables"`
		} `json:"data"`
	} `json:"table"`
	Overview struct {
		Playoff wcPlayoff `json:"playoff"`
		Season  string    `json:"season"`
	} `json:"overview"`
}

type wcPlayoff struct {
	Rounds  []wcPlayoffRound `json:"rounds"`
	Special []wcPlayoffRound `json:"special"`
}

type wcPlayoffRound struct {
	Stage    string         `json:"stage"`
	Matchups []wcMatchupRaw `json:"matchups"`
}

type wcMatchupRaw struct {
	HomeTeam          string `json:"homeTeam"`
	HomeTeamID        int    `json:"homeTeamId"`
	HomeTeamShortName string `json:"homeTeamShortName"`
	AwayTeam          string `json:"awayTeam"`
	AwayTeamID        int    `json:"awayTeamId"`
	AwayTeamShortName string `json:"awayTeamShortName"`
	HomeScore         int    `json:"homeScore"`
	AwayScore         int    `json:"awayScore"`
	Winner            int    `json:"winner"`
	TBDTeam1          bool   `json:"tbdTeam1"`
	TBDTeam2          bool   `json:"tbdTeam2"`
	Matches           []struct {
		Status struct {
			Finished *bool `json:"finished"`
		} `json:"status"`
	} `json:"matches"`
}

type wcTableRow struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	Idx         int    `json:"idx"`
	Played      int    `json:"played"`
	Wins        int    `json:"wins"`
	Draws       int    `json:"draws"`
	Losses      int    `json:"losses"`
	ScoresStr   string `json:"scoresStr"`
	GoalConDiff int    `json:"goalConDiff"`
	Pts         int    `json:"pts"`
}

func (r wcTableRow) toEntry() api.LeagueTableEntry {
	var gf, ga int
	_, _ = fmt.Sscanf(r.ScoresStr, "%d-%d", &gf, &ga)
	return api.LeagueTableEntry{
		Position:       r.Idx,
		Team:           api.Team{ID: r.ID, Name: r.Name, ShortName: r.ShortName},
		Played:         r.Played,
		Won:            r.Wins,
		Drawn:          r.Draws,
		Lost:           r.Losses,
		GoalsFor:       gf,
		GoalsAgainst:   ga,
		GoalDifference: r.GoalConDiff,
		Points:         r.Pts,
	}
}

// WorldCupData fetches FIFA World Cup groups and knockout bracket from FotMob.
func (c *Client) WorldCupData(ctx context.Context, season string) (*api.WorldCupData, error) {
	if season == "" {
		season = api.WCSeason2026
	}

	props, err := fetchWorldCupPage(ctx, c.httpClient, season)
	if err != nil {
		return nil, err
	}

	var resp wcPageResponse
	if err := json.Unmarshal(props, &resp); err != nil {
		return nil, fmt.Errorf("parse world cup data: %w", err)
	}

	groups := parseWCGroups(resp)
	rounds, bronze := parseWCBracket(resp.Overview.Playoff)

	s := season
	if resp.Overview.Season != "" {
		s = resp.Overview.Season
	}

	wc := &api.WorldCupData{
		Season:         s,
		Name:           fmt.Sprintf("FIFA World Cup %s", s),
		Groups:         groups,
		KnockoutRounds: rounds,
		BronzeFinal:    bronze,
	}
	wc.Champion, wc.RunnerUp = wc.DeriveFinalists()
	return wc, nil
}

func fetchWorldCupPage(ctx context.Context, client *http.Client, season string) (json.RawMessage, error) {
	url := "https://www.fotmob.com/leagues/77/overview/world-cup?season=" + season
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	setBrowserHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("world cup page status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return extractPageProps(string(body))
}

func parseWCGroups(resp wcPageResponse) []api.WCGroup {
	if len(resp.Table) == 0 {
		return nil
	}
	groups := make([]api.WCGroup, 0, len(resp.Table[0].Data.Tables))
	for _, t := range resp.Table[0].Data.Tables {
		if len(t.Table.All) == 0 {
			continue
		}
		letter := wcGroupLetter(t.LeagueName)
		if !isWCGroupLetter(letter) {
			continue
		}
		entries := make([]api.LeagueTableEntry, 0, len(t.Table.All))
		for _, row := range t.Table.All {
			if strings.TrimSpace(row.Name) == "" && strings.TrimSpace(row.ShortName) == "" {
				continue
			}
			entries = append(entries, row.toEntry())
		}
		groups = append(groups, api.WCGroup{
			ID: t.LeagueID, Letter: letter, Name: "Group " + letter, Teams: entries,
		})
	}
	return groups
}

func parseWCBracket(playoff wcPlayoff) ([]api.WCKnockoutRound, *api.WCMatchup) {
	order := map[string]int{"1/32": 0, "1/16": 1, "1/8": 2, "1/4": 3, "1/2": 4, "final": 5}
	labels := map[string]string{
		"1/32": "Round of 64", "1/16": "Round of 32", "1/8": "Round of 16",
		"1/4": "Quarterfinals", "1/2": "Semifinals", "final": "Final",
	}

	type indexed struct {
		order int
		round api.WCKnockoutRound
	}
	var items []indexed
	for _, r := range playoff.Rounds {
		label := labels[r.Stage]
		if label == "" {
			label = r.Stage
		}
		ord := order[r.Stage]
		if _, ok := order[r.Stage]; !ok {
			ord = 99
		}
		items = append(items, indexed{
			order: ord,
			round: api.WCKnockoutRound{Stage: r.Stage, Label: label, Matchups: convertWCMatchups(r.Matchups)},
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].order < items[j].order })

	rounds := make([]api.WCKnockoutRound, len(items))
	for i, it := range items {
		rounds[i] = it.round
	}

	var bronze *api.WCMatchup
	for _, s := range playoff.Special {
		if s.Stage == "bronze" && len(s.Matchups) > 0 {
			m := convertWCMatchup(s.Matchups[0])
			bronze = &m
			break
		}
	}
	return rounds, bronze
}

func convertWCMatchups(raw []wcMatchupRaw) []api.WCMatchup {
	out := make([]api.WCMatchup, len(raw))
	for i, r := range raw {
		out[i] = convertWCMatchup(r)
	}
	return out
}

func convertWCMatchup(r wcMatchupRaw) api.WCMatchup {
	m := api.WCMatchup{
		HomeTeam: r.HomeTeam, HomeTeamID: r.HomeTeamID, HomeShort: r.HomeTeamShortName,
		AwayTeam: r.AwayTeam, AwayTeamID: r.AwayTeamID, AwayShort: r.AwayTeamShortName,
		TBDHome: r.TBDTeam1, TBDAway: r.TBDTeam2,
	}
	if len(r.Matches) > 0 {
		if f := r.Matches[0].Status.Finished; f != nil && *f {
			m.HomeScore = wcIntPtr(r.HomeScore)
			m.AwayScore = wcIntPtr(r.AwayScore)
		}
	}
	if r.Winner != 0 {
		m.WinnerID = wcIntPtr(r.Winner)
		if m.HomeScore != nil && m.AwayScore != nil && *m.HomeScore == *m.AwayScore {
			m.IsPenalties = true
		}
	}
	return m
}

func wcGroupLetter(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return name
}

func isWCGroupLetter(s string) bool {
	return len(s) == 1 && s[0] >= 'A' && s[0] <= 'Z'
}

func wcIntPtr(v int) *int { i := v; return &i }
