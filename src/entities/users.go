package entities

import (
	"go.mongodb.org/mongo-driver/bson"
)

type CurrentlyBusy struct {
	IsBusy   bool
	BusyWith string
}

type User struct {
	Id        string `bson:"Id"`
	Name      string `bson:"Name"`
	BelongsTo string `bson:"BelongsTo"`
}

func CreateUser(userName, id, guildId string) *User {
	user := User{
		Name:      userName,
		Id:        id,
		BelongsTo: guildId,
	}

	return &user
}

func (user *User) ConvertUserToDocument() bson.D {
	return bson.D{
		{Key: "Id", Value: user.Id},
		{Key: "Name", Value: user.Name},
		{Key: "BelongsTo", Value: user.BelongsTo},
	}
}
