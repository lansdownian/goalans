package fotmob

import (
	"encoding/json"
	"testing"
)

func TestFlexIntUnmarshal(t *testing.T) {
	var v struct {
		ID flexInt `json:"id"`
	}
	if err := json.Unmarshal([]byte(`{"id":42}`), &v); err != nil || v.ID.int() != 42 {
		t.Fatalf("number: got %d err=%v", v.ID.int(), err)
	}
	if err := json.Unmarshal([]byte(`{"id":"4667751"}`), &v); err != nil || v.ID.int() != 4667751 {
		t.Fatalf("string: got %d err=%v", v.ID.int(), err)
	}
}
