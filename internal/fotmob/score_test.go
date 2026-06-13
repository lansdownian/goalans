package fotmob

import "testing"

func TestParseScoreStr(t *testing.T) {
	h, a := parseScoreStr("4 - 1")
	if h == nil || a == nil || *h != 4 || *a != 1 {
		t.Fatalf("got %v %v", h, a)
	}
	h, a = parseScoreStr("2-0")
	if h == nil || a == nil || *h != 2 || *a != 0 {
		t.Fatalf("got %v %v", h, a)
	}
	h, a = parseScoreStr("")
	if h != nil || a != nil {
		t.Fatalf("expected nil for empty")
	}
}
