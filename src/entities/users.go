package entities

import (
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

func InitializeUsers() {
	Users = make(map[string]map[string] *User);
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


func (user *User) ConvertUserToDocument() bson.D {
	return bson.D { 
		{ Key: "ID", Value: user.ID },
		{ Key: "Name", Value: user.Name },
		{ Key: "BelongsTo", Value: user.BelongsTo },
	}
}