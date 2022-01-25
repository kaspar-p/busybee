package update

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/utils"
)

func ReassignRoles(discord *discordgo.Session, guildId string, oldRoleId string, newRoleId string) {
	members, err := discord.GuildMembers(guildId, "", 1000);
	if err != nil {
		fmt.Println("Error getting members while reassigning roles in guild", guildId, "from role", oldRoleId, "to role", newRoleId);
		fmt.Println("Error: ", err);
		return;
	}

	for _, member := range members {
		hasOldBusyRole, _ := utils.StringInSlice(member.Roles, oldRoleId);
		if hasOldBusyRole {
			discord.GuildMemberRoleRemove(guildId, member.User.ID, oldRoleId);
			discord.GuildMemberRoleAdd(guildId, member.User.ID, newRoleId);
		}
	}
}

func DeleteRoleFromDiscord(discord *discordgo.Session, guildId string, roleId string) {
	err := discord.GuildRoleDelete(guildId, roleId);
	if err != nil {
		fmt.Println("Error removing role", roleId, "from guild", guildId, "during cleanup process.");
	}
}

func DeleteGuildRolePairFromDatabase(guildId string, roleId string) {
	err := database.DatabaseInstance.RemoveGuildRolePairByGuildAndRole(guildId, roleId);
	if err != nil {
		fmt.Println("Error deleting guild role pair from database: ", err);
		return;
	}
}

func ChangeGuildRoleMapEntry(guildId string, newRoleId string) {
	GuildRoleMap[guildId] = newRoleId
}

func CreateRoleInDiscord(discord *discordgo.Session, guildId string) string {
	fmt.Println("Creating busy role in guild: ", guildId)
	// There was no "busy" role - create one
	newRole, err := discord.GuildRoleCreate(guildId)
	if err != nil {
		fmt.Println("Error creating role: ", err)
	}

	return newRole.ID
}

func CreateRole(discord *discordgo.Session, guildId string) string {
	// Add the role to discord
	newRoleId := CreateRoleInDiscord(discord, guildId)

	// Add the role to the database
	database.DatabaseInstance.AddGuildRolePair(guildId, newRoleId)

	// Set global busy role ID to be the new ID
	ChangeGuildRoleMapEntry(guildId, newRoleId);

	_, err := discord.GuildRoleEdit(guildId, newRoleId, "busy :)", 12847710, false, 0, true)
	if err != nil {
		fmt.Println("Error editing role to have the correct properties. Error:", err)
	}

	return newRoleId
}
