package fotmob_test

import (
	"context"
	"testing"

	"github.com/lansdownian/goalans/internal/fotmob"
)

func TestFinishedMatchesReturnsRecentGames(t *testing.T) {
	if testing.Short() {
		t.Skip("network")
	}
	c := fotmob.NewClient()
	matches, err := c.FinishedMatches(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatal("expected recent finished matches (World Cup fixtures should be present)")
	}
	if matches[0].HomeScore == nil || matches[0].AwayScore == nil {
		t.Fatalf("expected scores in list, got home=%v away=%v", matches[0].HomeScore, matches[0].AwayScore)
	}
	d, err := c.MatchDetails(context.Background(), matches[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if d.HomeScore == nil || d.AwayScore == nil {
		t.Fatalf("expected scores in match details, got home=%v away=%v", d.HomeScore, d.AwayScore)
	}
}
