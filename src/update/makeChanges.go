package update

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/utils"
	"github.com/pkg/errors"
)

func ReassignRoles(discord *discordgo.Session, guildId, oldRoleId, newRoleId string) {
	defaultLimit := 1000

	members, err := discord.GuildMembers(guildId, "", defaultLimit)
	if err != nil {
		log.Panic("Error getting members while reassigning roles in guild", guildId, "from role", oldRoleId, "to role", newRoleId)
		log.Panic("Error: ", err)

		panic(&database.RemoveGuildRolePairError{})
	}

	for _, member := range members {
		hasOldBusyRole, _ := utils.StringInSlice(member.Roles, oldRoleId)
		if hasOldBusyRole {
			err = discord.GuildMemberRoleRemove(guildId, member.User.ID, oldRoleId)
			if err != nil {
				log.Panic("Error removing role in guild", guildId, "from role", oldRoleId,
					"to role", newRoleId, "to user", member.User.ID)
				panic(errors.Wrap(err, "Error removing role in guild "+guildId+" from role "+
					oldRoleId+" to role "+newRoleId+"to user"+member.User.ID))
			}

			err = discord.GuildMemberRoleAdd(guildId, member.User.ID, newRoleId)
			if err != nil {
				log.Panic("Error adding role in guild", guildId, "to role", newRoleId, "to user", member.User.ID)
				panic(errors.Wrap(err, "Error adding role in guild "+guildId+
					" to role "+newRoleId+"to user"+member.User.ID))
			}
		}
	}
}

func DeleteRoleFromDiscord(discord *discordgo.Session, guildId, roleId string) {
	err := discord.GuildRoleDelete(guildId, roleId)
	if err != nil {
		log.Panic("Error removing role", roleId, "from guild", guildId, "during cleanup process.")
		panic(&database.RemoveGuildRolePairError{})
	}
}

func DeleteGuildRolePairFromDatabase(guildId, roleId string) {
	err := database.DatabaseInstance.RemoveGuildRolePairByGuildAndRole(guildId, roleId)
	if err != nil {
		log.Panic("Error deleting guild role pair from database: ", err)
		panic(&database.RemoveGuildRolePairError{})
	}
}

func ChangeGuildRoleMapEntry(guildId, newRoleId string) {
	GuildRoleMap[guildId] = newRoleId
}

func CreateRoleInDiscord(discord *discordgo.Session, guildId string) string {
	// There was no "busy" role - create one
	fmt.Println("Creating busy role in guild: ", guildId)

	newRole, err := discord.GuildRoleCreate(guildId)
	if err != nil {
		log.Panic("Error creating role: ", err)
		panic(errors.Wrap(err, "Error creating role"))
	}

	return newRole.ID
}

func CreateRole(discord *discordgo.Session, guildId string) string {
	// Add the role to discord
	newRoleId := CreateRoleInDiscord(discord, guildId)

	// Add the role to the database
	database.DatabaseInstance.AddGuildRolePair(guildId, newRoleId)

	// Set global busy role ID to be the new ID
	ChangeGuildRoleMapEntry(guildId, newRoleId)

	yellowColor := 12847710

	_, err := discord.GuildRoleEdit(guildId, newRoleId, "busy :)", yellowColor, false, 0, true)
	if err != nil {
		log.Panic("Error editing role to have the correct properties. Error:", err)
		panic(errors.Wrap(err, "Error editing role to have the correct properties."))
	}

	return newRoleId
}