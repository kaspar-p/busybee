package persist

import (
	"context"
	"log"
	"time"

	"github.com/kaspar-p/busybee/src/entities"
	"github.com/kaspar-p/busybee/src/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (database *DatabaseType) getManyUsersWithFilter(filter bson.D) []*entities.User {
	cursor, err := database.users.Find(context.TODO(), filter)
	if err != nil {
		log.Panic("Error getting cursor when getting many users with filter", filter, ". Error: ", err)
		panic(&GetUserError{Err: err})
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Panic("Error getting cursor when getting many users with filter", filter, ". Error: ", err)
		panic(&GetUserError{Err: err})
	}

	// Create users out of the results
	users := make([]*entities.User, 0)

	for _, result := range results {
		user := ConvertDocumentToUser(result)
		users = append(users, user)
	}

	return users
}

func (database *DatabaseType) RemoveAllUsersInGuild(guildId string) error {
	if database == nil {
		return &DatabaseUninitializedError{}
	}

	filter := bson.D{{Key: "BelongsTo", Value: guildId}}
	deleteResult, err := database.users.DeleteMany(context.TODO(), filter)
	log.Println("Deleted", deleteResult.DeletedCount, "users that belonged to guild", guildId)

	return errors.Wrap(err, "Error removing all users from guild!")
}

func ConvertDocumentToUser(result bson.M) *entities.User {
	name, found := result["Name"].(string)
	if !found {
		log.Panic("Key 'Name' not found on result for user!")
		panic(&GetUserError{})
	}

	userId, found := result["Id"].(string)
	if !found {
		log.Panic("Key 'Id' not found on result for user!")
		panic(&GetUserError{})
	}

	guildId, found := result["BelongsTo"].(string)
	if !found {
		log.Panic("Key 'guildId' not found on result for user!")
		panic(&GetUserError{})
	}

	user := entities.CreateUser(name, userId, guildId)

	return user
}

func (database *DatabaseType) AddUser(newUser *entities.User) string {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	userDocument := newUser.ConvertUserToDocument()

	result, err := database.users.InsertOne(context.TODO(), userDocument)
	if err != nil {
		log.Panic("Error inserting user: ", newUser, ". Error: ", err)
		panic(&AddUserError{Err: err})
	}

	id := ObjectIdToString(result.InsertedID)

	return id
}

func (database *DatabaseType) GetCurrentlyBusyUsers() map[string]string {
	users := database.GetUsers()

	busyWithMap := make(map[string]string)

	for _, user := range users {
		busyTimes := database.GetBusyTimesForUser(user.Id)
		for _, busyTime := range busyTimes {
			if busyTime.Start.Before(time.Now()) && busyTime.End.After(time.Now()) {
				busyWithMap[user.Name] = busyTime.Title
			}
		}
	}

	return busyWithMap
}

func (database *DatabaseType) GetUsersInGuild(guildId string) []*entities.User {
	filter := bson.D{
		{Key: "BelongsTo", Value: guildId},
	}

	return database.getManyUsersWithFilter(filter)
}

func (database *DatabaseType) GetLatestEndTime(userId string) time.Time {
	allBusyTimes := database.GetBusyTimesForUser(userId)

	latestTime := time.Now()
	for _, busyTime := range allBusyTimes {
		if busyTime.End.After(latestTime) {
			latestTime = busyTime.End
		}
	}

	return latestTime
}

func (database *DatabaseType) UserIsCurrentlyBusy(userId string) bool {
	busyUsers := database.GetCurrentlyBusyUsers()

	keys := make([]string, 0, len(busyUsers))
	for k := range busyUsers {
		keys = append(keys, k)
	}

	found, _ := utils.StringInSlice(keys, userId)

	return found
}

func (database *DatabaseType) GetUser(guildId, userId string) (*entities.User, bool) {
	user := entities.User{}
	filter := bson.D{
		{Key: "Id", Value: userId},
		{Key: "BelongsTo", Value: guildId},
	}

	err := database.users.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, false
	} else if err != nil {
		log.Panic("Error finding or decoding a user. Error: ", err)
		panic(&GetUserError{Err: err})
	}

	return &user, true
}

func (database *DatabaseType) GetUsers() []*entities.User {
	return database.getManyUsersWithFilter(bson.D{})
}
