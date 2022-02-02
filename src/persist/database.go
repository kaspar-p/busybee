package persist

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(connectionUrl string) (*mongo.Client, context.Context) {
	clientOptions := options.Client().ApplyURI(connectionUrl)
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic("Error connecting to database: ", err)
	}

	return client, ctx
}

func InitializeDatabase(config *DatabaseConfig) (db *DatabaseType, disconnect DisconnectFunction) {
	client, _ := Connect(config.ConnectionUrl)

	disconnectFunction := func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Println("Error encountered while disconnecting: ", err)
			panic(errors.Wrap(err, "Error encountered while disconnecting!"))
		}
	}

	db = &DatabaseType{
		users:     client.Database(config.DatabaseName).Collection(config.CollectionNames.Users),
		busyTimes: client.Database(config.DatabaseName).Collection(config.CollectionNames.BusyTimes),
		guilds:    client.Database(config.DatabaseName).Collection(config.CollectionNames.Guilds),
	}

	return db, disconnectFunction
}

func ObjectIdToString(insertedId interface{}) string {
	objectId, _ := insertedId.(primitive.ObjectID)

	return objectId.String()
}
