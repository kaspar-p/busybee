package update

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/persist"
	"github.com/kaspar-p/bee/src/utils"
)

func CheckIfUserBusy(discord *discordgo.Session, user *entities.User, guildId string) {
	discordUser, err := discord.GuildMember(guildId, user.Id)
	if err != nil {
		log.Printf("Error getting discord user %s with ID %s. Error: %v\n", user.Name, user.Id, err)

		return
	}

	// Statuses of the user before the current check
	wasBusyBefore := user.CurrentlyBusy.IsBusy
	alreadyHasRole, _ := utils.StringInSlice(discordUser.Roles, GuildRoleMap[guildId])

	// Remove role and busy status
	user.CurrentlyBusy.BusyWith = ""
	user.CurrentlyBusy.IsBusy = false

	// Add role back if necessary
	for _, busyTime := range user.BusyTimes {
		now := time.Now()

		// Check if now is within the bounds of the event
		if busyTime.Start.Before(now) && busyTime.End.After(now) {
			log.Printf("\t\tUser %s has something going on right now: %s\n", user.Name, busyTime.Title)
			user.CurrentlyBusy.IsBusy = true
			user.CurrentlyBusy.BusyWith = busyTime.Title

			// Only add the role for the user if they don't already have it
			if !alreadyHasRole {
				err := discord.GuildMemberRoleAdd(guildId, user.Id, GuildRoleMap[guildId])
				if err != nil {
					log.Panic("Error adding role to user", user.Id, ", and title:", user.CurrentlyBusy.BusyWith, ". Error:", err)
				}
			}

			break
		}
	}

	// Remove the role - they are no longer busy
	if !wasBusyBefore && alreadyHasRole {
		err := discord.GuildMemberRoleRemove(guildId, user.Id, GuildRoleMap[guildId])
		if err != nil {
			log.Printf("Error removing role %s from user %s with ID %s.\n", GuildRoleMap[guildId], user.Name, user.Id)

			return
		}
	}
}

func UpdateSingleGuild(discord *discordgo.Session, guildId string) {
	log.Printf("Updating roles for guild %s and users: %d\n", guildId, len(entities.Users[guildId]))

	// For each user with a ID in the `users` map, change their role for the current time

	for _, user := range entities.Users[guildId] {
		if user.BelongsTo == guildId {
			log.Printf("\tUser %s has %d busy times!\n", user.Name, len(user.BusyTimes))

			// Check if the user should change busy status
			CheckIfUserBusy(discord, user, guildId)
		} else {
			log.Printf("This shouldn't happen! User %s has guildId %s but should have %s.\n", user.Name, user.BelongsTo, guildId)
		}
	}
}

func ValidateAllGuilds(database *persist.DatabaseType, discord *discordgo.Session) {
	log.Println("--------------------------- VALIDATING ALL GUILDS ----------------------------")

	var waitGroup sync.WaitGroup

	for guildId := range GuildRoleMap {
		waitGroup.Add(1)

		go func(discord *discordgo.Session, guildId string) {
			// After this function is finished, mark as done
			defer waitGroup.Done()

			// Run the process
			log.Printf("Validating data for guild %s and users: %d\n", guildId, len(entities.Users[guildId]))
			RunRoleValidityCheck(database, discord, guildId)
		}(discord, guildId)
	}
	// Wait for all validity checks to run
	waitGroup.Wait()
	log.Println("------------------------------------------------------------------------------")
}

func UpdateAllGuilds(database *persist.DatabaseType, discord *discordgo.Session) {
	// Boolean on whether or not to validate data.
	// Keep false until there is a mechanism for detecting disconnects from discord
	performValidation := false
	if performValidation {
		ValidateAllGuilds(database, discord)
	}

	log.Println("---------------------------- UPDATING ALL GUILDS -----------------------------")

	for guildId := range GuildRoleMap {
		UpdateSingleGuild(discord, guildId)
	}

	log.Println("------------------------------------------------------------------------------")
}
