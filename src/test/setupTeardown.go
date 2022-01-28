package test

import (
	"github.com/kaspar-p/bee/src/database"
	discordLib "github.com/kaspar-p/bee/src/discord"
	"github.com/kaspar-p/bee/src/environment"
)

type (
	SetupFunction    = func()
	TeardownFunction = func()
)

func Teardown(funcs ...func()) func() {
	return func() {
		for _, function := range funcs {
			function()
		}
	}
}

func SetupIntegrationTest() TeardownFunction {
	config := environment.InitializeViper(environment.TESTING)

	db, disconnect := database.InitializeDatabase(config.DatabaseConfig)
	discordLib.EstablishDiscordConnection(db, config.DiscordConfig)

	return Teardown(disconnect)
}

func SetupDatabaseRequiredTests() (db *database.Database, td TeardownFunction) {
	config := environment.InitializeViper(environment.TESTING)

	db, disconnect := database.InitializeDatabase(config.DatabaseConfig)

	// Potentially create new collections for these tests

	return db, Teardown(disconnect)
}
