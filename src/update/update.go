package update

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
)

func CheckIfUserBusy(discord *discordgo.Session, user *entities.User, guildId string) {
	// Remove role and busy status
	user.CurrentlyBusy.BusyWith = ""
	user.CurrentlyBusy.IsBusy = false

	err := discord.GuildMemberRoleRemove(guildId, user.Id, GuildRoleMap[guildId])
	if err != nil {
		log.Panic("Error removing role from user", user.Id, ". Error: ", err)
	}

	// Add role back if necessary
	for _, busyTime := range user.BusyTimes {
		now := time.Now()

		// Check if now is within the bounds of the event
		if busyTime.Start.Before(now) && busyTime.End.After(now) {
			fmt.Printf("\t\tUser %s has something going on right now: %s\n", user.Name, busyTime.Title)
			user.CurrentlyBusy.IsBusy = true
			user.CurrentlyBusy.BusyWith = busyTime.Title

			err := discord.GuildMemberRoleAdd(guildId, user.Id, GuildRoleMap[guildId])
			if err != nil {
				log.Panic("Error adding role to user", user.Id, ", and title:", user.CurrentlyBusy.BusyWith, ". Error:", err)
			}

			break
		}
	}
}

func UpdateSingleGuild(discord *discordgo.Session, guildId string) {
	fmt.Printf("Updating roles for guild %s and users: %d\n", guildId, len(entities.Users[guildId]))

	// For each user with a ID in the `users` map, change their role for the current time

	for _, user := range entities.Users[guildId] {
		if user.BelongsTo == guildId {
			fmt.Printf("\tUser %s has %d busy times!\n", user.Name, len(user.BusyTimes))

			// Check if the user should change busy status
			CheckIfUserBusy(discord, user, guildId)
		} else {
			fmt.Printf("This shouldn't happen! User %s has guildId %s but should have %s.\n", user.Name, user.BelongsTo, guildId)
		}
	}
}

func ValidateAllGuilds(discord *discordgo.Session) {
	fmt.Println("--------------------------- VALIDATING ALL GUILDS ----------------------------")

	var waitGroup sync.WaitGroup

	for guildId := range GuildRoleMap {
		waitGroup.Add(1)

		go func(discord *discordgo.Session, guildId string) {
			// After this function is finished, mark as done
			defer waitGroup.Done()

			// Run the process
			fmt.Printf("Validating data for guild %s and users: %d\n", guildId, len(entities.Users[guildId]))
			RunRoleValidityCheck(discord, guildId)
		}(discord, guildId)
	}
	// Wait for all validity checks to run
	waitGroup.Wait()
	fmt.Println("------------------------------------------------------------------------------")
}

func UpdateAllGuilds(discord *discordgo.Session) {
	// Boolean on whether or not to validate data.
	// Keep false until there is a mechanism for detecting disconnects from discord
	performValidation := false
	if performValidation {
		ValidateAllGuilds(discord);
	}

	fmt.Println("---------------------------- UPDATING ALL GUILDS -----------------------------")

	for guildId := range GuildRoleMap {
		UpdateSingleGuild(discord, guildId)
	}

	fmt.Println("------------------------------------------------------------------------------")
}
