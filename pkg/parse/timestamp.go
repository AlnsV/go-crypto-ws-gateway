package parse

import (
	"strings"
	"time"
)

func ParseTimestamp(timestamp string, tz *time.Location) (time.Time, error) {
	date, err := time.ParseInLocation(time.RFC3339Nano, timestamp, tz)
	if err != nil {
		date, errTwo := time.ParseInLocation(
			time.RFC3339Nano,
			strings.Replace(timestamp, " ", "T", 1),
			tz,
		)
		if errTwo != nil {
			dateStr := strings.Replace(timestamp, " ", "T", 1)
			date, errThree := time.ParseInLocation(
				time.RFC3339Nano,
				dateStr+"+00:00",
				tz,
			)
			if errThree != nil {
				return time.Now().UTC(), errThree
			}
			return date, nil
		}
		return date, nil
	}
	return date, nil
}
