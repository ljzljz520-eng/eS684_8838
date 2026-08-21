package clock

import "time"

type Clock struct {
	instant  time.Time
	location *time.Location
}

func New(instant time.Time) Clock { return Clock{instant: instant.UTC(), location: time.UTC} }

func Fixed() Clock { return New(time.Date(2026, 8, 21, 8, 0, 0, 0, time.UTC)) }

func (c Clock) Now() time.Time { return c.instant }

func (c Clock) NowString() string { return c.instant.Format(time.RFC3339) }

func (c Clock) Location() *time.Location { return c.location }

func (c Clock) Date() string { return c.instant.Format("2006-01-02") }

func (c Clock) AtDate(date string) time.Time {
	parsed, err := time.ParseInLocation("2006-01-02", date, c.location)
	if err != nil {
		return c.instant
	}
	return parsed
}
