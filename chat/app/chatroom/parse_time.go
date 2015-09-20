package chatroom

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	errInvalidIntervalFormat = errors.New("invalid format")
	reInterval               *regexp.Regexp

	intervals = map[string]time.Duration{
		"s":       time.Second,
		"sec":     time.Second,
		"seconds": time.Second,
		"m":       time.Minute,
		"min":     time.Minute,
		"minutes": time.Minute,
		"h":       time.Hour,
		"hour":    time.Hour,
		"hours":   time.Hour,
		"d":       time.Hour * 24,
		"day":     time.Hour * 24,
		"days":    time.Hour * 24,
		"w":       time.Hour * 24 * 7,
		"week":    time.Hour * 24 * 7,
		"weeks":   time.Hour * 24 * 7,
		"mon":     time.Hour * 24 * 30,
		"month":   time.Hour * 24 * 30,
		"months":  time.Hour * 24 * 30,
		"y":       time.Hour * 24 * 365,
		"year":    time.Hour * 24 * 365,
		"years":   time.Hour * 24 * 365,
	}
)

// parseInterval parses an interval string and returns the corresponding duration.
func parseInterval(s string) (time.Duration, error) {
	if m := reInterval.FindStringSubmatch(strings.TrimSpace(s)); len(m) != 0 {
		amount, err := strconv.ParseInt(m[1], 10, 32)

		if err != nil {
			return 0, err
		}

		interval := intervals[strings.ToLower(m[2])]
		return interval * time.Duration(amount), nil
	}

	return 0, errInvalidIntervalFormat
}
