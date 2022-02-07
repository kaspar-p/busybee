package entities

import (
	"time"

	"github.com/apognu/gocal"
	"github.com/kaspar-p/busybee/src/utils"
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

func EventsToBusyTimes(guildId, userId string, events []gocal.Event) []*BusyTime {
	busyTimes := make([]*BusyTime, 0)

	for i := 0; i < len(events); i++ {
		event := events[i]

		title := utils.ParseEventTitle(event.Summary)
		busyTime := CreateBusyTime(userId, guildId, title, *event.Start, *event.End)
		busyTimes = append(busyTimes, &busyTime)
	}

	return busyTimes
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
