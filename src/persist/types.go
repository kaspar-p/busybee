package persist

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	StopFunction       = func(context.Context) error
	DisconnectFunction = func()
)

type DatabaseType struct {
	users     *mongo.Collection
	busyTimes *mongo.Collection
	guilds    *mongo.Collection
}

type DatabaseConfig struct {
	ConnectionUrl   string
	DatabaseName    string
	CollectionNames *CollectionNames
}

// The names of the collections used in the backend.
type CollectionNames struct {
	Users     string
	BusyTimes string
	Guilds    string
}
