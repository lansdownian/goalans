package fotmob

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lansdownian/goalans/internal/api"
	"github.com/lansdownian/goalans/internal/data"
)

const (
	baseURL          = "https://www.fotmob.com"
	finishedDaysBack = 14
)

var nextDataRe = regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">(.+?)</script>`)

type Client struct {
	httpClient    *http.Client
	pageURLs      map[int]string
	pageURLsMu    sync.RWMutex
	maxConcurrent chan struct{}
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        20,
				MaxIdleConnsPerHost: 20,
			},
		},
		pageURLs:      make(map[int]string),
		maxConcurrent: make(chan struct{}, 6),
	}
}

func (c *Client) StorePageURL(matchID int, slug string) {
	if slug == "" {
		return
	}
	c.pageURLsMu.Lock()
	c.pageURLs[matchID] = slug
	c.pageURLsMu.Unlock()
}

func (c *Client) getPageURL(matchID int) string {
	c.pageURLsMu.RLock()
	defer c.pageURLsMu.RUnlock()
	return c.pageURLs[matchID]
}

func (c *Client) MatchesForDate(ctx context.Context, date time.Time, wantLive, wantFinished bool) ([]api.Match, error) {
	return c.MatchesInRange(ctx, date, date, wantLive, wantFinished)
}

func (c *Client) MatchesInRange(ctx context.Context, from, to time.Time, wantLive, wantFinished bool) ([]api.Match, error) {
	var mu sync.Mutex
	var all []api.Match
	var wg sync.WaitGroup

	for _, leagueID := range data.LeagueIDs() {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c.maxConcurrent <- struct{}{}
			defer func() { <-c.maxConcurrent }()

			matches, err := c.leagueMatchesInRange(ctx, id, from, to, wantLive, wantFinished)
			if err != nil || len(matches) == 0 {
				return
			}
			mu.Lock()
			all = append(all, matches...)
			mu.Unlock()
		}(leagueID)
	}
	wg.Wait()
	return all, nil
}

func (c *Client) leagueMatchesInRange(ctx context.Context, leagueID int, from, to time.Time, wantLive, wantFinished bool) ([]api.Match, error) {
	props, err := fetchLeaguePage(ctx, c.httpClient, leagueID)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Details struct {
			ID      flexInt `json:"id"`
			Name    string  `json:"name"`
			Country string  `json:"country"`
		} `json:"details"`
		Fixtures struct {
			AllMatches []rawMatch `json:"allMatches"`
		} `json:"fixtures"`
	}
	if err := json.Unmarshal(props, &resp); err != nil {
		return nil, err
	}

	out := make([]api.Match, 0, 8)
	for _, m := range resp.Fixtures.AllMatches {
		match, ok := m.toAPI(resp.Details.ID.int(), resp.Details.Name, resp.Details.Country, from, to, wantLive, wantFinished)
		if !ok {
			continue
		}
		c.StorePageURL(match.ID, match.PageURL)
		out = append(out, match)
	}
	return out, nil
}

func (c *Client) LiveMatches(ctx context.Context) ([]api.Match, error) {
	all, err := c.MatchesForDate(ctx, time.Now(), true, false)
	if err != nil {
		return nil, err
	}
	live := make([]api.Match, 0, len(all))
	for _, m := range all {
		if m.Status == api.MatchStatusLive {
			live = append(live, m)
		}
	}
	return live, nil
}

func (c *Client) FinishedMatches(ctx context.Context) ([]api.Match, error) {
	now := time.Now().UTC()
	from := now.AddDate(0, 0, -finishedDaysBack)
	all, err := c.MatchesInRange(ctx, from, now, false, true)
	if err != nil {
		return nil, err
	}
	finished := make([]api.Match, 0, len(all))
	for _, m := range all {
		if m.Status == api.MatchStatusFinished {
			finished = append(finished, m)
		}
	}
	sort.Slice(finished, func(i, j int) bool {
		if finished[i].MatchTime == nil {
			return false
		}
		if finished[j].MatchTime == nil {
			return true
		}
		return finished[i].MatchTime.After(*finished[j].MatchTime)
	})
	return finished, nil
}

func (c *Client) MatchDetails(ctx context.Context, matchID int) (*api.MatchDetails, error) {
	slug := c.getPageURL(matchID)
	if slug == "" {
		return nil, fmt.Errorf("no page URL for match %d — open from match list first", matchID)
	}
	return fetchMatchPage(ctx, c.httpClient, slug)
}

func fetchLeaguePage(ctx context.Context, client *http.Client, leagueID int) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/leagues/%d", baseURL, leagueID)
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
		return nil, fmt.Errorf("league page %d: status %d", leagueID, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return extractPageProps(string(body))
}

func fetchMatchPage(ctx context.Context, client *http.Client, slug string) (*api.MatchDetails, error) {
	url := baseURL + slug
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
		return nil, fmt.Errorf("match page: status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	props, err := extractPageProps(string(body))
	if err != nil {
		return nil, err
	}

	var raw rawMatchDetails
	if err := json.Unmarshal(props, &raw); err != nil {
		return nil, err
	}
	return raw.toAPI(), nil
}

func extractPageProps(html string) (json.RawMessage, error) {
	m := nextDataRe.FindStringSubmatch(html)
	if len(m) < 2 {
		return nil, fmt.Errorf("__NEXT_DATA__ not found")
	}
	var envelope struct {
		Props struct {
			PageProps json.RawMessage `json:"pageProps"`
		} `json:"props"`
	}
	if err := json.Unmarshal([]byte(m[1]), &envelope); err != nil {
		return nil, err
	}
	return envelope.Props.PageProps, nil
}

func setBrowserHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
}

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Parse("2006-01-02T15:04:05.000Z", s)
	}
	return t, nil
}

func scorePtr(v int) *int { return &v }

type rawMatch struct {
	ID     flexInt `json:"id"`
	Home   rawTeam `json:"home"`
	Away   rawTeam `json:"away"`
	Status struct {
		UTCTime   string `json:"utcTime"`
		Started   *bool  `json:"started"`
		Finished  *bool  `json:"finished"`
		Cancelled *bool  `json:"cancelled"`
		ScoreStr  string `json:"scoreStr"`
		LiveTime  struct {
			Short string `json:"short"`
		} `json:"liveTime"`
	} `json:"status"`
	PageURL string `json:"pageUrl"`
}

type rawTeam struct {
	ID        flexInt `json:"id"`
	Name      string  `json:"name"`
	ShortName string  `json:"shortName"`
	Score     *int    `json:"score"`
}

func utcDate(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func (m rawMatch) toAPI(leagueID int, leagueName, country string, from, to time.Time, wantLive, wantFinished bool) (api.Match, bool) {
	if m.Status.UTCTime == "" {
		return api.Match{}, false
	}
	if m.Status.Cancelled != nil && *m.Status.Cancelled {
		return api.Match{}, false
	}

	t, err := parseTime(m.Status.UTCTime)
	if err != nil {
		return api.Match{}, false
	}
	matchDay := utcDate(t)
	if matchDay.Before(utcDate(from)) || matchDay.After(utcDate(to)) {
		return api.Match{}, false
	}

	finished := m.Status.Finished != nil && *m.Status.Finished
	started := m.Status.Started != nil && *m.Status.Started

	var status api.MatchStatus
	switch {
	case finished:
		status = api.MatchStatusFinished
	case started:
		status = api.MatchStatusLive
	default:
		status = api.MatchStatusNotStarted
	}

	if status == api.MatchStatusLive && !wantLive {
		return api.Match{}, false
	}
	if status == api.MatchStatusFinished && !wantFinished {
		return api.Match{}, false
	}
	if status == api.MatchStatusNotStarted {
		return api.Match{}, false
	}

	match := api.Match{
		ID: m.ID.int(),
		League: api.League{ID: leagueID, Name: leagueName, Country: country},
		HomeTeam: api.Team{ID: m.Home.ID.int(), Name: m.Home.Name, ShortName: fallbackShort(m.Home.ShortName, m.Home.Name)},
		AwayTeam: api.Team{ID: m.Away.ID.int(), Name: m.Away.Name, ShortName: fallbackShort(m.Away.ShortName, m.Away.Name)},
		Status: status,
		HomeScore: m.Home.Score,
		AwayScore: m.Away.Score,
		MatchTime: &t,
		PageURL: m.PageURL,
	}
	if match.HomeScore == nil && match.AwayScore == nil {
		match.HomeScore, match.AwayScore = parseScoreStr(m.Status.ScoreStr)
	}
	if lt := strings.TrimSpace(m.Status.LiveTime.Short); lt != "" {
		match.LiveTime = &lt
	}
	return match, true
}

func fallbackShort(short, full string) string {
	if short != "" {
		return short
	}
	return full
}

func parseScoreStr(s string) (home, away *int) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return nil, nil
	}
	var h, a int
	if _, err := fmt.Sscanf(s, "%d - %d", &h, &a); err != nil {
		if _, err := fmt.Sscanf(s, "%d-%d", &h, &a); err != nil {
			return nil, nil
		}
	}
	return &h, &a
}

type rawMatchDetails struct {
	General struct {
		MatchID string `json:"matchId"`
		HomeTeam rawTeam `json:"homeTeam"`
		AwayTeam rawTeam `json:"awayTeam"`
		LeagueName string `json:"leagueName"`
		Started bool `json:"started"`
		Finished bool `json:"finished"`
		MatchTimeUTC     string `json:"matchTimeUTC"`
		MatchTimeUTCDate string `json:"matchTimeUTCDate"`
	} `json:"general"`
	Header struct {
		Teams []struct {
			ID int `json:"id"`
			Name string `json:"name"`
			Score *int `json:"score"`
		} `json:"teams"`
		Status struct {
			ScoreStr string `json:"scoreStr"`
			LiveTime struct {
				Short string `json:"short"`
			} `json:"liveTime"`
		} `json:"status"`
	} `json:"header"`
	Content struct {
		MatchFacts struct {
			InfoBox struct {
				Stadium struct{ Name string } `json:"Stadium"`
				Referee struct{ Text string } `json:"Referee"`
			} `json:"infoBox"`
		} `json:"matchFacts"`
		StatsPeriods []struct {
			Stats []struct {
				Title string `json:"title"`
				Stats []struct {
					Title string `json:"title"`
					Stats []struct {
						Title string `json:"title"`
						Stats []struct {
							Title string `json:"title"`
							Type  string `json:"type"`
							Stat  []struct {
								Value string `json:"value"`
							} `json:"stat"`
						} `json:"stats"`
					} `json:"stats"`
				} `json:"stats"`
			} `json:"stats"`
		} `json:"statsPeriods"`
	} `json:"content"`
	Events struct {
		Events []struct {
			Type string `json:"type"`
			Time int    `json:"time"`
			TeamID int  `json:"teamId"`
			Player struct {
				Name string `json:"name"`
			} `json:"player"`
			AssistStr string `json:"assistStr"`
			Card      string `json:"card"`
		} `json:"events"`
	} `json:"events"`
}

func (r rawMatchDetails) toAPI() *api.MatchDetails {
	g := r.General
	var matchTime *time.Time
	if t, err := parseTime(g.MatchTimeUTCDate); err == nil {
		matchTime = &t
	} else if t, err := parseTime(g.MatchTimeUTC); err == nil {
		matchTime = &t
	}

	status := api.MatchStatusNotStarted
	switch {
	case g.Finished:
		status = api.MatchStatusFinished
	case g.Started:
		status = api.MatchStatusLive
	}

	home := api.Team{ID: g.HomeTeam.ID.int(), Name: g.HomeTeam.Name, ShortName: fallbackShort(g.HomeTeam.ShortName, g.HomeTeam.Name)}
	away := api.Team{ID: g.AwayTeam.ID.int(), Name: g.AwayTeam.Name, ShortName: fallbackShort(g.AwayTeam.ShortName, g.AwayTeam.Name)}

	homeScore := g.HomeTeam.Score
	awayScore := g.AwayTeam.Score
	if len(r.Header.Teams) >= 2 {
		homeScore = r.Header.Teams[0].Score
		awayScore = r.Header.Teams[1].Score
	}
	if homeScore == nil && awayScore == nil {
		homeScore, awayScore = parseScoreStr(r.Header.Status.ScoreStr)
	}

	matchID := 0
	if id, err := strconv.Atoi(g.MatchID); err == nil {
		matchID = id
	}

	d := &api.MatchDetails{
		Match: api.Match{
			ID: matchID,
			League: api.League{Name: g.LeagueName},
			HomeTeam: home,
			AwayTeam: away,
			Status: status,
			HomeScore: homeScore,
			AwayScore: awayScore,
			MatchTime: matchTime,
		},
		Venue: r.Content.MatchFacts.InfoBox.Stadium.Name,
		Referee: strings.TrimSpace(r.Content.MatchFacts.InfoBox.Referee.Text),
	}

	if lt := strings.TrimSpace(r.Header.Status.LiveTime.Short); lt != "" {
		d.LiveTime = &lt
	}

	for _, e := range r.Events.Events {
		ev := api.MatchEvent{
			Minute: e.Time,
			Type: strings.ToLower(e.Type),
			Player: e.Player.Name,
			Assist: strings.TrimSpace(e.AssistStr),
		}
		if e.TeamID == home.ID {
			ev.Team = home
		} else {
			ev.Team = away
		}
		if ev.Type == "card" && e.Card != "" {
			ev.Type = "card"
		}
		d.Events = append(d.Events, ev)
	}

	// Flatten first stats period (full match)
	if len(r.Content.StatsPeriods) > 0 {
		for _, group := range r.Content.StatsPeriods[0].Stats {
			for _, row := range flattenStatRows(group.Stats) {
				if len(row.Stat) >= 2 {
					d.Statistics = append(d.Statistics, api.Stat{
						Label: row.Title,
						HomeValue: row.Stat[0].Value,
						AwayValue: row.Stat[1].Value,
					})
				}
			}
		}
	}

	return d
}

type flatStat struct {
	Title string
	Stat  []struct{ Value string `json:"value"` }
}

func flattenStatRows(rows []struct {
	Title string `json:"title"`
	Stats []struct {
		Title string `json:"title"`
		Stats []struct {
			Title string `json:"title"`
			Type  string `json:"type"`
			Stat  []struct {
				Value string `json:"value"`
			} `json:"stat"`
		} `json:"stats"`
	} `json:"stats"`
}) []flatStat {
	var out []flatStat
	for _, r := range rows {
		for _, s := range r.Stats {
			for _, leaf := range s.Stats {
				out = append(out, flatStat{Title: leaf.Title, Stat: leaf.Stat})
			}
		}
	}
	return out
}
