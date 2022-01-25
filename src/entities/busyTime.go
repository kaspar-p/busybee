package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type BusyTime struct {
	OwnerID string
	BelongsTo string
	Title string
	Start time.Time
	End time.Time
}

func CreateBusyTime(ownerID string, guildId string, title string, start time.Time, end time.Time) BusyTime {
	return BusyTime{
		BelongsTo: guildId,
		OwnerID: ownerID,
		Title: title,
		Start: start,
		End: end,
	}
}

func (busyTime *BusyTime) ConvertBusyTimeToDocument() bson.D {
	return bson.D {
		{ Key: "BelongsTo", Value: busyTime.BelongsTo },
		{ Key: "OwnerID", Value: busyTime.OwnerID },
		{ Key: "Title", Value: busyTime.Title },
		{ Key: "Start", Value: busyTime.Start },
		{ Key: "End", Value: busyTime.End },
	}
}