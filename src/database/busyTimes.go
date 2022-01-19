package database

import (
	"fmt"

	"github.com/kaspar-p/bee/src/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateBusyTimesDocuments(busyTimes []*entities.BusyTime) []interface{} {
	busyTimesDocuments := make([]interface{}, len(busyTimes));
	for index, busyTime := range busyTimes {
		busyTimesDocuments[index] = busyTime.ConvertBusyTimeToDocument();
	}
	return busyTimesDocuments;
}

func (database *Database) OverwriteUserBusyTimes(user *entities.User, busyTimes []*entities.BusyTime) {
	fmt.Println("Overwriting user", user.Name, "busy times with " + fmt.Sprint(len(busyTimes)) + " busy times.")
	
	// Delete all of the busy times associated with that user
	filter := bson.D{
		{ Key: "OwnerID" , Value: user.ID },
		{ Key: "BelongsTo", Value: user.BelongsTo },
	};
	deleteResult, err := database.busyTimes.DeleteMany(database.context, filter);
	if err != nil {
		fmt.Println("Error while deleting busy times tied to user", user.Name);
		return;
	}
	fmt.Println("Found and deleted", deleteResult.DeletedCount, "events tied to user", user.Name);

	// Add the new busy times
	database.AddBusyTimes(busyTimes)
}

func (database *Database) AddBusyTimes(busyTimes []*entities.BusyTime) {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	busyTimesDocuments := CreateBusyTimesDocuments(busyTimes);

	_, err := database.busyTimes.InsertMany(database.context, busyTimesDocuments);
	if err != nil {
		fmt.Println("Error inserting busy times: ", busyTimesDocuments, ". Error: ", err)
		panic(&AddBusyTimeError{Err: err})
	}
}

func (database *Database) GetBusyTimes() []*entities.BusyTime {
	cursor, err := database.busyTimes.Find(database.context, bson.D {{ }});
	if err != nil {
		fmt.Println("Error getting cursor when finding all busyTimes objects. Error: ", err);
		panic(&GetBusyTimeError{ Err: err });
	}

	var results []bson.M
	if err = cursor.All(database.context, &results); err != nil {
		fmt.Println("Error getting results from cursor when getting all busyTimes objects. Error: ", err);
		panic(&GetBusyTimeError{ Err: err });
	}
	
	// Create BusyTime's out of the results
	busyTimesArray := make([]*entities.BusyTime, 0);
	for _, result := range results {
		guildID := result["BelongsTo"].(string);
		ownerID := result["OwnerID"].(string);
		title := result["Title"].(string);
		start := (result["Start"].(primitive.DateTime)).Time();
		end := (result["End"].(primitive.DateTime)).Time();
		
		// Create new busyTime
		newBusyTime := entities.CreateBusyTime(ownerID, guildID, title, start, end);

		// Assign new busyTime to array corresponding to ONE user
		busyTimesArray = append(busyTimesArray, &newBusyTime);
	}

	fmt.Println("Gotten", len(busyTimesArray), "total busy times (events) from database!");

	return busyTimesArray;
}