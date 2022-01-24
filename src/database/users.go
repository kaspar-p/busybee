package database

import (
	"fmt"

	"github.com/kaspar-p/bee/src/entities"
	"go.mongodb.org/mongo-driver/bson"
)

func (database *Database) RemoveAllUsersInGuild(guildId string) error {
	if database == nil {
		return &DatabaseUninitializedError{};
	}

	filter := bson.D {{ Key: "BelongsTo", Value: guildId }};
	deleteResult, err := database.users.DeleteMany(database.context, filter);
	fmt.Println("Deleted", deleteResult.DeletedCount, "users that belonged to guild", guildId);

	return err;
}

func (database *Database) AddUser(newUser *entities.User) string {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	userDocument := newUser.ConvertUserToDocument()

	result, err := database.users.InsertOne(database.context, userDocument);
	if err != nil {
		fmt.Println("Error inserting user: ", newUser, ". Error: ", err);
		panic(&AddUserError{ Err: err })
	}

	id := ObjectIDToString(result.InsertedID);
	
	return id;
}

func (database *Database) GetUsers() []*entities.User {
	cursor, err := database.users.Find(database.context, bson.D{{ }} );
	if err != nil {
		fmt.Println("Error getting cursor when finding all users. Error: ", err);
		panic(&GetUserError{ Err: err });
	}

	var results []bson.M
	if err = cursor.All(database.context, &results); err != nil {
		fmt.Println("Error getting results from cursor when getting all users. Error: ", err);
		panic(&GetUserError{ Err: err });
	}
	
	// Create users out of the results
	users := make([]*entities.User, 0);
	for _, result := range results {
		name := result["Name"].(string)
		ID := result["ID"].(string)
		guildID := result["BelongsTo"].(string)
		user := entities.CreateUser(name, ID, guildID);

		users = append(users, user);
	}

	return users;
}

