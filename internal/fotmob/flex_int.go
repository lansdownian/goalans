package fotmob

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// flexInt unmarshals JSON numbers or numeric strings (FotMob uses both).
type flexInt int

func (f *flexInt) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		*f = 0
		return nil
	}
	var n int
	if err := json.Unmarshal(b, &n); err == nil {
		*f = flexInt(n)
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*f = 0
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("flexInt: %q", s)
	}
	*f = flexInt(n)
	return nil
}

func (f flexInt) int() int { return int(f) }
