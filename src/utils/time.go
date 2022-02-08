package utils

import "time"

func GetNow() time.Time {
	return time.Now().In(GetTz())
}

func GetTz() *time.Location {
	loc, _ := time.LoadLocation("America/New_York")

	return loc
}
