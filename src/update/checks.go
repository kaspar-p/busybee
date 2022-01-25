package update

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/utils"
)

func IsRoleIdInGuild(discord *discordgo.Session, guildId string, roleId string) bool {
	roles, err := discord.GuildRoles(guildId); 
	if err != nil {
		fmt.Println("Error while checking if role ID is in guild: ", err);
		return false;
	}

	roleIds := make([]string, len(roles));
	for index, role := range roles {
		roleIds[index] = role.ID;
	}

	found, _ := utils.StringInSlice(roleIds, roleId);
	return found;
}
