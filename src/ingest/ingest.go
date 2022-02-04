package ingest

import (
	"log"

	"github.com/apognu/gocal"
	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/busybee/src/entities"
	"github.com/kaspar-p/busybee/src/persist"
)

func IngestNewData(database *persist.DatabaseType, message *discordgo.MessageCreate, events []gocal.Event) {
	// Create a user if they do not already exist - overwrites BusyTimes
	user := GetOrCreateUser(database, message.Author.ID, message.Author.Username, message.GuildID)

	OverwriteUserEvents(database, user, events)
}

func GetOrCreateUser(database *persist.DatabaseType, userId, userName, guildId string) *entities.User {
	if user, ok := entities.Users[guildId][userId]; ok {
		log.Println("User found with ID: ", userId)

		return user
	} else {
		log.Println("User created with ID:", userId)
		// Create the new user
		user := entities.CreateUser(userName, userId, guildId)

		// Add the new user to the `users` map
		entities.Users[user.BelongsTo][user.Id] = user

		// Add the new user to the database
		database.AddUser(user)

		return user
	}
}

func OverwriteUserEvents(database *persist.DatabaseType, user *entities.User, events []gocal.Event) {
	// Overwrite the busyTimes in memory
	user.BusyTimes = make([]*entities.BusyTime, 0)

	for i := 0; i < len(events); i++ {
		event := events[i]

		title := ParseEventTitle(event.Summary)
		busyTime := entities.CreateBusyTime(user.Id, user.BelongsTo, title, *event.Start, *event.End)
		user.BusyTimes = append(user.BusyTimes, &busyTime)
	}

	user.SortBusyTimes()

	// Overwrite the busyTimes in the database
	database.OverwriteUserBusyTimes(user, user.BusyTimes)
}
