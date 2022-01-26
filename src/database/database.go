package database

import (
	"context"
	"log"
	"time"

	"github.com/kaspar-p/bee/src/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	context   context.Context
	users     *mongo.Collection
	busyTimes *mongo.Collection
	guilds    *mongo.Collection
}

var DatabaseInstance *Database

func Connect() (*mongo.Client, context.Context, context.CancelFunc) {
	hoursKeepAlive := 10

	clientOptions := options.Client().ApplyURI(constants.ConnectionURL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hoursKeepAlive)*time.Hour)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic("Error connecting to database: ", err)
		cancel()
	}

	return client, ctx, cancel
}

func InitializeDatabase() context.CancelFunc {
	client, ctx, cancel := Connect()

	DatabaseInstance = &Database{
		context:   ctx,
		users:     client.Database(constants.DatabaseName).Collection(constants.UsersCollectionName),
		busyTimes: client.Database(constants.DatabaseName).Collection(constants.BusyTimesCollectionName),
		guilds:    client.Database(constants.DatabaseName).Collection(constants.GuildsCollectionName),
	}

	return cancel
}

func ObjectIdToString(insertedId interface{}) string {
	objectId, _ := insertedId.(primitive.ObjectID)

	return objectId.String()
}
