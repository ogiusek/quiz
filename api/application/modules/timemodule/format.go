package timemodule

import "time"

const (
	// DateFormat string = "2006-01-02 15:04:05"
	DateFormat string = "2006-01-02T15:04:05Z07:00"
)

func FormatDate(date time.Time) string {
	// return date.Format(DateFormat)
	return date.Format(time.RFC3339Nano)
}
