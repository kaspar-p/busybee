package persist

import (
	"context"
	"log"
	"time"

	"github.com/kaspar-p/busybee/src/entities"
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

func GetStringValueFromDocument(result bson.M, key string) string {
	resultString, found := result[key].(string)
	if !found {
		log.Printf("No string with key %s found in GetBusyTimes()!\n", key)
		panic(&GetBusyTimeError{})
	}

	return resultString
}

func GetAndConvertTimeFromDocument(result bson.M, key string) time.Time {
	var resultTime time.Time

	interfaceTime, found := result[key]
	if !found {
		log.Panic("No end time found in GetBusyTimes()")
		panic(&GetBusyTimeError{})
	}

	if interfaceTime, ok := interfaceTime.(primitive.DateTime); !ok {
		log.Printf("Interface time with name %s and value %v is not convertible to primitive.Datetime!", key, interfaceTime)
		panic(&GetBusyTimeError{})
	} else {
		resultTime = interfaceTime.Time()
	}

	return resultTime
}

func CreateBusyTimeFromResult(result bson.M) entities.BusyTime {
	guildId := GetStringValueFromDocument(result, "BelongsTo")
	ownerId := GetStringValueFromDocument(result, "OwnerId")
	title := GetStringValueFromDocument(result, "Title")
	startTime := GetAndConvertTimeFromDocument(result, "Start")
	endTime := GetAndConvertTimeFromDocument(result, "End")

	// Create new busyTime
	newBusyTime := entities.CreateBusyTime(ownerId, guildId, title, startTime, endTime)

	return newBusyTime
}

func (database *DatabaseType) OverwriteUserBusyTimes(user *entities.User, busyTimes []*entities.BusyTime) {
	log.Printf("Overwriting user %s busy times with %d busy times.\n", user.Name, len(busyTimes))

	// Delete all of the busy times associated with that user
	filter := bson.D{
		{Key: "OwnerId", Value: user.Id},
		{Key: "BelongsTo", Value: user.BelongsTo},
	}

	deleteResult, err := database.busyTimes.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Panic("Error while deleting busy times tied to user", user.Name)

		return
	}

	log.Printf("Found and deleted %d events tied to user %s", deleteResult.DeletedCount, user.Name)

	// Add the new busy times
	database.AddBusyTimes(busyTimes)
}

func (database *DatabaseType) AddBusyTimes(busyTimes []*entities.BusyTime) {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	busyTimesDocuments := CreateBusyTimesDocuments(busyTimes)

	_, err := database.busyTimes.InsertMany(context.TODO(), busyTimesDocuments)
	if err != nil {
		log.Println("Error inserting busy times: ", busyTimesDocuments, ". Error: ", err)
		panic(&AddBusyTimeError{Err: err})
	}
}

func (database *DatabaseType) RemoveAllBusyTimesInGuild(guildId string) error {
	if database == nil {
		return &DatabaseUninitializedError{}
	}

	filter := bson.D{{Key: "BelongsTo", Value: guildId}}
	deleteResult, err := database.busyTimes.DeleteMany(context.TODO(), filter)
	log.Printf("Deleted %d users that belonged to guild %s.\n", deleteResult.DeletedCount, guildId)

	return errors.Wrap(err, "Error deleting all busy times from a guild")
}

func (database *DatabaseType) GetBusyTimes() []*entities.BusyTime {
	cursor, err := database.busyTimes.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Println("Error getting cursor when finding all busyTimes objects. Error: ", err)
		panic(&GetBusyTimeError{Err: err})
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Println("Error getting results from cursor when getting all busyTimes objects. Error: ", err)
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
