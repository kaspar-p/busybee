package test

import (
	"github.com/bwmarrin/discordgo"
	discordLib "github.com/kaspar-p/busybee/src/discord"
	"github.com/kaspar-p/busybee/src/environment"
	"github.com/kaspar-p/busybee/src/persist"
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
	_, closeDiscord := discordLib.EstablishDiscordConnection(db, config.DiscordConfig)

	return Teardown(disconnect, closeDiscord)
}

func SetupDiscordRequiredTests() (s *discordgo.Session, td TeardownFunction) {
	config := environment.InitializeViper(environment.TESTING)

	db, disconnect := persist.InitializeDatabase(config.DatabaseConfig)
	discord, closeDiscord := discordLib.EstablishDiscordConnection(db, config.DiscordConfig)

	return discord, Teardown(disconnect, closeDiscord)
}

func SetupDatabaseRequiredTests() (database *persist.DatabaseType, td TeardownFunction) {
	config := environment.InitializeViper(environment.TESTING)

	database, disconnect := persist.InitializeDatabase(config.DatabaseConfig)

	// Potentially create new collections for these tests

	return database, Teardown(disconnect)
}
