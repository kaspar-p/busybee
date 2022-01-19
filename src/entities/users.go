package entities

import "go.mongodb.org/mongo-driver/bson"

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
	}

	return &user;
}

func (user *User) ConvertUserToDocument() bson.D {
	return bson.D { 
		{ Key: "ID", Value: user.ID },
		{ Key: "Name", Value: user.Name },
		{ Key: "BelongsTo", Value: user.BelongsTo },
	}
}