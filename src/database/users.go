package database

import (
	"fmt"

	usersLib "github.com/kaspar-p/bee/src/users"
	"go.mongodb.org/mongo-driver/bson"
)

type UserDocument struct {
	UserID string
	Name string
}

func CreateUserDocumentFromUser(user *usersLib.User) UserDocument {
	userDocument := UserDocument { 
		UserID: user.UserID, 
		Name: user.Name, 
	}

	return userDocument;
}

func (database *Database) AddUser(newUser *usersLib.User) string {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	userDocument := CreateUserDocumentFromUser(newUser);

	result, err := database.users.InsertOne(database.context, userDocument);
	if err != nil {
		fmt.Println("Error inserting user: ", newUser, ". Error: ", err);
		panic(&AddUserError{ Err: err })
	}

	id := ObjectIDToString(result.InsertedID);
	
	return id;
}

func (database *Database) GetUsers() []*usersLib.User {
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
	users := make([]*usersLib.User, 0);
	for _, result := range results {
		user := usersLib.CreateUser(result["name"].(string), result["userid"].(string))

		users = append(users, user);
	}
	fmt.Println("Gotten", len(users), "users from database!");

	return users;
}

