package database

import (
	"context"
	"log"

	"github.com/kaspar-p/bee/src/entities"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (database *Database) RemoveAllUsersInGuild(guildId string) error {
	if database == nil {
		return &DatabaseUninitializedError{}
	}

	filter := bson.D{{Key: "BelongsTo", Value: guildId}}
	deleteResult, err := database.users.DeleteMany(context.Background(), filter)
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

func (database *Database) AddUser(newUser *entities.User) string {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	userDocument := newUser.ConvertUserToDocument()

	result, err := database.users.InsertOne(context.Background(), userDocument)
	if err != nil {
		log.Panic("Error inserting user: ", newUser, ". Error: ", err)
		panic(&AddUserError{Err: err})
	}

	id := ObjectIdToString(result.InsertedID)

	return id
}

func (database *Database) GetUsers() []*entities.User {
	cursor, err := database.users.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Panic("Error getting cursor when finding all users. Error: ", err)
		panic(&GetUserError{Err: err})
	}

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		log.Panic("Error getting results from cursor when getting all users. Error: ", err)
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
