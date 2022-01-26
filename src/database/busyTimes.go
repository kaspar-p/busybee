package database

import (
	"fmt"
	"log"

	"github.com/kaspar-p/bee/src/entities"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateBusyTimesDocuments(busyTimes []*entities.BusyTime) []interface{} {
	busyTimesDocuments := make([]interface{}, len(busyTimes))
	for index, busyTime := range busyTimes {
		busyTimesDocuments[index] = busyTime.ConvertBusyTimeToDocument()
	}

	return busyTimesDocuments
}

func CreateBusyTimeFromResult(result bson.M) entities.BusyTime {
	guildId, found := result["BelongsTo"].(string)
	if !found {
		log.Panic("No Guild ID found in GetBusyTimes()")
	}

	ownerId, found := result["OwnerId"].(string)
	if !found {
		log.Panic("No owner ID found in GetBusyTimes()")
	}

	title, found := result["Title"].(string)
	if !found {
		log.Panic("No title found in GetBusyTimes()")
	}

	start, found := result["Start"]
	if !found {
		log.Panic("No start time found in GetBusyTimes()")
		panic(&GetBusyTimeError{})
	}

	end, found := result["End"]
	if !found {
		log.Panic("No start time found in GetBusyTimes()")
		panic(&GetBusyTimeError{})
	}

	// Create new busyTime
	newBusyTime := entities.CreateBusyTime(ownerId, guildId, title,
		start.(primitive.DateTime).Time(), end.(primitive.DateTime).Time())

	return newBusyTime
}

func (database *Database) OverwriteUserBusyTimes(user *entities.User, busyTimes []*entities.BusyTime) {
	fmt.Println("Overwriting user", user.Name, "busy times with "+fmt.Sprint(len(busyTimes))+" busy times.")

	// Delete all of the busy times associated with that user
	filter := bson.D{
		{Key: "OwnerId", Value: user.Id},
		{Key: "BelongsTo", Value: user.BelongsTo},
	}

	deleteResult, err := database.busyTimes.DeleteMany(database.context, filter)
	if err != nil {
		log.Panic("Error while deleting busy times tied to user", user.Name)

		return
	}

	fmt.Println("Found and deleted", deleteResult.DeletedCount, "events tied to user", user.Name)

	// Add the new busy times
	database.AddBusyTimes(busyTimes)
}

func (database *Database) AddBusyTimes(busyTimes []*entities.BusyTime) {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	busyTimesDocuments := CreateBusyTimesDocuments(busyTimes)

	_, err := database.busyTimes.InsertMany(database.context, busyTimesDocuments)
	if err != nil {
		log.Panic("Error inserting busy times: ", busyTimesDocuments, ". Error: ", err)
		panic(&AddBusyTimeError{Err: err})
	}
}

func (database *Database) RemoveAllBusyTimesInGuild(guildId string) error {
	if database == nil {
		return &DatabaseUninitializedError{}
	}

	filter := bson.D{{Key: "BelongsTo", Value: guildId}}
	deleteResult, err := database.busyTimes.DeleteMany(database.context, filter)
	fmt.Println("Deleted", deleteResult.DeletedCount, "users that belonged to guild", guildId)

	return errors.Wrap(err, "Error deleting all busy times from a guild")
}

func (database *Database) GetBusyTimes() []*entities.BusyTime {
	cursor, err := database.busyTimes.Find(database.context, bson.D{{}})
	if err != nil {
		log.Panic("Error getting cursor when finding all busyTimes objects. Error: ", err)
		panic(&GetBusyTimeError{Err: err})
	}

	var results []bson.M
	if err = cursor.All(database.context, &results); err != nil {
		log.Panic("Error getting results from cursor when getting all busyTimes objects. Error: ", err)
		panic(&GetBusyTimeError{Err: err})
	}

	// Create BusyTime's out of the results
	busyTimesArray := make([]*entities.BusyTime, 0)

	for _, result := range results {
		newBusyTime := CreateBusyTimeFromResult(result)
		busyTimesArray = append(busyTimesArray, &newBusyTime)
	}

	return busyTimesArray
}
