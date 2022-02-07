package ingest

import (
	"log"

	"github.com/apognu/gocal"
	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/busybee/src/entities"
	"github.com/kaspar-p/busybee/src/persist"
)

func GetOrCreateUser(database *persist.DatabaseType, userId, userName, guildId string) *entities.User {
	if user, userExists := database.GetUser(guildId, userId); userExists {
		log.Println("User found with ID: ", userId)

		return user
	} else {
		log.Println("User created with ID:", userId)
		// Create the new user
		user := entities.CreateUser(userName, userId, guildId)

		// Add the new user to the database
		database.AddUser(user)

		return user
	}
}

func IngestNewData(database *persist.DatabaseType, message *discordgo.MessageCreate, events []gocal.Event) {
	// Create a user if they do not already exist - overwrites BusyTimes
	user := GetOrCreateUser(database, message.Author.ID, message.Author.Username, message.GuildID)
	busyTimes := entities.EventsToBusyTimes(message.GuildID, user.Id, events)

	database.OverwriteUserBusyTimes(user, busyTimes)
}
