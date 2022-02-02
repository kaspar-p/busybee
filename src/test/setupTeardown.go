package test

import (
	discordLib "github.com/kaspar-p/bee/src/discord"
	"github.com/kaspar-p/bee/src/environment"
	"github.com/kaspar-p/bee/src/persist"
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

	db, disconnect := persist.InitializeDatabase(config.DatabaseConfig)
	discordLib.EstablishDiscordConnection(db, config.DiscordConfig)

	return Teardown(disconnect)
}

func SetupDatabaseRequiredTests() (database *persist.DatabaseType, td TeardownFunction) {
	config := environment.InitializeViper(environment.TESTING)

	database, disconnect := persist.InitializeDatabase(config.DatabaseConfig)

	// Potentially create new collections for these tests

	return database, Teardown(disconnect)
}
