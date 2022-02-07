package update

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/busybee/src/entities"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/kaspar-p/busybee/src/utils"
)

func CheckIfUserBusy(database *persist.DatabaseType, discord *discordgo.Session, user *entities.User, guildId string) {
	roleId := database.GetRoleIdForGuild(guildId)

	discordUser, err := discord.GuildMember(guildId, user.Id)
	if err != nil {
		log.Printf("Error getting discord user %s with ID %s. Error: %v\n", user.Name, user.Id, err)

		return
	}

	// Statuses of the user before the current check
	alreadyHasRole, _ := utils.StringInSlice(discordUser.Roles, roleId)

	// Add role back if necessary
	userBusyTimes := database.GetBusyTimesForUser(user.Id)

	var busyWith string

	for _, busyTime := range userBusyTimes {
		now := time.Now()

		// Check if now is within the bounds of the event
		if busyTime.Start.Before(now) && busyTime.End.After(now) {
			log.Printf("\t\tUser %s has something going on right now: %s\n", user.Name, busyTime.Title)
			busyWith = busyTime.Title

			// Only add the role for the user if they don't already have it
			if !alreadyHasRole {
				err := discord.GuildMemberRoleAdd(guildId, user.Id, roleId)
				if err != nil {
					log.Panic("Error adding role to user", user.Id, ", and title:", busyTime.Title, ". Error:", err)
				}
			}

			break
		}
	}

	// Remove the role - they are no longer busy
	if busyWith == "" && alreadyHasRole {
		err := discord.GuildMemberRoleRemove(guildId, user.Id, roleId)
		if err != nil {
			log.Printf("Error removing role %s from user %s with ID %s.\n", roleId, user.Name, user.Id)

			return
		}
	}
}

func UpdateSingleGuild(database *persist.DatabaseType, discord *discordgo.Session, guildId string) {
	users := database.GetUsersInGuild(guildId)
	log.Printf("Updating roles for guild %s and users %d\n", guildId, len(users))

	// For each user with a ID in the `users` map, change their role for the current time
	for _, user := range users {
		if user.BelongsTo == guildId {
			log.Printf("\tChecking user %s", user.Name)

			// Check if the user should change busy status
			CheckIfUserBusy(database, discord, user, guildId)
		} else {
			log.Printf("This shouldn't happen! User %s has guildId %s but should have %s.\n", user.Name, user.BelongsTo, guildId)
		}
	}
}

func UpdateAllGuilds(database *persist.DatabaseType, discord *discordgo.Session) {
	log.Println("---------------------------- UPDATING ALL GUILDS -----------------------------")

	guildIds := database.GetAllGuildRolePairs()
	for guildId := range guildIds {
		UpdateSingleGuild(database, discord, guildId)
	}

	log.Println("------------------------------------------------------------------------------")
}
