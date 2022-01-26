package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type BusyTime struct {
	OwnerId   string
	BelongsTo string
	Title     string
	Start     time.Time
	End       time.Time
}

func CreateBusyTime(ownerID, guildId, title string, start, end time.Time) BusyTime {
	return BusyTime{
		BelongsTo: guildId,
		OwnerId:   ownerID,
		Title:     title,
		Start:     start,
		End:       end,
	}
}

func (busyTime *BusyTime) ConvertBusyTimeToDocument() bson.D {
	return bson.D{
		{Key: "BelongsTo", Value: busyTime.BelongsTo},
		{Key: "OwnerId", Value: busyTime.OwnerId},
		{Key: "Title", Value: busyTime.Title},
		{Key: "Start", Value: busyTime.Start},
		{Key: "End", Value: busyTime.End},
	}
}
