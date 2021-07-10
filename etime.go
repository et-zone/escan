package escan

import (
	"time"
)

func NewNowTimeFormat() string {
	return time.Now().Format(TIME_FORMAT)
}

func TimeStringToTime(s string) (time.Time, error) {
	return time.Parse(TIME_FORMAT, s)
}

func TimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}

	return t.Format(TIME_FORMAT)
}

func TimeToStringFormat(t *time.Time, format string) string {
	if t == nil {
		return ""
	}
	return t.Format(format)
}
