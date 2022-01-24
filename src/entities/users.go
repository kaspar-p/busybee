package entities

import (
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var Users map[string]map[string] *User

type CurrentlyBusy struct {
	IsBusy bool
	BusyWith string
}

type User struct {
	ID string
	Name string
	CurrentlyBusy CurrentlyBusy
	BusyTimes []*BusyTime
	BelongsTo string
}

func InitializeUsers(guildIds []string) {
	Users = make(map[string]map[string] *User);

	for _, guildId := range guildIds {
		Users[guildId] = make(map[string]*User);
	}
}

func CreateUser(userName string, ID string, guildID string) *User {
	user := User{
		Name: userName,
		ID: ID,
		CurrentlyBusy: CurrentlyBusy {
			IsBusy: false,
			BusyWith: "",
		},
		BelongsTo: guildID,
		BusyTimes: []*BusyTime{},
	}

	return &user;
}

func (user *User) IsBusy(t time.Time) bool {
	for _, busyTime := range user.BusyTimes {
		if (busyTime.Start.Before(t) && busyTime.End.After(t)) {
			return true;
		}
	}

	return false;
}

func (user *User) GetTodaysEvents() []*BusyTime {
	year, month, day := time.Now().Date();
	beginningOfDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local);
	endOfDay := time.Date(year, month, day, 23, 59, 59, 0, time.Local);

	todaysEvents := make([]*BusyTime, 0);
	for _, busyTime := range user.BusyTimes {
		if 	busyTime.Start.After(beginningOfDay) &&
			busyTime.End.Before(endOfDay) &&
			busyTime.End.After(time.Now()) {
				fmt.Println("Adding event", busyTime.Title, "starting at:", busyTime.Start, "and ending at:", busyTime.End);
				todaysEvents = append(todaysEvents, busyTime);
		}
	}

	// Sort them by start time
	sort.Slice(todaysEvents, func(i, j int) bool {
		return todaysEvents[i].Start.Unix() < todaysEvents[j].Start.Unix()
	})

	fmt.Println("Today's events for user", user.Name, "are:", todaysEvents);

	return todaysEvents;
}

func (user *User) GetLatestEndTime() time.Time {
	if len(user.BusyTimes) == 0 { 
		panic("No latest end time found for user: " + user.Name);
	}

	latest := user.BusyTimes[0].End;
	for _, busyTime := range user.BusyTimes {
		if busyTime.End.After(latest) {
			latest = busyTime.End;
		}
	}

	return latest;
}

func (user *User) GetLatestStartTime() time.Time {
	if len(user.BusyTimes) == 0 {
		panic("No latest start time found for user: " + user.Name);
	}

	latest := user.BusyTimes[0].Start;
	for _, busyTime := range user.BusyTimes {
		if busyTime.Start.After(latest) {
			latest = busyTime.Start;
		}
	}

	return latest;
}

func (user *User) IsBusyBetween(t1 time.Time, t2 time.Time) bool {
	for _, busyTime := range user.BusyTimes {
		if busyTime.Start.After(t1) && busyTime.Start.Before(t2) {
			return true;
		}
	}

	return false;
}

func (user *User) SortBusyTimes() {
	// Sort them by their start times
	sort.Slice(user.BusyTimes, func(i, j int) bool {
		return user.BusyTimes[i].Start.Before(user.BusyTimes[j].Start);
	})
}

func (user *User) ConvertUserToDocument() bson.D {
	return bson.D { 
		{ Key: "ID", Value: user.ID },
		{ Key: "Name", Value: user.Name },
		{ Key: "BelongsTo", Value: user.BelongsTo },
	}
}