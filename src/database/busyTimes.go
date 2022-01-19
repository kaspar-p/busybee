package database

import (
	"fmt"

	usersLib "github.com/kaspar-p/bee/src/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateBusyTimesDocuments(busyTimes []*usersLib.BusyTime) []interface{} {
	busyTimesDocuments := make([]interface{}, len(busyTimes));
	for index, busyTime := range busyTimes {
		busyTimesDocuments[index] = busyTime.ConvertBusyTimeToBsonD();
	}
	return busyTimesDocuments;
}

func (database *Database) OverwriteUserBusyTimes(userID string, busyTimes []*usersLib.BusyTime) {
	// Delete all of the busy times associated with that user
	filter := bson.D{{ Key: "OwnerID" , Value: userID }};
	database.busyTimes.DeleteMany(database.context, filter);

	// Add the new busy times
	database.AddBusyTimes(userID, busyTimes)
}

func (database *Database) AddBusyTimes(ownerID string, busyTimes []*usersLib.BusyTime) {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	busyTimesDocuments := CreateBusyTimesDocuments(busyTimes);

	_, err := database.busyTimes.InsertMany(database.context, busyTimesDocuments);
	if err != nil {
		fmt.Println("Error inserting busy times: ", busyTimesDocuments, ". Error: ", err)
		panic(&AddCourseError{Err: err})
	}
}

func (database *Database) GetBusyTimes() []*usersLib.BusyTime {
	cursor, err := database.busyTimes.Find(database.context, bson.D {{ }});
	if err != nil {
		fmt.Println("Error getting cursor when finding all busyTimes objects. Error: ", err);
		panic(&GetCourseError{ Err: err });
	}

	var results []bson.M
	if err = cursor.All(database.context, &results); err != nil {
		fmt.Println("Error getting results from cursor when getting all busyTimes objects. Error: ", err);
		panic(&GetCourseError{ Err: err });
	}
	
	// Create BusyTime's out of the results
	busyTimesArray := make([]*usersLib.BusyTime, 0);
	for _, result := range results {
		ownerID := result["OwnerID"].(string);
		// Filter data out of [courseCode, start, end] array
		courseCode := result["CourseCode"].(string);
		
		// Parse times
		start := (result["Start"].(primitive.DateTime)).Time();
		end := (result["End"].(primitive.DateTime)).Time();
		
		// Create new busyTime
		newBusyTime := usersLib.CreateBusyTime(ownerID, courseCode, start, end);

		// Assign new busyTime to array corresponding to ONE user
		busyTimesArray = append(busyTimesArray, &newBusyTime);
	}

	fmt.Println("Gotten", len(busyTimesArray), "total busy times (events) from database!");

	return busyTimesArray;
}