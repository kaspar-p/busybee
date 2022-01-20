package update

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
)

func CheckIfUserBusy(discord *discordgo.Session, user *entities.User, guildID string) {
	// Remove role and busy status
	user.CurrentlyBusy.BusyWith = "";
	user.CurrentlyBusy.IsBusy = false;
	err := discord.GuildMemberRoleRemove(guildID, user.ID, ServerRoleIDMap[guildID]);
	if err != nil {
		fmt.Println("Error removing role from user", user.ID, ". Error: ", err);
	}

	// Add role back if necessary
	for _, busyTime := range user.BusyTimes {
		now := time.Now();
		
		// Check if now is within the bounds of the event
		if busyTime.Start.Before(now) && busyTime.End.After(now) {
			fmt.Printf("\t\tUser %s has something going on right now: %s\n", user.Name, busyTime.Title);
			user.CurrentlyBusy.IsBusy = true;
			user.CurrentlyBusy.BusyWith = busyTime.Title;

			err := discord.GuildMemberRoleAdd(guildID, user.ID, ServerRoleIDMap[guildID]);
			if err != nil {
				fmt.Println("Error adding role to user", user.ID, ", and title:", user.CurrentlyBusy.BusyWith, ". Error:", err);
			}
			break;
		}
	}
}

func UpdateSingleServer(discord *discordgo.Session, guildID string) {
	KeepRolesUpdated(discord, guildID);

	// For each user with a ID in the `users` map, change their role for the current time
	index := 1
	for _, user := range entities.Users[guildID] {
		if user.BelongsTo == guildID {
			fmt.Printf("\t%d. User %s has %d busy times!\n", index, user.Name, len(user.BusyTimes));

			// Check if the user should change busy status
			CheckIfUserBusy(discord, user, guildID);
		} else {
			fmt.Printf("This shouldn't happen! User %s has guildID %s but should have %s.\n", user.Name, user.BelongsTo, guildID);
		}
		index++;
	}
}

func UpdateAllServers(discord *discordgo.Session) {
	fmt.Println("---------------------------- UPDATING ALL SERVERS ----------------------------");
	for guildID := range ServerRoleIDMap {
		fmt.Printf("Updating roles for guild with id %s and users: %d\n", guildID, len(entities.Users[guildID]));
		UpdateSingleServer(discord, guildID);
	}
	fmt.Println("------------------------------------------------------------------------------");
}