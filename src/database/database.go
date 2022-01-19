package database

import (
	"context"
	"fmt"
	"time"

	"github.com/kaspar-p/bee/src/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	context context.Context
	users *mongo.Collection
	busyTimes *mongo.Collection
}

var DatabaseInstance *Database;

func Connect() (*mongo.Client, context.Context, context.CancelFunc) {
	clientOptions := options.Client().ApplyURI(constants.ConnectionURL);
	context, cancel := context.WithTimeout(context.Background(), 10 * time.Hour);
	
	client, err := mongo.Connect(context, clientOptions);
	if err != nil {
		fmt.Println("Error connecting to database: ", err);
		cancel();
	}

	return client, context, cancel;
}

func InitializeDatabase() context.CancelFunc {
	client, context, cancel := Connect();

	DatabaseInstance = &Database{
		context: context,
		users: client.Database(constants.DatabaseName).Collection(constants.UsersCollectionName),
		busyTimes: client.Database(constants.DatabaseName).Collection(constants.BusyTimesCollectionName),
	};

	return cancel;
}

func ObjectIDToString(insertedID interface{}) string {
	objectID, _ := insertedID.(primitive.ObjectID);
	return objectID.String()
}

