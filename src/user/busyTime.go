package user

import "time"

type BusyTime struct {
	CourseCode string
	Start time.Time
	End time.Time
}