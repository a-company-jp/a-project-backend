package pkg_time

import (
	"time"
)

const baseYear = 1900

type DateOnly int64

func (d DateOnly) String() string {
	t := time.Date(baseYear, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, int(d))
	return t.Format("2006-01-02")
}

func FromString(s string) (DateOnly, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return 0, err
	}
	base := time.Date(baseYear, 1, 1, 0, 0, 0, 0, time.UTC)
	days := daysBetween(base, t)
	return DateOnly(days), nil
}

func daysBetween(a, b time.Time) int {
	return int(b.Sub(a).Hours() / 24)
}
