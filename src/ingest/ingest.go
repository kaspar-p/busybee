package ingest

import (
	"log"

	"github.com/apognu/gocal"
	"github.com/bwmarrin/discordgo"
	dbLib "github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/entities"
)

func IngestNewData(message *discordgo.MessageCreate, events []gocal.Event) {
	// Create a user if they do not already exist - overwrites BusyTimes
	user := GetOrCreateUser(message.Author.ID, message.Author.Username, message.GuildID)

	OverwriteUserEvents(user, events)
}

func GetOrCreateUser(userId, userName, guildId string) *entities.User {
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
		dbLib.DatabaseInstance.AddUser(user)

		return user
	}
}

func OverwriteUserEvents(user *entities.User, events []gocal.Event) {
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
	dbLib.DatabaseInstance.OverwriteUserBusyTimes(user, user.BusyTimes)
}
