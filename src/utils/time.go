package utils

import "time"

func GetNow() time.Time {
	loc, _ := time.LoadLocation("America/New_York")

	now := time.Now().In(loc)

	return now
}
