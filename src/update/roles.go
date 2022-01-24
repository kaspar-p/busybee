package update

import (
	"fmt"

	"github.com/kaspar-p/bee/src/constants"
	"github.com/kaspar-p/bee/src/database"

	"github.com/bwmarrin/discordgo"
)

var GuildRoleMap map[string]string;

func InitializeGuildRoleMap() {
	GuildRoleMap = make(map[string]string);
}

func KeepRolesUpdated(discord *discordgo.Session, guildId string) {
	roles, err := discord.GuildRoles(guildId);
	if err != nil {
		fmt.Println("Error getting roles: ", err);
		return;
	}

	found := false
	for _, role := range roles {
		if role.ID == GuildRoleMap[guildId] {
			found = true;
		}
	}

	if !found {
		fmt.Println("Found no busy role. Creating one.");
		CreateRoleInGuild(discord, guildId);
	}
}

func CreateRoleInGuild(discord *discordgo.Session, guildId string) {
	fmt.Println("Creating busy role in guild: ", guildId);
	// There was no "busy" role - create one
	newRole, err := discord.GuildRoleCreate(guildId);
	if err != nil {
		fmt.Println("Error creating role: ", err);
	}

	// Save that data into the database
	database.DatabaseInstance.AddGuildRolePair(guildId, newRole.ID);

	// Set global busy role ID to be the new ID
	GuildRoleMap[guildId] = newRole.ID;

	_, err = discord.GuildRoleEdit(guildId, newRole.ID, constants.BusyRoleName, 12847710, false, 0, true);
	if err != nil {
		fmt.Println("Error editing role to have the correct properties. Error:", err);
	}
}
