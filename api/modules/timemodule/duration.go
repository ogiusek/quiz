package timemodule

import (
	"encoding/json"
	"fmt"
	"time"
)

type Duration struct {
	duration time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsedDuration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}
	d.duration = parsedDuration
	return nil
}
func (d *Duration) Duration() time.Duration { return d.duration }
