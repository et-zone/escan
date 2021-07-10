package escan

import (
	"errors"
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

func TimeStringToUnix(s string) (int64, error) {
	t, err := time.Parse(TIME_FORMAT, s)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

//return like  2006-01-02
func DateTimeToDataString(s string) (string, error) {
	if s == "" {
		return "", errors.New("s not datetime")
	}
	return s[:10], nil
}

//return like  15:04:05
func DateTimeToTimeString(s string) (string, error) {
	if s == "" {
		return "", errors.New("s not datetime")
	}
	return s[11:], nil
}

/*
	Sunday Weekday = 0
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
*/
func DateTimeToWeek(s string) (time.Weekday, error) {
	t, err := time.Parse(TIME_FORMAT, s)
	if err != nil {
		return 0, err
	}
	return t.Weekday(), nil
}
