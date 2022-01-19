package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type BusyTime struct {
	OwnerID string
	CourseCode string
	Start time.Time
	End time.Time
}

func CreateBusyTime(ownerID string, CourseCode string, start time.Time, end time.Time) BusyTime {
	return BusyTime{
		OwnerID: ownerID,
		CourseCode: CourseCode,
		Start: start,
		End: end,
	}
}

func (busyTime *BusyTime) ConvertBusyTimeToBsonD() bson.D {
	return bson.D {
		{ Key: "OwnerID", Value: busyTime.OwnerID },
		{ Key: "CourseCode", Value: busyTime.CourseCode },
		{ Key: "Start", Value: busyTime.Start },
		{ Key: "End", Value: busyTime.End },
	}
}