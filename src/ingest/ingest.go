package ingest

import (
	"fmt"

	dbLib "github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/entities"

	"github.com/apognu/gocal"
	"github.com/bwmarrin/discordgo"
)

func FillMapsWithDatabaseData() {
	// Get data and fill the `users` map
	users := dbLib.DatabaseInstance.GetUsers();
	for _, user := range users {
		entities.Users[user.BelongsTo][user.ID] = user;
	}

	// Get data and fill the busyTimes of each user in the `users` map
	busyTimesArray := dbLib.DatabaseInstance.GetBusyTimes();
	for _, busyTime := range busyTimesArray {
		user := entities.Users[busyTime.BelongsTo][busyTime.OwnerID];
		user.BusyTimes = append(user.BusyTimes, busyTime);
	}

	for _, user := range users {
		user.SortBusyTimes();
	}

	fmt.Println("Got data: \n\tUsers:", len(users), "\n\tEvents:", len(busyTimesArray));
}

func IngestNewData(message *discordgo.MessageCreate, events []gocal.Event) {
	// Create a user if they do not already exist - overwrites BusyTimes
	user := GetOrCreateUser(message.Author.ID, message.Author.Username, message.GuildID);
	
	OverwriteUserEvents(user, events);
}

func GetOrCreateUser(ID string, userName string, guildID string) *entities.User {
	if user, ok := entities.Users[guildID][ID]; ok {
		fmt.Println("User found with ID: ", ID);
		return user;
	} else {
		fmt.Println("User created with ID:", ID);
		// Create the new user
		user := entities.CreateUser(userName, ID, guildID);

		// Add the new user to the `users` map
		entities.Users[user.BelongsTo][user.ID] = user;

		// Add the new user to the database
		dbLib.DatabaseInstance.AddUser(user);

		return user;
	}
}

func OverwriteUserEvents(user *entities.User, events []gocal.Event) {
	// Overwrite the busyTimes in memory
	user.BusyTimes = make([]*entities.BusyTime, 0);
	for _, event := range events {
		title := ParseEventTitle(event.Summary);
		busyTime := entities.CreateBusyTime(user.ID, user.BelongsTo, title, *event.Start, *event.End);
		user.BusyTimes = append(user.BusyTimes, &busyTime);
	}

	user.SortBusyTimes();

	// Overwrite the busyTimes in the database
	dbLib.DatabaseInstance.OverwriteUserBusyTimes(user, user.BusyTimes);
}