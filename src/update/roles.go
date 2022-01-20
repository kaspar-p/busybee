package update

import (
	"fmt"

	"github.com/kaspar-p/bee/src/constants"

	"github.com/bwmarrin/discordgo"
)

var ServerRoleIDMap map[string]string;

func InitializeServerRoleIDMap() {
	ServerRoleIDMap = make(map[string]string);
}

func KeepRolesUpdated(discord *discordgo.Session, guildID string) {
	roles, err := discord.GuildRoles(guildID);
	if err != nil {
		fmt.Println("Error getting roles: ", err);
		return;
	}

	foundBusyRoleID := "";
	for _, role := range roles {
		if role.Name == constants.BusyRoleName {
			foundBusyRoleID = role.ID;
		}
	}

	if foundBusyRoleID != "" {
		ServerRoleIDMap[guildID] = foundBusyRoleID;
	} else {
		fmt.Println("Found no busy role. Creating one.");
		
		// There was no "busy" role - create one
		newRole, err := discord.GuildRoleCreate(guildID);
		if err != nil {
			fmt.Println("Error creating role: ", err);
		}

		// Set global busy role ID to be the new ID
		ServerRoleIDMap[guildID] = newRole.ID;

		_, err = discord.GuildRoleEdit(guildID, newRole.ID, constants.BusyRoleName, 12847710, false, 0, true);
		if err != nil {
			fmt.Println("Error editing role to have the correct properties. Error:", err);
		}
	}
}
